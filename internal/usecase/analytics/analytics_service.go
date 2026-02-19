package analytics

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/mongodb"
)

// AnalyticsService handles bot analytics and statistics.
// It keeps in-memory caches and periodically flushes to MongoDB.
type AnalyticsService struct {
	serverStats  map[string]*entity.ServerStats
	userStats    map[string]*entity.UserStats
	globalStats  *entity.GlobalStats
	mu           sync.RWMutex
	db           *mongodb.Client
	autoSave     bool
	saveInterval time.Duration
	stopChan     chan struct{}
}

// NewAnalyticsService creates a new analytics service backed by MongoDB.
func NewAnalyticsService(db *mongodb.Client) *AnalyticsService {
	service := &AnalyticsService{
		serverStats:  make(map[string]*entity.ServerStats),
		userStats:    make(map[string]*entity.UserStats),
		globalStats:  entity.NewGlobalStats(),
		db:           db,
		autoSave:     true,
		saveInterval: 5 * time.Minute,
		stopChan:     make(chan struct{}),
	}

	// Load existing stats from MongoDB
	service.Load()

	// Start auto-save goroutine
	if service.autoSave {
		go service.autoSaveLoop()
	}

	return service
}

// RecordCommandUsage records a command execution
func (s *AnalyticsService) RecordCommandUsage(guildID, guildName, userID, username, commandName string, success bool, executionTime int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update server stats
	serverStats, exists := s.serverStats[guildID]
	if !exists {
		serverStats = entity.NewServerStats(guildID, guildName)
		s.serverStats[guildID] = serverStats
	}
	serverStats.RecordCommand(commandName, userID)

	// Update user stats
	userStats, exists := s.userStats[userID]
	if !exists {
		userStats = entity.NewUserStats(userID, username)
		s.userStats[userID] = userStats
	}
	userStats.RecordCommand(commandName)

	// Update global stats
	s.globalStats.TotalCommands++
	s.globalStats.TopCommands[commandName]++
	s.globalStats.TopGuilds[guildID]++
	s.globalStats.UpdatedAt = time.Now()

	// Update specific feature counts
	switch commandName {
	case "play":
		serverStats.RecordSong()
		userStats.SongsRequested++
		s.globalStats.TotalSongsPlayed++
	case "confess":
		serverStats.RecordConfession()
		userStats.ConfessionsPosted++
		s.globalStats.TotalConfessions++
	case "roast":
		serverStats.RecordRoast()
		userStats.RoastsReceived++
		s.globalStats.TotalRoasts++
	case "chat":
		serverStats.RecordChat()
		userStats.ChatMessages++
		s.globalStats.TotalChatMessages++
	case "news":
		serverStats.NewsRequests++
		userStats.NewsRequests++
	case "whale":
		serverStats.WhaleAlerts++
		userStats.WhaleChecks++
	}
}

// GetServerStats returns statistics for a specific server
func (s *AnalyticsService) GetServerStats(guildID string) (*entity.ServerStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats, exists := s.serverStats[guildID]
	if !exists {
		return nil, fmt.Errorf("no stats found for guild %s", guildID)
	}

	return stats, nil
}

// GetUserStats returns statistics for a specific user
func (s *AnalyticsService) GetUserStats(userID string) (*entity.UserStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats, exists := s.userStats[userID]
	if !exists {
		return nil, fmt.Errorf("no stats found for user %s", userID)
	}

	return stats, nil
}

// GetGlobalStats returns global bot statistics
func (s *AnalyticsService) GetGlobalStats() *entity.GlobalStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.globalStats.TotalGuilds = len(s.serverStats)
	s.globalStats.TotalUsers = len(s.userStats)
	s.globalStats.Uptime = time.Since(s.globalStats.StartTime)

	return s.globalStats
}

// GetTopServers returns the most active servers
func (s *AnalyticsService) GetTopServers(limit int) []*entity.ServerStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	servers := make([]*entity.ServerStats, 0, len(s.serverStats))
	for _, stats := range s.serverStats {
		servers = append(servers, stats)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].CommandsUsed > servers[j].CommandsUsed
	})

	if len(servers) > limit {
		servers = servers[:limit]
	}

	return servers
}

// GetTopUsers returns the most active users
func (s *AnalyticsService) GetTopUsers(limit int) []*entity.UserStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*entity.UserStats, 0, len(s.userStats))
	for _, stats := range s.userStats {
		users = append(users, stats)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CommandsUsed > users[j].CommandsUsed
	})

	if len(users) > limit {
		users = users[:limit]
	}

	return users
}

// Save persists analytics data to MongoDB.
func (s *AnalyticsService) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Save server stats
	serverColl := s.db.Collection("server_stats")
	for _, stats := range s.serverStats {
		filter := bson.M{"guild_id": stats.GuildID}
		opts := options.Replace().SetUpsert(true)
		if _, err := serverColl.ReplaceOne(ctx, filter, stats, opts); err != nil {
			return fmt.Errorf("failed to save server stats: %w", err)
		}
	}

	// Save user stats
	userColl := s.db.Collection("user_stats")
	for _, stats := range s.userStats {
		filter := bson.M{"user_id": stats.UserID}
		opts := options.Replace().SetUpsert(true)
		if _, err := userColl.ReplaceOne(ctx, filter, stats, opts); err != nil {
			return fmt.Errorf("failed to save user stats: %w", err)
		}
	}

	// Save global stats (single document, keyed by "global")
	globalColl := s.db.Collection("global_stats")
	filter := bson.M{"_id": "global"}
	opts := options.Replace().SetUpsert(true)
	if _, err := globalColl.ReplaceOne(ctx, filter, s.globalStats, opts); err != nil {
		return fmt.Errorf("failed to save global stats: %w", err)
	}

	return nil
}

// Load reads analytics data from MongoDB into memory.
func (s *AnalyticsService) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Load server stats
	serverColl := s.db.Collection("server_stats")
	cursor, err := serverColl.Find(ctx, bson.M{})
	if err == nil {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var stats entity.ServerStats
			if err := cursor.Decode(&stats); err == nil {
				s.serverStats[stats.GuildID] = &stats
			}
		}
	}

	// Load user stats
	userColl := s.db.Collection("user_stats")
	cursor2, err := userColl.Find(ctx, bson.M{})
	if err == nil {
		defer cursor2.Close(ctx)
		for cursor2.Next(ctx) {
			var stats entity.UserStats
			if err := cursor2.Decode(&stats); err == nil {
				s.userStats[stats.UserID] = &stats
			}
		}
	}

	// Load global stats
	globalColl := s.db.Collection("global_stats")
	var global entity.GlobalStats
	err = globalColl.FindOne(ctx, bson.M{"_id": "global"}).Decode(&global)
	if err == nil {
		s.globalStats = &global
	} else if err == mongo.ErrNoDocuments {
		s.globalStats = entity.NewGlobalStats()
	}

	return nil
}

// Stop stops the auto-save loop and saves data.
func (s *AnalyticsService) Stop() error {
	if s.autoSave {
		close(s.stopChan)
	}
	return s.Save()
}

// autoSaveLoop periodically saves data to MongoDB.
func (s *AnalyticsService) autoSaveLoop() {
	ticker := time.NewTicker(s.saveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.Save(); err != nil {
				fmt.Printf("Failed to auto-save analytics: %v\n", err)
			}
		case <-s.stopChan:
			return
		}
	}
}

// ResetServerStats resets statistics for a specific server
func (s *AnalyticsService) ResetServerStats(guildID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.serverStats, guildID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = s.db.Collection("server_stats").DeleteOne(ctx, bson.M{"guild_id": guildID})

	return nil
}

// ResetUserStats resets statistics for a specific user
func (s *AnalyticsService) ResetUserStats(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.userStats, userID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = s.db.Collection("user_stats").DeleteOne(ctx, bson.M{"user_id": userID})

	return nil
}
