package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	newspb "github.com/nerufuyo/nerubot/api/proto/news"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/usecase/news"
	"google.golang.org/grpc"
)

const (
	ServiceName     = "news-service"
	ServiceVersion  = "1.0.0"
	DefaultHTTPPort = "8085"
	DefaultGRPCPort = "50055"
)

type NewsServer struct {
	newspb.UnimplementedNewsServiceServer
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

	// Get ports from environment
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = DefaultHTTPPort
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = DefaultGRPCPort
	}

	// Initialize news service
	newsService := news.NewNewsService()
	if newsService == nil {
		log.Error("Failed to initialize news service")
		os.Exit(1)
	}

	// Create server
	server := &NewsServer{
		config:  cfg,
		logger:  log,
		service: newsService,
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

	log.Info("News service running", "http_port", httpPort, "grpc_port", grpcPort)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down News service...")
}

func (s *NewsServer) startHealthServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"news-service","version":"1.0.0"}`))
	})

	s.logger.Info("Starting health check server", "port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		s.logger.Error("Health server error", "error", err)
	}
}

func (s *NewsServer) startGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	newspb.RegisterNewsServiceServer(grpcServer, s)

	s.logger.Info("gRPC server listening", "port", port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *NewsServer) FetchNews(ctx context.Context, req *newspb.FetchNewsRequest) (*newspb.FetchNewsResponse, error) {
	s.logger.Info("Fetch news", "guild", req.GuildId, "limit", req.Limit)

	return &newspb.FetchNewsResponse{
		Articles:      []*newspb.Article{},
		TotalArticles: 0,
	}, nil
}

func (s *NewsServer) AddSource(ctx context.Context, req *newspb.AddSourceRequest) (*newspb.AddSourceResponse, error) {
	s.logger.Info("Add source", "guild", req.GuildId, "name", req.SourceName)

	return &newspb.AddSourceResponse{
		Success:  true,
		Message:  "Source added",
		SourceId: 1,
	}, nil
}

func (s *NewsServer) RemoveSource(ctx context.Context, req *newspb.RemoveSourceRequest) (*newspb.RemoveSourceResponse, error) {
	s.logger.Info("Remove source", "guild", req.GuildId, "source_id", req.SourceId)

	return &newspb.RemoveSourceResponse{
		Success: true,
		Message: "Source removed",
	}, nil
}

func (s *NewsServer) GetSources(ctx context.Context, req *newspb.GetSourcesRequest) (*newspb.GetSourcesResponse, error) {
	s.logger.Info("Get sources", "guild", req.GuildId)

	return &newspb.GetSourcesResponse{
		Sources: []*newspb.Source{},
	}, nil
}

func (s *NewsServer) PublishNews(ctx context.Context, req *newspb.PublishNewsRequest) (*newspb.PublishNewsResponse, error) {
	s.logger.Info("Publish news", "guild", req.GuildId, "article", req.ArticleId)

	return &newspb.PublishNewsResponse{
		Success:   true,
		Message:   "News published",
		MessageId: "",
	}, nil
}

func (s *NewsServer) HealthCheck(ctx context.Context, req *newspb.HealthCheckRequest) (*newspb.HealthCheckResponse, error) {
	return &newspb.HealthCheckResponse{
		Healthy: true,
		Service: ServiceName,
		Version: ServiceVersion,
	}, nil
}
