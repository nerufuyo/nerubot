package discord

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/ai"
	"github.com/nerufuyo/nerubot/internal/pkg/backend"
	lavalinkpkg "github.com/nerufuyo/nerubot/internal/pkg/lavalink"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/pkg/mongodb"
	redispkg "github.com/nerufuyo/nerubot/internal/pkg/redis"
	"github.com/nerufuyo/nerubot/internal/repository"
	"github.com/nerufuyo/nerubot/internal/usecase/analytics"
	"github.com/nerufuyo/nerubot/internal/usecase/chatbot"
	"github.com/nerufuyo/nerubot/internal/usecase/confession"
	"github.com/nerufuyo/nerubot/internal/usecase/fun"
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
	confessionService *confession.ConfessionService
	roastService      *roast.RoastService
	chatbotService    *chatbot.ChatbotService
	newsService       *news.NewsService
	whaleService      *whale.WhaleService
	analyticsService  *analytics.AnalyticsService
	reminderService   *reminder.ReminderService
	funService        *fun.FunService
	musicService      *music.MusicService
	lavalinkClient    *lavalinkpkg.Client
	mongoDB           *mongodb.Client
	redisClient       *redispkg.Client
	backendClient     *backend.Client
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
		bot.reminderService.SetMembersFunc(func(channelID string) []reminder.Member {
			var members []reminder.Member
			// Resolve the guild that owns this channel so we only mention members from that server
			ch, err := bot.session.Channel(channelID)
			if err != nil || ch == nil {
				log.Warn("Failed to resolve channel for members", "channel", channelID, "error", err)
				return members
			}
			guildMembers, err := bot.session.GuildMembers(ch.GuildID, "", 1000)
			if err != nil {
				log.Warn("Failed to fetch guild members", "guild", ch.GuildID, "error", err)
				return members
			}
			for _, m := range guildMembers {
				if m.User != nil && !m.User.Bot {
					members = append(members, reminder.Member{
						ID:       m.User.ID,
						Username: m.User.Username,
					})
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

	// Initialize music service if enabled (must be after Open so session.State.User is set)
	if b.config.Features.Music {
		b.logger.Info("Initializing music service...")
		lavalinkClient := lavalinkpkg.New(b.session, b.logger)
		musicRepo := repository.NewMusicRepository()
		musicSvc := music.NewMusicService(lavalinkClient, b.session, musicRepo, b.redisClient)
		musicSvc.SetSendFunc(func(channelID string, embed *discordgo.MessageEmbed) {
			if _, err := b.session.ChannelMessageSendEmbed(channelID, embed); err != nil {
				b.logger.Error("Failed to send music embed", "channel", channelID, "error", err)
			}
		})
		b.musicService = musicSvc
		b.lavalinkClient = lavalinkClient

		// Register Lavalink voice handlers now that client exists
		b.session.AddHandler(b.lavalinkClient.HandleVoiceStateUpdate)
		b.session.AddHandler(b.lavalinkClient.HandleVoiceServerUpdate)
		b.session.AddHandler(b.onVoiceStateUpdateAutoJoin)
		b.logger.Info("Music service initialized")
	}

	// Register slash commands
	if err := b.registerCommands(); err != nil {
		b.logger.Error("Failed to register commands", "error", err)
		return err
	}

	// Apply bot status from backend dashboard settings (falls back to config default).
	// Also register a callback so all profile changes sync whenever backend refreshes.
	b.applyBotStatus()
	b.backendClient.OnSettingsChange(func(s *backend.BotSettings) {
		if s == nil {
			return
		}
		b.syncBotProfile(s)
	})

	b.logger.Info("Bot is ready and running")

	// Sync guild list to backend
	go b.syncGuildsToBackend()

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

	// Connect to Lavalink node if music is enabled
	if b.lavalinkClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		addr := b.config.Music.LavalinkAddress
		secure := strings.HasSuffix(addr, ".railway.app") || strings.HasSuffix(addr, ".railway.app:443") || strings.Contains(addr, ".up.railway.app")
		lErr := b.lavalinkClient.AddNode(ctx, "main", addr, b.config.Music.LavalinkPassword, secure)
		cancel()
		if lErr != nil {
			b.logger.Error("Failed to connect to Lavalink", "error", lErr)
		} else {
			b.logger.Info("Lavalink node connected")
		}
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

	// Stop music (destroy all players)
	if b.lavalinkClient != nil {
		b.lavalinkClient.Link.Close()
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

// applyBotStatus sets the Discord presence from backend settings (or config fallback).
func (b *Bot) applyBotStatus() {
	status := b.config.Bot.Status
	if s := b.backendClient.GetSettings(); s != nil && s.BotStatus != "" {
		status = s.BotStatus
	}
	if err := b.session.UpdateGameStatus(0, status); err != nil {
		b.logger.Warn("Failed to set status", "error", err)
	}
}

// syncBotProfile synchronises Discord status and username from dashboard settings.
// Called whenever the backend client detects a settings change.
func (b *Bot) syncBotProfile(s *backend.BotSettings) {
	// Sync playing-status
	if s.BotStatus != "" {
		if err := b.session.UpdateGameStatus(0, s.BotStatus); err != nil {
			b.logger.Warn("Failed to update status from backend", "error", err)
		} else {
			b.logger.Info("Bot status synced from dashboard", "status", s.BotStatus)
		}
	}

	// Sync username if changed (Discord allows ~2 changes/hour for bots)
	if s.BotName != "" && b.session.State.User != nil && s.BotName != b.session.State.User.Username {
		_, err := b.session.UserUpdate(s.BotName, "", "")
		if err != nil {
			b.logger.Warn("Failed to update bot username", "name", s.BotName, "error", err)
		} else {
			b.logger.Info("Bot username synced from dashboard", "name", s.BotName)
		}
	}
}

// registerHandlers registers event handlers
func (b *Bot) registerHandlers() {
	b.session.AddHandler(b.onReady)
	b.session.AddHandler(b.onMessageCreate)
	b.session.AddHandler(b.onInteractionCreate)
	b.session.AddHandler(b.onGuildCreate)
}

// onReady handles the ready event
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	b.logger.Info("Bot ready event received",
		"guilds", len(event.Guilds),
		"user", event.User.Username,
	)
}

// onGuildCreate fires when the bot joins a new guild (or on reconnect for existing guilds).
// It registers slash commands so they're instantly available.
func (b *Bot) onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	if err := b.registerCommandsForGuild(g.ID); err != nil {
		b.logger.Warn("Failed to register commands for guild", "guild", g.ID, "error", err)
	}
	// Re-sync guild list when a new guild is joined
	go b.syncGuildsToBackend()
}

// syncGuildsToBackend collects all connected guilds and pushes them to the backend.
func (b *Bot) syncGuildsToBackend() {
	if b.backendClient == nil {
		return
	}

	guilds := make([]backend.GuildInfo, 0, len(b.session.State.Guilds))
	for _, g := range b.session.State.Guilds {
		iconURL := ""
		if g.Icon != "" {
			iconURL = fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", g.ID, g.Icon)
		}
		guilds = append(guilds, backend.GuildInfo{
			GuildID:     g.ID,
			GuildName:   g.Name,
			MemberCount: g.MemberCount,
			IconURL:     iconURL,
		})
	}
	b.backendClient.SyncGuilds(guilds)
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
	// Handle button/component interactions
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID
		if strings.HasPrefix(customID, "music_") {
			b.handleMusicButton(s, i)
		}
		return
	}

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
	// --- Music commands ---
	case "play":
		b.handlePlay(s, i)
	case "pause":
		b.handlePause(s, i)
	case "resume":
		b.handleResume(s, i)
	case "stop":
		b.handleStop(s, i)
	case "skip":
		b.handleSkip(s, i)
	case "nowplaying":
		b.handleNowPlaying(s, i)
	case "queue":
		b.handleQueue(s, i)
	case "volume":
		b.handleVolume(s, i)
	case "remove":
		b.handleRemove(s, i)
	case "clear":
		b.handleClear(s, i)
	case "shuffle":
		b.handleShuffle(s, i)
	case "move":
		b.handleMove(s, i)
	// --- Phase 2: Loop, Filters, Previous, Seek, Playlist ---
	case "loop":
		b.handleLoop(s, i)
	case "filter":
		b.handleFilter(s, i)
	case "previous":
		b.handlePrevious(s, i)
	case "seek":
		b.handleSeek(s, i)
	case "playlist":
		b.handlePlaylist(s, i)
	// --- Phase 3: DJ, VoteSkip, Lyrics, 24/7, Autoplay ---
	case "dj":
		b.handleDJ(s, i)
	case "voteskip":
		b.handleVoteSkip(s, i)
	case "lyrics":
		b.handleLyrics(s, i)
	case "247":
		b.handle247(s, i)
	case "autoplay":
		b.handleAutoplay(s, i)
	case "recommend":
		b.handleRecommend(s, i)
	case "radio":
		b.handleRadio(s, i)
	case "autojoin":
		b.handleAutoJoin(s, i)
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
	commands := b.buildCommands()

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

// registerCommandsForGuild registers slash commands for a single guild.
func (b *Bot) registerCommandsForGuild(guildID string) error {
	commands := b.buildCommands()
	_, err := b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, guildID, commands)
	return err
}

// buildCommands returns the full list of slash commands to register.
func (b *Bot) buildCommands() []*discordgo.ApplicationCommand {
	adminPermission := int64(discordgo.PermissionAdministrator)

	return []*discordgo.ApplicationCommand{
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
		// --- Music commands ---
		{
			Name:        "play",
			Description: "Play a song or add it to the queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "Song name or URL (YouTube, Spotify, SoundCloud)",
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
			Description: "Resume the paused song",
		},
		{
			Name:        "stop",
			Description: "Stop playback and clear the queue",
		},
		{
			Name:        "skip",
			Description: "Skip to the next song in the queue",
		},
		{
			Name:        "nowplaying",
			Description: "Show the currently playing song",
		},
		{
			Name:        "queue",
			Description: "Show the current song queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "page",
					Description: "Page number",
					Required:    false,
				},
			},
		},
		{
			Name:        "volume",
			Description: "Set the player volume (0-150)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "Volume level (0-150)",
					Required:    true,
					MinValue:    func() *float64 { v := 0.0; return &v }(),
					MaxValue:    150,
				},
			},
		},
		{
			Name:        "remove",
			Description: "Remove a song from the queue by position",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "position",
					Description: "Position in queue (1-based)",
					Required:    true,
					MinValue:    func() *float64 { v := 1.0; return &v }(),
				},
			},
		},
		{
			Name:        "clear",
			Description: "Clear the entire queue",
		},
		{
			Name:        "shuffle",
			Description: "Shuffle the current queue",
		},
		{
			Name:        "move",
			Description: "Move a song to a different position in the queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "from",
					Description: "Current position (1-based)",
					Required:    true,
					MinValue:    func() *float64 { v := 1.0; return &v }(),
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "to",
					Description: "New position (1-based)",
					Required:    true,
					MinValue:    func() *float64 { v := 1.0; return &v }(),
				},
			},
		},
		// --- Phase 2 commands ---
		{
			Name:        "loop",
			Description: "Set loop mode for playback",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mode",
					Description: "Loop mode",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Off", Value: "off"},
						{Name: "Song", Value: "song"},
						{Name: "Queue", Value: "queue"},
					},
				},
			},
		},
		{
			Name:        "filter",
			Description: "Apply an audio filter effect",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Filter name (bassboost, nightcore, vaporwave, karaoke, 8d, tremolo, clear)",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Bass Boost", Value: "bassboost"},
						{Name: "Nightcore", Value: "nightcore"},
						{Name: "Vaporwave", Value: "vaporwave"},
						{Name: "Karaoke", Value: "karaoke"},
						{Name: "8D Audio", Value: "8d"},
						{Name: "Tremolo", Value: "tremolo"},
						{Name: "Clear / Off", Value: "clear"},
					},
				},
			},
		},
		{
			Name:        "previous",
			Description: "Play the previous track from history",
		},
		{
			Name:        "seek",
			Description: "Seek to a position in the current song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "seconds",
					Description: "Position in seconds",
					Required:    true,
					MinValue:    func() *float64 { v := 0.0; return &v }(),
				},
			},
		},
		{
			Name:        "playlist",
			Description: "Manage your saved playlists",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "create",
					Description: "Create a playlist (saves current queue if playing)",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Playlist name",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "Add the current song to a playlist",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Playlist name",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "play",
					Description: "Load and play a saved playlist",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Playlist name",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "View all your playlists",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "delete",
					Description: "Delete a playlist",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Playlist name",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "show",
					Description: "View songs in a playlist",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Playlist name",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "import",
					Description: "Import a playlist from Spotify, YouTube, or SoundCloud URL",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "url",
							Description: "Playlist URL",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Name for the imported playlist",
							Required:    true,
						},
					},
				},
			},
		},
		// --- Phase 3 commands ---
		{
			Name:        "dj",
			Description: "Manage DJ role for music control",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set",
					Description: "Set the DJ role",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionRole,
							Name:        "role",
							Description: "The DJ role",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove",
					Description: "Remove the DJ role restriction",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "check",
					Description: "Check current DJ settings and your permissions",
				},
			},
		},
		{
			Name:        "voteskip",
			Description: "Vote to skip the current song (majority needed)",
		},
		{
			Name:        "lyrics",
			Description: "Show lyrics for the currently playing song",
		},
		{
			Name:                     "247",
			Description:              "Toggle 24/7 mode (stay in voice channel)",
			DefaultMemberPermissions: &adminPermission,
		},
		{
			Name:        "autoplay",
			Description: "Toggle autoplay (auto-queue similar songs when queue ends)",
		},
		{
			Name:        "recommend",
			Description: "Get song recommendations based on what's playing",
		},
		{
			Name:        "radio",
			Description: "Start nonstop radio for a genre",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "genre",
					Description: "Genre or mood (e.g., lofi, jazz, rock, chill)",
					Required:    true,
				},
			},
		},
		{
			Name:                     "autojoin",
			Description:              "Configure auto-join for a voice channel",
			DefaultMemberPermissions: &adminPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Voice channel to auto-join (omit to disable)",
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildVoice,
					},
				},
			},
		},
	}
}
