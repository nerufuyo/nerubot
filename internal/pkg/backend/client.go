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
// to retrieve RAG knowledge, projects, articles, experiences, and bot settings.
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger

	// Cached RAG context (refreshed periodically)
	ragContext     string
	ragMu          sync.RWMutex
	ragLastRefresh time.Time

	// Cached bot settings
	settings       *BotSettings
	settingsMu     sync.RWMutex
	settingsLast   time.Time
	knownVersion   int // last known settingsVersion from backend

	// Callback when settings change
	onSettingsChange SettingsChangeFunc
}

// BotSettings represents the configurable bot settings from the dashboard.
type BotSettings struct {
	ID              string           `json:"id"`
	BotName         string           `json:"botName"`
	BotStatus       string           `json:"botStatus"`
	BotDescription  string           `json:"botDescription"`  // About Me / profile bio
	SystemPrompt    string           `json:"systemPrompt"`
	RateLimitCount  int              `json:"rateLimitCount"`  // max messages per window
	RateLimitWindow int              `json:"rateLimitWindow"` // window in seconds
	MaxTokens       int              `json:"maxTokens"`
	Temperature     float64          `json:"temperature"`
	Features        BotFeatures      `json:"features"`
	WelcomeMessage  string           `json:"welcomeMessage"`
	SettingsVersion int              `json:"settingsVersion"`

	// Per-feature advanced settings
	ChatSettings       ChatFeatureSettings       `json:"chatSettings"`
	RoastSettings      RoastFeatureSettings      `json:"roastSettings"`
	ConfessionSettings ConfessionFeatureSettings `json:"confessionSettings"`
	NewsSettings       NewsFeatureSettings       `json:"newsSettings"`
	WhaleSettings      WhaleFeatureSettings      `json:"whaleSettings"`
	ReminderSettings   ReminderFeatureSettings   `json:"reminderSettings"`

	UpdatedAt       string           `json:"updatedAt"`
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
	EnableRAG          bool   `json:"enableRAG"`
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

// Experience from backend.
type Experience struct {
	Company     string `json:"company"`
	Role        string `json:"role"`
	Description string `json:"description"`
	Place       string `json:"place"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

// Project from backend.
type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	RepoURL     string   `json:"repoUrl"`
	LiveURL     string   `json:"liveUrl"`
	Tags        []string `json:"tags"`
	Stars       int      `json:"stars"`
	IsTop       bool     `json:"isTop"`
	Images      []string `json:"images"`
}

// Article from backend.
type Article struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"isPublished"`
}

// Knowledge from backend.
type Knowledge struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	Category  string `json:"category"`
	Content   string `json:"content"`
	Summary   string `json:"summary"`
	SourceURL string `json:"sourceUrl"`
	FileURL   string `json:"fileUrl"`
	IsActive  bool   `json:"isActive"`
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
	go c.refreshRAGContext()
	go c.refreshSettings()

	// Background: fast-poll version every 30s, full RAG refresh every 5 min
	go c.backgroundRefresh()

	return c
}

// OnSettingsChange registers a callback that fires whenever bot settings are
// refreshed from the backend. The bot uses this to sync Discord status, etc.
func (c *Client) OnSettingsChange(fn SettingsChangeFunc) {
	c.onSettingsChange = fn
}

// GetRAGContext returns the cached RAG context string for AI system prompts.
func (c *Client) GetRAGContext() string {
	c.ragMu.RLock()
	defer c.ragMu.RUnlock()
	return c.ragContext
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

			NewsEnabled:       true,
			WhaleEnabled:      true,
			ReminderEnabled:   true,
		},
		ChatSettings: ChatFeatureSettings{
			MaxHistoryMessages: 10,
			SessionTimeoutMins: 30,
			ModelName:          "deepseek-chat",
			EnableRAG:          true,
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
		WelcomeMessage: "Hey there! I'm Neru, your friendly AI companion. Ask me anything about Nerufuyo!",
	}
}

// backgroundRefresh fast-polls the settings version every 30s and does a
// full RAG context refresh every 5 minutes.
func (c *Client) backgroundRefresh() {
	versionTicker := time.NewTicker(30 * time.Second)
	ragTicker := time.NewTicker(5 * time.Minute)
	defer versionTicker.Stop()
	defer ragTicker.Stop()

	for {
		select {
		case <-versionTicker.C:
			c.checkVersionAndRefresh()
		case <-ragTicker.C:
			c.refreshRAGContext()
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

// refreshRAGContext fetches all knowledge data from the backend and builds a context string.
func (c *Client) refreshRAGContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var sections []string

	// Fetch knowledge base entries
	if knowledge, err := c.fetchKnowledge(ctx); err == nil && len(knowledge) > 0 {
		var lines []string
		lines = append(lines, "=== KNOWLEDGE BASE (from RAG database) ===")
		for _, k := range knowledge {
			if !k.IsActive {
				continue
			}
			line := fmt.Sprintf("--- %s | %s | %s ---", k.Type, k.Category, k.Title)
			if k.Content != "" {
				line += "\n" + k.Content
			}
			if k.Summary != "" {
				line += "\nSummary: " + k.Summary
			}
			if k.SourceURL != "" {
				line += "\nSource: " + k.SourceURL
			}
			lines = append(lines, line)
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}

	// Fetch experiences
	if experiences, err := c.fetchExperiences(ctx); err == nil && len(experiences) > 0 {
		var lines []string
		lines = append(lines, "=== WORK EXPERIENCE (from database) ===")
		for _, exp := range experiences {
			endDate := exp.EndDate
			if endDate == "" {
				endDate = "Present"
			}
			line := fmt.Sprintf("- %s at %s (%s) | %s – %s", exp.Role, exp.Company, exp.Place, exp.StartDate, endDate)
			if exp.Description != "" {
				line += "\n  Description: " + exp.Description
			}
			lines = append(lines, line)
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}

	// Fetch projects
	if projects, err := c.fetchProjects(ctx); err == nil && len(projects) > 0 {
		var lines []string
		lines = append(lines, "=== PROJECTS (from database) ===")
		for _, proj := range projects {
			line := fmt.Sprintf("- %s: %s", proj.Title, proj.Description)
			if len(proj.Tags) > 0 {
				line += " | Tech: " + strings.Join(proj.Tags, ", ")
			}
			if proj.RepoURL != "" {
				line += " | GitHub: " + proj.RepoURL
			}
			if proj.LiveURL != "" {
				line += " | Live: " + proj.LiveURL
			}
			if proj.Stars > 0 {
				line += fmt.Sprintf(" | ⭐ %d stars", proj.Stars)
			}
			if proj.IsTop {
				line += " | [TOP PROJECT]"
			}
			lines = append(lines, line)
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}

	// Fetch articles
	if articles, err := c.fetchArticles(ctx); err == nil && len(articles) > 0 {
		var lines []string
		lines = append(lines, "=== ARTICLES & BLOG POSTS (from database) ===")
		for _, art := range articles {
			if !art.IsPublished {
				continue
			}
			line := fmt.Sprintf("- \"%s\"", art.Title)
			if len(art.Tags) > 0 {
				line += " | Topics: " + strings.Join(art.Tags, ", ")
			}
			if art.Content != "" {
				excerpt := art.Content
				if len(excerpt) > 500 {
					excerpt = excerpt[:500] + "..."
				}
				line += "\n  Content: " + excerpt
			}
			if art.Slug != "" {
				line += "\n  URL: https://nerufuyo-workspace.com/articles/" + art.Slug
			}
			lines = append(lines, line)
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}

	ragContext := strings.Join(sections, "\n\n")

	c.ragMu.Lock()
	c.ragContext = ragContext
	c.ragLastRefresh = time.Now()
	c.ragMu.Unlock()

	if ragContext != "" {
		c.logger.Info("RAG context refreshed", "sections", len(sections), "length", len(ragContext))
	} else {
		c.logger.Warn("RAG context is empty — backend may be unreachable")
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

func (c *Client) fetchKnowledge(ctx context.Context) ([]Knowledge, error) {
	var resp struct {
		Data []Knowledge `json:"data"`
	}
	if err := c.get(ctx, "/api/v1/bot/knowledge", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) fetchExperiences(ctx context.Context) ([]Experience, error) {
	var resp struct {
		Data []Experience `json:"data"`
	}
	if err := c.get(ctx, "/api/v1/experiences", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) fetchProjects(ctx context.Context) ([]Project, error) {
	var resp struct {
		Data []Project `json:"data"`
	}
	if err := c.get(ctx, "/api/v1/projects", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) fetchArticles(ctx context.Context) ([]Article, error) {
	var resp struct {
		Data []Article `json:"data"`
	}
	if err := c.get(ctx, "/api/v1/articles", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

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
