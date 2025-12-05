package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/music"
)

const (
	ServiceName    = "music-service"
	ServiceVersion = "1.0.0"
	DefaultPort    = "8081"
)

type MusicServer struct {
	config  *config.Config
	logger  *logger.Logger
	service *music.MusicService
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

	log.Info("=== NeruBot Music Service v1.0.0 ===")
	log.Info("Starting Music service...")

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

	// Initialize music service (from existing usecase)
	musicService, err := music.NewMusicService()
	if err != nil {
		log.Error("Failed to initialize music service", "error", err)
		os.Exit(1)
	}

	// Create server
	server := &MusicServer{
		config:  cfg,
		logger:  log,
		service: musicService,
	}

	// Start HTTP server for health checks
	go server.startHealthServer(port)

	// TODO: Start gRPC server
	// For now, we'll just run the health server
	log.Info("Music service is running", "port", port)
	log.Info("Note: gRPC server implementation pending")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Info("Shutting down Music service...")
}

// startHealthServer starts HTTP health check server
func (s *MusicServer) startHealthServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"music-service","version":"1.0.0"}`))
	})

	s.logger.Info("Starting health check server", "port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		s.logger.Error("Health server error", "error", err)
	}
}

// TODO: Implement gRPC server
// This will require generated proto files and gRPC service implementation
func (s *MusicServer) startGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info("gRPC server listening", "port", port)

	// TODO: Create gRPC server and register service
	// grpcServer := grpc.NewServer()
	// pb.RegisterMusicServiceServer(grpcServer, s)
	// return grpcServer.Serve(lis)

	_ = lis // Suppress unused variable warning
	return nil
}
