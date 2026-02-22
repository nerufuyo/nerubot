package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/ai"
	"github.com/nerufuyo/nerubot/internal/pkg/backend"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/pkg/mongodb"
	redispkg "github.com/nerufuyo/nerubot/internal/pkg/redis"
	"github.com/nerufuyo/nerubot/internal/repository"
	"github.com/nerufuyo/nerubot/internal/usecase/analytics"
	"github.com/nerufuyo/nerubot/internal/usecase/chatbot"
	"github.com/nerufuyo/nerubot/internal/usecase/confession"
	"github.com/nerufuyo/nerubot/internal/usecase/fun"
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
	confessionService *confession.ConfessionService
	roastService      *roast.RoastService
	chatbotService    *chatbot.ChatbotService
	newsService       *news.NewsService
	whaleService      *whale.WhaleService
	analyticsService  *analytics.AnalyticsService
	reminderService   *reminder.ReminderService
	funService        *fun.FunService
	mongoDB           *mongodb.Client
	redisClient       *redispkg.Client
	backendClient     *backend.Client
	ollamaClient      *ai.OllamaClient
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

	// Initialize AI provider (shared between chatbot and reminder)
	var aiProvider ai.AIProvider
	if cfg.AI.DeepSeekKey != "" {
		aiProvider = ai.NewDeepSeekProvider(cfg.AI.DeepSeekKey)
		log.Info("AI provider initialized", "provider", "DeepSeek")
	}

	// Initialize backend API client for RAG knowledge and bot settings
	backendClient := backend.New(cfg.BackendURL)
	log.Info("Backend API client initialized", "url", cfg.BackendURL)

	bot := &Bot{
		session:           session,
		config:            cfg,
		logger:            log,
		confessionService: confession.NewConfessionService(backendClient),
		roastService:      roast.NewRoastService(backendClient),
		chatbotService:    chatbot.NewChatbotService(cfg.AI.DeepSeekKey, redisClient, backendClient),
		newsService:       news.NewNewsService(),
		whaleService:      whale.NewWhaleService(cfg.Crypto.WhaleAlertAPIKey),
		analyticsService:  analytics.NewAnalyticsService(mongoDB),
		mongoDB:           mongoDB,
		redisClient:       redisClient,
		backendClient:     backendClient,
	}

	// Initialize reminder service if enabled
	if cfg.Features.Reminder {
		bot.reminderService = reminder.NewReminderService(cfg.Reminder.ChannelID, aiProvider)
		bot.reminderService.SetSendFunc(func(channelID, message string) {
			if _, err := bot.session.ChannelMessageSend(channelID, message); err != nil {
				log.Error("Failed to send reminder", "error", err)
			}
		})
		bot.reminderService.SetMembersFunc(func() []reminder.Member {
			var members []reminder.Member
			for _, guild := range bot.session.State.Guilds {
				guildMembers, err := bot.session.GuildMembers(guild.ID, "", 1000)
				if err != nil {
					log.Warn("Failed to fetch guild members", "guild", guild.ID, "error", err)
					continue
				}
				for _, m := range guildMembers {
					if m.User != nil && !m.User.Bot {
						members = append(members, reminder.Member{
							ID:       m.User.ID,
							Username: m.User.Username,
						})
					}
				}
			}
			return members
		})
		if cfg.Reminder.ChannelID != "" {
			log.Info("Reminder service initialized", "channel", cfg.Reminder.ChannelID)
		} else {
			log.Info("Reminder service initialized (no channel set — use /reminder-set)")
		}
	}

	// Initialize Ollama client if configured
	if cfg.AI.OllamaURL != "" {
		bot.ollamaClient = ai.NewOllamaClient(cfg.AI.OllamaURL)
		if bot.ollamaClient.IsAvailable() {
			log.Info("Ollama client initialized", "url", cfg.AI.OllamaURL)
		} else {
			log.Warn("Ollama server not reachable", "url", cfg.AI.OllamaURL)
		}
	}

	// Initialize fun service (dad jokes + memes)
	bot.funService = fun.NewFunService()
	bot.funService.SetSendFunc(func(channelID string, embed *fun.FunEmbed) {
		discordEmbed := &discordgo.MessageEmbed{
			Title:       embed.Title,
			Description: embed.Description,
			Color:       embed.Color,
			URL:         embed.URL,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: embed.Footer,
			},
		}
		if embed.ImageURL != "" {
			discordEmbed.Image = &discordgo.MessageEmbedImage{URL: embed.ImageURL}
		}
		if embed.Content != "" {
			// Send with text content (for mentions) + embed
			msg := &discordgo.MessageSend{
				Content: embed.Content,
				Embeds:  []*discordgo.MessageEmbed{discordEmbed},
			}
			if _, err := bot.session.ChannelMessageSendComplex(channelID, msg); err != nil {
				log.Error("Failed to send scheduled fun message", "channel", channelID, "error", err)
			}
		} else {
			if _, err := bot.session.ChannelMessageSendEmbed(channelID, discordEmbed); err != nil {
				log.Error("Failed to send scheduled fun embed", "channel", channelID, "error", err)
			}
		}
	})
	log.Info("Fun service initialized (dad jokes + memes)")

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

	// Load saved guild configs from DB (restores reminder channels, etc.)
	b.loadSavedGuildConfigs()

	// Start background services
	if b.reminderService != nil {
		b.reminderService.Start()
	}

	// Start fun service scheduler (dad jokes + memes)
	if b.funService != nil {
		b.funService.Start()
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

	// Stop fun service scheduler
	if b.funService != nil {
		b.funService.Stop()
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
	case "reminder-stop":
		b.handleReminderStop(s, i)
	case "dadjoke":
		b.handleDadJoke(s, i)
	case "dadjoke-setup":
		b.handleDadJokeSetup(s, i)
	case "meme":
		b.handleMeme(s, i)
	case "meme-setup":
		b.handleMemeSetup(s, i)
	case "mentalhealth":
		b.handleMentalHealth(s, i)
	case "mentalhealth-setup":
		b.handleMentalHealthSetup(s, i)
	case "mentalhealth-stop":
		b.handleMentalHealthStop(s, i)
	case "ollama-models":
		b.handleOllamaModels(s, i)
	case "ollama-bench":
		b.handleOllamaBench(s, i)
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "Response language (default: EN)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "Response language (default: EN)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "News language/region (default: EN)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
				},
			},
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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "Help language (default: EN)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
				},
			},
		},
		{
			Name:        "reminder",
			Description: "View upcoming Indonesian holidays and Ramadan schedule",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "Response language (default: EN)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
				},
			},
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "Language for auto-reminders (default: random)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
				},
			},
		},
		{
			Name:                     "reminder-stop",
			Description:              "Stop automatic reminders",
			DefaultMemberPermissions: &adminPermission,
		},
		// --- Fun commands ---
		{
			Name:        "dadjoke",
			Description: "Get a random (clean) dad joke",
		},
		{
			Name:                     "dadjoke-setup",
			Description:              "Schedule automatic dad jokes in a channel",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to post dad jokes in",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "interval",
					Description: "Interval in minutes (e.g. 60 for hourly, 0 to disable)",
					Required:    true,
				},
			},
		},
		{
			Name:        "meme",
			Description: "Get a random meme from the internet",
		},
		{
			Name:                     "meme-setup",
			Description:              "Schedule automatic memes in a channel",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to post memes in",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "interval",
					Description: "Interval in minutes (e.g. 60 for hourly, 0 to disable)",
					Required:    true,
				},
			},
		},

		// --- Mental health ---
		{
			Name:        "mentalhealth",
			Description: "Get a mental health tip & self-care reminder",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "Language (default: EN)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
				},
			},
		},
		{
			Name:                     "mentalhealth-setup",
			Description:              "Schedule mental health reminders in a channel with mentioning",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to post mental health reminders in",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "interval",
					Description: "Interval in minutes (e.g. 60 for hourly, 0 to disable)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Name:        "tag",
					Description: "User or role to mention in each reminder (optional)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "everyone",
					Description: "Mention @everyone in each reminder (default: false)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "lang",
					Description: "Language for tips (default: EN)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "English", Value: "EN"},
						{Name: "Bahasa Indonesia", Value: "ID"},
						{Name: "日本語 (Japanese)", Value: "JP"},
						{Name: "한국어 (Korean)", Value: "KR"},
						{Name: "中文 (Chinese)", Value: "ZH"},
					},
				},
			},
		},
		{
			Name:                     "mentalhealth-stop",
			Description:              "Stop scheduled mental health reminders",
			DefaultMemberPermissions: &adminPermission,
		},

		// --- Ollama commands ---
		{
			Name:        "ollama-models",
			Description: "List available models on the Ollama server",
		},
		{
			Name:        "ollama-bench",
			Description: "Benchmark an Ollama model (tokens/sec, TTFT, etc.)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "model",
					Description: "Model name to benchmark (e.g. gemma3:1b, llama3.2:latest)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "Custom prompt (default: explain LLMs in 3 sentences)",
					Required:    false,
				},
			},
		},
	}

	// Bulk-overwrite global commands: atomically sets exactly these commands and removes any others.
	b.logger.Info("Registering slash commands", "count", len(commands))

	// Register commands per-guild for instant availability (global commands take up to 1h).
	for _, guild := range b.session.State.Guilds {
		_, err := b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, guild.ID, commands)
		if err != nil {
			b.logger.Warn("Failed to register guild commands", "guild", guild.ID, "error", err)
		} else {
			b.logger.Info("Guild commands registered", "guild", guild.ID)
		}
	}

	b.logger.Info("Slash commands registered successfully")
	return nil
}
