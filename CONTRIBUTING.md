# Contributing to NeruBot

Thank you for your interest in contributing to NeruBot! This document provides guidelines and information for contributors.

## 🚀 Quick Start for Contributors

### Prerequisites
- Python 3.8 or higher
- Git
- Basic knowledge of Discord.py
- Familiarity with async/await patterns

### Setup Development Environment
```bash
# Fork and clone the repository
git clone https://github.com/yourusername/nerubot.git
cd nerubot

# Set up development environment
./start.sh setup

# Activate virtual environment
source nerubot_env/bin/activate

# Install development dependencies
pip install -r requirements-dev.txt
```

## 📋 Code of Conduct

- Be respectful and inclusive
- Follow Python best practices
- Write clean, documented code
- Test your changes thoroughly
- Follow the project's architecture patterns

## 🏗️ Project Architecture

NeruBot follows a modular, feature-based architecture:

```
src/
├── main.py                 # Application entry point
├── config/                 # Configuration management
├── core/                   # Core utilities and shared code
├── features/               # Feature modules (music, help, etc.)
│   └── music/
│       ├── cogs/          # Discord commands
│       ├── services/      # Business logic
│       └── models/        # Data models
└── interfaces/            # External interfaces (Discord, etc.)
```

### Key Principles
- **Single Responsibility**: Each module has one clear purpose
- **Dependency Injection**: Services are injected, not imported directly
- **Interface Segregation**: Small, focused interfaces
- **DRY (Don't Repeat Yourself)**: Shared utilities in `core/`
- **KISS (Keep It Simple)**: Simple, readable code

## 🔧 Development Guidelines

### Code Style
- Follow PEP 8 guidelines
- Use Black for code formatting: `black src/`
- Use meaningful variable and function names
- Add type hints to all functions
- Maximum line length: 88 characters

### Documentation
- Add docstrings to all public functions and classes
- Use Google-style docstrings
- Update README.md for significant changes
- Document complex algorithms and business logic

### Testing
- Write unit tests for new features
- Test both success and error cases
- Use pytest for testing framework
- Aim for >80% code coverage

### Git Workflow
1. Create a feature branch: `git checkout -b feature/amazing-feature`
2. Make your changes following the guidelines
3. Test your changes thoroughly
4. Commit with clear, descriptive messages
5. Push to your fork and create a Pull Request

### Commit Message Format
```
type(scope): description

Longer explanation if needed

Fixes #123
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## 🎵 Adding New Music Features

### Creating a New Music Command
1. Add the command to `src/features/music/cogs/`
2. Implement business logic in `src/features/music/services/`
3. Add any data models to `src/features/music/models/`
4. Update tests and documentation

Example structure:
```python
# src/features/music/cogs/playlist_cog.py
from discord.ext import commands
from ..services.playlist_service import PlaylistService

class PlaylistCog(commands.Cog):
    def __init__(self, bot, playlist_service: PlaylistService):
        self.bot = bot
        self.playlist_service = playlist_service
    
    @commands.slash_command(name="playlist")
    async def playlist_command(self, ctx, action: str):
        """Manage playlists"""
        # Implementation here
```

### Adding New Music Sources
1. Create a new extractor in `src/features/music/services/extractors/`
2. Implement the `MusicExtractor` interface
3. Register the extractor in the music service
4. Add appropriate tests

## 🤖 Adding New Features

### Feature Module Structure
```
src/features/new_feature/
├── __init__.py
├── cogs/                   # Discord commands
├── services/               # Business logic
├── models/                 # Data models
└── tests/                  # Feature tests
```

### Integration Steps
1. Create the feature module following the structure above
2. Add the feature to `src/main.py`
3. Update configuration if needed
4. Add documentation and tests

## 🧪 Testing

### Running Tests
```bash
# Run all tests
python -m pytest

# Run specific feature tests
python -m pytest src/features/music/tests/

# Run with coverage
python -m pytest --cov=src/

# Run integration tests
python -m pytest tests/integration/
```

### Writing Tests
```python
import pytest
from unittest.mock import Mock, AsyncMock
from src.features.music.services.music_service import MusicService

class TestMusicService:
    @pytest.fixture
    def music_service(self):
        return MusicService()
    
    @pytest.mark.asyncio
    async def test_play_song(self, music_service):
        # Test implementation
        pass
```

## 📝 Documentation

### Code Documentation
- Use docstrings for all public methods
- Include parameter and return type information
- Provide usage examples for complex functions

```python
async def play_song(self, query: str, voice_channel) -> PlayResult:
    """Play a song from various sources.
    
    Args:
        query: Search query or direct URL
        voice_channel: Discord voice channel to join
        
    Returns:
        PlayResult containing song info and status
        
    Raises:
        MusicError: If song cannot be played
        
    Example:
        result = await music_service.play_song("Bohemian Rhapsody", channel)
    """
```

### README Updates
- Update feature lists for new functionality
- Add new configuration options
- Update installation instructions if needed

## 🐛 Bug Reports

### Before Reporting
1. Check existing issues
2. Test with the latest version
3. Gather relevant information

### Bug Report Template
```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Run command '...'
2. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g. Ubuntu 20.04]
- Python version: [e.g. 3.9.5]
- Bot version: [e.g. 2.0.1]

**Additional context**
Any other context about the problem.
```

## 💡 Feature Requests

### Feature Request Template
```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Alternative solutions or features you've considered.

**Additional context**
Any other context or screenshots about the feature request.
```

## 🔧 Development Tools

### Recommended IDE Setup
- **VS Code** with Python extension
- **PyCharm** Professional or Community
- Configure linting and formatting tools

### Useful Commands
```bash
# Code formatting
black src/
isort src/

# Linting
flake8 src/
pylint src/

# Type checking
mypy src/

# Security scanning
bandit src/

# Dependency checking
pip-audit
```

## 🚀 Release Process

### Version Numbering
We use Semantic Versioning (SemVer):
- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes

### Release Checklist
1. Update version numbers
2. Update CHANGELOG.md
3. Run full test suite
4. Update documentation
5. Create release tag
6. Deploy to production

## 📞 Getting Help

### Community
- **GitHub Discussions**: For questions and ideas
- **Issues**: For bug reports and feature requests
- **Discord**: Join our development server (link in README)

### Documentation
- **Project Wiki**: Detailed technical documentation
- **API Reference**: Auto-generated from docstrings
- **Architecture Guide**: Deep dive into project structure

## 🏆 Recognition

Contributors will be:
- Listed in the AUTHORS file
- Mentioned in release notes for significant contributions
- Invited to join the core team for sustained contributions

Thank you for contributing to NeruBot! 🎵
