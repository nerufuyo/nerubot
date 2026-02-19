package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Bot      BotConfig
	Limits   Limits
	Audio    AudioConfig
	Features FeatureFlags
	Discord  DiscordConfig
	AI       AIConfig
	Crypto   CryptoConfig
	Lavalink LavalinkConfig
	Reminder ReminderConfig
}

// BotConfig holds basic bot configuration
type BotConfig struct {
	Name        string
	Version     string
	Token       string
	Prefix      string
	Status      string
	Description string
	Author      string
	Website     string
}

// Limits holds rate limits and timeouts
type Limits struct {
	MaxQueueSize         int
	MaxSearchResults     int
	MaxSongDuration      time.Duration
	SearchTimeout        time.Duration
	ConversionTimeout    time.Duration
	IdleDisconnectTime   time.Duration
	VoiceConnectTimeout  time.Duration
	CommandsPerMinute    int
	SearchesPerMinute    int
	ConfessionCooldown   time.Duration
	RoastCooldown        time.Duration
}

// AudioConfig holds FFmpeg and audio settings
type AudioConfig struct {
	FFmpegPath      string
	YtdlpPath       string
	FFmpegOptions   FFmpegOptions
	OpusPaths       []string
	Bitrate         int
	SampleRate      int
	Channels        int
	DefaultVolume   float64
}

// FFmpegOptions holds FFmpeg command options
type FFmpegOptions struct {
	BeforeOptions string
	Options       string
}

// FeatureFlags controls which features are enabled
type FeatureFlags struct {
	Music         bool
	News          bool
	HelpSystem    bool
	Chatbot       bool
	Confession    bool
	Roast         bool
	WhaleAlerts   bool
	AutoDisconnect bool
	Mode247       bool
	Reminder      bool
}

// DiscordConfig holds Discord-specific configuration
type DiscordConfig struct {
	Colors                map[string]int
	SyncCommandsOnReady   bool
	SyncCommandsGlobally  bool
}

// AIConfig holds AI service configuration
type AIConfig struct {
	DeepSeekKey string
}

// CryptoConfig holds cryptocurrency feature configuration
type CryptoConfig struct {
	WhaleAlertAPIKey   string
	TwitterAPIKey      string
	TwitterAPISecret   string
	TwitterAccessToken string
	TwitterAccessSecret string
}

// LavalinkConfig holds Lavalink server configuration
type LavalinkConfig struct {
	Host     string
	Port     int
	Password string
	Enabled  bool
}

// ReminderConfig holds reminder feature configuration.
type ReminderConfig struct {
	ChannelID string // Discord channel ID for posting reminders
}

// MusicSources holds configuration for music source providers
type MusicSources struct {
	YouTube    MusicSource
	Spotify    MusicSource
	SoundCloud MusicSource
	Direct     MusicSource
}

// MusicSource holds configuration for a single music source
type MusicSource struct {
	Enabled  bool
	Emoji    string
	Name     string
	Priority int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get required configuration
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable is required")
	}

	// Build configuration with defaults
	cfg := &Config{
		Bot: BotConfig{
			Name:        AppName,
			Version:     AppVersion,
			Token:       token,
			Prefix:      getEnvOrDefault("COMMAND_PREFIX", "!"),
			Status:      "Ready to rock your server!",
			Description: "Your friendly Discord companion!",
			Author:      "nerufuyo",
			Website:     "https://github.com/nerufuyo/nerubot",
		},
		Limits: Limits{
			MaxQueueSize:         getEnvAsInt("MAX_QUEUE_SIZE", 100),
			MaxSearchResults:     5,
			MaxSongDuration:      time.Hour,
			SearchTimeout:        15 * time.Second,
			ConversionTimeout:    20 * time.Second,
			IdleDisconnectTime:   time.Duration(getEnvAsInt("IDLE_DISCONNECT_TIME", 300)) * time.Second,
			VoiceConnectTimeout:  30 * time.Second,
			CommandsPerMinute:    10,
			SearchesPerMinute:    5,
			ConfessionCooldown:   10 * time.Minute,
			RoastCooldown:        5 * time.Minute,
		},
		Audio: AudioConfig{
			FFmpegPath: os.Getenv("FFMPEG_PATH"),
			YtdlpPath:  os.Getenv("YTDLP_PATH"),
			FFmpegOptions: FFmpegOptions{
				BeforeOptions: "-reconnect 1 -reconnect_streamed 1 -reconnect_delay_max 5",
				Options:       "-vn -filter:a volume=0.5",
			},
			OpusPaths: []string{
				"/opt/homebrew/lib/libopus.dylib",
				"/usr/local/lib/libopus.dylib",
				"/opt/homebrew/lib/libopus.0.dylib",
				"/usr/local/lib/libopus.0.dylib",
				"/usr/lib/x86_64-linux-gnu/libopus.so.0",
				"/usr/lib/libopus.so.0",
				"libopus.so.0",
				"libopus.dylib",
				"opus",
			},
			Bitrate:       128,
			SampleRate:    48000,
			Channels:      2,
			DefaultVolume: 0.5,
		},
		Features: FeatureFlags{
			Music:          getEnvAsBool("ENABLE_MUSIC", false),
			News:           true,
			HelpSystem:     true,
			Chatbot:        hasAnyAIKey(),
			Confession:     getEnvAsBool("ENABLE_CONFESSION", true),
			Roast:          getEnvAsBool("ENABLE_ROAST", true),
			WhaleAlerts:    os.Getenv("WHALE_ALERT_API_KEY") != "",
			AutoDisconnect: true,
			Mode247:        getEnvAsBool("ENABLE_24_7", false),			Reminder:       getEnvAsBool("ENABLE_REMINDER", true),		},
		Discord: DiscordConfig{
			Colors: map[string]int{
				"primary":    0x0099FF,
				"secondary":  0x6C757D,
				"success":    0x00FF00,
				"error":      0xFF0000,
				"warning":    0xFFA500,
				"info":       0x0099FF,
				"music":      0x9932CC,
				"spotify":    0x1DB954,
				"youtube":    0xFF0000,
				"soundcloud": 0xFF7700,
			},
			SyncCommandsOnReady:  true,
			SyncCommandsGlobally: true,
		},
		AI: AIConfig{
			DeepSeekKey: os.Getenv("DEEPSEEK_API_KEY"),
		},
		Crypto: CryptoConfig{
			WhaleAlertAPIKey:    os.Getenv("WHALE_ALERT_API_KEY"),
			TwitterAPIKey:       os.Getenv("TWITTER_API_KEY"),
			TwitterAPISecret:    os.Getenv("TWITTER_API_SECRET"),
			TwitterAccessToken:  os.Getenv("TWITTER_ACCESS_TOKEN"),
			TwitterAccessSecret: os.Getenv("TWITTER_ACCESS_SECRET"),
		},
		Lavalink: LavalinkConfig{
			Host:     getEnvOrDefault("LAVALINK_HOST", "localhost"),
			Port:     getEnvAsInt("LAVALINK_PORT", 2333),
			Password: getEnvOrDefault("LAVALINK_PASSWORD", "youshallnotpass"),
			Enabled:  getEnvAsBool("LAVALINK_ENABLED", false),
		},
		Reminder: ReminderConfig{
			ChannelID: os.Getenv("REMINDER_CHANNEL_ID"),
		},
	}

	return cfg, nil
}

// Helper functions

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func hasAnyAIKey() bool {
	return os.Getenv("DEEPSEEK_API_KEY") != ""
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Bot.Token == "" {
		return fmt.Errorf("Discord bot token is required")
	}
	if c.Limits.MaxQueueSize <= 0 {
		return fmt.Errorf("max queue size must be positive")
	}
	return nil
}
