package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// Client wraps the Redis client with convenience methods.
type Client struct {
	rdb    *redis.Client
	logger *logger.Logger
}

// New connects to Redis and returns a ready Client.
func New(redisURL string) (*Client, error) {
	log := logger.New("redis")

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis parse url: %w", err)
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	log.Info("Connected to Redis")

	return &Client{
		rdb:    rdb,
		logger: log,
	}, nil
}

// Close gracefully closes the Redis connection.
func (c *Client) Close() error {
	c.logger.Info("Disconnecting from Redis")
	return c.rdb.Close()
}

// Underlying returns the raw *redis.Client for advanced usage.
func (c *Client) Underlying() *redis.Client {
	return c.rdb
}

// Set stores a value with optional TTL. Pass 0 for no expiration.
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis marshal: %w", err)
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value and unmarshals it into dest. Returns false if key does not exist.
func (c *Client) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("redis get: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("redis unmarshal: %w", err)
	}
	return true, nil
}

// Delete removes one or more keys.
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Exists checks if a key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// SetTTL updates the TTL of a key.
func (c *Client) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return c.rdb.Expire(ctx, key, ttl).Err()
}
