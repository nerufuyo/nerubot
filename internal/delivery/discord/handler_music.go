       // Block SARA/porn queries
       if containsBlockedKeyword(query) {
	       b.respondError(s, i, "Sorry, this request cannot be processed.")
	       return
       }
package discord

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/usecase/music"
)

// handlePlay searches and plays a song, or adds it to the queue.
func (b *Bot) handlePlay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	// User must be in a voice channel
	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel to play music")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a song name or URL")
		return
	}

	query := options[0].StringValue()
	if query == "" {
		b.respondError(s, i, "Please provide a song name or URL")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	song, queued, err := b.musicService.Play(ctx, i.GuildID, vs.ChannelID, i.ChannelID, i.Member.User.ID, query)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to play: %s", err.Error()))
		return
	}

	       var embed *discordgo.MessageEmbed
	       if queued {
		       embed = &discordgo.MessageEmbed{
			       Title:       config.EmojiQueue + " Added to Queue",
			       Description: fmt.Sprintf("**[%s](%s)**\nby %s • %s", song.Title, song.URL, song.Author, song.FormatDuration()),
			       Color:       0x00C9A7,
			       Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: song.Thumbnail},
			       Footer: &discordgo.MessageEmbedFooter{
				       Text: fmt.Sprintf("Position: #%d", b.musicService.QueueLength(i.GuildID)),
			       },
		       }
	       } else {
		       embed = &discordgo.MessageEmbed{
			       Title:       config.EmojiNowPlaying + " Now Playing",
			       Description: fmt.Sprintf("**[%s](%s)**\nby %s • %s", song.Title, song.URL, song.Author, song.FormatDuration()),
			       Color:       0x00C9A7,
			       Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: song.Thumbnail},
		       }
	       }

	       b.followUpEmbedWithButtons(s, i, embed, nowPlayingButtons())
}

// handlePause pauses the current track.
func (b *Bot) handlePause(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.musicService.Pause(ctx, i.GuildID); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, config.EmojiPause+" Paused")
}

// handleResume resumes the current track.
func (b *Bot) handleResume(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.musicService.Resume(ctx, i.GuildID); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, config.EmojiPlay+" Resumed")
}

// handleStop stops playback, clears queue, and disconnects.
func (b *Bot) handleStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.musicService.Stop(ctx, i.GuildID); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, config.EmojiStop+" Stopped and disconnected")
}

// handleSkip skips the current track.
func (b *Bot) handleSkip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nextSong, err := b.musicService.Skip(ctx, i.GuildID)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	if nextSong == nil {
		b.respond(s, i, config.EmojiSkip+" Skipped — queue is now empty")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiSkip + " Skipped",
		Description: fmt.Sprintf("Now playing: **[%s](%s)**\nby %s • %s", nextSong.Title, nextSong.URL, nextSong.Author, nextSong.FormatDuration()),
		Color:       0x00C9A7,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: nextSong.Thumbnail},
	}
	b.respondEmbed(s, i, embed)
}

// handleNowPlaying shows the currently playing track with progress bar.
func (b *Bot) handleNowPlaying(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	song, position, err := b.musicService.NowPlaying(i.GuildID)
	if err != nil {
		b.respondError(s, i, "Nothing is playing right now")
		return
	}

	progressBar := music.FormatProgressBar(position, song.Duration, 20)
	posStr := music.FormatDuration(position)
	durStr := music.FormatDuration(song.Duration)

	pauseLabel := ""
	if b.musicService.IsPaused(i.GuildID) {
		pauseLabel = " (Paused)"
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiNowPlaying + " Now Playing" + pauseLabel,
		Description: fmt.Sprintf("**[%s](%s)**\nby %s\n\n%s\n`%s / %s`", song.Title, song.URL, song.Author, progressBar, posStr, durStr),
		Color:       0x00C9A7,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: song.Thumbnail},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Source: %s • Volume: %d%%", song.Source, b.musicService.GetVolume(i.GuildID)),
		},
	}

	// Respond with buttons
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: nowPlayingComponents(),
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond with now playing", "error", err)
	}
}

// handleQueue shows the current queue.
func (b *Bot) handleQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	queue, currentIdx, err := b.musicService.GetQueue(i.GuildID)
	if err != nil {
		b.respondError(s, i, "Queue is empty")
		return
	}

	// Get page from options (default 1)
	page := 1
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "page" {
			page = int(opt.IntValue())
		}
	}

	pageSize := 10
	totalPages := (len(queue) + pageSize - 1) / pageSize
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(queue) {
		endIdx = len(queue)
	}

	var desc strings.Builder

	// Show currently playing
	if currentIdx < len(queue) {
		current := queue[currentIdx]
		desc.WriteString(fmt.Sprintf("**Now Playing:**\n%s [%s](%s) • %s\n\n", config.EmojiNowPlaying, current.Title, current.URL, current.FormatDuration()))
	}

	// Show queue items
	if len(queue) > 1 {
		desc.WriteString("**Up Next:**\n")
		for idx := startIdx; idx < endIdx; idx++ {
			if idx == currentIdx {
				continue // Skip current song in list
			}
			song := queue[idx]
			marker := ""
			if idx == currentIdx {
				marker = " ◄"
			}
			desc.WriteString(fmt.Sprintf("`%d.` [%s](%s) • %s%s\n", idx+1, song.Title, song.URL, song.FormatDuration(), marker))
		}
	}

	// Total duration
	var totalDuration int64
	for _, song := range queue {
		totalDuration += song.Duration
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiQueue + " Music Queue",
		Description: desc.String(),
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d • %d songs • Total: %s", page, totalPages, len(queue), music.FormatDuration(totalDuration)),
		},
	}

	b.respondEmbed(s, i, embed)
}

// handleVolume sets the playback volume.
func (b *Bot) handleVolume(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		// Show current volume
		vol := b.musicService.GetVolume(i.GuildID)
		b.respond(s, i, fmt.Sprintf("%s Volume: **%d%%**", config.EmojiVolume, vol))
		return
	}

	volume := int(options[0].IntValue())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.musicService.SetVolume(ctx, i.GuildID, volume); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, fmt.Sprintf("%s Volume set to **%d%%**", config.EmojiVolume, volume))
}

// handleRemove removes a song from the queue by position.
func (b *Bot) handleRemove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a position number")
		return
	}

	position := int(options[0].IntValue())
	song, err := b.musicService.RemoveFromQueue(i.GuildID, position)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	       // Block SARA/porn in song title
	       if containsBlockedKeyword(song.Title) {
		       b.respondError(s, i, "Sorry, this request cannot be processed.")
		       return
	       }
	       b.respond(s, i, fmt.Sprintf("Removed **%s** from the queue", song.Title))
}

// handleClear clears the queue except the current song.
func (b *Bot) handleClear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cleared, err := b.musicService.ClearQueue(ctx, i.GuildID)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, fmt.Sprintf("Cleared **%d** songs from the queue", cleared))
}

// handleShuffle shuffles the remaining queue.
func (b *Bot) handleShuffle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	if err := b.musicService.ShuffleQueue(i.GuildID); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, config.EmojiShuffle+" Queue shuffled!")
}

// handleMove moves a song in the queue.
func (b *Bot) handleMove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		b.respondError(s, i, "Please provide `from` and `to` positions")
		return
	}

	var from, to int
	for _, opt := range options {
		switch opt.Name {
		case "from":
			from = int(opt.IntValue())
		case "to":
			to = int(opt.IntValue())
		}
	}

	song, err := b.musicService.MoveInQueue(i.GuildID, from, to)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	       // Block SARA/porn in song title
	       if containsBlockedKeyword(song.Title) {
		       b.respondError(s, i, "Sorry, this request cannot be processed.")
		       return
	       }
	       b.respond(s, i, fmt.Sprintf("Moved **%s** from position %d to %d", song.Title, from, to))
}

// --- Helper: follow up with buttons ---

func (b *Bot) followUpEmbedWithButtons(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, buttons []discordgo.MessageComponent) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: buttons,
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up with buttons", "error", err)
	}
}

// --- Button definitions ---

func nowPlayingButtons() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "music_previous",
					Emoji: &discordgo.ComponentEmoji{
						Name: "⏮",
					},
					Style: discordgo.SecondaryButton,
				},
				discordgo.Button{
					CustomID: "music_pause_resume",
					Emoji: &discordgo.ComponentEmoji{
						Name: "⏯",
					},
					Style: discordgo.PrimaryButton,
				},
				discordgo.Button{
					CustomID: "music_skip",
					Emoji: &discordgo.ComponentEmoji{
						Name: "⏭",
					},
					Style: discordgo.SecondaryButton,
				},
				discordgo.Button{
					CustomID: "music_shuffle",
					Emoji: &discordgo.ComponentEmoji{
						Name: "🔀",
					},
					Style: discordgo.SecondaryButton,
				},
				discordgo.Button{
					CustomID: "music_stop",
					Emoji: &discordgo.ComponentEmoji{
						Name: "⏹",
					},
					Style: discordgo.DangerButton,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "music_loop",
					Emoji: &discordgo.ComponentEmoji{
						Name: "🔁",
					},
					Label: "Loop",
					Style: discordgo.SecondaryButton,
				},
				discordgo.Button{
					CustomID: "music_autoplay",
					Emoji: &discordgo.ComponentEmoji{
						Name: "🎵",
					},
					Label: "Autoplay",
					Style: discordgo.SecondaryButton,
				},
			},
		},
	}
}

func nowPlayingComponents() []discordgo.MessageComponent {
	return nowPlayingButtons()
}

// volumeBar returns a visual volume indicator.
func volumeBar(volume int) string {
	filled := volume / 10
	if filled > 20 {
		filled = 20
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", 20-filled)
}

// Ensure strconv is used (for future page parsing etc.)
var _ = strconv.Atoi
