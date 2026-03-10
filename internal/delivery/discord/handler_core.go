package discord

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// handlePing handles the /ping command - shows bot latency.
func (b *Bot) handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	start := time.Now()
	if err := b.deferResponse(s, i); err != nil {
		return
	}
	latency := time.Since(start).Milliseconds()
	wsLatency := s.HeartbeatLatency().Milliseconds()

	embed := &discordgo.MessageEmbed{
		Title: "🏓 Pong!",
		Color: config.ColorSuccess,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Bot Latency", Value: fmt.Sprintf("%dms", latency), Inline: true},
			{Name: "WebSocket", Value: fmt.Sprintf("%dms", wsLatency), Inline: true},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.followUpEmbed(s, i, embed)
}

// handleBotInfo handles the /botinfo command - shows bot information.
func (b *Bot) handleBotInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	uptime := time.Since(b.startedAt)
	uptimeStr := formatDuration(uptime)

	guildCount := len(s.State.Guilds)
	var totalMembers int
	for _, g := range s.State.Guilds {
		totalMembers += g.MemberCount
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("ℹ️ %s Info", b.config.Bot.Name),
		Description: b.config.Bot.Description,
		Color:       config.ColorPrimary,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Version", Value: b.config.Bot.Version, Inline: true},
			{Name: "Developer", Value: b.config.Bot.Author, Inline: true},
			{Name: "Language", Value: fmt.Sprintf("Go %s", runtime.Version()), Inline: true},
			{Name: "Uptime", Value: uptimeStr, Inline: true},
			{Name: "Servers", Value: fmt.Sprintf("%d", guildCount), Inline: true},
			{Name: "Users", Value: fmt.Sprintf("%d", totalMembers), Inline: true},
			{Name: "Memory", Value: fmt.Sprintf("%.1f MB", float64(memStats.Alloc)/1024/1024), Inline: true},
			{Name: "Library", Value: "discordgo", Inline: true},
			{Name: "Website", Value: b.config.Bot.Website, Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s v%s | %s", b.config.Bot.Name, b.config.Bot.Version, b.config.Bot.Author),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handleServerInfo handles the /serverinfo command - shows server details.
func (b *Bot) handleServerInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		b.respondError(s, i, "Failed to get server info")
		return
	}

	// Count channels by type
	var textChannels, voiceChannels, categories int
	channels, _ := s.GuildChannels(guild.ID)
	for _, ch := range channels {
		switch ch.Type {
		case discordgo.ChannelTypeGuildText:
			textChannels++
		case discordgo.ChannelTypeGuildVoice:
			voiceChannels++
		case discordgo.ChannelTypeGuildCategory:
			categories++
		}
	}

	// Count roles
	roleCount := len(guild.Roles)

	// Count emojis
	emojiCount := len(guild.Emojis)

	// Owner info
	owner, _ := s.User(guild.OwnerID)
	ownerStr := guild.OwnerID
	if owner != nil {
		ownerStr = owner.Username
	}

	// Server icon
	iconURL := ""
	if guild.Icon != "" {
		iconURL = fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png?size=256", guild.ID, guild.Icon)
	}

	// Boost info
	boostLevel := fmt.Sprintf("Level %d", guild.PremiumTier)
	boostCount := fmt.Sprintf("%d boosts", guild.PremiumSubscriptionCount)

	embed := &discordgo.MessageEmbed{
		Title: "🏠 " + guild.Name,
		Color: config.ColorPrimary,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Owner", Value: ownerStr, Inline: true},
			{Name: "Members", Value: fmt.Sprintf("%d", guild.MemberCount), Inline: true},
			{Name: "Roles", Value: fmt.Sprintf("%d", roleCount), Inline: true},
			{Name: "Text Channels", Value: fmt.Sprintf("%d", textChannels), Inline: true},
			{Name: "Voice Channels", Value: fmt.Sprintf("%d", voiceChannels), Inline: true},
			{Name: "Categories", Value: fmt.Sprintf("%d", categories), Inline: true},
			{Name: "Emojis", Value: fmt.Sprintf("%d", emojiCount), Inline: true},
			{Name: "Boost", Value: fmt.Sprintf("%s (%s)", boostLevel, boostCount), Inline: true},
			{Name: "Created", Value: fmt.Sprintf("<t:%d:R>", snowflakeToTime(guild.ID).Unix()), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Server ID: " + guild.ID,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	if iconURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: iconURL}
	}

	b.respondEmbed(s, i, embed)
}

// handleUserInfo handles the /userinfo command - shows user information.
func (b *Bot) handleUserInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get target user (default: command invoker)
	targetUser := i.Member.User
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "user" {
			targetUser = opt.UserValue(s)
		}
	}

	// Get member info for this guild
	member, err := s.GuildMember(i.GuildID, targetUser.ID)
	if err != nil {
		b.respondError(s, i, "Failed to get user info")
		return
	}

	// Build roles string
	rolesStr := "None"
	if len(member.Roles) > 0 {
		roles := ""
		for idx, roleID := range member.Roles {
			if idx > 0 {
				roles += ", "
			}
			roles += fmt.Sprintf("<@&%s>", roleID)
			if idx >= 9 { // limit display to 10 roles
				roles += fmt.Sprintf(" +%d more", len(member.Roles)-10)
				break
			}
		}
		rolesStr = roles
	}

	// Account creation time from snowflake
	createdAt := snowflakeToTime(targetUser.ID)

	embed := &discordgo.MessageEmbed{
		Title: "👤 " + targetUser.Username,
		Color: config.ColorPrimary,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: targetUser.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Username", Value: targetUser.Username, Inline: true},
			{Name: "Display Name", Value: memberDisplayName(member, targetUser), Inline: true},
			{Name: "Bot", Value: fmt.Sprintf("%t", targetUser.Bot), Inline: true},
			{Name: "Account Created", Value: fmt.Sprintf("<t:%d:R>", createdAt.Unix()), Inline: true},
			{Name: "Joined Server", Value: fmt.Sprintf("<t:%d:R>", member.JoinedAt.Unix()), Inline: true},
			{Name: "Roles", Value: fmt.Sprintf("%d roles", len(member.Roles)), Inline: true},
			{Name: "Role List", Value: rolesStr, Inline: false},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "User ID: " + targetUser.ID,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handleAvatar handles the /avatar command - displays user's avatar.
func (b *Bot) handleAvatar(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get target user (default: command invoker)
	targetUser := i.Member.User
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "user" {
			targetUser = opt.UserValue(s)
		}
	}

	avatarURL := targetUser.AvatarURL("1024")

	embed := &discordgo.MessageEmbed{
		Title: "🖼️ " + targetUser.Username + "'s Avatar",
		Color: config.ColorPrimary,
		Image: &discordgo.MessageEmbedImage{
			URL: avatarURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Click the image to open full size",
		},
	}
	b.respondEmbed(s, i, embed)
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// snowflakeToTime converts a Discord snowflake ID string to a time.Time.
func snowflakeToTime(id string) time.Time {
	sfID, _ := strconv.ParseInt(id, 10, 64)
	// Discord epoch: 2015-01-01T00:00:00.000Z = 1420070400000 ms
	ms := (sfID >> 22) + 1420070400000
	return time.Unix(0, ms*int64(time.Millisecond))
}

// memberDisplayName returns the display name for a guild member.
func memberDisplayName(member *discordgo.Member, user *discordgo.User) string {
	if member.Nick != "" {
		return member.Nick
	}
	if user.GlobalName != "" {
		return user.GlobalName
	}
	return user.Username
}
