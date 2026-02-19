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
	if err := b.deferResponse(s, i); err != nil {
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
		// First song — try to join voice and start playback via Lavalink
		b.logger.Info("Music playback initiated",
			"guild", i.GuildID, "channel", vs.ChannelID)

		if err := b.musicService.Play(i.GuildID); err != nil {
			b.followUpError(s, i, fmt.Sprintf("Failed to start playback: %s", err.Error()))
			return
		}

		// If Lavalink is enabled, join voice and play
		if b.config.Lavalink.Enabled && b.lavalinkClient != nil {
			go b.playWithLavalink(i.GuildID, vs.ChannelID, song)
		}

		message = fmt.Sprintf("%s Now playing: **%s** by **%s**", config.EmojiMusic, song.Title, song.Artist)
	} else {
		message = fmt.Sprintf("Added to queue: **%s** by **%s** (Position: #%d)",
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
		b.respond(s, i, fmt.Sprintf("%s Skipped to: **%s**", config.EmojiSkip, song.Title))
	} else {
		b.respond(s, i, fmt.Sprintf("%s Skipped. Queue is now empty.", config.EmojiSkip))
	}
}

// handleStop handles the stop command
func (b *Bot) handleStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := b.musicService.Stop(i.GuildID); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, fmt.Sprintf("%s Stopped playback and cleared queue", config.EmojiStop))
}

// handleQueue handles the queue command
func (b *Bot) handleQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	queue := b.musicService.GetQueue(i.GuildID, "", "")

	if queue.IsEmpty() {
		b.respond(s, i, "Queue is empty.")
		return
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("%s Music Queue", config.EmojiMusic),
		Color:  config.ColorMusic,
		Fields: make([]*discordgo.MessageEmbedField, 0),
	}

	// Add current song
	if current := queue.Current(); current != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Now Playing",
			Value:  fmt.Sprintf("**%s** by **%s**\n%s", current.Title, current.Artist, entity.FormatDuration(current.Duration)),
			Inline: false,
		})
	}

	// Add upcoming songs
	if queue.Remaining() > 0 {
		queueText := ""
		for idx := queue.CurrentIndex + 1; idx < queue.Size() && idx < queue.CurrentIndex+6; idx++ {
			song := queue.Songs[idx]
			queueText += fmt.Sprintf("%d. **%s** - %s\n", idx-queue.CurrentIndex, song.Title, entity.FormatDuration(song.Duration))
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

// playWithLavalink handles voice joining and playback via Lavalink
func (b *Bot) playWithLavalink(guildID, channelID string, song *entity.Song) {
	b.logger.Info("Attempting to play with Lavalink",
		"guild", guildID, "channel", channelID, "song", song.Title)

	// Get session ID for Lavalink
	sessionID := b.session.State.User.ID

	// Search for track on Lavalink
	query := fmt.Sprintf("ytsearch:%s %s", song.Title, song.Artist)
	tracks, err := b.lavalinkClient.SearchTracks(query)
	if err != nil {
		b.logger.Error("Failed to search tracks on Lavalink", "error", err)
		return
	}

	if len(tracks) == 0 {
		b.logger.Warn("No tracks found on Lavalink", "query", query)
		return
	}

	// Get the first track
	track := tracks[0]

	// Tell Lavalink to join voice and play the track
	if err := b.lavalinkClient.JoinVoice(guildID, sessionID, channelID, track.Encoded); err != nil {
		b.logger.Error("Failed to join voice channel with Lavalink", "error", err)
		return
	}

	b.logger.Info("Track playing on Lavalink", "guild", guildID, "track", track.Info.Title)

	// Auto-disconnect after track duration + buffer
	duration := time.Duration(track.Info.Length) * time.Millisecond
	go func() {
		time.Sleep(duration + 3*time.Second)
		if err := b.lavalinkClient.Stop(guildID, sessionID); err != nil {
			b.logger.Warn("Failed to stop playback", "error", err)
		}
		b.logger.Info("Auto-disconnected from voice", "guild", guildID)
	}()
}
