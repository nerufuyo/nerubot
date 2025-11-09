package entity
package entity

import (
	"fmt"
	"time"
)

// Source represents a music source type
type Source string

const (
	SourceYouTube    Source = "youtube"
	SourceSpotify    Source = "spotify"
	SourceSoundCloud Source = "soundcloud"
	SourceDirect     Source = "direct"
	SourceUnknown    Source = "unknown"
)

// LoopMode represents the queue loop mode
type LoopMode string

const (
	LoopModeOff    LoopMode = "off"
	LoopModeSingle LoopMode = "single"
	LoopModeQueue  LoopMode = "queue"
)

// Song represents a song/audio track
type Song struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Artist      string        `json:"artist"`
	Duration    time.Duration `json:"duration"`
	URL         string        `json:"url"`
	StreamURL   string        `json:"stream_url"`
	Thumbnail   string        `json:"thumbnail"`
	Source      Source        `json:"source"`
	RequestedBy string        `json:"requested_by"` // Discord user ID
	Webpage     string        `json:"webpage"`
	Channel     string        `json:"channel"`
	ViewCount   int64         `json:"view_count"`
	LikeCount   int64         `json:"like_count"`
	IsLive      bool          `json:"is_live"`
	AddedAt     time.Time     `json:"added_at"`
}

// Queue represents a music queue for a guild
type Queue struct {
	GuildID       string     `json:"guild_id"`
	Songs         []*Song    `json:"songs"`
	CurrentIndex  int        `json:"current_index"`
	LoopMode      LoopMode   `json:"loop_mode"`
	IsShuffled    bool       `json:"is_shuffled"`
	VoiceChannelID string    `json:"voice_channel_id"`
	TextChannelID string     `json:"text_channel_id"`
	Volume        float64    `json:"volume"`
	IsPaused      bool       `json:"is_paused"`
	IsPlaying     bool       `json:"is_playing"`
	Mode247       bool       `json:"mode_24_7"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Playlist represents a playlist from various sources
type Playlist struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Source   Source   `json:"source"`
	Songs    []*Song  `json:"songs"`
	URL      string   `json:"url"`
	Thumbnail string  `json:"thumbnail"`
	Author   string   `json:"author"`
	Count    int      `json:"count"`
}

// SearchResult represents a search result
type SearchResult struct {
	Query     string    `json:"query"`
	Results   []*Song   `json:"results"`
	Source    Source    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

// NewSong creates a new Song instance
func NewSong(title, artist, url string, duration time.Duration, source Source) *Song {
	return &Song{
		Title:    title,
		Artist:   artist,
		URL:      url,
		Duration: duration,
		Source:   source,
		AddedAt:  time.Now(),
	}
}

// NewQueue creates a new Queue instance
func NewQueue(guildID, voiceChannelID, textChannelID string) *Queue {
	return &Queue{
		GuildID:        guildID,
		Songs:          make([]*Song, 0),
		CurrentIndex:   -1,
		LoopMode:       LoopModeOff,
		IsShuffled:     false,
		VoiceChannelID: voiceChannelID,
		TextChannelID:  textChannelID,
		Volume:         0.5,
		IsPaused:       false,
		IsPlaying:      false,
		Mode247:        false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// Add adds a song to the queue
func (q *Queue) Add(song *Song) {
	q.Songs = append(q.Songs, song)
	q.UpdatedAt = time.Now()
}

// AddMultiple adds multiple songs to the queue
func (q *Queue) AddMultiple(songs []*Song) {
	q.Songs = append(q.Songs, songs...)
	q.UpdatedAt = time.Now()
}

// Remove removes a song at the given index
func (q *Queue) Remove(index int) *Song {
	if index < 0 || index >= len(q.Songs) {
		return nil
	}
	
	song := q.Songs[index]
	q.Songs = append(q.Songs[:index], q.Songs[index+1:]...)
	
	// Adjust current index if needed
	if index < q.CurrentIndex {
		q.CurrentIndex--
	} else if index == q.CurrentIndex && q.CurrentIndex >= len(q.Songs) {
		q.CurrentIndex = len(q.Songs) - 1
	}
	
	q.UpdatedAt = time.Now()
	return song
}

// Clear clears all songs from the queue
func (q *Queue) Clear() {
	q.Songs = make([]*Song, 0)
	q.CurrentIndex = -1
	q.UpdatedAt = time.Now()
}

// Current returns the current song
func (q *Queue) Current() *Song {
	if q.CurrentIndex < 0 || q.CurrentIndex >= len(q.Songs) {
		return nil
	}
	return q.Songs[q.CurrentIndex]
}

// Next returns the next song in the queue
func (q *Queue) Next() *Song {
	if len(q.Songs) == 0 {
		return nil
	}
	
	switch q.LoopMode {
	case LoopModeSingle:
		// Stay on current song
		return q.Current()
		
	case LoopModeQueue:
		// Loop back to start
		q.CurrentIndex = (q.CurrentIndex + 1) % len(q.Songs)
		
	default: // LoopModeOff
		q.CurrentIndex++
		if q.CurrentIndex >= len(q.Songs) {
			return nil
		}
	}
	
	q.UpdatedAt = time.Now()
	return q.Current()
}

// Previous returns the previous song in the queue
func (q *Queue) Previous() *Song {
	if len(q.Songs) == 0 {
		return nil
	}
	
	q.CurrentIndex--
	if q.CurrentIndex < 0 {
		if q.LoopMode == LoopModeQueue {
			q.CurrentIndex = len(q.Songs) - 1
		} else {
			q.CurrentIndex = 0
		}
	}
	
	q.UpdatedAt = time.Now()
	return q.Current()
}

// Skip skips the current song
func (q *Queue) Skip() *Song {
	if q.LoopMode == LoopModeSingle {
		// Temporarily disable single loop to skip
		q.LoopMode = LoopModeOff
		defer func() { q.LoopMode = LoopModeSingle }()
	}
	return q.Next()
}

// IsEmpty returns true if the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(q.Songs) == 0
}

// Size returns the number of songs in the queue
func (q *Queue) Size() int {
	return len(q.Songs)
}

// Shuffle shuffles the queue (excluding current song)
func (q *Queue) Shuffle() {
	if len(q.Songs) <= 1 {
		return
	}
	
	// Keep current song in place, shuffle the rest
	current := q.Current()
	if current == nil {
		// No current song, shuffle all
		shuffleSongs(q.Songs)
	} else {
		// Shuffle songs after current
		if q.CurrentIndex+1 < len(q.Songs) {
			shuffleSongs(q.Songs[q.CurrentIndex+1:])
		}
	}
	
	q.IsShuffled = true
	q.UpdatedAt = time.Now()
}

// Remaining returns the number of songs remaining in queue
func (q *Queue) Remaining() int {
	if q.CurrentIndex < 0 {
		return len(q.Songs)
	}
	return len(q.Songs) - q.CurrentIndex - 1
}

// TotalDuration returns the total duration of all songs in queue
func (q *Queue) TotalDuration() time.Duration {
	var total time.Duration
	for _, song := range q.Songs {
		total += song.Duration
	}
	return total
}

// RemainingDuration returns the total duration of remaining songs
func (q *Queue) RemainingDuration() time.Duration {
	if q.CurrentIndex < 0 {
		return q.TotalDuration()
	}
	
	var total time.Duration
	for i := q.CurrentIndex + 1; i < len(q.Songs); i++ {
		total += q.Songs[i].Duration
	}
	return total
}

// Helper function to shuffle songs
func shuffleSongs(songs []*Song) {
	// Fisher-Yates shuffle
	for i := len(songs) - 1; i > 0; i-- {
		j := int(time.Now().UnixNano()) % (i + 1)
		songs[i], songs[j] = songs[j], songs[i]
	}
}

// String returns the source name
func (s Source) String() string {
	return string(s)
}

// Emoji returns the emoji for the source
func (s Source) Emoji() string {
	switch s {
	case SourceYouTube:
		return "▶️"
	case SourceSpotify:
		return "💚"
	case SourceSoundCloud:
		return "🧡"
	case SourceDirect:
		return "🔗"
	default:
		return "🎵"
	}
}

// FormatDuration formats a duration as HH:MM:SS or MM:SS
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
