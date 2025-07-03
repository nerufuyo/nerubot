"""
Discord Music Bot - Constants and Configuration
DEPRECATED: This file is being phased out. Use src/config/ modules instead.
This file now imports from the new config system for backward compatibility.
"""

# Import from new config system
from src.config.settings import (
    BOT_CONFIG, LIMITS, AUDIO_CONFIG, DISCORD_CONFIG, 
    MUSIC_SOURCES, DEFAULTS, EMOJIS, SPOTIFY_SEARCH_STRATEGIES
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

# ============================
# CONFESSION SYSTEM CONSTANTS
# ============================

# Confession System Configuration
CONFESSION_CONSTANTS = {
    "ids": {
        "confession_prefix": "CONF",
        "reply_prefix": "REPLY",
        "id_format": "{prefix}-{id:03d}",
        "reply_format": "{prefix}-{confession_id:03d}-{letter}",
    },
    "messages": {
        "titles": {
            "confession": "{emoji} Confession #{id:03d}",
            "reply": "{emoji} {reply_id}",
            "setup_intro": "{emoji} Anonymous Confession System",
            "setup_success": "{emoji} Confession System Setup Complete",
            "settings": "{emoji} Confession Settings",
            "stats": "{emoji} Confession Statistics",
        },
        "descriptions": {
            "intro": (
                "Welcome to the anonymous confession system! This is a safe space where you can share your thoughts, "
                "feelings, and experiences completely anonymously.\n\n"
                "**How it works:**\n"
                "• Click the **Create New Confession** button to share anonymously\n"
                "• Each confession gets its own discussion thread\n"
                "• Anyone can reply anonymously using the **Reply** button\n"
                "• All messages are posted by the bot - complete anonymity guaranteed\n\n"
                "**Features:**\n"
                "• Unique ID system for easy reference (CONF-001, REPLY-001-A, etc.)\n"
                "• Support for text and attachments\n"
                "• Organized in threads for better discussions\n"
                "• No way to trace messages back to authors\n\n"
                "Click the button below to create your first confession!"
            ),
            "setup_success": "The confession system has been set up in {channel}",
        },
        "footers": {
            "confession": "ID: {id} | {reply_emoji} Reply",
            "reply": "ID: {id} | {reply_emoji} Reply",
            "intro": "All confessions and replies are completely anonymous",
        },
        "buttons": {
            "create_confession": "Create New Confession",
            "reply": "Reply",
        },
        "modals": {
            "confession_title": "Create New Confession",
            "reply_title": "Reply to Confession {id}",
            "fields": {
                "message_label": "Message",
                "message_placeholder": "Type your anonymous {type} here...",
                "attachments_label": "Attachments",
                "attachments_placeholder": "Paste image/video URLs here (optional, separate multiple with commas)",
                "confession_id_label": "Confession ID",
                "confession_id_placeholder": "Auto-populated",
            }
        },
        "success": {
            "confession_created": "Your confession has been submitted anonymously! (ID: `{id}`)",
            "reply_created": "Your reply has been posted anonymously in the confession thread! (ID: `{id}`)",
            "setup_complete": "Confession system setup complete",
        },
        "errors": {
            "system_unavailable": "Confession system is not available.",
            "empty_content": "Please provide some content for your {type}.",
            "confession_not_found": "Confession not found",
            "confession_not_in_guild": "Confession not found in this server",
            "content_too_long": "{type} too long (max {max} characters)",
            "no_confession_channel": "Confession channel is not set up for this server. Please ask an admin to set it up.",
            "channel_not_found": "Confession channel not found",
            "thread_not_found": "Thread not found",
            "no_permission": "No permission to access/post in {location}",
        }
    },
    "settings": {
        "defaults": {
            "max_confession_length": 2000,
            "max_reply_length": 1000,
            "thread_auto_archive_duration": 10080,  # 7 days
        },
        "limits": {
            "max_attachments": 10,
            "max_attachment_size": 8 * 1024 * 1024,  # 8MB
        }
    },
    "thread": {
        "name_format": "{emoji} Confession #{id:03d}",
        "auto_archive_duration": 10080,  # 7 days in minutes
    }
}

# File paths for confession system
CONFESSION_FILE_PATHS = {
    "data_dir": "data/confessions",
    "confessions_file": "data/confessions/confessions.json",
    "replies_file": "data/confessions/replies.json",
    "settings_file": "data/confessions/settings.json",
    "queue_file": "data/confessions/queue.json",
}

# Queue System Constants
QUEUE_CONSTANTS = {
    "max_queue_size": 1000,
    "queue_cleanup_interval": 3600,  # 1 hour in seconds
    "max_processing_time": 300,  # 5 minutes in seconds
    "retry_attempts": 3,
    "retry_delay": 5,  # seconds
}

# Logging Messages for Confession System
CONFESSION_LOG_MESSAGES = {
    "confession": {
        "created": "Created confession {id} for user {user_id} in guild {guild_id}",
        "posted": "Posted confession {id} to {guild_name} with thread {thread_id}",
        "failed_create": "Failed to create confession: {error}",
        "failed_post": "Error posting confession: {error}",
        "queued": "Queued confession {id} for processing",
    },
    "reply": {
        "created": "Created reply {id} for confession {confession_id} by user {user_id}",
        "posted": "Successfully posted reply {id} to confession {confession_id} thread",
        "failed_create": "Failed to create reply: {error}",
        "failed_post": "Error posting reply: {error}",
        "queued": "Queued reply {id} for processing",
    },
    "system": {
        "cog_loaded": "ConfessionCog loaded with persistent views",
        "guild_configured": "Guild {guild_id} configured for confessions",
        "queue_processed": "Processed {count} items from confession queue",
        "queue_error": "Error processing queue item: {error}",
        "queue_started": "Confession queue processor started",
        "queue_stopped": "Confession queue processor stopped",
    }
}
