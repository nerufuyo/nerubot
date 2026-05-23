-- Add reminder_type and title to reminders table
-- Types: holiday, work, standup, break, support, announcement, custom

ALTER TABLE reminders ADD COLUMN IF NOT EXISTS reminder_type TEXT NOT NULL DEFAULT 'custom'
    CHECK (reminder_type IN ('holiday', 'work', 'standup', 'break', 'support', 'announcement', 'custom'));

ALTER TABLE reminders ADD COLUMN IF NOT EXISTS title TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_reminders_type ON reminders(reminder_type, active, remind_at);
