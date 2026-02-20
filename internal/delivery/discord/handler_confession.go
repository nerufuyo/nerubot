package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// handleConfess handles the confess command
func (b *Bot) handleConfess(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if confessions are enabled from dashboard
	if b.backendClient != nil && !b.backendClient.GetSettings().Features.ConfessionEnabled {
		b.respondError(s, i, "Confession feature is currently disabled by the admin.")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide confession content")
		return
	}

	content := options[0].StringValue()

	ctx := context.Background()
	confession, err := b.confessionService.SubmitConfession(ctx, i.GuildID, i.Member.User.ID, content)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	// Respond privately (ephemeral)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Your confession has been submitted anonymously. (ID: #%d)", confession.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond", "error", err)
	}
}
