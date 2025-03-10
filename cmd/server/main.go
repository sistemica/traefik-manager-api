// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sistemica/traefik-manager/internal/api/server"
	"github.com/sistemica/traefik-manager/internal/config"
	"github.com/sistemica/traefik-manager/internal/logger"
	"github.com/sistemica/traefik-manager/internal/store"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	err = logger.Setup(logger.Config{
		Level:    cfg.Logger.Level,
		Format:   cfg.Logger.Format,
		FilePath: cfg.Logger.FilePath,
		UseColor: cfg.Logger.UseColors,
	})
	if err != nil {
		fmt.Printf("Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize store
	dataStore, err := store.NewFileStore(cfg.Storage.FilePath)
	if err != nil {
		logger.Fatal().Err(err).Str("path", cfg.Storage.FilePath).Msg("Failed to initialize store")
	}

	// Initialize and setup server
	server := server.New(cfg, dataStore)
	server.Setup()

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Setup automatic store saving
	if cfg.Storage.SaveInterval > 0 {
		go func() {
			ticker := time.NewTicker(cfg.Storage.SaveInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := dataStore.Save(); err != nil {
						logger.Error().Err(err).Msg("Failed to auto-save store data")
					} else {
						logger.Debug().Msg("Store data auto-saved")
					}
				}
			}
		}()
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server shutdown failed")
	}

	// Save store data
	logger.Info().Msg("Saving store data")
	if err := dataStore.Save(); err != nil {
		logger.Error().Err(err).Msg("Failed to save store data")
	}

	logger.Info().Msg("Server stopped")
}
