package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
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

	// TODO: Initialize Discord session
	// TODO: Load use cases and repositories
	// TODO: Register command handlers
	// TODO: Start bot

	log.Info("Bot is running. Press CTRL+C to exit.")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("Shutting down gracefully...")
	// TODO: Cleanup resources
	// TODO: Close Discord connection
	// TODO: Save state
	log.Info("Goodbye!")
}
