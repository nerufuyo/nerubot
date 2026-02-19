package entity

import (
	"time"
)

// UserProfile represents a user's Discord activity profile
type UserProfile struct {
	UserID           string                `json:"user_id" bson:"user_id"`
	GuildID          string                `json:"guild_id" bson:"guild_id"`
	Username         string                `json:"username" bson:"username"`
	MessageCount     int                   `json:"message_count" bson:"message_count"`
	VoiceMinutes     int                   `json:"voice_minutes" bson:"voice_minutes"`
	CommandsUsed     int                   `json:"commands_used" bson:"commands_used"`
	ReactionsGiven   int                   `json:"reactions_given" bson:"reactions_given"`
	MentionsReceived int                   `json:"mentions_received" bson:"mentions_received"`
	FirstSeen        time.Time             `json:"first_seen" bson:"first_seen"`
	LastSeen         time.Time             `json:"last_seen" bson:"last_seen"`
	ActivityHours    map[int]int           `json:"activity_hours" bson:"activity_hours"`    // Hour -> count
	ActiveDays       map[string]int        `json:"active_days" bson:"active_days"`       // Day name -> count
	ChannelActivity  map[string]int        `json:"channel_activity" bson:"channel_activity"`  // Channel ID -> count
	TopEmojis        map[string]int        `json:"top_emojis" bson:"top_emojis"`        // Emoji -> count
	Patterns         []string              `json:"patterns" bson:"patterns"`          // Detected patterns
	CreatedAt        time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at" bson:"updated_at"`
}

// RoastPattern represents a roast template category
type RoastPattern struct {
	ID          string   `json:"id" bson:"id"`
	Name        string   `json:"name" bson:"name"`
	Category    string   `json:"category" bson:"category"`
	Templates   []string `json:"templates" bson:"templates"`
	Conditions  []string `json:"conditions" bson:"conditions"`  // Conditions to trigger this pattern
	Severity    int      `json:"severity" bson:"severity"`    // 1-5, how harsh the roast is
	MinActivity int      `json:"min_activity" bson:"min_activity"` // Minimum activity required
}

// RoastCategory represents roast categories
type RoastCategory string

const (
	RoastCategoryNightOwl    RoastCategory = "night_owl"
	RoastCategorySpammer     RoastCategory = "spammer"
	RoastCategoryLurker      RoastCategory = "lurker"
	RoastCategoryCommandSpam RoastCategory = "command_spam"
	RoastCategoryEmojiAbuser RoastCategory = "emoji_abuser"
	RoastCategoryVoiceAddict RoastCategory = "voice_addict"
	RoastCategoryGhost       RoastCategory = "ghost"
	RoastCategoryNormal      RoastCategory = "normal"
)

// ActivityStats holds detailed activity statistics
type ActivityStats struct {
	GuildID         string             `json:"guild_id" bson:"guild_id"`
	UserID          string             `json:"user_id" bson:"user_id"`
	TotalMessages   int                `json:"total_messages" bson:"total_messages"`
	TotalVoiceTime  int                `json:"total_voice_time" bson:"total_voice_time"` // minutes
	TotalCommands   int                `json:"total_commands" bson:"total_commands"`
	AveragePerDay   float64            `json:"average_per_day" bson:"average_per_day"`
	MostActiveHour  int                `json:"most_active_hour" bson:"most_active_hour"`
	MostActiveDay   string             `json:"most_active_day" bson:"most_active_day"`
	LongestStreak   int                `json:"longest_streak" bson:"longest_streak"`   // days
	CurrentStreak   int                `json:"current_streak" bson:"current_streak"`
	ActivityScore   float64            `json:"activity_score" bson:"activity_score"`
	LastActivity    time.Time          `json:"last_activity" bson:"last_activity"`
	WeeklyActivity  map[string]int     `json:"weekly_activity" bson:"weekly_activity"`  // ISO week -> count
	MonthlyActivity map[string]int     `json:"monthly_activity" bson:"monthly_activity"` // YYYY-MM -> count
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

// RoastHistory represents a roast event
type RoastHistory struct {
	ID          int           `json:"id" bson:"id"`
	GuildID     string        `json:"guild_id" bson:"guild_id"`
	TargetID    string        `json:"target_id" bson:"target_id"`
	RequestedBy string        `json:"requested_by" bson:"requested_by"`
	Category    RoastCategory `json:"category" bson:"category"`
	Roast       string        `json:"roast" bson:"roast"`
	Severity    int           `json:"severity" bson:"severity"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
}

// UserRoastStats holds user-specific roast statistics
type UserRoastStats struct {
	UserID          string                   `json:"user_id" bson:"user_id"`
	GuildID         string                   `json:"guild_id" bson:"guild_id"`
	TimesRoasted    int                      `json:"times_roasted" bson:"times_roasted"`
	LastRoasted     *time.Time               `json:"last_roasted,omitempty" bson:"last_roasted,omitempty"`
	CooldownExpires *time.Time               `json:"cooldown_expires,omitempty" bson:"cooldown_expires,omitempty"`
	CategoryCounts  map[RoastCategory]int    `json:"category_counts" bson:"category_counts"`
	FavoriteRoast   string                   `json:"favorite_roast,omitempty" bson:"favorite_roast,omitempty"`
	UpdatedAt       time.Time                `json:"updated_at" bson:"updated_at"`
}

// NewUserProfile creates a new UserProfile instance
func NewUserProfile(userID, guildID, username string) *UserProfile {
	now := time.Now()
	return &UserProfile{
		UserID:          userID,
		GuildID:         guildID,
		Username:        username,
		MessageCount:    0,
		VoiceMinutes:    0,
		CommandsUsed:    0,
		ReactionsGiven:  0,
		MentionsReceived: 0,
		FirstSeen:       now,
		LastSeen:        now,
		ActivityHours:   make(map[int]int),
		ActiveDays:      make(map[string]int),
		ChannelActivity: make(map[string]int),
		TopEmojis:       make(map[string]int),
		Patterns:        make([]string, 0),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewActivityStats creates a new ActivityStats instance
func NewActivityStats(userID, guildID string) *ActivityStats {
	return &ActivityStats{
		GuildID:         guildID,
		UserID:          userID,
		TotalMessages:   0,
		TotalVoiceTime:  0,
		TotalCommands:   0,
		AveragePerDay:   0,
		MostActiveHour:  0,
		MostActiveDay:   "",
		LongestStreak:   0,
		CurrentStreak:   0,
		ActivityScore:   0,
		LastActivity:    time.Now(),
		WeeklyActivity:  make(map[string]int),
		MonthlyActivity: make(map[string]int),
		UpdatedAt:       time.Now(),
	}
}

// NewUserRoastStats creates a new UserRoastStats instance
func NewUserRoastStats(userID, guildID string) *UserRoastStats {
	return &UserRoastStats{
		UserID:         userID,
		GuildID:        guildID,
		TimesRoasted:   0,
		CategoryCounts: make(map[RoastCategory]int),
		UpdatedAt:      time.Now(),
	}
}

// RecordMessage records a message activity
func (p *UserProfile) RecordMessage(channelID string, emoji string) {
	p.MessageCount++
	p.LastSeen = time.Now()
	
	// Record hour
	hour := time.Now().Hour()
	p.ActivityHours[hour]++
	
	// Record day
	day := time.Now().Weekday().String()
	p.ActiveDays[day]++
	
	// Record channel
	if channelID != "" {
		p.ChannelActivity[channelID]++
	}
	
	// Record emoji if present
	if emoji != "" {
		p.TopEmojis[emoji]++
	}
	
	p.UpdatedAt = time.Now()
}

// RecordVoiceActivity records voice channel activity
func (p *UserProfile) RecordVoiceActivity(minutes int) {
	p.VoiceMinutes += minutes
	p.LastSeen = time.Now()
	p.UpdatedAt = time.Now()
}

// RecordCommand records a command usage
func (p *UserProfile) RecordCommand() {
	p.CommandsUsed++
	p.LastSeen = time.Now()
	p.UpdatedAt = time.Now()
}

// DetectPatterns analyzes the profile and detects behavior patterns
func (p *UserProfile) DetectPatterns() []RoastCategory {
	patterns := make([]RoastCategory, 0)
	
	// Night owl: most active between 10 PM and 4 AM
	nightActivity := 0
	for hour, count := range p.ActivityHours {
		if hour >= 22 || hour <= 4 {
			nightActivity += count
		}
	}
	if nightActivity > p.MessageCount/2 {
		patterns = append(patterns, RoastCategoryNightOwl)
	}
	
	// Spammer: high message count
	if p.MessageCount > 1000 {
		patterns = append(patterns, RoastCategorySpammer)
	}
	
	// Lurker: low message count but present
	if p.MessageCount < 50 && time.Since(p.FirstSeen).Hours() > 24*7 {
		patterns = append(patterns, RoastCategoryLurker)
	}
	
	// Command spammer
	if p.CommandsUsed > 100 {
		patterns = append(patterns, RoastCategoryCommandSpam)
	}
	
	// Voice addict
	if p.VoiceMinutes > 1000 {
		patterns = append(patterns, RoastCategoryVoiceAddict)
	}
	
	// Ghost: hasn't been seen in a while
	if time.Since(p.LastSeen).Hours() > 24*7 {
		patterns = append(patterns, RoastCategoryGhost)
	}
	
	// Normal if no patterns
	if len(patterns) == 0 {
		patterns = append(patterns, RoastCategoryNormal)
	}
	
	return patterns
}

// CalculateActivityScore calculates an activity score
func (s *ActivityStats) CalculateActivityScore() float64 {
	// Weight different activities
	messageScore := float64(s.TotalMessages) * 1.0
	voiceScore := float64(s.TotalVoiceTime) * 0.5
	commandScore := float64(s.TotalCommands) * 2.0
	streakScore := float64(s.LongestStreak) * 10.0
	
	s.ActivityScore = messageScore + voiceScore + commandScore + streakScore
	return s.ActivityScore
}

// IsOnCooldown checks if a user is on roast cooldown
func (s *UserRoastStats) IsOnCooldown() bool {
	if s.CooldownExpires == nil {
		return false
	}
	return time.Now().Before(*s.CooldownExpires)
}

// RemainingCooldown returns the remaining cooldown duration
func (s *UserRoastStats) RemainingCooldown() time.Duration {
	if !s.IsOnCooldown() {
		return 0
	}
	return time.Until(*s.CooldownExpires)
}

// RecordRoast records a roast event
func (s *UserRoastStats) RecordRoast(category RoastCategory, cooldown time.Duration) {
	s.TimesRoasted++
	now := time.Now()
	s.LastRoasted = &now
	
	// Set cooldown
	expires := now.Add(cooldown)
	s.CooldownExpires = &expires
	
	// Increment category count
	s.CategoryCounts[category]++
	
	s.UpdatedAt = time.Now()
}
