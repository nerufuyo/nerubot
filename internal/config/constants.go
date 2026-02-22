package config

// Emoji constants used as semantic status indicators in Discord messages
const (
	// Status
	EmojiSuccess = "[OK]"
	EmojiError = "[ERR]"
	EmojiWarning = "[WARN]"
)

// Color constants for Discord embeds (hex values)
const (
	ColorPrimary    = 0x0099FF
	ColorSecondary  = 0x6C757D
	ColorSuccess    = 0x00FF00
	ColorError      = 0xFF0000
	ColorWarning    = 0xFFA500
	ColorInfo       = 0x0099FF
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

// Language codes for multi-language support
const (
	LangEN = "EN" // English
	LangID = "ID" // Indonesian
	LangJP = "JP" // Japanese
	LangKR = "KR" // Korean
	LangZH = "ZH" // Chinese
)

// DefaultLang is the default language used when none is specified
const DefaultLang = LangEN

// SupportedLanguages is the list of all supported language codes
var SupportedLanguages = []string{LangEN, LangID, LangJP, LangKR, LangZH}

// LanguageNames maps language codes to their display names
var LanguageNames = map[string]string{
	LangEN: "English",
	LangID: "Bahasa Indonesia",
	LangJP: "日本語",
	LangKR: "한국어",
	LangZH: "中文",
}

// LanguagePromptInstruction returns an AI prompt instruction for the given language
func LanguagePromptInstruction(lang string) string {
	switch lang {
	case LangID:
		return "IMPORTANT: You MUST respond entirely in Bahasa Indonesia. Use natural, conversational Indonesian."
	case LangJP:
		return "IMPORTANT: You MUST respond entirely in Japanese (日本語). Use natural, conversational Japanese."
	case LangKR:
		return "IMPORTANT: You MUST respond entirely in Korean (한국어). Use natural, conversational Korean."
	case LangZH:
		return "IMPORTANT: You MUST respond entirely in Chinese (中文). Use natural, conversational Simplified Chinese."
	default:
		return "Respond in English."
	}
}

// LanguageAIName returns the full language name used in AI prompts for the given code.
func LanguageAIName(lang string) string {
	switch lang {
	case LangID:
		return "Indonesian"
	case LangJP:
		return "Japanese"
	case LangKR:
		return "Korean"
	case LangZH:
		return "Chinese"
	default:
		return "English"
	}
}
