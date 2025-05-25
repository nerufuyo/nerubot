# NeruBot - Simple Discord Bot

A clean, modular Discord bot built with Python and discord.py.

## ✨ Features

- 🎵 **Music**: Play music from YouTube with queue support
- 🎲 **Fun**: Dice rolling, coin flipping, jokes, and magic 8-ball
- 🔧 **Utility**: Calculator, user info, server info, avatars
- 📋 **General**: Ping, bot info, help commands

## 🚀 Quick Start

1. **Clone and Setup**
   ```bash
   git clone <your-repo>
   cd nerubot
   pip install -r requirements_new.txt
   ```

2. **Configure Bot**
   Create a `.env` file:
   ```
   DISCORD_TOKEN=your_discord_bot_token_here
   ```

3. **Run Bot**
   ```bash
   python3 bot.py
   ```

## 🏗️ Architecture

The bot uses a clean, modular architecture:

```
nerubot/
├── bot.py              # Main entry point
├── cogs/               # Feature modules
│   ├── general.py      # Basic commands
│   ├── music.py        # Music functionality  
│   ├── fun.py          # Entertainment
│   └── utility.py      # Tools & utilities
├── services/           # Business logic
│   └── music_service.py
├── models/             # Data models
│   └── song.py
└── config/             # Configuration
    └── settings.py
```

## 🎵 Music Commands

- `/play <song>` - Play a song from YouTube
- `/stop` - Stop music and clear queue
- `/skip` - Skip current song
- `/queue` - Show current queue
- `/pause` - Pause playback
- `/resume` - Resume playback
- `/join` - Join voice channel
- `/leave` - Leave voice channel

## 🎲 Fun Commands

- `/roll <sides>` - Roll a dice
- `/coinflip` - Flip a coin
- `/joke` - Get a random joke
- `/8ball <question>` - Ask the magic 8-ball

## 🔧 Utility Commands

- `/calculate <expression>` - Basic calculator
- `/userinfo [user]` - Get user information
- `/serverinfo` - Get server information
- `/avatar [user]` - Get user's avatar

## 📋 General Commands

- `/ping` - Check bot latency
- `/info` - Bot information
- `/help` - Show all commands

## 🛠️ Adding New Features

Adding new features is simple with the modular architecture:

1. **Create a new cog** in `cogs/` directory:
   ```python
   # cogs/example.py
   import discord
   from discord.ext import commands
   from discord import app_commands

   class Example(commands.Cog):
       def __init__(self, bot):
           self.bot = bot
       
       @app_commands.command(name="example", description="Example command")
       async def example(self, interaction: discord.Interaction):
           await interaction.response.send_message("Hello!")

   async def setup(bot):
       await bot.add_cog(Example(bot))
   ```

2. **Add business logic** in `services/` if needed
3. **Add models** in `models/` for data structures
4. The bot automatically loads all cogs!

## 📦 Dependencies

- `discord.py` - Discord API wrapper
- `yt-dlp` - YouTube audio extraction
- `python-dotenv` - Environment variable management
- `aiohttp` - Async HTTP requests
- `psutil` - System information

## 🔧 Configuration

Environment variables in `.env`:

```bash
# Required
DISCORD_TOKEN=your_bot_token

# Optional
COMMAND_PREFIX=!
MAX_QUEUE_SIZE=50
DEFAULT_VOLUME=0.5
LOG_LEVEL=INFO
WEATHER_API_KEY=your_weather_api_key
```

## 📝 License

This project is licensed under the MIT License.
