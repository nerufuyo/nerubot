package music

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/pkg/ytdlp"
)

// MusicService handles music playback operations
type MusicService struct {
	queues map[string]*entity.Queue // guildID -> Queue
	ytdlp  *ytdlp.YtDlp
	mu     sync.RWMutex
	logger *logger.Logger
}

// NewMusicService creates a new music service
func NewMusicService() (*MusicService, error) {
	yt, err := ytdlp.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize yt-dlp: %w", err)
	}

	return &MusicService{
		queues: make(map[string]*entity.Queue),
		ytdlp:  yt,
		logger: logger.New("music"),
	}, nil
}

// GetQueue gets or creates a queue for a guild
func (s *MusicService) GetQueue(guildID, voiceChannelID, textChannelID string) *entity.Queue {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		queue = entity.NewQueue(guildID, voiceChannelID, textChannelID)
		s.queues[guildID] = queue
	}

	return queue
}

// Search searches for songs
func (s *MusicService) Search(ctx context.Context, query string, maxResults int) ([]*entity.Song, error) {
	s.logger.Info("Searching for songs", "query", query, "max_results", maxResults)

	// Check if it's a URL or search query
	source := ytdlp.GetSource(query)
	
	if source == "unknown" {
		// It's a search query
		results, err := s.ytdlp.Search(ctx, query, maxResults)
		if err != nil {
			return nil, fmt.Errorf("search failed: %w", err)
		}

		songs := make([]*entity.Song, 0, len(results))
		for _, result := range results {
			song := s.convertToSong(&result, entity.SourceYouTube)
			songs = append(songs, song)
		}

		return songs, nil
	}

	// It's a URL - extract info
	opts := &ytdlp.ExtractOptions{
		NoPlaylist: true,
		Timeout:    15 * time.Second,
	}

	info, err := s.ytdlp.ExtractInfo(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to extract info: %w", err)
	}

	song := s.convertToSong(info, s.getSourceFromString(source))
	return []*entity.Song{song}, nil
}

// AddSong adds a song to the queue
func (s *MusicService) AddSong(ctx context.Context, guildID, voiceChannelID, textChannelID, query, requestedBy string) (*entity.Song, error) {
	queue := s.GetQueue(guildID, voiceChannelID, textChannelID)

	// Search for the song
	songs, err := s.Search(ctx, query, 1)
	if err != nil {
		return nil, err
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	song := songs[0]
	song.RequestedBy = requestedBy

	// Get stream URL
	streamURL, err := s.ytdlp.GetStreamURL(ctx, song.URL, &ytdlp.ExtractOptions{
		AudioOnly: true,
		Timeout:   15 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get stream URL: %w", err)
	}

	song.StreamURL = streamURL

	// Add to queue
	queue.Add(song)

	s.logger.Info("Song added to queue",
		"guild", guildID,
		"title", song.Title,
		"position", queue.Size(),
	)

	return song, nil
}

// Play starts playback
func (s *MusicService) Play(guildID string) error {
	s.mu.Lock()
	queue, exists := s.queues[guildID]
	if !exists {
		s.mu.Unlock()
		return fmt.Errorf("no queue for guild")
	}

	if queue.IsEmpty() {
		s.mu.Unlock()
		return fmt.Errorf("queue is empty")
	}

	if queue.CurrentIndex < 0 {
		queue.CurrentIndex = 0
	}

	// Get current song
	currentSong := queue.Current()
	if currentSong == nil {
		s.mu.Unlock()
		return fmt.Errorf("no current song")
	}

	queue.IsPlaying = true
	queue.IsPaused = false
	s.mu.Unlock()

	s.logger.Info("Playback started", "guild", guildID, "chat", queue.TextChannelID)
	return nil
}

// Pause pauses playback
func (s *MusicService) Pause(guildID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return fmt.Errorf("no queue for guild")
	}

	queue.IsPaused = true
	s.logger.Info("Playback paused", "guild", guildID)
	return nil
}

// Resume resumes playback
func (s *MusicService) Resume(guildID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return fmt.Errorf("no queue for guild")
	}

	queue.IsPaused = false
	s.logger.Info("Playback resumed", "guild", guildID)
	return nil
}

// Skip skips to the next song
func (s *MusicService) Skip(guildID string) (*entity.Song, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return nil, fmt.Errorf("no queue for guild")
	}

	next := queue.Skip()
	if next == nil {
		queue.IsPlaying = false
		return nil, fmt.Errorf("no more songs in queue")
	}

	s.logger.Info("Skipped to next song", "guild", guildID, "title", next.Title)
	return next, nil
}

// Stop stops playback and clears the queue
func (s *MusicService) Stop(guildID string) error {
	s.mu.Lock()
	queue, exists := s.queues[guildID]
	if !exists {
		s.mu.Unlock()
		return fmt.Errorf("no queue for guild")
	}

	queue.Clear()
	queue.IsPlaying = false
	queue.IsPaused = false
	s.mu.Unlock()

	s.logger.Info("Playback stopped", "guild", guildID)
	return nil
}

// SetLoopMode sets the loop mode
func (s *MusicService) SetLoopMode(guildID string, mode entity.LoopMode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return fmt.Errorf("no queue for guild")
	}

	queue.LoopMode = mode
	s.logger.Info("Loop mode changed", "guild", guildID, "mode", mode)
	return nil
}

// Shuffle shuffles the queue
func (s *MusicService) Shuffle(guildID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return fmt.Errorf("no queue for guild")
	}

	queue.Shuffle()
	s.logger.Info("Queue shuffled", "guild", guildID)
	return nil
}

// SetVolume sets the volume
func (s *MusicService) SetVolume(guildID string, volume float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return fmt.Errorf("no queue for guild")
	}

	if volume < 0 || volume > 1 {
		return fmt.Errorf("volume must be between 0 and 1")
	}

	queue.Volume = volume
	s.logger.Info("Volume changed", "guild", guildID, "volume", volume)
	return nil
}

// GetCurrentSong gets the current song
func (s *MusicService) GetCurrentSong(guildID string) (*entity.Song, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return nil, fmt.Errorf("no queue for guild")
	}

	return queue.Current(), nil
}

// RemoveFromQueue removes a song from the queue
func (s *MusicService) RemoveFromQueue(guildID string, index int) (*entity.Song, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return nil, fmt.Errorf("no queue for guild")
	}

	return queue.Remove(index), nil
}

// ClearQueue clears the queue
func (s *MusicService) ClearQueue(guildID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue, exists := s.queues[guildID]
	if !exists {
		return fmt.Errorf("no queue for guild")
	}

	queue.Clear()
	s.logger.Info("Queue cleared", "guild", guildID)
	return nil
}

// Helper functions

func (s *MusicService) convertToSong(info *ytdlp.VideoInfo, source entity.Source) *entity.Song {
	return &entity.Song{
		ID:        info.ID,
		Title:     info.Title,
		Artist:    info.Uploader,
		Duration:  time.Duration(info.Duration) * time.Second,
		URL:       info.Webpage,
		Thumbnail: info.Thumbnail,
		Source:    source,
		Webpage:   info.Webpage,
		Channel:   info.Channel,
		ViewCount: info.ViewCount,
		LikeCount: info.LikeCount,
		IsLive:    info.IsLive,
		AddedAt:   time.Now(),
	}
}

func (s *MusicService) getSourceFromString(source string) entity.Source {
	switch source {
	case "youtube":
		return entity.SourceYouTube
	case "spotify":
		return entity.SourceSpotify
	case "soundcloud":
		return entity.SourceSoundCloud
	case "direct":
		return entity.SourceDirect
	default:
		return entity.SourceUnknown
	}
}
