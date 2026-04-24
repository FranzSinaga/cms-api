package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FranzSinaga/blogcms/internal/shared"
	"github.com/FranzSinaga/blogcms/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	logger := shared.InitLogger()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("Warning: .env file not found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", slog.Any("error", err))
	}

	// Connect to database
	db, err := config.NewDatabase(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", slog.Any("error", err))
	}
	defer db.Close()
	logger.Info("Database connected!")

	// Initialize dependencies
	container := NewContainer(db, cfg)
	r := setupRouter(container, cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting...", slog.String("port", cfg.App.Port), slog.String("env", cfg.App.Env))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", slog.Any("error", err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("\nShutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", slog.Any("error", err))
	}

	logger.Info("Server stopped gracefully")
}
