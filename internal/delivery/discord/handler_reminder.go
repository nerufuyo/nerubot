package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// handleReminder shows upcoming Indonesian holidays and today's Ramadan schedule.
func (b *Bot) handleReminder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.reminderService == nil {
		b.respondError(s, i, "Reminder service is not available")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Your Reminders~",
		Description: "I've got everything lined up for you, don't worry about a thing!",
		Color:       config.ColorPrimary,
		Timestamp:   time.Now().Format(time.RFC3339),
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	// Upcoming holidays
	holidays := b.reminderService.GetUpcomingHolidays(5)
	if len(holidays) > 0 {
		var lines []string
		for _, h := range holidays {
			lines = append(lines, fmt.Sprintf("**%s** — %s", h.Date.Format("2 Jan 2006"), h.Name))
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Upcoming Holidays~ (days off with me!)",
			Value:  strings.Join(lines, "\n"),
			Inline: false,
		})
	}

	// Ramadan schedule for today
	schedule := b.reminderService.GetTodayRamadanSchedule()
	if schedule != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Today's Ramadan Schedule (WIB) ~stay strong!",
			Value: fmt.Sprintf("Sahoor: **%s** (I'll wake you up~)\nIftar: **%s** (you earned it!)",
				schedule.SahoorTime.Format("15:04"),
				schedule.BerbukaTime.Format("15:04"),
			),
			Inline: false,
		})
	}

	// Work schedule for today
	if b.reminderService.IsRamadanToday() {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Today's Work Hours (Ramadan)",
			Value:  "Start: **08:00 WIB**\nEnd: **16:00 WIB**\n\nDon't overdo it, okay? I need you healthy~",
			Inline: false,
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Today's Work Hours",
			Value:  "Start: **09:00 WIB**\nEnd: **17:30 WIB**\n\nWork hard but don't forget to take breaks for me~",
			Inline: false,
		})
	}

	// Daily reminders schedule
	now := time.Now().In(time.FixedZone("WIB", 7*60*60))
	weekday := now.Weekday()
	isWorkday := weekday >= time.Monday && weekday <= time.Friday

	var dailyLines []string
	if isWorkday {
		if b.reminderService.IsRamadanToday() {
			dailyLines = append(dailyLines, "Standup: **09:00 WIB** (Ramadan schedule)")
		} else {
			dailyLines = append(dailyLines, "Standup: **09:30 WIB**")
		}
		dailyLines = append(dailyLines, "Lunch Break: **12:00 WIB**")
		dailyLines = append(dailyLines, "Love Note: **11:00 & 15:00 WIB** (random surprise~)")
	}
	if weekday == time.Friday {
		dailyLines = append(dailyLines, "Friday Prayer: **11:30 WIB**")
	}
	if len(dailyLines) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Daily Reminders~",
			Value:  strings.Join(dailyLines, "\n"),
			Inline: false,
		})
	}

	if len(embed.Fields) == 0 {
		b.followUp(s, i, "No upcoming reminders at the moment.")
		return
	}

	b.followUpEmbed(s, i, embed)
}

// handleReminderSet allows admins to set the reminder channel via /reminder-set.
func (b *Bot) handleReminderSet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.reminderService == nil {
		b.respondError(s, i, "Reminder service is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please specify a channel.")
		return
	}

	channel := options[0].ChannelValue(s)
	if channel == nil {
		b.respondError(s, i, "Invalid channel.")
		return
	}

	// Only allow text channels
	if channel.Type != discordgo.ChannelTypeGuildText {
		b.respondError(s, i, "Please select a text channel.")
		return
	}

	b.reminderService.SetChannelID(channel.ID)
	b.respond(s, i, fmt.Sprintf("Reminders will now be sent to <#%s>.", channel.ID))
}
