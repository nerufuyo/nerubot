package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/chatbot"
)

const (
	ServiceName    = "chatbot-service"
	ServiceVersion = "1.0.0"
	DefaultPort    = "8084"
)

type ChatbotServer struct {
	config  *config.Config
	logger  *logger.Logger
	service *chatbot.ChatbotService
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

	log.Info("=== NeruBot Chatbot Service v1.0.0 ===")
	log.Info("Starting Chatbot service...")

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

	// Initialize chatbot service
	chatbotService := chatbot.NewChatbotService(
		cfg.AI.DeepSeekKey,
	)
	if chatbotService == nil {
		log.Error("Failed to initialize chatbot service")
		os.Exit(1)
	}
	log.Info("Chatbot service initialized successfully")

	// Create server
	server := &ChatbotServer{
		config:  cfg,
		logger:  log,
		service: chatbotService,
	}

	// Setup HTTP health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"chatbot-service","version":"1.0.0"}`))
	})

	// Start HTTP server in goroutine
	go func() {
		log.Info("Starting HTTP health check server", "port", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Error("HTTP server failed", "error", err)
		}
	}()

	log.Info("Chatbot service started successfully", "port", port)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down Chatbot service...")
	server.shutdown()
	log.Info("Chatbot service stopped")
}

func (s *ChatbotServer) shutdown() {
	s.logger.Info("Performing graceful shutdown...")
	// Add cleanup logic here if needed
}
