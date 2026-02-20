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
				if cfg.ReminderLang != "" {
					b.reminderService.SetLang(cfg.ReminderLang)
				}
				b.logger.Info("Restored reminder channel from DB",
					"guild", cfg.GuildID,
					"channel", cfg.ReminderChannelID,
					"lang", cfg.ReminderLang,
				)
			}
		}

		// Log restored mental health config if set
		if cfg.MentalHealthChannelID != "" && cfg.MentalHealthInterval > 0 {
			b.logger.Info("Mental health reminder active from DB",
				"guild", cfg.GuildID,
				"channel", cfg.MentalHealthChannelID,
				"interval", cfg.MentalHealthInterval,
				"tag", cfg.MentalHealthTag,
				"lang", cfg.MentalHealthLang,
			)
		}
	}

	b.logger.Info("Guild configs loaded from database", "count", len(configs))
}

// persistReminderChannel saves the reminder channel to guild config in DB
// so it survives redeployments.
func (b *Bot) persistReminderChannel(guildID, guildName, channelID, lang string) {
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
	cfg.ReminderLang = lang
	if err := repo.Save(cfg); err != nil {
		b.logger.Warn("Failed to persist reminder channel", "error", err)
	}
}

// --- Mental Health Commands ---

// handleMentalHealth handles the /mentalhealth command — returns a random mental health tip.
func (b *Bot) handleMentalHealth(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.funService == nil {
		b.respondError(s, i, "Fun service is not available")
		return
	}

	// Extract language option
	lang := "EN"
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "lang" {
			lang = opt.StringValue()
		}
	}

	title, tip := b.funService.GetMentalHealthTip(lang)

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: tip,
		Color:       0x57F287, // green
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Take care of your mental health 💚",
		},
	}

	b.respondEmbed(s, i, embed)
}

// handleMentalHealthSetup handles the /mentalhealth-setup command (admin only).
func (b *Bot) handleMentalHealthSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.funService == nil {
		b.respondError(s, i, "Fun service is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		b.respondError(s, i, "Please specify a channel and interval.")
		return
	}

	// Parse required options
	var channelID string
	var interval int64
	var tag string
	var everyone bool
	var lang string

	for _, opt := range options {
		switch opt.Name {
		case "channel":
			ch := opt.ChannelValue(s)
			if ch == nil || ch.Type != discordgo.ChannelTypeGuildText {
				b.respondError(s, i, "Please select a valid text channel.")
				return
			}
			channelID = ch.ID
		case "interval":
			interval = opt.IntValue()
		case "tag":
			// Mentionable can be a user or a role — build the mention string
			// The value is a snowflake ID. We check if it's a role or user.
			mentionID := opt.Value.(string)
			// Try to find as role first
			if guild, err := s.Guild(i.GuildID); err == nil {
				isRole := false
				for _, role := range guild.Roles {
					if role.ID == mentionID {
						tag = fmt.Sprintf("<@&%s>", mentionID)
						isRole = true
						break
					}
				}
				if !isRole {
					tag = fmt.Sprintf("<@%s>", mentionID)
				}
			} else {
				tag = fmt.Sprintf("<@%s>", mentionID)
			}
		case "everyone":
			everyone = opt.BoolValue()
		case "lang":
			lang = opt.StringValue()
		}
	}

	if channelID == "" {
		b.respondError(s, i, "Please specify a channel.")
		return
	}

	if interval < 0 {
		b.respondError(s, i, "Interval must be 0 (disabled) or a positive number of minutes.")
		return
	}

	// If everyone is true, override tag
	if everyone {
		tag = "@everyone"
	}

	// Get guild name
	guildName := i.GuildID
	if guild, err := s.Guild(i.GuildID); err == nil {
		guildName = guild.Name
	}

	// Save config
	b.persistMentalHealthConfig(i.GuildID, guildName, channelID, int(interval), tag, lang)

	if interval == 0 {
		b.respond(s, i, fmt.Sprintf("Mental health reminders to <#%s> have been **disabled**.", channelID))
	} else {
		tagInfo := ""
		if tag != "" {
			tagInfo = fmt.Sprintf(" (mentioning %s)", tag)
		}
		langInfo := ""
		if lang != "" {
			langInfo = fmt.Sprintf(" in **%s**", lang)
		}
		b.respond(s, i, fmt.Sprintf("🧠 Mental health reminders will be posted to <#%s> every **%d minutes**%s%s! 💚", channelID, interval, tagInfo, langInfo))
	}
}

// handleMentalHealthStop handles the /mentalhealth-stop command.
func (b *Bot) handleMentalHealthStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.funService == nil {
		b.respondError(s, i, "Fun service is not available")
		return
	}

	guildName := i.GuildID
	if guild, err := s.Guild(i.GuildID); err == nil {
		guildName = guild.Name
	}

	b.persistMentalHealthConfig(i.GuildID, guildName, "", 0, "", "")
	b.respond(s, i, "Mental health reminders have been **stopped**. Take care! 💚")
}

// persistMentalHealthConfig saves mental health reminder config to guild config in DB.
func (b *Bot) persistMentalHealthConfig(guildID, guildName, channelID string, interval int, tag, lang string) {
	repo := repository.NewGuildConfigRepository()
	cfg, err := repo.Get(guildID)
	if err != nil {
		b.logger.Warn("Failed to get guild config for mental health persist", "error", err)
		return
	}
	if cfg == nil {
		cfg = entity.NewGuildConfig(guildID, guildName)
	}
	cfg.MentalHealthChannelID = channelID
	cfg.MentalHealthInterval = interval
	cfg.MentalHealthTag = tag
	cfg.MentalHealthLang = lang
	if err := repo.Save(cfg); err != nil {
		b.logger.Warn("Failed to persist mental health config", "error", err)
	}
}
