# NeruBot - Python to Golang Migration Project Plan

## Project Information

**Project Name:** NeruBot Discord Bot - Python to Golang Migration  
**Version:** 1.0  
**Created:** November 9, 2025  
**Project Owner:** @nerufuyo  
**Methodology:** Agile/Iterative Development  
**Architecture:** Clean Architecture (Delivery → Use Case → Entity → Repository)

## Executive Summary

This project plan outlines the complete migration of NeruBot from Python (discord.py) to Golang (DiscordGo) using Clean Architecture principles. The migration will be executed in 5 phases over 10 weeks, maintaining feature parity while improving performance, reliability, and maintainability.

## Project Goals

### Primary Objectives
1. ✅ **Complete Migration** - Migrate all 8 features from Python to Golang
2. ✅ **Feature Parity** - Maintain 100% functionality of current Python implementation
3. ✅ **Clean Architecture** - Implement proper separation of concerns
4. ✅ **Performance** - Achieve 50% reduction in memory usage
5. ✅ **Reliability** - Improve error handling and recovery

### Secondary Objectives
1. 📚 **Documentation** - Comprehensive code documentation and guides
2. 🧪 **Testing** - Achieve >80% test coverage
3. 🚀 **Deployment** - Streamlined Docker and systemd deployment
4. 🔒 **Security** - Implement secure credential management
5. 📊 **Monitoring** - Add metrics and health checks

## Technology Stack

### Current (Python)
- **Language:** Python 3.8+
- **Discord Library:** discord.py 2.3.0
- **Concurrency:** asyncio
- **Dependencies:** 15+ packages
- **Runtime:** CPython interpreter

### Target (Golang)
- **Language:** Go 1.21+
- **Discord Library:** DiscordGo
- **Concurrency:** Goroutines + Channels
- **Dependencies:** ~10 packages (fewer, more stable)
- **Runtime:** Compiled binary (no interpreter)

### Supporting Tools
- **Version Control:** Git + GitHub
- **CI/CD:** GitHub Actions
- **Containerization:** Docker
- **Service Management:** systemd
- **Audio:** FFmpeg, yt-dlp (external binaries)

## Project Phases

### Phase 1: Foundation & Core Infrastructure (Week 1-2)
**Duration:** 2 weeks (40 hours)  
**Priority:** CRITICAL  
**Team:** 1 developer

#### Objectives
- Setup Go project structure
- Implement core packages
- Create base models
- Setup development environment

#### Deliverables
1. **Project Structure**
   - Go module initialization (`go.mod`)
   - Clean architecture directory structure
   - Package organization

2. **Core Packages**
   - `internal/config` - Configuration management
   - `internal/pkg/logger` - Structured logging
   - `internal/pkg/ffmpeg` - FFmpeg wrapper
   - `internal/pkg/ytdlp` - yt-dlp wrapper

3. **Base Models (Entity Layer)**
   - `entity/common.go` - Common types
   - `entity/error.go` - Custom error types
   - Basic model structures

4. **Development Setup**
   - `.env.example` file
   - `Makefile` for common tasks
   - VSCode/GoLand configuration
   - Pre-commit hooks

#### Tasks Breakdown
- [ ] Initialize Go module and project structure (4h)
- [ ] Implement config package with environment loading (6h)
- [ ] Create structured logger with levels and rotation (6h)
- [ ] Build FFmpeg wrapper with process management (8h)
- [ ] Build yt-dlp wrapper for music extraction (8h)
- [ ] Create base entity models (4h)
- [ ] Setup development tools and scripts (4h)

#### Success Criteria
- ✅ Clean architecture structure established
- ✅ All core utilities working and tested
- ✅ Configuration loading from .env
- ✅ Logger writing to console and file
- ✅ FFmpeg and yt-dlp executables detected

---

### Phase 2: High Priority Features (Week 3-5)
**Duration:** 3 weeks (96 hours)  
**Priority:** HIGH  
**Features:** Music, Confession, Help System

#### Feature 2.1: Music System (Week 3-4)
**Duration:** 40 hours  
**Complexity:** HIGH

##### Objectives
- Implement complete music playback system
- Support YouTube, Spotify, SoundCloud
- Queue management with loop modes
- Voice channel integration

##### Deliverables
1. **Entity Layer**
   - `entity/song.go` - Song model with metadata
   - `entity/queue.go` - Queue item model
   - `entity/source.go` - Source types enum

2. **Repository Layer** (if needed for playlists)
   - In-memory queue (no persistence for now)

3. **Use Case Layer**
   - `usecase/music/music_service.go` - Main music service
     - Queue management (add, remove, clear, shuffle)
     - Playback control (play, pause, resume, stop)
     - Loop modes (off, single, queue)
     - 24/7 mode management
   - `usecase/music/source_manager.go` - Multi-source handling
     - YouTube search and extraction
     - Spotify API with YouTube fallback
     - SoundCloud support
   - `usecase/music/voice_manager.go` - Voice connection handling
     - Join/leave voice channels
     - Audio streaming with DCA encoding
     - Voice state tracking

4. **Delivery Layer**
   - `delivery/discord/handlers/music_handler.go` - Music commands
     - `/play <song>` - Play or queue song
     - `/queue` - Show current queue
     - `/skip` - Skip current song
     - `/pause` / `/resume` - Control playback
     - `/stop` - Stop playback and clear queue
     - `/loop <mode>` - Set loop mode
     - `/247` - Toggle 24/7 mode
     - `/nowplaying` - Show current song
     - `/shuffle` - Shuffle queue

##### Tasks
- [ ] Implement song and queue entities (4h)
- [ ] Create YouTube source handler with yt-dlp (8h)
- [ ] Create Spotify source handler with API + fallback (8h)
- [ ] Create SoundCloud source handler (4h)
- [ ] Implement source manager with multi-source search (6h)
- [ ] Build music service with queue management (10h)
- [ ] Implement voice connection and streaming (12h)
- [ ] Create Discord command handlers (8h)
- [ ] Testing and debugging (8h)

##### Dependencies
- `github.com/bwmarrin/discordgo` - Discord API
- `github.com/bwmarrin/dca` - Audio encoding
- `github.com/zmb3/spotify` - Spotify API
- FFmpeg and yt-dlp binaries

#### Feature 2.2: Confession System (Week 4)
**Duration:** 24 hours  
**Complexity:** MEDIUM

##### Objectives
- Anonymous confession submission and posting
- Reply system with threading
- Queue-based async processing
- Guild-specific settings

##### Deliverables
1. **Entity Layer**
   - `entity/confession.go` - Confession model
   - `entity/confession_reply.go` - Reply model
   - `entity/confession_settings.go` - Guild settings
   - `entity/confession_status.go` - Status enum

2. **Repository Layer**
   - `repository/confession_repo.go` - JSON persistence
     - Save/load confessions
     - Save/load replies
     - Save/load guild settings
     - Thread-safe operations with sync.RWMutex
     - Atomic file writes

3. **Use Case Layer**
   - `usecase/confession/confession_service.go` - Main service
     - Create confession with ID generation
     - Post confession to channel
     - Add reply to confession
     - Manage cooldowns
     - Get confession by ID
   - `usecase/confession/queue_service.go` - Async queue
     - Channel-based queue
     - Worker pool for processing
     - Queue persistence on shutdown

4. **Delivery Layer**
   - `delivery/discord/handlers/confession_handler.go` - Commands
     - `/confess [image]` - Submit confession
     - `/reply <id> [image]` - Reply to confession
     - `/confession-setup <channel>` - Set channel (admin)
     - `/confession-stats` - View statistics

##### Tasks
- [ ] Create confession entity models (3h)
- [ ] Implement JSON repository with thread safety (6h)
- [ ] Build confession service with cooldowns (6h)
- [ ] Create queue service with workers (5h)
- [ ] Implement Discord handlers with modals (6h)
- [ ] Add image attachment support (3h)
- [ ] Testing and edge cases (4h)

##### Dependencies
- `encoding/json` - JSON operations
- `sync` - Thread safety

#### Feature 2.3: Help System (Week 5)
**Duration:** 12 hours  
**Complexity:** LOW

##### Objectives
- Interactive help with pagination
- Feature showcase
- Command reference

##### Deliverables
1. **Delivery Layer**
   - `delivery/discord/handlers/help_handler.go` - Help commands
     - `/help` - Main help with navigation
     - `/about` - Bot information
     - `/features` - Feature showcase
     - `/commands` - Command reference
   - Embed generation utilities
   - Button navigation handling

##### Tasks
- [ ] Create help content structure in config (2h)
- [ ] Implement embed generation utilities (3h)
- [ ] Build help command with pagination (4h)
- [ ] Create about, features, commands handlers (3h)
- [ ] Testing and polish (2h)

#### Phase 2 Success Criteria
- ✅ Music playback works across all sources
- ✅ Confessions can be submitted and replied to
- ✅ Help system provides clear documentation
- ✅ All commands respond within 100ms
- ✅ No memory leaks in continuous operation

---

### Phase 3: Medium Priority Features (Week 6-7)
**Duration:** 2 weeks (48 hours)  
**Priority:** MEDIUM  
**Features:** AI Chatbot, Roast System

#### Feature 3.1: AI Chatbot (Week 6)
**Duration:** 20 hours  
**Complexity:** MEDIUM

##### Objectives
- Multi-provider AI integration (Claude, Gemini, OpenAI)
- Smart fallback system
- Session management with timeout
- Personality-based responses

##### Deliverables
1. **Entity Layer**
   - `entity/chat_session.go` - Session model
   - `entity/chat_message.go` - Message model

2. **External Clients (pkg/ai)**
   - `pkg/ai/interface.go` - Common AI provider interface
   - `pkg/ai/claude.go` - Anthropic Claude client
   - `pkg/ai/gemini.go` - Google Gemini client
   - `pkg/ai/openai.go` - OpenAI client

3. **Use Case Layer**
   - `usecase/chatbot/chatbot_service.go` - Main service
     - Multi-provider fallback logic
     - Session management with TTL
     - Welcome and thank messages
     - Personality prompts

4. **Delivery Layer**
   - `delivery/discord/handlers/chatbot_handler.go` - Commands
     - `/chat <message>` - Send message to AI
     - `/reset-chat` - Reset session
   - `delivery/discord/events/message_handler.go` - Mention handler
     - Respond to @mentions
     - DM support

##### Tasks
- [ ] Define AI provider interface (2h)
- [ ] Implement Claude client with HTTP (4h)
- [ ] Implement Gemini client (4h)
- [ ] Implement OpenAI client (3h)
- [ ] Build chatbot service with fallback (5h)
- [ ] Add session management with cleanup (3h)
- [ ] Create Discord handlers (3h)
- [ ] Testing with all providers (3h)

##### Dependencies
- `net/http` - HTTP clients
- Provider API keys

#### Feature 3.2: Roast System (Week 7)
**Duration:** 28 hours  
**Complexity:** MEDIUM-HIGH

##### Objectives
- User behavior tracking
- AI-powered roast generation
- Activity pattern analysis
- Statistics and insights

##### Deliverables
1. **Entity Layer**
   - `entity/user_profile.go` - User activity profile
   - `entity/roast_pattern.go` - Roast template
   - `entity/activity_stats.go` - Statistics model

2. **Repository Layer**
   - `repository/roast_repo.go` - JSON persistence
     - User profiles
     - Activity logs
     - Roast history

3. **Use Case Layer**
   - `usecase/roast/roast_service.go` - Main service
     - Behavior analysis
     - Roast generation with AI
     - Pattern matching (8 categories)
     - Statistics calculation
   - `usecase/roast/activity_tracker.go` - Event tracking
     - Message tracking
     - Voice time tracking
     - Command usage tracking

4. **Delivery Layer**
   - `delivery/discord/handlers/roast_handler.go` - Commands
     - `/roast [target] [custom]` - Generate roast
     - `/roast-stats [user]` - View statistics
     - `/behavior-analysis [user]` - Detailed analysis
   - `delivery/discord/events/activity_listener.go` - Event listeners
     - Message events
     - Voice state events
     - Presence events

##### Tasks
- [ ] Create roast entity models (3h)
- [ ] Implement roast repository (5h)
- [ ] Build activity tracker with events (8h)
- [ ] Implement behavior analysis logic (6h)
- [ ] Create roast generation with AI (5h)
- [ ] Build Discord handlers (4h)
- [ ] Testing and cooldown logic (3h)

##### Dependencies
- AI service for roast generation
- Discord event listeners

#### Phase 3 Success Criteria
- ✅ AI chatbot responds with all 3 providers
- ✅ Fallback works when provider is down
- ✅ Roasts are generated based on real activity
- ✅ No performance impact from activity tracking

---

### Phase 4: Low Priority Features (Week 8)
**Duration:** 1 week (32 hours)  
**Priority:** LOW  
**Features:** News System, Whale Alerts

#### Feature 4.1: News System (Week 8, first half)
**Duration:** 16 hours  
**Complexity:** MEDIUM

##### Objectives
- Multi-source RSS news aggregation
- Auto-publishing with scheduler
- Manual control commands

##### Deliverables
1. **Entity Layer**
   - `entity/news_item.go` - News article model
   - `entity/news_source.go` - Source configuration

2. **Use Case Layer**
   - `usecase/news/news_service.go` - Main service
     - RSS feed fetching
     - News aggregation
     - Deduplication
     - Formatting
   - `usecase/news/scheduler.go` - Background scheduler
     - 10-minute ticker
     - Graceful shutdown
     - Channel management

3. **Delivery Layer**
   - `delivery/discord/handlers/news_handler.go` - Commands
     - `/news latest [count]` - Get latest news
     - `/news sources` - List sources
     - `/news set-channel <channel>` - Set channel (admin)
     - `/news start` / `/news stop` - Control auto-updates

##### Tasks
- [ ] Create news entity models (2h)
- [ ] Implement RSS feed fetching (4h)
- [ ] Build news aggregation logic (3h)
- [ ] Create background scheduler (4h)
- [ ] Implement Discord handlers (3h)
- [ ] Testing and edge cases (2h)

##### Dependencies
- `github.com/mmcdole/gofeed` - RSS parsing

#### Feature 4.2: Whale Alerts (Week 8, second half)
**Duration:** 16 hours  
**Complexity:** MEDIUM

##### Objectives
- Crypto whale transaction monitoring
- Guru tweet tracking
- Real-time alerts

##### Deliverables
1. **Entity Layer**
   - `entity/whale_transaction.go` - Transaction model
   - `entity/guru_tweet.go` - Tweet model

2. **External Clients**
   - `pkg/whale/whale_client.go` - Whale Alert API client
   - `pkg/whale/twitter_client.go` - Twitter API client

3. **Use Case Layer**
   - `usecase/whale/whale_service.go` - Whale monitoring
   - `usecase/whale/guru_service.go` - Guru tracking
   - Background polling with rate limiting

4. **Delivery Layer**
   - `delivery/discord/handlers/whale_handler.go` - Commands
     - `/whale setup [channel]` - Enable alerts
     - `/whale recent [limit]` - Show recent transactions
     - `/guru setup [channel]` - Enable guru tweets
     - `/guru accounts` - List monitored accounts

##### Tasks
- [ ] Create whale entity models (2h)
- [ ] Implement Whale Alert API client (4h)
- [ ] Implement Twitter API client (4h)
- [ ] Build whale monitoring service (3h)
- [ ] Create Discord handlers (3h)
- [ ] Testing and rate limiting (2h)

##### Dependencies
- Whale Alert API key
- Twitter API credentials
- `golang.org/x/time/rate` - Rate limiting

#### Phase 4 Success Criteria
- ✅ News updates post every 10 minutes
- ✅ Whale alerts trigger in real-time
- ✅ No duplicate news items
- ✅ Rate limits respected

---

### Phase 5: Testing, Deployment & Documentation (Week 9-10)
**Duration:** 2 weeks (40 hours)  
**Priority:** CRITICAL  
**Focus:** Quality Assurance, Production Readiness

#### Week 9: Testing & Quality Assurance
**Duration:** 24 hours

##### Objectives
- Comprehensive testing
- Performance validation
- Bug fixing

##### Deliverables
1. **Unit Tests**
   - Test all use cases (>80% coverage)
   - Mock repositories and external services
   - Edge case testing

2. **Integration Tests**
   - Test full command flow
   - Repository persistence tests
   - Discord interaction tests

3. **Performance Tests**
   - Memory usage profiling
   - Concurrent operation testing
   - Load testing with multiple guilds
   - Voice connection stability

4. **Bug Fixes**
   - Address all critical bugs
   - Fix memory leaks
   - Improve error handling

##### Tasks
- [ ] Write unit tests for all use cases (12h)
- [ ] Create integration test suite (6h)
- [ ] Run performance profiling (4h)
- [ ] Fix identified bugs and issues (8h)
- [ ] Code review and refactoring (4h)

##### Tools
- `testing` - Go standard testing
- `github.com/stretchr/testify` - Assertions
- `go test -cover` - Coverage reports
- `pprof` - Performance profiling

#### Week 10: Deployment & Documentation
**Duration:** 16 hours

##### Objectives
- Production deployment setup
- Complete documentation
- Migration guide

##### Deliverables
1. **Deployment Configuration**
   - `Dockerfile` - Multi-stage build
   - `docker-compose.yml` - Service orchestration
   - `deploy/systemd/nerubot.service` - Systemd service
   - `deploy/nginx/nerubot.conf` - Optional reverse proxy
   - CI/CD pipeline with GitHub Actions

2. **Documentation**
   - Update `README.md` for Go version
   - `docs/INSTALLATION.md` - Installation guide
   - `docs/DEVELOPMENT.md` - Development guide
   - `docs/API.md` - Internal API documentation
   - `docs/MIGRATION.md` - Migration from Python guide
   - Code documentation with godoc

3. **Deployment Scripts**
   - `deploy/setup.sh` - Automated VPS setup
   - `deploy/update.sh` - Update script
   - `deploy/monitor.sh` - Monitoring script
   - `deploy/backup.sh` - Backup script

4. **Migration Process**
   - Data migration scripts
   - Rollback plan
   - Validation checklist

##### Tasks
- [ ] Create production Dockerfile (3h)
- [ ] Update deployment scripts for Go (4h)
- [ ] Setup CI/CD pipeline (3h)
- [ ] Write comprehensive README (2h)
- [ ] Create migration guide (2h)
- [ ] Add code documentation (3h)
- [ ] Final production deployment (4h)

#### Phase 5 Success Criteria
- ✅ >80% test coverage achieved
- ✅ No critical bugs in production
- ✅ Performance targets met (50% memory reduction)
- ✅ Documentation complete and accurate
- ✅ Successful production deployment
- ✅ Zero downtime migration

---

## Resource Requirements

### Development Resources
- **Team:** 1 Senior Go Developer
- **Time:** 10 weeks (200 hours)
- **Hardware:** Development machine with 8GB+ RAM
- **Software:** 
  - Go 1.21+
  - Docker
  - Git
  - VSCode/GoLand
  - Discord Desktop (testing)

### Infrastructure
- **Testing:**
  - Discord test server
  - VPS for testing (2GB RAM, 1 vCPU)
  - Test bot token
  
- **Production:**
  - VPS (4GB RAM, 2 vCPU recommended)
  - FFmpeg installed
  - yt-dlp binary
  - 10GB+ storage

### External Services & API Keys
- Discord Bot Token (required)
- Spotify API credentials (optional but recommended)
- OpenAI API key (optional)
- Anthropic Claude API key (optional)
- Google Gemini API key (optional)
- Whale Alert API key (optional)
- Twitter API credentials (optional)

### Budget Estimate
- **Development:** 200 hours × $50/hr = $10,000
- **VPS Hosting:** $20/month (testing + production)
- **API Credits:** ~$50/month (AI APIs, if used)
- **Total:** ~$10,000 + $70/month

---

## Risk Management

### Technical Risks

#### Risk 1: Audio Streaming Complexity
**Probability:** HIGH  
**Impact:** HIGH  
**Mitigation:**
- Use proven DCA library
- Prototype early in Phase 2
- Allocate extra time for debugging
- Fallback to simple playback if needed

#### Risk 2: DiscordGo API Differences
**Probability:** MEDIUM  
**Impact:** MEDIUM  
**Mitigation:**
- Study DiscordGo documentation thoroughly
- Create proof-of-concept before full implementation
- Join DiscordGo community for support
- Budget time for API learning curve

#### Risk 3: Data Migration Issues
**Probability:** MEDIUM  
**Impact:** HIGH  
**Mitigation:**
- Create comprehensive backup before migration
- Write validation scripts for data integrity
- Test migration on copy of production data
- Have rollback plan ready

#### Risk 4: Performance Bottlenecks
**Probability:** LOW  
**Impact:** MEDIUM  
**Mitigation:**
- Profile early and often
- Use goroutines efficiently
- Implement proper resource cleanup
- Load test before production

### Project Risks

#### Risk 5: Scope Creep
**Probability:** MEDIUM  
**Impact:** MEDIUM  
**Mitigation:**
- Strict adherence to feature parity
- No new features until migration complete
- Clear phase boundaries
- Regular progress reviews

#### Risk 6: API Key Dependencies
**Probability:** LOW  
**Impact:** LOW  
**Mitigation:**
- Make all external APIs optional
- Implement graceful degradation
- Clear documentation on required vs optional keys
- Test without optional APIs

---

## Quality Assurance

### Code Quality Standards
- **Formatting:** `gofmt` and `goimports`
- **Linting:** `golangci-lint` with strict rules
- **Documentation:** godoc comments on all exported items
- **Error Handling:** Proper error wrapping and context
- **Testing:** Unit tests for all business logic
- **Coverage:** Minimum 80% test coverage

### Performance Targets
- **Memory:** <100MB base usage (vs ~200MB Python)
- **CPU:** <10% idle, <50% during music playback
- **Response Time:** <100ms for slash commands
- **Startup Time:** <5 seconds
- **Voice Latency:** <200ms

### Reliability Targets
- **Uptime:** 99.9% (less than 8.7 hours downtime/year)
- **Error Rate:** <0.1% of commands
- **Recovery:** Automatic restart on crashes
- **Data Integrity:** Zero data loss on restarts

---

## Communication Plan

### Status Updates
- **Daily:** Brief progress notes in commit messages
- **Weekly:** Phase completion summary
- **Milestones:** Detailed report at end of each phase

### Documentation
- **Code Comments:** Inline for complex logic
- **Commit Messages:** Follow format-commit.md guidelines
- **README:** Keep updated with progress
- **CHANGELOG:** Log all significant changes

---

## Success Metrics

### Functional Metrics
- ✅ All 8 features migrated and working
- ✅ 100% feature parity with Python version
- ✅ All slash commands functional
- ✅ Data successfully migrated

### Performance Metrics
- ✅ 50% reduction in memory usage
- ✅ Faster startup time
- ✅ Lower CPU usage
- ✅ Better voice quality

### Quality Metrics
- ✅ >80% test coverage
- ✅ Zero critical bugs in production
- ✅ Clean architecture maintained
- ✅ Well-documented codebase

### User Experience Metrics
- ✅ Command response time <100ms
- ✅ No functionality regressions
- ✅ Smooth migration (zero downtime)
- ✅ Positive user feedback

---

## Timeline Summary

| Phase | Duration | Dates | Key Deliverables |
|-------|----------|-------|------------------|
| **Phase 1: Foundation** | 2 weeks | Week 1-2 | Core packages, project structure |
| **Phase 2: High Priority** | 3 weeks | Week 3-5 | Music, Confession, Help |
| **Phase 3: Medium Priority** | 2 weeks | Week 6-7 | Chatbot, Roast |
| **Phase 4: Low Priority** | 1 week | Week 8 | News, Whale Alerts |
| **Phase 5: Testing & Deploy** | 2 weeks | Week 9-10 | Testing, deployment, docs |
| **Total** | **10 weeks** | | **All features migrated** |

---

## Next Steps

### Immediate Actions (Week 1)
1. ✅ Review and approve this project plan
2. ⏳ Create detailed implementation tickets
3. ⏳ Setup Go development environment
4. ⏳ Initialize Go module and project structure
5. ⏳ Begin Phase 1: Core infrastructure implementation

### Weekly Checklist
- [ ] Review completed tasks
- [ ] Update project plan if needed
- [ ] Commit all changes with proper messages
- [ ] Run tests and fix bugs
- [ ] Update documentation
- [ ] Plan next week's tasks

---

## Appendix

### A. Key Dependencies
```
github.com/bwmarrin/discordgo      # Discord API
github.com/bwmarrin/dca            # Audio encoding
github.com/zmb3/spotify            # Spotify API
github.com/mmcdole/gofeed          # RSS parsing
github.com/joho/godotenv           # Environment variables
github.com/stretchr/testify        # Testing assertions
golang.org/x/time/rate             # Rate limiting
```

### B. Useful Commands
```bash
# Development
go run cmd/nerubot/main.go         # Run bot
go test ./...                      # Run all tests
go test -cover ./...               # Test with coverage
go build -o nerubot cmd/nerubot/main.go  # Build binary

# Docker
docker build -t nerubot:latest .   # Build image
docker-compose up -d               # Start services

# Deployment
./deploy/setup.sh                  # Setup VPS
./deploy/update.sh                 # Update bot
systemctl status nerubot           # Check status
```

### C. References
- [DiscordGo Documentation](https://github.com/bwmarrin/discordgo)
- [Clean Architecture by Robert Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Effective Go](https://go.dev/doc/effective_go)
- [Project Breakdown](./project-breakdown.md)

---

**Document Version:** 1.0  
**Last Updated:** November 9, 2025  
**Status:** Approved - Ready for Implementation  
**Approved By:** @nerufuyo
