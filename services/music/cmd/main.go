package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	musicpb "github.com/nerufuyo/nerubot/api/proto/music"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/music"
	"google.golang.org/grpc"
)

const (
	ServiceName     = "music-service"
	ServiceVersion  = "1.0.0"
	DefaultHTTPPort = "8081"
	DefaultGRPCPort = "50051"
)

type MusicServer struct {
	musicpb.UnimplementedMusicServiceServer
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
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = DefaultHTTPPort
	}
	
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = DefaultGRPCPort
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
	go server.startHealthServer(httpPort)

	// Start gRPC server
	go func() {
		if err := server.startGRPCServer(grpcPort); err != nil {
			log.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("Music service is running", "http_port", httpPort, "grpc_port", grpcPort)

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

	grpcServer := grpc.NewServer()
	musicpb.RegisterMusicServiceServer(grpcServer, s)

	s.logger.Info("gRPC server listening", "port", port)
	
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Play implements the Play RPC method
func (s *MusicServer) Play(ctx context.Context, req *musicpb.PlayRequest) (*musicpb.PlayResponse, error) {
	s.logger.Info("Play request", "guild", req.GuildId, "url", req.SongUrl)
	
	// TODO: Implement actual music playback logic
	return &musicpb.PlayResponse{
		Success: true,
		Message: fmt.Sprintf("Playing: %s", req.SongUrl),
	}, nil
}

// Pause implements the Pause RPC method
func (s *MusicServer) Pause(ctx context.Context, req *musicpb.PauseRequest) (*musicpb.PauseResponse, error) {
	s.logger.Info("Pause request", "guild", req.GuildId)
	
	return &musicpb.PauseResponse{
		Success: true,
		Message: "Playback paused",
	}, nil
}

// Resume implements the Resume RPC method
func (s *MusicServer) Resume(ctx context.Context, req *musicpb.ResumeRequest) (*musicpb.ResumeResponse, error) {
	s.logger.Info("Resume request", "guild", req.GuildId)
	
	return &musicpb.ResumeResponse{
		Success: true,
		Message: "Playback resumed",
	}, nil
}

// Skip implements the Skip RPC method
func (s *MusicServer) Skip(ctx context.Context, req *musicpb.SkipRequest) (*musicpb.SkipResponse, error) {
	s.logger.Info("Skip request", "guild", req.GuildId)
	
	return &musicpb.SkipResponse{
		Success: true,
		Message: "Track skipped",
	}, nil
}

// Stop implements the Stop RPC method
func (s *MusicServer) Stop(ctx context.Context, req *musicpb.StopRequest) (*musicpb.StopResponse, error) {
	s.logger.Info("Stop request", "guild", req.GuildId)
	
	return &musicpb.StopResponse{
		Success: true,
		Message: "Playback stopped",
	}, nil
}

// Queue implements the Queue RPC method
func (s *MusicServer) Queue(ctx context.Context, req *musicpb.QueueRequest) (*musicpb.QueueResponse, error) {
	s.logger.Info("Queue request", "guild", req.GuildId, "url", req.SongUrl)
	
	return &musicpb.QueueResponse{
		Success: true,
		Message: "Added to queue",
	}, nil
}

// NowPlaying implements the NowPlaying RPC method
func (s *MusicServer) NowPlaying(ctx context.Context, req *musicpb.NowPlayingRequest) (*musicpb.NowPlayingResponse, error) {
	s.logger.Info("NowPlaying request", "guild", req.GuildId)
	
	return &musicpb.NowPlayingResponse{
		IsPlaying: false,
	}, nil
}

// GetQueue implements the GetQueue RPC method
func (s *MusicServer) GetQueue(ctx context.Context, req *musicpb.GetQueueRequest) (*musicpb.GetQueueResponse, error) {
	s.logger.Info("GetQueue request", "guild", req.GuildId)
	
	return &musicpb.GetQueueResponse{
		Songs: []*musicpb.Song{},
	}, nil
}

// SetLoop implements the SetLoop RPC method
func (s *MusicServer) SetLoop(ctx context.Context, req *musicpb.SetLoopRequest) (*musicpb.SetLoopResponse, error) {
	s.logger.Info("SetLoop request", "guild", req.GuildId, "loop_mode", req.LoopMode)
	
	return &musicpb.SetLoopResponse{
		Success: true,
		Message: "Loop setting updated",
	}, nil
}

// SetVolume implements the SetVolume RPC method
func (s *MusicServer) SetVolume(ctx context.Context, req *musicpb.SetVolumeRequest) (*musicpb.SetVolumeResponse, error) {
	s.logger.Info("SetVolume request", "guild", req.GuildId, "volume", req.Volume)
	
	return &musicpb.SetVolumeResponse{
		Success: true,
		Message: fmt.Sprintf("Volume set to %d", req.Volume),
	}, nil
}

// HealthCheck implements the HealthCheck RPC method
func (s *MusicServer) HealthCheck(ctx context.Context, req *musicpb.HealthCheckRequest) (*musicpb.HealthCheckResponse, error) {
	return &musicpb.HealthCheckResponse{
		Healthy: true,
		Service: ServiceName,
		Version: ServiceVersion,
	}, nil
}
