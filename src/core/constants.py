"""
Discord Music Bot - Constants and Configuration
DEPRECATED: This file is being phased out. Use src/config/ modules instead.
This file now imports from the new config system for backward compatibility.
"""

# Import from new config system
from src.config.settings import (
    BOT_CONFIG, LIMITS, AUDIO_CONFIG, DISCORD_CONFIG, 
    MUSIC_SOURCES, DEFAULTS, EMOJIS
)
from src.config.messages import (
    MSG_SUCCESS, MSG_ERROR, MSG_INFO, MSG_HELP, 
    CMD_DESCRIPTIONS, LOG_MSG
)

# ============================
# BACKWARD COMPATIBILITY
# Export old constant names for existing code
# ============================

# Bot configuration
BOT_NAME = BOT_CONFIG["name"]
BOT_VERSION = BOT_CONFIG["version"]  
BOT_DEFAULT_STATUS = BOT_CONFIG["status"]
BOT_DEFAULT_ACTIVITY_TYPE = BOT_CONFIG["activity_type"]

# Timeouts and limits
TIMEOUT_SPOTIFY_API = LIMITS["spotify_api_timeout"]
TIMEOUT_SEARCH = LIMITS["search_timeout"]
TIMEOUT_CONVERSION = LIMITS["conversion_timeout"]
TIMEOUT_DISCORD_INTERACTION = LIMITS["discord_timeout"]
TIMEOUT_PLAY_COMMAND = LIMITS["play_command_timeout"]
MAX_QUEUE_SIZE = LIMITS["max_queue_size"]
MAX_SEARCH_RESULTS = LIMITS["max_search_results"]
IDLE_DISCONNECT_TIME = LIMITS["idle_disconnect_time"]

# Default values
DEFAULT_UNKNOWN_DURATION = DEFAULTS["unknown_duration"]
DEFAULT_UNKNOWN_ARTIST = DEFAULTS["unknown_artist"]
DEFAULT_UNKNOWN_TITLE = DEFAULTS["unknown_title"]
DEFAULT_UNKNOWN_ALBUM = DEFAULTS["unknown_album"]

# Colors
COLOR_SUCCESS = DISCORD_CONFIG["colors"]["success"]
COLOR_ERROR = DISCORD_CONFIG["colors"]["error"]
COLOR_WARNING = DISCORD_CONFIG["colors"]["warning"]
COLOR_INFO = DISCORD_CONFIG["colors"]["info"]
COLOR_MUSIC = DISCORD_CONFIG["colors"]["music"]
COLOR_SPOTIFY = DISCORD_CONFIG["colors"]["spotify"]
COLOR_YOUTUBE = DISCORD_CONFIG["colors"]["youtube"]

# Emojis (backward compatibility)
EMOJI_MUSIC = EMOJIS["music"]
EMOJI_PLAYING = EMOJIS["play"]
EMOJI_PAUSED = EMOJIS["pause"]
EMOJI_STOPPED = EMOJIS["stop"]
EMOJI_SKIPPED = EMOJIS["skip"]
EMOJI_LOOP_OFF = EMOJIS["loop_off"]
EMOJI_LOOP_SINGLE = EMOJIS["loop_single"]
EMOJI_LOOP_QUEUE = EMOJIS["loop_queue"]
EMOJI_SHUFFLE = EMOJIS["shuffle"]
EMOJI_VOLUME = EMOJIS["volume"]
EMOJI_JOINED = EMOJIS["joined"]
EMOJI_LEFT = EMOJIS["left"]
EMOJI_SUCCESS = EMOJIS["success"]
EMOJI_ERROR = EMOJIS["error"]
EMOJI_WARNING = EMOJIS["warning"]
EMOJI_INFO = EMOJIS["info"]
EMOJI_LOADING = EMOJIS["loading"]

# Source emojis
SOURCE_EMOJI = {k: v["emoji"] for k, v in MUSIC_SOURCES.items()}
SOURCE_EMOJI["unknown"] = EMOJIS["music"]

# Loop mode emojis (for backward compatibility)
LOOP_EMOJI = {
    "off": EMOJIS["loop_off"],
    "single": EMOJIS["loop_single"], 
    "queue": EMOJIS["loop_queue"]
}

# ============================
# SUCCESS MESSAGES
# ============================

MSG_SUCCESS = {
    "joined_channel": f"{EMOJI_JOINED} Joined Voice Channel\nConnected to **{{channel_name}}**",
    "left_channel": f"{EMOJI_LEFT} Left Voice Channel\nDisconnected from **{{channel_name}}**",
    "song_added_queue": f"{EMOJI_SUCCESS} Added to Queue\n**{{title}}**\nPosition: {{position}}",
    "song_playing": f"{EMOJI_PLAYING} Now Playing\n**{{title}}**",
    "song_skipped": f"{EMOJI_SKIPPED} Skipped\nSkipped to the next song",
    "song_skipped_last": f"{EMOJI_SKIPPED} Skipped\nSkipped the last song. Queue is now empty.",
    "music_paused": f"{EMOJI_PAUSED} Paused\nMusic paused",
    "music_resumed": f"{EMOJI_PLAYING} Resumed\nMusic resumed",
    "music_stopped": f"{EMOJI_STOPPED} Stopped\nMusic stopped and queue cleared",
    "queue_cleared": f"{EMOJI_SUCCESS} Queue Cleared\nThe music queue has been cleared.",
    "loop_changed": f"{EMOJI_LOOP_OFF} Loop Mode\nLoop mode set to: **{{mode}}**",
    "247_enabled": f"{EMOJI_SUCCESS} 24/7 Mode Enabled\nBot will stay connected even when not playing music.",
    "247_disabled": f"{EMOJI_WARNING} 24/7 Mode Disabled\nBot will disconnect after 5 minutes of inactivity."
}
# ============================
# AUDIO CONFIGURATION (kept for convenience)
# ============================

FFMPEG_OPTIONS = AUDIO_CONFIG["ffmpeg_options"]
OPUS_PATHS = AUDIO_CONFIG["opus_paths"]

# ============================
# DEPRECATED - Use config.messages and config.settings instead
# Remove this section when all code is updated
# ============================
