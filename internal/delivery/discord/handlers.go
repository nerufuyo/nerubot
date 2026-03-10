package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// helpPage represents a single page of the help menu.
type helpPage struct {
	Title  string
	Emoji  string
	Fields []*discordgo.MessageEmbedField
}

// buildHelpPages creates all paginated help pages in the specified language.
func buildHelpPages(cfg *config.Config, description string, lang string) []helpPage {
	h := helpText[lang]
	if h == nil {
		h = helpText[config.DefaultLang]
	}

	return []helpPage{
		{
			Title: h["title"] + " — " + h["overview"],
			Emoji: "📖",
			Fields: []*discordgo.MessageEmbedField{
				{Name: "", Value: description, Inline: false},
				{
					Name: h["core"],
					Value: "`/ping` - " + h["ping"] + "\n" +
						"`/botinfo` - " + h["botinfo"] + "\n" +
						"`/serverinfo` - " + h["serverinfo"] + "\n" +
						"`/userinfo [user]` - " + h["userinfo"] + "\n" +
						"`/avatar [user]` - " + h["avatar"] + "\n" +
						"`/help [lang]` - " + h["help_desc"],
					Inline: false,
				},
			},
		},
		{
			Title: h["title"] + " — " + h["fun"],
			Emoji: "🎮",
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
					Name: h["fun"],
					Value: "`/dadjoke` - " + h["dadjoke"] + "\n" +
						"`/meme` - " + h["meme"] + "\n" +
						"`/coinflip` - " + h["coinflip"] + "\n" +
						"`/8ball <question>` - " + h["eightball"],
					Inline: false,
				},
			},
		},
		{
			Title: h["title"] + " — " + h["utility"],
			Emoji: "🔧",
			Fields: []*discordgo.MessageEmbedField{
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
					Name:   h["calc"],
					Value:  "`/calc <a> <operation> <b>` - " + h["calc_desc"],
					Inline: false,
				},
				{
					Name:   h["poll_title"],
					Value:  "`/poll <question> <option1> <option2> ...` - " + h["poll_desc"],
					Inline: false,
				},
			},
		},
		{
			Title: h["title"] + " — " + h["moderation"],
			Emoji: "🛡️",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: h["moderation"],
					Value: "`/kick <user> [reason]` - " + h["kick"] + "\n" +
						"`/ban <user> [reason]` - " + h["ban"] + "\n" +
						"`/timeout <user> <duration> [reason]` - " + h["timeout"] + "\n" +
						"`/purge <amount>` - " + h["purge"] + "\n" +
						"`/warn <user> [reason]` - " + h["warn"] + "\n" +
						"`/warnings <user>` - " + h["warnings_desc"] + "\n" +
						"`/clearwarnings <user>` - " + h["clearwarnings"],
					Inline: false,
				},
			},
		},
		{
			Title: h["title"] + " — " + h["scheduling"],
			Emoji: "⏰",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: h["reminder"],
					Value: "`/reminder` - " + h["reminder_view"] + "\n" +
						"`/reminder-set <channel>` - " + h["reminder_set"] + "\n" +
						"`/reminder-stop` - " + h["reminder_stop"],
					Inline: false,
				},
				{
					Name: h["scheduled_fun"],
					Value: "`/dadjoke-setup <channel> <interval>` - " + h["dadjoke_setup"] + "\n" +
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
			},
		},
		{
			Title: h["title"] + " — " + h["music_title"],
			Emoji: "🎵",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: h["music_playback"],
					Value: "`/play <query>` - " + h["play"] + "\n" +
						"`/pause` - " + h["pause"] + "\n" +
						"`/resume` - " + h["resume"] + "\n" +
						"`/stop` - " + h["stop"] + "\n" +
						"`/skip` - " + h["skip"] + "\n" +
						"`/previous` - " + h["previous"] + "\n" +
						"`/seek <seconds>` - " + h["seek"],
					Inline: false,
				},
				{
					Name: h["music_queue"],
					Value: "`/queue` - " + h["queue"] + "\n" +
						"`/nowplaying` - " + h["nowplaying"] + "\n" +
						"`/shuffle` - " + h["shuffle"] + "\n" +
						"`/remove <pos>` - " + h["remove"] + "\n" +
						"`/move <from> <to>` - " + h["move"] + "\n" +
						"`/clear` - " + h["clear_queue"] + "\n" +
						"`/volume <0-150>` - " + h["volume"],
					Inline: false,
				},
				{
					Name: h["music_advanced"],
					Value: "`/loop <off|song|queue>` - " + h["loop"] + "\n" +
						"`/filter <name>` - " + h["filter"] + "\n" +
						"`/playlist` - " + h["playlist"] + "\n" +
						"`/lyrics` - " + h["lyrics"] + "\n" +
						"`/recommend` - " + h["recommend"] + "\n" +
						"`/radio <genre>` - " + h["radio"] + "\n" +
						"`/autoplay` - " + h["autoplay"] + "\n" +
						"`/voteskip` - " + h["voteskip"],
					Inline: false,
				},
			},
		},
	}
}

// handleHelp handles the help command with paginated embeds and navigation buttons.
func (b *Bot) handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lang := config.DefaultLang
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "lang" {
			lang = opt.StringValue()
		}
	}

	description := b.config.Bot.Description
	if bs := b.backendClient.GetSettings(); bs != nil && bs.BotDescription != "" {
		description = bs.BotDescription
	}

	pages := buildHelpPages(b.config, description, lang)
	embed := buildHelpPageEmbed(b.config, pages, 0)
	buttons := buildHelpButtons(0, len(pages))

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: buttons,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond to help", "error", err)
	}
}

// handleHelpButton handles help pagination button clicks.
func (b *Bot) handleHelpButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID

	// Parse page number from custom ID (format: help_page_N)
	var page int
	_, _ = fmt.Sscanf(customID, "help_page_%d", &page)

	lang := config.DefaultLang // default; we don't persist lang in button state
	description := b.config.Bot.Description
	if bs := b.backendClient.GetSettings(); bs != nil && bs.BotDescription != "" {
		description = bs.BotDescription
	}

	pages := buildHelpPages(b.config, description, lang)
	if page < 0 {
		page = 0
	}
	if page >= len(pages) {
		page = len(pages) - 1
	}

	embed := buildHelpPageEmbed(b.config, pages, page)
	buttons := buildHelpButtons(page, len(pages))

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: buttons,
		},
	})
	if err != nil {
		b.logger.Error("Failed to update help page", "error", err)
	}
}

// buildHelpPageEmbed creates a Discord embed for a specific help page.
func buildHelpPageEmbed(cfg *config.Config, pages []helpPage, page int) *discordgo.MessageEmbed {
	p := pages[page]
	return &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("%s %s %s", p.Emoji, cfg.Bot.Name, p.Title),
		Color:  config.ColorPrimary,
		Fields: p.Fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d • %s v%s | %s", page+1, len(pages), cfg.Bot.Name, cfg.Bot.Version, cfg.Bot.Author),
		},
	}
}

// buildHelpButtons creates navigation buttons for help pagination.
func buildHelpButtons(currentPage, totalPages int) []discordgo.MessageComponent {
	var buttons []discordgo.MessageComponent

	// Previous button
	prevDisabled := currentPage == 0
	buttons = append(buttons, discordgo.Button{
		Label:    "◀ Prev",
		Style:    discordgo.SecondaryButton,
		CustomID: fmt.Sprintf("help_page_%d", currentPage-1),
		Disabled: prevDisabled,
	})

	// Page indicator button (non-interactive)
	buttons = append(buttons, discordgo.Button{
		Label:    fmt.Sprintf("%d / %d", currentPage+1, totalPages),
		Style:    discordgo.SecondaryButton,
		CustomID: "help_page_indicator",
		Disabled: true,
	})

	// Next button
	nextDisabled := currentPage >= totalPages-1
	buttons = append(buttons, discordgo.Button{
		Label:    "Next ▶",
		Style:    discordgo.SecondaryButton,
		CustomID: fmt.Sprintf("help_page_%d", currentPage+1),
		Disabled: nextDisabled,
	})

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: buttons},
	}
}

// helpText contains translated help strings per language
var helpText = map[string]map[string]string{
	"EN": {
		"title": "Help", "overview": "Overview",
		"core": "Core Commands", "ping": "Check bot latency", "botinfo": "Show bot information & uptime",
		"serverinfo": "Show server details", "userinfo": "Show user information", "avatar": "Display user's avatar",
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
		"meme_setup": "Schedule memes (admin only)", "coinflip": "Flip a coin", "eightball": "Ask the magic 8-ball",
		"mentalhealth_title": "Mental Health", "mentalhealth": "Get a mental health tip",
		"mentalhealth_setup": "Schedule mental health reminders with tag (admin only)", "mentalhealth_stop": "Stop mental health reminders (admin only)",
		"other": "Other Commands", "help_desc": "Show this help message",
		"utility": "Utility & Info", "calc": "Calculator", "calc_desc": "Simple math calculations",
		"poll_title": "Poll", "poll_desc": "Create a simple poll with reactions",
		"moderation": "Moderation", "kick": "Kick a user from the server",
		"ban": "Ban a user from the server", "timeout": "Temporarily mute a user",
		"purge": "Delete multiple messages", "warn": "Warn a user and record it",
		"warnings_desc": "View warnings for a user", "clearwarnings": "Clear all warnings (admin only)",
		"scheduling": "Scheduling & Reminders", "scheduled_fun": "Scheduled Fun",
		"music_title": "Music", "music_playback": "Playback",
		"play": "Play a song from query or link", "pause": "Pause the current song",
		"resume": "Resume the paused song", "stop": "Stop music and clear queue",
		"skip": "Skip to the next song", "previous": "Play the previous song", "seek": "Seek to a position",
		"music_queue": "Queue Management",
		"queue": "Show the music queue", "nowplaying": "Show the currently playing song",
		"shuffle": "Shuffle the queue", "remove": "Remove a song from queue",
		"move": "Move a song position", "clear_queue": "Clear the entire queue", "volume": "Adjust the volume",
		"music_advanced": "Advanced Music",
		"loop": "Set loop mode", "filter": "Apply audio filters",
		"playlist": "Manage playlists", "lyrics": "Show song lyrics",
		"recommend": "Get song recommendations", "radio": "Start genre radio",
		"autoplay": "Toggle autoplay", "voteskip": "Vote to skip the current song",
	},
	"ID": {
		"title": "Bantuan", "overview": "Ringkasan",
		"core": "Perintah Utama", "ping": "Cek latensi bot", "botinfo": "Info bot & uptime",
		"serverinfo": "Info server", "userinfo": "Info pengguna", "avatar": "Tampilkan avatar pengguna",
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
		"meme_setup": "Jadwalkan meme (admin)", "coinflip": "Lempar koin", "eightball": "Tanya bola ajaib 8",
		"mentalhealth_title": "Kesehatan Mental", "mentalhealth": "Dapatkan tips kesehatan mental",
		"mentalhealth_setup": "Jadwalkan pengingat kesehatan mental dengan tag (admin)", "mentalhealth_stop": "Hentikan pengingat kesehatan mental (admin)",
		"other": "Perintah Lainnya", "help_desc": "Tampilkan pesan bantuan ini",
		"utility": "Utilitas & Info", "calc": "Kalkulator", "calc_desc": "Perhitungan matematika sederhana",
		"poll_title": "Polling", "poll_desc": "Buat polling sederhana dengan reaksi",
		"moderation": "Moderasi", "kick": "Tendang pengguna dari server",
		"ban": "Ban pengguna dari server", "timeout": "Bisukan pengguna sementara",
		"purge": "Hapus banyak pesan", "warn": "Beri peringatan kepada pengguna",
		"warnings_desc": "Lihat peringatan pengguna", "clearwarnings": "Hapus semua peringatan (admin)",
		"scheduling": "Penjadwalan & Pengingat", "scheduled_fun": "Fun Terjadwal",
		"music_title": "Musik", "music_playback": "Pemutaran",
		"play": "Putar lagu dari query atau link", "pause": "Jeda lagu saat ini",
		"resume": "Lanjutkan lagu", "stop": "Hentikan musik dan bersihkan antrian",
		"skip": "Lewati ke lagu berikutnya", "previous": "Putar lagu sebelumnya", "seek": "Loncat ke posisi tertentu",
		"music_queue": "Manajemen Antrian",
		"queue": "Tampilkan antrian musik", "nowplaying": "Tampilkan lagu yang sedang diputar",
		"shuffle": "Acak antrian", "remove": "Hapus lagu dari antrian",
		"move": "Pindahkan posisi lagu", "clear_queue": "Bersihkan seluruh antrian", "volume": "Atur volume",
		"music_advanced": "Musik Lanjutan",
		"loop": "Atur mode loop", "filter": "Terapkan filter audio",
		"playlist": "Kelola playlist", "lyrics": "Tampilkan lirik lagu",
		"recommend": "Dapatkan rekomendasi lagu", "radio": "Mulai radio genre",
		"autoplay": "Toggle autoplay", "voteskip": "Vote untuk lewati lagu",
	},
	"JP": {
		"title": "ヘルプ", "overview": "概要",
		"core": "基本コマンド", "ping": "ボットの遅延を確認", "botinfo": "ボット情報とアップタイムを表示",
		"serverinfo": "サーバー詳細を表示", "userinfo": "ユーザー情報を表示", "avatar": "ユーザーのアバターを表示",
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
		"meme_setup": "ミームをスケジュール（管理者）", "coinflip": "コイントス", "eightball": "マジック8ボールに質問",
		"mentalhealth_title": "メンタルヘルス", "mentalhealth": "メンタルヘルスのヒントを取得",
		"mentalhealth_setup": "メンタルヘルスリマインダーをスケジュール（管理者）", "mentalhealth_stop": "メンタルヘルスリマインダーを停止（管理者）",
		"other": "その他のコマンド", "help_desc": "このヘルプメッセージを表示",
		"utility": "ユーティリティ＆情報", "calc": "電卓", "calc_desc": "簡単な数学計算",
		"poll_title": "投票", "poll_desc": "リアクション付きの簡単な投票を作成",
		"moderation": "モデレーション", "kick": "ユーザーをサーバーからキック",
		"ban": "ユーザーをサーバーからBAN", "timeout": "ユーザーを一時的にミュート",
		"purge": "複数のメッセージを削除", "warn": "ユーザーに警告を与えて記録",
		"warnings_desc": "ユーザーの警告を表示", "clearwarnings": "すべての警告をクリア（管理者）",
		"scheduling": "スケジュール＆リマインダー", "scheduled_fun": "スケジュールされた楽しみ",
		"music_title": "音楽", "music_playback": "再生",
		"play": "曲を再生", "pause": "一時停止",
		"resume": "再開", "stop": "停止してキューをクリア",
		"skip": "次の曲にスキップ", "previous": "前の曲を再生", "seek": "位置にシーク",
		"music_queue": "キュー管理",
		"queue": "キューを表示", "nowplaying": "再生中の曲を表示",
		"shuffle": "シャッフル", "remove": "キューから削除",
		"move": "曲の位置を移動", "clear_queue": "キューをクリア", "volume": "音量調整",
		"music_advanced": "高度な音楽",
		"loop": "ループモード設定", "filter": "オーディオフィルター適用",
		"playlist": "プレイリスト管理", "lyrics": "歌詞を表示",
		"recommend": "おすすめ曲を取得", "radio": "ジャンルラジオを開始",
		"autoplay": "自動再生切替", "voteskip": "スキップに投票",
	},
	"KR": {
		"title": "도움말", "overview": "개요",
		"core": "기본 명령어", "ping": "봇 지연 시간 확인", "botinfo": "봇 정보 및 가동 시간 표시",
		"serverinfo": "서버 정보 표시", "userinfo": "사용자 정보 표시", "avatar": "사용자 아바타 표시",
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
		"meme_setup": "밈 예약 (관리자)", "coinflip": "동전 던지기", "eightball": "마법의 8볼에게 질문",
		"mentalhealth_title": "정신 건강", "mentalhealth": "정신 건강 팁 가져오기",
		"mentalhealth_setup": "정신 건강 리마인더 예약 (관리자)", "mentalhealth_stop": "정신 건강 리마인더 중지 (관리자)",
		"other": "기타 명령어", "help_desc": "이 도움말 메시지 표시",
		"utility": "유틸리티 및 정보", "calc": "계산기", "calc_desc": "간단한 수학 계산",
		"poll_title": "투표", "poll_desc": "리액션으로 간단한 투표 만들기",
		"moderation": "관리", "kick": "서버에서 사용자 추방",
		"ban": "서버에서 사용자 차단", "timeout": "사용자를 일시적으로 음소거",
		"purge": "여러 메시지 삭제", "warn": "사용자에게 경고 부여",
		"warnings_desc": "사용자 경고 조회", "clearwarnings": "모든 경고 삭제 (관리자)",
		"scheduling": "예약 및 리마인더", "scheduled_fun": "예약된 재미",
		"music_title": "음악", "music_playback": "재생",
		"play": "노래 재생", "pause": "일시 정지",
		"resume": "재개", "stop": "음악 정지 및 대기열 초기화",
		"skip": "다음 곡으로 건너뛰기", "previous": "이전 곡 재생", "seek": "위치로 이동",
		"music_queue": "대기열 관리",
		"queue": "대기열 표시", "nowplaying": "현재 재생 중인 곡 표시",
		"shuffle": "셔플", "remove": "대기열에서 제거",
		"move": "곡 위치 이동", "clear_queue": "대기열 초기화", "volume": "볼륨 조절",
		"music_advanced": "고급 음악",
		"loop": "루프 모드 설정", "filter": "오디오 필터 적용",
		"playlist": "플레이리스트 관리", "lyrics": "가사 표시",
		"recommend": "추천 곡 가져오기", "radio": "장르 라디오 시작",
		"autoplay": "자동 재생 전환", "voteskip": "건너뛰기 투표",
	},
	"ZH": {
		"title": "帮助", "overview": "概览",
		"core": "核心命令", "ping": "检查机器人延迟", "botinfo": "显示机器人信息和运行时间",
		"serverinfo": "显示服务器详情", "userinfo": "显示用户信息", "avatar": "显示用户头像",
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
		"meme_setup": "安排表情包（管理员）", "coinflip": "抛硬币", "eightball": "问魔法8球",
		"mentalhealth_title": "心理健康", "mentalhealth": "获取心理健康小贴士",
		"mentalhealth_setup": "安排心理健康提醒（管理员）", "mentalhealth_stop": "停止心理健康提醒（管理员）",
		"other": "其他命令", "help_desc": "显示此帮助信息",
		"utility": "工具与信息", "calc": "计算器", "calc_desc": "简单的数学计算",
		"poll_title": "投票", "poll_desc": "用表情创建简单投票",
		"moderation": "管理", "kick": "将用户踢出服务器",
		"ban": "禁止用户进入服务器", "timeout": "临时禁言用户",
		"purge": "批量删除消息", "warn": "警告用户并记录",
		"warnings_desc": "查看用户的警告", "clearwarnings": "清除所有警告（管理员）",
		"scheduling": "计划与提醒", "scheduled_fun": "定时趣味",
		"music_title": "音乐", "music_playback": "播放",
		"play": "播放歌曲", "pause": "暂停",
		"resume": "恢复", "stop": "停止音乐并清空队列",
		"skip": "跳到下一首", "previous": "播放上一首", "seek": "跳转到指定位置",
		"music_queue": "队列管理",
		"queue": "显示音乐队列", "nowplaying": "显示当前播放的歌曲",
		"shuffle": "随机打乱", "remove": "从队列中移除",
		"move": "移动歌曲位置", "clear_queue": "清空整个队列", "volume": "调整音量",
		"music_advanced": "高级音乐",
		"loop": "设置循环模式", "filter": "应用音频滤镜",
		"playlist": "管理播放列表", "lyrics": "显示歌词",
		"recommend": "获取歌曲推荐", "radio": "启动风格电台",
		"autoplay": "切换自动播放", "voteskip": "投票跳过当前歌曲",
	},
}

// Ensure strings import is used (used in button handling)
var _ = strings.HasPrefix

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
