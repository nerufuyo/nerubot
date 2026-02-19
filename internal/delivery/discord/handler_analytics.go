package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// handleStats handles showing server statistics
func (b *Bot) handleStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.analyticsService == nil {
		b.respondError(s, i, "Analytics service is not available")
		return
	}

	// Defer response
	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Get server stats
	stats, err := b.analyticsService.GetServerStats(i.GuildID)
	if err != nil {
		b.followUpError(s, i, "No statistics available yet. Start using commands to generate stats!")
		return
	}

	// Get most used command and most active user
	topCmd, topCmdCount := stats.GetMostUsedCommand()
	topUser, topUserCount := stats.GetMostActiveUser()

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title:     "Server Statistics",
		Color:     0x00ff00,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Total Commands",
				Value:  fmt.Sprintf("%d", stats.CommandsUsed),
				Inline: true,
			},
			{
				Name:   "Songs Played",
				Value:  fmt.Sprintf("%d", stats.SongsPlayed),
				Inline: true,
			},
			{
				Name:   "Confessions",
				Value:  fmt.Sprintf("%d", stats.ConfessionsTotal),
				Inline: true,
			},
			{
				Name:   "Roasts Generated",
				Value:  fmt.Sprintf("%d", stats.RoastsGenerated),
				Inline: true,
			},
			{
				Name:   "Chat Messages",
				Value:  fmt.Sprintf("%d", stats.ChatMessages),
				Inline: true,
			},
			{
				Name:   "News Requests",
				Value:  fmt.Sprintf("%d", stats.NewsRequests),
				Inline: true,
			},
		},
	}

	if topCmd != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Most Used Command",
			Value:  fmt.Sprintf("/%s (%d times)", topCmd, topCmdCount),
			Inline: false,
		})
	}

	if topUser != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Most Active User",
			Value:  fmt.Sprintf("<@%s> (%d commands)", topUser, topUserCount),
			Inline: false,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Server: %s | Active since %s", stats.GuildName, stats.FirstSeen.Format("Jan 2, 2006")),
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

// handleProfile handles showing user statistics
func (b *Bot) handleProfile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.analyticsService == nil {
		b.respondError(s, i, "Analytics service is not available")
		return
	}

	// Get target user (default to command user)
	targetUser := i.Member.User
	options := i.ApplicationCommandData().Options
	if len(options) > 0 && options[0].Type == discordgo.ApplicationCommandOptionUser {
		targetUser = options[0].UserValue(s)
	}

	// Defer response
	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Get user stats
	stats, err := b.analyticsService.GetUserStats(targetUser.ID)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("No statistics available for <@%s> yet.", targetUser.ID))
		return
	}

	// Get favorite command
	favCmd, favCount := stats.GetFavoriteCommand()

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("User Profile: %s", targetUser.Username),
		Color:     0x0099ff,
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: targetUser.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Total Commands",
				Value:  fmt.Sprintf("%d", stats.CommandsUsed),
				Inline: true,
			},
			{
				Name:   "Songs Requested",
				Value:  fmt.Sprintf("%d", stats.SongsRequested),
				Inline: true,
			},
			{
				Name:   "Confessions Posted",
				Value:  fmt.Sprintf("%d", stats.ConfessionsPosted),
				Inline: true,
			},
			{
				Name:   "Roasts Received",
				Value:  fmt.Sprintf("%d", stats.RoastsReceived),
				Inline: true,
			},
			{
				Name:   "Chat Messages",
				Value:  fmt.Sprintf("%d", stats.ChatMessages),
				Inline: true,
			},
			{
				Name:   "News Checks",
				Value:  fmt.Sprintf("%d", stats.NewsRequests),
				Inline: true,
			},
		},
	}

	if favCmd != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Favorite Command",
			Value:  fmt.Sprintf("/%s (%d times)", favCmd, favCount),
			Inline: false,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Active since %s", stats.FirstSeen.Format("Jan 2, 2006")),
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}
