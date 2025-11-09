# Contributing to NeruBot

Thank you for your interest in contributing to NeruBot! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment, discrimination, or trolling
- Publishing others' private information
- Other conduct which could reasonably be considered inappropriate

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.21 or higher
- Git
- FFmpeg
- yt-dlp
- A code editor (VS Code, GoLand, etc.)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/nerubot.git
   cd nerubot
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/nerufuyo/nerubot.git
   ```

## Development Setup

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools (optional)
go install golang.org/x/tools/cmd/goimports@latest
```

### Configure Environment

```bash
cp .env.example .env
# Edit .env with your Discord bot token and other credentials
```

### Build and Run

```bash
# Build the project
make build

# Run the bot
./build/nerubot

# Or run directly
go run cmd/nerubot/main.go
```

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates.

**When submitting a bug report, include:**

- Clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Go version (`go version`)
- Operating system
- Relevant logs or error messages

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:

- Clear and descriptive title
- Detailed description of the proposed feature
- Why this enhancement would be useful
- Possible implementation approach (optional)

### Your First Code Contribution

Not sure where to start? Look for issues labeled:

- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `documentation` - Improvements to docs

## Coding Standards

### Go Code Style

Follow standard Go conventions:

```bash
# Format your code
go fmt ./...

# Check for common mistakes
go vet ./...

# Run imports formatter
goimports -w .
```

### Architecture Guidelines

NeruBot follows Clean Architecture. Please read [docs/format-architecture.md](docs/format-architecture.md) for detailed guidelines.

**Key Principles:**

1. **KISS** - Keep It Simple, Stupid
2. **Separation of Concerns** - Each layer has a single responsibility
3. **Dependency Rule** - Dependencies point inward
4. **Interface Segregation** - Use small, focused interfaces

### File Organization

```
internal/
├── config/          # Configuration only
├── entity/          # Domain models only
├── repository/      # Data access only
├── usecase/         # Business logic only
├── delivery/        # External interfaces only
└── pkg/             # Shared utilities only
```

### Code Quality

- **Keep functions small** - Ideally under 50 lines
- **Single Responsibility** - One function does one thing
- **Error Handling** - Always handle errors, never ignore them
- **Comments** - Add comments for complex logic
- **Naming** - Use clear, descriptive names

### Example Code Style

```go
// Good
func (s *MusicService) AddToQueue(song *entity.Song) error {
    if song == nil {
        return ErrInvalidSong
    }
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.queue = append(s.queue, song)
    return nil
}

// Bad
func (s *MusicService) add(x *entity.Song) error {
    s.mu.Lock()
    s.queue = append(s.queue, x)
    s.mu.Unlock()
    return nil
}
```

## Commit Guidelines

We follow conventional commits. See [docs/format-commit.md](docs/format-commit.md) for details.

### Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

### Types

- `feat` - A new feature
- `fix` - A bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting, etc.)
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks

### Examples

```bash
feat: add Spotify playlist support to music system

Implemented Spotify API integration for playlist fetching.
Added fallback to YouTube search for unavailable tracks.

Closes #123

---

fix: resolve memory leak in session cleanup

The cleanup goroutine wasn't properly closing channels,
causing goroutines to leak over time.

---

docs: update README with new deployment instructions

Added Docker Compose section and updated environment
variable documentation.
```

## Pull Request Process

### Before Submitting

1. **Update your fork:**
   ```bash
   git fetch upstream
   git rebase upstream/master
   ```

2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes:**
   - Write clean, documented code
   - Follow coding standards
   - Add tests if applicable

4. **Test your changes:**
   ```bash
   go test ./...
   go build ./...
   ```

5. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add your feature"
   ```

6. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

### Submitting the Pull Request

1. Go to the [NeruBot repository](https://github.com/nerufuyo/nerubot)
2. Click "New Pull Request"
3. Select your fork and branch
4. Fill in the PR template with:
   - **Description** - What changes did you make?
   - **Motivation** - Why are these changes needed?
   - **Testing** - How did you test the changes?
   - **Screenshots** - If applicable

### PR Review Process

- Maintainers will review your PR
- Address any feedback or requested changes
- Once approved, your PR will be merged
- Your contribution will be credited in the changelog

### PR Checklist

- [ ] Code follows the project's style guidelines
- [ ] Self-review of the code completed
- [ ] Comments added for complex logic
- [ ] Documentation updated (if needed)
- [ ] No new warnings generated
- [ ] Tests added/updated (if applicable)
- [ ] All tests passing
- [ ] Commit messages follow conventions

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/usecase/music/

# Run tests with verbose output
go test -v ./...
```

### Writing Tests

```go
func TestMusicService_AddToQueue(t *testing.T) {
    service := NewMusicService()
    song := &entity.Song{
        Title: "Test Song",
        URL:   "https://youtube.com/watch?v=test",
    }
    
    err := service.AddToQueue(song)
    if err != nil {
        t.Errorf("AddToQueue() error = %v", err)
    }
    
    if len(service.queue) != 1 {
        t.Errorf("Queue length = %d, want 1", len(service.queue))
    }
}
```

### Test Coverage

Aim for:
- **Overall:** >70% coverage
- **Critical paths:** >90% coverage
- **New features:** >80% coverage

## Documentation

### Code Documentation

- **Exported functions** - Must have doc comments
- **Packages** - Must have package doc comments
- **Complex logic** - Should have inline comments

```go
// Package music provides music playback functionality for the Discord bot.
package music

// MusicService handles music playback, queue management, and voice connections.
type MusicService struct {
    queue []*entity.Song
    mu    sync.RWMutex
}

// AddToQueue adds a song to the playback queue.
// Returns an error if the song is nil or invalid.
func (s *MusicService) AddToQueue(song *entity.Song) error {
    // Implementation
}
```

### README Updates

If your contribution affects usage or setup:

- Update the README.md
- Update relevant documentation in docs/
- Add examples if applicable

### Changelog

Maintainers will update the CHANGELOG.md when releasing new versions.

## Questions?

If you have questions, feel free to:

- Open an issue with the `question` label
- Ask in discussions (if enabled)
- Contact the maintainers

## License

By contributing to NeruBot, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to NeruBot! 🎉
