package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// handleChat handles the chatbot command
func (b *Bot) handleChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.chatbotService == nil {
		b.respondError(s, i, "Chatbot service is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a message")
		return
	}

	message := options[0].StringValue()

	// Defer response since AI might take time
	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Get AI response
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := b.chatbotService.Chat(ctx, i.Member.User.ID, message)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("AI error: %s", err.Error()))
		return
	}

	// Send response
	embed := &discordgo.MessageEmbed{
		Title:       "AI Response",
		Description: response,
		Color:       0x00ff00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Providers: %v", b.chatbotService.GetAvailableProviders()),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

// handleChatReset handles clearing chat history
func (b *Bot) handleChatReset(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.chatbotService == nil {
		b.respondError(s, i, "Chatbot service is not available")
		return
	}

	b.chatbotService.ResetSession(i.Member.User.ID)
	b.respond(s, i, "Chat history cleared.")
}
