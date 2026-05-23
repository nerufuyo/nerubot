-- NeruBot PostgreSQL Schema
-- Simplified version (no whale, music, confession, news)

-- ── Guild Config ────────────────────────────────────────
CREATE TABLE IF NOT EXISTS guild_config (
    guild_id         BIGINT PRIMARY KEY,
    prefix           TEXT NOT NULL DEFAULT '!',
    language         TEXT NOT NULL DEFAULT 'id',
    welcome_channel  BIGINT,
    log_channel      BIGINT,
    reminder_channel BIGINT,
    enable_roast     BOOLEAN NOT NULL DEFAULT true,
    enable_reminder  BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── Chat History (AI) ───────────────────────────────────
CREATE TABLE IF NOT EXISTS chat_history (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     BIGINT NOT NULL,
    guild_id    BIGINT NOT NULL,
    role        TEXT NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_chat_history_user ON chat_history(user_id, guild_id, created_at DESC);

-- ── Roasts ──────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS roasts (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    guild_id    BIGINT NOT NULL,
    target_id   BIGINT NOT NULL,
    author_id   BIGINT NOT NULL,
    content     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── Moderation ──────────────────────────────────────────
CREATE TABLE IF NOT EXISTS mod_logs (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    guild_id    BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    moderator_id BIGINT NOT NULL,
    action      TEXT NOT NULL CHECK (action IN ('kick', 'ban', 'timeout', 'warn', 'purge')),
    reason      TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_mod_logs_user ON mod_logs(guild_id, user_id, created_at DESC);

-- ── Warnings ────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS warnings (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    guild_id    BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    moderator_id BIGINT NOT NULL,
    reason      TEXT NOT NULL,
    active      BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_warnings_user ON warnings(guild_id, user_id, active);

-- ── Analytics ───────────────────────────────────────────
CREATE TABLE IF NOT EXISTS message_stats (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    guild_id    BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    channel_id  BIGINT NOT NULL,
    message_count BIGINT NOT NULL DEFAULT 1,
    date        DATE NOT NULL DEFAULT CURRENT_DATE,
    UNIQUE(guild_id, user_id, channel_id, date)
);
CREATE INDEX IF NOT EXISTS idx_stats_guild ON message_stats(guild_id, date DESC);

CREATE TABLE IF NOT EXISTS command_usage (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    guild_id    BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    command     TEXT NOT NULL,
    used_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_cmd_usage ON command_usage(guild_id, command, used_at DESC);

-- ── Reminders ───────────────────────────────────────────
CREATE TABLE IF NOT EXISTS reminders (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    guild_id    BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    channel_id  BIGINT NOT NULL,
    content     TEXT NOT NULL,
    remind_at   TIMESTAMPTZ NOT NULL,
    recurring   TEXT CHECK (recurring IN ('daily', 'weekly', 'monthly')),
    active      BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_reminders_active ON reminders(active, remind_at);

-- ── Polls ───────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS polls (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    guild_id    BIGINT NOT NULL,
    channel_id  BIGINT NOT NULL,
    message_id  BIGINT,
    author_id   BIGINT NOT NULL,
    question    TEXT NOT NULL,
    options     JSONB NOT NULL DEFAULT '[]',
    votes       JSONB NOT NULL DEFAULT '{}',
    active      BOOLEAN NOT NULL DEFAULT true,
    ends_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── Blocked Users ───────────────────────────────────────
CREATE TABLE IF NOT EXISTS blocked_users (
    guild_id    BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    reason      TEXT,
    blocked_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (guild_id, user_id)
);
