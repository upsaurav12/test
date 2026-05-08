package server

import (
	"fmt"
	"net/http"
	"time"

	"hello_world/internal/config"
	database "hello_world/internal/db"
)

// Server holds the application's dependencies needed to build HTTP routes.
type Server struct {
	cfg *config.Config
	db  database.Service
}

// New constructs and returns a fully configured *http.Server.
func New(cfg *config.Config, db database.Service) *http.Server {
	srv := &Server{cfg: cfg, db: db}

	return &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           srv.registerRoutes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
	}
}
