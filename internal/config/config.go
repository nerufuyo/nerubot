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
	Bot        BotConfig
	Limits     Limits
	Features   FeatureFlags
	Discord    DiscordConfig
	AI         AIConfig
	Crypto     CryptoConfig
	Reminder   ReminderConfig
	Mongo      MongoConfig
	Redis      RedisConfig
	BackendURL string // nerufuyo-workspace-backend API URL for RAG & settings
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
	CommandsPerMinute  int
	ConfessionCooldown time.Duration
	RoastCooldown      time.Duration
}

// FeatureFlags controls which features are enabled
type FeatureFlags struct {
	News       bool
	HelpSystem bool
	Chatbot    bool
	Confession bool
	Roast      bool
	WhaleAlerts bool
	Reminder   bool
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
	OllamaURL   string // Ollama server base URL (e.g. https://ai.infantai.tech)
}

// CryptoConfig holds cryptocurrency feature configuration
type CryptoConfig struct {
	WhaleAlertAPIKey   string
	TwitterAPIKey      string
	TwitterAPISecret   string
	TwitterAccessToken string
	TwitterAccessSecret string
}

// ReminderConfig holds reminder feature configuration.
type ReminderConfig struct {
	ChannelID string // Discord channel ID for posting reminders
}

// MongoConfig holds MongoDB connection configuration.
type MongoConfig struct {
	URL      string
	Database string
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	URL string
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
			CommandsPerMinute:  10,
			ConfessionCooldown: 10 * time.Minute,
			RoastCooldown:      5 * time.Minute,
		},
		Features: FeatureFlags{
			News:        true,
			HelpSystem:  true,
			Chatbot:     hasAnyAIKey(),
			Confession:  getEnvAsBool("ENABLE_CONFESSION", true),
			Roast:       getEnvAsBool("ENABLE_ROAST", true),
			WhaleAlerts: os.Getenv("WHALE_ALERT_API_KEY") != "",
			Reminder:    getEnvAsBool("ENABLE_REMINDER", true),
		},
		Discord: DiscordConfig{
			Colors: map[string]int{
				"primary":    0x0099FF,
				"secondary":  0x6C757D,
				"success":    0x00FF00,
				"error":      0xFF0000,
				"warning":    0xFFA500,
				"info":       0x0099FF,
				},
			SyncCommandsOnReady:  true,
			SyncCommandsGlobally: true,
		},
		AI: AIConfig{
			DeepSeekKey: os.Getenv("DEEPSEEK_API_KEY"),
			OllamaURL:   getEnvOrDefault("OLLAMA_URL", ""),
		},
		Crypto: CryptoConfig{
			WhaleAlertAPIKey:    os.Getenv("WHALE_ALERT_API_KEY"),
			TwitterAPIKey:       os.Getenv("TWITTER_API_KEY"),
			TwitterAPISecret:    os.Getenv("TWITTER_API_SECRET"),
			TwitterAccessToken:  os.Getenv("TWITTER_ACCESS_TOKEN"),
			TwitterAccessSecret: os.Getenv("TWITTER_ACCESS_SECRET"),
		},
		Reminder: ReminderConfig{
			ChannelID: os.Getenv("REMINDER_CHANNEL_ID"),
		},
		Mongo: MongoConfig{
			URL:      os.Getenv("MONGO_URL"),
			Database: getEnvOrDefault("MONGO_DB", "nerufuyo"),
		},
		Redis: RedisConfig{
			URL: os.Getenv("REDIS_URL"),
		},
		BackendURL: getEnvOrDefault("BACKEND_URL", "https://api.nerufuyo-workspace.com"),
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
	return os.Getenv("DEEPSEEK_API_KEY") != "" || os.Getenv("OLLAMA_URL") != ""
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Bot.Token == "" {
		return fmt.Errorf("Discord bot token is required")
	}
	return nil
}
