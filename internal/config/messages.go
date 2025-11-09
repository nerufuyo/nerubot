package config

// Messages holds all user-facing message templates
type Messages struct {
	Info    map[string]string
	Error   map[string]string
	Success map[string]string
	Help    map[string]string
}

// DefaultMessages returns the default message templates
func DefaultMessages() *Messages {
	return &Messages{
		Info: map[string]string{
			// Bot status
			"bot_starting":    "🚀 Starting %s v%s...",
			"bot_ready":       "✅ %s is ready!",
			"bot_connecting":  "🔌 Connecting to Discord...",
			"bot_disconnected": "🔌 Disconnected from Discord",
			
			// Music
			"now_playing":     "🎵 Now playing: **%s** by **%s**",
			"added_to_queue":  "➕ Added to queue: **%s**",
			"queue_position":  "📋 Position in queue: #%d",
			"queue_empty":     "📋 Queue is empty",
			"joined_voice":    "🔊 Joined voice channel: %s",
			"left_voice":      "👋 Left voice channel",
			
			// Confession
			"confession_submitted": "✅ Confession submitted successfully! ID: #%d",
			"confession_posted":    "📝 New confession posted!",
			"reply_posted":         "💬 Reply posted successfully!",
			
			// General
			"processing":      "⏱️ Processing...",
			"please_wait":     "⏱️ Please wait...",
		},
		
		Error: map[string]string{
			// General
			"unknown_error":        "❌ An unknown error occurred",
			"invalid_input":        "❌ Invalid input provided",
			"unauthorized":         "❌ You don't have permission to do that",
			"not_found":            "❌ Not found",
			
			// Music
			"no_results":           "❌ No results found for: %s",
			"search_timeout":       "⏱️ Search timed out. Please try again.",
			"conversion_timeout":   "⏱️ Audio conversion timed out",
			"not_in_voice":         "❌ You must be in a voice channel!",
			"different_voice":      "❌ You must be in the same voice channel as the bot",
			"queue_full":           "❌ Queue is full! Maximum: %d songs",
			"song_too_long":        "❌ Song is too long! Maximum: %s",
			"playback_error":       "❌ Error during playback: %s",
			"no_permission_voice":  "❌ I don't have permission to join that voice channel",
			
			// Confession
			"confession_not_found": "❌ Confession #%d not found",
			"cooldown_active":      "⏱️ Please wait before submitting again (%s remaining)",
			"content_too_long":     "❌ Content is too long! Maximum: %d characters",
			"content_empty":        "❌ Content cannot be empty",
			
			// AI Chatbot
			"ai_unavailable":       "❌ AI service is currently unavailable",
			"ai_error":             "❌ AI service error: %s",
			"all_providers_failed": "❌ All AI providers failed. Please try again later.",
			
			// Crypto
			"whale_api_error":      "❌ Error fetching whale transactions",
			"twitter_api_error":    "❌ Error fetching tweets",
		},
		
		Success: map[string]string{
			"command_executed":     "✅ Command executed successfully",
			"settings_updated":     "✅ Settings updated",
			"channel_set":          "✅ Channel set to %s",
			"feature_enabled":      "✅ Feature enabled",
			"feature_disabled":     "✅ Feature disabled",
		},
		
		Help: map[string]string{
			"main_title":           "📚 %s Help",
			"main_description":     "Select a category to view commands",
			"command_usage":        "**Usage:** %s",
			"command_description":  "%s",
			"no_permission":        "🔒 Requires: %s",
			
			// Command descriptions
			"play":                 "Play a song or add it to the queue",
			"queue":                "Show the current music queue",
			"skip":                 "Skip the current song",
			"stop":                 "Stop playback and clear the queue",
			"pause":                "Pause the current song",
			"resume":               "Resume playback",
			"nowplaying":           "Show currently playing song",
			"loop":                 "Set loop mode (off/single/queue)",
			"shuffle":              "Shuffle the queue",
			"247":                  "Toggle 24/7 mode",
			
			"confess":              "Submit an anonymous confession",
			"reply":                "Reply to a confession anonymously",
			"confession_setup":     "Set up confession channel (Admin)",
			"confession_stats":     "View confession statistics",
			
			"chat":                 "Chat with the AI bot",
			"reset_chat":           "Reset your chat session",
			
			"roast":                "Get roasted based on your Discord activity",
			"roast_stats":          "View roast statistics",
			"behavior_analysis":    "Analyze Discord behavior patterns",
			
			"news":                 "Get latest news updates",
			"whale":                "View crypto whale transactions",
			"guru":                 "View crypto guru tweets",
			
			"help":                 "Show this help message",
			"about":                "Show bot information",
			"features":             "Show all available features",
			"commands":             "Show quick command reference",
		},
	}
}

// LogMessages holds log message templates
var LogMessages = map[string]string{
	"bot_starting":      "Starting %s v%s",
	"bot_ready":         "Bot ready as: %s",
	"cog_loaded":        "Loaded cog: %s",
	"cog_failed":        "Failed to load cog %s: %v",
	"command_executed":  "Command executed: %s by %s",
	"error_occurred":    "Error in %s: %v",
	"connecting_db":     "Connecting to database...",
	"db_connected":      "Database connected",
	"shutdown":          "Shutting down...",
}
