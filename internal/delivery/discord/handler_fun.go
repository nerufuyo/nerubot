package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// handleDadJoke handles the /dadjoke command.
func (b *Bot) handleDadJoke(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.funService == nil {
		b.respondError(s, i, "Fun service is not available")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	joke, err := b.funService.FetchDadJoke()
	if err != nil {
		b.followUpError(s, i, "Failed to fetch a dad joke: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🤣 Dad Joke",
		Description: joke.Punchline,
		Color:       0xFFD700,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Powered by icanhazdadjoke.com",
		},
	}

	b.followUpEmbed(s, i, embed)
}

// handleDadJokeSetup handles the /dadjoke-setup command (admin only).
func (b *Bot) handleDadJokeSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.funService == nil {
		b.respondError(s, i, "Fun service is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		b.respondError(s, i, "Please specify a channel and interval.")
		return
	}

	channel := options[0].ChannelValue(s)
	if channel == nil || channel.Type != discordgo.ChannelTypeGuildText {
		b.respondError(s, i, "Please select a valid text channel.")
		return
	}

	interval := options[1].IntValue()
	if interval < 0 {
		b.respondError(s, i, "Interval must be 0 (disabled) or a positive number of minutes.")
		return
	}

	// Get guild name
	guildName := i.GuildID
	if guild, err := s.Guild(i.GuildID); err == nil {
		guildName = guild.Name
	}

	// Load or create config
	cfg, err := b.funService.GetGuildConfig(i.GuildID, guildName)
	if err != nil {
		b.respondError(s, i, "Failed to load guild config: "+err.Error())
		return
	}

	cfg.DadJokeChannelID = channel.ID
	cfg.DadJokeInterval = int(interval)

	if err := b.funService.SaveGuildConfig(cfg); err != nil {
		b.respondError(s, i, "Failed to save settings: "+err.Error())
		return
	}

	if interval == 0 {
		b.respond(s, i, fmt.Sprintf("Dad jokes scheduled posting to <#%s> has been **disabled**.", channel.ID))
	} else {
		b.respond(s, i, fmt.Sprintf("Dad jokes will be posted to <#%s> every **%d minutes**! 🤣", channel.ID, interval))
	}
}

// handleMeme handles the /meme command.
func (b *Bot) handleMeme(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.funService == nil {
		b.respondError(s, i, "Fun service is not available")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	meme, err := b.funService.FetchMeme()
	if err != nil {
		b.followUpError(s, i, "Failed to fetch a meme: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:     "😂 " + meme.Title,
		URL:       meme.PostLink,
		Color:     0xFF4500,
		Timestamp: time.Now().Format(time.RFC3339),
		Image: &discordgo.MessageEmbedImage{
			URL: meme.URL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("r/%s • by u/%s", meme.Subreddit, meme.Author),
		},
	}

	b.followUpEmbed(s, i, embed)
}

// handleMemeSetup handles the /meme-setup command (admin only).
func (b *Bot) handleMemeSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.funService == nil {
		b.respondError(s, i, "Fun service is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		b.respondError(s, i, "Please specify a channel and interval.")
		return
	}

	channel := options[0].ChannelValue(s)
	if channel == nil || channel.Type != discordgo.ChannelTypeGuildText {
		b.respondError(s, i, "Please select a valid text channel.")
		return
	}

	interval := options[1].IntValue()
	if interval < 0 {
		b.respondError(s, i, "Interval must be 0 (disabled) or a positive number of minutes.")
		return
	}

	// Get guild name
	guildName := i.GuildID
	if guild, err := s.Guild(i.GuildID); err == nil {
		guildName = guild.Name
	}

	// Load or create config
	cfg, err := b.funService.GetGuildConfig(i.GuildID, guildName)
	if err != nil {
		b.respondError(s, i, "Failed to load guild config: "+err.Error())
		return
	}

	cfg.MemeChannelID = channel.ID
	cfg.MemeInterval = int(interval)

	if err := b.funService.SaveGuildConfig(cfg); err != nil {
		b.respondError(s, i, "Failed to save settings: "+err.Error())
		return
	}

	if interval == 0 {
		b.respond(s, i, fmt.Sprintf("Meme scheduled posting to <#%s> has been **disabled**.", channel.ID))
	} else {
		b.respond(s, i, fmt.Sprintf("Memes will be posted to <#%s> every **%d minutes**! 😂", channel.ID, interval))
	}
}

// loadSavedGuildConfigs loads all saved guild configs from DB and restores
// reminder channels and other persisted settings on startup.
func (b *Bot) loadSavedGuildConfigs() {
	repo := repository.NewGuildConfigRepository()
	configs, err := repo.GetAll()
	if err != nil {
		b.logger.Warn("Failed to load guild configs from DB", "error", err)
		return
	}

	for _, cfg := range configs {
		// Restore reminder channel if not already set via env and we have it saved
		if b.reminderService != nil && cfg.ReminderChannelID != "" {
			currentCh := b.reminderService.GetChannelID()
			if currentCh == "" {
				b.reminderService.SetChannelID(cfg.ReminderChannelID)
				b.logger.Info("Restored reminder channel from DB",
					"guild", cfg.GuildID,
					"channel", cfg.ReminderChannelID,
				)
			}
		}
	}

	b.logger.Info("Guild configs loaded from database", "count", len(configs))
}

// persistReminderChannel saves the reminder channel to guild config in DB
// so it survives redeployments.
func (b *Bot) persistReminderChannel(guildID, guildName, channelID string) {
	repo := repository.NewGuildConfigRepository()
	cfg, err := repo.Get(guildID)
	if err != nil {
		b.logger.Warn("Failed to get guild config for reminder persist", "error", err)
		return
	}
	if cfg == nil {
		cfg = entity.NewGuildConfig(guildID, guildName)
	}
	cfg.ReminderChannelID = channelID
	if err := repo.Save(cfg); err != nil {
		b.logger.Warn("Failed to persist reminder channel", "error", err)
	}
}
