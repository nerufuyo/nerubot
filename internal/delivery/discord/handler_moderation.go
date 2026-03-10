package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// handleKick handles the /kick command (admin/mod only).
func (b *Bot) handleKick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User
	reason := "No reason provided"

	for _, opt := range options {
		switch opt.Name {
		case "user":
			targetUser = opt.UserValue(s)
		case "reason":
			reason = opt.StringValue()
		}
	}

	if targetUser == nil {
		b.respondError(s, i, "Please specify a user to kick.")
		return
	}

	// Prevent kicking self or the bot
	if targetUser.ID == i.Member.User.ID {
		b.respondError(s, i, "You can't kick yourself!")
		return
	}
	if targetUser.ID == s.State.User.ID {
		b.respondError(s, i, "I can't kick myself!")
		return
	}

	err := s.GuildMemberDeleteWithReason(i.GuildID, targetUser.ID, reason)
	if err != nil {
		b.respondError(s, i, "Failed to kick user: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "👢 User Kicked",
		Color: config.ColorWarning,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "User", Value: targetUser.Username, Inline: true},
			{Name: "Moderator", Value: i.Member.User.Username, Inline: true},
			{Name: "Reason", Value: reason, Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handleBan handles the /ban command (admin/mod only).
func (b *Bot) handleBan(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User
	reason := "No reason provided"
	deleteDays := 0

	for _, opt := range options {
		switch opt.Name {
		case "user":
			targetUser = opt.UserValue(s)
		case "reason":
			reason = opt.StringValue()
		case "delete_days":
			deleteDays = int(opt.IntValue())
			if deleteDays > 7 {
				deleteDays = 7
			}
		}
	}

	if targetUser == nil {
		b.respondError(s, i, "Please specify a user to ban.")
		return
	}

	if targetUser.ID == i.Member.User.ID {
		b.respondError(s, i, "You can't ban yourself!")
		return
	}
	if targetUser.ID == s.State.User.ID {
		b.respondError(s, i, "I can't ban myself!")
		return
	}

	err := s.GuildBanCreateWithReason(i.GuildID, targetUser.ID, reason, deleteDays)
	if err != nil {
		b.respondError(s, i, "Failed to ban user: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "🔨 User Banned",
		Color: config.ColorError,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "User", Value: targetUser.Username, Inline: true},
			{Name: "Moderator", Value: i.Member.User.Username, Inline: true},
			{Name: "Reason", Value: reason, Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handleTimeout handles the /timeout command (admin/mod only).
func (b *Bot) handleTimeout(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User
	reason := "No reason provided"
	duration := int64(5) // default 5 minutes

	for _, opt := range options {
		switch opt.Name {
		case "user":
			targetUser = opt.UserValue(s)
		case "duration":
			duration = opt.IntValue()
		case "reason":
			reason = opt.StringValue()
		}
	}

	if targetUser == nil {
		b.respondError(s, i, "Please specify a user to timeout.")
		return
	}

	if targetUser.ID == i.Member.User.ID {
		b.respondError(s, i, "You can't timeout yourself!")
		return
	}

	// Calculate timeout end time
	timeoutUntil := time.Now().Add(time.Duration(duration) * time.Minute)

	_, err := s.GuildMemberEdit(i.GuildID, targetUser.ID, &discordgo.GuildMemberParams{
		CommunicationDisabledUntil: &timeoutUntil,
	})
	if err != nil {
		b.respondError(s, i, "Failed to timeout user: "+err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "🔇 User Timed Out",
		Color: config.ColorWarning,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "User", Value: targetUser.Username, Inline: true},
			{Name: "Duration", Value: fmt.Sprintf("%d minutes", duration), Inline: true},
			{Name: "Moderator", Value: i.Member.User.Username, Inline: true},
			{Name: "Reason", Value: reason, Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handlePurge handles the /purge command - delete multiple messages (admin/mod only).
func (b *Bot) handlePurge(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	amount := int64(10) // default

	for _, opt := range options {
		if opt.Name == "amount" {
			amount = opt.IntValue()
		}
	}

	if amount < 1 || amount > 100 {
		b.respondError(s, i, "Amount must be between 1 and 100.")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Fetch messages
	messages, err := s.ChannelMessages(i.ChannelID, int(amount), "", "", "")
	if err != nil {
		b.followUpError(s, i, "Failed to fetch messages: "+err.Error())
		return
	}

	if len(messages) == 0 {
		b.followUp(s, i, "No messages to delete.")
		return
	}

	// Collect message IDs (Discord only allows bulk delete for messages < 14 days old)
	var messageIDs []string
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.ID)
	}

	err = s.ChannelMessagesBulkDelete(i.ChannelID, messageIDs)
	if err != nil {
		b.followUpError(s, i, "Failed to delete messages: "+err.Error())
		return
	}

	b.followUp(s, i, fmt.Sprintf("🗑️ Successfully deleted **%d** messages.", len(messageIDs)))
}

var moderationRepo = repository.NewModerationRepository()

// handleWarn handles the /warn command (admin/mod only).
func (b *Bot) handleWarn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User
	reason := "No reason provided"

	for _, opt := range options {
		switch opt.Name {
		case "user":
			targetUser = opt.UserValue(s)
		case "reason":
			reason = opt.StringValue()
		}
	}

	if targetUser == nil {
		b.respondError(s, i, "Please specify a user to warn.")
		return
	}

	warning := &entity.Warning{
		GuildID:  i.GuildID,
		UserID:   targetUser.ID,
		Username: targetUser.Username,
		Reason:   reason,
		IssuedBy: i.Member.User.Username,
		IssuedAt: time.Now(),
	}

	if err := moderationRepo.AddWarning(warning); err != nil {
		b.respondError(s, i, "Failed to save warning: "+err.Error())
		return
	}

	// Get total warnings
	warnings, _ := moderationRepo.GetWarnings(i.GuildID, targetUser.ID)
	totalWarnings := len(warnings)

	embed := &discordgo.MessageEmbed{
		Title: "⚠️ User Warned",
		Color: config.ColorWarning,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "User", Value: targetUser.Username, Inline: true},
			{Name: "Moderator", Value: i.Member.User.Username, Inline: true},
			{Name: "Total Warnings", Value: fmt.Sprintf("%d", totalWarnings), Inline: true},
			{Name: "Reason", Value: reason, Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handleWarnings handles the /warnings command - view warnings for a user.
func (b *Bot) handleWarnings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User

	for _, opt := range options {
		if opt.Name == "user" {
			targetUser = opt.UserValue(s)
		}
	}

	if targetUser == nil {
		b.respondError(s, i, "Please specify a user.")
		return
	}

	warnings, err := moderationRepo.GetWarnings(i.GuildID, targetUser.ID)
	if err != nil {
		b.respondError(s, i, "Failed to get warnings: "+err.Error())
		return
	}

	if len(warnings) == 0 {
		b.respond(s, i, fmt.Sprintf("✅ **%s** has no warnings.", targetUser.Username))
		return
	}

	desc := ""
	for idx, w := range warnings {
		desc += fmt.Sprintf("**%d.** %s\n   By: %s | <t:%d:R>\n", idx+1, w.Reason, w.IssuedBy, w.IssuedAt.Unix())
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("⚠️ Warnings for %s (%d total)", targetUser.Username, len(warnings)),
		Description: desc,
		Color:       config.ColorWarning,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handleClearWarnings handles the /clearwarnings command - clear warnings for a user (admin only).
func (b *Bot) handleClearWarnings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var targetUser *discordgo.User

	for _, opt := range options {
		if opt.Name == "user" {
			targetUser = opt.UserValue(s)
		}
	}

	if targetUser == nil {
		b.respondError(s, i, "Please specify a user.")
		return
	}

	if err := moderationRepo.ClearWarnings(i.GuildID, targetUser.ID); err != nil {
		b.respondError(s, i, "Failed to clear warnings: "+err.Error())
		return
	}

	b.respond(s, i, fmt.Sprintf("✅ Cleared all warnings for **%s**.", targetUser.Username))
}
