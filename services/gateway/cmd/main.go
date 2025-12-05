package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

const (
	ServiceName    = "api-gateway"
	ServiceVersion = "1.0.0"
)

type Gateway struct {
	config    *config.Config
	logger    *logger.Logger
	discord   *discordgo.Session
	musicURL  string
	confURL   string
	roastURL  string
	chatURL   string
	newsURL   string
	whaleURL  string
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

	log.Info("=== NeruBot API Gateway v1.0.0 ===")
	log.Info("Starting API Gateway service...")

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

	// Get service URLs from environment
	musicURL := getEnvOrDefault("MUSIC_SERVICE_URL", "localhost:8081")
	confURL := getEnvOrDefault("CONFESSION_SERVICE_URL", "localhost:8082")
	roastURL := getEnvOrDefault("ROAST_SERVICE_URL", "localhost:8083")
	chatURL := getEnvOrDefault("CHATBOT_SERVICE_URL", "localhost:8084")
	newsURL := getEnvOrDefault("NEWS_SERVICE_URL", "localhost:8085")
	whaleURL := getEnvOrDefault("WHALE_SERVICE_URL", "localhost:8086")

	log.Info("Service URLs configured",
		"music", musicURL,
		"confession", confURL,
		"roast", roastURL,
		"chatbot", chatURL,
		"news", newsURL,
		"whale", whaleURL,
	)

	// Create Discord session
	discord, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		log.Error("Failed to create Discord session", "error", err)
		os.Exit(1)
	}

	// Create gateway instance
	gw := &Gateway{
		config:   cfg,
		logger:   log,
		discord:  discord,
		musicURL: musicURL,
		confURL:  confURL,
		roastURL: roastURL,
		chatURL:  chatURL,
		newsURL:  newsURL,
		whaleURL: whaleURL,
	}

	// Register Discord handlers
	gw.registerHandlers()

	// Open Discord connection
	if err := discord.Open(); err != nil {
		log.Error("Failed to open Discord connection", "error", err)
		os.Exit(1)
	}
	defer discord.Close()

	log.Info("Discord bot connected successfully")

	// Register slash commands
	if err := gw.registerCommands(); err != nil {
		log.Error("Failed to register commands", "error", err)
		os.Exit(1)
	}

	log.Info("Slash commands registered successfully")

	// Start HTTP health check server
	go gw.startHealthServer()

	log.Info("API Gateway is running. Press CTRL+C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Info("Shutting down API Gateway...")
}

// registerHandlers registers Discord event handlers
func (gw *Gateway) registerHandlers() {
	gw.discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		gw.logger.Info("Bot is ready",
			"username", s.State.User.Username,
			"discriminator", s.State.User.Discriminator,
		)
	})

	// Register interaction handler for slash commands
	gw.discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		gw.handleInteraction(s, i)
	})
}

// registerCommands registers all slash commands
func (gw *Gateway) registerCommands() error {
	commands := []*discordgo.ApplicationCommand{
		// Music commands
		{
			Name:        "play",
			Description: "Play a song from YouTube",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "song",
					Description: "YouTube URL or search query",
					Required:    true,
				},
			},
		},
		{
			Name:        "pause",
			Description: "Pause the current song",
		},
		{
			Name:        "resume",
			Description: "Resume playback",
		},
		{
			Name:        "skip",
			Description: "Skip the current song",
		},
		{
			Name:        "stop",
			Description: "Stop playback and clear queue",
		},
		{
			Name:        "queue",
			Description: "Show the current queue",
		},
		{
			Name:        "nowplaying",
			Description: "Show currently playing song",
		},
		// Confession commands
		{
			Name:        "confess",
			Description: "Submit an anonymous confession",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Your confession",
					Required:    true,
				},
			},
		},
		// Roast commands
		{
			Name:        "roast",
			Description: "Get roasted based on your activity",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to roast (default: yourself)",
					Required:    false,
				},
			},
		},
		{
			Name:        "profile",
			Description: "View your activity profile",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to view (default: yourself)",
					Required:    false,
				},
			},
		},
		// Health check
		{
			Name:        "ping",
			Description: "Check bot status",
		},
	}

	gw.logger.Info("Registering slash commands", "count", len(commands))

	for _, cmd := range commands {
		_, err := gw.discord.ApplicationCommandCreate(gw.discord.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("failed to create command %s: %w", cmd.Name, err)
		}
		gw.logger.Info("Registered command", "name", cmd.Name)
	}

	return nil
}

// handleInteraction handles Discord slash command interactions
func (gw *Gateway) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer interaction response
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get command name
	cmdName := i.ApplicationCommandData().Name

	gw.logger.Info("Received command",
		"command", cmdName,
		"user", i.Member.User.Username,
		"guild", i.GuildID,
	)

	// Route to appropriate handler
	switch cmdName {
	case "play":
		gw.handlePlayCommand(ctx, s, i)
	case "pause":
		gw.handlePauseCommand(ctx, s, i)
	case "resume":
		gw.handleResumeCommand(ctx, s, i)
	case "skip":
		gw.handleSkipCommand(ctx, s, i)
	case "stop":
		gw.handleStopCommand(ctx, s, i)
	case "queue":
		gw.handleQueueCommand(ctx, s, i)
	case "nowplaying":
		gw.handleNowPlayingCommand(ctx, s, i)
	case "confess":
		gw.handleConfessCommand(ctx, s, i)
	case "roast":
		gw.handleRoastCommand(ctx, s, i)
	case "profile":
		gw.handleProfileCommand(ctx, s, i)
	case "ping":
		gw.handlePingCommand(ctx, s, i)
	default:
		gw.respondError(s, i, "Unknown command")
	}
}

// Placeholder command handlers (to be implemented with gRPC calls)
func (gw *Gateway) handlePlayCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Implement gRPC call to music service
	gw.respondMessage(s, i, "🎵 Music feature coming soon! (Will connect to Music Service)")
}

func (gw *Gateway) handlePauseCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "⏸️ Pause feature coming soon!")
}

func (gw *Gateway) handleResumeCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "▶️ Resume feature coming soon!")
}

func (gw *Gateway) handleSkipCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "⏭️ Skip feature coming soon!")
}

func (gw *Gateway) handleStopCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "⏹️ Stop feature coming soon!")
}

func (gw *Gateway) handleQueueCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "📋 Queue feature coming soon!")
}

func (gw *Gateway) handleNowPlayingCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "🎵 Now Playing feature coming soon!")
}

func (gw *Gateway) handleConfessCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "🤐 Confession feature coming soon! (Will connect to Confession Service)")
}

func (gw *Gateway) handleRoastCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "🔥 Roast feature coming soon! (Will connect to Roast Service)")
}

func (gw *Gateway) handleProfileCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "📊 Profile feature coming soon!")
}

func (gw *Gateway) handlePingCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	gw.respondMessage(s, i, "🏓 Pong! API Gateway is running.")
}

// respondMessage sends a simple message response
func (gw *Gateway) respondMessage(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		gw.logger.Error("Failed to respond to interaction", "error", err)
	}
}

// respondError sends an error response
func (gw *Gateway) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		gw.logger.Error("Failed to respond to interaction", "error", err)
	}
}

// startHealthServer starts HTTP health check server
func (gw *Gateway) startHealthServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"api-gateway","version":"1.0.0"}`))
	})

	port := getEnvOrDefault("PORT", "8080")
	gw.logger.Info("Starting health check server", "port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		gw.logger.Error("Health server error", "error", err)
	}
}

// getEnvOrDefault gets environment variable or returns default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
