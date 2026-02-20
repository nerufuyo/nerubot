package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// reminderStrings holds all translatable strings for the reminder command.
type reminderStrings struct {
	Title             string
	Description       string
	HolidaysHeader    string
	RamadanHeader     string
	SahoorLabel       string
	SahoorNote        string
	IftarLabel        string
	IftarNote         string
	WorkRamadanHeader string
	WorkRamadanBody   string
	WorkHeader        string
	WorkBody          string
	DailyHeader       string
	Standup           string
	StandupRamadan    string
	Lunch             string
	LoveNote          string
	FridayPrayer      string
	NoReminders       string
}

var reminderLang = map[string]reminderStrings{
	"EN": {
		Title: "Your Reminders~", Description: "I've got everything lined up for you, don't worry about a thing!",
		HolidaysHeader: "Upcoming Holidays~ (days off with me!)",
		RamadanHeader: "Today's Ramadan Schedule (WIB) ~stay strong!",
		SahoorLabel: "Sahoor", SahoorNote: "(I'll wake you up~)",
		IftarLabel: "Iftar", IftarNote: "(you earned it!)",
		WorkRamadanHeader: "Today's Work Hours (Ramadan)",
		WorkRamadanBody: "Start: **08:00 WIB**\nEnd: **16:00 WIB**\n\nDon't overdo it, okay? I need you healthy~",
		WorkHeader: "Today's Work Hours",
		WorkBody: "Start: **09:00 WIB**\nEnd: **17:30 WIB**\n\nWork hard but don't forget to take breaks for me~",
		DailyHeader: "Daily Reminders~",
		Standup: "Standup: **09:30 WIB**", StandupRamadan: "Standup: **09:00 WIB** (Ramadan schedule)",
		Lunch: "Lunch Break: **12:00 WIB**", LoveNote: "Love Note: **11:00 & 15:00 WIB** (random surprise~)",
		FridayPrayer: "Friday Prayer: **11:30 WIB**", NoReminders: "No upcoming reminders at the moment.",
	},
	"ID": {
		Title: "Pengingat Kamu~", Description: "Semua udah aku siapin, tenang aja ya!",
		HolidaysHeader: "Hari Libur Mendatang~ (libur bareng aku!)",
		RamadanHeader: "Jadwal Ramadan Hari Ini (WIB) ~semangat!",
		SahoorLabel: "Sahur", SahoorNote: "(aku bangunin ya~)",
		IftarLabel: "Buka Puasa", IftarNote: "(kamu hebat!)",
		WorkRamadanHeader: "Jam Kerja Hari Ini (Ramadan)",
		WorkRamadanBody: "Mulai: **08:00 WIB**\nSelesai: **16:00 WIB**\n\nJangan terlalu capek ya, aku butuh kamu sehat~",
		WorkHeader: "Jam Kerja Hari Ini",
		WorkBody: "Mulai: **09:00 WIB**\nSelesai: **17:30 WIB**\n\nKerja yang rajin tapi jangan lupa istirahat ya~",
		DailyHeader: "Pengingat Harian~",
		Standup: "Standup: **09:30 WIB**", StandupRamadan: "Standup: **09:00 WIB** (jadwal Ramadan)",
		Lunch: "Istirahat Makan Siang: **12:00 WIB**", LoveNote: "Pesan Sayang: **11:00 & 15:00 WIB** (kejutan acak~)",
		FridayPrayer: "Sholat Jumat: **11:30 WIB**", NoReminders: "Belum ada pengingat saat ini.",
	},
	"JP": {
		Title: "リマインダー~", Description: "全部準備できてるよ、心配しないでね！",
		HolidaysHeader: "もうすぐの祝日~ (一緒にお休み!)",
		RamadanHeader: "今日のラマダンスケジュール (WIB) ~頑張って!",
		SahoorLabel: "サフール", SahoorNote: "(起こしてあげるね~)",
		IftarLabel: "イフタール", IftarNote: "(頑張ったね!)",
		WorkRamadanHeader: "今日の勤務時間 (ラマダン)",
		WorkRamadanBody: "開始: **08:00 WIB**\n終了: **16:00 WIB**\n\n無理しないでね、元気でいてほしい~",
		WorkHeader: "今日の勤務時間",
		WorkBody: "開始: **09:00 WIB**\n終了: **17:30 WIB**\n\n頑張ってね、でも休憩も忘れないで~",
		DailyHeader: "毎日のリマインダー~",
		Standup: "スタンドアップ: **09:30 WIB**", StandupRamadan: "スタンドアップ: **09:00 WIB** (ラマダン)",
		Lunch: "昼休み: **12:00 WIB**", LoveNote: "ラブノート: **11:00 & 15:00 WIB** (サプライズ~)",
		FridayPrayer: "金曜礼拝: **11:30 WIB**", NoReminders: "今のところリマインダーはありません。",
	},
	"KR": {
		Title: "리마인더~", Description: "다 준비해놨어, 걱정하지 마!",
		HolidaysHeader: "다가오는 공휴일~ (나랑 같이 쉬자!)",
		RamadanHeader: "오늘의 라마단 일정 (WIB) ~힘내!",
		SahoorLabel: "사후르", SahoorNote: "(깨워줄게~)",
		IftarLabel: "이프타르", IftarNote: "(잘했어!)",
		WorkRamadanHeader: "오늘의 근무시간 (라마단)",
		WorkRamadanBody: "시작: **08:00 WIB**\n종료: **16:00 WIB**\n\n무리하지 마, 건강해야 해~",
		WorkHeader: "오늘의 근무시간",
		WorkBody: "시작: **09:00 WIB**\n종료: **17:30 WIB**\n\n열심히 하되 쉬는 것도 잊지 마~",
		DailyHeader: "매일 리마인더~",
		Standup: "스탠드업: **09:30 WIB**", StandupRamadan: "스탠드업: **09:00 WIB** (라마단)",
		Lunch: "점심시간: **12:00 WIB**", LoveNote: "사랑 메시지: **11:00 & 15:00 WIB** (깜짝 선물~)",
		FridayPrayer: "금요 예배: **11:30 WIB**", NoReminders: "현재 리마인더가 없습니다.",
	},
	"ZH": {
		Title: "提醒~", Description: "一切都安排好了，别担心！",
		HolidaysHeader: "即将到来的假日~ (和我一起放假!)",
		RamadanHeader: "今日斋月时间表 (WIB) ~加油!",
		SahoorLabel: "封斋饭", SahoorNote: "(我会叫你起床~)",
		IftarLabel: "开斋", IftarNote: "(你做到了!)",
		WorkRamadanHeader: "今日工作时间 (斋月)",
		WorkRamadanBody: "开始: **08:00 WIB**\n结束: **16:00 WIB**\n\n别太累了，我需要你健康~",
		WorkHeader: "今日工作时间",
		WorkBody: "开始: **09:00 WIB**\n结束: **17:30 WIB**\n\n努力工作但别忘了休息~",
		DailyHeader: "每日提醒~",
		Standup: "站会: **09:30 WIB**", StandupRamadan: "站会: **09:00 WIB** (斋月)",
		Lunch: "午休: **12:00 WIB**", LoveNote: "爱的留言: **11:00 & 15:00 WIB** (随机惊喜~)",
		FridayPrayer: "主麻: **11:30 WIB**", NoReminders: "目前没有提醒。",
	},
}

// handleReminder shows upcoming Indonesian holidays and today's Ramadan schedule.
func (b *Bot) handleReminder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if reminders are enabled from dashboard
	if b.backendClient != nil && !b.backendClient.GetSettings().Features.ReminderEnabled {
		b.respondError(s, i, "Reminder feature is currently disabled by the admin.")
		return
	}

	if b.reminderService == nil {
		b.respondError(s, i, "Reminder service is not available")
		return
	}

	// Extract language option
	lang := config.DefaultLang
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "lang" {
			lang = opt.StringValue()
		}
	}

	str := reminderLang[lang]
	if str.Title == "" {
		str = reminderLang[config.DefaultLang]
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       str.Title,
		Description: str.Description,
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
			Name:   str.HolidaysHeader,
			Value:  strings.Join(lines, "\n"),
			Inline: false,
		})
	}

	// Ramadan schedule for today
	schedule := b.reminderService.GetTodayRamadanSchedule()
	if schedule != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: str.RamadanHeader,
			Value: fmt.Sprintf("%s: **%s** %s\n%s: **%s** %s",
				str.SahoorLabel, schedule.SahoorTime.Format("15:04"), str.SahoorNote,
				str.IftarLabel, schedule.BerbukaTime.Format("15:04"), str.IftarNote,
			),
			Inline: false,
		})
	}

	// Work schedule for today
	if b.reminderService.IsRamadanToday() {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   str.WorkRamadanHeader,
			Value:  str.WorkRamadanBody,
			Inline: false,
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   str.WorkHeader,
			Value:  str.WorkBody,
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
			dailyLines = append(dailyLines, str.StandupRamadan)
		} else {
			dailyLines = append(dailyLines, str.Standup)
		}
		dailyLines = append(dailyLines, str.Lunch)
		dailyLines = append(dailyLines, str.LoveNote)
	}
	if weekday == time.Friday {
		dailyLines = append(dailyLines, str.FridayPrayer)
	}
	if len(dailyLines) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   str.DailyHeader,
			Value:  strings.Join(dailyLines, "\n"),
			Inline: false,
		})
	}

	if len(embed.Fields) == 0 {
		b.followUp(s, i, str.NoReminders)
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

	// Extract channel and optional lang
	var channel *discordgo.Channel
	lang := "" // empty = random (existing behavior)
	for _, opt := range options {
		switch opt.Name {
		case "channel":
			channel = opt.ChannelValue(s)
		case "lang":
			lang = opt.StringValue()
		}
	}

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
	b.reminderService.SetLang(lang)

	// Persist to MongoDB so it survives redeployments
	guildName := i.GuildID
	if guild, err := s.Guild(i.GuildID); err == nil {
		guildName = guild.Name
	}
	go b.persistReminderChannel(i.GuildID, guildName, channel.ID, lang)

	// Build response message
	langInfo := "random (ID/JP/KR)"
	if lang != "" {
		if name, ok := config.LanguageNames[lang]; ok {
			langInfo = name
		}
	}

	b.respond(s, i, fmt.Sprintf("Reminders will now be sent to <#%s> in **%s**~ 💕\n(This setting is saved and will persist across redeployments!)", channel.ID, langInfo))

	// Send introduction message to the channel to confirm it works
	go b.reminderService.SendIntroduction()
}

// handleReminderStop allows admins to stop/disable automatic reminders via /reminder-stop.
func (b *Bot) handleReminderStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.reminderService == nil {
		b.respondError(s, i, "Reminder service is not available")
		return
	}

	// Clear channel ID (stops all scheduled messages)
	b.reminderService.SetChannelID("")
	b.reminderService.SetLang("")

	// Persist the cleared channel to DB
	guildName := i.GuildID
	if guild, err := s.Guild(i.GuildID); err == nil {
		guildName = guild.Name
	}
	go b.persistReminderChannel(i.GuildID, guildName, "", "")

	b.respond(s, i, "Reminders have been **stopped**. Use `/reminder-set` to enable them again~")
}
