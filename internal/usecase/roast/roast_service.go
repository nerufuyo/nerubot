package roast

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/backend"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// RoastService handles roast operations
type RoastService struct {
	repo          *repository.RoastRepository
	stats         map[string]*entity.UserRoastStats
	logger        *logger.Logger
	backendClient *backend.Client
}

// NewRoastService creates a new roast service
func NewRoastService(backendClient *backend.Client) *RoastService {
	return &RoastService{
		repo:          repository.NewRoastRepository(),
		stats:         make(map[string]*entity.UserRoastStats),
		logger:        logger.New("roast"),
		backendClient: backendClient,
	}
}

// GenerateRoast generates a roast for a user
func (s *RoastService) GenerateRoast(ctx context.Context, userID, guildID, username string) (string, error) {
	// Get configurable values from dashboard settings
	cooldown := 5 * time.Minute
	minMessages := 10
	if s.backendClient != nil {
		rs := s.backendClient.GetSettings().RoastSettings
		if rs.CooldownMinutes > 0 {
			cooldown = time.Duration(rs.CooldownMinutes) * time.Minute
		}
		if rs.MinMessages > 0 {
			minMessages = rs.MinMessages
		}
	}

	// Check cooldown
	if s.isOnCooldown(userID, guildID) {
		remaining := s.getRemainingCooldown(userID, guildID)
		return "", fmt.Errorf("roast cooldown: %s remaining", remaining)
	}

	// Get or create profile
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return "", err
	}

	// Check if enough data
	if profile.MessageCount < minMessages {
		return "", fmt.Errorf("not enough data to roast. Need at least %d messages!", minMessages)
	}

	// Detect patterns
	categories := profile.DetectPatterns()
	if len(categories) == 0 {
		categories = []entity.RoastCategory{entity.RoastCategoryNormal}
	}

	// Pick a random category
	category := categories[rand.Intn(len(categories))]

	// Get patterns
	patterns, err := s.repo.GetRoastPatterns()
	if err != nil {
		return "", err
	}

	// Find matching pattern
	var roastTemplate string
	for _, pattern := range patterns {
		if pattern.Category == string(category) {
			// Pick random template
			if len(pattern.Templates) > 0 {
				roastTemplate = pattern.Templates[rand.Intn(len(pattern.Templates))]
				break
			}
		}
	}

	if roastTemplate == "" {
		roastTemplate = "You're so normal, %s, even I can't think of a good roast!"
	}

	// Format roast with username
	roast := fmt.Sprintf(roastTemplate, username)

	// Record roast
	history := &entity.RoastHistory{
		GuildID:     guildID,
		TargetID:    userID,
		RequestedBy: userID,
		Category:    category,
		Roast:       roast,
		Severity:    2,
	}
	if err := s.repo.AddRoastHistory(history); err != nil {
		s.logger.Warn("Failed to save roast history", "error", err)
	}

	// Set cooldown from dashboard settings
	s.setCooldown(userID, guildID, cooldown)

	s.logger.Info("Roast generated",
		"user", userID,
		"guild", guildID,
		"category", category,
	)

	return roast, nil
}

// TrackMessage tracks a message for roast analysis
func (s *RoastService) TrackMessage(userID, guildID, username, channelID string) error {
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return err
	}

	profile.RecordMessage(channelID, "")
	
	return s.repo.SaveUserProfile(profile)
}

// TrackVoiceActivity tracks voice activity
func (s *RoastService) TrackVoiceActivity(userID, guildID, username string, minutes int) error {
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return err
	}

	profile.RecordVoiceActivity(minutes)
	
	return s.repo.SaveUserProfile(profile)
}

// TrackCommand tracks command usage
func (s *RoastService) TrackCommand(userID, guildID, username string) error {
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return err
	}

	profile.RecordCommand()
	
	return s.repo.SaveUserProfile(profile)
}

// GetUserProfile gets a user's profile
func (s *RoastService) GetUserProfile(userID, guildID string) (*entity.UserProfile, error) {
	return s.repo.GetUserProfile(userID, guildID)
}

// GetActivityStats gets activity statistics
func (s *RoastService) GetActivityStats(userID, guildID string) (*entity.ActivityStats, error) {
	stats, err := s.repo.GetOrCreateStats(userID, guildID)
	if err != nil {
		return nil, err
	}

	// Update stats from profile
	profile, err := s.repo.GetUserProfile(userID, guildID)
	if err == nil {
		stats.TotalMessages = profile.MessageCount
		stats.TotalVoiceTime = profile.VoiceMinutes
		stats.TotalCommands = profile.CommandsUsed
		stats.LastActivity = profile.LastSeen
		stats.CalculateActivityScore()
		
		if err := s.repo.SaveActivityStats(stats); err != nil {
			s.logger.Warn("Failed to save stats", "error", err)
		}
	}

	return stats, nil
}

// GetRoastHistory gets roast history for a user
func (s *RoastService) GetRoastHistory(userID, guildID string, limit int) ([]*entity.RoastHistory, error) {
	return s.repo.GetRoastHistory(userID, guildID, limit)
}

// Cooldown helpers

func (s *RoastService) isOnCooldown(userID, guildID string) bool {
	key := userID + ":" + guildID
	stat, exists := s.stats[key]
	if !exists {
		return false
	}
	return stat.IsOnCooldown()
}

func (s *RoastService) getRemainingCooldown(userID, guildID string) time.Duration {
	key := userID + ":" + guildID
	stat, exists := s.stats[key]
	if !exists {
		return 0
	}
	return stat.RemainingCooldown()
}

func (s *RoastService) setCooldown(userID, guildID string, duration time.Duration) {
	key := userID + ":" + guildID
	stat, exists := s.stats[key]
	if !exists {
		stat = entity.NewUserRoastStats(userID, guildID)
		s.stats[key] = stat
	}
	stat.RecordRoast(entity.RoastCategoryNormal, duration)
}
