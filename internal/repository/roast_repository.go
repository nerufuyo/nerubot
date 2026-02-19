package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// RoastRepository handles roast data persistence via MongoDB.
type RoastRepository struct {
	logger *logger.Logger
}

// NewRoastRepository creates a new roast repository.
func NewRoastRepository() *RoastRepository {
	repo := &RoastRepository{
		logger: logger.New("roast-repo"),
	}

	// Seed default patterns if collection is empty
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := repo.patterns().CountDocuments(ctx, bson.M{})
	if err == nil && count == 0 {
		repo.initializeDefaultPatterns()
	}

	return repo
}

func (r *RoastRepository) profiles() *mongo.Collection {
	return MongoDB.Collection("roast_profiles")
}

func (r *RoastRepository) patterns() *mongo.Collection {
	return MongoDB.Collection("roast_patterns")
}

func (r *RoastRepository) stats() *mongo.Collection {
	return MongoDB.Collection("roast_stats")
}

func (r *RoastRepository) history() *mongo.Collection {
	return MongoDB.Collection("roast_history")
}

// GetUserProfile retrieves a user profile.
func (r *RoastRepository) GetUserProfile(userID, guildID string) (*entity.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "guild_id": guildID}
	var profile entity.UserProfile
	err := r.profiles().FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("profile not found for user: %s", userID)
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	return &profile, nil
}

// SaveUserProfile upserts a user profile.
func (r *RoastRepository) SaveUserProfile(profile *entity.UserProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profile.UpdatedAt = time.Now()
	filter := bson.M{"user_id": profile.UserID, "guild_id": profile.GuildID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.profiles().ReplaceOne(ctx, filter, profile, opts)
	if err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}
	return nil
}

// GetOrCreateProfile gets or creates a user profile.
func (r *RoastRepository) GetOrCreateProfile(userID, guildID, username string) (*entity.UserProfile, error) {
	profile, err := r.GetUserProfile(userID, guildID)
	if err != nil {
		profile = entity.NewUserProfile(userID, guildID, username)
		if err := r.SaveUserProfile(profile); err != nil {
			return nil, err
		}
	}
	return profile, nil
}

// GetActivityStats retrieves activity stats.
func (r *RoastRepository) GetActivityStats(userID, guildID string) (*entity.ActivityStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "guild_id": guildID}
	var stats entity.ActivityStats
	err := r.stats().FindOne(ctx, filter).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("stats not found for user: %s", userID)
		}
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return &stats, nil
}

// SaveActivityStats upserts activity stats.
func (r *RoastRepository) SaveActivityStats(stats *entity.ActivityStats) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats.UpdatedAt = time.Now()
	filter := bson.M{"user_id": stats.UserID, "guild_id": stats.GuildID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.stats().ReplaceOne(ctx, filter, stats, opts)
	if err != nil {
		return fmt.Errorf("failed to save stats: %w", err)
	}
	return nil
}

// GetOrCreateStats gets or creates activity stats.
func (r *RoastRepository) GetOrCreateStats(userID, guildID string) (*entity.ActivityStats, error) {
	stats, err := r.GetActivityStats(userID, guildID)
	if err != nil {
		stats = entity.NewActivityStats(userID, guildID)
		if err := r.SaveActivityStats(stats); err != nil {
			return nil, err
		}
	}
	return stats, nil
}

// GetRoastPatterns retrieves all roast patterns.
func (r *RoastRepository) GetRoastPatterns() ([]*entity.RoastPattern, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.patterns().Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query patterns: %w", err)
	}
	defer cursor.Close(ctx)

	var result []*entity.RoastPattern
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AddRoastHistory adds a roast to history.
func (r *RoastRepository) AddRoastHistory(history *entity.RoastHistory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := MongoDB.GetNextSequence(ctx, "roast_history_id")
	if err != nil {
		return fmt.Errorf("failed to get next roast history ID: %w", err)
	}
	history.ID = id
	history.CreatedAt = time.Now()

	_, err = r.history().InsertOne(ctx, history)
	if err != nil {
		return fmt.Errorf("failed to save roast history: %w", err)
	}
	return nil
}

// GetRoastHistory retrieves roast history for a user (newest first, limited).
func (r *RoastRepository) GetRoastHistory(userID, guildID string, limit int) ([]*entity.RoastHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"target_id": userID, "guild_id": guildID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.history().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query roast history: %w", err)
	}
	defer cursor.Close(ctx)

	var result []*entity.RoastHistory
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// initializeDefaultPatterns seeds default roast patterns into MongoDB.
func (r *RoastRepository) initializeDefaultPatterns() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defaults := []interface{}{
		&entity.RoastPattern{
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
		&entity.RoastPattern{
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
		&entity.RoastPattern{
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
		&entity.RoastPattern{
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

	if _, err := r.patterns().InsertMany(ctx, defaults); err != nil {
		r.logger.Warn("Failed to seed default patterns", "error", err)
	}
}
