package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// handleRoast handles the roast command
func (b *Bot) handleRoast(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if roast is enabled from dashboard
	if b.backendClient != nil && !b.backendClient.GetSettings().Features.RoastEnabled {
		b.respondError(s, i, "Roast feature is currently disabled by the admin.")
		return
	}

	// Determine target user and language
	targetUser := i.Member.User
	lang := config.DefaultLang
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		switch opt.Name {
		case "user":
			targetUser = opt.UserValue(s)
		case "lang":
			lang = opt.StringValue()
		}
	}

	if targetUser.Bot {
		b.respondError(s, i, "I don't roast bots, they're already broken enough!")
		return
	}

	ctx := context.Background()
	roast, err := b.roastService.GenerateRoast(ctx, targetUser.ID, i.GuildID, targetUser.Username, lang)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	langLabel := config.LanguageNames[lang]
	if langLabel == "" {
		langLabel = "English"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Roast",
		Description: roast,
		Color:       config.ColorError,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s | Lang: %s", i.Member.User.Username, langLabel),
		},
	}

	b.respondEmbed(s, i, embed)
}
