package entity

import (
	"time"
)

// NewsArticle represents a news article
type NewsArticle struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"image_url"`
	Source      string    `json:"source"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	FetchedAt   time.Time `json:"fetched_at"`
}

// NewsSource represents a news source configuration
type NewsSource struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Enabled     bool     `json:"enabled"`
	Category    string   `json:"category"`
	Icon        string   `json:"icon"`
	Color       int      `json:"color"`
	Priority    int      `json:"priority"`
	Languages   []string `json:"languages"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GuildNewsSettings holds news settings for a guild
type GuildNewsSettings struct {
	GuildID           string        `json:"guild_id"`
	Enabled           bool          `json:"enabled"`
	ChannelID         string        `json:"channel_id"`
	UpdateInterval    time.Duration `json:"update_interval"`
	MaxArticles       int           `json:"max_articles"`
	EnabledSources    []string      `json:"enabled_sources"`
	EnabledCategories []string      `json:"enabled_categories"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// NewNewsArticle creates a new NewsArticle instance
func NewNewsArticle(title, url, source string) *NewsArticle {
	return &NewsArticle{
		Title:       title,
		URL:         url,
		Source:      source,
		Tags:        make([]string, 0),
		FetchedAt:   time.Now(),
	}
}

// NewNewsSource creates a new NewsSource instance
func NewNewsSource(id, name, url string) *NewsSource {
	return &NewsSource{
		ID:        id,
		Name:      name,
		URL:       url,
		Enabled:   true,
		Priority:  0,
		Languages: []string{"en"},
		UpdatedAt: time.Now(),
	}
}

// NewGuildNewsSettings creates new settings with defaults
func NewGuildNewsSettings(guildID, channelID string) *GuildNewsSettings {
	return &GuildNewsSettings{
		GuildID:           guildID,
		Enabled:           true,
		ChannelID:         channelID,
		UpdateInterval:    10 * time.Minute,
		MaxArticles:       5,
		EnabledSources:    make([]string, 0),
		EnabledCategories: make([]string, 0),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}

// IsSourceEnabled checks if a source is enabled
func (s *GuildNewsSettings) IsSourceEnabled(sourceID string) bool {
	// If no sources specified, all are enabled
	if len(s.EnabledSources) == 0 {
		return true
	}
	
	for _, id := range s.EnabledSources {
		if id == sourceID {
			return true
		}
	}
	return false
}

// IsCategoryEnabled checks if a category is enabled
func (s *GuildNewsSettings) IsCategoryEnabled(category string) bool {
	// If no categories specified, all are enabled
	if len(s.EnabledCategories) == 0 {
		return true
	}
	
	for _, cat := range s.EnabledCategories {
		if cat == category {
			return true
		}
	}
	return false
}
