"""
Settings and configuration constants
All configurable values, limits, timeouts, and technical settings
"""

# ============================
# BOT SETTINGS
# ============================

BOT_CONFIG = {
    "name": "NeruBot",
    "version": "2.1.0",
    "prefix": "!",
    "description": "A Discord music bot with clean architecture and high-quality audio",
    "status": "🎵 Ready to play music!",
    "activity_type": "listening",  # listening, playing, watching, streaming
}

# ============================
# FEATURE FLAGS
# ============================

FEATURES = {
    "music": True,
    "news": True,
    "help_system": True,
    "advanced_logging": True,
    "auto_disconnect": True,
    "24_7_mode": True,
    "volume_control": False,  # Disabled until implemented
    "playlist_support": True,
}

# ============================
# LIMITS AND TIMEOUTS
# ============================

LIMITS = {
    # Music limits
    "max_queue_size": 100,
    "max_search_results": 5,
    "max_song_duration": 3600,  # 1 hour in seconds
    
    # API timeouts (seconds)
    "search_timeout": 15.0,
    "conversion_timeout": 20.0,
    "discord_timeout": 30.0,
    "play_command_timeout": 25.0,
    "spotify_api_timeout": 10.0,
    
    # Connection timeouts
    "idle_disconnect_time": 300,  # 5 minutes
    "voice_connect_timeout": 30,
    
    # Rate limits
    "commands_per_minute": 10,
    "searches_per_minute": 5,
}

# ============================
# AUDIO SETTINGS
# ============================

AUDIO_CONFIG = {
    # FFmpeg settings
    "ffmpeg_options": {
        "before_options": "-reconnect 1 -reconnect_streamed 1 -reconnect_delay_max 5",
        "options": "-vn -filter:a volume=0.5"  # Default volume at 50%
    },
    
    # Opus library paths (for voice)
    "opus_paths": [
        "/opt/homebrew/lib/libopus.dylib",  # Apple Silicon Homebrew
        "/usr/local/lib/libopus.dylib",     # Intel Homebrew
        "/opt/homebrew/lib/libopus.0.dylib",
        "/usr/local/lib/libopus.0.dylib",
        "/usr/lib/x86_64-linux-gnu/libopus.so.0",  # Linux
        "/usr/lib/libopus.so.0",
        "libopus.so.0",
        "libopus.dylib",
        "opus"
    ],
    
    # Audio quality settings
    "bitrate": 128,  # kbps
    "sample_rate": 48000,  # Hz
    "channels": 2,  # Stereo
}

# ============================
# LOGGING SETTINGS
# ============================

LOGGING_CONFIG = {
    "level": "INFO",  # DEBUG, INFO, WARNING, ERROR, CRITICAL
    "file": "bot.log",
    "max_file_size": 10 * 1024 * 1024,  # 10MB
    "backup_count": 5,
    "format": "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    "date_format": "%Y-%m-%d %H:%M:%S",
}

# ============================
# DISCORD SETTINGS
# ============================

DISCORD_CONFIG = {
    # Required intents
    "intents": {
        "guilds": True,
        "guild_messages": True,
        "guild_voice_states": True,
        "message_content": True,
    },
    
    # Embed colors (hex values)
    "colors": {
        "success": 0x00FF00,    # Green
        "error": 0xFF0000,      # Red  
        "warning": 0xFFA500,    # Orange
        "info": 0x0099FF,       # Blue
        "music": 0x9932CC,      # Purple
        "spotify": 0x1DB954,    # Spotify Green
        "youtube": 0xFF0000,    # YouTube Red
        "soundcloud": 0xFF7700, # SoundCloud Orange
    },
    
    # Command sync settings
    "sync_commands_on_ready": True,
    "sync_commands_globally": True,
}

# ============================
# MUSIC SOURCES SETTINGS
# ============================

MUSIC_SOURCES = {
    "youtube": {
        "enabled": True,
        "emoji": "▶️",
        "name": "YouTube",
        "priority": 1,
    },
    "spotify": {
        "enabled": True,
        "emoji": "💚", 
        "name": "Spotify",
        "priority": 2,
        "requires_youtube_fallback": True,
    },
    "soundcloud": {
        "enabled": True,
        "emoji": "🧡",
        "name": "SoundCloud", 
        "priority": 3,
    },
    "direct": {
        "enabled": True,
        "emoji": "🔗",
        "name": "Direct Links",
        "priority": 4,
    }
}

# ============================
# HELP SYSTEM SETTINGS
# ============================

HELP_CONFIG = {
    "page_timeout": 60,  # seconds
    "max_fields_per_page": 10,
    "thumbnail_url": "https://imgur.com/a/gjS2Rea",
    "categories": [
        "🎵 Music - Playback",
        "🎵 Music - Voice", 
        "🎵 Music - Queue",
        "🎵 Music - Info",
        "🤖 General",
    ]
}

# ============================
# DEFAULT VALUES
# ============================

DEFAULTS = {
    "unknown_duration": "Unknown",
    "unknown_artist": "Unknown Artist", 
    "unknown_title": "Unknown Title",
    "unknown_album": "Unknown Album",
    "volume": 50,  # Default volume percentage
    "loop_mode": "off",
    "24_7_mode": False,
}

# ============================
# EMOJI MAPPING
# ============================

EMOJIS = {
    # Music control
    "play": "▶️",
    "pause": "⏸️", 
    "stop": "⏹️",
    "skip": "⏭️",
    "previous": "⏮️",
    "volume": "🔊",
    "mute": "🔇",
    
    # Loop modes
    "loop_off": "🔁",
    "loop_single": "🔂",
    "loop_queue": "🔁",
    "shuffle": "🔀",
    
    # Status indicators
    "success": "✅",
    "error": "❌", 
    "warning": "⚠️",
    "info": "ℹ️",
    "loading": "⏱️",
    "music": "🎵",
    
    # Voice
    "joined": "🔊",
    "left": "👋",
    "deafened": "🔇",
    
    # Navigation
    "left_arrow": "⬅️",
    "right_arrow": "➡️", 
    "close": "❌",
    
    # Sources (from MUSIC_SOURCES)
    **{k: v["emoji"] for k, v in MUSIC_SOURCES.items()}
}
