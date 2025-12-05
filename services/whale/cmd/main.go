package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	whalepb "github.com/nerufuyo/nerubot/api/proto/whale"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/whale"
	"google.golang.org/grpc"
)

const (
	ServiceName     = "whale-service"
	ServiceVersion  = "1.0.0"
	DefaultHTTPPort = "8086"
	DefaultGRPCPort = "50056"
)

type WhaleServer struct {
	whalepb.UnimplementedWhaleServiceServer
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

	// Get ports from environment
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = DefaultHTTPPort
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = DefaultGRPCPort
	}

	// Initialize whale service
	whaleService := whale.NewWhaleService(cfg.Crypto.WhaleAlertAPIKey)
	if whaleService == nil {
		log.Error("Failed to initialize whale service")
		os.Exit(1)
	}

	// Create server
	server := &WhaleServer{
		config:  cfg,
		logger:  log,
		service: whaleService,
	}

	// Start HTTP health server
	go server.startHealthServer(httpPort)

	// Start gRPC server
	go func() {
		if err := server.startGRPCServer(grpcPort); err != nil {
			log.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("Whale service running", "http_port", httpPort, "grpc_port", grpcPort)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down Whale service...")
}

func (s *WhaleServer) startHealthServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"whale-service","version":"1.0.0"}`))
	})

	s.logger.Info("Starting health check server", "port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		s.logger.Error("Health server error", "error", err)
	}
}

func (s *WhaleServer) startGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	whalepb.RegisterWhaleServiceServer(grpcServer, s)

	s.logger.Info("gRPC server listening", "port", port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *WhaleServer) GetTransactions(ctx context.Context, req *whalepb.GetTransactionsRequest) (*whalepb.GetTransactionsResponse, error) {
	s.logger.Info("Get transactions", "blockchain", req.Blockchain, "limit", req.Limit)

	return &whalepb.GetTransactionsResponse{
		Transactions: []*whalepb.Transaction{},
	}, nil
}

func (s *WhaleServer) GetAlertSettings(ctx context.Context, req *whalepb.GetAlertSettingsRequest) (*whalepb.GetAlertSettingsResponse, error) {
	s.logger.Info("Get alert settings", "guild", req.GuildId)

	return &whalepb.GetAlertSettingsResponse{
		Settings: &whalepb.AlertSettings{
			GuildId:      req.GuildId,
			ChannelId:    "",
			MinAmountUsd: 1000000,
			Enabled:      true,
			Blockchains:  []string{"ethereum", "bitcoin"},
		},
	}, nil
}

func (s *WhaleServer) UpdateAlertSettings(ctx context.Context, req *whalepb.UpdateAlertSettingsRequest) (*whalepb.UpdateAlertSettingsResponse, error) {
	s.logger.Info("Update alert settings", "guild", req.GuildId)

	return &whalepb.UpdateAlertSettingsResponse{
		Success: true,
		Message: "Settings updated",
	}, nil
}

func (s *WhaleServer) HealthCheck(ctx context.Context, req *whalepb.HealthCheckRequest) (*whalepb.HealthCheckResponse, error) {
	return &whalepb.HealthCheckResponse{
		Healthy: true,
		Service: ServiceName,
		Version: ServiceVersion,
	}, nil
}
