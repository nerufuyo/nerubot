package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/entity"
)

// RoastRepository handles roast data persistence
type RoastRepository struct {
	profiles   *JSONRepository
	patterns   *JSONRepository
	stats      *JSONRepository
	activities *JSONRepository
	
	profileData  map[string]*entity.UserProfile
	patternData  []*entity.RoastPattern
	statsData    map[string]*entity.ActivityStats
	historyData  []*entity.RoastHistory
	
	mu           sync.RWMutex
}

// NewRoastRepository creates a new roast repository
func NewRoastRepository() *RoastRepository {
	repo := &RoastRepository{
		profiles:    NewJSONRepository(config.RoastProfilesFile),
		patterns:    NewJSONRepository(config.RoastPatternsFile),
		stats:       NewJSONRepository(config.RoastStatsFile),
		activities:  NewJSONRepository(config.RoastActivitiesFile),
		profileData: make(map[string]*entity.UserProfile),
		patternData: make([]*entity.RoastPattern, 0),
		statsData:   make(map[string]*entity.ActivityStats),
		historyData: make([]*entity.RoastHistory, 0),
	}
	
	// Load existing data
	_ = repo.Load()
	
	// Initialize default patterns if empty
	if len(repo.patternData) == 0 {
		repo.initializeDefaultPatterns()
	}
	
	return repo
}

// Load loads all roast data from files
func (r *RoastRepository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Load profiles
	if err := r.profiles.Load(&r.profileData); err != nil {
		return fmt.Errorf("failed to load profiles: %w", err)
	}
	
	// Load patterns
	if err := r.patterns.Load(&r.patternData); err != nil {
		return fmt.Errorf("failed to load patterns: %w", err)
	}
	
	// Load stats
	if err := r.stats.Load(&r.statsData); err != nil {
		return fmt.Errorf("failed to load stats: %w", err)
	}
	
	// Load activities
	if err := r.activities.Load(&r.historyData); err != nil {
		return fmt.Errorf("failed to load activities: %w", err)
	}
	
	return nil
}

// GetUserProfile retrieves a user profile
func (r *RoastRepository) GetUserProfile(userID, guildID string) (*entity.UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	key := userID + ":" + guildID
	profile, exists := r.profileData[key]
	if !exists {
		return nil, fmt.Errorf("profile not found for user: %s", userID)
	}
	
	return profile, nil
}

// SaveUserProfile saves a user profile
func (r *RoastRepository) SaveUserProfile(profile *entity.UserProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	key := profile.UserID + ":" + profile.GuildID
	profile.UpdatedAt = time.Now()
	r.profileData[key] = profile
	
	return r.profiles.Save(r.profileData)
}

// GetOrCreateProfile gets or creates a user profile
func (r *RoastRepository) GetOrCreateProfile(userID, guildID, username string) (*entity.UserProfile, error) {
	profile, err := r.GetUserProfile(userID, guildID)
	if err != nil {
		// Create new profile
		profile = entity.NewUserProfile(userID, guildID, username)
		if err := r.SaveUserProfile(profile); err != nil {
			return nil, err
		}
	}
	return profile, nil
}

// GetActivityStats retrieves activity stats
func (r *RoastRepository) GetActivityStats(userID, guildID string) (*entity.ActivityStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	key := userID + ":" + guildID
	stats, exists := r.statsData[key]
	if !exists {
		return nil, fmt.Errorf("stats not found for user: %s", userID)
	}
	
	return stats, nil
}

// SaveActivityStats saves activity stats
func (r *RoastRepository) SaveActivityStats(stats *entity.ActivityStats) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	key := stats.UserID + ":" + stats.GuildID
	stats.UpdatedAt = time.Now()
	r.statsData[key] = stats
	
	return r.stats.Save(r.statsData)
}

// GetOrCreateStats gets or creates activity stats
func (r *RoastRepository) GetOrCreateStats(userID, guildID string) (*entity.ActivityStats, error) {
	stats, err := r.GetActivityStats(userID, guildID)
	if err != nil {
		// Create new stats
		stats = entity.NewActivityStats(userID, guildID)
		if err := r.SaveActivityStats(stats); err != nil {
			return nil, err
		}
	}
	return stats, nil
}

// GetRoastPatterns retrieves all roast patterns
func (r *RoastRepository) GetRoastPatterns() ([]*entity.RoastPattern, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.patternData, nil
}

// AddRoastHistory adds a roast to history
func (r *RoastRepository) AddRoastHistory(history *entity.RoastHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	history.ID = len(r.historyData) + 1
	history.CreatedAt = time.Now()
	r.historyData = append(r.historyData, history)
	
	return r.activities.Save(r.historyData)
}

// GetRoastHistory retrieves roast history for a user
func (r *RoastRepository) GetRoastHistory(userID, guildID string, limit int) ([]*entity.RoastHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var history []*entity.RoastHistory
	for i := len(r.historyData) - 1; i >= 0 && len(history) < limit; i-- {
		h := r.historyData[i]
		if h.TargetID == userID && h.GuildID == guildID {
			history = append(history, h)
		}
	}
	
	return history, nil
}

// initializeDefaultPatterns creates default roast patterns
func (r *RoastRepository) initializeDefaultPatterns() {
	r.patternData = []*entity.RoastPattern{
		{
			ID:       "night_owl",
			Name:     "Night Owl",
			Category: string(entity.RoastCategoryNightOwl),
			Templates: []string{
				"Do you even know what the sun looks like, %s?",
				"Your sleep schedule is more broken than your code!",
				"I bet vampires are jealous of your nocturnal lifestyle.",
			},
			Severity:    2,
			MinActivity: 50,
		},
		{
			ID:       "spammer",
			Name:     "Spammer",
			Category: string(entity.RoastCategorySpammer),
			Templates: []string{
				"%s, your fingers must be exhausted from all that typing!",
				"We get it, you have a keyboard. You don't need to prove it every second!",
				"Discord's servers are working overtime just to handle your messages.",
			},
			Severity:    3,
			MinActivity: 100,
		},
		{
			ID:       "lurker",
			Name:     "Lurker",
			Category: string(entity.RoastCategoryLurker),
			Templates: []string{
				"%s, I forgot you were even in this server!",
				"Is your keyboard broken, or are you just too cool to talk to us?",
				"The FBI could learn a thing or two from your surveillance techniques.",
			},
			Severity:    2,
			MinActivity: 10,
		},
		{
			ID:       "command_spam",
			Name:     "Command Spammer",
			Category: string(entity.RoastCategoryCommandSpam),
			Templates: []string{
				"%s treating me like a personal assistant. I'm not Siri!",
				"Do you get paid per command or something?",
				"My circuits are exhausted from your constant requests!",
			},
			Severity:    3,
			MinActivity: 50,
		},
	}
	
	_ = r.patterns.Save(r.patternData)
}
