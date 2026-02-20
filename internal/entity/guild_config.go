package entity

import "time"

// GuildConfig holds persisted bot configuration for a specific guild.
// This ensures settings survive redeployments without requiring re-setup.
type GuildConfig struct {
	GuildID   string    `json:"guild_id" bson:"guild_id"`
	GuildName string    `json:"guild_name" bson:"guild_name"`

	// Reminder settings
	ReminderChannelID string `json:"reminder_channel_id" bson:"reminder_channel_id"`

	// Fun feature channels & scheduling
	DadJokeChannelID  string `json:"dad_joke_channel_id" bson:"dad_joke_channel_id"`
	DadJokeInterval   int    `json:"dad_joke_interval" bson:"dad_joke_interval"`     // interval in minutes (0 = disabled)
	MemeChannelID     string `json:"meme_channel_id" bson:"meme_channel_id"`
	MemeInterval      int    `json:"meme_interval" bson:"meme_interval"`             // interval in minutes (0 = disabled)

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewGuildConfig creates a new GuildConfig with defaults.
func NewGuildConfig(guildID, guildName string) *GuildConfig {
	return &GuildConfig{
		GuildID:   guildID,
		GuildName: guildName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
