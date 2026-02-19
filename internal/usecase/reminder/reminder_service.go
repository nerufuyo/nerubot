package reminder

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// jakartaTZ is UTC+7 (WIB).
var jakartaTZ = time.FixedZone("WIB", 7*60*60)

// SendFunc is a callback the service calls when a reminder fires.
type SendFunc func(channelID, message string)

// ReminderService manages scheduled reminders for Indonesian holidays
// and Ramadan Sahoor / Berbuka times.
type ReminderService struct {
	mu        sync.RWMutex
	channelID string
	sendFn    SendFunc
	logger    *logger.Logger
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

// NewReminderService creates a new service.
func NewReminderService(channelID string) *ReminderService {
	return &ReminderService{
		channelID: channelID,
		logger:    logger.New("reminder"),
		stopCh:    make(chan struct{}),
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
			msg := fmt.Sprintf(
				"@everyone\n\nHappy **%s**!\n\nWishing everyone a wonderful holiday. Enjoy this national holiday with your loved ones.\n\nWarm regards from NeruBot.",
				h.Name,
			)
			s.send(msg)
		}
	}
}

func (s *ReminderService) checkRamadan(now time.Time, fired map[string]bool) {
	for _, r := range ramadanSchedules(now.Year()) {
		if r.Date.Month() != now.Month() || r.Date.Day() != now.Day() {
			continue
		}

		// Sahoor reminder
		if now.Hour() == r.SahoorTime.Hour() && now.Minute() == r.SahoorTime.Minute() {
			key := "sahoor"
			if !fired[key] {
				fired[key] = true
				msg := fmt.Sprintf(
					"@everyone\n\n"+
						"Hey sweetie... wake up, don't oversleep!\n"+
						"Sahoor time is almost over, eat something so you have the strength to fast today.\n\n"+
						"Set your intention from the heart, may today's fast be full of blessings.\n"+
						"Imsak at **%s WIB** — don't be late!\n\n"+
						"Sahoor greetings from NeruBot.",
					r.SahoorTime.Add(30*time.Minute).Format("15:04"),
				)
				s.send(msg)
			}
		}

		// Berbuka reminder
		if now.Hour() == r.BerbukaTime.Hour() && now.Minute() == r.BerbukaTime.Minute() {
			key := "berbuka"
			if !fired[key] {
				fired[key] = true
				msg := fmt.Sprintf(
					"@everyone\n\n"+
						"Alhamdulillah... you made it through today's fast!\n"+
						"Time to break your fast — start with something sweet, just like your smile.\n\n"+
						"Maghrib at **%s WIB**.\n"+
						"May today's worship be accepted. You're amazing!\n\n"+
						"Iftar greetings from NeruBot.",
					r.BerbukaTime.Format("15:04"),
				)
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
			var msg string
			if ramadan {
				msg = "@everyone\n\n" +
					"Good morning! Time to start work.\n" +
					"Ramadan work hours: **08:00 - 16:00 WIB**\n\n" +
					"Stay strong while fasting. Have a blessed day!"
			} else {
				msg = "@everyone\n\n" +
					"Good morning! Time to start work.\n" +
					"Work hours today: **09:00 - 17:30 WIB**\n\n" +
					"Have a productive day!"
			}
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
			var msg string
			if ramadan {
				msg = "@everyone\n\n" +
					"Work hours are over!\n" +
					"Time to head home and get ready for iftar.\n\n" +
					"Stay safe on the way. See you tomorrow!"
			} else {
				msg = "@everyone\n\n" +
					"Work hours are over!\n" +
					"Time to head home and rest up.\n\n" +
					"Stay safe on the way. See you tomorrow!"
			}
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
