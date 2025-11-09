package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/delivery/discord"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/repository"
	"github.com/nerufuyo/nerubot/internal/usecase/confession"
	"github.com/nerufuyo/nerubot/internal/usecase/music"
	"github.com/nerufuyo/nerubot/internal/usecase/roast"
)

func main() {
	// Initialize logger
	logCfg := logger.DefaultConfig()
	logCfg.Level = logger.LevelInfo
	log, err := logger.Init(logCfg)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("=== NeruBot v3.0.0 (Golang Edition) ===")
	log.Info("Starting Discord bot...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	log.Info("Configuration loaded successfully",
		"bot", cfg.Bot.Name,
		"version", cfg.Bot.Version,
	)

	// Log enabled features
	log.Info("Features enabled",
		"music", cfg.Features.Music,
		"chatbot", cfg.Features.Chatbot,
		"confession", cfg.Features.Confession,
		"roast", cfg.Features.Roast,
		"news", cfg.Features.News,
		"whale_alerts", cfg.Features.WhaleAlerts,
	)

	// Initialize repositories
	confessionRepo := repository.NewConfessionRepository()
	roastRepo := repository.NewRoastRepository()

	log.Info("Repositories initialized successfully")

	// Initialize use cases
	if cfg.Features.Music {
		if _, err := music.NewMusicService(); err != nil {
			log.Error("Failed to initialize music service", "error", err)
			os.Exit(1)
		}
		log.Info("Music service initialized")
	}

	if cfg.Features.Confession {
		_ = confession.NewConfessionService()
		log.Info("Confession service initialized")
	}

	if cfg.Features.Roast {
		_ = roast.NewRoastService()
		log.Info("Roast service initialized")
	}

	// Initialize Discord bot
	bot, err := discord.New(cfg)
	if err != nil {
		log.Error("Failed to create Discord bot", "error", err)
		os.Exit(1)
	}

	// Start the bot
	ctx := context.Background()
	if err := bot.Start(ctx); err != nil {
		log.Error("Failed to start Discord bot", "error", err)
		os.Exit(1)
	}

	log.Info("Bot is running. Press CTRL+C to exit.")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Shutting down gracefully...")
	
	// Cleanup resources
	if err := bot.Stop(); err != nil {
		log.Error("Error stopping bot", "error", err)
	}

	// Save repositories
	_ = confessionRepo
	_ = roastRepo

	log.Info("Goodbye!")
}
