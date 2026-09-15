package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/dmi3midd/lw-macauth/internal/shared/apierror"
	"github.com/rs/xid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = xid.New().String()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "panic recovered",
					slog.Any("panic", rec),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"code":500,"message":"Internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		slog.InfoContext(r.Context(), "http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /generate", apierror.ErrorHandler(s.GenerateTokens))
	mux.HandleFunc("POST /validate/access", apierror.ErrorHandler(s.ValidateAccessToken))
	mux.HandleFunc("POST /validate/refresh", apierror.ErrorHandler(s.ValidateRefreshToken))
	mux.HandleFunc("GET /public", apierror.ErrorHandler(s.GetPublicKey))

	var handler http.Handler = mux
	handler = loggingMiddleware(handler)
	handler = recovererMiddleware(handler)
	handler = requestIDMiddleware(handler)

	return handler
}
