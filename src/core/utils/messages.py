"""
Message constants for the Discord music bot
This file contains all user-facing messages to make them easier to update
"""

# Bot general messages
BOT_STARTED = "=== NeruBot Discord Music Bot ==="
BOT_SHUTDOWN = "Shutting down the bot..."
BOT_SETUP_COMPLETE = "Bot setup complete"
BOT_LOGGED_IN = "Logged in as {name} (ID: {id})"
BOT_GUILD_COUNT = "Connected to {count} guilds"

# Music player messages
SONG_ADDED = "Added to queue"
SONG_ADDED_TO_QUEUE = "Added to queue: **{title}**"
NOW_PLAYING = "Now playing"
NOW_PLAYING_TITLE = "Now playing: **{title}**"
QUEUE_TITLE = "Music Queue"
QUEUE_NOW_PLAYING = "Now Playing:"
QUEUE_UP_NEXT = "Up Next:"
QUEUE_PAGE_INFO = "Page {page}/{max_pages} | {total} songs in queue"
QUEUE_EMPTY = "The queue is empty."
CANNOT_REMOVE_CURRENT = "Cannot remove the currently playing song. Use /skip instead."
LOOP_NOT_AVAILABLE = "Loop feature is not available in the current version."

# Loop mode messages
LOOP_MODE_INVALID = "Invalid mode. Use 'off', 'song', or 'queue'."
LOOP_MODES = {
    "off": "Loop mode disabled",
    "song": "Looping current song",
    "queue": "Looping entire queue"
}

# Error messages
ERROR_COMMAND = "An error occurred: {error}"
ERROR_MISSING_ARG = "Missing required argument: {param}"
ERROR_BAD_ARG = "Bad argument: {error}"
ERROR_NO_RESULTS = "No results found"

# Setup and configuration messages
CONFIG_TOKEN_MISSING = "Error: .env file not found! Creating one now..."
CONFIG_TOKEN_CREATE = "Please add your Discord token to the .env file at {path}"
CONFIG_TOKEN_PROMPT = "Please enter your Discord bot token: "
CONFIG_TOKEN_EMPTY = "No token provided. Please add your Discord token to the .env file."
CONFIG_TOKEN_SAVED = "Token saved to .env file."
CONFIG_DEPS_MISSING = "Missing dependencies: {deps}"
CONFIG_DEPS_INSTALL_PROMPT = "Would you like to install them now? (y/n): "
CONFIG_DEPS_SKIPPED = "Dependencies not installed. The bot may not work correctly."
CONFIG_FFMPEG_NOT_FOUND = "ffmpeg not found!"
CONFIG_FFMPEG_INSTALL = "Please install ffmpeg:"
CONFIG_FFMPEG_WIN = "  - Download from https://ffmpeg.org/download.html\n  - Add it to your PATH"
CONFIG_FFMPEG_MAC = "  - Install with Homebrew: brew install ffmpeg"
CONFIG_FFMPEG_LINUX = "  - Ubuntu/Debian: sudo apt-get install ffmpeg\n  - CentOS/RHEL: sudo yum install ffmpeg"
CONFIG_FFMPEG_PATH = "Found FFmpeg at: {path}"
CONFIG_FFMPEG_VERSION = "FFmpeg version: {version}"
CONFIG_ISSUES = "Please resolve the issues above and try again."
CONFIG_START = "Starting NeruBot..."

# VPS setup messages
VPS_SETUP_TITLE = "\n=== Setting up NeruBot for VPS deployment ==="
VPS_USERNAME_PROMPT = "Enter your VPS username: "
VPS_GROUP_PROMPT = "Enter your VPS user group (usually the same as username): "
VPS_PATH_PROMPT = "Enter the absolute path where the bot will be located on the VPS (e.g. /home/{username}/nerubot): "
VPS_SERVICE_CREATED = "\nService file created at: {path}"
VPS_DEPLOY_HEADER = "\nTo deploy on your VPS:"
VPS_DEPLOY_COPY = "1. Copy the project to your VPS: scp -r {project_dir} {username}@your-vps-ip:{bot_dir}"
VPS_DEPLOY_DEPS = "2. Install dependencies on your VPS: pip3 install -r {bot_dir}/requirements.txt"
VPS_DEPLOY_FFMPEG = "3. Install ffmpeg on your VPS: sudo apt-get install ffmpeg"
VPS_DEPLOY_SERVICE = "4. Copy the service file: sudo cp {bot_dir}/nerubot.service /etc/systemd/system/"
VPS_DEPLOY_ENABLE = "5. Enable and start the service: sudo systemctl enable nerubot && sudo systemctl start nerubot"
VPS_MONITOR_HEADER = "\nTo monitor the bot:"
VPS_MONITOR_STATUS = "  - Check status: sudo systemctl status nerubot"
VPS_MONITOR_LOGS = "  - View logs: sudo journalctl -u nerubot -f"

# FFmpeg messages
FFMPEG_CREATING_SOURCE = "Creating FFmpegPCMAudio with explicit path: {path}"
FFMPEG_DEFAULT_PATH = "Creating FFmpegPCMAudio with system default path"
FFMPEG_ERROR = "Error creating FFmpegPCMAudio: {error}"
FFMPEG_FALLBACK = "Trying fallback method for FFmpegPCMAudio"
FFMPEG_FOUND_ENV = "Found FFmpeg path in environment: {path}"
FFMPEG_FOUND_CONFIG = "Found FFmpeg path in config: {path}"
FFMPEG_FOUND_DEFAULT = "Found FFmpeg at default path: {path}"
FFMPEG_FOUND_PATH = "Found FFmpeg in PATH"
FFMPEG_NOT_FOUND = "Could not find FFmpeg path, will rely on system PATH"
FFMPEG_USING_PATH = "Using FFmpeg path: {path}"

# Channel messages
VOICE_JOINED = "Joined {channel_name}"
VOICE_NOT_CONNECTED = "I'm not connected to a voice channel."
VOICE_JOIN_FAILED = "Failed to join voice channel: {error}"
VOICE_DISCONNECTED = "Disconnected from voice channel"
VOICE_LEAVE_FAILED = "Failed to leave voice channel: {error}"
USER_NOT_IN_CHANNEL = "You are not connected to a voice channel."

# Music control messages
PLAYBACK_ERROR = "Failed to play: {error}"
SONG_SKIPPED_NEXT = "⏭️ Skipped! Next up: **{title}**"
SONG_SKIPPED_NO_MORE = "⏭️ Skipped! No more songs in queue."
SONG_REMOVED = "Removed **{title}** from the queue."
SONG_REMOVE_INVALID = "Invalid song index."
QUEUE_CLEARED = "Stopped playing and cleared the queue."
QUEUE_SHUFFLED = "Shuffled the queue!"
NOTHING_PLAYING = "I'm not playing anything right now."
NOTHING_PAUSED = "Nothing is paused right now."
PAUSED = "Paused playback."
RESUMED = "Resumed playback."
VOLUME_SET = "Volume set to {volume}%"
VOLUME_RANGE = "Volume must be between 0 and 100."

# Help command messages
HELP_TITLE = "NeruBot Help"
HELP_DESCRIPTION = "Here are the available commands:"
HELP_MUSIC_COMMANDS_TITLE = "Music Commands"
HELP_JOIN_DESC = "Join your voice channel"
HELP_LEAVE_DESC = "Leave the voice channel"
HELP_PLAY_DESC = "Play a song from URL or search query"
HELP_STOP_DESC = "Stop playing and clear the queue"
HELP_PAUSE_DESC = "Pause the current song"
HELP_RESUME_DESC = "Resume the current song"
HELP_SKIP_DESC = "Skip to the next song"
HELP_VOLUME_DESC = "Set the volume"
HELP_NOW_DESC = "Show the current song"
HELP_QUEUE_DESC = "Show the song queue"
HELP_REMOVE_DESC = "Remove a song from the queue"
HELP_SHUFFLE_DESC = "Shuffle the queue"
HELP_LOOP_DESC = "Set loop mode"

# Music service messages
PLAYLIST_NO_VIDEOS = "❌ No playable videos found in this playlist."
PLAYLIST_FOUND = "🎵 Found playlist with {count} available videos. Playing first video: **{title}**"
NO_AUDIO_STREAM = "❌ Could not find a playable audio stream for this video."
NOW_PLAYING_NOTIFICATION = "Now playing: **{title}**"