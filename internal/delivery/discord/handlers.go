package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// handleHelp handles the help command.
func (b *Bot) handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       b.config.Bot.Name + " Help",
		Description: b.config.Bot.Description,
		Color:       config.ColorPrimary,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Music Commands",
				Value: "`/play <query>` - Play a song or add to queue\n" +
					"`/skip` - Skip current song\n" +
					"`/stop` - Stop playback and clear queue\n" +
					"`/queue` - Show music queue",
				Inline: false,
			},
			{
				Name:   "Confession Commands",
				Value:  "`/confess <content>` - Submit an anonymous confession",
				Inline: false,
			},
			{
				Name:   "Roast Commands",
				Value:  "`/roast [user]` - Get roasted based on Discord activity",
				Inline: false,
			},
			{
				Name: "AI Chatbot Commands",
				Value: "`/chat <message>` - Chat with AI\n" +
					"`/chat-reset` - Reset your chat history",
				Inline: false,
			},
			{
				Name:   "News Commands",
				Value:  "`/news` - Get latest news from multiple sources",
				Inline: false,
			},
			{
				Name:   "Whale Alert Commands",
				Value:  "`/whale` - Get recent whale cryptocurrency transactions",
				Inline: false,
			},
			{
				Name: "Analytics Commands",
				Value: "`/stats` - View server statistics\n" +
					"`/profile [user]` - View user profile",
				Inline: false,
			},
			{
				Name: "Reminder Commands",
				Value: "`/reminder` - View upcoming holidays and Ramadan schedule\n" +
					"`/reminder-set <channel>` - Set reminder channel (admin only)",
				Inline: false,
			},
			{
				Name:   "Other Commands",
				Value:  "`/help` - Show this help message",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s v%s | %s", b.config.Bot.Name, b.config.Bot.Version, b.config.Bot.Author),
		},
	}

	b.respondEmbed(s, i, embed)
}

// --- Response helpers ---

// deferResponse sends a deferred response to the interaction.
func (b *Bot) deferResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
	}
	return err
}

// respond sends a text response to the interaction.
func (b *Bot) respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond to interaction", "error", err)
	}
}

// respondEmbed sends an embed response to the interaction.
func (b *Bot) respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond to interaction", "error", err)
	}
}

// respondError sends an error text response to the interaction.
func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	b.respond(s, i, config.EmojiError+" "+message)
}

// followUp sends a follow-up message after a deferred response.
func (b *Bot) followUp(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

// followUpEmbed sends a follow-up embed message after a deferred response.
func (b *Bot) followUpEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up embed", "error", err)
	}
}

// followUpError sends an error follow-up message.
func (b *Bot) followUpError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	b.followUp(s, i, config.EmojiError+" "+message)
}

// findUserVoiceState returns the voice state of a user in a guild.
func (b *Bot) findUserVoiceState(guildID, userID string) *discordgo.VoiceState {
	guild, err := b.session.State.Guild(guildID)
	if err != nil {
		return nil
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs
		}
	}

	return nil
}
