# NeruBot Microservices Migration Plan

## Project Information
- **Project:** NeruBot Microservices Architecture
- **Start Date:** December 5, 2025
- **Target Completion:** December 26, 2025
- **Platform:** Railway
- **Current State:** Monolithic Go application
- **Target State:** Microservices with API Gateway

---

## Table of Contents
1. [Migration Strategy](#migration-strategy)
2. [Architecture Overview](#architecture-overview)
3. [Service Specifications](#service-specifications)
4. [Database Design](#database-design)
5. [Railway Deployment](#railway-deployment)
6. [Implementation Phases](#implementation-phases)
7. [Testing Strategy](#testing-strategy)
8. [Rollback Plan](#rollback-plan)

---

## Migration Strategy

### Approach: Strangler Fig Pattern
We'll use the Strangler Fig pattern to gradually migrate from monolithic to microservices:
1. Build new microservices alongside the monolith
2. Route new features to microservices
3. Gradually migrate existing functionality
4. Decommission monolithic components incrementally

### Key Principles
- ✅ **Zero Downtime:** Maintain service availability throughout migration
- ✅ **Backward Compatibility:** Ensure existing features continue working
- ✅ **Incremental Rollout:** Deploy and test one service at a time
- ✅ **Data Consistency:** Maintain data integrity during migration
- ✅ **Monitoring First:** Implement observability before migration

---

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Discord API                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                     API Gateway Service                      │
│  - Discord Bot Handler                                       │
│  - Command Router                                            │
│  - Response Aggregator                                       │
│  Port: 8080                                                  │
└─────────┬──────┬──────┬──────┬──────┬──────┬──────┬─────────┘
          │      │      │      │      │      │      │
          ▼      ▼      ▼      ▼      ▼      ▼      ▼
    ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐
    │Music│ │Conf │ │Roast│ │Chat │ │News │ │Whale│ │Shared│
    │:8081│ │:8082│ │:8083│ │:8084│ │:8085│ │:8086│ │:8087│
    └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘
       │       │       │       │       │       │       │
       ▼       ▼       ▼       ▼       ▼       ▼       ▼
    ┌────────────────────────────────────────────────────┐
    │              PostgreSQL Database                   │
    │  - music_db                                        │
    │  - confession_db                                   │
    │  - roast_db                                        │
    │  - chat_db                                         │
    │  - news_db                                         │
    │  - whale_db                                        │
    └────────────────────────────────────────────────────┘
                          │
                          ▼
                    ┌──────────┐
                    │  Redis   │
                    │  Cache   │
                    └──────────┘
```

### Communication Protocol
- **Internal:** gRPC (high performance, type-safe)
- **External:** REST API (Discord webhook compatibility)
- **Caching:** Redis for session management and rate limiting

---

## Service Specifications

### 1. API Gateway Service

**Purpose:** Main entry point for Discord bot, routes commands to appropriate services

**Responsibilities:**
- Discord WebSocket connection management
- Slash command registration and handling
- Request routing to microservices
- Response aggregation and formatting
- Rate limiting and authentication
- Health check aggregation

**Tech Stack:**
- Language: Go
- Framework: DiscordGo
- Port: 8080
- Database: None (stateless)

**Endpoints:**
- `POST /gateway/health` - Health check
- `POST /gateway/metrics` - Metrics endpoint
- WebSocket: Discord Gateway connection

**Environment Variables:**
```env
DISCORD_TOKEN=xxx
DISCORD_GUILD_ID=xxx
MUSIC_SERVICE_URL=music-service:8081
CONFESSION_SERVICE_URL=confession-service:8082
ROAST_SERVICE_URL=roast-service:8083
CHATBOT_SERVICE_URL=chatbot-service:8084
NEWS_SERVICE_URL=news-service:8085
WHALE_SERVICE_URL=whale-service:8086
REDIS_URL=redis://redis:6379
LOG_LEVEL=INFO
```

**Docker:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gateway .
EXPOSE 8080
CMD ["./gateway"]
```

---

### 2. Music Service

**Purpose:** Handle music streaming, queue management, and playback control

**Responsibilities:**
- YouTube audio extraction (yt-dlp)
- FFmpeg audio processing
- Queue management (per guild)
- Playback controls (play, pause, skip, stop)
- Voice channel connection
- Loop modes and shuffle

**Tech Stack:**
- Language: Go
- Dependencies: FFmpeg, yt-dlp
- Port: 8081
- Database: PostgreSQL (queue persistence)
- Cache: Redis (current playing state)

**gRPC Service Definition:**
```protobuf
service MusicService {
  rpc Play(PlayRequest) returns (PlayResponse);
  rpc Pause(PauseRequest) returns (PauseResponse);
  rpc Skip(SkipRequest) returns (SkipResponse);
  rpc Stop(StopRequest) returns (StopResponse);
  rpc Queue(QueueRequest) returns (QueueResponse);
  rpc NowPlaying(NowPlayingRequest) returns (NowPlayingResponse);
}
```

**Database Schema:**
```sql
CREATE TABLE music_queues (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    song_url TEXT NOT NULL,
    song_title VARCHAR(255),
    requested_by VARCHAR(32),
    added_at TIMESTAMP DEFAULT NOW(),
    position INTEGER,
    status VARCHAR(20) DEFAULT 'queued'
);

CREATE TABLE playback_state (
    guild_id VARCHAR(32) PRIMARY KEY,
    current_song_id INTEGER REFERENCES music_queues(id),
    is_playing BOOLEAN DEFAULT false,
    loop_mode VARCHAR(20) DEFAULT 'none',
    volume INTEGER DEFAULT 100,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Environment Variables:**
```env
DATABASE_URL=postgresql://user:pass@postgres:5432/music_db
REDIS_URL=redis://redis:6379
FFMPEG_PATH=/usr/bin/ffmpeg
YTDLP_PATH=/usr/bin/yt-dlp
MAX_QUEUE_SIZE=100
PORT=8081
```

---

### 3. Confession Service

**Purpose:** Manage anonymous confessions, moderation queue, and replies

**Responsibilities:**
- Confession submission (anonymous)
- Image attachment handling
- Moderation queue (approve/reject)
- Reply system
- Per-guild settings
- Confession numbering

**Tech Stack:**
- Language: Go
- Port: 8082
- Database: PostgreSQL
- Storage: S3-compatible (Railway volumes)

**gRPC Service Definition:**
```protobuf
service ConfessionService {
  rpc Submit(SubmitRequest) returns (SubmitResponse);
  rpc Approve(ApproveRequest) returns (ApproveResponse);
  rpc Reject(RejectRequest) returns (RejectResponse);
  rpc Reply(ReplyRequest) returns (ReplyResponse);
  rpc GetSettings(GetSettingsRequest) returns (GetSettingsResponse);
  rpc UpdateSettings(UpdateSettingsRequest) returns (UpdateSettingsResponse);
}
```

**Database Schema:**
```sql
CREATE TABLE confessions (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    confession_number INTEGER,
    content TEXT NOT NULL,
    image_url TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    submitted_at TIMESTAMP DEFAULT NOW(),
    approved_at TIMESTAMP,
    channel_id VARCHAR(32),
    message_id VARCHAR(32)
);

CREATE TABLE confession_replies (
    id SERIAL PRIMARY KEY,
    confession_id INTEGER REFERENCES confessions(id),
    reply_content TEXT NOT NULL,
    replied_at TIMESTAMP DEFAULT NOW(),
    message_id VARCHAR(32)
);

CREATE TABLE confession_settings (
    guild_id VARCHAR(32) PRIMARY KEY,
    enabled BOOLEAN DEFAULT true,
    mod_channel_id VARCHAR(32),
    confession_channel_id VARCHAR(32),
    require_approval BOOLEAN DEFAULT true,
    allow_images BOOLEAN DEFAULT true,
    max_length INTEGER DEFAULT 2000
);
```

**Environment Variables:**
```env
DATABASE_URL=postgresql://user:pass@postgres:5432/confession_db
STORAGE_URL=s3://bucket/confessions
MAX_IMAGE_SIZE=10485760
PORT=8082
```

---

### 4. Roast Service

**Purpose:** Track user activity and generate personalized roasts

**Responsibilities:**
- Activity tracking (messages, reactions, voice)
- Pattern detection (spammer, lurker, etc.)
- Roast generation based on profile
- Statistics and leaderboards
- Cooldown management

**Tech Stack:**
- Language: Go
- Port: 8083
- Database: PostgreSQL
- Cache: Redis (cooldowns)

**gRPC Service Definition:**
```protobuf
service RoastService {
  rpc TrackActivity(ActivityRequest) returns (ActivityResponse);
  rpc GenerateRoast(RoastRequest) returns (RoastResponse);
  rpc GetProfile(ProfileRequest) returns (ProfileResponse);
  rpc GetLeaderboard(LeaderboardRequest) returns (LeaderboardResponse);
  rpc GetStats(StatsRequest) returns (StatsResponse);
}
```

**Database Schema:**
```sql
CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    user_id VARCHAR(32) NOT NULL,
    message_count INTEGER DEFAULT 0,
    reaction_count INTEGER DEFAULT 0,
    voice_minutes INTEGER DEFAULT 0,
    command_count INTEGER DEFAULT 0,
    last_seen TIMESTAMP DEFAULT NOW(),
    UNIQUE(guild_id, user_id)
);

CREATE TABLE roast_history (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    user_id VARCHAR(32) NOT NULL,
    roast_content TEXT NOT NULL,
    roast_category VARCHAR(50),
    roasted_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roast_patterns (
    id SERIAL PRIMARY KEY,
    pattern_name VARCHAR(50) UNIQUE,
    pattern_description TEXT,
    min_messages INTEGER,
    min_reactions INTEGER,
    min_voice_minutes INTEGER,
    roast_templates TEXT[]
);
```

**Environment Variables:**
```env
DATABASE_URL=postgresql://user:pass@postgres:5432/roast_db
REDIS_URL=redis://redis:6379
ROAST_COOLDOWN_SECONDS=300
PORT=8083
```

---

### 5. Chatbot Service

**Purpose:** AI-powered chatbot with multi-provider support

**Responsibilities:**
- Multi-provider AI integration (OpenAI, Claude, Gemini)
- Automatic provider fallback
- Session management (30-min timeout)
- Context-aware conversations
- Token tracking

**Tech Stack:**
- Language: Go
- Port: 8084
- Database: PostgreSQL (session history)
- Cache: Redis (active sessions)

**gRPC Service Definition:**
```protobuf
service ChatbotService {
  rpc Chat(ChatRequest) returns (ChatResponse);
  rpc GetSession(SessionRequest) returns (SessionResponse);
  rpc ClearSession(ClearSessionRequest) returns (ClearSessionResponse);
  rpc GetProviderStatus(ProviderStatusRequest) returns (ProviderStatusResponse);
}
```

**Database Schema:**
```sql
CREATE TABLE chat_sessions (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    user_id VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    last_interaction TIMESTAMP DEFAULT NOW(),
    message_count INTEGER DEFAULT 0
);

CREATE TABLE chat_messages (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES chat_sessions(id),
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    provider VARCHAR(20),
    tokens_used INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE provider_stats (
    provider_name VARCHAR(20) PRIMARY KEY,
    total_requests INTEGER DEFAULT 0,
    successful_requests INTEGER DEFAULT 0,
    failed_requests INTEGER DEFAULT 0,
    total_tokens INTEGER DEFAULT 0,
    last_used TIMESTAMP
);
```

**Environment Variables:**
```env
DATABASE_URL=postgresql://user:pass@postgres:5432/chat_db
REDIS_URL=redis://redis:6379
OPENAI_API_KEY=xxx
ANTHROPIC_API_KEY=xxx
GEMINI_API_KEY=xxx
SESSION_TIMEOUT_MINUTES=30
MAX_CONTEXT_MESSAGES=10
PORT=8084
```

---

### 6. News Service

**Purpose:** RSS feed aggregation and news publishing

**Responsibilities:**
- RSS feed fetching from multiple sources
- Concurrent news aggregation
- Article deduplication
- Auto-publishing to Discord
- Source management

**Tech Stack:**
- Language: Go
- Port: 8085
- Database: PostgreSQL
- Cache: Redis (published articles)

**gRPC Service Definition:**
```protobuf
service NewsService {
  rpc FetchNews(FetchNewsRequest) returns (FetchNewsResponse);
  rpc AddSource(AddSourceRequest) returns (AddSourceResponse);
  rpc RemoveSource(RemoveSourceRequest) returns (RemoveSourceResponse);
  rpc GetSources(GetSourcesRequest) returns (GetSourcesResponse);
  rpc PublishNews(PublishNewsRequest) returns (PublishNewsResponse);
}
```

**Database Schema:**
```sql
CREATE TABLE news_sources (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    source_name VARCHAR(100),
    source_url TEXT NOT NULL,
    feed_type VARCHAR(20) DEFAULT 'rss',
    enabled BOOLEAN DEFAULT true,
    last_fetched TIMESTAMP
);

CREATE TABLE news_articles (
    id SERIAL PRIMARY KEY,
    source_id INTEGER REFERENCES news_sources(id),
    article_title VARCHAR(255),
    article_url TEXT UNIQUE,
    article_content TEXT,
    published_at TIMESTAMP,
    fetched_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE published_articles (
    id SERIAL PRIMARY KEY,
    article_id INTEGER REFERENCES news_articles(id),
    guild_id VARCHAR(32),
    channel_id VARCHAR(32),
    message_id VARCHAR(32),
    published_at TIMESTAMP DEFAULT NOW()
);
```

**Environment Variables:**
```env
DATABASE_URL=postgresql://user:pass@postgres:5432/news_db
REDIS_URL=redis://redis:6379
FETCH_INTERVAL_MINUTES=15
MAX_ARTICLES_PER_FETCH=50
PORT=8085
```

---

### 7. Whale Service

**Purpose:** Cryptocurrency whale transaction monitoring

**Responsibilities:**
- Monitor large crypto transactions
- Real-time alerts
- Configurable thresholds
- Multi-blockchain support
- Transaction history

**Tech Stack:**
- Language: Go
- Port: 8086
- Database: PostgreSQL
- External API: Whale Alert API

**gRPC Service Definition:**
```protobuf
service WhaleService {
  rpc MonitorTransactions(MonitorRequest) returns (stream Transaction);
  rpc GetTransactionHistory(HistoryRequest) returns (HistoryResponse);
  rpc UpdateThreshold(ThresholdRequest) returns (ThresholdResponse);
  rpc GetAlertSettings(AlertSettingsRequest) returns (AlertSettingsResponse);
}
```

**Database Schema:**
```sql
CREATE TABLE whale_transactions (
    id SERIAL PRIMARY KEY,
    blockchain VARCHAR(20),
    transaction_hash VARCHAR(128) UNIQUE,
    amount DECIMAL(20,2),
    amount_usd DECIMAL(20,2),
    from_address VARCHAR(128),
    to_address VARCHAR(128),
    timestamp TIMESTAMP,
    detected_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE whale_alerts (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    channel_id VARCHAR(32),
    min_amount_usd DECIMAL(20,2) DEFAULT 1000000,
    enabled BOOLEAN DEFAULT true,
    blockchains TEXT[]
);

CREATE TABLE alerted_transactions (
    transaction_id INTEGER REFERENCES whale_transactions(id),
    guild_id VARCHAR(32),
    message_id VARCHAR(32),
    alerted_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (transaction_id, guild_id)
);
```

**Environment Variables:**
```env
DATABASE_URL=postgresql://user:pass@postgres:5432/whale_db
WHALE_ALERT_API_KEY=xxx
CHECK_INTERVAL_SECONDS=60
DEFAULT_MIN_AMOUNT_USD=1000000
PORT=8086
```

---

### 8. Shared Service (Optional)

**Purpose:** Shared utilities and common functionality

**Responsibilities:**
- Centralized logging
- Metrics collection
- Configuration management
- Health checks
- Service discovery

**Tech Stack:**
- Language: Go
- Port: 8087
- Database: None

---

## Database Design

### Database Strategy

**Option 1: Database per Service (Recommended)**
- Each service has its own PostgreSQL database
- Complete data isolation
- Independent scaling
- Railway: Multiple PostgreSQL instances

**Option 2: Shared Database with Schemas**
- Single PostgreSQL with multiple schemas
- Cost-effective
- Simpler management
- Railway: Single PostgreSQL instance

**Selected:** Option 1 (Database per Service)

### Migration from JSON to PostgreSQL

**Phase 1: Dual Write**
1. Write to both JSON and PostgreSQL
2. Read from JSON (existing behavior)
3. Validate data consistency

**Phase 2: Dual Read**
1. Write to both JSON and PostgreSQL
2. Read from PostgreSQL with JSON fallback
3. Monitor for issues

**Phase 3: PostgreSQL Only**
1. Write to PostgreSQL only
2. Read from PostgreSQL only
3. Archive JSON files

**Phase 4: Cleanup**
1. Remove JSON read/write code
2. Delete JSON files
3. Complete migration

---

## Railway Deployment

### Railway Configuration

**Project Structure:**
```
nerubot-api-gateway/        # Service 1
nerubot-music-service/      # Service 2
nerubot-confession-service/ # Service 3
nerubot-roast-service/      # Service 4
nerubot-chatbot-service/    # Service 5
nerubot-news-service/       # Service 6
nerubot-whale-service/      # Service 7
nerubot-postgres/           # PostgreSQL
nerubot-redis/              # Redis
```

### Railway Services Setup

**1. API Gateway Service**
```toml
# railway.toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "services/gateway/Dockerfile"

[deploy]
healthcheckPath = "/health"
healthcheckTimeout = 100
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 10

[[services]]
name = "api-gateway"
port = 8080

[env]
DISCORD_TOKEN = "${{secrets.DISCORD_TOKEN}}"
MUSIC_SERVICE_URL = "${{services.music-service.url}}"
CONFESSION_SERVICE_URL = "${{services.confession-service.url}}"
```

**2. Music Service**
```toml
# railway.toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "services/music/Dockerfile"

[deploy]
healthcheckPath = "/health"
healthcheckTimeout = 100

[[services]]
name = "music-service"
port = 8081

[env]
DATABASE_URL = "${{Postgres.DATABASE_URL}}"
REDIS_URL = "${{Redis.REDIS_URL}}"
```

**Similar configuration for other services...**

### Environment Variables Management

**Shared Variables (Railway Project Level):**
- `DISCORD_TOKEN`
- `DISCORD_GUILD_ID`
- `LOG_LEVEL`
- `ENVIRONMENT` (production/staging)

**Service-Specific Variables:**
- Each service has its own database credentials
- API keys for external services
- Service-specific configuration

### Networking

**Internal Communication:**
- Railway provides private networking
- Services communicate via internal DNS
- Format: `service-name.railway.internal:port`

**External Access:**
- Only API Gateway exposes public endpoint
- All other services are private

### Monitoring and Logging

**Railway Built-in:**
- Service metrics (CPU, Memory, Network)
- Application logs
- Deployment logs

**Custom Monitoring:**
- Health check endpoints on all services
- Prometheus metrics export
- Grafana dashboards (optional)

**Log Aggregation:**
- Structured JSON logging
- Railway log viewer
- Log retention: 30 days

---

## Implementation Phases

### Phase 1: Project Setup (Days 1-2)

**Tasks:**
1. ✅ Create project documentation
2. ☐ Setup Railway project
3. ☐ Create repository structure
4. ☐ Setup CI/CD pipeline
5. ☐ Configure development environment

**Deliverables:**
- Project brief document
- Migration plan document
- Railway project created
- GitHub Actions workflows
- Development docker-compose.yml

**Success Criteria:**
- All documentation complete
- Railway project accessible
- Local development environment working

---

### Phase 2: Shared Infrastructure (Days 3-4)

**Tasks:**
1. ☐ Setup PostgreSQL on Railway
2. ☐ Setup Redis on Railway
3. ☐ Create database schemas
4. ☐ Implement shared packages
5. ☐ Create proto definitions for gRPC

**Deliverables:**
- PostgreSQL databases for each service
- Redis instance configured
- Shared Go packages (logger, config, metrics)
- gRPC proto files compiled

**Success Criteria:**
- Databases accessible and tested
- Redis working correctly
- Shared packages reusable across services

---

### Phase 3: API Gateway Service (Days 5-6)

**Tasks:**
1. ☐ Create gateway service structure
2. ☐ Implement Discord bot connection
3. ☐ Setup command routing
4. ☐ Implement service discovery
5. ☐ Add health checks and metrics
6. ☐ Deploy to Railway

**Deliverables:**
- API Gateway service running
- Discord bot connected
- Command routing framework
- Railway deployment successful

**Success Criteria:**
- Bot responds to Discord events
- Commands route to placeholder handlers
- Health checks returning 200
- Railway deployment stable

---

### Phase 4: Music Service (Days 7-8)

**Tasks:**
1. ☐ Extract music logic to separate service
2. ☐ Implement gRPC server
3. ☐ Migrate queue management
4. ☐ Integrate FFmpeg and yt-dlp
5. ☐ Database migration (JSON → PostgreSQL)
6. ☐ Deploy to Railway
7. ☐ Integration testing

**Deliverables:**
- Music service deployed
- Database migrated
- gRPC endpoints working
- Integration with API Gateway

**Success Criteria:**
- Music playback working
- Queue management functional
- Database persistence working
- No data loss from migration

---

### Phase 5: Confession Service (Days 9-10)

**Tasks:**
1. ☐ Extract confession logic
2. ☐ Implement gRPC server
3. ☐ Setup image storage
4. ☐ Database migration
5. ☐ Implement moderation queue
6. ☐ Deploy to Railway
7. ☐ Integration testing

**Deliverables:**
- Confession service deployed
- Image storage working
- Database migrated
- Moderation flow functional

**Success Criteria:**
- Anonymous submissions working
- Image attachments supported
- Moderation queue functional
- Reply system working

---

### Phase 6: Roast Service (Days 11-12)

**Tasks:**
1. ☐ Extract roast logic
2. ☐ Implement gRPC server
3. ☐ Activity tracking system
4. ☐ Database migration
5. ☐ Pattern detection logic
6. ☐ Deploy to Railway
7. ☐ Integration testing

**Deliverables:**
- Roast service deployed
- Activity tracking working
- Database migrated
- Statistics functional

**Success Criteria:**
- Activity tracking accurate
- Roast generation working
- Cooldowns enforced
- Statistics displaying correctly

---

### Phase 7: Chatbot Service (Days 13-14)

**Tasks:**
1. ☐ Extract chatbot logic
2. ☐ Implement gRPC server
3. ☐ Multi-provider integration
4. ☐ Session management
5. ☐ Database setup
6. ☐ Deploy to Railway
7. ☐ Integration testing

**Deliverables:**
- Chatbot service deployed
- All AI providers integrated
- Session management working
- Fallback mechanism functional

**Success Criteria:**
- Chat conversations working
- Provider fallback functional
- Sessions timeout correctly
- Context maintained

---

### Phase 8: News & Whale Services (Days 15-16)

**Tasks:**
1. ☐ Implement News service
2. ☐ Implement Whale service
3. ☐ RSS feed integration
4. ☐ Whale Alert API integration
5. ☐ Deploy both services
6. ☐ Integration testing

**Deliverables:**
- News service deployed
- Whale service deployed
- RSS feeds working
- Whale alerts functional

**Success Criteria:**
- News fetching and publishing working
- Whale transaction alerts working
- Both services stable on Railway

---

### Phase 9: Testing & Optimization (Days 17-19)

**Tasks:**
1. ☐ End-to-end testing
2. ☐ Load testing
3. ☐ Performance optimization
4. ☐ Security audit
5. ☐ Documentation updates
6. ☐ Monitoring setup

**Deliverables:**
- Complete test suite
- Performance benchmarks
- Security audit report
- Updated documentation
- Monitoring dashboards

**Success Criteria:**
- All tests passing
- Performance meets SLAs
- No security vulnerabilities
- Documentation complete

---

### Phase 10: Production Launch (Days 20-21)

**Tasks:**
1. ☐ Production deployment
2. ☐ DNS configuration
3. ☐ Monitoring verification
4. ☐ Backup verification
5. ☐ User announcement
6. ☐ Post-launch monitoring

**Deliverables:**
- Production environment live
- All services healthy
- Monitoring active
- Backups configured

**Success Criteria:**
- Zero downtime migration
- All features working
- Monitoring showing green
- Users experiencing no issues

---

## Testing Strategy

### Unit Testing
- Test coverage target: >80%
- Test each service independently
- Mock external dependencies
- Use table-driven tests

### Integration Testing
- Test service-to-service communication
- Test database operations
- Test gRPC endpoints
- Test Discord bot commands

### End-to-End Testing
- Test complete user workflows
- Test failure scenarios
- Test rate limiting
- Test concurrent operations

### Load Testing
- Simulate high message volume
- Test queue management under load
- Test database connection pooling
- Identify bottlenecks

### Tools
- Go testing package
- Testify for assertions
- gomock for mocking
- k6 for load testing

---

## Rollback Plan

### Rollback Triggers
1. Critical bugs affecting core features
2. Data corruption or loss
3. Performance degradation >50%
4. Service unavailability >5 minutes
5. Security vulnerabilities

### Rollback Procedure

**Step 1: Immediate Actions**
1. Stop new deployments
2. Revert to previous Railway deployment
3. Switch DNS back to monolithic service
4. Notify team and users

**Step 2: Data Recovery**
1. Restore database from backup
2. Verify data integrity
3. Re-sync if needed

**Step 3: Analysis**
1. Identify root cause
2. Document issues
3. Plan corrective actions
4. Update rollback procedures

**Step 4: Communication**
1. Status page update
2. User notification
3. Incident report
4. Lessons learned document

---

## Success Metrics

### Performance Metrics
- **Response Time:** <500ms for 95th percentile
- **Availability:** >99.9% uptime
- **Error Rate:** <0.1% of requests
- **Database Queries:** <100ms for 95th percentile

### Operational Metrics
- **Deployment Frequency:** Multiple times per day
- **Mean Time to Recovery:** <30 minutes
- **Change Failure Rate:** <5%
- **Lead Time for Changes:** <1 hour

### Business Metrics
- **User Satisfaction:** >4.5/5 rating
- **Feature Adoption:** >80% of users use new features
- **Support Tickets:** <10 per week
- **Community Growth:** Positive trend

---

## Risk Assessment

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Service communication latency | High | Medium | Implement caching, optimize gRPC |
| Database migration data loss | Low | High | Dual-write period, thorough testing |
| Railway service limits | Medium | Medium | Monitor usage, plan scaling |
| Discord API rate limits | Medium | High | Implement rate limiting, request queuing |
| FFmpeg/yt-dlp compatibility | Low | Medium | Container with pre-installed tools |

### Operational Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Increased operational complexity | High | Medium | Comprehensive documentation, monitoring |
| Cost overrun | Medium | Medium | Regular cost review, optimization |
| Team knowledge gap | Low | Medium | Training, documentation |
| Vendor lock-in (Railway) | Low | Low | Use standard containers, portable config |

---

## Maintenance Plan

### Regular Maintenance
- **Daily:** Monitor service health, check logs
- **Weekly:** Review metrics, update dependencies
- **Monthly:** Security patches, performance review
- **Quarterly:** Architecture review, cost optimization

### Backup Strategy
- **Database:** Daily automated backups, 30-day retention
- **Configuration:** Version controlled in Git
- **Logs:** 30-day retention on Railway
- **Docker Images:** Tagged and stored in registry

### Monitoring Checklist
- ✅ All services responding to health checks
- ✅ Database connections healthy
- ✅ Redis cache hit rate >80%
- ✅ Discord bot latency <500ms
- ✅ Error rate <0.1%
- ✅ CPU usage <70%
- ✅ Memory usage <80%
- ✅ Disk usage <70%

---

## Conclusion

This migration plan provides a comprehensive roadmap for transforming NeruBot from a monolithic application to a modern microservices architecture deployed on Railway. By following the phased approach and adhering to the testing and rollback strategies, we can ensure a smooth transition with minimal disruption to users.

### Key Success Factors
1. Thorough planning and documentation
2. Incremental migration approach
3. Comprehensive testing at each phase
4. Robust monitoring and alerting
5. Clear rollback procedures
6. Regular communication with stakeholders

### Next Steps
1. Review and approve this plan
2. Setup Railway project and infrastructure
3. Begin Phase 1: Project Setup
4. Execute phases incrementally
5. Monitor progress and adjust as needed

---

**Document Version:** 1.0.0  
**Last Updated:** December 5, 2025  
**Owner:** @nerufuyo  
**Status:** Ready for Implementation
