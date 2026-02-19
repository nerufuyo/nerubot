package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// handleRoast handles the roast command
func (b *Bot) handleRoast(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Determine target user
	targetUser := i.Member.User
	options := i.ApplicationCommandData().Options
	if len(options) > 0 && options[0].Type == discordgo.ApplicationCommandOptionUser {
		targetUser = options[0].UserValue(s)
	}

	if targetUser.Bot {
		b.respondError(s, i, "I don't roast bots, they're already broken enough!")
		return
	}

	ctx := context.Background()
	roast, err := b.roastService.GenerateRoast(ctx, targetUser.ID, i.GuildID, targetUser.Username)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Roast",
		Description: roast,
		Color:       config.ColorError,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
	}

	b.respondEmbed(s, i, embed)
}
