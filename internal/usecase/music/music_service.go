package music

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/nerufuyo/nerubot/internal/entity"
	lavalinkpkg "github.com/nerufuyo/nerubot/internal/pkg/lavalink"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	redispkg "github.com/nerufuyo/nerubot/internal/pkg/redis"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// SendEmbedFunc is a callback to send an embed to a Discord text channel.
type SendEmbedFunc func(channelID string, embed *discordgo.MessageEmbed)

// MusicService manages per-guild music playback state and coordinates with Lavalink.
type MusicService struct {
	lavalink *lavalinkpkg.Client
	session  *discordgo.Session
	repo     *repository.MusicRepository
	redis    *redispkg.Client
	logger   *logger.Logger

	players  map[string]*entity.GuildPlayer // guildID -> player state
	mu       sync.RWMutex

	sendEmbed SendEmbedFunc
}

// NewMusicService creates a new MusicService.
func NewMusicService(
	lc *lavalinkpkg.Client,
	session *discordgo.Session,
	repo *repository.MusicRepository,
	redis *redispkg.Client,
) *MusicService {
	log := logger.New("music")

	s := &MusicService{
		lavalink: lc,
		session:  session,
		repo:     repo,
		redis:    redis,
		logger:   log,
		players:  make(map[string]*entity.GuildPlayer),
	}

	// Register Lavalink event callbacks
	lc.OnTrackEnd(s.onTrackEnd)
	lc.OnTrackException(s.onTrackException)
	lc.OnTrackStuck(s.onTrackStuck)

	return s
}

// SetSendFunc sets the callback for sending embeds to Discord channels.
func (s *MusicService) SetSendFunc(fn SendEmbedFunc) {
	s.sendEmbed = fn
}

// --- Player state management ---

func (s *MusicService) getOrCreatePlayer(guildID string) *entity.GuildPlayer {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.players[guildID]
	if !ok {
		p = &entity.GuildPlayer{
			GuildID:  guildID,
			Volume:   100,
			LoopMode: entity.LoopOff,
			Queue:    make([]*entity.Song, 0),
			History:  make([]*entity.Song, 0),
		}
		s.players[guildID] = p
	}
	return p
}

func (s *MusicService) getPlayer(guildID string) *entity.GuildPlayer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.players[guildID]
}

func (s *MusicService) removePlayer(guildID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.players, guildID)
}

// --- Core playback operations ---

// Play searches for a track and starts playback, or adds to queue if already playing.
func (s *MusicService) Play(ctx context.Context, guildID, voiceChID, textChID, requesterID, query string) (*entity.Song, bool, error) {
	// Resolve the query into a Lavalink search string
	searchQuery := resolveQuery(query)

	result, err := s.lavalink.LoadTracks(ctx, searchQuery)
	if err != nil {
		return nil, false, fmt.Errorf("failed to search tracks: %w", err)
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		return nil, false, fmt.Errorf("no results found for: %s", query)
	}

	track := tracks[0]
	song := trackToSong(track, requesterID)

	gp := s.getOrCreatePlayer(guildID)
	gp.TextChID = textChID
	gp.ChannelID = voiceChID

	// Add to queue
	gp.Queue = append(gp.Queue, song)

	// If not currently playing, start playback
	if !gp.IsPlaying {
		gp.Current = len(gp.Queue) - 1
		if err := s.startPlayback(ctx, guildID, gp, track); err != nil {
			return nil, false, err
		}
		return song, false, nil // false = not queued, playing now
	}

	s.saveQueueToRedis(guildID, gp)
	return song, true, nil // true = added to queue
}

// PlayMultiple adds multiple tracks (e.g., from a playlist) to the queue.
func (s *MusicService) PlayMultiple(ctx context.Context, guildID, voiceChID, textChID, requesterID, query string) ([]*entity.Song, error) {
	result, err := s.lavalink.LoadTracks(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to load playlist: %w", err)
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		return nil, fmt.Errorf("no tracks found")
	}

	gp := s.getOrCreatePlayer(guildID)
	gp.TextChID = textChID
	gp.ChannelID = voiceChID

	var songs []*entity.Song
	for _, t := range tracks {
		song := trackToSong(t, requesterID)
		gp.Queue = append(gp.Queue, song)
		songs = append(songs, song)
	}

	// If not currently playing, start the first track
	if !gp.IsPlaying {
		gp.Current = len(gp.Queue) - len(songs)
		if err := s.startPlayback(ctx, guildID, gp, tracks[0]); err != nil {
			return nil, err
		}
	}

	s.saveQueueToRedis(guildID, gp)
	return songs, nil
}

// Pause pauses playback.
func (s *MusicService) Pause(ctx context.Context, guildID string) error {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return fmt.Errorf("nothing is playing")
	}

	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithPaused(true)); err != nil {
		return fmt.Errorf("failed to pause: %w", err)
	}
	gp.IsPaused = true
	return nil
}

// Resume resumes playback.
func (s *MusicService) Resume(ctx context.Context, guildID string) error {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return fmt.Errorf("nothing is playing")
	}

	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithPaused(false)); err != nil {
		return fmt.Errorf("failed to resume: %w", err)
	}
	gp.IsPaused = false
	return nil
}

// Skip skips the current track and plays the next one.
func (s *MusicService) Skip(ctx context.Context, guildID string) (*entity.Song, error) {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return nil, fmt.Errorf("nothing is playing")
	}

	nextSong, nextTrack := s.advanceQueue(gp)
	if nextSong == nil {
		// No more songs, stop
		s.stopPlayback(ctx, guildID, gp)
		return nil, nil
	}

	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithTrack(nextTrack)); err != nil {
		return nil, fmt.Errorf("failed to skip: %w", err)
	}
	gp.IsPaused = false
	s.saveQueueToRedis(guildID, gp)
	return nextSong, nil
}

// Stop stops playback, clears the queue, and disconnects from voice.
func (s *MusicService) Stop(ctx context.Context, guildID string) error {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return fmt.Errorf("no active player")
	}

	s.stopPlayback(ctx, guildID, gp)
	s.disconnectVoice(guildID)
	s.removePlayer(guildID)
	s.deleteQueueFromRedis(guildID)
	return nil
}

// NowPlaying returns the currently playing song and position.
func (s *MusicService) NowPlaying(guildID string) (*entity.Song, int64, error) {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return nil, 0, fmt.Errorf("nothing is playing")
	}

	song := gp.NowPlaying()
	if song == nil {
		return nil, 0, fmt.Errorf("nothing is playing")
	}

	// Get position from Lavalink player
	player := s.lavalink.ExistingPlayer(guildID)
	var position int64
	if player != nil {
		position = int64(player.Position().Milliseconds())
	}

	return song, position, nil
}

// GetQueue returns the current queue for a guild.
func (s *MusicService) GetQueue(guildID string) ([]*entity.Song, int, error) {
	gp := s.getPlayer(guildID)
	if gp == nil || len(gp.Queue) == 0 {
		return nil, 0, fmt.Errorf("queue is empty")
	}
	return gp.Queue, gp.Current, nil
}

// RemoveFromQueue removes a song at the given position (1-based index).
func (s *MusicService) RemoveFromQueue(guildID string, position int) (*entity.Song, error) {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return nil, fmt.Errorf("no active player")
	}

	idx := position - 1 // Convert to 0-based
	if idx < 0 || idx >= len(gp.Queue) {
		return nil, fmt.Errorf("invalid position: %d", position)
	}
	if idx == gp.Current {
		return nil, fmt.Errorf("cannot remove the currently playing song — use /skip instead")
	}

	removed := gp.Queue[idx]
	gp.Queue = append(gp.Queue[:idx], gp.Queue[idx+1:]...)

	// Adjust current index if needed
	if idx < gp.Current {
		gp.Current--
	}

	s.saveQueueToRedis(guildID, gp)
	return removed, nil
}

// ClearQueue clears all songs except the currently playing one.
func (s *MusicService) ClearQueue(ctx context.Context, guildID string) (int, error) {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return 0, fmt.Errorf("no active player")
	}

	if len(gp.Queue) <= 1 {
		return 0, nil
	}

	current := gp.NowPlaying()
	cleared := len(gp.Queue) - 1

	if current != nil && gp.IsPlaying {
		gp.Queue = []*entity.Song{current}
		gp.Current = 0
	} else {
		gp.Queue = make([]*entity.Song, 0)
		gp.Current = 0
	}

	s.saveQueueToRedis(guildID, gp)
	return cleared, nil
}

// ShuffleQueue shuffles the remaining songs in the queue (after current).
func (s *MusicService) ShuffleQueue(guildID string) error {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return fmt.Errorf("no active player")
	}

	remaining := gp.Queue[gp.Current+1:]
	if len(remaining) < 2 {
		return fmt.Errorf("not enough songs to shuffle")
	}

	rand.Shuffle(len(remaining), func(i, j int) {
		remaining[i], remaining[j] = remaining[j], remaining[i]
	})

	s.saveQueueToRedis(guildID, gp)
	return nil
}

// MoveInQueue moves a song from one position to another (1-based).
func (s *MusicService) MoveInQueue(guildID string, from, to int) (*entity.Song, error) {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return nil, fmt.Errorf("no active player")
	}

	fromIdx := from - 1
	toIdx := to - 1
	if fromIdx < 0 || fromIdx >= len(gp.Queue) || toIdx < 0 || toIdx >= len(gp.Queue) {
		return nil, fmt.Errorf("invalid position")
	}
	if fromIdx == gp.Current || toIdx == gp.Current {
		return nil, fmt.Errorf("cannot move the currently playing song")
	}

	song := gp.Queue[fromIdx]

	// Remove from old position
	gp.Queue = append(gp.Queue[:fromIdx], gp.Queue[fromIdx+1:]...)

	// Adjust current if needed
	if fromIdx < gp.Current && toIdx >= gp.Current {
		gp.Current--
	} else if fromIdx > gp.Current && toIdx <= gp.Current {
		gp.Current++
	}

	// Insert at new position
	newQueue := make([]*entity.Song, 0, len(gp.Queue)+1)
	newQueue = append(newQueue, gp.Queue[:toIdx]...)
	newQueue = append(newQueue, song)
	newQueue = append(newQueue, gp.Queue[toIdx:]...)
	gp.Queue = newQueue

	s.saveQueueToRedis(guildID, gp)
	return song, nil
}

// SetVolume sets the volume for a guild (0-200).
func (s *MusicService) SetVolume(ctx context.Context, guildID string, volume int) error {
	if volume < 0 || volume > 200 {
		return fmt.Errorf("volume must be between 0 and 200")
	}

	gp := s.getPlayer(guildID)
	if gp == nil {
		return fmt.Errorf("no active player")
	}

	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithVolume(volume)); err != nil {
		return fmt.Errorf("failed to set volume: %w", err)
	}
	gp.Volume = volume
	return nil
}

// GetVolume returns the current volume for a guild.
func (s *MusicService) GetVolume(guildID string) int {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return 100
	}
	return gp.Volume
}

// IsPlaying returns whether music is currently playing in a guild.
func (s *MusicService) IsPlaying(guildID string) bool {
	gp := s.getPlayer(guildID)
	return gp != nil && gp.IsPlaying
}

// IsPaused returns whether playback is paused in a guild.
func (s *MusicService) IsPaused(guildID string) bool {
	gp := s.getPlayer(guildID)
	return gp != nil && gp.IsPaused
}

// QueueLength returns the number of songs in the queue.
func (s *MusicService) QueueLength(guildID string) int {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return 0
	}
	return len(gp.Queue)
}

// --- Loop control ---

// SetLoop cycles or sets the loop mode for a guild.
func (s *MusicService) SetLoop(guildID string, mode entity.LoopMode) entity.LoopMode {
	gp := s.getOrCreatePlayer(guildID)
	gp.LoopMode = mode
	s.saveQueueToRedis(guildID, gp)
	return gp.LoopMode
}

// GetLoop returns the current loop mode for a guild.
func (s *MusicService) GetLoop(guildID string) entity.LoopMode {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return entity.LoopOff
	}
	return gp.LoopMode
}

// --- Audio filters ---

// SetFilter applies a named audio filter preset to the Lavalink player.
func (s *MusicService) SetFilter(ctx context.Context, guildID, filterName string) error {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return fmt.Errorf("nothing is playing")
	}

	filters, err := getFilterPreset(filterName)
	if err != nil {
		return err
	}

	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithFilters(filters)); err != nil {
		return fmt.Errorf("failed to apply filter: %w", err)
	}
	return nil
}

// ClearFilters removes all audio filters.
func (s *MusicService) ClearFilters(ctx context.Context, guildID string) error {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return fmt.Errorf("nothing is playing")
	}

	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithFilters(lavalink.Filters{})); err != nil {
		return fmt.Errorf("failed to clear filters: %w", err)
	}
	return nil
}

// --- Previous track ---

// Previous goes back to the previous song from history.
func (s *MusicService) Previous(ctx context.Context, guildID string) (*entity.Song, error) {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return nil, fmt.Errorf("nothing is playing")
	}

	if len(gp.History) == 0 {
		return nil, fmt.Errorf("no previous track in history")
	}

	// Pop last song from history
	prevSong := gp.History[len(gp.History)-1]
	gp.History = gp.History[:len(gp.History)-1]

	// Insert current song back at current position +1 (so it becomes "next")
	// and insert the prev song at current position
	gp.Queue = append(gp.Queue[:gp.Current], append([]*entity.Song{prevSong}, gp.Queue[gp.Current:]...)...)

	// Load and play the previous song
	query := resolveQuery(prevSong.URL)
	result, err := s.lavalink.LoadTracks(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to load previous track: %w", err)
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		return nil, fmt.Errorf("previous track not found")
	}

	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithTrack(tracks[0])); err != nil {
		return nil, fmt.Errorf("failed to play previous track: %w", err)
	}
	gp.IsPaused = false
	s.saveQueueToRedis(guildID, gp)
	return prevSong, nil
}

// --- Playlist operations ---

// SaveCurrentAsPlaylist saves the current queue as a named playlist.
func (s *MusicService) SaveCurrentAsPlaylist(ctx context.Context, userID, name string, songs []*entity.Song) error {
	playlist := &entity.Playlist{
		ID:     userID + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "_"),
		UserID: userID,
		Name:   name,
		Songs:  songs,
	}
	return s.repo.SavePlaylist(ctx, playlist)
}

// LoadPlaylist retrieves a playlist by ID.
func (s *MusicService) LoadPlaylist(ctx context.Context, playlistID string) (*entity.Playlist, error) {
	return s.repo.GetPlaylist(ctx, playlistID)
}

// GetUserPlaylists returns all playlists for a user.
func (s *MusicService) GetUserPlaylists(ctx context.Context, userID string) ([]*entity.Playlist, error) {
	return s.repo.GetUserPlaylists(ctx, userID)
}

// DeletePlaylist removes a playlist.
func (s *MusicService) DeletePlaylist(ctx context.Context, playlistID string) error {
	return s.repo.DeletePlaylist(ctx, playlistID)
}

// PlayPlaylist loads a saved playlist into the queue.
func (s *MusicService) PlayPlaylist(ctx context.Context, guildID, voiceChID, textChID, requesterID string, playlist *entity.Playlist) (int, error) {
	if len(playlist.Songs) == 0 {
		return 0, fmt.Errorf("playlist is empty")
	}

	gp := s.getOrCreatePlayer(guildID)
	gp.TextChID = textChID
	gp.ChannelID = voiceChID

	wasPlaying := gp.IsPlaying
	startIdx := len(gp.Queue)

	for _, song := range playlist.Songs {
		clone := *song
		clone.RequesterID = requesterID
		gp.Queue = append(gp.Queue, &clone)
	}

	if !wasPlaying {
		gp.Current = startIdx
		// Load first track
		firstSong := gp.Queue[startIdx]
		query := resolveQuery(firstSong.URL)
		result, err := s.lavalink.LoadTracks(ctx, query)
		if err != nil {
			return 0, fmt.Errorf("failed to load first playlist track: %w", err)
		}
		tracks := extractTracks(result)
		if len(tracks) == 0 {
			return 0, fmt.Errorf("first playlist track not found")
		}
		if err := s.startPlayback(ctx, guildID, gp, tracks[0]); err != nil {
			return 0, err
		}
	}

	s.saveQueueToRedis(guildID, gp)
	return len(playlist.Songs), nil
}

// AddCurrentToPlaylist adds the currently playing song to a playlist.
func (s *MusicService) AddCurrentToPlaylist(ctx context.Context, guildID, userID, playlistID string) (*entity.Song, error) {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return nil, fmt.Errorf("nothing is playing")
	}

	song := gp.NowPlaying()
	if song == nil {
		return nil, fmt.Errorf("nothing is playing")
	}

	playlist, err := s.repo.GetPlaylist(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("playlist not found")
	}
	if playlist.UserID != userID {
		return nil, fmt.Errorf("you don't own this playlist")
	}

	playlist.Songs = append(playlist.Songs, song)
	if err := s.repo.SavePlaylist(ctx, playlist); err != nil {
		return nil, fmt.Errorf("failed to save playlist: %w", err)
	}
	return song, nil
}

// --- DJ role & guild settings ---

// GetGuildMusicSettings retrieves music settings for a guild.
func (s *MusicService) GetGuildMusicSettings(ctx context.Context, guildID string) (*entity.GuildMusicSettings, error) {
	return s.repo.GetGuildMusicSettings(ctx, guildID)
}

// SaveGuildMusicSettings saves music settings for a guild.
func (s *MusicService) SaveGuildMusicSettings(ctx context.Context, settings *entity.GuildMusicSettings) error {
	return s.repo.SaveGuildMusicSettings(ctx, settings)
}

// CheckDJPermission returns true if the user has DJ permission or there's no DJ role set.
func (s *MusicService) CheckDJPermission(ctx context.Context, guildID string, memberRoles []string) bool {
	settings, err := s.repo.GetGuildMusicSettings(ctx, guildID)
	if err != nil || settings == nil || settings.DJRoleID == "" {
		return true // No DJ role configured, everyone can use
	}

	for _, roleID := range memberRoles {
		if roleID == settings.DJRoleID {
			return true
		}
	}
	return false
}

// --- 24/7 mode ---

// Is247 returns whether a guild has 24/7 mode enabled.
func (s *MusicService) Is247(ctx context.Context, guildID string) bool {
	settings, err := s.repo.GetGuildMusicSettings(ctx, guildID)
	if err != nil || settings == nil {
		return false
	}
	return settings.Stay247
}

// Set247 enables or disables 24/7 mode.
func (s *MusicService) Set247(ctx context.Context, guildID string, enabled bool) error {
	settings, err := s.repo.GetGuildMusicSettings(ctx, guildID)
	if err != nil || settings == nil {
		settings = &entity.GuildMusicSettings{GuildID: guildID, DefaultVolume: 100, MaxQueueLen: 500}
	}
	settings.Stay247 = enabled
	return s.repo.SaveGuildMusicSettings(ctx, settings)
}

// --- Seek ---

// Seek changes playback position.
func (s *MusicService) Seek(ctx context.Context, guildID string, positionMs int64) error {
	gp := s.getPlayer(guildID)
	if gp == nil || !gp.IsPlaying {
		return fmt.Errorf("nothing is playing")
	}

	player := s.lavalink.Player(guildID)
	position := lavalink.Duration(positionMs * int64(time.Millisecond))
	if err := player.Update(ctx, lavalink.WithPosition(position)); err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}
	return nil
}

// --- Vote skip ---

// GetVoiceChannelUserCount returns the number of non-bot users in the bot's voice channel.
func (s *MusicService) GetVoiceChannelUserCount(guildID string) int {
	gp := s.getPlayer(guildID)
	if gp == nil || gp.ChannelID == "" {
		return 0
	}

	guild, err := s.session.State.Guild(guildID)
	if err != nil {
		return 0
	}

	count := 0
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == gp.ChannelID && vs.UserID != s.session.State.User.ID {
			count++
		}
	}
	return count
}

// --- Autoplay ---

// ToggleAutoplay toggles autoplay mode and returns the new state.
func (s *MusicService) ToggleAutoplay(guildID string) bool {
	gp := s.getOrCreatePlayer(guildID)
	gp.Autoplay = !gp.Autoplay
	s.saveQueueToRedis(guildID, gp)
	return gp.Autoplay
}

// IsAutoplay returns whether autoplay is enabled.
func (s *MusicService) IsAutoplay(guildID string) bool {
	gp := s.getPlayer(guildID)
	if gp == nil {
		return false
	}
	return gp.Autoplay
}

// handleAutoplay searches for a similar track and queues it when autoplay is on.
func (s *MusicService) handleAutoplay(guildID string, gp *entity.GuildPlayer, player disgolink.Player) {
	// Use the last played song to find something similar
	var lastSong *entity.Song
	if len(gp.History) > 0 {
		lastSong = gp.History[len(gp.History)-1]
	} else if gp.NowPlaying() != nil {
		lastSong = gp.NowPlaying()
	}

	if lastSong == nil {
		gp.IsPlaying = false
		s.disconnectVoice(guildID)
		s.removePlayer(guildID)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Search for similar music using the author and general keywords
	query := fmt.Sprintf("ytsearch:%s mix", lastSong.Author)
	result, err := s.lavalink.LoadTracks(ctx, query)
	if err != nil || result == nil {
		s.logger.Warn("Autoplay search failed", "error", err)
		gp.IsPlaying = false
		s.disconnectVoice(guildID)
		s.removePlayer(guildID)
		return
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		gp.IsPlaying = false
		s.disconnectVoice(guildID)
		s.removePlayer(guildID)
		return
	}

	// Pick a track that's not already in recent history
	var selectedTrack lavalink.Track
	var selectedSong *entity.Song
	for _, t := range tracks {
		if !s.isInRecentHistory(gp, t.Info.Identifier) {
			selectedTrack = t
			selectedSong = trackToSong(t, s.session.State.User.ID)
			break
		}
	}

	if selectedSong == nil {
		// All results already played, just pick the first one
		selectedTrack = tracks[0]
		selectedSong = trackToSong(tracks[0], s.session.State.User.ID)
	}

	// Add to queue and play
	gp.Queue = append(gp.Queue, selectedSong)
	gp.Current = len(gp.Queue) - 1

	if err := player.Update(ctx, lavalink.WithTrack(selectedTrack)); err != nil {
		s.logger.Error("Autoplay failed to play track", "error", err)
		gp.IsPlaying = false
		return
	}

	s.saveQueueToRedis(guildID, gp)

	if s.sendEmbed != nil && gp.TextChID != "" {
		s.sendEmbed(gp.TextChID, &discordgo.MessageEmbed{
			Title:       "🎵 Autoplay",
			Description: fmt.Sprintf("**[%s](%s)**\nby %s", selectedSong.Title, selectedSong.URL, selectedSong.Author),
			Color:       0x00C9A7,
		})
	}
}

func (s *MusicService) isInRecentHistory(gp *entity.GuildPlayer, identifier string) bool {
	// Check last 10 songs in history
	start := len(gp.History) - 10
	if start < 0 {
		start = 0
	}
	for _, h := range gp.History[start:] {
		if h.Identifier == identifier {
			return true
		}
	}
	return false
}

// --- Lyrics ---

// GetLyrics retrieves lyrics for the currently playing track via Lavalink REST API.
func (s *MusicService) GetLyrics(ctx context.Context, guildID string) (string, error) {
	return s.lavalink.GetLyrics(ctx, guildID)
}

// --- Playlist Import ---

// ImportPlaylist imports an external playlist URL (Spotify, YouTube, SoundCloud) and saves it.
func (s *MusicService) ImportPlaylist(ctx context.Context, userID, name, url string) (int, error) {
	result, err := s.lavalink.LoadTracks(ctx, url)
	if err != nil {
		return 0, fmt.Errorf("failed to load playlist: %w", err)
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		return 0, fmt.Errorf("no tracks found at that URL")
	}

	var songs []*entity.Song
	for _, t := range tracks {
		songs = append(songs, trackToSong(t, userID))
	}

	if err := s.SaveCurrentAsPlaylist(ctx, userID, name, songs); err != nil {
		return 0, fmt.Errorf("failed to save imported playlist: %w", err)
	}

	return len(songs), nil
}

// --- Recommend ---

// Recommend finds similar tracks based on the currently playing song.
func (s *MusicService) Recommend(ctx context.Context, guildID string, count int) ([]*entity.Song, error) {
	s.mu.RLock()
	gp := s.players[guildID]
	s.mu.RUnlock()

	if gp == nil || gp.NowPlaying() == nil {
		return nil, fmt.Errorf("nothing is playing")
	}

	current := gp.NowPlaying()
	query := fmt.Sprintf("ytsearch:%s %s", current.Author, current.Title)
	result, err := s.lavalink.LoadTracks(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search for recommendations: %w", err)
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		return nil, fmt.Errorf("no recommendations found")
	}

	// Skip the first track if it matches the current (exact same song)
	var recommendations []*entity.Song
	for _, t := range tracks {
		if t.Info.Identifier == current.Identifier {
			continue
		}
		recommendations = append(recommendations, trackToSong(t, ""))
		if len(recommendations) >= count {
			break
		}
	}

	if len(recommendations) == 0 {
		return nil, fmt.Errorf("no recommendations found")
	}

	return recommendations, nil
}

// --- Radio ---

// StartRadio begins continuous playback for a genre or search term.
func (s *MusicService) StartRadio(ctx context.Context, guildID, voiceChID, textChID, requesterID, genre string) (*entity.Song, error) {
	query := fmt.Sprintf("ytsearch:%s radio mix", genre)
	result, err := s.lavalink.LoadTracks(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search radio: %w", err)
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		return nil, fmt.Errorf("no stations found for genre: %s", genre)
	}

	// Pick a track (pick first or random from top results)
	idx := 0
	if len(tracks) > 3 {
		idx = rand.Intn(3)
	}
	track := tracks[idx]
	song := trackToSong(track, requesterID)

	gp := s.getOrCreatePlayer(guildID)
	gp.TextChID = textChID
	gp.ChannelID = voiceChID
	gp.Autoplay = true // radio always has autoplay on

	// Clear existing queue and start fresh
	gp.Queue = []*entity.Song{song}
	gp.Current = 0

	if err := s.startPlayback(ctx, guildID, gp, track); err != nil {
		return nil, err
	}

	return song, nil
}

// --- AutoJoin ---

// GetAutoJoinChannel returns the configured auto-join voice channel for a guild.
func (s *MusicService) GetAutoJoinChannel(ctx context.Context, guildID string) string {
	settings, err := s.repo.GetGuildMusicSettings(ctx, guildID)
	if err != nil || settings == nil {
		return ""
	}
	return settings.AutoJoinChID
}

// SetAutoJoinChannel configures or clears the auto-join voice channel.
func (s *MusicService) SetAutoJoinChannel(ctx context.Context, guildID, channelID string) error {
	settings, err := s.repo.GetGuildMusicSettings(ctx, guildID)
	if err != nil || settings == nil {
		settings = &entity.GuildMusicSettings{
			GuildID:       guildID,
			DefaultVolume: 100,
			MaxQueueLen:   500,
		}
	}
	settings.AutoJoinChID = channelID
	return s.repo.SaveGuildMusicSettings(ctx, settings)
}

// --- Internal playback helpers ---

func (s *MusicService) startPlayback(ctx context.Context, guildID string, gp *entity.GuildPlayer, track lavalink.Track) error {
	// Join voice channel via Discord gateway (not Lavalink)
	if err := s.joinVoice(guildID, gp.ChannelID); err != nil {
		return fmt.Errorf("failed to join voice channel: %w", err)
	}

	// Wait briefly for the voice connection to establish
	time.Sleep(500 * time.Millisecond)

	// Start track on Lavalink player
	player := s.lavalink.Player(guildID)
	if err := player.Update(ctx, lavalink.WithTrack(track), lavalink.WithVolume(gp.Volume)); err != nil {
		return fmt.Errorf("failed to start playback: %w", err)
	}

	gp.IsPlaying = true
	gp.IsPaused = false
	s.saveQueueToRedis(guildID, gp)
	return nil
}

func (s *MusicService) stopPlayback(ctx context.Context, guildID string, gp *entity.GuildPlayer) {
	player := s.lavalink.ExistingPlayer(guildID)
	if player != nil {
		_ = player.Update(ctx, lavalink.WithNullTrack())
	}
	s.lavalink.RemovePlayer(guildID)
	gp.IsPlaying = false
	gp.IsPaused = false
	gp.Queue = make([]*entity.Song, 0)
	gp.Current = 0
}

func (s *MusicService) joinVoice(guildID, channelID string) error {
	_, err := s.session.ChannelVoiceJoin(guildID, channelID, false, true) // not muted, deafened
	return err
}

func (s *MusicService) disconnectVoice(guildID string) {
	// Lavalink handles the voice connection, but we also need to tell Discord
	// to disconnect the gateway-side voice state.
	_ = s.session.ChannelVoiceJoinManual(guildID, "", false, false)
}

// advanceQueue moves to the next song in the queue based on loop mode.
// Returns the next song and its Lavalink track, or nil if queue is finished.
func (s *MusicService) advanceQueue(gp *entity.GuildPlayer) (*entity.Song, lavalink.Track) {
	// Save current to history
	if current := gp.NowPlaying(); current != nil {
		gp.History = append(gp.History, current)
		if len(gp.History) > 50 {
			gp.History = gp.History[len(gp.History)-50:]
		}
	}

	switch gp.LoopMode {
	case entity.LoopSong:
		// Replay the same song
		song := gp.NowPlaying()
		if song == nil {
			return nil, lavalink.Track{}
		}
		track := songToSearchQuery(song)
		return song, track

	case entity.LoopQueue:
		gp.Current++
		if gp.Current >= len(gp.Queue) {
			gp.Current = 0
		}
		song := gp.NowPlaying()
		if song == nil {
			return nil, lavalink.Track{}
		}
		track := songToSearchQuery(song)
		return song, track

	default: // LoopOff
		gp.Current++
		if gp.Current >= len(gp.Queue) {
			return nil, lavalink.Track{}
		}
		song := gp.NowPlaying()
		if song == nil {
			return nil, lavalink.Track{}
		}
		track := songToSearchQuery(song)
		return song, track
	}
}

// --- Lavalink event handlers ---

func (s *MusicService) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	guildID := player.GuildID().String()

	// Only advance if the track finished naturally or was replaced
	if event.Reason != lavalink.TrackEndReasonFinished && event.Reason != lavalink.TrackEndReasonLoadFailed {
		return
	}

	gp := s.getPlayer(guildID)
	if gp == nil {
		return
	}

	nextSong, _ := s.advanceQueue(gp)
	if nextSong == nil {
		// Try autoplay before stopping
		if gp.Autoplay {
			s.handleAutoplay(guildID, gp, player)
			return
		}

		// Queue finished
		gp.IsPlaying = false
		gp.IsPaused = false
		s.saveQueueToRedis(guildID, gp)

		// Notify in text channel
		if s.sendEmbed != nil && gp.TextChID != "" {
			s.sendEmbed(gp.TextChID, &discordgo.MessageEmbed{
				Description: "Queue finished! Add more songs with `/play`",
				Color:       0x00C9A7,
			})
		}

		// Disconnect if not in 24/7 mode
		ctx247, cancel247 := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel247()
		if !s.Is247(ctx247, guildID) {
			s.disconnectVoice(guildID)
			s.removePlayer(guildID)
			s.deleteQueueFromRedis(guildID)
		}
		return
	}

	// Play next track by searching again (since we don't store encoded track data)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := resolveQuery(nextSong.URL)
	result, err := s.lavalink.LoadTracks(ctx, query)
	if err != nil {
		s.logger.Error("Failed to load next track", "error", err, "song", nextSong.Title)
		return
	}

	tracks := extractTracks(result)
	if len(tracks) == 0 {
		s.logger.Warn("Next track not found, skipping", "song", nextSong.Title)
		// Try to skip to the next one
		s.onTrackEnd(player, event)
		return
	}

	if err := player.Update(ctx, lavalink.WithTrack(tracks[0])); err != nil {
		s.logger.Error("Failed to play next track", "error", err)
	}

	s.saveQueueToRedis(guildID, gp)
}

func (s *MusicService) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	guildID := player.GuildID().String()
	s.logger.Error("Track exception",
		"guild", guildID,
		"track", event.Track.Info.Title,
		"message", event.Exception.Message,
	)

	gp := s.getPlayer(guildID)
	if gp != nil && gp.TextChID != "" && s.sendEmbed != nil {
		s.sendEmbed(gp.TextChID, &discordgo.MessageEmbed{
			Description: fmt.Sprintf("Error playing **%s**: %s\nSkipping...", event.Track.Info.Title, event.Exception.Message),
			Color:       0xFF6B6B,
		})
	}
}

func (s *MusicService) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	guildID := player.GuildID().String()
	s.logger.Warn("Track stuck",
		"guild", guildID,
		"track", event.Track.Info.Title,
		"threshold", event.Threshold,
	)
}

// --- Redis queue persistence ---

func (s *MusicService) saveQueueToRedis(guildID string, gp *entity.GuildPlayer) {
	if s.redis == nil {
		return
	}
	data, err := json.Marshal(gp)
	if err != nil {
		s.logger.Error("Failed to marshal queue for Redis", "error", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = s.redis.Set(ctx, "music:queue:"+guildID, string(data), 24*time.Hour)
}

func (s *MusicService) deleteQueueFromRedis(guildID string) {
	if s.redis == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = s.redis.Delete(ctx, "music:queue:"+guildID)
}

// --- Track/song conversion helpers ---

func trackToSong(track lavalink.Track, requesterID string) *entity.Song {
	info := track.Info
	var url string
	if info.URI != nil {
		url = *info.URI
	}
	var thumbnail string
	if info.ArtworkURL != nil {
		thumbnail = *info.ArtworkURL
	}
	return &entity.Song{
		Title:       info.Title,
		URI:         track.Encoded,
		URL:         url,
		Duration:    int64(info.Length.Milliseconds()),
		Author:      info.Author,
		Thumbnail:   thumbnail,
		Source:      info.SourceName,
		RequesterID: requesterID,
		Identifier:  info.Identifier,
	}
}

func songToSearchQuery(song *entity.Song) lavalink.Track {
	// Return an empty track — the caller will re-search using the URL
	return lavalink.Track{}
}

func resolveQuery(query string) string {
	query = strings.TrimSpace(query)

	// If it's already a URL, use it directly
	if strings.HasPrefix(query, "http://") || strings.HasPrefix(query, "https://") {
		return query
	}

	// Otherwise, search YouTube
	return "ytsearch:" + query
}

func extractTracks(result *lavalink.LoadResult) []lavalink.Track {
	switch data := result.Data.(type) {
	case lavalink.Track:
		return []lavalink.Track{data}
	case lavalink.Search:
		return []lavalink.Track(data)
	case lavalink.Playlist:
		return data.Tracks
	default:
		return nil
	}
}

// --- Filter presets ---

// FilterPresets contains all available audio filter presets.
var FilterPresets = []string{"bassboost", "nightcore", "vaporwave", "karaoke", "8d", "tremolo", "clear"}

func getFilterPreset(name string) (lavalink.Filters, error) {
	switch strings.ToLower(name) {
	case "bassboost":
		eq := lavalink.Equalizer{}
		eq[0] = 0.6
		eq[1] = 0.67
		eq[2] = 0.67
		eq[3] = 0.4
		eq[4] = -0.5
		eq[5] = 0.15
		eq[6] = -0.45
		eq[7] = 0.23
		eq[8] = 0.35
		eq[9] = 0.45
		eq[10] = 0.55
		eq[11] = 0.6
		eq[12] = 0.55
		eq[13] = 0
		return lavalink.Filters{Equalizer: &eq}, nil

	case "nightcore":
		return lavalink.Filters{
			Timescale: &lavalink.Timescale{Speed: 1.3, Pitch: 1.3, Rate: 1.0},
		}, nil

	case "vaporwave":
		return lavalink.Filters{
			Timescale: &lavalink.Timescale{Speed: 0.8, Pitch: 0.8, Rate: 1.0},
			Tremolo:   &lavalink.Tremolo{Frequency: 14.0, Depth: 0.3},
		}, nil

	case "karaoke":
		return lavalink.Filters{
			Karaoke: &lavalink.Karaoke{Level: 1.0, MonoLevel: 1.0, FilterBand: 220.0, FilterWidth: 100.0},
		}, nil

	case "8d":
		return lavalink.Filters{
			Rotation: &lavalink.Rotation{RotationHz: 2},
		}, nil

	case "tremolo":
		return lavalink.Filters{
			Tremolo: &lavalink.Tremolo{Frequency: 10.0, Depth: 0.5},
		}, nil

	case "clear":
		return lavalink.Filters{}, nil

	default:
		return lavalink.Filters{}, fmt.Errorf("unknown filter: %s. Available: %s", name, strings.Join(FilterPresets, ", "))
	}
}

// --- Progress bar helper ---

// FormatProgressBar creates a visual progress bar string.
func FormatProgressBar(position, duration int64, length int) string {
	if duration <= 0 {
		return strings.Repeat("▬", length)
	}

	progress := float64(position) / float64(duration)
	if progress > 1 {
		progress = 1
	}
	if progress < 0 {
		progress = 0
	}

	filled := int(progress * float64(length))
	if filled > length {
		filled = length
	}

	bar := strings.Repeat("▬", filled) + "🔘" + strings.Repeat("▬", length-filled)
	return bar
}

// FormatDuration formats milliseconds to a human-readable string.
func FormatDuration(ms int64) string {
	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	if minutes >= 60 {
		hours := minutes / 60
		minutes = minutes % 60
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
