package discord

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/usecase/music"
)

// --- Phase 2: Loop ---

func (b *Bot) handleLoop(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		b.respondError(s, i, "Please select a loop mode")
		return
	}

	modeStr := options[0].StringValue()
	var mode entity.LoopMode
	switch modeStr {
	case "off":
		mode = entity.LoopOff
	case "song":
		mode = entity.LoopSong
	case "queue":
		mode = entity.LoopQueue
	default:
		b.respondError(s, i, "Invalid mode. Use: off, song, or queue")
		return
	}

	result := b.musicService.SetLoop(i.GuildID, mode)

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

	b.respond(s, i, fmt.Sprintf("%s %s", emoji, label))
}

// --- Phase 2: Filters ---

func (b *Bot) handleFilter(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		b.respond(s, i, fmt.Sprintf("%s Available filters: %s", config.EmojiFilter, strings.Join(music.FilterPresets, ", ")))
		return
	}

	filterName := options[0].StringValue()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if filterName == "clear" || filterName == "off" || filterName == "reset" {
		if err := b.musicService.ClearFilters(ctx, i.GuildID); err != nil {
			b.respondError(s, i, err.Error())
			return
		}
		b.respond(s, i, config.EmojiFilter+" Filters cleared")
		return
	}

	if err := b.musicService.SetFilter(ctx, i.GuildID, filterName); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, fmt.Sprintf("%s Applied filter: **%s**", config.EmojiFilter, filterName))
}

// --- Phase 2: Previous ---

func (b *Bot) handlePrevious(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	song, err := b.musicService.Previous(ctx, i.GuildID)
	if err != nil {
		b.followUpError(s, i, err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiPrevious + " Playing Previous",
		Description: fmt.Sprintf("**[%s](%s)**\nby %s", song.Title, song.URL, song.Author),
		Color:       0x00C9A7,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: song.Thumbnail},
	}

	b.followUpEmbed(s, i, embed)
}

// --- Phase 2: Seek ---

func (b *Bot) handleSeek(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		b.respondError(s, i, "Please provide a position in seconds")
		return
	}

	seconds := options[0].IntValue()
	positionMs := seconds * 1000

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.musicService.Seek(ctx, i.GuildID, positionMs); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, fmt.Sprintf("%s Seeked to **%s**", config.EmojiSeek, music.FormatDuration(positionMs)))
}

// --- Phase 2: Playlist commands ---

func (b *Bot) handlePlaylist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please select a subcommand")
		return
	}

	sub := options[0]
	switch sub.Name {
	case "create":
		b.handlePlaylistCreate(s, i, sub.Options)
	case "add":
		b.handlePlaylistAdd(s, i, sub.Options)
	case "play":
		b.handlePlaylistPlay(s, i, sub.Options)
	case "list":
		b.handlePlaylistList(s, i)
	case "delete":
		b.handlePlaylistDelete(s, i, sub.Options)
	case "show":
		b.handlePlaylistShow(s, i, sub.Options)
	case "import":
		b.handlePlaylistImport(s, i, sub.Options)
	default:
		b.respondError(s, i, "Unknown playlist subcommand")
	}
}

func (b *Bot) handlePlaylistCreate(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a playlist name")
		return
	}

	name := options[0].StringValue()
	if name == "" {
		b.respondError(s, i, "Please provide a playlist name")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if we should save current queue or just create empty
	var songs []*entity.Song
	queue, _, err := b.musicService.GetQueue(i.GuildID)
	if err == nil && len(queue) > 0 {
		songs = queue
	}

	if err := b.musicService.SaveCurrentAsPlaylist(ctx, i.Member.User.ID, name, songs); err != nil {
		b.respondError(s, i, fmt.Sprintf("Failed to create playlist: %s", err.Error()))
		return
	}

	songCount := len(songs)
	msg := fmt.Sprintf("%s Created playlist **%s**", config.EmojiPlaylist, name)
	if songCount > 0 {
		msg += fmt.Sprintf(" with **%d** songs from current queue", songCount)
	}
	b.respond(s, i, msg)
}

func (b *Bot) handlePlaylistAdd(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a playlist name")
		return
	}

	name := options[0].StringValue()
	playlistID := i.Member.User.ID + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "_")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	song, err := b.musicService.AddCurrentToPlaylist(ctx, i.GuildID, i.Member.User.ID, playlistID)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, fmt.Sprintf("%s Added **%s** to playlist **%s**", config.EmojiPlaylist, song.Title, name))
}

func (b *Bot) handlePlaylistPlay(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	if len(options) == 0 {
		b.respondError(s, i, "Please provide a playlist name")
		return
	}

	name := options[0].StringValue()
	playlistID := i.Member.User.ID + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "_")

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	playlist, err := b.musicService.LoadPlaylist(ctx, playlistID)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Playlist **%s** not found", name))
		return
	}

	count, err := b.musicService.PlayPlaylist(ctx, i.GuildID, vs.ChannelID, i.ChannelID, i.Member.User.ID, playlist)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to load playlist: %s", err.Error()))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiPlaylist + " Playlist Loaded",
		Description: fmt.Sprintf("**%s**\nLoaded **%d** songs into the queue", playlist.Name, count),
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
	}
	b.followUpEmbed(s, i, embed)
}

func (b *Bot) handlePlaylistList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playlists, err := b.musicService.GetUserPlaylists(ctx, i.Member.User.ID)
	if err != nil || len(playlists) == 0 {
		b.respondError(s, i, "You don't have any playlists. Create one with `/playlist create <name>`")
		return
	}

	var desc strings.Builder
	for idx, pl := range playlists {
		desc.WriteString(fmt.Sprintf("`%d.` **%s** — %d songs\n", idx+1, pl.Name, len(pl.Songs)))
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiPlaylist + " Your Playlists",
		Description: desc.String(),
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d playlists", len(playlists)),
		},
	}
	b.respondEmbed(s, i, embed)
}

func (b *Bot) handlePlaylistDelete(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a playlist name")
		return
	}

	name := options[0].StringValue()
	playlistID := i.Member.User.ID + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "_")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify ownership
	playlist, err := b.musicService.LoadPlaylist(ctx, playlistID)
	if err != nil || playlist == nil {
		b.respondError(s, i, fmt.Sprintf("Playlist **%s** not found", name))
		return
	}
	if playlist.UserID != i.Member.User.ID {
		b.respondError(s, i, "You can only delete your own playlists")
		return
	}

	if err := b.musicService.DeletePlaylist(ctx, playlistID); err != nil {
		b.respondError(s, i, fmt.Sprintf("Failed to delete: %s", err.Error()))
		return
	}

	b.respond(s, i, fmt.Sprintf("%s Deleted playlist **%s**", config.EmojiPlaylist, name))
}

func (b *Bot) handlePlaylistShow(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a playlist name")
		return
	}

	name := options[0].StringValue()
	playlistID := i.Member.User.ID + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "_")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playlist, err := b.musicService.LoadPlaylist(ctx, playlistID)
	if err != nil || playlist == nil {
		b.respondError(s, i, fmt.Sprintf("Playlist **%s** not found", name))
		return
	}

	var desc strings.Builder
	maxShow := 15
	for idx, song := range playlist.Songs {
		if idx >= maxShow {
			desc.WriteString(fmt.Sprintf("\n... and **%d** more songs", len(playlist.Songs)-maxShow))
			break
		}
		desc.WriteString(fmt.Sprintf("`%d.` [%s](%s) • %s\n", idx+1, song.Title, song.URL, song.FormatDuration()))
	}

	var totalDuration int64
	for _, song := range playlist.Songs {
		totalDuration += song.Duration
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiPlaylist + " " + playlist.Name,
		Description: desc.String(),
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d songs • Total: %s", len(playlist.Songs), music.FormatDuration(totalDuration)),
		},
	}
	b.respondEmbed(s, i, embed)
}

// --- Phase 3: DJ role ---

func (b *Bot) handleDJ(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please select a subcommand")
		return
	}

	sub := options[0]
	switch sub.Name {
	case "set":
		b.handleDJSet(s, i, sub.Options)
	case "remove":
		b.handleDJRemove(s, i)
	case "check":
		b.handleDJCheck(s, i)
	default:
		b.respondError(s, i, "Unknown DJ subcommand")
	}
}

func (b *Bot) handleDJSet(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondError(s, i, "Please mention a role")
		return
	}

	roleID := options[0].RoleValue(s, i.GuildID).ID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settings, err := b.musicService.GetGuildMusicSettings(ctx, i.GuildID)
	if err != nil || settings == nil {
		settings = &entity.GuildMusicSettings{GuildID: i.GuildID, DefaultVolume: 100, MaxQueueLen: 500}
	}
	settings.DJRoleID = roleID

	if err := b.musicService.SaveGuildMusicSettings(ctx, settings); err != nil {
		b.respondError(s, i, "Failed to save DJ role")
		return
	}

	b.respond(s, i, fmt.Sprintf("%s DJ role set to <@&%s>. Only members with this role can control music.", config.EmojiDJ, roleID))
}

func (b *Bot) handleDJRemove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settings, err := b.musicService.GetGuildMusicSettings(ctx, i.GuildID)
	if err != nil || settings == nil {
		b.respond(s, i, config.EmojiDJ+" No DJ role is configured")
		return
	}

	settings.DJRoleID = ""
	if err := b.musicService.SaveGuildMusicSettings(ctx, settings); err != nil {
		b.respondError(s, i, "Failed to remove DJ role")
		return
	}

	b.respond(s, i, config.EmojiDJ+" DJ role removed. Everyone can now control music.")
}

func (b *Bot) handleDJCheck(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settings, err := b.musicService.GetGuildMusicSettings(ctx, i.GuildID)
	if err != nil || settings == nil || settings.DJRoleID == "" {
		b.respond(s, i, config.EmojiDJ+" No DJ role configured. Everyone can control music.")
		return
	}

	hasDJ := b.musicService.CheckDJPermission(ctx, i.GuildID, i.Member.Roles)
	status := "❌ You do **not** have DJ permissions"
	if hasDJ {
		status = "✅ You have DJ permissions"
	}

	b.respond(s, i, fmt.Sprintf("%s DJ Role: <@&%s>\n%s", config.EmojiDJ, settings.DJRoleID, status))
}

// --- Phase 3: Vote Skip ---

func (b *Bot) handleVoteSkip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	if !b.musicService.IsPlaying(i.GuildID) {
		b.respondError(s, i, "Nothing is playing")
		return
	}

	// Register vote via the service (use a simple approach with skip votes map)
	voted := b.registerSkipVote(i.GuildID, i.Member.User.ID)
	if !voted {
		b.respondError(s, i, "You already voted to skip")
		return
	}

	voteCount := b.getSkipVoteCount(i.GuildID)
	totalUsers := b.musicService.GetVoiceChannelUserCount(i.GuildID)
	needed := (totalUsers / 2) + 1
	if needed < 1 {
		needed = 1
	}

	if voteCount >= needed {
		// Enough votes, skip the track
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		b.clearSkipVotes(i.GuildID)
		nextSong, err := b.musicService.Skip(ctx, i.GuildID)
		if err != nil {
			b.respondError(s, i, err.Error())
			return
		}

		if nextSong == nil {
			b.respond(s, i, config.EmojiVoteSkip+" Vote skip passed — queue is now empty")
		} else {
			b.respond(s, i, fmt.Sprintf("%s Vote skip passed! Now playing: **%s**", config.EmojiVoteSkip, nextSong.Title))
		}
		return
	}

	b.respond(s, i, fmt.Sprintf("%s Vote skip: **%d/%d** votes (need %d)", config.EmojiVoteSkip, voteCount, totalUsers, needed))
}

// Skip vote helpers (stored in-memory on bot, keyed by guildID)
var (
	skipVotes   = make(map[string]map[string]bool)
	skipVotesMu = &sync.Mutex{}
)

func (b *Bot) registerSkipVote(guildID, userID string) bool {
	skipVotesMu.Lock()
	defer skipVotesMu.Unlock()

	if skipVotes[guildID] == nil {
		skipVotes[guildID] = make(map[string]bool)
	}
	if skipVotes[guildID][userID] {
		return false
	}
	skipVotes[guildID][userID] = true
	return true
}

func (b *Bot) getSkipVoteCount(guildID string) int {
	skipVotesMu.Lock()
	defer skipVotesMu.Unlock()
	return len(skipVotes[guildID])
}

func (b *Bot) clearSkipVotes(guildID string) {
	skipVotesMu.Lock()
	defer skipVotesMu.Unlock()
	delete(skipVotes, guildID)
}

// --- Phase 3: Lyrics ---

func (b *Bot) handleLyrics(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	song, _, err := b.musicService.NowPlaying(i.GuildID)
	if err != nil {
		b.respondError(s, i, "Nothing is playing right now")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Try to load lyrics from Lavalink's LavaLyrics plugin
	lyrics, err := b.musicService.GetLyrics(ctx, i.GuildID)
	if err != nil || lyrics == "" {
		b.followUpError(s, i, fmt.Sprintf("Lyrics not found for **%s**", song.Title))
		return
	}

	// Truncate if too long for Discord embed (max 4096 chars)
	if len(lyrics) > 3900 {
		lyrics = lyrics[:3900] + "\n\n*... truncated*"
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiLyrics + " " + song.Title,
		Description: lyrics,
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("by %s", song.Author),
		},
	}
	b.followUpEmbed(s, i, embed)
}

// --- Phase 3: 24/7 ---

func (b *Bot) handle247(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	current := b.musicService.Is247(ctx, i.GuildID)
	newState := !current

	if err := b.musicService.Set247(ctx, i.GuildID, newState); err != nil {
		b.respondError(s, i, "Failed to toggle 24/7 mode")
		return
	}

	if newState {
		b.respond(s, i, config.Emoji247+" **24/7 mode enabled** — I'll stay in voice even when alone")
	} else {
		b.respond(s, i, config.Emoji247+" **24/7 mode disabled** — I'll leave when the queue ends")
	}
}

// --- Phase 3: Autoplay ---

func (b *Bot) handleAutoplay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You need to be in a voice channel")
		return
	}

	enabled := b.musicService.ToggleAutoplay(i.GuildID)

	if enabled {
		b.respond(s, i, config.EmojiMusic+" **Autoplay enabled** — I'll auto-queue similar songs when the queue ends")
	} else {
		b.respond(s, i, config.EmojiMusic+" **Autoplay disabled**")
	}
}

// --- Phase 2: Playlist Import ---

func (b *Bot) handlePlaylistImport(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a URL and name")
		return
	}

	var url, name string
	for _, opt := range options {
		switch opt.Name {
		case "url":
			url = opt.StringValue()
		case "name":
			name = opt.StringValue()
		}
	}

	if url == "" {
		b.respondError(s, i, "Please provide a playlist URL (Spotify, YouTube, or SoundCloud)")
		return
	}
	if name == "" {
		b.respondError(s, i, "Please provide a name for the imported playlist")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	count, err := b.musicService.ImportPlaylist(ctx, i.Member.User.ID, name, url)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to import: %s", err.Error()))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiPlaylist + " Playlist Imported",
		Description: fmt.Sprintf("Saved **%s** with **%d** tracks", name, count),
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Imported by %s", i.Member.User.Username),
		},
	}
	b.followUpEmbed(s, i, embed)
}

// --- Phase 3: Recommend ---

func (b *Bot) handleRecommend(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	if !b.musicService.IsPlaying(i.GuildID) {
		b.respondError(s, i, "Nothing is playing right now")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	songs, err := b.musicService.Recommend(ctx, i.GuildID, 5)
	if err != nil {
		b.followUpError(s, i, err.Error())
		return
	}

	var desc strings.Builder
	for idx, song := range songs {
		desc.WriteString(fmt.Sprintf("`%d.` [%s](%s) • %s — %s\n", idx+1, song.Title, song.URL, song.Author, song.FormatDuration()))
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiRadio + " Recommended Songs",
		Description: desc.String(),
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use /play to queue any of these",
		},
	}
	b.followUpEmbed(s, i, embed)
}

// --- Phase 3: Radio ---

func (b *Bot) handleRadio(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		b.respondError(s, i, "Please provide a genre")
		return
	}

	genre := options[0].StringValue()

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	song, err := b.musicService.StartRadio(ctx, i.GuildID, vs.ChannelID, i.ChannelID, i.Member.User.ID, genre)
	if err != nil {
		b.followUpError(s, i, err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       config.EmojiRadio + " Radio: " + genre,
		Description: fmt.Sprintf("Now playing: **[%s](%s)**\nby %s\n\nAutoplay is enabled — music will keep playing!", song.Title, song.URL, song.Author),
		Color:       0x00C9A7,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
	}
	b.followUpEmbed(s, i, embed)
}

// --- Phase 3: AutoJoin ---

func (b *Bot) handleAutoJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.musicService == nil {
		b.respondError(s, i, "Music is not available")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := i.ApplicationCommandData().Options

	// If no options, toggle off (clear the auto-join channel)
	if len(options) == 0 {
		current := b.musicService.GetAutoJoinChannel(ctx, i.GuildID)
		if current == "" {
			b.respondError(s, i, "No auto-join channel configured. Use `/autojoin channel:<voice channel>` to set one.")
			return
		}
		// Clear auto-join
		if err := b.musicService.SetAutoJoinChannel(ctx, i.GuildID, ""); err != nil {
			b.respondError(s, i, "Failed to clear auto-join")
			return
		}
		b.respond(s, i, config.EmojiMusic+" **Auto-join disabled**")
		return
	}

	channelID := options[0].ChannelValue(s).ID
	if err := b.musicService.SetAutoJoinChannel(ctx, i.GuildID, channelID); err != nil {
		b.respondError(s, i, "Failed to configure auto-join")
		return
	}

	b.respond(s, i, fmt.Sprintf("%s **Auto-join enabled** for <#%s> — I'll join when someone enters this channel", config.EmojiMusic, channelID))
}

// onVoiceStateUpdateAutoJoin handles auto-join logic when users enter a configured voice channel.
func (b *Bot) onVoiceStateUpdateAutoJoin(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if b.musicService == nil {
		return
	}

	// Ignore bot's own voice state changes
	if v.UserID == s.State.User.ID {
		return
	}

	// Only care about joins (ChannelID is set, and it's different from before)
	if v.ChannelID == "" {
		return
	}

	// Check if this guild has auto-join configured
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	autoJoinCh := b.musicService.GetAutoJoinChannel(ctx, v.GuildID)
	if autoJoinCh == "" || autoJoinCh != v.ChannelID {
		return
	}

	// Check if bot is already in a voice channel in this guild
	if b.musicService.IsPlaying(v.GuildID) {
		return
	}

	// Join the voice channel silently
	if err := s.ChannelVoiceJoinManual(v.GuildID, autoJoinCh, false, true); err != nil {
		b.logger.Warn("Auto-join failed", "guild", v.GuildID, "channel", autoJoinCh, "error", err)
	} else {
		b.logger.Info("Auto-joined voice channel", "guild", v.GuildID, "channel", autoJoinCh)
	}
}
