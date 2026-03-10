package discord

import (
	"regexp"
	"strings"
)

// Contains keywords to block for SARA and porn topics
var blockedKeywords = []string{
	// SARA (ethnic, religious, racial, intergroup)
	"agama", "ras", "suku", "etnis", "rasis", "rasisme", "agama tertentu", "provokasi", "intoleran", "diskriminasi",
	// Porn/NSFW
	"sex", "seks", "porno", "porn", "bokep", "mesum", "bugil", "nude", "masturbasi", "orgasme", "payudara", "vagina", "penis", "kontol", "memek", "ngentot", "jilmek", "jav", "hentai", "bdsm", "fetish", "anal", "cumshot", "blowjob", "handjob", "tits", "boobs", "pussy", "dick", "cock", "ejakulasi", "sperma",
}

// Regex patterns for detecting malicious or suspicious links
var maliciousLinkPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)discord\.gift|discordapp\.gift`),                 // fake gift links
	regexp.MustCompile(`(?i)(?:free|steam|nitro)[\w-]*\.(?:ru|xyz|tk|ml)`),  // phishing domains
	regexp.MustCompile(`(?i)@everyone\s+https?://`),                         // spam @everyone with link
	regexp.MustCompile(`(?i)https?://\S*(?:grabify|iplogger|2no\.co)\S*`),   // IP loggers
}

// Checks if a message contains blocked keywords (case-insensitive)
func containsBlockedKeyword(msg string) bool {
	msgLower := strings.ToLower(msg)
	for _, word := range blockedKeywords {
		if strings.Contains(msgLower, word) {
			return true
		}
	}
	return false
}

// containsMaliciousLink checks if a message contains a suspicious or malicious link.
func containsMaliciousLink(msg string) bool {
	for _, pattern := range maliciousLinkPatterns {
		if pattern.MatchString(msg) {
			return true
		}
	}
	return false
}

// spamTracker tracks recent message timestamps per user for spam detection.
type spamTracker struct {
	// userMessages maps "guildID:userID" → list of unix timestamps
	userMessages map[string][]int64
}

var globalSpamTracker = &spamTracker{
	userMessages: make(map[string][]int64),
}

// recordMessage records a message timestamp for spam detection.
// Returns true if the user is spamming (more than maxMessages in windowSeconds).
func (st *spamTracker) recordMessage(guildID, userID string, timestamp int64, maxMessages int, windowSeconds int64) bool {
	key := guildID + ":" + userID

	// Clean old entries
	cutoff := timestamp - windowSeconds
	filtered := make([]int64, 0)
	for _, ts := range st.userMessages[key] {
		if ts > cutoff {
			filtered = append(filtered, ts)
		}
	}
	filtered = append(filtered, timestamp)
	st.userMessages[key] = filtered

	return len(filtered) > maxMessages
}

