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

// ConfessionRepository handles confession data persistence via MongoDB.
type ConfessionRepository struct {
	logger *logger.Logger
}

// NewConfessionRepository creates a new confession repository.
func NewConfessionRepository() *ConfessionRepository {
	return &ConfessionRepository{
		logger: logger.New("confession-repo"),
	}
}

func (r *ConfessionRepository) confessions() *mongo.Collection {
	return MongoDB.Collection("confessions")
}

func (r *ConfessionRepository) replies() *mongo.Collection {
	return MongoDB.Collection("confession_replies")
}

func (r *ConfessionRepository) settingsColl() *mongo.Collection {
	return MongoDB.Collection("confession_settings")
}

func (r *ConfessionRepository) queue() *mongo.Collection {
	return MongoDB.Collection("confession_queue")
}

// SaveConfession upserts a confession.
func (r *ConfessionRepository) SaveConfession(confession *entity.Confession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Assign ID if new
	if confession.ID == 0 {
		id, err := MongoDB.GetNextSequence(ctx, "confession_id")
		if err != nil {
			return fmt.Errorf("failed to get next confession ID: %w", err)
		}
		confession.ID = id
	}

	filter := bson.M{"id": confession.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.confessions().ReplaceOne(ctx, filter, confession, opts)
	if err != nil {
		return fmt.Errorf("failed to save confession: %w", err)
	}
	return nil
}

// GetConfession retrieves a confession by ID.
func (r *ConfessionRepository) GetConfession(id int) (*entity.Confession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var confession entity.Confession
	err := r.confessions().FindOne(ctx, bson.M{"id": id}).Decode(&confession)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("confession not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get confession: %w", err)
	}
	return &confession, nil
}

// GetConfessionsByGuild retrieves all confessions for a guild.
func (r *ConfessionRepository) GetConfessionsByGuild(guildID string) ([]*entity.Confession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.confessions().Find(ctx, bson.M{"guild_id": guildID})
	if err != nil {
		return nil, fmt.Errorf("failed to query confessions: %w", err)
	}
	defer cursor.Close(ctx)

	var confessions []*entity.Confession
	if err := cursor.All(ctx, &confessions); err != nil {
		return nil, err
	}
	return confessions, nil
}

// GetPendingConfessions retrieves pending confessions for a guild.
func (r *ConfessionRepository) GetPendingConfessions(guildID string) ([]*entity.Confession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"guild_id": guildID,
		"status":   string(entity.ConfessionStatusPending),
	}
	cursor, err := r.confessions().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending confessions: %w", err)
	}
	defer cursor.Close(ctx)

	var confessions []*entity.Confession
	if err := cursor.All(ctx, &confessions); err != nil {
		return nil, err
	}
	return confessions, nil
}

// SaveReply saves a confession reply.
func (r *ConfessionRepository) SaveReply(reply *entity.ConfessionReply) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if reply.ID == 0 {
		id, err := MongoDB.GetNextSequence(ctx, "confession_reply_id")
		if err != nil {
			return fmt.Errorf("failed to get next reply ID: %w", err)
		}
		reply.ID = id
	}

	_, err := r.replies().InsertOne(ctx, reply)
	if err != nil {
		return fmt.Errorf("failed to save reply: %w", err)
	}
	return nil
}

// GetReplies retrieves all replies for a confession.
func (r *ConfessionRepository) GetReplies(confessionID int) ([]*entity.ConfessionReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.replies().Find(ctx, bson.M{"confession_id": confessionID})
	if err != nil {
		return nil, fmt.Errorf("failed to query replies: %w", err)
	}
	defer cursor.Close(ctx)

	var result []*entity.ConfessionReply
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SaveSettings upserts guild confession settings.
func (r *ConfessionRepository) SaveSettings(settings *entity.GuildConfessionSettings) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"guild_id": settings.GuildID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.settingsColl().ReplaceOne(ctx, filter, settings, opts)
	if err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}
	return nil
}

// GetSettings retrieves guild confession settings.
func (r *ConfessionRepository) GetSettings(guildID string) (*entity.GuildConfessionSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var settings entity.GuildConfessionSettings
	err := r.settingsColl().FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("settings not found for guild: %s", guildID)
		}
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}
	return &settings, nil
}

// AddToQueue adds a confession to the queue.
func (r *ConfessionRepository) AddToQueue(item *entity.ConfessionQueue) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.queue().InsertOne(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to add to queue: %w", err)
	}
	return nil
}

// GetQueue retrieves the confession queue.
func (r *ConfessionRepository) GetQueue() ([]*entity.ConfessionQueue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.queue().Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query queue: %w", err)
	}
	defer cursor.Close(ctx)

	var result []*entity.ConfessionQueue
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// RemoveFromQueue removes an item from the queue by confession ID.
func (r *ConfessionRepository) RemoveFromQueue(confession *entity.Confession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.queue().DeleteOne(ctx, bson.M{"confession.id": confession.ID})
	if err != nil {
		return fmt.Errorf("failed to remove from queue: %w", err)
	}
	return nil
}

// DeleteConfession deletes a confession and its replies.
func (r *ConfessionRepository) DeleteConfession(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := r.confessions().DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return fmt.Errorf("failed to delete confession: %w", err)
	}
	if _, err := r.replies().DeleteMany(ctx, bson.M{"confession_id": id}); err != nil {
		return fmt.Errorf("failed to delete replies: %w", err)
	}
	return nil
}
