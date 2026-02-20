package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// handleChat handles the chatbot command with rate limiting
func (b *Bot) handleChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.chatbotService == nil {
		b.respondError(s, i, "Chatbot service is not available")
		return
	}

	// Check if chat is enabled from backend settings
	if b.backendClient != nil {
		settings := b.backendClient.GetSettings()
		if !settings.Features.ChatEnabled {
			b.respondError(s, i, "AI chat is currently disabled by the administrator")
			return
		}
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a message")
		return
	}

	message := options[0].StringValue()

	// Check rate limit (5 messages per 3 minutes per user)
	allowed, remaining, resetSeconds := b.chatbotService.CheckRateLimit(i.Member.User.ID)
	if !allowed {
		embed := &discordgo.MessageEmbed{
			Title:       "Rate Limit Exceeded",
			Description: fmt.Sprintf("You've reached the AI chat limit. Please wait **%d seconds** before sending another message.\n\nThis helps us manage API costs while keeping the bot free for everyone! 🙏", resetSeconds),
			Color:       0xFF6B6B,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Rate limit: 5 messages per 3 minutes",
			},
		}
		b.respondEmbed(s, i, embed)
		return
	}

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

	// Build footer with rate limit info and providers
	footerText := fmt.Sprintf("Powered by RAG | %d messages remaining", remaining)
	providers := b.chatbotService.GetAvailableProviders()
	if len(providers) > 0 {
		footerText += fmt.Sprintf(" | Provider: %s", providers[0])
	}

	// Send response
	embed := &discordgo.MessageEmbed{
		Title:       "🤖 Neru AI",
		Description: response,
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footerText,
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
	b.respond(s, i, "✅ Chat history cleared. Start a fresh conversation with `/chat`!")
}
