package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/entity"
)

// handlePlay handles the play command
func (b *Bot) handlePlay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a song name or URL")
		return
	}

	query := options[0].StringValue()

	// Check if user is in voice channel
	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You must be in a voice channel!")
		return
	}

	// Defer response since search might take time
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
		return
	}

	// Search and add song
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	song, err := b.musicService.AddSong(ctx, i.GuildID, vs.ChannelID, i.ChannelID, query, i.Member.User.ID)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to add song: %s", err.Error()))
		return
	}

	// Get queue to check position
	queue := b.musicService.GetQueue(i.GuildID, vs.ChannelID, i.ChannelID)

	var message string
	if queue.Size() == 1 {
		// First song, start playing
		if err := b.musicService.Play(i.GuildID); err != nil {
			b.followUpError(s, i, fmt.Sprintf("Failed to start playback: %s", err.Error()))
			return
		}
		message = fmt.Sprintf("🎵 Now playing: **%s** by **%s**", song.Title, song.Artist)
	} else {
		message = fmt.Sprintf("➕ Added to queue: **%s** by **%s** (Position: #%d)", 
			song.Title, song.Artist, queue.Size())
	}

	b.followUp(s, i, message)
}

// handleSkip handles the skip command
func (b *Bot) handleSkip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	song, err := b.musicService.Skip(i.GuildID)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	if song != nil {
		b.respond(s, i, fmt.Sprintf("⏭️ Skipped to: **%s**", song.Title))
	} else {
		b.respond(s, i, "⏭️ Skipped. Queue is now empty.")
	}
}

// handleStop handles the stop command
func (b *Bot) handleStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := b.musicService.Stop(i.GuildID); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, "⏹️ Stopped playback and cleared queue")
}

// handleQueue handles the queue command
func (b *Bot) handleQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	queue := b.musicService.GetQueue(i.GuildID, "", "")
	
	if queue.IsEmpty() {
		b.respond(s, i, "📭 Queue is empty")
		return
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title: "🎵 Music Queue",
		Color: config.ColorMusic,
		Fields: []*discordgo.MessageEmbedField{},
	}

	// Add current song
	if current := queue.Current(); current != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Now Playing",
			Value:  fmt.Sprintf("**%s** by **%s**\n%s", current.Title, current.Artist, entity.FormatDuration(current.Duration)),
			Inline: false,
		})
	}

	// Add queue
	if queue.Remaining() > 0 {
		queueText := ""
		for i := queue.CurrentIndex + 1; i < queue.Size() && i < queue.CurrentIndex+6; i++ {
			song := queue.Songs[i]
			queueText += fmt.Sprintf("%d. **%s** - %s\n", i-queue.CurrentIndex, song.Title, entity.FormatDuration(song.Duration))
		}
		
		if queue.Remaining() > 5 {
			queueText += fmt.Sprintf("\n...and %d more songs", queue.Remaining()-5)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Up Next",
			Value:  queueText,
			Inline: false,
		})
	}

	// Add footer
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Total: %d songs | Duration: %s | Loop: %s",
			queue.Size(),
			entity.FormatDuration(queue.TotalDuration()),
			queue.LoopMode,
		),
	}

	b.respondEmbed(s, i, embed)
}

// handleConfess handles the confess command
func (b *Bot) handleConfess(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	// Respond privately
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Your confession has been submitted anonymously! (ID: #%d)", confession.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond", "error", err)
	}
}

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
		Title:       "🔥 Roast",
		Description: roast,
		Color:       config.ColorError,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
	}

	b.respondEmbed(s, i, embed)
}

// handleHelp handles the help command
func (b *Bot) handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "📚 " + b.config.Bot.Name + " Help",
		Description: b.config.Bot.Description,
		Color:       config.ColorPrimary,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "🎵 Music Commands",
				Value: "`/play <query>` - Play a song or add to queue\n" +
					"`/skip` - Skip current song\n" +
					"`/stop` - Stop playback and clear queue\n" +
					"`/queue` - Show music queue",
				Inline: false,
			},
			{
				Name: "📝 Confession Commands",
				Value: "`/confess <content>` - Submit an anonymous confession",
				Inline: false,
			},
			{
				Name: "🔥 Roast Commands",
				Value: "`/roast [user]` - Get roasted based on Discord activity",
				Inline: false,
			},
			{
				Name: "ℹ️ Other Commands",
				Value: "`/help` - Show this help message",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s v%s | %s", b.config.Bot.Name, b.config.Bot.Version, b.config.Bot.Author),
		},
	}

	b.respondEmbed(s, i, embed)
}

// Helper methods

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

func (b *Bot) followUp(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

func (b *Bot) followUpError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	b.followUp(s, i, "❌ "+message)
}
