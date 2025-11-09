package repository

import (
	"fmt"
	"sync"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/entity"
)

// ConfessionRepository handles confession data persistence
type ConfessionRepository struct {
	confessions *JSONRepository
	replies     *JSONRepository
	settings    *JSONRepository
	queue       *JSONRepository
	
	confessionData map[int]*entity.Confession
	replyData      map[int][]*entity.ConfessionReply
	settingsData   map[string]*entity.GuildConfessionSettings
	queueData      []*entity.ConfessionQueue
	
	mu             sync.RWMutex
	nextID         int
}

// NewConfessionRepository creates a new confession repository
func NewConfessionRepository() *ConfessionRepository {
	repo := &ConfessionRepository{
		confessions:    NewJSONRepository(config.ConfessionsFile),
		replies:        NewJSONRepository(config.RepliesFile),
		settings:       NewJSONRepository(config.ConfessionSettingsFile),
		queue:          NewJSONRepository(config.ConfessionQueueFile),
		confessionData: make(map[int]*entity.Confession),
		replyData:      make(map[int][]*entity.ConfessionReply),
		settingsData:   make(map[string]*entity.GuildConfessionSettings),
		queueData:      make([]*entity.ConfessionQueue, 0),
		nextID:         1,
	}
	
	// Load existing data
	_ = repo.Load()
	
	return repo
}

// Load loads all confession data from files
func (r *ConfessionRepository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Load confessions
	if err := r.confessions.Load(&r.confessionData); err != nil {
		return fmt.Errorf("failed to load confessions: %w", err)
	}
	
	// Load replies
	if err := r.replies.Load(&r.replyData); err != nil {
		return fmt.Errorf("failed to load replies: %w", err)
	}
	
	// Load settings
	if err := r.settings.Load(&r.settingsData); err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}
	
	// Load queue
	if err := r.queue.Load(&r.queueData); err != nil {
		return fmt.Errorf("failed to load queue: %w", err)
	}
	
	// Update next ID
	for id := range r.confessionData {
		if id >= r.nextID {
			r.nextID = id + 1
		}
	}
	
	return nil
}

// SaveConfession saves a confession
func (r *ConfessionRepository) SaveConfession(confession *entity.Confession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Assign ID if new
	if confession.ID == 0 {
		confession.ID = r.nextID
		r.nextID++
	}
	
	r.confessionData[confession.ID] = confession
	return r.confessions.Save(r.confessionData)
}

// GetConfession retrieves a confession by ID
func (r *ConfessionRepository) GetConfession(id int) (*entity.Confession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	confession, exists := r.confessionData[id]
	if !exists {
		return nil, fmt.Errorf("confession not found: %d", id)
	}
	
	return confession, nil
}

// GetConfessionsByGuild retrieves all confessions for a guild
func (r *ConfessionRepository) GetConfessionsByGuild(guildID string) ([]*entity.Confession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var confessions []*entity.Confession
	for _, confession := range r.confessionData {
		if confession.GuildID == guildID {
			confessions = append(confessions, confession)
		}
	}
	
	return confessions, nil
}

// GetPendingConfessions retrieves pending confessions for a guild
func (r *ConfessionRepository) GetPendingConfessions(guildID string) ([]*entity.Confession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var confessions []*entity.Confession
	for _, confession := range r.confessionData {
		if confession.GuildID == guildID && confession.IsPending() {
			confessions = append(confessions, confession)
		}
	}
	
	return confessions, nil
}

// SaveReply saves a confession reply
func (r *ConfessionRepository) SaveReply(reply *entity.ConfessionReply) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Assign ID if new
	if reply.ID == 0 {
		// Find max ID for this confession
		maxID := 0
		for _, r := range r.replyData[reply.ConfessionID] {
			if r.ID > maxID {
				maxID = r.ID
			}
		}
		reply.ID = maxID + 1
	}
	
	r.replyData[reply.ConfessionID] = append(r.replyData[reply.ConfessionID], reply)
	return r.replies.Save(r.replyData)
}

// GetReplies retrieves all replies for a confession
func (r *ConfessionRepository) GetReplies(confessionID int) ([]*entity.ConfessionReply, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	replies, exists := r.replyData[confessionID]
	if !exists {
		return []*entity.ConfessionReply{}, nil
	}
	
	return replies, nil
}

// SaveSettings saves guild confession settings
func (r *ConfessionRepository) SaveSettings(settings *entity.GuildConfessionSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.settingsData[settings.GuildID] = settings
	return r.settings.Save(r.settingsData)
}

// GetSettings retrieves guild confession settings
func (r *ConfessionRepository) GetSettings(guildID string) (*entity.GuildConfessionSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	settings, exists := r.settingsData[guildID]
	if !exists {
		return nil, fmt.Errorf("settings not found for guild: %s", guildID)
	}
	
	return settings, nil
}

// AddToQueue adds a confession to the queue
func (r *ConfessionRepository) AddToQueue(item *entity.ConfessionQueue) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.queueData = append(r.queueData, item)
	return r.queue.Save(r.queueData)
}

// GetQueue retrieves the confession queue
func (r *ConfessionRepository) GetQueue() ([]*entity.ConfessionQueue, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.queueData, nil
}

// RemoveFromQueue removes an item from the queue
func (r *ConfessionRepository) RemoveFromQueue(confession *entity.Confession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	for i, item := range r.queueData {
		if item.Confession.ID == confession.ID {
			r.queueData = append(r.queueData[:i], r.queueData[i+1:]...)
			return r.queue.Save(r.queueData)
		}
	}
	
	return nil
}

// DeleteConfession deletes a confession and its replies
func (r *ConfessionRepository) DeleteConfession(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	delete(r.confessionData, id)
	delete(r.replyData, id)
	
	if err := r.confessions.Save(r.confessionData); err != nil {
		return err
	}
	
	return r.replies.Save(r.replyData)
}
