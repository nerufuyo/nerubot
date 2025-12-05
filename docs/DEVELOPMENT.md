# NeruBot Microservices Development Guide

## Table of Contents
1. [Getting Started](#getting-started)
2. [Project Structure](#project-structure)
3. [Development Environment](#development-environment)
4. [Building and Running](#building-and-running)
5. [Testing](#testing)
6. [Contributing](#contributing)
7. [Troubleshooting](#troubleshooting)

---

## Getting Started

### Prerequisites

**Required:**
- Go 1.21 or higher
- Docker and Docker Compose
- Git
- Make (for build automation)
- Protocol Buffers compiler (protoc)

**Optional:**
- PostgreSQL client (for database access)
- Redis CLI (for cache debugging)
- Postman or similar (for API testing)

### Installation

#### 1. Install Go
```bash
# Download from https://golang.org/dl/
# Or use package manager:

# macOS
brew install go

# Ubuntu/Debian
sudo apt-get install golang-go

# Windows
# Download installer from golang.org
```

#### 2. Install Protocol Buffers
```bash
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# Windows
# Download from https://github.com/protocolbuffers/protobuf/releases
# Add to PATH
```

#### 3. Install Go Tools
```bash
# Install protoc Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add to PATH (if not already)
export PATH="$PATH:$(go env GOPATH)/bin"
```

#### 4. Clone Repository
```bash
git clone https://github.com/nerufuyo/nerubot.git
cd nerubot
```

#### 5. Setup Development Environment
```bash
# Install dependencies
make -f Makefile.microservices deps

# Generate gRPC code
make -f Makefile.microservices proto

# Setup environment variables
cp .env.example .env
# Edit .env with your configuration
```

---

## Project Structure

```
nerubot/
├── api/
│   └── proto/               # Protocol Buffer definitions
│       ├── music.proto
│       ├── confession.proto
│       ├── roast.proto
│       ├── chatbot.proto
│       ├── news.proto
│       └── whale.proto
├── cmd/
│   └── nerubot/            # Monolithic application (legacy)
│       └── main.go
├── services/               # Microservices
│   ├── gateway/           # API Gateway
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── Dockerfile
│   │   └── railway.toml
│   ├── music/             # Music Service
│   │   ├── cmd/
│   │   ├── Dockerfile
│   │   └── railway.toml
│   ├── confession/        # Confession Service
│   ├── roast/            # Roast Service
│   ├── chatbot/          # Chatbot Service
│   ├── news/             # News Service
│   └── whale/            # Whale Service
├── internal/              # Shared internal packages
│   ├── config/           # Configuration
│   ├── entity/           # Domain models
│   ├── pkg/              # Shared utilities
│   │   ├── logger/
│   │   ├── ai/
│   │   ├── ffmpeg/
│   │   └── ytdlp/
│   ├── repository/       # Data access
│   └── usecase/          # Business logic
├── data/                 # JSON data storage (legacy)
├── scripts/              # Utility scripts
│   └── init-db.sql      # Database initialization
├── docs/                 # Documentation
│   ├── DEPLOYMENT.md
│   └── RAILWAY_DEPLOYMENT.md
├── projects/             # Project planning
│   ├── project_brief.md
│   └── project_plan.md
├── docker-compose.yml           # Legacy monolith
├── docker-compose.microservices.yml  # Microservices
├── Makefile                     # Monolith build
├── Makefile.microservices       # Microservices build
├── railway.toml                 # Railway config
├── go.mod
└── go.sum
```

---

## Development Environment

### Environment Variables

Create a `.env` file in the project root:

```env
# Discord Configuration
DISCORD_TOKEN=your_discord_bot_token
DISCORD_GUILD_ID=your_guild_id

# Service URLs (for local development)
MUSIC_SERVICE_URL=localhost:8081
CONFESSION_SERVICE_URL=localhost:8082
ROAST_SERVICE_URL=localhost:8083
CHATBOT_SERVICE_URL=localhost:8084
NEWS_SERVICE_URL=localhost:8085
WHALE_SERVICE_URL=localhost:8086

# Database (Docker)
DATABASE_URL=postgresql://nerubot:nerubot_dev@localhost:5432

# Redis (Docker)
REDIS_URL=redis://localhost:6379

# Logging
LOG_LEVEL=INFO

# AI Providers (Optional)
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_anthropic_key
GEMINI_API_KEY=your_gemini_key

# External Services (Optional)
WHALE_ALERT_API_KEY=your_whale_alert_key
```

### Local Development Setup

#### Option 1: Docker Compose (Recommended)

Start all services with Docker:

```bash
# Start infrastructure (PostgreSQL, Redis)
docker-compose -f docker-compose.microservices.yml up -d postgres redis

# Wait for databases to be ready
sleep 10

# Initialize databases
docker-compose -f docker-compose.microservices.yml exec postgres \
  psql -U nerubot -f /docker-entrypoint-initdb.d/init.sql

# Start all services
docker-compose -f docker-compose.microservices.yml up -d

# View logs
docker-compose -f docker-compose.microservices.yml logs -f
```

#### Option 2: Manual (Native)

Run services individually for debugging:

```bash
# Terminal 1: Start PostgreSQL and Redis
docker-compose -f docker-compose.microservices.yml up postgres redis

# Terminal 2: Build and run API Gateway
make -f Makefile.microservices build-gateway
./build/gateway/gateway

# Terminal 3: Run other services as needed
# ... (similar for other services)
```

---

## Building and Running

### Build Commands

```bash
# Build all services
make -f Makefile.microservices build

# Build specific service
make -f Makefile.microservices build-gateway

# Clean build artifacts
make -f Makefile.microservices clean

# Run tests
make -f Makefile.microservices test

# Format code
make -f Makefile.microservices fmt

# Tidy dependencies
make -f Makefile.microservices tidy
```

### Running Services

#### Run API Gateway
```bash
# Build and run
make -f Makefile.microservices run-gateway

# Or manually
./build/gateway/gateway
```

#### Run with Hot Reload (Air)

Install Air for hot reloading during development:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with Air
air -c .air.toml
```

Create `.air.toml`:
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./services/gateway/cmd"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

---

## Testing

### Unit Tests

```bash
# Run all tests
make -f Makefile.microservices test

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# View coverage
go tool cover -html=coverage.txt

# Test specific package
go test -v ./services/gateway/...
```

### Integration Tests

```bash
# Start test environment
docker-compose -f docker-compose.microservices.yml up -d

# Run integration tests
go test -v -tags=integration ./tests/integration/...
```

### Manual Testing

#### Test API Gateway Health
```bash
curl http://localhost:8080/health
```

#### Test Discord Commands
1. Invite bot to your Discord server
2. Use slash commands: `/ping`, `/play`, `/confess`, etc.
3. Check logs for debugging

---

## Contributing

### Development Workflow

1. **Create a branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make changes**
- Follow Go best practices
- Write tests for new functionality
- Update documentation

3. **Test your changes**
```bash
make -f Makefile.microservices test
make -f Makefile.microservices fmt
```

4. **Commit changes**
Follow commit format from `docs/format-commit.md`:
```bash
git add .
git commit -m "feat: Add new feature"
```

5. **Push and create PR**
```bash
git push origin feature/your-feature-name
```

### Code Style Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions small and focused
- Write descriptive variable names
- Add unit tests for new code

### Commit Message Format

```
<type>: <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Add or update tests
- `chore`: Maintenance tasks

Examples:
```bash
feat: Add music queue management
fix: Resolve confession approval bug
docs: Update API documentation
refactor: Improve error handling in roast service
```

---

## Troubleshooting

### Common Issues

#### 1. Proto generation fails
```bash
# Error: protoc not found
# Solution: Install protobuf compiler

# macOS
brew install protobuf

# Ubuntu
sudo apt-get install protobuf-compiler

# Then install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### 2. Build fails with missing dependencies
```bash
# Solution: Download dependencies
make -f Makefile.microservices deps
go mod tidy
```

#### 3. Docker services won't start
```bash
# Check if ports are already in use
netstat -an | grep 8080
netstat -an | grep 5432
netstat -an | grep 6379

# Stop conflicting services or change ports
docker-compose -f docker-compose.microservices.yml down

# Restart
docker-compose -f docker-compose.microservices.yml up -d
```

#### 4. Database connection fails
```bash
# Check PostgreSQL is running
docker-compose -f docker-compose.microservices.yml ps

# Check connection string
echo $DATABASE_URL

# Verify databases exist
docker-compose -f docker-compose.microservices.yml exec postgres \
  psql -U nerubot -c "\l"
```

#### 5. Discord bot not responding
```bash
# Check bot token is correct
echo $DISCORD_TOKEN

# Check bot is online in Discord
# Check logs
docker-compose -f docker-compose.microservices.yml logs api-gateway

# Verify slash commands are registered
# May take up to 1 hour for Discord to update
```

### Debug Mode

Enable debug logging:
```env
LOG_LEVEL=DEBUG
```

View detailed logs:
```bash
# All services
docker-compose -f docker-compose.microservices.yml logs -f

# Specific service
docker-compose -f docker-compose.microservices.yml logs -f api-gateway
```

### Performance Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# View profile
go tool pprof cpu.prof
```

---

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [DiscordGo Documentation](https://github.com/bwmarrin/discordgo)
- [Railway Documentation](https://docs.railway.app/)
- [Project Brief](../projects/project_brief.md)
- [Migration Plan](../projects/project_plan.md)
- [Architecture](../ARCHITECTURE.md)

---

## Getting Help

- **GitHub Issues:** [Report bugs or request features](https://github.com/nerufuyo/nerubot/issues)
- **Discussions:** [Ask questions and share ideas](https://github.com/nerufuyo/nerubot/discussions)
- **Discord Server:** [Join our community](#) (Coming soon)

---

**Last Updated:** December 5, 2025  
**Version:** 1.0.0  
**Maintainer:** @nerufuyo
