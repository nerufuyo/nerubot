package analytics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
)

// AnalyticsService handles bot analytics and statistics
type AnalyticsService struct {
	serverStats  map[string]*entity.ServerStats
	userStats    map[string]*entity.UserStats
	globalStats  *entity.GlobalStats
	mu           sync.RWMutex
	dataDir      string
	autoSave     bool
	saveInterval time.Duration
	stopChan     chan struct{}
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(dataDir string) *AnalyticsService {
	service := &AnalyticsService{
		serverStats:  make(map[string]*entity.ServerStats),
		userStats:    make(map[string]*entity.UserStats),
		globalStats:  entity.NewGlobalStats(),
		dataDir:      dataDir,
		autoSave:     true,
		saveInterval: 5 * time.Minute,
		stopChan:     make(chan struct{}),
	}

	// Load existing stats
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

	// Update counts
	s.globalStats.TotalGuilds = len(s.serverStats)
	s.globalStats.TotalUsers = len(s.userStats)
	s.globalStats.Uptime = time.Since(s.globalStats.StartTime)

	return s.globalStats
}

// GetTopServers returns the most active servers
func (s *AnalyticsService) GetTopServers(limit int) []*entity.ServerStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice
	servers := make([]*entity.ServerStats, 0, len(s.serverStats))
	for _, stats := range s.serverStats {
		servers = append(servers, stats)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].CommandsUsed > servers[j].CommandsUsed
	})

	// Limit results
	if len(servers) > limit {
		servers = servers[:limit]
	}

	return servers
}

// GetTopUsers returns the most active users
func (s *AnalyticsService) GetTopUsers(limit int) []*entity.UserStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice
	users := make([]*entity.UserStats, 0, len(s.userStats))
	for _, stats := range s.userStats {
		users = append(users, stats)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CommandsUsed > users[j].CommandsUsed
	})

	// Limit results
	if len(users) > limit {
		users = users[:limit]
	}

	return users
}

// Save persists analytics data to disk
func (s *AnalyticsService) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Ensure data directory exists
	if err := os.MkdirAll(s.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Save server stats
	serverFile := filepath.Join(s.dataDir, "server_stats.json")
	if err := s.saveJSON(serverFile, s.serverStats); err != nil {
		return fmt.Errorf("failed to save server stats: %w", err)
	}

	// Save user stats
	userFile := filepath.Join(s.dataDir, "user_stats.json")
	if err := s.saveJSON(userFile, s.userStats); err != nil {
		return fmt.Errorf("failed to save user stats: %w", err)
	}

	// Save global stats
	globalFile := filepath.Join(s.dataDir, "global_stats.json")
	if err := s.saveJSON(globalFile, s.globalStats); err != nil {
		return fmt.Errorf("failed to save global stats: %w", err)
	}

	return nil
}

// Load reads analytics data from disk
func (s *AnalyticsService) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Load server stats
	serverFile := filepath.Join(s.dataDir, "server_stats.json")
	if err := s.loadJSON(serverFile, &s.serverStats); err != nil {
		// It's okay if file doesn't exist yet
		s.serverStats = make(map[string]*entity.ServerStats)
	}

	// Load user stats
	userFile := filepath.Join(s.dataDir, "user_stats.json")
	if err := s.loadJSON(userFile, &s.userStats); err != nil {
		s.userStats = make(map[string]*entity.UserStats)
	}

	// Load global stats
	globalFile := filepath.Join(s.dataDir, "global_stats.json")
	if err := s.loadJSON(globalFile, &s.globalStats); err != nil {
		s.globalStats = entity.NewGlobalStats()
	}

	return nil
}

// Stop stops the auto-save loop and saves data
func (s *AnalyticsService) Stop() error {
	if s.autoSave {
		close(s.stopChan)
	}
	return s.Save()
}

// autoSaveLoop periodically saves data to disk
func (s *AnalyticsService) autoSaveLoop() {
	ticker := time.NewTicker(s.saveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.Save(); err != nil {
				// Log error but don't crash
				fmt.Printf("Failed to auto-save analytics: %v\n", err)
			}
		case <-s.stopChan:
			return
		}
	}
}

// saveJSON saves data to a JSON file
func (s *AnalyticsService) saveJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// loadJSON loads data from a JSON file
func (s *AnalyticsService) loadJSON(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(data)
}

// ResetServerStats resets statistics for a specific server
func (s *AnalyticsService) ResetServerStats(guildID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.serverStats, guildID)
	return nil
}

// ResetUserStats resets statistics for a specific user
func (s *AnalyticsService) ResetUserStats(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.userStats, userID)
	return nil
}
