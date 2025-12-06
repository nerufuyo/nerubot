package chatbot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/pkg/ai"
)

// ChatSession represents a user's chat session
type ChatSession struct {
	UserID    string
	Messages  []ai.Message
	CreatedAt time.Time
	LastUsed  time.Time
}

// ChatbotService handles AI chatbot functionality
type ChatbotService struct {
	providers     []ai.AIProvider
	sessions      map[string]*ChatSession
	sessionMutex  sync.RWMutex
	timeout       time.Duration
	systemPrompt  string
}

// NewChatbotService creates a new chatbot service
func NewChatbotService(deepseekKey string) *ChatbotService {
	providers := make([]ai.AIProvider, 0)

	// Add DeepSeek provider
	if deepseekKey != "" {
		providers = append(providers, ai.NewDeepSeekProvider(deepseekKey))
	}

	service := &ChatbotService{
		providers:    providers,
		sessions:     make(map[string]*ChatSession),
		timeout:      30 * time.Minute,
		systemPrompt: getNeruPersonality(),
	}

	// Start session cleanup goroutine
	go service.cleanupSessions()

	return service
}

// getNeruPersonality returns Neru's personality prompt
func getNeruPersonality() string {
	return `You are Neru, a friendly and helpful AI companion integrated into a Discord bot called NeruBot. Here's your personality:

CORE TRAITS:
- Friendly and approachable, like talking to a good friend
- Enthusiastic about music, technology, and helping people
- Smart but not arrogant - you explain things clearly without being condescending
- Occasionally playful and witty, but never mean-spirited
- Genuinely interested in what users are saying

COMMUNICATION STYLE:
- Keep responses conversational and natural
- Use emojis sparingly (1-2 per message maximum)
- Be concise - aim for 2-3 sentences unless more detail is specifically requested
- If you don't know something, admit it honestly
- Avoid overly formal language - be casual but respectful

KNOWLEDGE AREAS:
- Discord bot features and commands
- Music recommendations and trivia
- General technology and programming
- Gaming culture
- Crypto and blockchain basics
- Current events and news

BOUNDARIES:
- Don't pretend to be human
- Don't make promises about features you can't deliver
- Don't engage with inappropriate or harmful content
- Direct technical issues to the bot developer (@nerufuyo)

SPECIAL NOTES:
- You're part of NeruBot, which has music streaming, confessions, roasts, news, and crypto alerts
- You remember context within a conversation session
- Users can reset their chat history with /chat-reset

Be yourself, be helpful, and most importantly - be a good friend to the community!`
}

// Chat sends a message and returns the AI response
func (s *ChatbotService) Chat(ctx context.Context, userID, message string) (string, error) {
	if len(s.providers) == 0 {
		return "", fmt.Errorf("no AI providers configured")
	}

	// Get or create session
	session := s.getOrCreateSession(userID)

	// Build messages with system prompt
	messages := make([]ai.Message, 0, len(session.Messages)+2)
	
	// Add system prompt if this is the first message
	if len(session.Messages) == 0 {
		messages = append(messages, ai.Message{
			Role:    "system",
			Content: s.systemPrompt,
		})
	}
	
	// Add conversation history
	messages = append(messages, session.Messages...)
	
	// Add new user message
	messages = append(messages, ai.Message{
		Role:    "user",
		Content: message,
	})

	// Try each provider in order
	var lastErr error
	for _, provider := range s.providers {
		if !provider.IsAvailable() {
			continue
		}

		response, err := provider.Chat(ctx, messages)
		if err != nil {
			lastErr = err
			continue
		}

		// Add user message and assistant response to session
		session.Messages = append(session.Messages, ai.Message{
			Role:    "user",
			Content: message,
		})
		session.Messages = append(session.Messages, ai.Message{
			Role:    "assistant",
			Content: response,
		})
		session.LastUsed = time.Now()

		return response, nil
	}

	if lastErr != nil {
		return "", fmt.Errorf("all AI providers failed: %w", lastErr)
	}

	return "", fmt.Errorf("no available AI providers")
}

// ResetSession clears a user's chat history
func (s *ChatbotService) ResetSession(userID string) {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()

	delete(s.sessions, userID)
}

// GetSessionInfo returns information about a user's session
func (s *ChatbotService) GetSessionInfo(userID string) (messageCount int, createdAt, lastUsed time.Time, exists bool) {
	s.sessionMutex.RLock()
	defer s.sessionMutex.RUnlock()

	session, exists := s.sessions[userID]
	if !exists {
		return 0, time.Time{}, time.Time{}, false
	}

	return len(session.Messages), session.CreatedAt, session.LastUsed, true
}

// getOrCreateSession retrieves or creates a chat session
func (s *ChatbotService) getOrCreateSession(userID string) *ChatSession {
	s.sessionMutex.Lock()
	defer s.sessionMutex.Unlock()

	session, exists := s.sessions[userID]
	if !exists {
		session = &ChatSession{
			UserID:    userID,
			Messages:  make([]ai.Message, 0),
			CreatedAt: time.Now(),
			LastUsed:  time.Now(),
		}
		s.sessions[userID] = session
	}

	return session
}

// cleanupSessions removes expired sessions
func (s *ChatbotService) cleanupSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.sessionMutex.Lock()
		now := time.Now()
		for userID, session := range s.sessions {
			if now.Sub(session.LastUsed) > s.timeout {
				delete(s.sessions, userID)
			}
		}
		s.sessionMutex.Unlock()
	}
}

// GetAvailableProviders returns the list of configured providers
func (s *ChatbotService) GetAvailableProviders() []string {
	providers := make([]string, 0)
	for _, p := range s.providers {
		if p.IsAvailable() {
			providers = append(providers, p.Name())
		}
	}
	return providers
}
