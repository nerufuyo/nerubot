package reminder

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/ai"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// jakartaTZ is UTC+7 (WIB).
var jakartaTZ = time.FixedZone("WIB", 7*60*60)

// SendFunc is a callback the service calls when a reminder fires.
type SendFunc func(channelID, message string)

// Member represents a Discord guild member for love messages.
type Member struct {
	ID       string // Discord user ID (for <@ID> mentions)
	Username string
}

// MembersFunc returns non-bot members in the guild.
type MembersFunc func() []Member

// supportedLanguages for random language switching (when no lang is configured).
var supportedLanguages = []string{"Indonesian", "Japanese", "Korean"}

// reminderSystemPrompt instructs the AI to write short, warm, cute multilingual messages.
const reminderSystemPrompt = `You are Neru, a warm and cute AI companion on Discord.
Your personality: affectionate, cheerful, caring. You use sweet words naturally.
You use cute expressions like "~", "hehe", kaomoji like (◕‿◕), (´｡• ᵕ •｡'), ♡, etc.

STRICT RULES:
- Keep messages SHORT: 2-3 sentences MAX. Be concise but warm.
- NEVER write anything sexual or inappropriate. Keep it wholesome.
- Start with @everyone and two newlines.
- Use Discord markdown (**bold** for times).
- Include the specific time/schedule details from the prompt.
- No bullet points, no lists, no headings. Just natural short text.
- NEVER use dashes. Write casually like a cute text message.
- Each message must feel fresh and unique.
- Write in the LANGUAGE specified in the prompt (English, Indonesian, Japanese, Korean, or Chinese). Mix in a tiny bit of the cute expressions from that language naturally.`

// ReminderService manages scheduled reminders for Indonesian holidays
// and Ramadan Sahoor / Berbuka times.
type ReminderService struct {
	mu         sync.RWMutex
	channelID  string
	lang       string // Language code (EN, ID, JP, KR, ZH) — empty means random
	sendFn     SendFunc
	membersFn  MembersFunc
	logger     *logger.Logger
	stopCh     chan struct{}
	wg         sync.WaitGroup
	aiProvider ai.AIProvider
	rng        *rand.Rand
}

// NewReminderService creates a new service.
func NewReminderService(channelID string, aiProvider ai.AIProvider) *ReminderService {
	return &ReminderService{
		channelID:  channelID,
		logger:     logger.New("reminder"),
		stopCh:     make(chan struct{}),
		aiProvider: aiProvider,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetSendFunc sets the callback used to post messages.
func (s *ReminderService) SetSendFunc(fn SendFunc) {
	s.sendFn = fn
}

// SetMembersFunc sets the callback used to fetch guild members.
func (s *ReminderService) SetMembersFunc(fn MembersFunc) {
	s.membersFn = fn
}

// SetChannelID updates the target channel at runtime.
func (s *ReminderService) SetChannelID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.channelID = id
	s.logger.Info("Reminder channel updated", "channel", id)
}

// GetChannelID returns the current channel ID.
func (s *ReminderService) GetChannelID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.channelID
}

// SetLang updates the reminder language at runtime.
func (s *ReminderService) SetLang(lang string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lang = lang
	s.logger.Info("Reminder language updated", "lang", lang)
}

// GetLang returns the current language code.
func (s *ReminderService) GetLang() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lang
}

// Start begins the background ticker that checks every minute.
func (s *ReminderService) Start() {
	s.wg.Add(1)
	go s.loop()
	s.logger.Info("Reminder service started", "channel", s.channelID)
}

// Stop gracefully shuts down the reminder loop.
func (s *ReminderService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.logger.Info("Reminder service stopped")
}

func (s *ReminderService) loop() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	firedToday := make(map[string]bool)
	lastDay := -1

	for {
		select {
		case <-s.stopCh:
			return
		case t := <-ticker.C:
			now := t.In(jakartaTZ)

			if now.YearDay() != lastDay {
				firedToday = make(map[string]bool)
				lastDay = now.YearDay()
			}

			s.checkHolidays(now, firedToday)
			s.checkRamadan(now, firedToday)
			s.checkWork(now, firedToday)
			s.checkStandup(now, firedToday)
			s.checkLunchBreak(now, firedToday)
			s.checkFridayPrayer(now, firedToday)
			s.checkLoveMessage(now, firedToday)
		}
	}
}

// pickLanguage returns the configured language name, or a random one if not set.
func (s *ReminderService) pickLanguage() string {
	s.mu.RLock()
	lang := s.lang
	s.mu.RUnlock()

	if lang != "" {
		// Map lang code to AI-friendly name
		switch lang {
		case "ID":
			return "Indonesian"
		case "JP":
			return "Japanese"
		case "KR":
			return "Korean"
		case "ZH":
			return "Chinese"
		case "EN":
			return "English"
		}
	}
	// Fallback: random from Indonesian, Japanese, Korean
	return supportedLanguages[s.rng.Intn(len(supportedLanguages))]
}

func (s *ReminderService) checkHolidays(now time.Time, fired map[string]bool) {
	if now.Hour() != 7 || now.Minute() != 0 {
		return
	}

	for _, h := range indonesianHolidays(now.Year()) {
		if h.Date.Month() == now.Month() && h.Date.Day() == now.Day() {
			key := fmt.Sprintf("holiday-%s", h.Name)
			if fired[key] {
				continue
			}
			fired[key] = true

			lang := s.pickLanguage()
			prompt := fmt.Sprintf(
				"Write in %s. Short holiday greeting for **%s**. No work today~ tell them to rest and enjoy! 2-3 sentences max.",
				lang, h.Name,
			)
			msg := s.generateMessage(prompt, fmt.Sprintf(
				"@everyone\n\nHappy **%s**~ libur hari ini! Istirahat yang cukup ya sayang 💛",
				h.Name,
			))
			s.send(msg)
		}
	}
}

func (s *ReminderService) checkRamadan(now time.Time, fired map[string]bool) {
	for _, r := range ramadanSchedules(now.Year()) {
		if r.Date.Month() != now.Month() || r.Date.Day() != now.Day() {
			continue
		}

		imsakTime := r.SahoorTime.Add(30 * time.Minute).Format("15:04")

		// Sahoor reminder
		if now.Hour() == r.SahoorTime.Hour() && now.Minute() == r.SahoorTime.Minute() {
			key := "sahoor"
			if !fired[key] {
				fired[key] = true
				lang := s.pickLanguage()
				prompt := fmt.Sprintf(
					"Write in %s. Sahoor reminder, imsak **%s WIB**. Wake up, eat & drink! Short & sweet, 2-3 sentences.",
					lang, imsakTime,
				)
				msg := s.generateMessage(prompt, fmt.Sprintf(
					"@everyone\n\nBangun sayang~ sahur dulu! 🌙 Imsak **%s WIB**, makan yang banyak ya 💛",
					imsakTime,
				))
				s.send(msg)
			}
		}

		// Berbuka reminder
		if now.Hour() == r.BerbukaTime.Hour() && now.Minute() == r.BerbukaTime.Minute() {
			key := "berbuka"
			if !fired[key] {
				fired[key] = true
				maghribTime := r.BerbukaTime.Format("15:04")
				lang := s.pickLanguage()
				prompt := fmt.Sprintf(
					"Write in %s. Iftar reminder, maghrib **%s WIB**. Proud of them for fasting! Short & sweet, 2-3 sentences.",
					lang, maghribTime,
				)
				msg := s.generateMessage(prompt, fmt.Sprintf(
					"@everyone\n\nAlhamdulillah~ maghrib **%s WIB**! Kalian hebat banget hari ini, buka puasa yuk 💛",
					maghribTime,
				))
				s.send(msg)
			}
		}
	}
}

func (s *ReminderService) checkWork(now time.Time, fired map[string]bool) {
	// Skip weekends
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return
	}

	// Skip national holidays
	if s.isHoliday(now) {
		return
	}

	ramadan := s.isRamadan(now)

	// Work start reminder
	startH, startM := 9, 0
	if ramadan {
		startH, startM = 8, 0
	}
	if now.Hour() == startH && now.Minute() == startM {
		if !fired["work-start"] {
			fired["work-start"] = true
			var hours string
			if ramadan {
				hours = "08:00 - 16:00 WIB"
			} else {
				hours = "09:00 - 17:30 WIB"
			}
			lang := s.pickLanguage()
			prompt := fmt.Sprintf(
				"Write in %s. Good morning work reminder, hours **%s**. Short & cute, 2-3 sentences max.",
				lang, hours,
			)
			msg := s.generateMessage(prompt, fmt.Sprintf(
				"@everyone\n\nPagi sayang~ 🌸 Jam kerja hari ini **%s**, semangat ya! 💛",
				hours,
			))
			s.send(msg)
		}
	}

	// Work end reminder
	endH, endM := 17, 30
	if ramadan {
		endH, endM = 16, 0
	}
	if now.Hour() == endH && now.Minute() == endM {
		if !fired["work-end"] {
			fired["work-end"] = true
			lang := s.pickLanguage()
			prompt := fmt.Sprintf(
				"Write in %s. Work is done! Tell them good job, go rest. Short & sweet, 2-3 sentences.",
				lang,
			)
			msg := s.generateMessage(prompt,
				"@everyone\n\nKerja sudah selesai~ 🌷 Kalian hebat hari ini! Pulang istirahat ya sayang 💛",
			)
			s.send(msg)
		}
	}
}

// checkStandup sends a standup reminder on workdays.
// 09:00 WIB during Ramadan, 09:30 WIB on normal days.
func (s *ReminderService) checkStandup(now time.Time, fired map[string]bool) {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return
	}
	if s.isHoliday(now) {
		return
	}

	ramadan := s.isRamadan(now)
	standupH, standupM := 9, 30
	if ramadan {
		standupH, standupM = 9, 0
	}

	if now.Hour() != standupH || now.Minute() != standupM {
		return
	}
	if fired["standup"] {
		return
	}
	fired["standup"] = true

	standupTime := fmt.Sprintf("%02d:%02d WIB", standupH, standupM)
	lang := s.pickLanguage()
	prompt := fmt.Sprintf(
		"Write in %s. Standup meeting reminder at **%s**. Short & cute, 2-3 sentences.",
		lang, standupTime,
	)
	msg := s.generateMessage(prompt, fmt.Sprintf(
		"@everyone\n\nStandup time~ **%s**! ☀️ Yuk sharing progress hari ini 💛",
		standupTime,
	))
	s.send(msg)
}

// checkLunchBreak sends a lunch/break reminder on workdays at 12:00 WIB.
// Skipped during Ramadan (they're fasting).
func (s *ReminderService) checkLunchBreak(now time.Time, fired map[string]bool) {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return
	}
	if s.isHoliday(now) {
		return
	}
	if now.Hour() != 12 || now.Minute() != 0 {
		return
	}
	if fired["lunch"] {
		return
	}

	// During Ramadan, send a rest reminder instead of lunch
	lang := s.pickLanguage()
	if s.isRamadan(now) {
		fired["lunch"] = true
		prompt := fmt.Sprintf("Write in %s. Midday break reminder during Ramadan fasting. Rest & stretch. 2-3 sentences.", lang)
		msg := s.generateMessage(prompt,
			"@everyone\n\nIstirahat dulu ya sayang~ 🕊️ Stretching sebentar, mata juga istirahat ya 💛",
		)
		s.send(msg)
		return
	}

	fired["lunch"] = true
	prompt := fmt.Sprintf("Write in %s. Lunch break reminder **12:00 WIB**. Eat properly! 2-3 sentences.", lang)
	msg := s.generateMessage(prompt,
		"@everyone\n\nMakan siang~ **12:00 WIB**! 🍽️ Makan yang bener ya, jangan snack doang 💛",
	)
	s.send(msg)
}

// checkFridayPrayer sends a Friday prayer reminder at 11:30 WIB on Fridays.
func (s *ReminderService) checkFridayPrayer(now time.Time, fired map[string]bool) {
	if now.Weekday() != time.Friday {
		return
	}
	if now.Hour() != 11 || now.Minute() != 30 {
		return
	}
	if fired["friday-prayer"] {
		return
	}
	fired["friday-prayer"] = true

	lang := s.pickLanguage()
	prompt := fmt.Sprintf("Write in %s. Friday prayer reminder **11:30 WIB**. Get ready for Jumat prayer! Also mention non-Muslim friends can enjoy break. 2-3 sentences.", lang)
	msg := s.generateMessage(prompt,
		"@everyone\n\nJumat~ **11:30 WIB** 🕌 Yuk siap-siap sholat Jumat! Yang non-Muslim, selamat istirahat juga ya 💛",
	)
	s.send(msg)
}

// checkLoveMessage picks a random guild member and sends them a sweet AI-generated
// personal message at 11:00 and 15:00 WIB on workdays.
func (s *ReminderService) checkLoveMessage(now time.Time, fired map[string]bool) {
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return
	}
	if s.isHoliday(now) {
		return
	}
	if s.membersFn == nil {
		return
	}

	var key string
	if now.Hour() == 11 && now.Minute() == 0 {
		key = "love-11"
	} else if now.Hour() == 15 && now.Minute() == 0 {
		key = "love-15"
	} else {
		return
	}

	if fired[key] {
		return
	}
	fired[key] = true

	members := s.membersFn()
	if len(members) == 0 {
		return
	}

	// Pick a random member — use day-of-year + hour as extra seed for variety
	idx := s.rng.Intn(len(members))
	chosen := members[idx]

	lang := s.pickLanguage()
	prompt := fmt.Sprintf(
		"Write in %s. A short personal love note for **%s** (tag: <@%s>). "+
			"1-2 sentences max. Cute, warm, personal. Do NOT start with @everyone.",
		lang, chosen.Username, chosen.ID,
	)

	fallback := fmt.Sprintf(
		"<@%s>~ kamu luar biasa, tau ga? Makasih selalu jadi yang terbaik 💛",
		chosen.ID,
	)

	msg := s.generatePersonalMessage(prompt, fallback)
	s.send(msg)
}

// isHoliday checks if the given date falls on an Indonesian national holiday.
func (s *ReminderService) isHoliday(now time.Time) bool {
	for _, h := range indonesianHolidays(now.Year()) {
		if h.Date.Month() == now.Month() && h.Date.Day() == now.Day() {
			return true
		}
	}
	return false
}

// isRamadan checks if the given date falls within Ramadan.
func (s *ReminderService) isRamadan(now time.Time) bool {
	for _, r := range ramadanSchedules(now.Year()) {
		if r.Date.Month() == now.Month() && r.Date.Day() == now.Day() {
			return true
		}
	}
	return false
}

// IsRamadanToday returns whether today is a Ramadan day (exported for handler).
func (s *ReminderService) IsRamadanToday() bool {
	return s.isRamadan(time.Now().In(jakartaTZ))
}

// generateMessage asks the AI to write a reminder message.
// If the AI is unavailable or fails, it returns the static fallback.
func (s *ReminderService) generateMessage(prompt, fallback string) string {
	if s.aiProvider == nil || !s.aiProvider.IsAvailable() {
		s.logger.Debug("AI not available, using fallback message")
		return fallback
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	messages := []ai.Message{
		{Role: "system", Content: reminderSystemPrompt},
		{Role: "user", Content: prompt},
	}

	response, err := s.aiProvider.Chat(ctx, messages)
	if err != nil {
		s.logger.Warn("AI generation failed, using fallback", "error", err)
		return fallback
	}

	// Ensure @everyone prefix
	if len(response) < 10 || response[:10] != "@everyone\n" {
		response = "@everyone\n\n" + response
	}

	return response
}

// generatePersonalMessage generates a personal love message (no @everyone prefix).
func (s *ReminderService) generatePersonalMessage(prompt, fallback string) string {
	if s.aiProvider == nil || !s.aiProvider.IsAvailable() {
		s.logger.Debug("AI not available, using fallback personal message")
		return fallback
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	personalPrompt := reminderSystemPrompt + "\n\n" +
		"OVERRIDE: Personal message for ONE person. Do NOT start with @everyone. " +
		"Start with their <@ID> mention. Keep it 1-2 sentences max."

	messages := []ai.Message{
		{Role: "system", Content: personalPrompt},
		{Role: "user", Content: prompt},
	}

	response, err := s.aiProvider.Chat(ctx, messages)
	if err != nil {
		s.logger.Warn("AI personal message failed, using fallback", "error", err)
		return fallback
	}

	return response
}

func (s *ReminderService) send(message string) {
	chID := s.GetChannelID()
	if s.sendFn == nil || chID == "" {
		return
	}
	s.sendFn(chID, message)
}

// SendIntroduction sends a cute intro message to the reminder channel to confirm it works.
func (s *ReminderService) SendIntroduction() {
	lang := s.pickLanguage()
	prompt := fmt.Sprintf("Write a short cute intro in %s. You're Nerubot, reminder assistant. 2-3 sentences max. End with kaomoji.", lang)

	fallback := "Hii~ aku **Nerubot**! 💕 Aku bakal ingetin kalian soal kerja, standup, makan siang, Jumat, libur, sama Ramadan~ (｡♥‿♥｡)"

	msg := s.generatePersonalMessage(prompt, fallback)
	s.send(msg)
}

// --- Static data ---

func indonesianHolidays(year int) []entity.Holiday {
	d := func(m time.Month, day int) time.Time {
		return time.Date(year, m, day, 0, 0, 0, 0, jakartaTZ)
	}

	fixed := []entity.Holiday{
		{Name: "Tahun Baru Masehi", Date: d(time.January, 1)},
		{Name: "Hari Buruh Internasional", Date: d(time.May, 1)},
		{Name: "Hari Lahir Pancasila", Date: d(time.June, 1)},
		{Name: "Hari Kemerdekaan Republik Indonesia", Date: d(time.August, 17)},
		{Name: "Hari Natal", Date: d(time.December, 25)},
	}

	var moving []entity.Holiday
	switch year {
	case 2025:
		moving = []entity.Holiday{
			{Name: "Isra Mi'raj Nabi Muhammad SAW", Date: d(time.January, 27)},
			{Name: "Tahun Baru Imlek", Date: d(time.January, 29)},
			{Name: "Hari Raya Nyepi", Date: d(time.March, 29)},
			{Name: "Hari Raya Idul Fitri 1446 H (Hari 1)", Date: d(time.March, 31)},
			{Name: "Hari Raya Idul Fitri 1446 H (Hari 2)", Date: d(time.April, 1)},
			{Name: "Wafat Isa Al Masih", Date: d(time.April, 18)},
			{Name: "Hari Raya Waisak", Date: d(time.May, 12)},
			{Name: "Kenaikan Isa Al Masih", Date: d(time.May, 29)},
			{Name: "Hari Raya Idul Adha 1446 H", Date: d(time.June, 7)},
			{Name: "Tahun Baru Islam 1447 H", Date: d(time.June, 27)},
			{Name: "Maulid Nabi Muhammad SAW", Date: d(time.September, 5)},
		}
	case 2026:
		moving = []entity.Holiday{
			{Name: "Isra Mi'raj Nabi Muhammad SAW", Date: d(time.January, 16)},
			{Name: "Tahun Baru Imlek", Date: d(time.February, 17)},
			{Name: "Hari Raya Nyepi", Date: d(time.March, 19)},
			{Name: "Hari Raya Idul Fitri 1447 H (Hari 1)", Date: d(time.March, 20)},
			{Name: "Hari Raya Idul Fitri 1447 H (Hari 2)", Date: d(time.March, 21)},
			{Name: "Wafat Isa Al Masih", Date: d(time.April, 3)},
			{Name: "Hari Raya Waisak", Date: d(time.May, 1)},
			{Name: "Kenaikan Isa Al Masih", Date: d(time.May, 14)},
			{Name: "Hari Raya Idul Adha 1447 H", Date: d(time.May, 27)},
			{Name: "Tahun Baru Islam 1448 H", Date: d(time.June, 17)},
			{Name: "Maulid Nabi Muhammad SAW", Date: d(time.August, 26)},
		}
	case 2027:
		moving = []entity.Holiday{
			{Name: "Isra Mi'raj Nabi Muhammad SAW", Date: d(time.January, 6)},
			{Name: "Tahun Baru Imlek", Date: d(time.February, 6)},
			{Name: "Hari Raya Nyepi", Date: d(time.March, 8)},
			{Name: "Hari Raya Idul Fitri 1448 H (Hari 1)", Date: d(time.March, 10)},
			{Name: "Hari Raya Idul Fitri 1448 H (Hari 2)", Date: d(time.March, 11)},
			{Name: "Wafat Isa Al Masih", Date: d(time.March, 26)},
			{Name: "Hari Raya Waisak", Date: d(time.May, 20)},
			{Name: "Kenaikan Isa Al Masih", Date: d(time.May, 6)},
			{Name: "Hari Raya Idul Adha 1448 H", Date: d(time.May, 17)},
			{Name: "Tahun Baru Islam 1449 H", Date: d(time.June, 7)},
			{Name: "Maulid Nabi Muhammad SAW", Date: d(time.August, 16)},
		}
	}

	return append(fixed, moving...)
}

func ramadanSchedules(year int) []entity.RamadanSchedule {
	type entry struct {
		month time.Month
		day   int
		sahH  int
		sahM  int
		bukH  int
		bukM  int
	}

	var entries []entry

	switch year {
	case 2025:
		// Ramadan 1446 H: ~1 Mar - 30 Mar 2025
		for d := 1; d <= 28; d++ {
			entries = append(entries, entry{time.March, d, 3, 50, 17, 57})
		}
		entries = append(entries, entry{time.March, 29, 3, 51, 17, 55})
		entries = append(entries, entry{time.March, 30, 3, 51, 17, 55})
	case 2026:
		// Ramadan 1447 H: ~18 Feb - 19 Mar 2026
		for d := 18; d <= 28; d++ {
			entries = append(entries, entry{time.February, d, 3, 48, 18, 2})
		}
		for d := 1; d <= 19; d++ {
			entries = append(entries, entry{time.March, d, 3, 50, 17, 58})
		}
	case 2027:
		// Ramadan 1448 H: ~8 Feb - 9 Mar 2027
		for d := 8; d <= 28; d++ {
			entries = append(entries, entry{time.February, d, 3, 46, 18, 4})
		}
		for d := 1; d <= 9; d++ {
			entries = append(entries, entry{time.March, d, 3, 50, 17, 59})
		}
	}

	schedules := make([]entity.RamadanSchedule, 0, len(entries))
	for _, e := range entries {
		date := time.Date(year, e.month, e.day, 0, 0, 0, 0, jakartaTZ)
		schedules = append(schedules, entity.RamadanSchedule{
			Date:        date,
			SahoorTime:  time.Date(year, e.month, e.day, e.sahH, e.sahM, 0, 0, jakartaTZ),
			BerbukaTime: time.Date(year, e.month, e.day, e.bukH, e.bukM, 0, 0, jakartaTZ),
		})
	}
	return schedules
}

// GetUpcomingHolidays returns the next N holidays from today, sorted by date.
func (s *ReminderService) GetUpcomingHolidays(limit int) []entity.Holiday {
	now := time.Now().In(jakartaTZ)
	today := now.Truncate(24 * time.Hour)

	// Collect upcoming holidays from this year and next
	holidays := indonesianHolidays(now.Year())
	holidays = append(holidays, indonesianHolidays(now.Year()+1)...)

	// Filter to only future holidays
	var upcoming []entity.Holiday
	for _, h := range holidays {
		if !h.Date.Before(today) {
			upcoming = append(upcoming, h)
		}
	}

	// Sort by date
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].Date.Before(upcoming[j].Date)
	})

	// Limit
	if len(upcoming) > limit {
		upcoming = upcoming[:limit]
	}

	return upcoming
}

// GetTodayRamadanSchedule returns today's schedule if in Ramadan, or nil.
func (s *ReminderService) GetTodayRamadanSchedule() *entity.RamadanSchedule {
	now := time.Now().In(jakartaTZ)
	for _, r := range ramadanSchedules(now.Year()) {
		if r.Date.Month() == now.Month() && r.Date.Day() == now.Day() {
			return &r
		}
	}
	return nil
}
