package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	confessionpb "github.com/nerufuyo/nerubot/api/proto/confession"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/confession"
	"google.golang.org/grpc"
)

const (
	ServiceName     = "confession-service"
	ServiceVersion  = "1.0.0"
	DefaultHTTPPort = "8082"
	DefaultGRPCPort = "50052"
)

type ConfessionServer struct {
	confessionpb.UnimplementedConfessionServiceServer
	config  *config.Config
	logger  *logger.Logger
	service *confession.ConfessionService
}

func main() {
	logCfg := logger.DefaultConfig()
	logCfg.Level = logger.LevelInfo
	log, err := logger.Init(logCfg)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("=== NeruBot Confession Service v1.0.0 ===")

	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = DefaultHTTPPort
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = DefaultGRPCPort
	}

	confessionService := confession.NewConfessionService()

	server := &ConfessionServer{
		config:  cfg,
		logger:  log,
		service: confessionService,
	}

	go server.startHealthServer(httpPort)

	go func() {
		if err := server.startGRPCServer(grpcPort); err != nil {
			log.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("Confession service running", "http_port", httpPort, "grpc_port", grpcPort)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Info("Shutting down...")
}

func (s *ConfessionServer) startHealthServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"confession-service","version":"1.0.0"}`))
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		s.logger.Error("Health server error", "error", err)
	}
}

func (s *ConfessionServer) startGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	confessionpb.RegisterConfessionServiceServer(grpcServer, s)

	s.logger.Info("gRPC server listening", "port", port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *ConfessionServer) Submit(ctx context.Context, req *confessionpb.SubmitRequest) (*confessionpb.SubmitResponse, error) {
	s.logger.Info("Submit confession", "guild", req.GuildId)

	return &confessionpb.SubmitResponse{
		Success:      true,
		Message:      "Confession submitted",
		ConfessionId: 1,
	}, nil
}

func (s *ConfessionServer) Approve(ctx context.Context, req *confessionpb.ApproveRequest) (*confessionpb.ApproveResponse, error) {
	s.logger.Info("Approve confession", "id", req.ConfessionId)

	return &confessionpb.ApproveResponse{
		Success: true,
		Message: "Approved",
	}, nil
}

func (s *ConfessionServer) Reject(ctx context.Context, req *confessionpb.RejectRequest) (*confessionpb.RejectResponse, error) {
	s.logger.Info("Reject confession", "id", req.ConfessionId)

	return &confessionpb.RejectResponse{
		Success: true,
		Message: "Rejected",
	}, nil
}

func (s *ConfessionServer) Reply(ctx context.Context, req *confessionpb.ReplyRequest) (*confessionpb.ReplyResponse, error) {
	s.logger.Info("Reply to confession", "number", req.ConfessionNumber)

	return &confessionpb.ReplyResponse{
		Success:   true,
		Message:   "Reply sent",
		MessageId: "",
	}, nil
}

func (s *ConfessionServer) GetPending(ctx context.Context, req *confessionpb.GetPendingRequest) (*confessionpb.GetPendingResponse, error) {
	s.logger.Info("Get pending", "guild", req.GuildId)

	return &confessionpb.GetPendingResponse{
		Confessions: []*confessionpb.Confession{},
	}, nil
}

func (s *ConfessionServer) GetSettings(ctx context.Context, req *confessionpb.GetSettingsRequest) (*confessionpb.GetSettingsResponse, error) {
	s.logger.Info("Get settings", "guild", req.GuildId)

	return &confessionpb.GetSettingsResponse{
		Settings: &confessionpb.Settings{
			GuildId:             req.GuildId,
			Enabled:             true,
			RequireApproval:     true,
			AllowImages:         false,
			MaxLength:           500,
		},
	}, nil
}

func (s *ConfessionServer) UpdateSettings(ctx context.Context, req *confessionpb.UpdateSettingsRequest) (*confessionpb.UpdateSettingsResponse, error) {
	s.logger.Info("Update settings", "guild", req.GuildId)

	return &confessionpb.UpdateSettingsResponse{
		Success: true,
		Message: "Settings updated",
	}, nil
}

func (s *ConfessionServer) HealthCheck(ctx context.Context, req *confessionpb.HealthCheckRequest) (*confessionpb.HealthCheckResponse, error) {
	return &confessionpb.HealthCheckResponse{
		Healthy: true,
		Service: ServiceName,
		Version: ServiceVersion,
	}, nil
}
