package server

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/dmi3midd/lw-macauth/internal/config"
	"github.com/dmi3midd/lw-macauth/internal/service"
)

type Server struct {
	cfg          *config.Config
	tokenService service.TokenService
}

func NewServer(cfg *config.Config) *http.Server {
	tokenService := service.NewTokenService(cfg.Keys)
	s := &Server{
		cfg:          cfg,
		tokenService: tokenService,
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
