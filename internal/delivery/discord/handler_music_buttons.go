package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/usecase/music"
)

// handleMusicButton routes button interactions for the music player.
func (b *Bot) handleMusicButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondButtonError(s, i, "Music is not available")
		return
	}

	// User must be in voice
	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondButtonError(s, i, "You need to be in a voice channel")
		return
	}

	customID := i.MessageComponentData().CustomID
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch customID {
	case "music_pause_resume":
		b.handleButtonPauseResume(ctx, s, i)
	case "music_skip":
		b.handleButtonSkip(ctx, s, i)
	case "music_previous":
		b.handleButtonPrevious(ctx, s, i)
	case "music_shuffle":
		b.handleButtonShuffle(s, i)
	case "music_stop":
		b.handleButtonStop(ctx, s, i)
	case "music_loop":
		b.handleButtonLoop(s, i)
	case "music_autoplay":
		b.handleButtonAutoplay(s, i)
	default:
		b.respondButtonError(s, i, "Unknown button action")
	}
}

func (b *Bot) handleButtonPauseResume(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService.IsPaused(i.GuildID) {
		if err := b.musicService.Resume(ctx, i.GuildID); err != nil {
			b.respondButtonError(s, i, err.Error())
			return
		}
		b.respondButton(s, i, config.EmojiPlay+" Resumed")
	} else {
		if err := b.musicService.Pause(ctx, i.GuildID); err != nil {
			b.respondButtonError(s, i, err.Error())
			return
		}
		b.respondButton(s, i, config.EmojiPause+" Paused")
	}
}

func (b *Bot) handleButtonSkip(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	nextSong, err := b.musicService.Skip(ctx, i.GuildID)
	if err != nil {
		b.respondButtonError(s, i, err.Error())
		return
	}

	if nextSong == nil {
		b.respondButton(s, i, config.EmojiSkip+" Skipped — queue is now empty")
		return
	}

	b.respondButton(s, i, fmt.Sprintf("%s Now playing: **%s**", config.EmojiSkip, nextSong.Title))
}

func (b *Bot) handleButtonPrevious(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	song, err := b.musicService.Previous(ctx, i.GuildID)
	if err != nil {
		b.respondButtonError(s, i, err.Error())
		return
	}
	b.respondButton(s, i, fmt.Sprintf("%s Playing previous: **%s**", config.EmojiPrevious, song.Title))
}

func (b *Bot) handleButtonShuffle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := b.musicService.ShuffleQueue(i.GuildID); err != nil {
		b.respondButtonError(s, i, err.Error())
		return
	}
	b.respondButton(s, i, config.EmojiShuffle+" Queue shuffled!")
}

func (b *Bot) handleButtonStop(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := b.musicService.Stop(ctx, i.GuildID); err != nil {
		b.respondButtonError(s, i, err.Error())
		return
	}
	b.respondButton(s, i, config.EmojiStop+" Stopped and disconnected")
}

func (b *Bot) handleButtonLoop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Cycle: Off -> Song -> Queue -> Off
	current := b.musicService.GetLoop(i.GuildID)
	var next entity.LoopMode
	switch current {
	case entity.LoopOff:
		next = entity.LoopSong
	case entity.LoopSong:
		next = entity.LoopQueue
	default:
		next = entity.LoopOff
	}

	result := b.musicService.SetLoop(i.GuildID, next)

	var emoji, label string
	switch result {
	case entity.LoopSong:
		emoji = config.EmojiLoopOne
		label = "Looping current song"
	case entity.LoopQueue:
		emoji = config.EmojiLoop
		label = "Looping entire queue"
	default:
		emoji = config.EmojiLoop
		label = "Loop disabled"
	}
	b.respondButton(s, i, fmt.Sprintf("%s %s", emoji, label))
}

func (b *Bot) handleButtonAutoplay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	enabled := b.musicService.ToggleAutoplay(i.GuildID)
	if enabled {
		b.respondButton(s, i, config.EmojiMusic+" Autoplay enabled")
	} else {
		b.respondButton(s, i, config.EmojiMusic+" Autoplay disabled")
	}
}

// --- Now Playing embed updater ---

// buildNowPlayingEmbed creates a fresh Now Playing embed snapshot.
func (b *Bot) buildNowPlayingEmbed(guildID string) *discordgo.MessageEmbed {
	song, position, err := b.musicService.NowPlaying(guildID)
	if err != nil {
		return &discordgo.MessageEmbed{
			Description: "Nothing is playing",
			Color:       0x6C757D,
		}
	}

	progressBar := music.FormatProgressBar(position, song.Duration, 20)
	posStr := music.FormatDuration(position)
	durStr := music.FormatDuration(song.Duration)

	pauseLabel := ""
	if b.musicService.IsPaused(guildID) {
		pauseLabel = " (Paused)"
	}

	return &discordgo.MessageEmbed{
		Title:       config.EmojiNowPlaying + " Now Playing" + pauseLabel,
		Description: fmt.Sprintf("**[%s](%s)**\nby %s\n\n%s\n`%s / %s`", song.Title, song.URL, song.Author, progressBar, posStr, durStr),
		Color:       0x00C9A7,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: song.Thumbnail},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Source: %s • Volume: %d%%", song.Source, b.musicService.GetVolume(guildID)),
		},
	}
}

// --- Button response helpers ---

// respondButton sends an ephemeral text response to a button interaction.
func (b *Bot) respondButton(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond to button", "error", err)
	}
}

// respondButtonError sends an ephemeral error response to a button interaction.
func (b *Bot) respondButtonError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	b.respondButton(s, i, config.EmojiError+" "+message)
}
