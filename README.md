# NeruBot - Enhanced Discord Bot

A sophisticated, modular Discord bot with music capabilities and advanced features including news, AI quotes, user profiles, and anonymous confessions.

## ✨ Features

### 🎵 Music (Ready)
- Play music from YouTube with queue support
- Advanced queue management and controls
- High-quality audio streaming

### 📰 News (Ready) ⭐ NEW!
- Real-time RSS news feeds from multiple sources
- Category-based news filtering (Technology, World, General)
- Support for major news sources (BBC, CNN, TechCrunch, Reuters, etc.)
- Clean, formatted news display with links

### 💭 AI Quotes (Coming Soon)
- AI-powered inspirational quotes using DeepSeek
- Category-based quote generation
- Mood-based suggestions
- Multi-language support

### 👤 User Profiles (Coming Soon)
- Custom user profiles with stats tracking
- Activity monitoring and achievements
- Preference management and social features

### 🤐 Anonymous Confessions (Coming Soon)
- Secure anonymous messaging system
- Optional moderation and content filtering
- Server-specific configuration

## 🚀 Quick Start

1. **Install Dependencies**
   ```bash
   pip install -r requirements.txt
   ```

2. **Configure Environment**
   ```bash
   # Create .env file
   DISCORD_TOKEN=your_discord_bot_token_here
   ```

3. **Run the Bot**
   ```bash
   python3 -m src.main
   ```

## 📁 New Improved Architecture

The bot now uses a feature-based modular architecture following DRY and KISS principles:

```
src/
├── main.py                    # Entry point
├── interfaces/               # Discord interface layer
│   └── discord/
│       ├── bot.py           # Main bot class
│       └── music_cog.py     # Music commands
├── features/                # Feature modules (NEW!)
│   ├── news/               # News feature
│   │   ├── cogs/           # Discord commands
│   │   ├── services/       # Business logic
│   │   └── models/         # Data models
│   ├── quotes/             # AI quotes feature
│   ├── profile/            # User profiles
│   └── confession/         # Anonymous confessions
├── shared/                 # Shared utilities
│   ├── services/           # Common services
│   ├── models/             # Shared models
│   └── utils/              # Utility functions
└── application/            # Core application logic
    └── services/
        └── music_service.py
```

## 🆕 News Commands

### `/news` - Get Latest News
Get the latest news articles from multiple sources.

**Usage:**
- `/news` - General news
- `/news category:technology` - Tech news
- `/news category:world count:3` - World news (3 articles)

### `/news-categories` - Available Categories
View all available news categories and sources.

### `/news-source` - Source-Specific News
Get news from a specific source.

**Usage:**
- `/news-source source:BBC World`
- `/news-source source:TechCrunch count:5`

**Available Sources:**
- BBC World
- CNN Top Stories
- TechCrunch
- Hacker News
- Reuters World
- AP News

## 🎵 Music Commands

- `/play <song>` - Play a song from YouTube
- `/queue` - Show current queue
- `/skip` - Skip current song
- `/stop` - Stop music and clear queue
- `/pause` / `/resume` - Control playback
- `/join` / `/leave` - Voice channel management

## 🔮 Upcoming Features

### `/quote` - AI Quotes (Coming Soon)
Get AI-generated inspirational quotes using DeepSeek AI.
- Category-based quotes (wisdom, motivation, etc.)
- Mood-aware suggestions
- Personalized content

### `/profile` - User Profiles (Coming Soon)
Comprehensive user profile system.
- Custom bios and preferences
- Activity statistics and achievements
- Social features and reputation

### `/confess` - Anonymous Confessions (Coming Soon)
Secure anonymous messaging system.
- Complete anonymity and privacy
- Optional moderation system
- Server-specific configuration

## 🛠️ Development

### Adding New Features

The modular architecture makes adding features simple:

1. **Create feature directory:**
   ```
   src/features/my_feature/
   ├── cogs/          # Discord commands
   ├── services/      # Business logic
   └── models/        # Data structures
   ```

2. **Implement the cog:**
   ```python
   # src/features/my_feature/cogs/my_cog.py
   from discord.ext import commands
   
   class MyCog(commands.Cog):
       @app_commands.command()
       async def my_command(self, interaction):
           # Your command logic
   ```

3. **Add to bot:** The bot automatically loads feature cogs!

### Architecture Benefits

- **🧩 Modular:** Each feature is independent
- **🔧 Maintainable:** Clear separation of concerns
- **📈 Scalable:** Easy to add/remove features
- **🧪 Testable:** Services can be tested independently
- **♻️ DRY:** Shared utilities prevent code duplication
- **💋 KISS:** Simple, clean interfaces

## 🔧 Configuration

Environment variables in `.env`:

```bash
# Required
DISCORD_TOKEN=your_bot_token

# Optional
COMMAND_PREFIX=!
LOG_LEVEL=INFO

# Future: API keys for additional features
DEEPSEEK_API_KEY=your_deepseek_key  # For AI quotes
```

## 📊 Testing

Test the news feature:
```bash
python3 test_news.py
```

## 🤝 Contributing

1. Follow the feature-based architecture
2. Use the shared utilities for common functionality
3. Write tests for new services
4. Follow Python best practices and type hints

## 📝 License

MIT License - see LICENSE file for details.

---

**NeruBot v2.0** - Now with advanced features and improved architecture! 🚀
