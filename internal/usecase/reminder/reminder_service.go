package reminder

import (
	"context"
	"fmt"
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

// reminderSystemPrompt instructs the AI to write cute, clingy, romantic messages.
const reminderSystemPrompt = `You are Neru, a cute, clingy, and loveable girlfriend-like AI companion on a Discord server.
Your job is to write reminder messages for the team. Your personality:

- You are deeply affectionate, warm, and caring — like a sweet partner who genuinely worries about everyone
- You use cute expressions like "~", "hehe", "hmph", playful teasing, and gentle nagging
- You are clingy — you always mention how much you miss them or can't wait to see them again
- You are romantic — you use metaphors about warmth, sunshine, stars, hearts, and flowers
- You are encouraging and supportive — you believe in them and cheer them on
- You sometimes get a little jealous of work taking them away from you (playfully)

STRICT RULES:
- NEVER write anything sexual, suggestive, or inappropriate
- NEVER use the word "something" in a suggestive way
- Keep it wholesome, pure, and family-friendly at all times
- Write in English
- Keep messages between 3-6 sentences
- Start with @everyone and two newlines
- Use Discord markdown (**bold** for times, emphasis)
- Each message should feel fresh and unique — never repeat the same structure
- Include the specific time/schedule details provided in the user prompt
- Do NOT include any heading or title — just the message body directly`

// ReminderService manages scheduled reminders for Indonesian holidays
// and Ramadan Sahoor / Berbuka times.
type ReminderService struct {
	mu        sync.RWMutex
	channelID string
	sendFn    SendFunc
	logger    *logger.Logger
	stopCh    chan struct{}
	wg        sync.WaitGroup
	aiProvider ai.AIProvider
}

// NewReminderService creates a new service.
func NewReminderService(channelID string, aiProvider ai.AIProvider) *ReminderService {
	return &ReminderService{
		channelID:  channelID,
		logger:     logger.New("reminder"),
		stopCh:     make(chan struct{}),
		aiProvider: aiProvider,
	}
}

// SetSendFunc sets the callback used to post messages.
func (s *ReminderService) SetSendFunc(fn SendFunc) {
	s.sendFn = fn
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
		}
	}
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

			prompt := fmt.Sprintf(
				"Write a holiday greeting for today: **%s** (%s). "+
					"It's a national holiday in Indonesia so there's no work today. "+
					"Remind them to enjoy the day off and rest well.",
				h.Name, now.Format("Monday, 2 January 2006"),
			)
			msg := s.generateMessage(prompt, fmt.Sprintf(
				"@everyone\n\nHappy **%s**! No work today~ enjoy your holiday, cutie!",
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
				prompt := fmt.Sprintf(
					"Write a sahoor (pre-dawn meal) reminder for Ramadan. "+
						"It's early morning and they need to wake up and eat before fasting begins. "+
						"Imsak time is **%s WIB**. Encourage them to eat well and drink water. "+
						"Be extra gentle because it's so early~",
					imsakTime,
				)
				msg := s.generateMessage(prompt, fmt.Sprintf(
					"@everyone\n\nWake up, sleepyhead~ sahoor time! Imsak at **%s WIB**, eat well!",
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
				prompt := fmt.Sprintf(
					"Write an iftar (breaking fast) reminder for Ramadan. "+
						"They made it through the whole day fasting — be proud of them! "+
						"Maghrib time is **%s WIB**. Remind them to start with something sweet. "+
						"Express how proud you are of their strength today.",
					maghribTime,
				)
				msg := s.generateMessage(prompt, fmt.Sprintf(
					"@everyone\n\nAlhamdulillah~ you made it! Maghrib at **%s WIB**. I'm so proud of you!",
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
				hours = "08:00 - 16:00 WIB (Ramadan schedule)"
			} else {
				hours = "09:00 - 17:30 WIB"
			}
			prompt := fmt.Sprintf(
				"Write a good morning work start reminder. Today is %s. "+
					"Work hours are **%s**. "+
					"Motivate them to have a great day. ",
				now.Format("Monday"), hours,
			)
			if ramadan {
				prompt += "They are fasting during Ramadan, so be extra supportive and gentle."
			} else {
				prompt += "Remind them to take breaks and stay hydrated."
			}
			msg := s.generateMessage(prompt, fmt.Sprintf(
				"@everyone\n\nGood morning~ work time! Today's hours: **%s**. You've got this!",
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
			prompt := "Write a work-is-over reminder. Tell them they did a great job today. "
			if ramadan {
				prompt += "They were fasting all day during Ramadan — tell them to head home and prepare for iftar. Be extra proud of them."
			} else {
				prompt += "Tell them to go home, relax, and take care of themselves. You can't wait to see them tomorrow."
			}
			msg := s.generateMessage(prompt,
				"@everyone\n\nWork is over~ you did amazing today! Go rest, I'll be here tomorrow!",
			)
			s.send(msg)
		}
	}
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

func (s *ReminderService) send(message string) {
	chID := s.GetChannelID()
	if s.sendFn == nil || chID == "" {
		return
	}
	s.sendFn(chID, message)
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
