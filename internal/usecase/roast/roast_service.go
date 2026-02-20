package roast

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/backend"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// RoastService handles roast operations
type RoastService struct {
	repo          *repository.RoastRepository
	stats         map[string]*entity.UserRoastStats
	logger        *logger.Logger
	backendClient *backend.Client
}

// NewRoastService creates a new roast service
func NewRoastService(backendClient *backend.Client) *RoastService {
	return &RoastService{
		repo:          repository.NewRoastRepository(),
		stats:         make(map[string]*entity.UserRoastStats),
		logger:        logger.New("roast"),
		backendClient: backendClient,
	}
}

// GenerateRoast generates a roast for a user
func (s *RoastService) GenerateRoast(ctx context.Context, userID, guildID, username, lang string) (string, error) {
	// Get configurable values from dashboard settings
	cooldown := 5 * time.Minute
	minMessages := 10
	if s.backendClient != nil {
		rs := s.backendClient.GetSettings().RoastSettings
		if rs.CooldownMinutes > 0 {
			cooldown = time.Duration(rs.CooldownMinutes) * time.Minute
		}
		if rs.MinMessages > 0 {
			minMessages = rs.MinMessages
		}
	}

	// Check cooldown
	if s.isOnCooldown(userID, guildID) {
		remaining := s.getRemainingCooldown(userID, guildID)
		return "", fmt.Errorf("roast cooldown: %s remaining", remaining)
	}

	// Get or create profile
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return "", err
	}

	// Check if enough data
	if profile.MessageCount < minMessages {
		return "", fmt.Errorf("not enough data to roast. Need at least %d messages!", minMessages)
	}

	// Detect patterns
	categories := profile.DetectPatterns()
	if len(categories) == 0 {
		categories = []entity.RoastCategory{entity.RoastCategoryNormal}
	}

	// Pick a random category
	category := categories[rand.Intn(len(categories))]

	// Get patterns
	patterns, err := s.repo.GetRoastPatterns()
	if err != nil {
		return "", err
	}

	// Find matching pattern
	var roastTemplate string
	for _, pattern := range patterns {
		if pattern.Category == string(category) {
			// Pick random template
			if len(pattern.Templates) > 0 {
				roastTemplate = pattern.Templates[rand.Intn(len(pattern.Templates))]
				break
			}
		}
	}

	if roastTemplate == "" {
		roastTemplate = "You're so normal, %s, even I can't think of a good roast!"
	}

	// Format roast with username
	roast := fmt.Sprintf(roastTemplate, username)

	// Translate roast if non-English language is requested
	if lang != "" && lang != "EN" {
		roast = translateRoast(string(category), username, lang)
	}

	// Record roast
	history := &entity.RoastHistory{
		GuildID:     guildID,
		TargetID:    userID,
		RequestedBy: userID,
		Category:    category,
		Roast:       roast,
		Severity:    2,
	}
	if err := s.repo.AddRoastHistory(history); err != nil {
		s.logger.Warn("Failed to save roast history", "error", err)
	}

	// Set cooldown from dashboard settings
	s.setCooldown(userID, guildID, cooldown)

	s.logger.Info("Roast generated",
		"user", userID,
		"guild", guildID,
		"category", category,
	)

	return roast, nil
}

// TrackMessage tracks a message for roast analysis
func (s *RoastService) TrackMessage(userID, guildID, username, channelID string) error {
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return err
	}

	profile.RecordMessage(channelID, "")
	
	return s.repo.SaveUserProfile(profile)
}

// TrackVoiceActivity tracks voice activity
func (s *RoastService) TrackVoiceActivity(userID, guildID, username string, minutes int) error {
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return err
	}

	profile.RecordVoiceActivity(minutes)
	
	return s.repo.SaveUserProfile(profile)
}

// TrackCommand tracks command usage
func (s *RoastService) TrackCommand(userID, guildID, username string) error {
	profile, err := s.repo.GetOrCreateProfile(userID, guildID, username)
	if err != nil {
		return err
	}

	profile.RecordCommand()
	
	return s.repo.SaveUserProfile(profile)
}

// GetUserProfile gets a user's profile
func (s *RoastService) GetUserProfile(userID, guildID string) (*entity.UserProfile, error) {
	return s.repo.GetUserProfile(userID, guildID)
}

// GetActivityStats gets activity statistics
func (s *RoastService) GetActivityStats(userID, guildID string) (*entity.ActivityStats, error) {
	stats, err := s.repo.GetOrCreateStats(userID, guildID)
	if err != nil {
		return nil, err
	}

	// Update stats from profile
	profile, err := s.repo.GetUserProfile(userID, guildID)
	if err == nil {
		stats.TotalMessages = profile.MessageCount
		stats.TotalVoiceTime = profile.VoiceMinutes
		stats.TotalCommands = profile.CommandsUsed
		stats.LastActivity = profile.LastSeen
		stats.CalculateActivityScore()
		
		if err := s.repo.SaveActivityStats(stats); err != nil {
			s.logger.Warn("Failed to save stats", "error", err)
		}
	}

	return stats, nil
}

// GetRoastHistory gets roast history for a user
func (s *RoastService) GetRoastHistory(userID, guildID string, limit int) ([]*entity.RoastHistory, error) {
	return s.repo.GetRoastHistory(userID, guildID, limit)
}

// Cooldown helpers

func (s *RoastService) isOnCooldown(userID, guildID string) bool {
	key := userID + ":" + guildID
	stat, exists := s.stats[key]
	if !exists {
		return false
	}
	return stat.IsOnCooldown()
}

func (s *RoastService) getRemainingCooldown(userID, guildID string) time.Duration {
	key := userID + ":" + guildID
	stat, exists := s.stats[key]
	if !exists {
		return 0
	}
	return stat.RemainingCooldown()
}

func (s *RoastService) setCooldown(userID, guildID string, duration time.Duration) {
	key := userID + ":" + guildID
	stat, exists := s.stats[key]
	if !exists {
		stat = entity.NewUserRoastStats(userID, guildID)
		s.stats[key] = stat
	}
	stat.RecordRoast(entity.RoastCategoryNormal, duration)
}

// roastTranslations maps language codes to localized roast templates by category.
// Each category has a list of templates with %s as the username placeholder.
var roastTranslations = map[string]map[string][]string{
	"ID": {
		string(entity.RoastCategoryNightOwl): {
			"Emang kamu tau matahari itu bentuknya kayak apa, %s?",
			"Jadwal tidurmu lebih berantakan dari kode-mu!",
			"Vampir aja iri sama gaya hidupmu yang nocturnal.",
		},
		string(entity.RoastCategorySpammer): {
			"%s, jari-jarimu pasti capek banget ngetik terus!",
			"Iya iya, kamu punya keyboard. Gak perlu dibuktiin tiap detik!",
			"Server Discord lembur gara-gara pesanmu doang.",
		},
		string(entity.RoastCategoryLurker): {
			"%s, gue lupa lo masih di server ini!",
			"Keyboardmu rusak, atau emang terlalu keren buat ngomong sama kita?",
			"FBI bisa belajar dari teknik pengintaianmu.",
		},
		string(entity.RoastCategoryCommandSpam): {
			"%s nganggep gue asisten pribadi. Gue bukan Siri!",
			"Dibayar per command atau gimana sih?",
			"Sirkuit gue udah capek gara-gara permintaanmu yang terus-terusan!",
		},
		"normal": {
			"Kamu terlalu normal, %s, bahkan gue gak bisa mikirin roast yang bagus!",
		},
	},
	"JP": {
		string(entity.RoastCategoryNightOwl): {
			"%s、太陽がどんな形か知ってる？",
			"君の睡眠スケジュール、コードよりバグってるよ！",
			"吸血鬼も君の夜型生活に嫉妬してるよ。",
		},
		string(entity.RoastCategorySpammer): {
			"%s、その指まだ動くの？タイピングしすぎでしょ！",
			"キーボード持ってるのはわかったから、毎秒証明しなくていいよ！",
			"Discordのサーバーが君のメッセージだけで残業してるよ。",
		},
		string(entity.RoastCategoryLurker): {
			"%s、このサーバーにいたこと忘れてたわ！",
			"キーボード壊れてるの？それとも僕らと話すのがカッコ悪い？",
			"FBIも君の監視テクニックを学べるね。",
		},
		string(entity.RoastCategoryCommandSpam): {
			"%s、僕をパーソナルアシスタントだと思ってるでしょ。Siriじゃないんだけど！",
			"コマンドごとに給料もらってるの？",
			"君のリクエストが多すぎて回路が疲れちゃったよ！",
		},
		"normal": {
			"%s、普通すぎて良いロースト思いつかないよ！",
		},
	},
	"KR": {
		string(entity.RoastCategoryNightOwl): {
			"%s, 태양이 어떻게 생겼는지 알기는 해?",
			"네 수면 스케줄은 네 코드보다 더 엉망이야!",
			"뱀파이어도 네 야행성 생활을 부러워하겠다.",
		},
		string(entity.RoastCategorySpammer): {
			"%s, 손가락 안 아파? 타이핑을 너무 많이 하잖아!",
			"키보드 있는 거 알겠어. 매초마다 증명 안 해도 돼!",
			"Discord 서버가 네 메시지 때문에 야근하고 있어.",
		},
		string(entity.RoastCategoryLurker): {
			"%s, 이 서버에 있는 거 까먹고 있었어!",
			"키보드 고장났어, 아니면 우리랑 얘기하기 너무 쿨해서?",
			"FBI도 네 잠복 기술에서 배울 게 있을 거야.",
		},
		string(entity.RoastCategoryCommandSpam): {
			"%s 나를 개인 비서로 취급하네. 시리 아니거든!",
			"명령어마다 돈 받는 거야?",
			"네 끊임없는 요청에 내 회로가 지쳤어!",
		},
		"normal": {
			"너무 평범해서, %s, 좋은 로스트가 생각이 안 나!",
		},
	},
	"ZH": {
		string(entity.RoastCategoryNightOwl): {
			"%s，你知道太阳长什么样吗？",
			"你的作息比你的代码还乱！",
			"吸血鬼都羡慕你的夜猫子生活。",
		},
		string(entity.RoastCategorySpammer): {
			"%s，你的手指打字打到不累吗！",
			"好了好了，你有键盘。不用每秒都证明一次！",
			"Discord的服务器就因为你的消息在加班。",
		},
		string(entity.RoastCategoryLurker): {
			"%s，我都忘了你还在这个服务器里！",
			"你的键盘坏了，还是太酷了不屑跟我们说话？",
			"FBI都可以跟你学潜伏技术了。",
		},
		string(entity.RoastCategoryCommandSpam): {
			"%s把我当成私人助理了。我不是Siri！",
			"你是按命令数拿工资的吗？",
			"你不停的请求让我的电路都累了！",
		},
		"normal": {
			"你太普通了，%s，我都想不出好的吐槽！",
		},
	},
}

// translateRoast picks a translated roast template for the given language and category.
// If no translation exists for the category, it picks from the "normal" fallback.
func translateRoast(category, username, lang string) string {
	langTemplates, ok := roastTranslations[lang]
	if !ok {
		return fmt.Sprintf("You're so normal, %s, even I can't think of a good roast!", username)
	}

	// Try to match the exact category first
	templates, ok := langTemplates[category]
	if !ok || len(templates) == 0 {
		// Fall back to "normal" category
		templates, ok = langTemplates["normal"]
		if !ok || len(templates) == 0 {
			return fmt.Sprintf("You're so normal, %s, even I can't think of a good roast!", username)
		}
	}

	template := templates[rand.Intn(len(templates))]
	// Only format if template contains %s
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, username)
	}
	return template
}
