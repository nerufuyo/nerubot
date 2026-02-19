package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/lavalink"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/pkg/mongodb"
	redispkg "github.com/nerufuyo/nerubot/internal/pkg/redis"
	"github.com/nerufuyo/nerubot/internal/repository"
	"github.com/nerufuyo/nerubot/internal/usecase/analytics"
	"github.com/nerufuyo/nerubot/internal/usecase/chatbot"
	"github.com/nerufuyo/nerubot/internal/usecase/confession"
	"github.com/nerufuyo/nerubot/internal/usecase/music"
	"github.com/nerufuyo/nerubot/internal/usecase/news"
	"github.com/nerufuyo/nerubot/internal/usecase/reminder"
	"github.com/nerufuyo/nerubot/internal/usecase/roast"
	"github.com/nerufuyo/nerubot/internal/usecase/whale"
)

// Bot represents the Discord bot
type Bot struct {
	session           *discordgo.Session
	config            *config.Config
	logger            *logger.Logger
	musicService      *music.MusicService
	confessionService *confession.ConfessionService
	roastService      *roast.RoastService
	chatbotService    *chatbot.ChatbotService
	newsService       *news.NewsService
	whaleService      *whale.WhaleService
	analyticsService  *analytics.AnalyticsService
	reminderService   *reminder.ReminderService
	lavalinkClient    *lavalink.Client
	mongoDB           *mongodb.Client
	redisClient       *redispkg.Client
}

// New creates a new Discord bot instance
func New(cfg *config.Config) (*Bot, error) {
	log := logger.New("discord")

	// --- Connect to MongoDB ---
	mongoDB, err := mongodb.New(cfg.Mongo.URL, cfg.Mongo.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ensure indexes
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := mongoDB.EnsureIndexes(ctx); err != nil {
		log.Warn("Failed to ensure MongoDB indexes", "error", err)
	}

	// Set the shared MongoDB client for repositories
	repository.SetMongo(mongoDB)

	// --- Connect to Redis ---
	var redisClient *redispkg.Client
	if cfg.Redis.URL != "" {
		rc, err := redispkg.New(cfg.Redis.URL)
		if err != nil {
			log.Warn("Redis unavailable, sessions will be in-memory only", "error", err)
		} else {
			redisClient = rc
		}
	}

	// Create Discord session
	session, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// Set intents
	session.Identify.Intents = discordgo.IntentsAll

	// Initialize services
	var musicService *music.MusicService
	if cfg.Features.Music {
		ms, err := music.NewMusicService()
		if err != nil {
			log.Warn("Music service disabled", "error", err)
		} else {
			musicService = ms
			log.Info("Music service initialized")
		}
	}

	// Initialize Lavalink client if enabled
	var lavalinkClient *lavalink.Client
	if cfg.Lavalink.Enabled {
		lavalinkClient = lavalink.NewClient(
			cfg.Lavalink.Host,
			cfg.Lavalink.Port,
			cfg.Lavalink.Password,
			"", // UserID will be set after bot connects
		)
		log.Info("Lavalink client initialized", "host", cfg.Lavalink.Host, "port", cfg.Lavalink.Port)
	}

	bot := &Bot{
		session:           session,
		config:            cfg,
		logger:            log,
		musicService:      musicService,
		confessionService: confession.NewConfessionService(),
		roastService:      roast.NewRoastService(),
		chatbotService:    chatbot.NewChatbotService(cfg.AI.DeepSeekKey, redisClient),
		newsService:       news.NewNewsService(),
		whaleService:      whale.NewWhaleService(cfg.Crypto.WhaleAlertAPIKey),
		analyticsService:  analytics.NewAnalyticsService(mongoDB),
		lavalinkClient:    lavalinkClient,
		mongoDB:           mongoDB,
		redisClient:       redisClient,
	}

	// Initialize reminder service if enabled
	if cfg.Features.Reminder {
		bot.reminderService = reminder.NewReminderService(cfg.Reminder.ChannelID)
		bot.reminderService.SetSendFunc(func(channelID, message string) {
			if _, err := bot.session.ChannelMessageSend(channelID, message); err != nil {
				log.Error("Failed to send reminder", "error", err)
			}
		})
		if cfg.Reminder.ChannelID != "" {
			log.Info("Reminder service initialized", "channel", cfg.Reminder.ChannelID)
		} else {
			log.Info("Reminder service initialized (no channel set — use /reminder-set)")
		}
	}

	// Register event handlers
	bot.registerHandlers()

	return bot, nil
}

// Start opens the Discord connection, registers commands, and sets bot status.
// It returns once the bot is ready. Shutdown is handled externally via context cancellation + Stop().
func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("Starting Discord bot...")

	// Open WebSocket connection
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord connection: %w", err)
	}

	b.logger.Info("Discord bot connected",
		"user", b.session.State.User.Username,
		"id", b.session.State.User.ID,
	)

	// Register slash commands
	if err := b.registerCommands(); err != nil {
		b.logger.Error("Failed to register commands", "error", err)
		return err
	}

	// Set bot status
	if err := b.session.UpdateGameStatus(0, b.config.Bot.Status); err != nil {
		b.logger.Warn("Failed to set status", "error", err)
	}

	b.logger.Info("Bot is ready and running")

	// Start background services
	if b.reminderService != nil {
		b.reminderService.Start()
	}

	return nil
}

// Stop stops the bot gracefully
func (b *Bot) Stop() error {
	b.logger.Info("Shutting down bot...")

	// Save analytics data
	if b.analyticsService != nil {
		if err := b.analyticsService.Stop(); err != nil {
			b.logger.Error("Failed to save analytics", "error", err)
		}
	}

	// Stop reminder service
	if b.reminderService != nil {
		b.reminderService.Stop()
	}

	// Close Discord connection
	if err := b.session.Close(); err != nil {
		b.logger.Error("Failed to close Discord connection", "error", err)
	}

	// Disconnect Redis
	if b.redisClient != nil {
		if err := b.redisClient.Close(); err != nil {
			b.logger.Error("Failed to close Redis", "error", err)
		}
	}

	// Disconnect MongoDB
	if b.mongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := b.mongoDB.Disconnect(ctx); err != nil {
			b.logger.Error("Failed to close MongoDB", "error", err)
		}
	}

	b.logger.Info("Bot stopped successfully")
	return nil
}

// registerHandlers registers event handlers
func (b *Bot) registerHandlers() {
	b.session.AddHandler(b.onReady)
	b.session.AddHandler(b.onMessageCreate)
	b.session.AddHandler(b.onInteractionCreate)
}

// onReady handles the ready event
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	b.logger.Info("Bot ready event received",
		"guilds", len(event.Guilds),
		"user", event.User.Username,
	)
}

// onMessageCreate handles message creation events
func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Track activity for roast system
	if m.GuildID != "" {
		_ = b.roastService.TrackMessage(
			m.Author.ID,
			m.GuildID,
			m.Author.Username,
			m.ChannelID,
		)
	}
}

// onInteractionCreate handles slash command interactions
func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Get command name
	cmdName := i.ApplicationCommandData().Name
	startTime := time.Now()

	b.logger.Debug("Command received",
		"command", cmdName,
		"user", i.Member.User.Username,
		"guild", i.GuildID,
	)

	// Route to appropriate handler
	switch cmdName {
	case "play":
		b.handlePlay(s, i)
	case "skip":
		b.handleSkip(s, i)
	case "stop":
		b.handleStop(s, i)
	case "queue":
		b.handleQueue(s, i)
	case "confess":
		b.handleConfess(s, i)
	case "roast":
		b.handleRoast(s, i)
	case "chat":
		b.handleChat(s, i)
	case "chat-reset":
		b.handleChatReset(s, i)
	case "news":
		b.handleNews(s, i)
	case "whale":
		b.handleWhale(s, i)
	case "stats":
		b.handleStats(s, i)
	case "profile":
		b.handleProfile(s, i)
	case "help":
		b.handleHelp(s, i)
	case "reminder":
		b.handleReminder(s, i)
	case "reminder-set":
		b.handleReminderSet(s, i)
	default:
		b.respondError(s, i, "Unknown command")
	}

	// Record command usage in analytics
	if b.analyticsService != nil {
		executionTime := time.Since(startTime).Milliseconds()
		guildName := i.GuildID
		if guild, err := s.Guild(i.GuildID); err == nil {
			guildName = guild.Name
		}
		b.analyticsService.RecordCommandUsage(
			i.GuildID,
			guildName,
			i.Member.User.ID,
			i.Member.User.Username,
			cmdName,
			true,
			executionTime,
		)
	}
}


// registerCommands registers slash commands with Discord
func (b *Bot) registerCommands() error {
	adminPermission := int64(discordgo.PermissionAdministrator)

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Play a song or add to queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "Song name or URL",
					Required:    true,
				},
			},
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
			Description: "Show the music queue",
		},
		{
			Name:        "confess",
			Description: "Submit an anonymous confession",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "Your confession",
					Required:    true,
				},
			},
		},
		{
			Name:        "roast",
			Description: "Get roasted based on your activity",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to roast (optional)",
					Required:    false,
				},
			},
		},
		{
			Name:        "chat",
			Description: "Chat with AI (supports Claude, Gemini, OpenAI)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Your message to the AI",
					Required:    true,
				},
			},
		},
		{
			Name:        "chat-reset",
			Description: "Reset your chat history",
		},
		{
			Name:        "news",
			Description: "Get latest news from multiple sources",
		},
		{
			Name:        "whale",
			Description: "Get recent whale cryptocurrency transactions",
		},
		{
			Name:        "stats",
			Description: "View server statistics and analytics",
		},
		{
			Name:        "profile",
			Description: "View user statistics and activity",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to view (optional, defaults to you)",
					Required:    false,
				},
			},
		},
		{
			Name:        "help",
			Description: "Show help information",
		},
		{
			Name:        "reminder",
			Description: "View upcoming Indonesian holidays and Ramadan schedule",
		},
		{
			Name:                     "reminder-set",
			Description:              "Set the channel for automatic reminders",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to send reminders to",
					Required:    true,
				},
			},
		},
	}

	// Bulk-overwrite: atomically sets exactly these commands and removes any others.
	b.logger.Info("Registering slash commands", "count", len(commands))

	_, err := b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", commands)
	if err != nil {
		return fmt.Errorf("failed to bulk-overwrite commands: %w", err)
	}

	b.logger.Info("Slash commands registered successfully")
	return nil
}
