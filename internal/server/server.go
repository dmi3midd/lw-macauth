package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"lw-macauth/internal/config"
	"lw-macauth/internal/services"
)

type Server struct {
	cfg          *config.Config
	tokenService *services.TokenService
}

func NewServer(cfg *config.Config) *http.Server {
	tokenService := services.NewTokenService(cfg.Keys)
	s := &Server{
		cfg:          cfg,
		tokenService: &tokenService,
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
