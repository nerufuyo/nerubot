package ai

import "context"

// Message represents a chat message
type Message struct {
	Role    string // "user" or "assistant"
	Content string
}

// AIProvider defines the interface for AI service providers
type AIProvider interface {
	// Name returns the provider name
	Name() string

	// Chat sends a message and returns the response
	Chat(ctx context.Context, messages []Message) (string, error)

	// IsAvailable checks if the provider is configured and available
	IsAvailable() bool
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages []Message
	MaxTokens int
	Temperature float64
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Content  string
	Provider string
	Error    error
}
