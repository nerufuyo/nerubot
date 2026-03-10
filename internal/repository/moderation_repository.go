package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// ModerationRepository handles moderation data persistence via MongoDB.
type ModerationRepository struct {
	logger *logger.Logger
}

// NewModerationRepository creates a new moderation repository.
func NewModerationRepository() *ModerationRepository {
	return &ModerationRepository{
		logger: logger.New("moderation-repo"),
	}
}

func (r *ModerationRepository) warnings() *mongo.Collection {
	return MongoDB.Collection("warnings")
}

// AddWarning saves a new warning.
func (r *ModerationRepository) AddWarning(w *entity.Warning) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.warnings().InsertOne(ctx, w)
	return err
}

// GetWarnings returns all warnings for a user in a guild.
func (r *ModerationRepository) GetWarnings(guildID, userID string) ([]*entity.Warning, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"guild_id": guildID, "user_id": userID}
	cursor, err := r.warnings().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var warnings []*entity.Warning
	if err := cursor.All(ctx, &warnings); err != nil {
		return nil, err
	}
	return warnings, nil
}

// ClearWarnings removes all warnings for a user in a guild.
func (r *ModerationRepository) ClearWarnings(guildID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"guild_id": guildID, "user_id": userID}
	_, err := r.warnings().DeleteMany(ctx, filter)
	return err
}
