package chatbot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/pkg/ai"
	"github.com/nerufuyo/nerubot/internal/pkg/backend"
	redispkg "github.com/nerufuyo/nerubot/internal/pkg/redis"
)

// ChatSession represents a user's chat session
type ChatSession struct {
	UserID    string       `json:"user_id"`
	Messages  []ai.Message `json:"messages"`
	CreatedAt time.Time    `json:"created_at"`
	LastUsed  time.Time    `json:"last_used"`
}

// RateLimitEntry tracks per-user rate limiting
type RateLimitEntry struct {
	Timestamps []time.Time
}

// ChatbotService handles AI chatbot functionality
type ChatbotService struct {
	providers      []ai.AIProvider
	sessions       map[string]*ChatSession
	sessionMutex   sync.RWMutex
	timeout        time.Duration
	systemPrompt   string
	redis          *redispkg.Client
	backendClient  *backend.Client

	// Rate limiting
	rateLimits     map[string]*RateLimitEntry
	rateLimitMutex sync.Mutex
}

// NewChatbotService creates a new chatbot service
func NewChatbotService(deepseekKey string, redis *redispkg.Client, backendClient *backend.Client) *ChatbotService {
	providers := make([]ai.AIProvider, 0)

	// Add DeepSeek provider
	if deepseekKey != "" {
		providers = append(providers, ai.NewDeepSeekProvider(deepseekKey))
	}

	service := &ChatbotService{
		providers:     providers,
		sessions:      make(map[string]*ChatSession),
		timeout:       30 * time.Minute,
		systemPrompt:  getNeruPersonality(),
		redis:         redis,
		backendClient: backendClient,
		rateLimits:    make(map[string]*RateLimitEntry),
	}

	// Start session cleanup goroutine
	go service.cleanupSessions()

	return service
}

// getNeruPersonality returns Neru's base personality prompt
func getNeruPersonality() string {
	return `You are Neru, a friendly and helpful AI companion integrated into a Discord bot called NeruBot, powered by Nerufuyo's personal knowledge base and real data.

CORE TRAITS:
- Friendly and approachable, like talking to a good friend
- Enthusiastic about technology, software engineering, and helping people
- Smart but not arrogant — you explain things clearly without being condescending
- Occasionally playful and witty, but never mean-spirited
- Genuinely interested in what users are saying
- You know everything about Nerufuyo (Listyo Adi) — his projects, experience, skills, and work

COMMUNICATION STYLE:
- Keep responses conversational and natural
- Use emojis sparingly (1-2 per message maximum)
- Be concise — aim for 2-3 sentences unless more detail is specifically requested
- If you don't know something, admit it honestly
- Avoid overly formal language — be casual but respectful

KNOWLEDGE:
- You have access to Nerufuyo's real data: projects, work experience, articles, and knowledge base
- Reference specific projects, technologies, and experiences when relevant
- Link to nerufuyo-workspace.com for more details
- You can answer questions about Neru's skills, projects, experience, and services
- For general topics, you can help too — but you shine when talking about Neru's work

BOUNDARIES:
- Don't pretend to be human
- Don't make promises about features you can't deliver
- Don't engage with inappropriate or harmful content
- Direct technical issues to the bot developer (@nerufuyo)

SPECIAL NOTES:
- You're part of NeruBot, which has music streaming, confessions, roasts, news, and crypto alerts
- You remember context within a conversation session
- Users can reset their chat history with /chat-reset
- You're powered by RAG (Retrieval Augmented Generation) with live data from Nerufuyo's database

Be yourself, be helpful, and most importantly — be knowledgeable about Neru's work!`
}

// buildSystemPrompt builds the full system prompt with RAG knowledge context
func (s *ChatbotService) buildSystemPrompt() string {
	base := s.systemPrompt

	// Override with backend settings prompt if available
	if s.backendClient != nil {
		settings := s.backendClient.GetSettings()
		if settings.SystemPrompt != "" {
			base = settings.SystemPrompt
		}
	}

	// Inject RAG knowledge context from backend (if enabled)
	if s.backendClient != nil {
		settings := s.backendClient.GetSettings()
		enableRAG := true // default
		if settings.ChatSettings.MaxHistoryMessages > 0 {
			enableRAG = settings.ChatSettings.EnableRAG
		}
		if enableRAG {
			ragContext := s.backendClient.GetRAGContext()
			if ragContext != "" {
				base = base + "\n\nKNOWLEDGE BASE (Real data from Nerufuyo's database — use this to answer questions accurately):\n" + ragContext
			}
		}
	}

	// Add current timestamp
	base = fmt.Sprintf("CURRENT DATE: %s\n\n%s", time.Now().UTC().Format("January 2, 2006"), base)

	return base
}

// CheckRateLimit checks if a user has exceeded the rate limit.
// Returns (allowed, remaining, resetSeconds).
func (s *ChatbotService) CheckRateLimit(userID string) (bool, int, int) {
	maxMessages := 5
	windowSeconds := 180 // 3 minutes

	// Override from backend settings
	if s.backendClient != nil {
		settings := s.backendClient.GetSettings()
		if settings.RateLimitCount > 0 {
			maxMessages = settings.RateLimitCount
		}
		if settings.RateLimitWindow > 0 {
			windowSeconds = settings.RateLimitWindow
		}
	}

	window := time.Duration(windowSeconds) * time.Second
	now := time.Now()

	s.rateLimitMutex.Lock()
	defer s.rateLimitMutex.Unlock()

	entry, exists := s.rateLimits[userID]
	if !exists {
		entry = &RateLimitEntry{}
		s.rateLimits[userID] = entry
	}

	// Remove expired timestamps
	var valid []time.Time
	for _, ts := range entry.Timestamps {
		if now.Sub(ts) < window {
			valid = append(valid, ts)
		}
	}
	entry.Timestamps = valid

	remaining := maxMessages - len(entry.Timestamps)
	if remaining < 0 {
		remaining = 0
	}

	if len(entry.Timestamps) >= maxMessages {
		// Calculate reset time from the oldest entry
		oldestInWindow := entry.Timestamps[0]
		resetAt := oldestInWindow.Add(window)
		resetSeconds := int(resetAt.Sub(now).Seconds()) + 1
		if resetSeconds < 1 {
			resetSeconds = 1
		}
		return false, 0, resetSeconds
	}

	// Record this usage
	entry.Timestamps = append(entry.Timestamps, now)
	remaining = maxMessages - len(entry.Timestamps)

	return true, remaining, 0
}

// Chat sends a message and returns the AI response
func (s *ChatbotService) Chat(ctx context.Context, userID, message, lang string) (string, error) {
	if len(s.providers) == 0 {
		return "", fmt.Errorf("no AI providers configured")
	}

	// Get or create session
	session := s.getOrCreateSession(userID)

	// Build messages with RAG-enhanced system prompt
	messages := make([]ai.Message, 0, len(session.Messages)+2)
	
	// Always include the system prompt with latest RAG context + language instruction
	systemPrompt := s.buildSystemPrompt()
	if lang != "" && lang != config.DefaultLang {
		systemPrompt += "\n\n" + config.LanguagePromptInstruction(lang)
	}
	messages = append(messages, ai.Message{
		Role:    "system",
		Content: systemPrompt,
	})
	
	// Add conversation history (limit from dashboard settings)
	history := session.Messages
	maxHistory := 10 // default
	if s.backendClient != nil {
		cs := s.backendClient.GetSettings().ChatSettings
		if cs.MaxHistoryMessages > 0 {
			maxHistory = cs.MaxHistoryMessages
		}
	}
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}
	messages = append(messages, history...)
	
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

		// Persist session to Redis
		s.saveSessionToRedis(userID, session)

		// Truncate if too long for Discord embed description (4096 chars)
		if len(response) > 4000 {
			response = response[:4000] + "..."
		}

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

	// Remove from Redis
	if s.redis != nil {
		_ = s.redis.Delete(context.Background(), "chat:session:"+userID)
	}
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
	if exists {
		return session
	}

	// Try loading from Redis
	if s.redis != nil {
		var cached ChatSession
		found, err := s.redis.Get(context.Background(), "chat:session:"+userID, &cached)
		if err == nil && found {
			s.sessions[userID] = &cached
			return &cached
		}
	}

	session = &ChatSession{
		UserID:    userID,
		Messages:  make([]ai.Message, 0),
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}
	s.sessions[userID] = session
	return session
}

// saveSessionToRedis persists a session to Redis with TTL.
func (s *ChatbotService) saveSessionToRedis(userID string, session *ChatSession) {
	if s.redis == nil {
		return
	}
	_ = s.redis.Set(context.Background(), "chat:session:"+userID, session, s.timeout)
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
