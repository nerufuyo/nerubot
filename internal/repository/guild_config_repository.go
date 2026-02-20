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

// GuildConfigRepository handles guild configuration persistence via MongoDB.
type GuildConfigRepository struct {
	logger *logger.Logger
}

// NewGuildConfigRepository creates a new guild config repository.
func NewGuildConfigRepository() *GuildConfigRepository {
	return &GuildConfigRepository{
		logger: logger.New("guild-config-repo"),
	}
}

func (r *GuildConfigRepository) collection() *mongo.Collection {
	return MongoDB.Collection("guild_configs")
}

// Save upserts a guild configuration.
func (r *GuildConfigRepository) Save(config *entity.GuildConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config.UpdatedAt = time.Now()

	filter := bson.M{"guild_id": config.GuildID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.collection().ReplaceOne(ctx, filter, config, opts)
	if err != nil {
		return fmt.Errorf("failed to save guild config: %w", err)
	}
	return nil
}

// Get retrieves a guild configuration by guild ID.
func (r *GuildConfigRepository) Get(guildID string) (*entity.GuildConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var config entity.GuildConfig
	err := r.collection().FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // not found, not an error
		}
		return nil, fmt.Errorf("failed to get guild config: %w", err)
	}
	return &config, nil
}

// GetAll retrieves all guild configurations.
func (r *GuildConfigRepository) GetAll() ([]*entity.GuildConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query guild configs: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []*entity.GuildConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, fmt.Errorf("failed to decode guild configs: %w", err)
	}
	return configs, nil
}

// Delete removes a guild configuration.
func (r *GuildConfigRepository) Delete(guildID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection().DeleteOne(ctx, bson.M{"guild_id": guildID})
	if err != nil {
		return fmt.Errorf("failed to delete guild config: %w", err)
	}
	return nil
}
