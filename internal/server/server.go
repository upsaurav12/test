package server

import (
	"fmt"
	"net/http"
	"time"

	database "hello_world/internal/db"
	"hello_world/internal/config"
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
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      srv.registerRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}
}
