package entity

import "time"

// LoopMode represents the loop mode for music playback
type LoopMode int

const (
	LoopOff   LoopMode = iota // No looping
	LoopSong                  // Loop current song
	LoopQueue                 // Loop entire queue
)

// Song represents a music track
type Song struct {
	Title       string `json:"title" bson:"title"`
	URI         string `json:"uri" bson:"uri"`
	URL         string `json:"url" bson:"url"`           // User-facing URL (e.g., YouTube link)
	Duration    int64  `json:"duration" bson:"duration"` // Milliseconds
	Author      string `json:"author" bson:"author"`
	Thumbnail   string `json:"thumbnail" bson:"thumbnail"`
	Source      string `json:"source" bson:"source"` // youtube, spotify, soundcloud
	RequesterID string `json:"requesterId" bson:"requesterId"`
	Identifier  string `json:"identifier" bson:"identifier"` // Platform-specific ID
}

// FormatDuration returns a human-readable duration string (e.g., "3:45")
func (s *Song) FormatDuration() string {
	totalSeconds := s.Duration / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	if minutes >= 60 {
		hours := minutes / 60
		minutes = minutes % 60
		return formatTime(hours) + ":" + formatTime(minutes) + ":" + formatTime(seconds)
	}
	return formatTime(minutes) + ":" + formatTime(seconds)
}

func formatTime(n int64) string {
	if n < 10 {
		return "0" + intToStr(n)
	}
	return intToStr(n)
}

func intToStr(n int64) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// GuildPlayer represents the active music state for a guild
type GuildPlayer struct {
	GuildID   string          `json:"guildId"`
	ChannelID string          `json:"channelId"` // Voice channel ID
	TextChID  string          `json:"textChId"`  // Text channel for notifications
	Queue     []*Song         `json:"queue"`
	Current   int             `json:"current"` // Current song index in queue
	Volume    int             `json:"volume"`  // 0-200
	LoopMode  LoopMode        `json:"loopMode"`
	IsPlaying bool            `json:"isPlaying"`
	IsPaused  bool            `json:"isPaused"`
	History   []*Song         `json:"history"`  // Previously played songs
	Autoplay  bool            `json:"autoplay"` // Auto-queue similar songs when queue ends
	SkipVotes map[string]bool `json:"-"`        // User IDs that voted to skip (not persisted)
}

// NowPlaying returns the currently playing song, or nil if none
func (gp *GuildPlayer) NowPlaying() *Song {
	if gp.Current < 0 || gp.Current >= len(gp.Queue) {
		return nil
	}
	return gp.Queue[gp.Current]
}

// Remaining returns the number of songs remaining after current
func (gp *GuildPlayer) Remaining() int {
	if gp.Current >= len(gp.Queue) {
		return 0
	}
	return len(gp.Queue) - gp.Current - 1
}

// Playlist represents a user-saved playlist stored in MongoDB
type Playlist struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"userId" bson:"userId"`
	Name      string    `json:"name" bson:"name"`
	Songs     []*Song   `json:"songs" bson:"songs"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// GuildMusicSettings represents per-guild music configuration
type GuildMusicSettings struct {
	GuildID       string `json:"guildId" bson:"guildId"`
	DJRoleID      string `json:"djRoleId" bson:"djRoleId"`
	DefaultVolume int    `json:"defaultVolume" bson:"defaultVolume"`
	MaxQueueLen   int    `json:"maxQueueLen" bson:"maxQueueLen"`
	Stay247       bool   `json:"stay247" bson:"stay247"`
	AutoJoinChID  string `json:"autoJoinChId" bson:"autoJoinChId"` // Voice channel to auto-join
}
