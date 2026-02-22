package config

// Messages holds all user-facing message templates.
type Messages struct {
	Info    map[string]string
	Error   map[string]string
	Success map[string]string
	Help    map[string]string
}

// DefaultMessages returns the default message templates.
func DefaultMessages() *Messages {
	return &Messages{
		Info: map[string]string{
			"bot_starting":    "Starting %s v%s...",
			"bot_ready":       "%s is ready!",
			"bot_connecting":  "Connecting to Discord...",
			"bot_disconnected": "Disconnected from Discord",

			"confession_submitted": "Confession submitted. ID: #%d",
			"confession_posted":    "New confession posted.",
			"reply_posted":         "Reply posted.",

			"processing":      "Processing...",
			"please_wait":     "Please wait...",
		},

		Error: map[string]string{
			"unknown_error":        "An unknown error occurred",
			"invalid_input":        "Invalid input provided",
			"unauthorized":         "You don't have permission to do that",
			"not_found":            "Not found",

			"no_results":           "No results found for: %s",

			"confession_not_found": "Confession #%d not found",
			"cooldown_active":      "Please wait before submitting again (%s remaining)",
			"content_too_long":     "Content is too long! Maximum: %d characters",
			"content_empty":        "Content cannot be empty",

			"ai_unavailable":       "AI service is currently unavailable",
			"ai_error":             "AI service error: %s",
			"all_providers_failed": "All AI providers failed. Please try again later.",

			"whale_api_error":      "Error fetching whale transactions",
		},

		Success: map[string]string{
			"command_executed":     "Command executed successfully",
			"settings_updated":     "Settings updated",
			"channel_set":          "Channel set to %s",
			"feature_enabled":      "Feature enabled",
			"feature_disabled":     "Feature disabled",
		},

		Help: map[string]string{
			"main_title":           "%s Help",
			"main_description":     "Select a category to view commands",
			"command_usage":        "**Usage:** %s",
			"command_description":  "%s",
			"no_permission":        "Requires: %s",

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

			"help":                 "Show this help message",
			"about":                "Show bot information",
			"features":             "Show all available features",
			"commands":             "Show quick command reference",
		},
	}
}

// LogMessages holds log message templates.
var LogMessages = map[string]string{
	"bot_starting":      "Starting %s v%s",
	"bot_ready":         "Bot ready as: %s",
	"cog_loaded":        "Loaded cog: %s",
	"cog_failed":        "Failed to load cog %s: %v",
	"command_executed":  "Command executed: %s by %s",
	"error_occurred":    "Error in %s: %v",
	"shutdown":          "Shutting down...",
}
