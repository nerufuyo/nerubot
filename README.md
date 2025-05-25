# NeruBot - Discord Music Bot

A feature-rich Discord music bot built with Python using clean architecture principles.

## Features

- Play music from YouTube via URL or search query
- Queue management (add, remove, view queue)
- Playback controls (pause, resume, skip, stop)
- Volume control
- Loop modes (off, song, queue)
- Shuffle functionality
- Clean architecture design (DRY and KISS principles)
- Easy setup for both local and VPS deployment

## Architecture

This project follows clean architecture principles with distinct layers:

- **Core**: Contains business logic, entities, and use cases
- **Infrastructure**: Implements repositories and external services
- **Interfaces**: Handles user interaction via Discord commands
- **Application**: Contains services that coordinate between layers

## Requirements

- Python 3.8 or higher
- ffmpeg
- Discord Bot Token

## Quick Start

### One-Line Command

Simply run this command in your terminal:

```bash
./start_bot.sh
```

Or alternatively:

```bash
python3 run_bot.py
```

### What it does automatically:

The script will:
- Check for required dependencies and offer to install them
- Prompt for your Discord bot token if not found
- Check for ffmpeg installation and provide installation instructions
- Start the bot automatically

### VPS Deployment

To set up for VPS deployment:

```bash
python run_bot.py --setup-vps
```

The script will create a systemd service file and provide instructions for deploying to a VPS.

## Manual Setup

If you prefer manual setup:

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set up your Discord token in `.env`:
```
DISCORD_TOKEN=your_token_here
```

3. Run the bot:
```bash
python -m src.main
```

## Commands

- `!join` - Join your voice channel
- `!leave` - Leave the voice channel
- `!play <song>` - Play a song by URL or search term
- `!stop` - Stop playback and clear the queue
- `!pause` - Pause the current song
- `!resume` - Resume the paused song
- `!skip` - Skip to the next song
- `!volume <0-100>` - Set playback volume
- `!now` - Show the currently playing song
- `!queue [page]` - Show songs in the queue
- `!remove <index>` - Remove a song from the queue
- `!shuffle` - Shuffle the queue
- `!loop <off/song/queue>` - Set loop mode
- `!help` - Show all available commands

## Creating Your Own Discord Bot

To create your own Discord bot:

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to the "Bot" tab and click "Add Bot"
4. Under the "Token" section, click "Copy" to copy your bot token
5. Under "Privileged Gateway Intents", enable "Message Content Intent"
6. Go to OAuth2 > URL Generator:
   - Select "bot" and "applications.commands" scopes
   - Select permissions: "Send Messages", "Connect", "Speak", "Use Voice Activity"
7. Use the generated URL to invite the bot to your server

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
