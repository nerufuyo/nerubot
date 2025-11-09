package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// JSONRepository provides base functionality for JSON-based repositories
type JSONRepository struct {
	filePath string
	mu       sync.RWMutex
	logger   *logger.Logger
}

// NewJSONRepository creates a new JSON repository
func NewJSONRepository(filePath string) *JSONRepository {
	return &JSONRepository{
		filePath: filePath,
		logger:   logger.New("repository"),
	}
}

// Load loads data from JSON file into the provided interface
func (r *JSONRepository) Load(v interface{}) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		// File doesn't exist, initialize with empty data
		r.logger.Debug("File does not exist, initializing", "path", r.filePath)
		return nil
	}

	// Read file
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	if len(data) > 0 {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	return nil
}

// Save saves data to JSON file
func (r *JSONRepository) Save(v interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to temporary file first
	tempFile := r.filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Rename temporary file to actual file (atomic operation)
	if err := os.Rename(tempFile, r.filePath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

// Exists checks if the file exists
func (r *JSONRepository) Exists() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, err := os.Stat(r.filePath)
	return err == nil
}

// Delete deletes the file
func (r *JSONRepository) Delete() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.Remove(r.filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
