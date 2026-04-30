package server

import (
	errs "lw-macauth/internal/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/v1/generate", errs.ErrorHandler(s.GenerateTokens))
	r.Get("/v1/validate/access", errs.ErrorHandler(s.ValidateAccessToken))
	r.Post("/v1/validate/refresh", errs.ErrorHandler(s.ValidateRefreshToken))
	r.Get("/v1/public", errs.ErrorHandler(s.GetPublicKey))

	return r
}
