# NeruBot Architecture Guide

## Overview

NeruBot is built following **Clean Architecture** principles, ensuring maintainability, testability, and scalability. This document explains the system design, layer responsibilities, and key architectural decisions.

---

## 📐 Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Discord API                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   Delivery Layer (Discord)                   │
│  - Bot initialization & connection                           │
│  - Slash command handlers                                    │
│  - Event listeners                                           │
│  - Response formatting                                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                     Use Case Layer                           │
│  ┌─────────┐ ┌────────────┐ ┌───────┐ ┌─────────┐          │
│  │  Music  │ │ Confession │ │ Roast │ │ Chatbot │          │
│  │ Service │ │  Service   │ │Service│ │ Service │          │
│  └─────────┘ └────────────┘ └───────┘ └─────────┘          │
│  - Business logic                                            │
│  - Workflow orchestration                                    │
│  - Feature-specific operations                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Repository Layer                          │
│  - Data persistence abstraction                              │
│  - JSON file operations                                      │
│  - Future: Database operations                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Entity Layer                            │
│  - Domain models (Confession, Music, Roast, etc.)           │
│  - Pure business objects                                     │
│  - No external dependencies                                  │
└─────────────────────────────────────────────────────────────┘

         ┌──────────────────────────────────┐
         │     Infrastructure (pkg/)         │
         │  - AI Providers (DeepSeek)       │
         │  - FFmpeg wrapper                │
         │  - yt-dlp wrapper                │
         │  - Logger utility                │
         └──────────────────────────────────┘
```

---

## 🏗️ Layer Architecture

### 1. Entity Layer (`internal/entity/`)

**Purpose:** Define core business objects and domain models

**Responsibilities:**
- Pure data structures
- Business rules and validation
- No external dependencies
- Framework-agnostic

**Examples:**
```go
// entity/music.go
type Song struct {
    URL       string
    Title     string
    Duration  time.Duration
    Requester string
}

// entity/confession.go
type Confession struct {
    ID        int
    GuildID   string
    Content   string
    ImageURL  string
    Status    ConfessionStatus
    Timestamp time.Time
}
```

**Principles:**
- ✅ No imports from other layers
- ✅ Contains only domain logic
- ✅ Independent and reusable
- ✅ Easy to test

---

### 2. Use Case Layer (`internal/usecase/`)

**Purpose:** Implement business logic and orchestrate workflows

**Responsibilities:**
- Feature-specific operations
- Coordinate between entities and repositories
- Enforce business rules
- Handle complex workflows

**Structure:**
```
usecase/
├── music/
│   └── music_service.go       # Music streaming logic
├── confession/
│   └── confession_service.go  # Confession management
├── roast/
│   └── roast_service.go       # Roast generation
├── chatbot/
│   └── chatbot_service.go     # AI chatbot logic
├── news/
│   └── news_service.go        # News aggregation
└── whale/
    └── whale_service.go       # Whale alerts
```

**Example:**
```go
// usecase/music/music_service.go
type MusicService struct {
    queues map[string]*Queue
    ytdlp  *ytdlp.YtDlp
    ffmpeg *ffmpeg.FFmpeg
}

func (s *MusicService) Play(guildID, url string) error {
    // 1. Extract video info
    // 2. Add to queue
    // 3. Start playback if not playing
    // 4. Return result
}
```

**Principles:**
- ✅ Framework-independent
- ✅ Uses repository interfaces
- ✅ Testable with mocks
- ✅ Single Responsibility

---

### 3. Repository Layer (`internal/repository/`)

**Purpose:** Abstract data persistence operations

**Responsibilities:**
- File/database operations
- Data access abstraction
- CRUD operations
- Data format conversion

**Structure:**
```
repository/
├── repository.go               # Base repository interface
├── confession_repository.go    # Confession data access
└── roast_repository.go         # Roast data access
```

**Example:**
```go
// repository/confession_repository.go
type ConfessionRepository interface {
    Save(confession *entity.Confession) error
    FindByID(id int) (*entity.Confession, error)
    FindAll(guildID string) ([]*entity.Confession, error)
    Update(confession *entity.Confession) error
    Delete(id int) error
}

type JSONConfessionRepository struct {
    dataPath string
}

func (r *JSONConfessionRepository) Save(confession *entity.Confession) error {
    // Write to JSON file
}
```

**Principles:**
- ✅ Interface-based design
- ✅ Encapsulates storage details
- ✅ Easily swappable (JSON → DB)
- ✅ Thread-safe operations

---

### 4. Delivery Layer (`internal/delivery/discord/`)

**Purpose:** Handle external interfaces and framework-specific code

**Responsibilities:**
- Discord bot connection
- Command registration and handling
- Event listening
- Response formatting
- Error handling

**Structure:**
```
delivery/discord/
├── bot.go          # Bot initialization
└── handlers.go     # Command handlers
```

**Example:**
```go
// delivery/discord/bot.go
type Bot struct {
    session           *discordgo.Session
    config            *config.Config
    musicService      *music.MusicService
    confessionService *confession.ConfessionService
    roastService      *roast.RoastService
}

func (b *Bot) handlePlayCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
    // 1. Extract command options
    // 2. Call music service
    // 3. Format response
    // 4. Send to Discord
}
```

**Principles:**
- ✅ Framework-specific code isolated
- ✅ Thin layer (delegates to use cases)
- ✅ Converts external requests to domain operations
- ✅ Handles presentation logic

---

### 5. Infrastructure Layer (`internal/pkg/`)

**Purpose:** Provide shared utilities and external service wrappers

**Responsibilities:**
- External service integration
- Logging utilities
- Helper functions
- Third-party library wrappers

**Structure:**
```
pkg/
├── ai/
│   └── deepseek.go    # DeepSeek AI integration
├── ffmpeg/
│   └── ffmpeg.go      # FFmpeg wrapper
├── ytdlp/
│   └── ytdlp.go       # yt-dlp wrapper
└── logger/
    └── logger.go      # Logging utility
```

**Example:**
```go
// pkg/ytdlp/ytdlp.go
type YtDlp struct {
    path   string
    logger *logger.Logger
}

func (y *YtDlp) ExtractInfo(url string) (*VideoInfo, error) {
    // Execute yt-dlp command
    // Parse output
    // Return video information
}
```

**Principles:**
- ✅ Reusable across features
- ✅ Wraps external dependencies
- ✅ Provides clean interfaces
- ✅ Handles complexity

---

## 🔄 Data Flow

### Example: Music Play Command

```
1. User executes /play command in Discord
   ↓
2. Discord API → Delivery Layer (bot.go)
   - Receives interaction event
   - Extracts command parameters
   ↓
3. Delivery Layer → Use Case Layer (music_service.go)
   - Calls Play(guildID, url) method
   ↓
4. Use Case Layer → Infrastructure (ytdlp.go)
   - Extracts video information
   ↓
5. Use Case Layer → Repository Layer (optional)
   - Saves queue state
   ↓
6. Use Case Layer → Infrastructure (ffmpeg.go)
   - Starts audio streaming
   ↓
7. Use Case Layer → Delivery Layer
   - Returns success/error
   ↓
8. Delivery Layer → Discord API
   - Formats and sends response
   ↓
9. User sees confirmation message
```

---

## 🎯 Design Principles

### 1. Dependency Inversion Principle

**Rule:** High-level modules should not depend on low-level modules. Both should depend on abstractions.

**Implementation:**
- Use interfaces for repositories
- Inject dependencies via constructors
- Use cases don't know about Discord

```go
// Good: Use case depends on interface
type MusicService struct {
    repo Repository  // Interface, not concrete implementation
}

// Bad: Use case depends on concrete implementation
type MusicService struct {
    repo *JSONRepository  // Tight coupling
}
```

### 2. Single Responsibility Principle

**Rule:** A class should have only one reason to change.

**Implementation:**
- Each service handles one feature
- Separate command handlers
- Split complex logic into functions

```go
// Good: Single responsibility
type ConfessionService struct {
    // Only handles confession logic
}

type RoastService struct {
    // Only handles roast logic
}

// Bad: Multiple responsibilities
type BotService struct {
    // Handles confessions, roasts, music, etc.
}
```

### 3. Interface Segregation Principle

**Rule:** Clients should not be forced to depend on interfaces they don't use.

**Implementation:**
- Small, focused interfaces
- Feature-specific abstractions
- No "god" interfaces

```go
// Good: Focused interfaces
type ConfessionReader interface {
    FindByID(id int) (*Confession, error)
    FindAll(guildID string) ([]*Confession, error)
}

type ConfessionWriter interface {
    Save(confession *Confession) error
    Update(confession *Confession) error
}

// Bad: Large interface
type ConfessionRepository interface {
    Save(*Confession) error
    Update(*Confession) error
    Delete(int) error
    FindByID(int) (*Confession, error)
    FindAll(string) ([]*Confession, error)
    FindByStatus(string, Status) ([]*Confession, error)
    // ... 20 more methods
}
```

### 4. Open/Closed Principle

**Rule:** Software entities should be open for extension but closed for modification.

**Implementation:**
- Use interfaces for extensibility
- Plugin-based AI providers
- Strategy pattern for roast generation

```go
// Good: Open for extension
type AIProvider interface {
    Chat(prompt string) (string, error)
}

type DeepSeekProvider struct {}
func (d *DeepSeekProvider) Chat(prompt string) (string, error) { /*...*/ }

// Can add new providers without modifying existing code
type ClaudeProvider struct {}
func (c *ClaudeProvider) Chat(prompt string) (string, error) { /*...*/ }
```

---

## 🔧 Configuration Management

### Configuration Layer (`internal/config/`)

**Structure:**
```
config/
├── config.go       # Main configuration
├── constants.go    # Constants and defaults
└── messages.go     # Response messages
```

**Design:**
- Environment variable based
- Type-safe configuration
- Default values for all settings
- Validation on load

**Example:**
```go
// config/config.go
type Config struct {
    Bot      BotConfig
    Features FeatureFlags
    Limits   Limits
    Audio    AudioConfig
}

func Load() (*Config, error) {
    // Load from .env
    // Validate configuration
    // Return config or error
}
```

---

## 🧪 Testing Strategy

### Unit Testing

**Approach:**
- Test each layer independently
- Mock external dependencies
- Focus on business logic

```go
// Example: Testing use case
func TestMusicService_Play(t *testing.T) {
    mockRepo := &MockRepository{}
    mockYtdlp := &MockYtDlp{}
    
    service := music.NewMusicService(mockRepo, mockYtdlp)
    
    err := service.Play("guild123", "https://youtube.com/watch?v=...")
    
    assert.NoError(t, err)
    assert.True(t, mockYtdlp.ExtractInfoCalled)
}
```

### Integration Testing

**Approach:**
- Test layer interactions
- Use real dependencies (with test data)
- Verify end-to-end workflows

### Test Coverage Goals

- Entity Layer: >90%
- Use Case Layer: >80%
- Repository Layer: >70%
- Delivery Layer: >60%

---

## 🚀 Future Architecture Improvements

### Planned Enhancements

1. **Database Migration**
   - Move from JSON to PostgreSQL
   - Repository pattern already supports this
   - No changes to use case layer required

2. **Microservices Architecture**
   - Split services into separate deployments
   - gRPC for inter-service communication
   - API Gateway pattern

3. **Event-Driven Architecture**
   - Event bus for service communication
   - Async processing for heavy operations
   - Better scalability

4. **Caching Layer**
   - Redis for session management
   - Cache frequently accessed data
   - Improve performance

5. **Observability**
   - Structured logging (already implemented)
   - Metrics collection (Prometheus)
   - Distributed tracing
   - Health checks

---

## 📚 Best Practices

### Code Organization

1. **Package by Feature**
   - Group related code together
   - Each feature is self-contained
   - Easy to understand and maintain

2. **Clear Boundaries**
   - Strict layer separation
   - No circular dependencies
   - Use interfaces for communication

3. **Consistent Naming**
   - Entity: `Confession`, `Song`, `Roast`
   - Use Case: `ConfessionService`, `MusicService`
   - Repository: `ConfessionRepository`, `RoastRepository`

### Error Handling

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to play song: %w", err)
}

// Use custom error types for business logic
type ValidationError struct {
    Field   string
    Message string
}
```

### Logging

```go
// Use structured logging
log.Info("Song added to queue",
    "guild_id", guildID,
    "song_title", song.Title,
    "requester", requester,
)
```

---

## 🔗 Related Documentation

- [Project Structure](PROJECT_STRUCTURE.md) - Detailed file organization
- [Deployment Guide](DEPLOYMENT.md) - Production deployment
- [Contributing Guide](../CONTRIBUTING.md) - Development guidelines

---

**Last Updated:** December 6, 2025  
**Version:** 3.0.0  
**Author:** @nerufuyo
