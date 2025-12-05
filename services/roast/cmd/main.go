package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/repository"
	"github.com/nerufuyo/nerubot/internal/usecase/roast"
)

const (
	ServiceName    = "roast-service"
	ServiceVersion = "1.0.0"
	DefaultPort    = "8083"
)

type RoastServer struct {
	config  *config.Config
	logger  *logger.Logger
	service *roast.RoastService
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

	log.Info("=== NeruBot Roast Service v1.0.0 ===")
	log.Info("Starting Roast service...")

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

	// Initialize repository
	_ = repository.NewRoastRepository()

	// Initialize roast service
	roastService := roast.NewRoastService()

	// Create server
	server := &RoastServer{
		config:  cfg,
		logger:  log,
		service: roastService,
	}

	// Start HTTP server for health checks
	go server.startHealthServer(port)

	log.Info("Roast service is running", "port", port)
	log.Info("Note: gRPC server implementation pending")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Info("Shutting down Roast service...")
}

// startHealthServer starts HTTP health check server
func (s *RoastServer) startHealthServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"roast-service","version":"1.0.0"}`))
	})

	s.logger.Info("Starting health check server", "port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		s.logger.Error("Health server error", "error", err)
	}
}
