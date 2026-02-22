package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// handleHelp handles the help command.
func (b *Bot) handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract language option
	lang := config.DefaultLang
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "lang" {
			lang = opt.StringValue()
		}
	}

	// Use description from backend dashboard if available
	description := b.config.Bot.Description
	if bs := b.backendClient.GetSettings(); bs != nil && bs.BotDescription != "" {
		description = bs.BotDescription
	}

	embed := buildHelpEmbed(b.config, description, lang)
	b.respondEmbed(s, i, embed)
}

// buildHelpEmbed creates a help embed in the specified language
func buildHelpEmbed(cfg *config.Config, description string, lang string) *discordgo.MessageEmbed {
	h := helpText[lang]
	if h == nil {
		h = helpText[config.DefaultLang]
	}

	return &discordgo.MessageEmbed{
		Title:       cfg.Bot.Name + " " + h["title"],
		Description: description,
		Color:       config.ColorPrimary,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   h["confession"],
				Value:  "`/confess <content>` - " + h["confess"],
				Inline: false,
			},
			{
				Name:   h["roast"],
				Value:  "`/roast [user] [lang]` - " + h["roast_desc"],
				Inline: false,
			},
			{
				Name: h["ai"],
				Value: "`/chat <message> [lang]` - " + h["chat"] + "\n" +
					"`/chat-reset` - " + h["chat_reset"],
				Inline: false,
			},
			{
				Name:   h["news_title"],
				Value:  "`/news [lang]` - " + h["news_desc"],
				Inline: false,
			},
			{
				Name:   h["whale_title"],
				Value:  "`/whale` - " + h["whale_desc"],
				Inline: false,
			},
			{
				Name: h["analytics"],
				Value: "`/stats` - " + h["stats"] + "\n" +
					"`/profile [user]` - " + h["profile"],
				Inline: false,
			},
			{
				Name: h["reminder"],
				Value: "`/reminder` - " + h["reminder_view"] + "\n" +
					"`/reminder-set <channel>` - " + h["reminder_set"] + "\n" +
					"`/reminder-stop` - " + h["reminder_stop"],
				Inline: false,
			},
			{
				Name: h["fun"],
				Value: "`/dadjoke` - " + h["dadjoke"] + "\n" +
					"`/dadjoke-setup <channel> <interval>` - " + h["dadjoke_setup"] + "\n" +
					"`/meme` - " + h["meme"] + "\n" +
					"`/meme-setup <channel> <interval>` - " + h["meme_setup"],
				Inline: false,
			},
			{
				Name: h["mentalhealth_title"],
				Value: "`/mentalhealth [lang]` - " + h["mentalhealth"] + "\n" +
					"`/mentalhealth-setup <channel> <interval> [tag] [everyone] [lang]` - " + h["mentalhealth_setup"] + "\n" +
					"`/mentalhealth-stop` - " + h["mentalhealth_stop"],
				Inline: false,
			},
			{
				Name: h["ollama_title"],
				Value: "`/ollama-models` - " + h["ollama_models"] + "\n" +
					"`/ollama-bench <model> [prompt]` - " + h["ollama_bench"],
				Inline: false,
			},
			{
				Name:   h["other"],
				Value:  "`/help [lang]` - " + h["help_desc"],
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s v%s | %s", cfg.Bot.Name, cfg.Bot.Version, cfg.Bot.Author),
		},
	}
}

// helpText contains translated help strings per language
var helpText = map[string]map[string]string{
	"EN": {
		"title": "Help",
		"confession": "Confession Commands", "confess": "Submit an anonymous confession",
		"roast": "Roast Commands", "roast_desc": "Get roasted based on Discord activity",
		"ai": "AI Chatbot Commands", "chat": "Chat with AI", "chat_reset": "Reset your chat history",
		"news_title": "News Commands", "news_desc": "Get latest news from multiple sources",
		"whale_title": "Whale Alert Commands", "whale_desc": "Get recent whale cryptocurrency transactions",
		"analytics": "Analytics Commands", "stats": "View server statistics", "profile": "View user profile",
		"reminder": "Reminder Commands", "reminder_view": "View upcoming holidays and Ramadan schedule",
		"reminder_set": "Set/change reminder channel (admin only)", "reminder_stop": "Stop automatic reminders (admin only)",
		"fun": "Fun Commands", "dadjoke": "Get a random (clean) dad joke",
		"dadjoke_setup": "Schedule dad jokes (admin only)", "meme": "Get a random meme from the internet",
		"meme_setup": "Schedule memes (admin only)",
		"mentalhealth_title": "Mental Health", "mentalhealth": "Get a mental health tip",
		"mentalhealth_setup": "Schedule mental health reminders with tag (admin only)", "mentalhealth_stop": "Stop mental health reminders (admin only)",
		"ollama_title": "Ollama Commands", "ollama_models": "List available Ollama models", "ollama_bench": "Benchmark an Ollama model",
		"other": "Other Commands", "help_desc": "Show this help message",
	},
	"ID": {
		"title": "Bantuan",
		"confession": "Perintah Konfesi", "confess": "Kirim konfesi anonim",
		"roast": "Perintah Roast", "roast_desc": "Diroast berdasarkan aktivitas Discord",
		"ai": "Perintah AI Chatbot", "chat": "Ngobrol dengan AI", "chat_reset": "Reset riwayat chat",
		"news_title": "Perintah Berita", "news_desc": "Dapatkan berita terbaru dari berbagai sumber",
		"whale_title": "Perintah Whale Alert", "whale_desc": "Lihat transaksi whale crypto terbaru",
		"analytics": "Perintah Analitik", "stats": "Lihat statistik server", "profile": "Lihat profil pengguna",
		"reminder": "Perintah Pengingat", "reminder_view": "Lihat hari libur & jadwal Ramadan",
		"reminder_set": "Atur channel pengingat (admin)", "reminder_stop": "Hentikan pengingat otomatis (admin)",
		"fun": "Perintah Fun", "dadjoke": "Dapatkan dad joke acak",
		"dadjoke_setup": "Jadwalkan dad joke (admin)", "meme": "Dapatkan meme acak dari internet",
		"meme_setup": "Jadwalkan meme (admin)",
		"mentalhealth_title": "Kesehatan Mental", "mentalhealth": "Dapatkan tips kesehatan mental",
		"mentalhealth_setup": "Jadwalkan pengingat kesehatan mental dengan tag (admin)", "mentalhealth_stop": "Hentikan pengingat kesehatan mental (admin)",
		"ollama_title": "Perintah Ollama", "ollama_models": "Daftar model Ollama yang tersedia", "ollama_bench": "Benchmark model Ollama",
		"other": "Perintah Lainnya", "help_desc": "Tampilkan pesan bantuan ini",
	},
	"JP": {
		"title": "ヘルプ",
		"confession": "告白コマンド", "confess": "匿名の告白を送信",
		"roast": "ローストコマンド", "roast_desc": "Discordの活動に基づいてロースト",
		"ai": "AIチャットボットコマンド", "chat": "AIとチャット", "chat_reset": "チャット履歴をリセット",
		"news_title": "ニュースコマンド", "news_desc": "複数ソースから最新ニュースを取得",
		"whale_title": "ホエールアラートコマンド", "whale_desc": "最新のクジラ暗号通貨取引を取得",
		"analytics": "分析コマンド", "stats": "サーバー統計を表示", "profile": "ユーザープロフィールを表示",
		"reminder": "リマインダーコマンド", "reminder_view": "祝日とラマダンスケジュールを表示",
		"reminder_set": "リマインダーチャンネルを設定（管理者）", "reminder_stop": "自動リマインダーを停止（管理者）",
		"fun": "お楽しみコマンド", "dadjoke": "ランダムなダジャレを取得",
		"dadjoke_setup": "ダジャレをスケジュール（管理者）", "meme": "ランダムなミームを取得",
		"meme_setup": "ミームをスケジュール（管理者）",
		"mentalhealth_title": "メンタルヘルス", "mentalhealth": "メンタルヘルスのヒントを取得",
		"mentalhealth_setup": "メンタルヘルスリマインダーをスケジュール（管理者）", "mentalhealth_stop": "メンタルヘルスリマインダーを停止（管理者）",
		"ollama_title": "Ollamaコマンド", "ollama_models": "利用可能なOllamaモデル一覧", "ollama_bench": "Ollamaモデルのベンチマーク",
		"other": "その他のコマンド", "help_desc": "このヘルプメッセージを表示",
	},
	"KR": {
		"title": "도움말",
		"confession": "고백 명령어", "confess": "익명 고백 제출",
		"roast": "로스트 명령어", "roast_desc": "Discord 활동 기반으로 로스트",
		"ai": "AI 챗봇 명령어", "chat": "AI와 채팅", "chat_reset": "채팅 기록 초기화",
		"news_title": "뉴스 명령어", "news_desc": "여러 소스에서 최신 뉴스 가져오기",
		"whale_title": "고래 알림 명령어", "whale_desc": "최근 고래 암호화폐 거래 조회",
		"analytics": "분석 명령어", "stats": "서버 통계 보기", "profile": "사용자 프로필 보기",
		"reminder": "리마인더 명령어", "reminder_view": "다가오는 공휴일 및 라마단 일정 보기",
		"reminder_set": "리마인더 채널 설정 (관리자)", "reminder_stop": "자동 리마인더 중지 (관리자)",
		"fun": "재미 명령어", "dadjoke": "랜덤 아재 개그 가져오기",
		"dadjoke_setup": "아재 개그 예약 (관리자)", "meme": "인터넷에서 랜덤 밈 가져오기",
		"meme_setup": "밈 예약 (관리자)",
		"mentalhealth_title": "정신 건강", "mentalhealth": "정신 건강 팁 가져오기",
		"mentalhealth_setup": "정신 건강 리마인더 예약 (관리자)", "mentalhealth_stop": "정신 건강 리마인더 중지 (관리자)",
		"ollama_title": "Ollama 명령어", "ollama_models": "사용 가능한 Ollama 모델 목록", "ollama_bench": "Ollama 모델 벤치마크",
		"other": "기타 명령어", "help_desc": "이 도움말 메시지 표시",
	},
	"ZH": {
		"title": "帮助",
		"confession": "告白命令", "confess": "提交匿名告白",
		"roast": "吐槽命令", "roast_desc": "根据Discord活动进行吐槽",
		"ai": "AI聊天机器人命令", "chat": "与AI聊天", "chat_reset": "重置聊天记录",
		"news_title": "新闻命令", "news_desc": "从多个来源获取最新新闻",
		"whale_title": "鲸鱼提醒命令", "whale_desc": "获取最近的鲸鱼加密货币交易",
		"analytics": "分析命令", "stats": "查看服务器统计", "profile": "查看用户资料",
		"reminder": "提醒命令", "reminder_view": "查看即将到来的假日和斋月时间表",
		"reminder_set": "设置提醒频道（管理员）", "reminder_stop": "停止自动提醒（管理员）",
		"fun": "趣味命令", "dadjoke": "获取随机冷笑话",
		"dadjoke_setup": "安排冷笑话（管理员）", "meme": "从互联网获取随机表情包",
		"meme_setup": "安排表情包（管理员）",
		"mentalhealth_title": "心理健康", "mentalhealth": "获取心理健康小贴士",
		"mentalhealth_setup": "安排心理健康提醒（管理员）", "mentalhealth_stop": "停止心理健康提醒（管理员）",
		"ollama_title": "Ollama 命令", "ollama_models": "列出可用的 Ollama 模型", "ollama_bench": "对 Ollama 模型进行基准测试",
		"other": "其他命令", "help_desc": "显示此帮助信息",
	},
}

// --- Response helpers ---

// deferResponse sends a deferred response to the interaction.
func (b *Bot) deferResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
	}
	return err
}

// respond sends a text response to the interaction.
func (b *Bot) respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond to interaction", "error", err)
	}
}

// respondEmbed sends an embed response to the interaction.
func (b *Bot) respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond to interaction", "error", err)
	}
}

// respondError sends an error text response to the interaction.
func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	b.respond(s, i, config.EmojiError+" "+message)
}

// followUp sends a follow-up message after a deferred response.
func (b *Bot) followUp(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

// followUpEmbed sends a follow-up embed message after a deferred response.
func (b *Bot) followUpEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up embed", "error", err)
	}
}

// followUpError sends an error follow-up message.
func (b *Bot) followUpError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	b.followUp(s, i, config.EmojiError+" "+message)
}

// findUserVoiceState returns the voice state of a user in a guild.
func (b *Bot) findUserVoiceState(guildID, userID string) *discordgo.VoiceState {
	guild, err := b.session.State.Guild(guildID)
	if err != nil {
		return nil
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs
		}
	}

	return nil
}
