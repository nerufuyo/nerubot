# Architecture Documentation

## Overview

NeruBot is built using **Clean Architecture** principles, ensuring a maintainable, testable, and scalable codebase. The architecture separates concerns into distinct layers, each with clear responsibilities and dependencies flowing inward.

## Architecture Diagram

![Architecture Diagram](docs/architecture.png)

## Layer Structure

### 1. Config Layer (Outermost)
**Location:** `internal/config/`

The configuration layer handles all application settings and environment variables.

**Components:**
- `config.go` - Main configuration struct and loading logic
- `constants.go` - Application-wide constants
- `messages.go` - User-facing message templates

**Responsibilities:**
- Load environment variables
- Provide default values
- Validate configuration
- Expose settings to other layers

### 2. Delivery Layer (External Interface)
**Location:** `internal/delivery/discord/`

The delivery layer handles external communication with Discord.

**Components:**
- `bot.go` - Discord bot initialization and lifecycle
- `handlers.go` - Slash command handlers

**Responsibilities:**
- Initialize Discord session
- Register slash commands
- Handle user interactions
- Transform Discord events to internal calls
- Format responses for Discord

**Dependencies:** Use Case Layer

### 3. Use Case Layer (Business Logic)
**Location:** `internal/usecase/`

The use case layer contains the core business logic for each feature.

**Components:**
- `music/music_service.go` - Music playback and queue management
- `confession/confession_service.go` - Anonymous confession handling
- `roast/roast_service.go` - User roasting with activity tracking
- `chatbot/chatbot_service.go` - AI chatbot with multi-provider support
- `news/news_service.go` - News aggregation from RSS feeds
- `whale/whale_service.go` - Cryptocurrency transaction monitoring

**Responsibilities:**
- Implement business rules
- Orchestrate operations
- Manage sessions and state
- Call repositories for data
- Integrate with external services

**Dependencies:** Entity Layer, Repository Layer, Pkg Layer

### 4. Entity Layer (Domain Models)
**Location:** `internal/entity/`

The entity layer defines the core domain models.

**Components:**
- `music.go` - Song, Queue models
- `confession.go` - Confession, Reply, Settings models
- `roast.go` - UserProfile, RoastPattern, Stats models
- `news.go` - NewsArticle model
- `whale.go` - Transaction, Address models

**Responsibilities:**
- Define domain structures
- Enforce business rules at model level
- Provide validation methods

**Dependencies:** None (core domain)

### 5. Repository Layer (Data Access)
**Location:** `internal/repository/`

The repository layer handles data persistence.

**Components:**
- `confession_repository.go` - Confession data storage (JSON)
- `roast_repository.go` - Roast data storage (JSON)
- `repository.go` - Base repository interface

**Responsibilities:**
- CRUD operations
- JSON file management
- Thread-safe data access
- Data serialization/deserialization

**Dependencies:** Entity Layer

### 6. Pkg Layer (Shared Utilities)
**Location:** `internal/pkg/`

The pkg layer provides shared utilities and external integrations.

**Components:**
- `ai/` - AI provider interface and implementations (Claude, Gemini, OpenAI)
- `ffmpeg/` - FFmpeg wrapper for audio processing
- `ytdlp/` - yt-dlp wrapper for YouTube downloads
- `logger/` - Structured logging utilities

**Responsibilities:**
- Provide reusable utilities
- Abstract external dependencies
- Handle external API calls
- Process audio/video

**Dependencies:** Config Layer

## Data Flow

### Request Flow (Inward)
```
User Input (Discord)
    ↓
Delivery Layer (handlers.go)
    ↓
Use Case Layer (service.go)
    ↓
Repository Layer (repository.go)
    ↓
Entity Layer (models)
```

### Response Flow (Outward)
```
Entity Layer (data models)
    ↓
Repository Layer (loaded data)
    ↓
Use Case Layer (processed result)
    ↓
Delivery Layer (formatted response)
    ↓
Discord (user sees result)
```

## Dependency Rules

1. **Inner layers never depend on outer layers**
   - Use Case doesn't know about Delivery
   - Entity doesn't know about Repository

2. **Dependencies point inward**
   - Delivery → Use Case → Entity
   - Repository → Entity

3. **Interfaces for abstraction**
   - Use Cases define repository interfaces
   - Repositories implement those interfaces

## Key Design Patterns

### 1. Dependency Injection
Services are injected into handlers:
```go
bot := &Bot{
    musicService: music.NewMusicService(),
    confessionService: confession.NewConfessionService(repo),
}
```

### 2. Repository Pattern
Data access abstracted through interfaces:
```go
type ConfessionRepository interface {
    Save(confession *entity.Confession) error
    FindByID(id int64) (*entity.Confession, error)
}
```

### 3. Service Pattern
Business logic encapsulated in services:
```go
type MusicService struct {
    queue *Queue
    ffmpeg *ffmpeg.FFmpeg
}
```

### 4. Strategy Pattern (AI Providers)
Multiple AI providers with fallback:
```go
type AIProvider interface {
    Chat(ctx context.Context, messages []Message) (string, error)
    IsAvailable() bool
}
```

## Concurrency

### Goroutines
- **Session Cleanup** - Background goroutine cleans up expired sessions
- **News Fetching** - Concurrent fetching from multiple RSS sources
- **Whale Monitoring** - Background monitoring for transactions

### Thread Safety
- **sync.RWMutex** - Used for thread-safe access to shared data
- **Channels** - Used for goroutine communication
- **Context** - Used for cancellation and timeouts

## Error Handling

### Error Propagation
Errors bubble up from inner layers to outer layers:
```go
// Repository returns error
err := repo.Save(confession)
if err != nil {
    return err
}

// Use Case handles and wraps error
if err != nil {
    return fmt.Errorf("failed to save confession: %w", err)
}

// Delivery formats error for user
if err != nil {
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "❌ Failed to submit confession. Please try again.",
        },
    })
}
```

### Error Types
- **Domain Errors** - Business rule violations (entity layer)
- **Persistence Errors** - Data access failures (repository layer)
- **External Errors** - API/service failures (pkg layer)
- **User Errors** - Invalid input (delivery layer)

## Testing Strategy

### Unit Tests
- **Entity Layer** - Test business rules and validation
- **Use Case Layer** - Test business logic with mocked repositories
- **Repository Layer** - Test data persistence
- **Pkg Layer** - Test utility functions

### Integration Tests
- **Delivery + Use Case** - Test command handlers end-to-end
- **Use Case + Repository** - Test data flow

### Mocking
- Mock repositories for use case tests
- Mock AI providers for chatbot tests
- Mock Discord session for handler tests

## Scalability Considerations

### Horizontal Scaling
- Stateless design (session data in repository)
- Multiple bot instances can run simultaneously
- Shared data store (future: database instead of JSON)

### Vertical Scaling
- Efficient goroutine usage
- Connection pooling for external APIs
- Caching for frequently accessed data

## Future Enhancements

### Database Migration
Replace JSON files with PostgreSQL/MongoDB:
- Implement database repositories
- No changes needed to use cases or delivery
- Swap repository implementations via dependency injection

### Microservices
Break into services:
- Music service
- Moderation service (confessions)
- AI service
- Each service with its own database

### Event-Driven
Add event bus:
- Services communicate via events
- Better decoupling
- Easier to add new features

## Best Practices

1. **Keep It Simple** - Don't over-engineer
2. **Single Responsibility** - Each layer/package has one job
3. **Dependency Inversion** - Depend on abstractions, not concretions
4. **Interface Segregation** - Small, focused interfaces
5. **Open/Closed** - Open for extension, closed for modification

---

For more details on specific components, see the code documentation in each package.
