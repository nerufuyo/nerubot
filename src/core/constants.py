"""
Discord Music Bot - All Constants and Messages
This file contains all configurable values, strings, messages, and settings
for easy maintenance and localization.
"""

# ============================
# BOT CONFIGURATION
# ============================

BOT_NAME = "NeruBot"
BOT_VERSION = "2.0.0"
BOT_DEFAULT_STATUS = "🎵 Ready to play music!"
BOT_DEFAULT_ACTIVITY_TYPE = "listening"  # listening, playing, watching, streaming

# ============================
# TIMEOUTS AND LIMITS
# ============================

# Music processing timeouts (seconds)
TIMEOUT_SPOTIFY_API = 10.0
TIMEOUT_SEARCH = 15.0
TIMEOUT_CONVERSION = 20.0
TIMEOUT_DISCORD_INTERACTION = 30.0
TIMEOUT_PLAY_COMMAND = 25.0  # Timeout for play command processing

# Music limits
MAX_QUEUE_SIZE = 100
MAX_SEARCH_RESULTS = 5
IDLE_DISCONNECT_TIME = 300  # 5 minutes

# ============================
# DEFAULT VALUES
# ============================

DEFAULT_UNKNOWN_DURATION = "Unknown"
DEFAULT_UNKNOWN_ARTIST = "Unknown Artist"
DEFAULT_UNKNOWN_TITLE = "Unknown Title"
DEFAULT_UNKNOWN_ALBUM = "Unknown Album"

# ============================
# EMOJI AND SYMBOLS
# ============================

EMOJI_MUSIC = "🎵"
EMOJI_PLAYING = "▶️"
EMOJI_PAUSED = "⏸️"
EMOJI_STOPPED = "⏹️"
EMOJI_SKIPPED = "⏭️"
EMOJI_LOOP_OFF = "🔁"
EMOJI_LOOP_SINGLE = "🔂"
EMOJI_LOOP_QUEUE = "🔁"
EMOJI_SHUFFLE = "🔀"
EMOJI_VOLUME = "🔊"
EMOJI_JOINED = "🔊"
EMOJI_LEFT = "👋"
EMOJI_SUCCESS = "✅"
EMOJI_ERROR = "❌"
EMOJI_WARNING = "⚠️"
EMOJI_INFO = "ℹ️"
EMOJI_LOADING = "⏱️"

# Source emojis
SOURCE_EMOJI = {
    'youtube': '▶️',
    'spotify': '💚',
    'soundcloud': '🧡',
    'direct': '🔗',
    'unknown': '🎵'
}

# Loop mode emojis (for backward compatibility with music cog)
LOOP_EMOJI = {
    "off": EMOJI_LOOP_OFF,
    "single": EMOJI_LOOP_SINGLE,
    "queue": EMOJI_LOOP_QUEUE
}

# ============================
# DISCORD EMBED COLORS
# ============================

COLOR_SUCCESS = 0x00FF00    # Green
COLOR_ERROR = 0xFF0000      # Red
COLOR_WARNING = 0xFFA500    # Orange
COLOR_INFO = 0x0099FF       # Blue
COLOR_MUSIC = 0x9932CC      # Purple
COLOR_SPOTIFY = 0x1DB954    # Spotify Green
COLOR_YOUTUBE = 0xFF0000    # YouTube Red

# ============================
# COMMAND DESCRIPTIONS
# ============================

CMD_DESCRIPTIONS = {
    # Music commands
    "join": "Join your voice channel",
    "leave": "Leave the voice channel",
    "play": "Play music from YouTube, Spotify, or SoundCloud",
    "stop": "Stop music and clear the queue",
    "pause": "Pause the current song",
    "resume": "Resume the current song",
    "skip": "Skip the current song",
    "queue": "Show the music queue",
    "nowplaying": "Show the currently playing song",
    "loop": "Toggle loop mode (off/single/queue)",
    "clear": "Clear the music queue",
    "sources": "Show available music sources",
    "247": "Toggle 24/7 mode (stay connected)",
    
    # Help commands
    "help": "Show help information",
    "about": "Show bot information",
    "commands": "Show all available commands",
    "features": "Show bot features",
    
    # News commands
    "news": "Get latest news",
    "news_sources": "List news sources",
    "news_status": "Show news configuration"
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
# ERROR MESSAGES
# ============================

MSG_ERROR = {
    "not_in_voice": f"{EMOJI_ERROR} You need to be in a voice channel!",
    "bot_not_connected": f"{EMOJI_ERROR} I'm not connected to a voice channel!",
    "nothing_playing": f"{EMOJI_ERROR} Nothing is currently playing!",
    "nothing_paused": f"{EMOJI_ERROR} Nothing is paused right now!",
    "queue_empty": f"{EMOJI_INFO} The queue is empty.",
    "no_results": f"{EMOJI_ERROR} Could not find any songs matching your query. Try a different search term or check the URL.",
    "search_timeout": f"{EMOJI_LOADING} Request timed out. The song might be from a slow source or playlist. Please try again with a simpler query.",
    "processing_timeout": f"{EMOJI_LOADING} Processing timed out. This might be a complex playlist or the service is slow. Please try again.",
    "conversion_failed": f"{EMOJI_ERROR} Could not find any information for that song. Try a different search term or check if the link is valid.",
    "join_failed": f"{EMOJI_ERROR} Failed to join voice channel: {{error}}",
    "unexpected_error": f"{EMOJI_ERROR} An unexpected error occurred. Please try again or contact support.",
    "invalid_loop_mode": f"{EMOJI_ERROR} Invalid loop mode. Use: off, single, or queue",
    "spotify_unavailable": f"{EMOJI_WARNING} Spotify is currently unavailable. Trying YouTube instead..."
}

# ============================
# INFO MESSAGES
# ============================

MSG_INFO = {
    "queue_title": f"{EMOJI_MUSIC} Music Queue",
    "now_playing_prefix": f"{EMOJI_PLAYING} **Now Playing:**",
    "up_next_prefix": f"{EMOJI_MUSIC} **Up Next:**",
    "queue_position": "Position {{position}} in queue",
    "processing": f"{EMOJI_LOADING} Processing your request...",
    "searching": f"{EMOJI_LOADING} Searching for music...",
    "converting": f"{EMOJI_LOADING} Converting to playable format...",
    "goodbye": f"{EMOJI_LEFT} Disconnected due to inactivity. Use `/join` to reconnect!",
    "bot_ready": f"{EMOJI_SUCCESS} {{bot_name}} is ready! Use `/help` to get started."
}

# ============================
# HELP MESSAGES
# ============================

MSG_HELP = {
    "tips_spotify": "\n\n💡 **Tips:**\n• Try searching with simpler terms\n• Check if the Spotify/URL link is valid\n• Try a YouTube link instead",
    "tips_timeout": "\n\n💡 **Try:** Using a direct YouTube link or simpler search terms",
    "sources_info": "You can play music from these sources:",
    "youtube_desc": "Play music from YouTube links or search for songs\nExample: `/play despacito` or `/play https://www.youtube.com/watch?v=kJQP7kiw5Fk`",
    "spotify_desc": "Play music from Spotify links (tracks, albums, playlists, artists)\nExample: `/play https://open.spotify.com/track/6habFhsOp2NvshLv26jCFK`",
    "soundcloud_desc": "Play music from SoundCloud links (tracks, playlists, users)\nExample: `/play https://soundcloud.com/artist/track`",
    "direct_desc": "Play music directly from audio file links\nExample: `/play https://example.com/music.mp3`",
    "sources_footer": "Use /play followed by a search term or URL to play music from any of these sources!"
}

# ============================
# LOG MESSAGES
# ============================

LOG_MSG = {
    "bot_starting": "Starting {bot_name} v{version}...",
    "bot_ready": "Bot is ready! Logged in as {user}",
    "bot_disconnected": "Bot disconnected",
    "cog_loaded": "Loaded cog: {cog_name}",
    "cog_failed": "Failed to load cog {cog_name}: {error}",
    "music_added_queue": "Added to queue: {title} in guild {guild_id}",
    "music_now_playing": "Now playing: {title} in guild {guild_id} from {source}",
    "music_search_start": "Searching for: {query}",
    "music_search_results": "Found {count} results for '{query}' from {source}",
    "music_search_timeout": "Search timeout for query: {query}",
    "music_conversion_success": "Successfully converted to playable: {title}",
    "music_conversion_failed": "Failed to convert to playable: {title}",
    "voice_joined": "Joined voice channel: {channel} in guild {guild_id}",
    "voice_left": "Left voice channel: {channel} in guild {guild_id}",
    "idle_timer_start": "Starting idle timer for guild {guild_id}",
    "idle_disconnect": "Disconnecting due to inactivity for guild {guild_id}",
    "emergency_cleanup": "Performed emergency cleanup for guild {guild_id}",
    "spotify_timeout": "Spotify API timeout for: {query}",
    "youtube_search": "Trying YouTube search with: {query}",
    "error_general": "Error in {function}: {error}",
    "error_critical": "Critical error in {function} for guild {guild_id}: {error}",
    
    # Source adapter specific messages
    "source_search_start": "Searching for: {query}",
    "source_determined": "Determined source: {source}",
    "source_search_results": "Found {count} results for '{query}' from {source}",
    "source_no_results": "No results found for '{query}' from {source}",
    "source_search_error": "Error searching '{query}' from {source}: {error}",
    "source_convert_start": "Converting to playable: {title} from {source}",
    "source_convert_none": "Cannot get playable version of None result",
    "source_youtube_playable": "YouTube result already playable: {title}",
    "source_spotify_converted": "Successfully converted Spotify to playable: {title}",
    "source_spotify_failed": "Failed to convert Spotify to playable: {title}",
    "source_using_original": "Using original result as playable: {title}",
    "source_convert_error": "Error converting '{title}' to playable: {error}",
    
    # Spotify specific messages
    "spotify_no_credentials": "Spotify credentials not found. Spotify support will be limited.",
    "spotify_not_initialized": "Spotify adapter not initialized. Falling back to YouTube.",
    "spotify_search_timeout": "Spotify search timeout for query: {query}",
    "spotify_no_results": "No Spotify results found for: {query}",
    "spotify_results_found": "Found {count} Spotify results for: {query}",
    "spotify_search_error": "Spotify search error for '{query}': {error}",
    "spotify_track_timeout": "Spotify track lookup timeout for ID: {track_id}",
    "spotify_track_error": "Error processing Spotify track URL '{url}': {error}",
    "spotify_album_error": "Error processing Spotify album URL: {error}",
    "spotify_playlist_error": "Error processing Spotify playlist URL: {error}",
    "spotify_artist_error": "Error processing Spotify artist URL: {error}",
    "spotify_convert_error": "Error converting Spotify track: {error}",
    "spotify_convert_none": "Cannot convert None Spotify result to playable",
    "spotify_converting": "Converting Spotify track to playable: {title}",
    "spotify_youtube_trying": "Trying YouTube search with: {query}",
    "spotify_youtube_found": "Found YouTube match: {title}",
    "spotify_youtube_no_results": "No results for query: {query}",
    "spotify_youtube_failed": "YouTube search failed for '{query}': {error}",
    "spotify_no_youtube": "Could not find YouTube equivalent for Spotify track: {title}",
    "spotify_conversion_success": "Successfully converted Spotify track to playable: {title}",
    
    # YouTube specific messages
    "youtube_search_error": "YouTube search error: {error}",
    "youtube_process_error": "Error processing YouTube data: {error}",
    
    # SoundCloud specific messages
    "soundcloud_using_ytdlp": "Using yt-dlp for SoundCloud playback",
    "soundcloud_search_error": "SoundCloud search error: {error}",
    "soundcloud_fallback_error": "Error in SoundCloud fallback search: {error}",
    "soundcloud_convert_error": "Error converting SoundCloud to playable: {error}"
}

# ============================
# LOOP MODE CONFIGURATION
# ============================

LOOP_MODES = {
    "off": {
        "name": "Off",
        "emoji": EMOJI_LOOP_OFF,
        "description": "Loop mode disabled"
    },
    "single": {
        "name": "Single",
        "emoji": EMOJI_LOOP_SINGLE,
        "description": "Looping current song"
    },
    "queue": {
        "name": "Queue",
        "emoji": EMOJI_LOOP_QUEUE,
        "description": "Looping entire queue"
    }
}

# ============================
# SPOTIFY SEARCH STRATEGIES
# ============================

SPOTIFY_SEARCH_STRATEGIES = [
    "{title} {artist} audio",      # Original strategy
    "{title} {artist}",            # Simple strategy
    "{artist} {title} official",   # With "official"
    "{title} {artist} lyrics"      # With "lyrics"
]

# ============================
# FFMPEG CONFIGURATION
# ============================

FFMPEG_OPTIONS = {
    'before_options': '-reconnect 1 -reconnect_streamed 1 -reconnect_delay_max 5',
    'options': '-vn'
}

# Opus library paths (for Discord voice)
OPUS_PATHS = [
    '/opt/homebrew/lib/libopus.dylib',     # Homebrew on Apple Silicon
    '/usr/local/lib/libopus.dylib',        # Homebrew on Intel
    '/opt/homebrew/lib/libopus.0.dylib',
    '/usr/local/lib/libopus.0.dylib',
    'libopus.dylib',                       # System path
    'opus'                                 # Let the system find it
]
