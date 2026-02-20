package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// Client wraps the MongoDB client with convenience methods.
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	logger   *logger.Logger
}

// New connects to MongoDB and returns a ready Client.
func New(mongoURL, dbName string) (*Client, error) {
	log := logger.New("mongodb")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongodb connect: %w", err)
	}

	// Verify connectivity
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongodb ping: %w", err)
	}

	log.Info("Connected to MongoDB", "database", dbName)

	return &Client{
		client:   client,
		database: client.Database(dbName),
		logger:   log,
	}, nil
}

// Database returns the underlying *mongo.Database.
func (c *Client) Database() *mongo.Database {
	return c.database
}

// Collection returns a collection handle.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// Disconnect gracefully closes the MongoDB connection.
func (c *Client) Disconnect(ctx context.Context) error {
	c.logger.Info("Disconnecting from MongoDB")
	return c.client.Disconnect(ctx)
}

// EnsureIndexes creates indexes for all collections used by the bot.
func (c *Client) EnsureIndexes(ctx context.Context) error {
	indexes := map[string][]mongo.IndexModel{
		"confessions": {
			{Keys: bson.D{{Key: "guild_id", Value: 1}}},
			{Keys: bson.D{{Key: "status", Value: 1}}},
		},
		"confession_replies": {
			{Keys: bson.D{{Key: "confession_id", Value: 1}}},
		},
		"confession_settings": {
			{Keys: bson.D{{Key: "guild_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"roast_profiles": {
			{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "guild_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"roast_patterns": {
			{Keys: bson.D{{Key: "category", Value: 1}}},
		},
		"roast_stats": {
			{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "guild_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"roast_history": {
			{Keys: bson.D{{Key: "target_id", Value: 1}, {Key: "guild_id", Value: 1}}},
		},
		"server_stats": {
			{Keys: bson.D{{Key: "guild_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"user_stats": {
			{Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"bot_config": {
			{Keys: bson.D{{Key: "key", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"guild_configs": {
			{Keys: bson.D{{Key: "guild_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
	}

	for coll, models := range indexes {
		if _, err := c.Collection(coll).Indexes().CreateMany(ctx, models); err != nil {
			return fmt.Errorf("create indexes for %s: %w", coll, err)
		}
	}

	c.logger.Info("MongoDB indexes ensured")
	return nil
}

// GetNextSequence atomically increments and returns the next integer ID for the given sequence name.
// This replaces the JSON auto-increment pattern.
func (c *Client) GetNextSequence(ctx context.Context, name string) (int, error) {
	coll := c.Collection("counters")
	filter := bson.M{"_id": name}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("get next sequence %s: %w", name, err)
	}
	return result.Seq, nil
}
