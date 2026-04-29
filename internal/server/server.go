package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"lw-macauth/internal/config"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *http.Server {
	s := &Server{
		cfg: cfg,
	}

	router := s.RegisterRoutes()
	return &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}
}
