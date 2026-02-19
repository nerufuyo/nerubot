package entity

import (
	"time"
)

// ServerStats represents analytics for a Discord server
type ServerStats struct {
	GuildID           string            `json:"guild_id" bson:"guild_id"`
	GuildName         string            `json:"guild_name" bson:"guild_name"`
	MemberCount       int               `json:"member_count" bson:"member_count"`
	CommandsUsed      int               `json:"commands_used" bson:"commands_used"`
	MessagesProcessed int               `json:"messages_processed" bson:"messages_processed"`
	SongsPlayed       int               `json:"songs_played" bson:"songs_played"`
	ConfessionsTotal  int               `json:"confessions_total" bson:"confessions_total"`
	RoastsGenerated   int               `json:"roasts_generated" bson:"roasts_generated"`
	ChatMessages      int               `json:"chat_messages" bson:"chat_messages"`
	NewsRequests      int               `json:"news_requests" bson:"news_requests"`
	WhaleAlerts       int               `json:"whale_alerts" bson:"whale_alerts"`
	TopCommands       map[string]int    `json:"top_commands" bson:"top_commands"`
	TopUsers          map[string]int    `json:"top_users" bson:"top_users"`
	ActiveDays        []string          `json:"active_days" bson:"active_days"`
	FirstSeen         time.Time         `json:"first_seen" bson:"first_seen"`
	LastActive        time.Time         `json:"last_active" bson:"last_active"`
	UpdatedAt         time.Time         `json:"updated_at" bson:"updated_at"`
}

// UserStats represents analytics for a Discord user
type UserStats struct {
	UserID            string         `json:"user_id" bson:"user_id"`
	Username          string         `json:"username" bson:"username"`
	CommandsUsed      int            `json:"commands_used" bson:"commands_used"`
	SongsRequested    int            `json:"songs_requested" bson:"songs_requested"`
	ConfessionsPosted int            `json:"confessions_posted" bson:"confessions_posted"`
	RoastsReceived    int            `json:"roasts_received" bson:"roasts_received"`
	ChatMessages      int            `json:"chat_messages" bson:"chat_messages"`
	NewsRequests      int            `json:"news_requests" bson:"news_requests"`
	WhaleChecks       int            `json:"whale_checks" bson:"whale_checks"`
	FavoriteCommands  map[string]int `json:"favorite_commands" bson:"favorite_commands"`
	FirstSeen         time.Time      `json:"first_seen" bson:"first_seen"`
	LastActive        time.Time      `json:"last_active" bson:"last_active"`
	UpdatedAt         time.Time      `json:"updated_at" bson:"updated_at"`
}

// CommandUsage tracks usage of a specific command
type CommandUsage struct {
	CommandName   string    `json:"command_name" bson:"command_name"`
	UserID        string    `json:"user_id" bson:"user_id"`
	GuildID       string    `json:"guild_id" bson:"guild_id"`
	ChannelID     string    `json:"channel_id" bson:"channel_id"`
	Success       bool      `json:"success" bson:"success"`
	ErrorMsg      string    `json:"error_msg,omitempty" bson:"error_msg,omitempty"`
	ExecutionTime int64     `json:"execution_time_ms" bson:"execution_time_ms"`
	Timestamp     time.Time `json:"timestamp" bson:"timestamp"`
}

// GlobalStats represents bot-wide statistics
type GlobalStats struct {
	TotalGuilds       int            `json:"total_guilds" bson:"total_guilds"`
	TotalUsers        int            `json:"total_users" bson:"total_users"`
	TotalCommands     int            `json:"total_commands" bson:"total_commands"`
	TotalSongsPlayed  int            `json:"total_songs_played" bson:"total_songs_played"`
	TotalConfessions  int            `json:"total_confessions" bson:"total_confessions"`
	TotalRoasts       int            `json:"total_roasts" bson:"total_roasts"`
	TotalChatMessages int            `json:"total_chat_messages" bson:"total_chat_messages"`
	Uptime            time.Duration  `json:"uptime" bson:"uptime"`
	TopGuilds         map[string]int `json:"top_guilds" bson:"top_guilds"`
	TopCommands       map[string]int `json:"top_commands" bson:"top_commands"`
	StartTime         time.Time      `json:"start_time" bson:"start_time"`
	UpdatedAt         time.Time      `json:"updated_at" bson:"updated_at"`
}

// NewServerStats creates a new ServerStats instance
func NewServerStats(guildID, guildName string) *ServerStats {
	return &ServerStats{
		GuildID:      guildID,
		GuildName:    guildName,
		TopCommands:  make(map[string]int),
		TopUsers:     make(map[string]int),
		ActiveDays:   make([]string, 0),
		FirstSeen:    time.Now(),
		LastActive:   time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// NewUserStats creates a new UserStats instance
func NewUserStats(userID, username string) *UserStats {
	return &UserStats{
		UserID:           userID,
		Username:         username,
		FavoriteCommands: make(map[string]int),
		FirstSeen:        time.Now(),
		LastActive:       time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// NewGlobalStats creates a new GlobalStats instance
func NewGlobalStats() *GlobalStats {
	return &GlobalStats{
		TopGuilds:   make(map[string]int),
		TopCommands: make(map[string]int),
		StartTime:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// RecordCommand records a command usage in server stats
func (s *ServerStats) RecordCommand(commandName, userID string) {
	s.CommandsUsed++
	s.TopCommands[commandName]++
	s.TopUsers[userID]++
	s.LastActive = time.Now()
	s.UpdatedAt = time.Now()
	
	// Add today to active days if not already present
	today := time.Now().Format("2006-01-02")
	found := false
	for _, day := range s.ActiveDays {
		if day == today {
			found = true
			break
		}
	}
	if !found {
		s.ActiveDays = append(s.ActiveDays, today)
	}
}

// RecordCommand records a command usage in user stats
func (u *UserStats) RecordCommand(commandName string) {
	u.CommandsUsed++
	u.FavoriteCommands[commandName]++
	u.LastActive = time.Now()
	u.UpdatedAt = time.Now()
}

// RecordSong increments song play count
func (s *ServerStats) RecordSong() {
	s.SongsPlayed++
	s.UpdatedAt = time.Now()
}

// RecordConfession increments confession count
func (s *ServerStats) RecordConfession() {
	s.ConfessionsTotal++
	s.UpdatedAt = time.Now()
}

// RecordRoast increments roast count
func (s *ServerStats) RecordRoast() {
	s.RoastsGenerated++
	s.UpdatedAt = time.Now()
}

// RecordChat increments chat message count
func (s *ServerStats) RecordChat() {
	s.ChatMessages++
	s.UpdatedAt = time.Now()
}

// GetMostUsedCommand returns the most frequently used command
func (s *ServerStats) GetMostUsedCommand() (string, int) {
	maxCmd := ""
	maxCount := 0
	for cmd, count := range s.TopCommands {
		if count > maxCount {
			maxCmd = cmd
			maxCount = count
		}
	}
	return maxCmd, maxCount
}

// GetMostActiveUser returns the most active user
func (s *ServerStats) GetMostActiveUser() (string, int) {
	maxUser := ""
	maxCount := 0
	for user, count := range s.TopUsers {
		if count > maxCount {
			maxUser = user
			maxCount = count
		}
	}
	return maxUser, maxCount
}

// GetFavoriteCommand returns the user's most used command
func (u *UserStats) GetFavoriteCommand() (string, int) {
	maxCmd := ""
	maxCount := 0
	for cmd, count := range u.FavoriteCommands {
		if count > maxCount {
			maxCmd = cmd
			maxCount = count
		}
	}
	return maxCmd, maxCount
}
