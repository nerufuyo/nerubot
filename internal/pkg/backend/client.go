package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// SettingsChangeFunc is called whenever bot settings are refreshed from the backend.
type SettingsChangeFunc func(settings *BotSettings)

// Client communicates with the nerufuyo-workspace-backend API
// to retrieve bot settings.
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger

	// Cached bot settings
	settings     *BotSettings
	settingsMu   sync.RWMutex
	settingsLast time.Time
	knownVersion int // last known settingsVersion from backend

	// Callback when settings change
	onSettingsChange SettingsChangeFunc
}

// BotSettings represents the configurable bot settings from the dashboard.
type BotSettings struct {
	ID              string      `json:"id"`
	BotName         string      `json:"botName"`
	BotStatus       string      `json:"botStatus"`
	BotDescription  string      `json:"botDescription"` // About Me / profile bio
	SystemPrompt    string      `json:"systemPrompt"`
	RateLimitCount  int         `json:"rateLimitCount"`  // max messages per window
	RateLimitWindow int         `json:"rateLimitWindow"` // window in seconds
	MaxTokens       int         `json:"maxTokens"`
	Temperature     float64     `json:"temperature"`
	Features        BotFeatures `json:"features"`
	WelcomeMessage  string      `json:"welcomeMessage"`
	SettingsVersion int         `json:"settingsVersion"`

	// Per-feature advanced settings
	ChatSettings       ChatFeatureSettings       `json:"chatSettings"`
	RoastSettings      RoastFeatureSettings      `json:"roastSettings"`
	ConfessionSettings ConfessionFeatureSettings `json:"confessionSettings"`
	NewsSettings       NewsFeatureSettings       `json:"newsSettings"`
	WhaleSettings      WhaleFeatureSettings      `json:"whaleSettings"`
	ReminderSettings   ReminderFeatureSettings   `json:"reminderSettings"`

	UpdatedAt string `json:"updatedAt"`
}

// BotFeatures toggles for bot features managed from dashboard.
type BotFeatures struct {
	ChatEnabled       bool `json:"chatEnabled"`
	RoastEnabled      bool `json:"roastEnabled"`
	ConfessionEnabled bool `json:"confessionEnabled"`
	NewsEnabled       bool `json:"newsEnabled"`
	WhaleEnabled      bool `json:"whaleEnabled"`
	ReminderEnabled   bool `json:"reminderEnabled"`
}

// ChatFeatureSettings holds advanced settings for AI chat.
type ChatFeatureSettings struct {
	MaxHistoryMessages int    `json:"maxHistoryMessages"`
	SessionTimeoutMins int    `json:"sessionTimeoutMins"`
	ModelName          string `json:"modelName"`
}

// RoastFeatureSettings holds advanced settings for roast.
type RoastFeatureSettings struct {
	CooldownMinutes int `json:"cooldownMinutes"`
	MinMessages     int `json:"minMessages"`
	MaxSeverity     int `json:"maxSeverity"`
}

// ConfessionFeatureSettings holds advanced settings for confessions.
type ConfessionFeatureSettings struct {
	CooldownMinutes int  `json:"cooldownMinutes"`
	MaxLength       int  `json:"maxLength"`
	RequireApproval bool `json:"requireApproval"`
	AllowImages     bool `json:"allowImages"`
	AllowReplies    bool `json:"allowReplies"`
}

// NewsFeatureSettings holds advanced settings for news.
type NewsFeatureSettings struct {
	MaxArticles    int      `json:"maxArticles"`
	UpdateInterval int      `json:"updateInterval"`
	Sources        []string `json:"sources"`
}

// WhaleFeatureSettings holds advanced settings for whale alerts.
type WhaleFeatureSettings struct {
	MinAmountUSD float64  `json:"minAmountUSD"`
	Blockchains  []string `json:"blockchains"`
	Symbols      []string `json:"symbols"`
}

// ReminderFeatureSettings holds advanced settings for reminders.
type ReminderFeatureSettings struct {
	Timezone  string `json:"timezone"`
	ChannelID string `json:"channelID"`
}

// New creates a new backend API client.
func New(baseURL string) *Client {
	log := logger.New("backend-client")

	// Normalize base URL
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.nerufuyo-workspace.com"
	}

	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: log,
	}

	// Initial fetch
	go c.refreshSettings()

	// Background: fast-poll version every 30s
	go c.backgroundRefresh()

	return c
}

// OnSettingsChange registers a callback that fires whenever bot settings are
// refreshed from the backend. The bot uses this to sync Discord status, etc.
func (c *Client) OnSettingsChange(fn SettingsChangeFunc) {
	c.onSettingsChange = fn
}

// GetSettings returns the cached bot settings.
func (c *Client) GetSettings() *BotSettings {
	c.settingsMu.RLock()
	defer c.settingsMu.RUnlock()
	if c.settings != nil {
		return c.settings
	}
	// Return defaults
	return &BotSettings{
		BotName:         "NeruBot",
		BotStatus:       "Ready to rock your server!",
		BotDescription:  "Meet NERU — your all-in-one Discord buddy! From news to confessions and messages, NERU's got your back.",
		RateLimitCount:  5,
		RateLimitWindow: 180, // 3 minutes
		MaxTokens:       1024,
		Temperature:     0.7,
		Features: BotFeatures{
			ChatEnabled:       true,
			RoastEnabled:      true,
			ConfessionEnabled: true,

			NewsEnabled:     true,
			WhaleEnabled:    true,
			ReminderEnabled: true,
		},
		ChatSettings: ChatFeatureSettings{
			MaxHistoryMessages: 10,
			SessionTimeoutMins: 30,
			ModelName:          "deepseek-chat",
		},
		RoastSettings: RoastFeatureSettings{
			CooldownMinutes: 5,
			MinMessages:     10,
			MaxSeverity:     3,
		},
		ConfessionSettings: ConfessionFeatureSettings{
			CooldownMinutes: 10,
			MaxLength:       2000,
			RequireApproval: false,
			AllowImages:     true,
			AllowReplies:    true,
		},
		NewsSettings: NewsFeatureSettings{
			MaxArticles:    10,
			UpdateInterval: 10,
			Sources:        []string{"BBC News", "CNN", "Reuters", "TechCrunch", "The Verge"},
		},
		WhaleSettings: WhaleFeatureSettings{
			MinAmountUSD: 1000000,
			Blockchains:  []string{"bitcoin", "ethereum"},
			Symbols:      []string{"BTC", "ETH"},
		},
		ReminderSettings: ReminderFeatureSettings{
			Timezone:  "Asia/Jakarta",
			ChannelID: "",
		},
		WelcomeMessage: "Hiii~! ✨ Neru is here! Ask me anything — tech stuff, cool ideas, or just chat! Neru's ready to help! 🎉",
	}
}

// backgroundRefresh fast-polls the settings version every 30s.
func (c *Client) backgroundRefresh() {
	versionTicker := time.NewTicker(30 * time.Second)
	defer versionTicker.Stop()

	for {
		select {
		case <-versionTicker.C:
			c.checkVersionAndRefresh()
		}
	}
}

// checkVersionAndRefresh polls the lightweight /bot/settings/version endpoint.
// If the version changed, it does a full settings refresh.
func (c *Client) checkVersionAndRefresh() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var resp struct {
		Version int `json:"version"`
	}
	if err := c.get(ctx, "/api/v1/bot/settings/version", &resp); err != nil {
		// Fall back silently; full refresh will happen on the next cycle
		return
	}

	if resp.Version != c.knownVersion {
		c.logger.Info("Settings version changed, refreshing",
			"old", c.knownVersion, "new", resp.Version)
		c.refreshSettings()
	}
}

// refreshSettings fetches bot settings from the backend.
func (c *Client) refreshSettings() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	settings, err := c.fetchBotSettings(ctx)
	if err != nil {
		c.logger.Warn("Failed to fetch bot settings from backend", "error", err)
		return
	}

	c.settingsMu.Lock()
	c.settings = settings
	c.settingsLast = time.Now()
	c.knownVersion = settings.SettingsVersion
	c.settingsMu.Unlock()

	c.logger.Info("Bot settings refreshed from backend",
		"version", settings.SettingsVersion)

	// Notify listener (e.g. bot syncs Discord status)
	if c.onSettingsChange != nil {
		c.onSettingsChange(settings)
	}
}

// --- API fetch methods ---

func (c *Client) fetchBotSettings(ctx context.Context) (*BotSettings, error) {
	var resp struct {
		Data BotSettings `json:"data"`
	}
	if err := c.get(ctx, "/api/v1/bot/settings", &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// get performs an HTTP GET request and decodes the JSON response.
func (c *Client) get(ctx context.Context, path string, dest interface{}) error {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http get %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http get %s: status %d: %s", path, resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response from %s: %w", path, err)
	}

	return nil
}

// post performs an HTTP POST request with a JSON body.
func (c *Client) post(ctx context.Context, path string, body interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal body: %w", err)
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http post %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http post %s: status %d: %s", path, resp.StatusCode, string(respBody))
	}
	return nil
}

// GuildInfo represents a Discord guild for syncing with the backend.
type GuildInfo struct {
	GuildID     string `json:"guildId"`
	GuildName   string `json:"guildName"`
	MemberCount int    `json:"memberCount"`
	IconURL     string `json:"iconUrl"`
}

// SyncGuilds pushes the bot's current guild list to the backend.
func (c *Client) SyncGuilds(guilds []GuildInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payload := struct {
		Guilds []GuildInfo `json:"guilds"`
	}{Guilds: guilds}

	if err := c.post(ctx, "/api/v1/bot/guilds/sync", payload); err != nil {
		c.logger.Warn("Failed to sync guilds to backend", "error", err)
		return
	}
	c.logger.Info("Guilds synced to backend", "count", len(guilds))
}
