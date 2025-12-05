package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/news"
)

const (
	ServiceName    = "news-service"
	ServiceVersion = "1.0.0"
	DefaultPort    = "8085"
)

type NewsServer struct {
	config  *config.Config
	logger  *logger.Logger
	service *news.NewsService
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

	log.Info("=== NeruBot News Service v1.0.0 ===")
	log.Info("Starting News service...")

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

	// Initialize news service
	newsService := news.NewNewsService()
	if newsService == nil {
		log.Error("Failed to initialize news service")
		os.Exit(1)
	}
	log.Info("News service initialized successfully")

	// Create server
	server := &NewsServer{
		config:  cfg,
		logger:  log,
		service: newsService,
	}

	// Setup HTTP health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"news-service","version":"1.0.0"}`))
	})

	// Start HTTP server in goroutine
	go func() {
		log.Info("Starting HTTP health check server", "port", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Error("HTTP server failed", "error", err)
		}
	}()

	log.Info("News service started successfully", "port", port)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down News service...")
	server.shutdown()
	log.Info("News service stopped")
}

func (s *NewsServer) shutdown() {
	s.logger.Info("Performing graceful shutdown...")
	// Add cleanup logic here if needed
}
