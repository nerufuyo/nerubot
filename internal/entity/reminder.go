package entity

import "time"

// ReminderType distinguishes between holiday and Ramadan reminders.
type ReminderType string

const (
	ReminderHoliday ReminderType = "holiday"
	ReminderSahoor  ReminderType = "sahoor"
	ReminderBerbuka ReminderType = "berbuka"
)

// Holiday represents an Indonesian national holiday or religious celebration.
type Holiday struct {
	Name string
	Date time.Time // year-specific; month + day are used for matching
}

// RamadanSchedule holds the Sahoor and Berbuka times for a specific date.
type RamadanSchedule struct {
	Date        time.Time
	SahoorTime  time.Time // typically ~04:20 WIB
	BerbukaTime time.Time // typically ~17:55 WIB
}

// Reminder is a scheduled reminder that the bot will post.
type Reminder struct {
	Type    ReminderType
	Title   string
	Message string
	FireAt  time.Time
}
