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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	chatbotpb "github.com/nerufuyo/nerubot/api/proto/chatbot"
	confessionpb "github.com/nerufuyo/nerubot/api/proto/confession"
	musicpb "github.com/nerufuyo/nerubot/api/proto/music"
	newspb "github.com/nerufuyo/nerubot/api/proto/news"
	roastpb "github.com/nerufuyo/nerubot/api/proto/roast"
	whalepb "github.com/nerufuyo/nerubot/api/proto/whale"
)

const (
	ServiceName    = "api-gateway"
	ServiceVersion = "1.0.0"
)

type Gateway struct {
	config        *config.Config
	logger        *logger.Logger
	discord       *discordgo.Session
	musicClient   musicpb.MusicServiceClient
	confClient    confessionpb.ConfessionServiceClient
	roastClient   roastpb.RoastServiceClient
	chatClient    chatbotpb.ChatbotServiceClient
	newsClient    newspb.NewsServiceClient
	whaleClient   whalepb.WhaleServiceClient
	musicConn     *grpc.ClientConn
	confConn      *grpc.ClientConn
	roastConn     *grpc.ClientConn
	chatConn      *grpc.ClientConn
	newsConn      *grpc.ClientConn
	whaleConn     *grpc.ClientConn
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
	musicURL := getEnvOrDefault("MUSIC_SERVICE_URL", "localhost:50051")
	confURL := getEnvOrDefault("CONFESSION_SERVICE_URL", "localhost:50052")
	roastURL := getEnvOrDefault("ROAST_SERVICE_URL", "localhost:50053")
	chatURL := getEnvOrDefault("CHATBOT_SERVICE_URL", "localhost:50054")
	newsURL := getEnvOrDefault("NEWS_SERVICE_URL", "localhost:50055")
	whaleURL := getEnvOrDefault("WHALE_SERVICE_URL", "localhost:50056")

	log.Info("Service URLs configured",
		"music", musicURL,
		"confession", confURL,
		"roast", roastURL,
		"chatbot", chatURL,
		"news", newsURL,
		"whale", whaleURL,
	)

	// Connect to gRPC services
	musicConn, err := grpc.NewClient(musicURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to Music service", "error", err)
		os.Exit(1)
	}
	defer musicConn.Close()

	confConn, err := grpc.NewClient(confURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to Confession service", "error", err)
		os.Exit(1)
	}
	defer confConn.Close()

	roastConn, err := grpc.NewClient(roastURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to Roast service", "error", err)
		os.Exit(1)
	}
	defer roastConn.Close()

	chatConn, err := grpc.NewClient(chatURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to Chatbot service", "error", err)
		os.Exit(1)
	}
	defer chatConn.Close()

	newsConn, err := grpc.NewClient(newsURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to News service", "error", err)
		os.Exit(1)
	}
	defer newsConn.Close()

	whaleConn, err := grpc.NewClient(whaleURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to Whale service", "error", err)
		os.Exit(1)
	}
	defer whaleConn.Close()

	log.Info("gRPC connections established to all 6 backend services")

	// Create Discord session
	discord, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		log.Error("Failed to create Discord session", "error", err)
		os.Exit(1)
	}

	// Create gateway instance
	gw := &Gateway{
		config:      cfg,
		logger:      log,
		discord:     discord,
		musicClient: musicpb.NewMusicServiceClient(musicConn),
		confClient:  confessionpb.NewConfessionServiceClient(confConn),
		roastClient: roastpb.NewRoastServiceClient(roastConn),
		chatClient:  chatbotpb.NewChatbotServiceClient(chatConn),
		newsClient:  newspb.NewNewsServiceClient(newsConn),
		whaleClient: whalepb.NewWhaleServiceClient(whaleConn),
		musicConn:   musicConn,
		confConn:    confConn,
		roastConn:   roastConn,
		chatConn:    chatConn,
		newsConn:    newsConn,
		whaleConn:   whaleConn,
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
	// Get song parameter
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		gw.respondError(s, i, "Please provide a song URL or search query")
		return
	}

	song := options[0].StringValue()
	guildID := i.GuildID
	userID := i.Member.User.ID

	// Defer response to avoid timeout
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		gw.logger.Error("Failed to defer response", "error", err)
		return
	}

	// Call Music service
	resp, err := gw.musicClient.Play(ctx, &musicpb.PlayRequest{
		GuildId:     guildID,
		SongUrl:     song,
		RequestedBy: userID,
		ChannelId:   i.ChannelID,
	})
	if err != nil {
		gw.followUp(s, i, fmt.Sprintf("❌ Failed to play song: %v", err))
		return
	}

	if !resp.Success {
		gw.followUp(s, i, fmt.Sprintf("❌ %s", resp.Message))
		return
	}

	// Send response
	message := fmt.Sprintf("🎵 %s", resp.Message)
	if resp.Song != nil {
		message = fmt.Sprintf("🎵 Now playing: **%s**\nRequested by: <@%s>", resp.Song.Title, resp.Song.RequestedBy)
	}
	gw.followUp(s, i, message)
}

func (gw *Gateway) handlePauseCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := gw.musicClient.Pause(ctx, &musicpb.PauseRequest{
		GuildId: i.GuildID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to pause: %v", err))
		return
	}
	gw.respondMessage(s, i, fmt.Sprintf("⏸️ %s", resp.Message))
}

func (gw *Gateway) handleResumeCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := gw.musicClient.Resume(ctx, &musicpb.ResumeRequest{
		GuildId: i.GuildID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to resume: %v", err))
		return
	}
	gw.respondMessage(s, i, fmt.Sprintf("▶️ %s", resp.Message))
}

func (gw *Gateway) handleSkipCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := gw.musicClient.Skip(ctx, &musicpb.SkipRequest{
		GuildId: i.GuildID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to skip: %v", err))
		return
	}
	message := fmt.Sprintf("⏭️ %s", resp.Message)
	if resp.NextSong != nil {
		message += fmt.Sprintf("\nNow playing: **%s**", resp.NextSong.Title)
	}
	gw.respondMessage(s, i, message)
}

func (gw *Gateway) handleStopCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := gw.musicClient.Stop(ctx, &musicpb.StopRequest{
		GuildId: i.GuildID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to stop: %v", err))
		return
	}
	gw.respondMessage(s, i, fmt.Sprintf("⏹️ %s", resp.Message))
}

func (gw *Gateway) handleQueueCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := gw.musicClient.GetQueue(ctx, &musicpb.GetQueueRequest{
		GuildId: i.GuildID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to get queue: %v", err))
		return
	}

	if resp.TotalSongs == 0 {
		gw.respondMessage(s, i, "📋 Queue is empty")
		return
	}

	message := fmt.Sprintf("📋 **Queue** (%d songs)\n\n", resp.TotalSongs)
	for idx, song := range resp.Songs {
		if idx >= 10 { // Limit to 10 songs
			message += fmt.Sprintf("\n... and %d more", resp.TotalSongs-10)
			break
		}
		message += fmt.Sprintf("%d. **%s**\n", idx+1, song.Title)
	}
	gw.respondMessage(s, i, message)
}

func (gw *Gateway) handleNowPlayingCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := gw.musicClient.NowPlaying(ctx, &musicpb.NowPlayingRequest{
		GuildId: i.GuildID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to get now playing: %v", err))
		return
	}

	if !resp.IsPlaying || resp.CurrentSong == nil {
		gw.respondMessage(s, i, "🎵 Nothing is playing")
		return
	}

	message := fmt.Sprintf("🎵 **Now Playing**\n%s\nRequested by: <@%s>\nVolume: %d%%\nLoop: %s",
		resp.CurrentSong.Title,
		resp.CurrentSong.RequestedBy,
		resp.Volume,
		resp.LoopMode,
	)
	gw.respondMessage(s, i, message)
}

func (gw *Gateway) handleConfessCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		gw.respondError(s, i, "Please provide a confession message")
		return
	}

	content := options[0].StringValue()
	userID := i.Member.User.ID

	resp, err := gw.confClient.Submit(ctx, &confessionpb.SubmitRequest{
		GuildId:     i.GuildID,
		Content:     content,
		SubmittedBy: userID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to submit confession: %v", err))
		return
	}

	if !resp.Success {
		gw.respondError(s, i, resp.Message)
		return
	}

	gw.respondMessage(s, i, fmt.Sprintf("🤐 %s\nConfession ID: #%d", resp.Message, resp.ConfessionId))
}

func (gw *Gateway) handleRoastCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get target user (default to command author)
	targetID := i.Member.User.ID
	options := i.ApplicationCommandData().Options
	if len(options) > 0 && options[0].Type == discordgo.ApplicationCommandOptionUser {
		targetID = options[0].UserValue(s).ID
	}

	resp, err := gw.roastClient.GenerateRoast(ctx, &roastpb.RoastRequest{
		GuildId: i.GuildID,
		UserId:  targetID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to generate roast: %v", err))
		return
	}

	if !resp.Success {
		if resp.CooldownRemaining > 0 {
			gw.respondError(s, i, fmt.Sprintf("Roast on cooldown! Wait %d seconds", resp.CooldownRemaining))
		} else {
			gw.respondError(s, i, "Failed to generate roast")
		}
		return
	}

	message := fmt.Sprintf("🔥 **Roast for <@%s>**\n\n%s\n\n*Category: %s*", targetID, resp.RoastContent, resp.RoastCategory)
	gw.respondMessage(s, i, message)
}

func (gw *Gateway) handleProfileCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get target user (default to command author)
	targetID := i.Member.User.ID
	options := i.ApplicationCommandData().Options
	if len(options) > 0 && options[0].Type == discordgo.ApplicationCommandOptionUser {
		targetID = options[0].UserValue(s).ID
	}

	resp, err := gw.roastClient.GetProfile(ctx, &roastpb.ProfileRequest{
		GuildId: i.GuildID,
		UserId:  targetID,
	})
	if err != nil {
		gw.respondError(s, i, fmt.Sprintf("Failed to get profile: %v", err))
		return
	}

	if resp.Profile == nil {
		gw.respondMessage(s, i, fmt.Sprintf("📊 No profile found for <@%s>", targetID))
		return
	}

	p := resp.Profile
	message := fmt.Sprintf("📊 **Profile for <@%s>**\n\n"+
		"Messages: %d\n"+
		"Reactions: %d\n"+
		"Voice Time: %d minutes\n"+
		"Commands: %d\n"+
		"Last Seen: %s",
		targetID,
		p.MessageCount,
		p.ReactionCount,
		p.VoiceMinutes,
		p.CommandCount,
		p.LastSeen,
	)
	gw.respondMessage(s, i, message)
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

// followUp sends a follow-up message after deferred response
func (gw *Gateway) followUp(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: message,
	})
	if err != nil {
		gw.logger.Error("Failed to send follow-up message", "error", err)
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
