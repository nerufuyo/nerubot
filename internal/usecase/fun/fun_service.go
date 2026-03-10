package fun

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// SendFunc is a callback to send a message to a Discord channel.
// Returns an error if the channel is invalid or the message could not be sent.
type SendFunc func(channelID string, embed *FunEmbed) error

// FunEmbed holds data for an embed message.
type FunEmbed struct {
	Title       string
	Description string
	ImageURL    string
	Footer      string
	Color       int
	URL         string
	Content     string // text sent alongside the embed (e.g. mentions)
}

// FunService manages dad jokes, memes, and their scheduled delivery.
type FunService struct {
	mu         sync.RWMutex
	repo       *repository.GuildConfigRepository
	logger     *logger.Logger
	httpClient *http.Client
	sendFn     SendFunc
	stopCh     chan struct{}
	wg         sync.WaitGroup
	rng        *rand.Rand

	// Track last fire times per guild to avoid duplicates
	lastJoke         map[string]time.Time
	lastMeme         map[string]time.Time
	lastMentalHealth map[string]time.Time
}

// NewFunService creates a new fun service.
func NewFunService() *FunService {
	return &FunService{
		repo:   repository.NewGuildConfigRepository(),
		logger: logger.New("fun"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		stopCh:           make(chan struct{}),
		rng:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lastJoke:         make(map[string]time.Time),
		lastMeme:         make(map[string]time.Time),
		lastMentalHealth: make(map[string]time.Time),
	}
}

// SetSendFunc sets the callback used to post embeds.
func (s *FunService) SetSendFunc(fn SendFunc) {
	s.sendFn = fn
}

// Start begins the background scheduler that checks every minute for scheduled jokes/memes.
func (s *FunService) Start() {
	s.wg.Add(1)
	go s.loop()
	s.logger.Info("Fun service scheduler started")
}

// Stop gracefully shuts down the scheduler.
func (s *FunService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.logger.Info("Fun service scheduler stopped")
}

func (s *FunService) loop() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.checkScheduled()
		}
	}
}

func (s *FunService) checkScheduled() {
	configs, err := s.repo.GetAll()
	if err != nil {
		s.logger.Warn("Failed to load guild configs for scheduler", "error", err)
		return
	}

	now := time.Now()

	for _, cfg := range configs {
		// Check dad jokes schedule
		if cfg.DadJokeChannelID != "" && cfg.DadJokeInterval > 0 {
			lastFired, ok := s.lastJoke[cfg.GuildID]
			interval := time.Duration(cfg.DadJokeInterval) * time.Minute
			if !ok || now.Sub(lastFired) >= interval {
				s.lastJoke[cfg.GuildID] = now
				go s.sendScheduledJoke(cfg.DadJokeChannelID, cfg.GuildID)
			}
		}

		// Check memes schedule
		if cfg.MemeChannelID != "" && cfg.MemeInterval > 0 {
			lastFired, ok := s.lastMeme[cfg.GuildID]
			interval := time.Duration(cfg.MemeInterval) * time.Minute
			if !ok || now.Sub(lastFired) >= interval {
				s.lastMeme[cfg.GuildID] = now
				go s.sendScheduledMeme(cfg.MemeChannelID, cfg.GuildID)
			}
		}

		// Check mental health reminders schedule
		if cfg.MentalHealthChannelID != "" && cfg.MentalHealthInterval > 0 {
			lastFired, ok := s.lastMentalHealth[cfg.GuildID]
			interval := time.Duration(cfg.MentalHealthInterval) * time.Minute
			if !ok || now.Sub(lastFired) >= interval {
				s.lastMentalHealth[cfg.GuildID] = now
				go s.sendScheduledMentalHealth(cfg.MentalHealthChannelID, cfg.MentalHealthTag, cfg.MentalHealthLang, cfg.GuildID)
			}
		}
	}
}

func (s *FunService) sendScheduledJoke(channelID, guildID string) {
	joke, err := s.FetchDadJoke()
	if err != nil {
		s.logger.Warn("Scheduled dad joke fetch failed", "error", err)
		return
	}

	embed := &FunEmbed{
		Title:       "🤣 Dad Joke of the Hour",
		Description: joke.Punchline,
		Footer:      "Powered by icanhazdadjoke.com",
		Color:       0xFFD700, // gold
	}
	if s.sendFn != nil {
		if err := s.sendFn(channelID, embed); err != nil {
			s.handleSendError(guildID, "dadjoke", channelID, err)
		}
	}
}

func (s *FunService) sendScheduledMeme(channelID, guildID string) {
	meme, err := s.FetchMeme()
	if err != nil {
		s.logger.Warn("Scheduled meme fetch failed", "error", err)
		return
	}

	embed := &FunEmbed{
		Title:    "😂 " + meme.Title,
		ImageURL: meme.URL,
		Footer:   fmt.Sprintf("r/%s • by u/%s", meme.Subreddit, meme.Author),
		Color:    0xFF4500, // reddit orange
		URL:      meme.PostLink,
	}
	if s.sendFn != nil {
		if err := s.sendFn(channelID, embed); err != nil {
			s.handleSendError(guildID, "meme", channelID, err)
		}
	}
}

// --- Mental Health Tips ---

// mentalHealthTips holds curated mental health tips per language.
var mentalHealthTips = map[string][]struct {
	Title string
	Tip   string
}{
	"EN": {
		{Title: "Take a Deep Breath", Tip: "Try the 4-7-8 breathing technique: inhale for 4 seconds, hold for 7, exhale for 8. This activates your parasympathetic nervous system and calms anxiety."},
		{Title: "Move Your Body", Tip: "Even a 10-minute walk can boost your mood. Physical activity releases endorphins, your brain's natural feel-good chemicals."},
		{Title: "Stay Hydrated", Tip: "Dehydration can affect your mood and concentration. Aim for at least 8 glasses of water today. Your brain is about 75% water!"},
		{Title: "Screen Break", Tip: "Follow the 20-20-20 rule: every 20 minutes, look at something 20 feet away for 20 seconds. Your eyes and mind will thank you."},
		{Title: "Practice Gratitude", Tip: "Write down 3 things you're grateful for today. Gratitude rewires your brain to focus on the positive and reduces stress hormones."},
		{Title: "It's OK to Rest", Tip: "Rest is not laziness. Your brain needs downtime to process information and recharge. Give yourself permission to take a break."},
		{Title: "Connect With Someone", Tip: "Reach out to a friend or family member today. Social connection is one of the strongest protectors of mental health."},
		{Title: "Mindful Moment", Tip: "Pause for 60 seconds. Notice 5 things you can see, 4 you can hear, 3 you can touch, 2 you can smell, and 1 you can taste."},
		{Title: "Sleep Matters", Tip: "Try to maintain a consistent sleep schedule. Quality sleep is essential for emotional regulation and mental clarity."},
		{Title: "Be Kind to Yourself", Tip: "Talk to yourself the way you'd talk to a good friend. Self-compassion reduces anxiety and builds emotional resilience."},
		{Title: "Set Boundaries", Tip: "It's okay to say no. Protecting your energy is not selfish — it's necessary for your well-being."},
		{Title: "Celebrate Small Wins", Tip: "Did you get out of bed? Eat a meal? Complete a task? Every small step counts. Progress is progress, no matter how small."},
	},
	"ID": {
		{Title: "Tarik Napas Dalam", Tip: "Coba teknik pernapasan 4-7-8: tarik napas 4 detik, tahan 7 detik, hembuskan 8 detik. Ini mengaktifkan sistem saraf parasimpatik dan menenangkan kecemasan."},
		{Title: "Gerakkan Tubuhmu", Tip: "Jalan kaki 10 menit saja bisa meningkatkan mood. Aktivitas fisik melepaskan endorfin, zat kimia alami yang membuat otak merasa senang."},
		{Title: "Jaga Hidrasi", Tip: "Dehidrasi bisa mempengaruhi suasana hati dan konsentrasi. Usahakan minum minimal 8 gelas air hari ini. Otakmu 75% terdiri dari air!"},
		{Title: "Istirahat dari Layar", Tip: "Ikuti aturan 20-20-20: setiap 20 menit, lihat sesuatu yang berjarak 20 kaki selama 20 detik. Mata dan pikiranmu akan berterima kasih."},
		{Title: "Latih Rasa Syukur", Tip: "Tulis 3 hal yang kamu syukuri hari ini. Rasa syukur melatih otakmu untuk fokus pada hal positif dan mengurangi hormon stres."},
		{Title: "Boleh Istirahat", Tip: "Istirahat bukan kemalasan. Otakmu butuh waktu untuk memproses informasi dan mengisi ulang energi. Beri dirimu izin untuk beristirahat."},
		{Title: "Terhubung dengan Seseorang", Tip: "Hubungi teman atau keluarga hari ini. Koneksi sosial adalah salah satu pelindung terkuat kesehatan mental."},
		{Title: "Momen Mindfulness", Tip: "Berhenti sejenak selama 60 detik. Perhatikan 5 hal yang bisa kamu lihat, 4 yang bisa kamu dengar, 3 yang bisa kamu sentuh, 2 yang bisa kamu cium, dan 1 yang bisa kamu rasakan."},
		{Title: "Tidur Itu Penting", Tip: "Usahakan jadwal tidur yang konsisten. Tidur berkualitas sangat penting untuk regulasi emosi dan kejernihan mental."},
		{Title: "Sayangi Dirimu", Tip: "Bicaralah pada dirimu seperti kamu berbicara pada sahabat. Belas kasih pada diri sendiri mengurangi kecemasan dan membangun ketahanan emosional."},
		{Title: "Tetapkan Batasan", Tip: "Tidak apa-apa mengatakan tidak. Melindungi energimu bukan egois — itu perlu untuk kesejahteraanmu."},
		{Title: "Rayakan Kemenangan Kecil", Tip: "Bangun dari tempat tidur? Makan? Menyelesaikan tugas? Setiap langkah kecil berarti. Kemajuan tetap kemajuan, sekecil apa pun."},
	},
	"JP": {
		{Title: "深呼吸しよう", Tip: "4-7-8呼吸法を試してみて：4秒吸って、7秒止めて、8秒で吐く。副交感神経が活性化され、不安が和らぎます。"},
		{Title: "体を動かそう", Tip: "10分の散歩でも気分が上がります。運動はエンドルフィンという脳の天然の幸福物質を放出します。"},
		{Title: "水分補給を忘れずに", Tip: "脱水は気分や集中力に影響します。今日は最低8杯の水を飲みましょう。脳の約75%は水でできています！"},
		{Title: "画面から離れよう", Tip: "20-20-20ルールを実践：20分ごとに、20フィート先のものを20秒間見ましょう。目と心が感謝します。"},
		{Title: "感謝を実践しよう", Tip: "今日感謝していることを3つ書き出しましょう。感謝の気持ちは脳をポジティブに配線し直し、ストレスホルモンを減らします。"},
		{Title: "休んでも大丈夫", Tip: "休むことは怠けることではありません。脳には情報を処理してリチャージする時間が必要です。自分に休む許可を与えましょう。"},
		{Title: "誰かとつながろう", Tip: "今日、友人や家族に連絡してみましょう。社会的なつながりは、メンタルヘルスの最も強い守り手の一つです。"},
		{Title: "マインドフルな瞬間", Tip: "60秒間立ち止まって。見えるもの5つ、聞こえるもの4つ、触れるもの3つ、嗅げるもの2つ、味わえるもの1つに注目しましょう。"},
		{Title: "睡眠は大切", Tip: "一定の睡眠スケジュールを維持しましょう。質の良い睡眠は感情の調整と精神的な明晰さに不可欠です。"},
		{Title: "自分に優しく", Tip: "親友に話すように自分に話しかけましょう。セルフコンパッションは不安を減らし、感情的な回復力を築きます。"},
		{Title: "境界線を設定しよう", Tip: "断っても大丈夫です。自分のエネルギーを守ることは自己中心的ではなく、幸福のために必要なことです。"},
		{Title: "小さな勝利を祝おう", Tip: "ベッドから出れた？食事した？タスクを完了した？どんな小さな一歩も大切です。進歩は進歩、どんなに小さくても。"},
	},
	"KR": {
		{Title: "깊게 숨을 쉬세요", Tip: "4-7-8 호흡법을 시도해보세요: 4초 동안 들이쉬고, 7초 동안 참고, 8초 동안 내쉬세요. 부교감 신경계를 활성화하고 불안을 진정시킵니다."},
		{Title: "몸을 움직이세요", Tip: "10분만 걸어도 기분이 좋아집니다. 신체 활동은 뇌의 천연 기분 전환 화학물질인 엔도르핀을 방출합니다."},
		{Title: "수분을 충분히", Tip: "탈수는 기분과 집중력에 영향을 줄 수 있습니다. 오늘 최소한 8잔의 물을 마시세요. 뇌의 약 75%가 물입니다!"},
		{Title: "화면에서 벗어나세요", Tip: "20-20-20 규칙을 따르세요: 20분마다 20피트 떨어진 곳을 20초 동안 바라보세요. 눈과 마음이 감사할 것입니다."},
		{Title: "감사를 연습하세요", Tip: "오늘 감사한 것 3가지를 적어보세요. 감사는 뇌를 긍정적으로 재배선하고 스트레스 호르몬을 줄여줍니다."},
		{Title: "쉬어도 괜찮아요", Tip: "쉬는 것은 게으름이 아닙니다. 뇌는 정보를 처리하고 재충전할 시간이 필요합니다. 자신에게 쉬는 것을 허락하세요."},
		{Title: "누군가와 연결하세요", Tip: "오늘 친구나 가족에게 연락해보세요. 사회적 연결은 정신 건강의 가장 강력한 보호 요소 중 하나입니다."},
		{Title: "마음챙김 순간", Tip: "60초 동안 멈추세요. 볼 수 있는 것 5가지, 들을 수 있는 것 4가지, 만질 수 있는 것 3가지, 냄새 맡을 수 있는 것 2가지, 맛볼 수 있는 것 1가지에 주목하세요."},
		{Title: "수면이 중요합니다", Tip: "일정한 수면 스케줄을 유지하세요. 양질의 수면은 감정 조절과 정신적 명확성에 필수적입니다."},
		{Title: "자신에게 친절하세요", Tip: "좋은 친구에게 말하듯 자신에게 말하세요. 자기 연민은 불안을 줄이고 감정적 회복력을 키웁니다."},
		{Title: "경계를 설정하세요", Tip: "거절해도 괜찮습니다. 에너지를 보호하는 것은 이기적인 것이 아니라 행복을 위해 필요한 것입니다."},
		{Title: "작은 승리를 축하하세요", Tip: "침대에서 일어났나요? 식사했나요? 할 일을 마쳤나요? 아무리 작은 발걸음이라도 중요합니다. 진전은 진전입니다."},
	},
	"ZH": {
		{Title: "深呼吸", Tip: "试试4-7-8呼吸法：吸气4秒，屏气7秒，呼气8秒。这能激活副交感神经系统，缓解焦虑。"},
		{Title: "动起来", Tip: "即使只是10分钟的散步也能改善心情。身体活动会释放内啡肽，大脑天然的快乐化学物质。"},
		{Title: "保持水分", Tip: "脱水会影响你的情绪和注意力。今天至少喝8杯水吧。大脑大约75%是水！"},
		{Title: "屏幕休息", Tip: "遵循20-20-20法则：每20分钟，看20英尺远的东西20秒。你的眼睛和大脑会感谢你的。"},
		{Title: "练习感恩", Tip: "写下今天你感激的3件事。感恩能重新连接你的大脑，让它关注积极的事物，减少压力荷尔蒙。"},
		{Title: "休息也没关系", Tip: "休息不是懒惰。你的大脑需要停机时间来处理信息和充电。允许自己休息一下吧。"},
		{Title: "和人聊聊", Tip: "今天联系一位朋友或家人。社交联系是心理健康最强的保护因素之一。"},
		{Title: "正念时刻", Tip: "暂停60秒。注意你能看到的5样东西、听到的4样、触碰到的3样、闻到的2样、尝到的1样。"},
		{Title: "睡眠很重要", Tip: "尽量保持一致的睡眠时间表。优质的睡眠对情绪调节和思维清晰至关重要。"},
		{Title: "善待自己", Tip: "像和好朋友说话一样和自己说话。自我同情能减少焦虑，建立情感韧性。"},
		{Title: "设立界限", Tip: "拒绝也没关系。保护你的精力不是自私——这是你幸福所必需的。"},
		{Title: "庆祝小胜利", Tip: "你起床了？吃饭了？完成了一项任务？每一小步都很重要。进步就是进步，无论多小。"},
	},
}

// mentalHealthEmoji returns a themed emoji for mental health embeds.
var mentalHealthEmojis = []string{"🧠", "💚", "🌿", "🌻", "🕊️", "💆", "🧘", "🌈", "✨", "🫂", "💛", "🌸"}

// GetMentalHealthTip returns a random mental health tip in the given language.
func (s *FunService) GetMentalHealthTip(lang string) (string, string) {
	tips, ok := mentalHealthTips[lang]
	if !ok {
		tips = mentalHealthTips["EN"]
	}
	tip := tips[s.rng.Intn(len(tips))]
	emoji := mentalHealthEmojis[s.rng.Intn(len(mentalHealthEmojis))]
	return emoji + " " + tip.Title, tip.Tip
}

func (s *FunService) sendScheduledMentalHealth(channelID, tag, lang, guildID string) {
	if lang == "" {
		lang = "EN"
	}
	title, tip := s.GetMentalHealthTip(lang)

	embed := &FunEmbed{
		Title:       title,
		Description: tip,
		Footer:      "Take care of your mental health 💚",
		Color:       0x57F287, // green
		Content:     tag,      // mention string (@everyone, <@&roleID>, <@userID>) or empty
	}
	if s.sendFn != nil {
		if err := s.sendFn(channelID, embed); err != nil {
			s.handleSendError(guildID, "mentalhealth", channelID, err)
		}
	}
}

// handleSendError checks if a send failure is due to a deleted/unknown channel
// and auto-clears the channel config to prevent recurring errors.
func (s *FunService) handleSendError(guildID, feature, channelID string, err error) {
	errMsg := err.Error()
	if !strings.Contains(errMsg, "Unknown Channel") && !strings.Contains(errMsg, "404") {
		s.logger.Warn("Failed to send scheduled fun message", "feature", feature, "channel", channelID, "error", err)
		return
	}

	s.logger.Warn("Channel no longer exists, auto-disabling scheduled feature",
		"feature", feature, "guild", guildID, "channel", channelID)

	cfg, getErr := s.repo.Get(guildID)
	if getErr != nil || cfg == nil {
		return
	}

	switch feature {
	case "dadjoke":
		cfg.DadJokeChannelID = ""
		cfg.DadJokeInterval = 0
	case "meme":
		cfg.MemeChannelID = ""
		cfg.MemeInterval = 0
	case "mentalhealth":
		cfg.MentalHealthChannelID = ""
		cfg.MentalHealthInterval = 0
	}

	if saveErr := s.repo.Save(cfg); saveErr != nil {
		s.logger.Error("Failed to auto-disable channel config", "guild", guildID, "error", saveErr)
	}
}

// --- Dad Joke API ---

// icanhazdadjoke response
type dadJokeAPIResponse struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

// FetchDadJoke fetches a random clean dad joke from icanhazdadjoke.com.
func (s *FunService) FetchDadJoke() (*entity.DadJoke, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NeruBot Discord Bot (https://github.com/nerufuyo/nerubot)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dad joke: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dad joke API returned status %d", resp.StatusCode)
	}

	var apiResp dadJokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode dad joke response: %w", err)
	}

	return &entity.DadJoke{
		ID:        apiResp.ID,
		Punchline: apiResp.Joke,
		Source:    "icanhazdadjoke.com",
		FetchedAt: time.Now(),
	}, nil
}

// --- Meme API ---

// memeAPIResponse is the response from the meme-api.com endpoint.
type memeAPIResponse struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	NSFW      bool   `json:"nsfw"`
	Spoiler   bool   `json:"spoiler"`
	Author    string `json:"author"`
	Ups       int    `json:"ups"`
}

// memeAPIMultiResponse is the response for multiple memes.
type memeAPIMultiResponse struct {
	Count int               `json:"count"`
	Memes []memeAPIResponse `json:"memes"`
}

// FetchMeme fetches a random SFW meme from Reddit via meme-api.com.
func (s *FunService) FetchMeme() (*entity.Meme, error) {
	// Use clean meme subreddits only
	subreddits := []string{"memes", "dankmemes", "wholesomememes", "me_irl", "ProgrammerHumor"}
	sub := subreddits[s.rng.Intn(len(subreddits))]

	url := fmt.Sprintf("https://meme-api.com/gimme/%s", sub)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "NeruBot Discord Bot (https://github.com/nerufuyo/nerubot)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meme: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meme API returned status %d", resp.StatusCode)
	}

	var apiResp memeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode meme response: %w", err)
	}

	// Reject NSFW content
	if apiResp.NSFW || apiResp.Spoiler {
		// Try once more with wholesomememes as fallback
		return s.fetchMemeFromSubreddit("wholesomememes")
	}

	return &entity.Meme{
		Title:     apiResp.Title,
		URL:       apiResp.URL,
		PostLink:  apiResp.PostLink,
		Subreddit: apiResp.Subreddit,
		Author:    apiResp.Author,
		NSFW:      apiResp.NSFW,
		FetchedAt: time.Now(),
	}, nil
}

func (s *FunService) fetchMemeFromSubreddit(subreddit string) (*entity.Meme, error) {
	url := fmt.Sprintf("https://meme-api.com/gimme/%s", subreddit)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "NeruBot Discord Bot (https://github.com/nerufuyo/nerubot)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meme: %w", err)
	}
	defer resp.Body.Close()

	var apiResp memeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode meme response: %w", err)
	}

	if apiResp.NSFW {
		return nil, fmt.Errorf("could not find a SFW meme")
	}

	return &entity.Meme{
		Title:     apiResp.Title,
		URL:       apiResp.URL,
		PostLink:  apiResp.PostLink,
		Subreddit: apiResp.Subreddit,
		Author:    apiResp.Author,
		NSFW:      false,
		FetchedAt: time.Now(),
	}, nil
}

// --- Guild Config helpers ---

// GetGuildConfig retrieves (or creates) a guild config.
func (s *FunService) GetGuildConfig(guildID, guildName string) (*entity.GuildConfig, error) {
	cfg, err := s.repo.Get(guildID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = entity.NewGuildConfig(guildID, guildName)
	}
	return cfg, nil
}

// SaveGuildConfig persists a guild config.
func (s *FunService) SaveGuildConfig(cfg *entity.GuildConfig) error {
	return s.repo.Save(cfg)
}
