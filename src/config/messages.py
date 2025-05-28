"""
Messages and strings configuration
All user-facing strings, error messages, help text, and localizable content
"""

# ============================
# BOT INFORMATION MESSAGES
# ============================

BOT_INFO = {
    "ready": "🎵 {bot_name} is online and ready to make your server amazing!",
    "disconnected": "Bot disconnected from Discord",
    "shutdown": "Bot is shutting down...",
    "welcome": "🎉 Thank you for adding {bot_name} to your server! I'm your new friendly companion ready to bring music, fun, and community features to your Discord! Use `/help` to discover all my amazing capabilities, or try `/about` to learn more about me!",
    "greeting": "👋 Hey there! I'm NeruBot, your friendly Discord companion created by **nerufuyo**! I'm here to make your server more entertaining with music, confessions, news, and much more!",
}

# ============================
# SUCCESS MESSAGES
# ============================

MSG_SUCCESS = {
    "music_added": "✅ Added **{title}** to the queue",
    "music_playing": "▶️ Now playing: **{title}**",
    "music_paused": "⏸️ Music paused",
    "music_resumed": "▶️ Music resumed",
    "music_stopped": "⏹️ Music stopped and queue cleared",
    "music_skipped": "⏭️ Skipped: **{title}**",
    "voice_joined": "🔊 Joined **{channel}**",
    "voice_left": "👋 Left voice channel",
    "queue_cleared": "🗑️ Queue cleared",
    "loop_set": "🔁 Loop mode set to: **{mode}**",
    "247_enabled": "🔄 24/7 mode enabled - I'll stay connected",
    "247_disabled": "🔄 24/7 mode disabled - I'll auto-disconnect when idle",
    "volume_set": "🔊 Volume set to **{volume}%**",
}

# ============================
# ERROR MESSAGES
# ============================

MSG_ERROR = {
    "not_in_voice": "❌ You need to be in a voice channel to use this command",
    "not_connected": "❌ I'm not connected to a voice channel",
    "nothing_playing": "❌ Nothing is currently playing",
    "queue_empty": "❌ The queue is empty",
    "invalid_query": "❌ Please provide a valid search query or URL",
    "search_failed": "❌ Could not find any results for: **{query}**",
    "search_timeout": "⏱️ Search timed out. Please try again",
    "conversion_failed": "❌ Failed to process audio for: **{title}**",
    "playback_error": "❌ An error occurred during playback",
    "permission_error": "❌ I don't have permission to join your voice channel",
    "queue_full": "❌ Queue is full (max {max_size} songs)",
    "invalid_volume": "❌ Volume must be between 0 and 100",
    "invalid_loop_mode": "❌ Invalid loop mode. Use: off, single, or queue",
    "command_error": "❌ An error occurred while processing your command",
    "no_results": "❌ No results found for your search",
    "already_connected": "❌ I'm already connected to a voice channel",
    "different_channel": "❌ You need to be in the same voice channel as me",
    "no_spotify_support": "⚠️ Spotify support requires additional setup",
}

# ============================
# INFO MESSAGES
# ============================

MSG_INFO = {
    "searching": "🔍 Searching for: **{query}**...",
    "processing": "⏱️ Processing your request...",
    "converting": "🔄 Converting audio...",
    "queue_position": "📍 Position in queue: **#{position}**",
    "queue_duration": "⏱️ Queue duration: **{duration}**",
    "now_playing_info": "▶️ **Now Playing:** {title}\n⏱️ **Duration:** {duration}\n👤 **Requested by:** {requester}",
    "idle_warning": "⚠️ I'll leave in 5 minutes due to inactivity",
    "auto_disconnect": "👋 Left voice channel due to inactivity",
    "spotify_track": "💚 Found Spotify track, searching on YouTube...",
    "playlist_detected": "📋 Playlist detected, adding {count} songs...",
    "bot_ready": "🎵 {bot_name} is ready to play music!",
}

# ============================
# HELP MESSAGES
# ============================

MSG_HELP = {
    "main_description": "👋 **Welcome to NeruBot!** I'm your friendly Discord companion created by **nerufuyo**!\n\n"
                       "I'm here to make your server amazing with music, community features, and much more! "
                       "Browse through the help pages using the buttons below to discover everything I can do.\n\n"
                       "**🎯 What I Offer:**\n"
                       "• 🎵 **Premium Music** - High-quality streaming from multiple sources\n"
                       "• 📝 **Anonymous Confessions** - Safe community sharing space\n"
                       "• 📰 **News Updates** - Stay informed with RSS feeds\n"
                       "• 🤖 **User-Friendly** - Easy commands with beautiful interfaces\n\n"
                       "Use the navigation buttons below to explore! ⬇️",
    "music_description": "Complete music streaming solution with high-quality audio",
    "general_description": "General bot commands and information",
    "usage_tips": [
        "💡 **Pro Tips:**",
        "• Use Spotify, YouTube, or SoundCloud links with `/play`",
        "• Try `/loop queue` to repeat your entire playlist",
        "• Use `/247` to keep the bot in your voice channel",
        "• Use `/sources` to see all supported music platforms"
    ],
    "commands": {
        "play": "Play music from YouTube, Spotify, or SoundCloud",
        "pause": "Pause the current song",
        "resume": "Resume the current song",
        "stop": "Stop music and clear queue",
        "skip": "Skip the current song",
        "join": "Join your voice channel",
        "leave": "Leave the voice channel",
        "volume": "Set the volume level",
        "queue": "Show the music queue",
        "nowplaying": "Show currently playing song",
        "clear": "Clear the music queue",
        "loop": "Toggle loop mode",
        "247": "Toggle 24/7 mode (stays in voice channel)",
        "sources": "Show all available music sources",
        "help": "Show this help menu",
        "about": "Learn about NeruBot and meet the creator nerufuyo!",
        "features": "Display detailed bot features and capabilities",
        "commands": "Show compact command reference card",
        # Confession commands
        "confess": "Submit an anonymous confession",
        "reply": "Reply to a confession anonymously",
        "confession-setup": "Set up confession channel (Admin only)",
        "confession-settings": "View confession settings (Admin only)",
        "confession-stats": "View confession statistics",
        # News commands
        "news-latest": "Get the latest news items",
        "news-sources": "List all configured news sources",
        "news-status": "Show current news configuration and status",
        "news-help": "Show help for news commands",
        "news-set-channel": "Set channel for automatic news updates (Admin only)",
        "news-start": "Start automatic news updates (Admin only)",
        "news-stop": "Stop automatic news updates (Admin only)",
        "news-add": "Add a news source (Admin only)",
        "news-remove": "Remove a news source (Admin only)"
    },
    "about": {
        "features": "• 🎵 Multi-source Music (YouTube, Spotify, SoundCloud)\n"
                   "• 📝 Anonymous Confession System\n"
                   "• 📰 News & RSS Feed Integration\n"
                   "• 🔄 Advanced Queue Management\n"
                   "• 🎛️ High-quality Audio\n"
                   "• 🏗️ Clean Architecture\n"
                   "• ⚡ Lightning-fast Performance\n"
                   "• 🛡️ Reliable & Stable\n"
                   "• 🎨 Beautiful User Interface",
        "links": "• [GitHub Repository](https://github.com/nerufuyo/nerubot)\n"
                "• [Invite NeruBot](https://discord.com/oauth2/authorize?client_id=yourid&permissions=8&scope=bot%20applications.commands)\n"
                "• [Support & Community](https://discord.gg/yourserver)\n"
                "• [Documentation](https://github.com/nerufuyo/nerubot/wiki)",
        "footer": "Crafted with ❤️ by nerufuyo | The friendliest bot in Discord!",
        "description": "Your all-in-one Discord companion for music, community, and entertainment!",
        "author_info": "Created by **nerufuyo** - A passionate developer dedicated to creating amazing Discord experiences!"
    },
    "features": {
        "title": "🚀 NeruBot - Your Amazing Discord Companion!",
        "description": "Discover all the incredible features that make NeruBot the perfect addition to your Discord server! Created with love by **nerufuyo**.",
        "current": (
            "**🎵 Premium Music Experience**\n"
            "• Multi-platform streaming (YouTube, Spotify, SoundCloud)\n"
            "• Crystal-clear audio quality with advanced processing\n"
            "• Smart queue management with loop modes\n"
            "• 24/7 mode - I'll stay in your voice channel!\n"
            "• Volume control and audio filters\n"
            "• Playlist support with bulk importing\n\n"
            
            "**📝 Anonymous Confession System**\n"
            "• Safe and secure anonymous confessions\n"
            "• Anonymous replies and interactions\n"
            "• Smart moderation with ID tracking\n"
            "• Customizable confession channels\n"
            "• Cooldown protection and anti-spam\n"
            "• Image support for confessions\n\n"
            
            "**📰 Intelligent News System**\n"
            "• Real-time RSS feed integration\n"
            "• Automatic news updates to your channels\n"
            "• Multiple configurable news sources\n"
            "• Smart filtering and formatting\n"
            "• Admin controls for news management\n"
            "• Breaking news notifications\n\n"
            
            "**🤖 User-Friendly Interface**\n"
            "• Modern slash commands throughout\n"
            "• Interactive help system with navigation\n"
            "• Beautiful embeds and rich formatting\n"
            "• Comprehensive error handling\n"
            "• Intuitive command structure\n"
            "• Fast and responsive performance\n"
        ),
        "sources": (
            "• ▶️ **YouTube** - Direct streaming with high quality\n"
            "• 💚 **Spotify** - Your favorite playlists and tracks\n"
            "• 🧡 **SoundCloud** - Discover independent artists\n"
            "• 🔗 **Direct Links** - MP3, MP4, and more formats\n"
        ),
        "upcoming": (
            "• 🎮 **Gaming Integration** - Discord activities and games\n"
            "• 📊 **Analytics Dashboard** - Server insights and statistics\n"
            "• 🎨 **Custom Themes** - Personalize your bot experience\n"
            "• 🌐 **Web Interface** - Manage your bot from anywhere\n"
            "• 🔔 **Smart Notifications** - Advanced alert system\n"
            "• 🎪 **Event Management** - Server events and scheduling\n"
        ),
        "footer": "More exciting features coming soon! | Created with ❤️ by nerufuyo"
    },
    "sources": {
        "youtube": "🎵 **YouTube**\nDirect playback from YouTube videos and playlists",
        "spotify": "💚 **Spotify**\nTrack/playlist support (searches equivalent on YouTube)",
        "soundcloud": "🧡 **SoundCloud**\nDirect streaming from SoundCloud tracks",
        "direct": "🔗 **Direct Links**\nSupports MP3, MP4, and other audio formats",
        "footer": "Use /play with any of these sources!"
    },
    "search_tips": "💡 **Tip:** Try being more specific with your search terms",
    "timeout_tips": "⏱️ **Tip:** Try a simpler search query or check your connection",
    
    # Command reference card content
    "command_card": {
        "title": "📋 NeruBot Command Reference",
        "description": "Quick reference for all available commands",
        "music_commands": (
            "`/play` - Play a song\n"
            "`/pause` - Pause playback\n"
            "`/resume` - Resume playback\n"
            "`/skip` - Skip current song\n"
            "`/stop` - Stop and clear queue\n"
            "`/queue` - Show music queue\n"
            "`/nowplaying` - Show current song\n"
            "`/clear` - Clear queue\n"
            "`/loop` - Toggle loop mode\n"
            "`/247` - Toggle 24/7 mode\n"
            "`/join` - Join voice channel\n"
            "`/leave` - Leave voice channel\n"
            "`/sources` - Show music sources"
        ),
        "confession_commands": (
            "`/confess [image]` - Submit anonymous confession\n"
            "`/reply <id> [image]` - Reply to a confession\n"
            "`/confession-setup` - Setup channel (Admin)\n"
            "`/confession-settings` - View settings (Admin)\n"
            "`/confession-stats` - View statistics"
        ),
        "news_commands": (
            "`/news latest [count]` - Get latest news\n"
            "`/news sources` - List news sources\n"
            "`/news status` - Show configuration\n"
            "`/news set-channel` - Set channel (Admin)\n"
            "`/news start/stop` - Control auto-posting (Admin)"
        ),
        "general_commands": (
            "`/help` - Detailed help pages\n"
            "`/commands` - This command card\n"
            "`/about` - Bot information\n"
            "`/features` - Feature showcase"
        ),
        "tips": (
            "**Pro Tips:**\n"
            "• Use `/play` with Spotify, YouTube or SoundCloud links\n"
            "• Try `/loop queue` to repeat your playlist\n"
            "• Use `/sources` to see all supported music sources"
        ),
        "footer": "Use /help for more detailed command information"
    }
}

# ============================
# COMMAND DESCRIPTIONS
# ============================

CMD_DESCRIPTIONS = {
    # Music playback commands
    "play": "Play music from YouTube, Spotify, or SoundCloud",
    "pause": "Pause the current song",
    "resume": "Resume the current song", 
    "stop": "Stop music and clear the queue",
    "skip": "Skip the current song",
    
    # Voice channel commands
    "join": "Join your voice channel",
    "leave": "Leave the voice channel", 
    "volume": "Set the volume (0-100)",
    
    # Queue management commands
    "queue": "Show the music queue",
    "nowplaying": "Show the currently playing song",
    "clear": "Clear the music queue",
    "loop": "Toggle loop mode (off/single/queue)",
    "247": "Toggle 24/7 mode (stay connected)",
    
    # Information commands
    "sources": "Show available music sources",
    
    # General commands
    "help": "Show help information with categories",
    "about": "Show bot information and statistics", 
    "commands": "Show compact command reference",
    "features": "Show bot features and capabilities",
    
    # Confession commands
    "confess": "Submit an anonymous confession (with optional image)",
    "reply": "Reply to a confession anonymously (with optional image)",
    "confession-setup": "Set up confession channel (Admin only)",
    "confession-settings": "View confession settings (Admin only)",
    "confession-stats": "View confession statistics",
    
    # News commands (if enabled)
    "news-latest": "Get the latest news items",
    "news-sources": "List all configured news sources",
    "news-status": "Show current news configuration and status",
    "news-help": "Show help for news commands",
    "news-set-channel": "Set channel for automatic news updates (Admin only)",
    "news-start": "Start automatic news updates (Admin only)",
    "news-stop": "Stop automatic news updates (Admin only)",
    "news-add": "Add a news source (Admin only)",
    "news-remove": "Remove a news source (Admin only)"
}

# ============================
# STATUS MESSAGES
# ============================

STATUS_MESSAGES = {
    "online": "🎵 Ready to play music!",
    "idle": "😴 Waiting for commands...",
    "playing": "🎵 Playing music",
    "maintenance": "🔧 Under maintenance",
}

# ============================
# LOG MESSAGES (for developers)
# ============================

LOG_MSG = {
    # Bot lifecycle
    "bot_starting": "Starting {bot_name} v{version}...",
    "bot_ready": "Bot is ready! Logged in as {user}",
    "bot_disconnected": "Bot disconnected",
    "cog_loaded": "Loaded cog: {cog_name}",
    "cog_failed": "Failed to load cog {cog_name}: {error}",
    
    # Music core
    "music_added_queue": "Added to queue: {title} in guild {guild_id}",
    "music_now_playing": "Now playing: {title} in guild {guild_id} from {source}",
    "music_search_start": "Searching for: {query}",
    "music_search_results": "Found {count} results for '{query}' from {source}",
    "music_search_timeout": "Search timeout for query: {query}",
    "music_conversion_success": "Successfully converted to playable: {title}",
    "music_conversion_failed": "Failed to convert to playable: {title}",
    
    # Voice channel
    "voice_joined": "Joined voice channel: {channel} in guild {guild_id}",
    "voice_left": "Left voice channel in guild {guild_id}",
    "idle_timer_started": "Started idle timer for guild {guild_id}",
    "idle_disconnect": "Auto-disconnected from guild {guild_id} due to inactivity",
    
    # Source management
    "source_search_start": "Starting search for: {query}",
    "source_determined": "Determined source: {source}",
    "source_search_results": "Found {count} results for '{query}' from {source}",
    "source_no_results": "No results found for '{query}' from {source}",
    "source_search_error": "Search error for '{query}' from {source}: {error}",
    "source_convert_none": "Cannot convert None result to playable",
    "source_convert_start": "Converting {title} from {source} to playable",
    "source_youtube_playable": "YouTube result {title} is already playable",
    "source_spotify_converted": "Successfully converted Spotify track: {title}",
    "source_spotify_failed": "Failed to convert Spotify track: {title}",
    "source_using_original": "Using original result for: {title}",
    "source_convert_error": "Error converting {title} to playable: {error}",
    
    # Spotify specific
    "spotify_no_credentials": "Spotify credentials not found - Spotify search disabled",
    "spotify_not_initialized": "Spotify adapter not initialized",
    "spotify_search_timeout": "Spotify search timeout for: {query}",
    "spotify_no_results": "No Spotify results for: {query}",
    "spotify_results_found": "Found {count} Spotify results for: {query}",
    "spotify_search_error": "Spotify search error for '{query}': {error}",
    "spotify_track_timeout": "Spotify track timeout for ID: {track_id}",
    "spotify_track_error": "Spotify track error for {url}: {error}",
    "spotify_album_error": "Spotify album error: {error}",
    "spotify_playlist_error": "Spotify playlist error: {error}",
    "spotify_artist_error": "Spotify artist error: {error}",
    "spotify_convert_none": "Cannot convert None Spotify result",
    "spotify_converting": "Converting Spotify track to YouTube: {title}",
    "spotify_youtube_trying": "Trying YouTube search with: {query}",
    "spotify_youtube_found": "Found YouTube match: {title}",
    "spotify_youtube_no_results": "No YouTube results for: {query}",
    "spotify_youtube_failed": "YouTube search failed for '{query}': {error}",
    "spotify_no_youtube": "No YouTube match found for Spotify track: {title}",
    "spotify_conversion_success": "Successfully converted to playable: {title}",
    "spotify_convert_error": "Error converting Spotify result: {error}",
    
    # SoundCloud specific
    "soundcloud_using_ytdlp": "Using yt-dlp for SoundCloud search functionality",
    "soundcloud_search_error": "SoundCloud search error: {error}",
    "soundcloud_fallback_error": "SoundCloud fallback search error: {error}",
    "soundcloud_convert_error": "SoundCloud conversion error: {error}",
    
    # YouTube specific
    "youtube_search_error": "YouTube search error for '{query}': {error}",
    
    # General
    "command_used": "Command {command} used by {user} in guild {guild_id}",
    "error_occurred": "Error in {location}: {error}",
}

# ============================
# NEWS MESSAGES
# ============================

MSG_NEWS = {
    "breaking_news": "📰 **Breaking News!**",
    "no_items_available": "No news items available yet. Please try again later.",
    "no_sources_configured": "No news sources configured.",
    "source_already_exists": "News source already exists: {name}",
    "source_added": "Added news source: {name}",
    "source_removed": "Removed news source: {name}",
    "source_not_found": "News source not found: {name}",
    "channel_set": "News updates will be sent to {channel}. Auto-posting is now enabled.",
    "set_channel_first": "Please set a news channel first using `/news set-channel`.",
    "auto_post_started": "Automatic news updates have been started!",
    "auto_post_stopped": "Automatic news updates have been stopped.",
    "auto_post_not_enabled": "Automatic news updates are not currently enabled.",
    "specify_subcommand": "Please specify a news subcommand. Use `/news help` for more information.",
    "help": {
        "title": "News Commands",
        "description": "Commands for the news feature",
        "latest": "Get the latest news items. Optionally specify how many items to show (default: 5).",
        "sources": "List all configured news sources.",
        "status": "Show current news configuration and status.",
        "set_channel": "Set the channel for automatic news updates and enable auto-posting. If no channel is specified, the current channel is used.",
        "start": "Start automatic news updates (admin only).",
        "stop": "Stop automatic news updates (admin only).",
        "add": "Add a news source (admin only).",
        "remove": "Remove a news source (admin only).",
    },
    "status": {
        "title": "News Configuration Status",
        "channel": "News Channel",
        "auto_posting": "Auto-posting",
        "service": "News Service",
        "sources": "News Sources",
        "available": "Available News",
        "not_set": "Not set",
        "enabled": "Enabled",
        "disabled": "Disabled",
        "running": "Running",
        "stopped": "Stopped",
        "sources_count": "{count} configured",
        "items_count": "{count} items",
        "channel_not_found": "Channel not found",
    },
    "sources": {
        "title": "News Sources",
        "description": "List of configured news sources",
    }
}

# ============================
# CONFESSION MESSAGES
# ============================

MSG_CONFESSION = {
    "confession_submitted": "✅ Your confession has been submitted anonymously! (ID: `{confession_id}`)",
    "reply_submitted": "✅ Your reply has been posted anonymously!",
    "confession_not_found": "❌ No confession found with ID `{confession_id}` in this server.",
    "channel_not_set": "❌ Confession channel is not set up for this server. Please ask an admin to set it up.",
    "content_too_long": "❌ Content too long! Maximum {max_length} characters allowed.",
    "on_cooldown": "❌ You're on cooldown! Please wait {time} before submitting another confession.",
    "channel_set": "✅ Anonymous confessions will now be posted to {channel}",
    "no_content": "❌ Please provide some content for your confession.",
    "no_confessions": "📊 No confessions found for this server.",
    "image_too_large": "❌ Image too large! Please use an image smaller than 8MB.",
    "invalid_image": "❌ Please attach a valid image file (PNG, JPG, GIF, etc.)",
    "help": {
        "title": "📝 Anonymous Confession System",
        "description": "Submit and reply to anonymous confessions safely with optional image attachments",
        "confess": "Submit an anonymous confession (with optional image)",
        "reply": "Reply to a confession using its ID (with optional image)",
        "setup": "Set up confession channel (Admin only)",
        "settings": "View confession settings (Admin only)",
        "stats": "View confession statistics"
    }
}
