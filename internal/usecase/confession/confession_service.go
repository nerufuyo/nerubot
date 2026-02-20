package confession

import (
	"context"
	"fmt"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/backend"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// ConfessionService handles confession operations
type ConfessionService struct {
	repo          *repository.ConfessionRepository
	cooldowns     map[string]*entity.UserConfessionCooldown
	logger        *logger.Logger
	backendClient *backend.Client
}

// NewConfessionService creates a new confession service
func NewConfessionService(backendClient *backend.Client) *ConfessionService {
	return &ConfessionService{
		repo:          repository.NewConfessionRepository(),
		cooldowns:     make(map[string]*entity.UserConfessionCooldown),
		logger:        logger.New("confession"),
		backendClient: backendClient,
	}
}

// SubmitConfession submits a new confession
func (s *ConfessionService) SubmitConfession(ctx context.Context, guildID, authorID, content string) (*entity.Confession, error) {
	// Get guild-level settings
	settings, err := s.repo.GetSettings(guildID)
	if err != nil {
		// Create default settings
		settings = entity.NewGuildConfessionSettings(guildID, "")
		if err := s.repo.SaveSettings(settings); err != nil {
			return nil, fmt.Errorf("failed to create settings: %w", err)
		}
	}

	// Override guild settings with global dashboard settings if available
	if s.backendClient != nil {
		cs := s.backendClient.GetSettings().ConfessionSettings
		if cs.CooldownMinutes > 0 {
			settings.Cooldown = time.Duration(cs.CooldownMinutes) * time.Minute
		}
		if cs.MaxLength > 0 {
			settings.MaxLength = cs.MaxLength
		}
		settings.RequireApproval = cs.RequireApproval
		settings.AllowImages = cs.AllowImages
		settings.AllowReplies = cs.AllowReplies
	}

	// Check if enabled
	if !settings.Enabled {
		return nil, fmt.Errorf("confessions are disabled in this server")
	}

	// Check cooldown
	if s.isOnCooldown(authorID, guildID) {
		remaining := s.getRemainingCooldown(authorID, guildID)
		return nil, fmt.Errorf("please wait %s before confessing again", remaining)
	}

	// Validate content
	if content == "" {
		return nil, fmt.Errorf("confession cannot be empty")
	}

	if len(content) > settings.MaxLength {
		return nil, fmt.Errorf("confession is too long (max %d characters)", settings.MaxLength)
	}

	// Create confession
	confession := entity.NewConfession(guildID, authorID, content)

	// Set status based on approval requirement
	if settings.RequireApproval {
		confession.Status = entity.ConfessionStatusPending
	} else {
		confession.Status = entity.ConfessionStatusApproved
	}

	// Save confession
	if err := s.repo.SaveConfession(confession); err != nil {
		return nil, fmt.Errorf("failed to save confession: %w", err)
	}

	// Set cooldown
	s.setCooldown(authorID, guildID, settings.Cooldown)

	// Add to queue if approved
	if !settings.RequireApproval {
		queue := &entity.ConfessionQueue{
			Confession: confession,
			Priority:   0,
			AddedAt:    time.Now(),
		}
		if err := s.repo.AddToQueue(queue); err != nil {
			s.logger.Warn("Failed to add to queue", "error", err)
		}
	}

	s.logger.Info("Confession submitted",
		"id", confession.ID,
		"guild", guildID,
		"requires_approval", settings.RequireApproval,
	)

	return confession, nil
}

// ApproveConfession approves a confession
func (s *ConfessionService) ApproveConfession(confessionID int, moderatorID string) error {
	confession, err := s.repo.GetConfession(confessionID)
	if err != nil {
		return err
	}

	confession.Approve(moderatorID)

	if err := s.repo.SaveConfession(confession); err != nil {
		return err
	}

	// Add to queue
	queue := &entity.ConfessionQueue{
		Confession: confession,
		Priority:   0,
		AddedAt:    time.Now(),
	}
	if err := s.repo.AddToQueue(queue); err != nil {
		s.logger.Warn("Failed to add to queue", "error", err)
	}

	s.logger.Info("Confession approved", "id", confessionID, "moderator", moderatorID)
	return nil
}

// RejectConfession rejects a confession
func (s *ConfessionService) RejectConfession(confessionID int, moderatorID string) error {
	confession, err := s.repo.GetConfession(confessionID)
	if err != nil {
		return err
	}

	confession.Reject(moderatorID)

	if err := s.repo.SaveConfession(confession); err != nil {
		return err
	}

	s.logger.Info("Confession rejected", "id", confessionID, "moderator", moderatorID)
	return nil
}

// PostConfession marks a confession as posted
func (s *ConfessionService) PostConfession(confessionID int, messageID, channelID string) error {
	confession, err := s.repo.GetConfession(confessionID)
	if err != nil {
		return err
	}

	confession.Post(messageID, channelID)

	if err := s.repo.SaveConfession(confession); err != nil {
		return err
	}

	// Remove from queue
	if err := s.repo.RemoveFromQueue(confession); err != nil {
		s.logger.Warn("Failed to remove from queue", "error", err)
	}

	s.logger.Info("Confession posted", "id", confessionID, "message", messageID)
	return nil
}

// ReplyToConfession adds a reply to a confession
func (s *ConfessionService) ReplyToConfession(confessionID int, guildID, authorID, content string) (*entity.ConfessionReply, error) {
	// Get settings
	settings, err := s.repo.GetSettings(guildID)
	if err != nil {
		return nil, fmt.Errorf("confession system not setup")
	}

	if !settings.AllowReplies {
		return nil, fmt.Errorf("replies are disabled")
	}

	// Get confession
	confession, err := s.repo.GetConfession(confessionID)
	if err != nil {
		return nil, err
	}

	if !confession.IsPosted() {
		return nil, fmt.Errorf("cannot reply to unposted confession")
	}

	// Create reply
	reply := entity.NewConfessionReply(confessionID, guildID, authorID, content)

	// Save reply
	if err := s.repo.SaveReply(reply); err != nil {
		return nil, fmt.Errorf("failed to save reply: %w", err)
	}

	// Increment reply count
	confession.IncrementReplyCount()
	if err := s.repo.SaveConfession(confession); err != nil {
		s.logger.Warn("Failed to update reply count", "error", err)
	}

	s.logger.Info("Reply added", "confession_id", confessionID, "reply_id", reply.ID)
	return reply, nil
}

// GetPendingConfessions gets pending confessions for moderation
func (s *ConfessionService) GetPendingConfessions(guildID string) ([]*entity.Confession, error) {
	return s.repo.GetPendingConfessions(guildID)
}

// GetConfession gets a confession by ID
func (s *ConfessionService) GetConfession(id int) (*entity.Confession, error) {
	return s.repo.GetConfession(id)
}

// GetReplies gets replies for a confession
func (s *ConfessionService) GetReplies(confessionID int) ([]*entity.ConfessionReply, error) {
	return s.repo.GetReplies(confessionID)
}

// GetSettings gets confession settings for a guild
func (s *ConfessionService) GetSettings(guildID string) (*entity.GuildConfessionSettings, error) {
	settings, err := s.repo.GetSettings(guildID)
	if err != nil {
		// Return default settings
		return entity.NewGuildConfessionSettings(guildID, ""), nil
	}
	return settings, nil
}

// UpdateSettings updates confession settings
func (s *ConfessionService) UpdateSettings(settings *entity.GuildConfessionSettings) error {
	return s.repo.SaveSettings(settings)
}

// GetQueue gets the confession queue
func (s *ConfessionService) GetQueue() ([]*entity.ConfessionQueue, error) {
	return s.repo.GetQueue()
}

// Cooldown helpers

func (s *ConfessionService) isOnCooldown(userID, guildID string) bool {
	key := userID + ":" + guildID
	cooldown, exists := s.cooldowns[key]
	if !exists {
		return false
	}
	return cooldown.IsOnCooldown()
}

func (s *ConfessionService) getRemainingCooldown(userID, guildID string) time.Duration {
	key := userID + ":" + guildID
	cooldown, exists := s.cooldowns[key]
	if !exists {
		return 0
	}
	return cooldown.RemainingCooldown()
}

func (s *ConfessionService) setCooldown(userID, guildID string, duration time.Duration) {
	key := userID + ":" + guildID
	cooldown, exists := s.cooldowns[key]
	if !exists {
		cooldown = &entity.UserConfessionCooldown{
			UserID:  userID,
			GuildID: guildID,
		}
		s.cooldowns[key] = cooldown
	}
	cooldown.SetCooldown(duration)
}
