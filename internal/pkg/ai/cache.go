package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	redispkg "github.com/nerufuyo/nerubot/internal/pkg/redis"
)

// CachedProvider wraps an AIProvider with Redis-based response caching.
// It caches responses keyed by a hash of the non-system messages, so
// identical questions get instant replies without consuming tokens.
type CachedProvider struct {
	inner  AIProvider
	redis  *redispkg.Client
	ttl    time.Duration
	logger *logger.Logger
}

// CachedResponse is the Redis-stored value.
type CachedResponse struct {
	Content  string `json:"content"`
	CachedAt string `json:"cached_at"`
}

// NewCachedProvider wraps an existing provider with Redis caching.
// ttl controls how long a cached response is valid (e.g. 1 hour).
func NewCachedProvider(inner AIProvider, redis *redispkg.Client, ttl time.Duration) *CachedProvider {
	return &CachedProvider{
		inner:  inner,
		redis:  redis,
		ttl:    ttl,
		logger: logger.New("ai-cache"),
	}
}

// Name returns the underlying provider name with a cache tag.
func (c *CachedProvider) Name() string {
	return c.inner.Name()
}

// IsAvailable delegates to the underlying provider.
func (c *CachedProvider) IsAvailable() bool {
	return c.inner.IsAvailable()
}

// Chat checks the cache first; on miss, calls the real provider and stores the result.
// Only user messages are hashed (system prompts change with date/settings and shouldn't
// invalidate the cache for the same user question).
func (c *CachedProvider) Chat(ctx context.Context, messages []Message) (string, error) {
	key := c.cacheKey(messages)

	// Try cache
	if c.redis != nil {
		var cached CachedResponse
		found, err := c.redis.Get(ctx, key, &cached)
		if err == nil && found && cached.Content != "" {
			c.logger.Info("Cache hit", "key_suffix", key[len(key)-8:])
			return cached.Content, nil
		}
	}

	// Cache miss — call real provider
	response, err := c.inner.Chat(ctx, messages)
	if err != nil {
		return "", err
	}

	// Store in cache
	if c.redis != nil {
		entry := CachedResponse{
			Content:  response,
			CachedAt: time.Now().UTC().Format(time.RFC3339),
		}
		if setErr := c.redis.Set(ctx, key, entry, c.ttl); setErr != nil {
			c.logger.Warn("Failed to cache response", "error", setErr)
		}
	}

	return response, nil
}

// cacheKey builds a deterministic Redis key from the user/assistant messages.
// System messages are excluded so that minor prompt changes (date, settings)
// don't bust the cache for the same conversation.
func (c *CachedProvider) cacheKey(messages []Message) string {
	h := sha256.New()
	for _, m := range messages {
		if m.Role == "system" {
			continue
		}
		fmt.Fprintf(h, "%s:%s\n", m.Role, strings.TrimSpace(m.Content))
	}
	hash := hex.EncodeToString(h.Sum(nil))[:16]
	return "ai:cache:" + c.inner.Name() + ":" + hash
}
