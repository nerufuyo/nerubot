package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/whale"
)

const (
	ServiceName    = "whale-service"
	ServiceVersion = "1.0.0"
	DefaultPort    = "8086"
)

type WhaleServer struct {
	config  *config.Config
	logger  *logger.Logger
	service *whale.WhaleService
}

func main() {
	// Initialize logger
	logCfg := logger.DefaultConfig()
	logCfg.Level = logger.LevelInfo
	log, err := logger.Init(logCfg)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("=== NeruBot Whale Service v1.0.0 ===")
	log.Info("Starting Whale service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Initialize whale service
	whaleService := whale.NewWhaleService(cfg.Crypto.WhaleAlertAPIKey)
	if whaleService == nil {
		log.Error("Failed to initialize whale service")
		os.Exit(1)
	}
	log.Info("Whale service initialized successfully")

	// Create server
	server := &WhaleServer{
		config:  cfg,
		logger:  log,
		service: whaleService,
	}

	// Setup HTTP health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"whale-service","version":"1.0.0"}`))
	})

	// Start HTTP server in goroutine
	go func() {
		log.Info("Starting HTTP health check server", "port", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Error("HTTP server failed", "error", err)
		}
	}()

	log.Info("Whale service started successfully", "port", port)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down Whale service...")
	server.shutdown()
	log.Info("Whale service stopped")
}

func (s *WhaleServer) shutdown() {
	s.logger.Info("Performing graceful shutdown...")
	// Add cleanup logic here if needed
}
