package main

import (
"context"
"errors"
"log/slog"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/joho/godotenv"

"hello_world/internal/config"
database "hello_world/internal/db"
"hello_world/internal/server"
)

func main() {
// Load .env if present (development convenience; no-op in production).
if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
slog.Warn("could not load .env file – this is expected in production", "error", err)
}

cfg, err := config.New()
if err != nil {
slog.Error("invalid configuration", "error", err)
os.Exit(1)
}

dbService, err := database.New(cfg)
if err != nil {
slog.Error("failed to connect to database", "error", err)
os.Exit(1)
}

httpServer := server.New(cfg, dbService)

// Run graceful shutdown in a separate goroutine.
idleConnsClosed := make(chan struct{})
go func() {
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

slog.Info("shutdown signal received – draining connections...")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := httpServer.Shutdown(ctx); err != nil {
slog.Error("server shutdown error", "error", err)
}

if err := dbService.Close(); err != nil {
slog.Error("database close error", "error", err)
}

close(idleConnsClosed)
}()

slog.Info("server starting", "addr", httpServer.Addr)
if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
slog.Error("http server error", "error", err)
os.Exit(1)
}

<-idleConnsClosed
slog.Info("server stopped gracefully")
}
