package config

// Emoji constants used as semantic status indicators in Discord messages
const (
	// Music controls
	EmojiPlay = ">>" // play indicator
	EmojiPause = "||" // pause indicator
	EmojiStop = "[]" // stop indicator
	EmojiSkip = ">>|" // skip indicator

	// Status
	EmojiSuccess = "[OK]"
	EmojiError = "[ERR]"
	EmojiWarning = "[WARN]"
	EmojiMusic = "[Music]"
)

// Color constants for Discord embeds (hex values)
const (
	ColorPrimary    = 0x0099FF
	ColorSecondary  = 0x6C757D
	ColorSuccess    = 0x00FF00
	ColorError      = 0xFF0000
	ColorWarning    = 0xFFA500
	ColorInfo       = 0x0099FF
	ColorMusic      = 0x9932CC
	ColorSpotify    = 0x1DB954
	ColorYouTube    = 0xFF0000
	ColorSoundCloud = 0xFF7700
)

// Command permission levels
const (
	PermissionEveryone = 0
	PermissionDJ       = 1
	PermissionModerator = 2
	PermissionAdmin    = 3
	PermissionOwner    = 4
)

// Discord limits and constraints
const (
	MaxEmbedFields       = 25
	MaxEmbedDescription  = 4096
	MaxEmbedTitle        = 256
	MaxEmbedFieldName    = 256
	MaxEmbedFieldValue   = 1024
	MaxMessageLength     = 2000
	MaxEmbedFooter       = 2048
	MaxEmbedAuthor       = 256
)

// Application constants
const (
	AppName        = "NeruBot"
	AppVersion     = "4.0.0"
	AppDescription = "Your friendly Discord companion!"
	AppAuthor      = "nerufuyo"
	AppWebsite     = "https://github.com/nerufuyo/nerubot"
	AppRepository  = "https://github.com/nerufuyo/nerubot"
)

// Data file paths
const (
	DataDir                = "data"
	ConfessionDir          = "data/confessions"
	RoastDir               = "data/roasts"
	ConfessionsFile        = "data/confessions/confessions.json"
	RepliesFile            = "data/confessions/replies.json"
	ConfessionSettingsFile = "data/confessions/settings.json"
	ConfessionQueueFile    = "data/confessions/queue.json"
	RoastProfilesFile      = "data/roasts/profiles.json"
	RoastActivitiesFile    = "data/roasts/activities.json"
	RoastStatsFile         = "data/roasts/stats.json"
	RoastPatternsFile      = "data/roasts/patterns.json"
)
