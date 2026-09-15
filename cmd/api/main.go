package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmi3midd/lw-macauth/internal/config"
	"github.com/dmi3midd/lw-macauth/internal/logger"
	"github.com/dmi3midd/lw-macauth/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.Setup(cfg.Log.Level)

	srv := server.NewServer(cfg)

	slog.Info(
		"server is running",
		slog.String("address", cfg.HTTPServer.Address),
	)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to run server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	slog.Info("Graceful shutdown complete")
}
