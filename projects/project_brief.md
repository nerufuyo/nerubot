# NeruBot Project Brief

## Project Overview

**Project Name:** NeruBot - Discord Companion Bot  
**Version:** 3.0.0 (Golang Edition)  
**Created by:** [@nerufuyo](https://github.com/nerufuyo)  
**Current Status:** Monolithic Architecture → Migrating to Microservices  
**Target Deployment:** Railway Platform

## Executive Summary

NeruBot is a feature-rich Discord bot built with Go, following Clean Architecture principles. The project aims to migrate from a monolithic architecture to a microservices-based architecture while maintaining code quality, scalability, and reliability. The bot provides multiple features including music streaming, anonymous confessions, user roasting, AI chatbot, news aggregation, and cryptocurrency whale alerts.

## Current Architecture

### Technology Stack
- **Language:** Go 1.25.1
- **Framework:** DiscordGo v0.29.0
- **Architecture:** Clean Architecture (Monolithic)
- **Storage:** JSON-based file storage
- **External Tools:** FFmpeg, yt-dlp
- **Deployment:** Docker, Docker Compose

### Project Structure
```
nerubot/
├── cmd/nerubot/           # Application entry point
├── internal/
│   ├── config/            # Configuration layer
│   ├── delivery/discord/  # Discord bot interface
│   ├── entity/            # Domain models
│   ├── pkg/              # Shared packages (AI, FFmpeg, Logger, yt-dlp)
│   ├── repository/        # Data persistence layer
│   └── usecase/          # Business logic layer
│       ├── chatbot/      # AI chatbot service
│       ├── confession/   # Anonymous confession system
│       ├── music/        # Music streaming service
│       ├── news/         # RSS news aggregation
│       ├── roast/        # User roasting system
│       └── whale/        # Crypto whale alerts
├── data/                 # JSON data storage
├── deploy/               # Deployment scripts and configs
└── docs/                 # Documentation
```

## Core Features

### 1. 🎵 Music System
- **Description:** YouTube audio streaming with queue management
- **Key Capabilities:**
  - High-quality audio streaming via yt-dlp
  - Queue management (add, skip, stop, shuffle)
  - Loop modes (none, single, queue)
  - Voice state detection
  - Thread-safe operations
- **Dependencies:** FFmpeg, yt-dlp
- **Status:** ✅ Implemented

### 2. 🤐 Confession System
- **Description:** Anonymous confession submission and management
- **Key Capabilities:**
  - Complete anonymity
  - Image attachment support
  - Reply system for engagement
  - Moderation queue with approval/rejection
  - Per-server settings
  - Thread-safe JSON storage
- **Dependencies:** None
- **Status:** ✅ Implemented

### 3. 🔥 Roast System
- **Description:** AI-powered user roasting based on activity tracking
- **Key Capabilities:**
  - Activity tracking (messages, reactions, voice time)
  - 8 roast categories (spammer, lurker, etc.)
  - Profile analysis and statistics
  - Safety systems (cooldowns, friendly content)
  - Persistent data storage
- **Dependencies:** None
- **Status:** ✅ Implemented

### 4. 🤖 AI Chatbot (Coming Soon)
- **Description:** Multi-provider AI chatbot integration
- **Key Capabilities:**
  - Multi-provider support (Claude, Gemini, OpenAI)
  - Automatic fallback between providers
  - Session management (30-min timeout)
  - Context-aware conversations
- **Dependencies:** AI API keys
- **Status:** 🚧 Planned

### 5. 📰 News System (Coming Soon)
- **Description:** RSS feed aggregation and publishing
- **Key Capabilities:**
  - Multiple source aggregation
  - Concurrent news fetching
  - Customizable sources
  - Auto-publishing capability
- **Dependencies:** None
- **Status:** 🚧 Planned

### 6. 🐋 Whale Alerts (Coming Soon)
- **Description:** Cryptocurrency whale transaction monitoring
- **Key Capabilities:**
  - Real-time transaction alerts
  - Configurable minimum threshold
  - Multi-blockchain support
  - Transaction tracking
- **Dependencies:** Crypto API
- **Status:** 🚧 Planned

## Migration Goals

### Primary Objectives

1. **Microservices Architecture**
   - Separate each feature into independent microservices
   - Implement API Gateway for Discord interaction
   - Enable independent scaling and deployment
   - Improve fault isolation

2. **Railway Deployment**
   - Configure Railway deployment for each microservice
   - Setup environment variables and secrets
   - Configure automatic deployments from Git
   - Setup monitoring and logging

3. **Service Communication**
   - Implement gRPC for inter-service communication
   - Setup service discovery mechanism
   - Implement circuit breakers and retries
   - Add distributed tracing

4. **Data Management**
   - Migrate from file-based to database storage
   - Setup PostgreSQL for persistent data
   - Implement Redis for caching and sessions
   - Design database schema per service

### Target Microservices

1. **API Gateway Service** (Port: 8080)
   - Discord bot interface
   - Command routing
   - Response aggregation

2. **Music Service** (Port: 8081)
   - Audio streaming
   - Queue management
   - FFmpeg processing

3. **Confession Service** (Port: 8082)
   - Confession management
   - Moderation queue
   - Reply system

4. **Roast Service** (Port: 8083)
   - Activity tracking
   - Roast generation
   - Statistics

5. **Chatbot Service** (Port: 8084)
   - AI provider integration
   - Session management
   - Context handling

6. **News Service** (Port: 8085)
   - RSS aggregation
   - Article fetching
   - Publishing

7. **Whale Service** (Port: 8086)
   - Transaction monitoring
   - Alert generation
   - Threshold management

## Technical Requirements

### Development Environment
- Go 1.21+
- Docker & Docker Compose
- FFmpeg
- yt-dlp
- Git

### Production Infrastructure
- Railway Platform
- PostgreSQL (managed)
- Redis (managed)
- Object Storage (for media files)

### External Services
- Discord API
- OpenAI API (optional)
- Anthropic API (optional)
- Gemini API (optional)
- Crypto Alert API (optional)

## Success Criteria

1. ✅ Each feature runs as an independent microservice
2. ✅ Successfully deployed to Railway platform
3. ✅ All existing features continue to work
4. ✅ Improved scalability and fault tolerance
5. ✅ Comprehensive monitoring and logging
6. ✅ Database migration from JSON to PostgreSQL
7. ✅ Automated CI/CD pipeline
8. ✅ Complete documentation for each service

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Service communication latency | High | Implement caching, optimize gRPC calls |
| Data migration complexity | High | Phased migration, comprehensive testing |
| Increased operational complexity | Medium | Comprehensive monitoring, documentation |
| Cost increase from multiple services | Medium | Optimize resource allocation, auto-scaling |
| Discord API rate limits | Medium | Implement rate limiting, request queuing |

## Timeline Estimate

- **Phase 1 - Planning & Setup:** 2-3 days
- **Phase 2 - Service Separation:** 5-7 days
- **Phase 3 - Database Migration:** 3-4 days
- **Phase 4 - Railway Deployment:** 2-3 days
- **Phase 5 - Testing & Optimization:** 3-4 days
- **Total Estimated Time:** 15-21 days

## Project Principles

1. **Clean Architecture:** Maintain separation of concerns
2. **SOLID Principles:** Write maintainable, extensible code
3. **Testability:** Comprehensive unit and integration tests
4. **Documentation:** Clear, concise documentation for all components
5. **Security:** Secure API endpoints, encrypted secrets
6. **Performance:** Optimize for low latency and high throughput
7. **Reliability:** Implement proper error handling and recovery

## Team & Ownership

- **Lead Developer:** @nerufuyo
- **Repository:** github.com/nerufuyo/nerubot
- **License:** MIT
- **Support:** GitHub Issues

## References

- [Architecture Documentation](../ARCHITECTURE.md)
- [Deployment Guide](../docs/DEPLOYMENT.md)
- [Contributing Guidelines](../CONTRIBUTING.md)
- [Changelog](../CHANGELOG.md)

---

**Last Updated:** December 5, 2025  
**Document Version:** 1.0.0
