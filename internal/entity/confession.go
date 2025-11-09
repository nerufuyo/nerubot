package entity

import (
	"time"
)

// ConfessionStatus represents the status of a confession
type ConfessionStatus string

const (
	ConfessionStatusPending  ConfessionStatus = "pending"
	ConfessionStatusApproved ConfessionStatus = "approved"
	ConfessionStatusRejected ConfessionStatus = "rejected"
	ConfessionStatusPosted   ConfessionStatus = "posted"
)

// Confession represents an anonymous confession
type Confession struct {
	ID           int              `json:"id"`
	GuildID      string           `json:"guild_id"`
	AuthorID     string           `json:"author_id"` // Anonymous, stored for moderation
	Content      string           `json:"content"`
	ImageURL     string           `json:"image_url,omitempty"`
	Status       ConfessionStatus `json:"status"`
	MessageID    string           `json:"message_id,omitempty"` // Discord message ID when posted
	ChannelID    string           `json:"channel_id,omitempty"`
	ReplyCount   int              `json:"reply_count"`
	CreatedAt    time.Time        `json:"created_at"`
	PostedAt     *time.Time       `json:"posted_at,omitempty"`
	ModeratedBy  string           `json:"moderated_by,omitempty"`
	ModeratedAt  *time.Time       `json:"moderated_at,omitempty"`
}

// ConfessionReply represents a reply to a confession
type ConfessionReply struct {
	ID            int       `json:"id"`
	ConfessionID  int       `json:"confession_id"`
	GuildID       string    `json:"guild_id"`
	AuthorID      string    `json:"author_id"` // Anonymous
	Content       string    `json:"content"`
	MessageID     string    `json:"message_id,omitempty"`
	ThreadID      string    `json:"thread_id,omitempty"` // Discord thread ID
	CreatedAt     time.Time `json:"created_at"`
}

// GuildConfessionSettings holds confession settings for a guild
type GuildConfessionSettings struct {
	GuildID             string        `json:"guild_id"`
	Enabled             bool          `json:"enabled"`
	ChannelID           string        `json:"channel_id"`
	RequireApproval     bool          `json:"require_approval"`
	ModeratorRoleID     string        `json:"moderator_role_id,omitempty"`
	AllowImages         bool          `json:"allow_images"`
	AllowReplies        bool          `json:"allow_replies"`
	Cooldown            time.Duration `json:"cooldown"`
	MaxLength           int           `json:"max_length"`
	AllowedRoleIDs      []string      `json:"allowed_role_ids,omitempty"`
	BannedUserIDs       []string      `json:"banned_user_ids,omitempty"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

// ConfessionQueue represents a queued confession waiting to be posted
type ConfessionQueue struct {
	Confession *Confession `json:"confession"`
	Priority   int         `json:"priority"`
	AddedAt    time.Time   `json:"added_at"`
}

// UserConfessionCooldown tracks user cooldowns
type UserConfessionCooldown struct {
	UserID    string    `json:"user_id"`
	GuildID   string    `json:"guild_id"`
	LastPost  time.Time `json:"last_post"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ConfessionStats holds statistics for confessions
type ConfessionStats struct {
	GuildID          string    `json:"guild_id"`
	TotalConfessions int       `json:"total_confessions"`
	TotalReplies     int       `json:"total_replies"`
	PendingCount     int       `json:"pending_count"`
	ApprovedCount    int       `json:"approved_count"`
	RejectedCount    int       `json:"rejected_count"`
	PostedCount      int       `json:"posted_count"`
	TopContributor   string    `json:"top_contributor,omitempty"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NewConfession creates a new Confession instance
func NewConfession(guildID, authorID, content string) *Confession {
	return &Confession{
		GuildID:    guildID,
		AuthorID:   authorID,
		Content:    content,
		Status:     ConfessionStatusPending,
		ReplyCount: 0,
		CreatedAt:  time.Now(),
	}
}

// NewConfessionReply creates a new ConfessionReply instance
func NewConfessionReply(confessionID int, guildID, authorID, content string) *ConfessionReply {
	return &ConfessionReply{
		ConfessionID: confessionID,
		GuildID:      guildID,
		AuthorID:     authorID,
		Content:      content,
		CreatedAt:    time.Now(),
	}
}

// NewGuildConfessionSettings creates new settings with defaults
func NewGuildConfessionSettings(guildID, channelID string) *GuildConfessionSettings {
	return &GuildConfessionSettings{
		GuildID:         guildID,
		Enabled:         true,
		ChannelID:       channelID,
		RequireApproval: false,
		AllowImages:     true,
		AllowReplies:    true,
		Cooldown:        10 * time.Minute,
		MaxLength:       2000,
		AllowedRoleIDs:  make([]string, 0),
		BannedUserIDs:   make([]string, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// Approve approves the confession
func (c *Confession) Approve(moderatorID string) {
	c.Status = ConfessionStatusApproved
	c.ModeratedBy = moderatorID
	now := time.Now()
	c.ModeratedAt = &now
}

// Reject rejects the confession
func (c *Confession) Reject(moderatorID string) {
	c.Status = ConfessionStatusRejected
	c.ModeratedBy = moderatorID
	now := time.Now()
	c.ModeratedAt = &now
}

// Post marks the confession as posted
func (c *Confession) Post(messageID, channelID string) {
	c.Status = ConfessionStatusPosted
	c.MessageID = messageID
	c.ChannelID = channelID
	now := time.Now()
	c.PostedAt = &now
}

// IncrementReplyCount increments the reply count
func (c *Confession) IncrementReplyCount() {
	c.ReplyCount++
}

// IsPending returns true if confession is pending
func (c *Confession) IsPending() bool {
	return c.Status == ConfessionStatusPending
}

// IsApproved returns true if confession is approved
func (c *Confession) IsApproved() bool {
	return c.Status == ConfessionStatusApproved
}

// IsRejected returns true if confession is rejected
func (c *Confession) IsRejected() bool {
	return c.Status == ConfessionStatusRejected
}

// IsPosted returns true if confession is posted
func (c *Confession) IsPosted() bool {
	return c.Status == ConfessionStatusPosted
}

// IsUserBanned checks if a user is banned
func (s *GuildConfessionSettings) IsUserBanned(userID string) bool {
	for _, id := range s.BannedUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// IsUserAllowed checks if a user is allowed (if role restrictions exist)
func (s *GuildConfessionSettings) IsUserAllowed(userID string, userRoleIDs []string) bool {
	// If no role restrictions, everyone is allowed
	if len(s.AllowedRoleIDs) == 0 {
		return true
	}
	
	// Check if user has any of the allowed roles
	for _, allowedRole := range s.AllowedRoleIDs {
		for _, userRole := range userRoleIDs {
			if allowedRole == userRole {
				return true
			}
		}
	}
	
	return false
}

// BanUser adds a user to the banned list
func (s *GuildConfessionSettings) BanUser(userID string) {
	if !s.IsUserBanned(userID) {
		s.BannedUserIDs = append(s.BannedUserIDs, userID)
		s.UpdatedAt = time.Now()
	}
}

// UnbanUser removes a user from the banned list
func (s *GuildConfessionSettings) UnbanUser(userID string) {
	for i, id := range s.BannedUserIDs {
		if id == userID {
			s.BannedUserIDs = append(s.BannedUserIDs[:i], s.BannedUserIDs[i+1:]...)
			s.UpdatedAt = time.Now()
			break
		}
	}
}

// IsOnCooldown checks if a user is on cooldown
func (c *UserConfessionCooldown) IsOnCooldown() bool {
	return time.Now().Before(c.ExpiresAt)
}

// RemainingCooldown returns the remaining cooldown duration
func (c *UserConfessionCooldown) RemainingCooldown() time.Duration {
	if !c.IsOnCooldown() {
		return 0
	}
	return time.Until(c.ExpiresAt)
}

// SetCooldown sets a new cooldown
func (c *UserConfessionCooldown) SetCooldown(duration time.Duration) {
	c.LastPost = time.Now()
	c.ExpiresAt = time.Now().Add(duration)
}
