package discord

import "strings"

// Contains keywords to block for SARA and porn topics
var blockedKeywords = []string{
	// SARA (ethnic, religious, racial, intergroup)
	"agama", "ras", "suku", "etnis", "rasis", "rasisme", "agama tertentu", "provokasi", "intoleran", "diskriminasi",
	// Porn/NSFW
	"sex", "seks", "porno", "porn", "bokep", "mesum", "bugil", "nude", "masturbasi", "orgasme", "payudara", "vagina", "penis", "kontol", "memek", "ngentot", "jilmek", "jav", "hentai", "bdsm", "fetish", "anal", "cumshot", "blowjob", "handjob", "tits", "boobs", "pussy", "dick", "cock", "ejakulasi", "sperma",
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
