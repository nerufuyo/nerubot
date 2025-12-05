-- PostgreSQL Initialization Script for NeruBot Microservices
-- This script creates all required databases and schemas

-- Create databases for each service
CREATE DATABASE music_db;
CREATE DATABASE confession_db;
CREATE DATABASE roast_db;
CREATE DATABASE chat_db;
CREATE DATABASE news_db;
CREATE DATABASE whale_db;

-- Connect to music_db and create schema
\c music_db;

CREATE TABLE IF NOT EXISTS music_queues (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    song_url TEXT NOT NULL,
    song_title VARCHAR(255),
    requested_by VARCHAR(32),
    added_at TIMESTAMP DEFAULT NOW(),
    position INTEGER,
    status VARCHAR(20) DEFAULT 'queued',
    INDEX idx_guild_status (guild_id, status)
);

CREATE TABLE IF NOT EXISTS playback_state (
    guild_id VARCHAR(32) PRIMARY KEY,
    current_song_id INTEGER REFERENCES music_queues(id),
    is_playing BOOLEAN DEFAULT false,
    loop_mode VARCHAR(20) DEFAULT 'none',
    volume INTEGER DEFAULT 100,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Connect to confession_db and create schema
\c confession_db;

CREATE TABLE IF NOT EXISTS confessions (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    confession_number INTEGER,
    content TEXT NOT NULL,
    image_url TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    submitted_at TIMESTAMP DEFAULT NOW(),
    approved_at TIMESTAMP,
    channel_id VARCHAR(32),
    message_id VARCHAR(32),
    INDEX idx_guild_status (guild_id, status),
    INDEX idx_guild_number (guild_id, confession_number)
);

CREATE TABLE IF NOT EXISTS confession_replies (
    id SERIAL PRIMARY KEY,
    confession_id INTEGER REFERENCES confessions(id) ON DELETE CASCADE,
    reply_content TEXT NOT NULL,
    replied_at TIMESTAMP DEFAULT NOW(),
    message_id VARCHAR(32)
);

CREATE TABLE IF NOT EXISTS confession_settings (
    guild_id VARCHAR(32) PRIMARY KEY,
    enabled BOOLEAN DEFAULT true,
    mod_channel_id VARCHAR(32),
    confession_channel_id VARCHAR(32),
    require_approval BOOLEAN DEFAULT true,
    allow_images BOOLEAN DEFAULT true,
    max_length INTEGER DEFAULT 2000
);

-- Connect to roast_db and create schema
\c roast_db;

CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    user_id VARCHAR(32) NOT NULL,
    message_count INTEGER DEFAULT 0,
    reaction_count INTEGER DEFAULT 0,
    voice_minutes INTEGER DEFAULT 0,
    command_count INTEGER DEFAULT 0,
    last_seen TIMESTAMP DEFAULT NOW(),
    UNIQUE(guild_id, user_id),
    INDEX idx_guild_user (guild_id, user_id)
);

CREATE TABLE IF NOT EXISTS roast_history (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    user_id VARCHAR(32) NOT NULL,
    roast_content TEXT NOT NULL,
    roast_category VARCHAR(50),
    roasted_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_guild_user (guild_id, user_id)
);

CREATE TABLE IF NOT EXISTS roast_patterns (
    id SERIAL PRIMARY KEY,
    pattern_name VARCHAR(50) UNIQUE,
    pattern_description TEXT,
    min_messages INTEGER,
    min_reactions INTEGER,
    min_voice_minutes INTEGER,
    roast_templates TEXT[]
);

-- Insert default roast patterns
INSERT INTO roast_patterns (pattern_name, pattern_description, min_messages, min_reactions, min_voice_minutes, roast_templates)
VALUES 
    ('spammer', 'User who sends too many messages', 100, 0, 0, ARRAY['You type faster than you think!', 'Ever heard of quality over quantity?']),
    ('lurker', 'User who barely participates', 0, 0, 0, ARRAY['Are you even here?', 'Breaking your silence I see!']),
    ('voice_addict', 'User who spends too much time in voice', 0, 0, 60, ARRAY['Living in voice chat, I see!', 'Do you sleep in VC?']),
    ('reactor', 'User who reacts a lot but doesnt talk', 0, 50, 0, ARRAY['All reactions, no action!', 'Use your words!']),
    ('ghost', 'Inactive user', 0, 0, 0, ARRAY['Long time no see!', 'Back from the dead?']),
    ('balanced', 'Well-balanced activity', 20, 10, 10, ARRAY['Perfectly balanced, as all things should be!', 'You actually have a life!']),
    ('newbie', 'New user with low activity', 0, 0, 0, ARRAY['Fresh meat!', 'Welcome to the chaos!']),
    ('commander', 'User who uses many commands', 0, 0, 0, ARRAY['Bot enthusiast, eh?', 'The bot whisperer!'])
ON CONFLICT (pattern_name) DO NOTHING;

-- Connect to chat_db and create schema
\c chat_db;

CREATE TABLE IF NOT EXISTS chat_sessions (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    user_id VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    last_interaction TIMESTAMP DEFAULT NOW(),
    message_count INTEGER DEFAULT 0,
    INDEX idx_guild_user (guild_id, user_id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    provider VARCHAR(20),
    tokens_used INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS provider_stats (
    provider_name VARCHAR(20) PRIMARY KEY,
    total_requests INTEGER DEFAULT 0,
    successful_requests INTEGER DEFAULT 0,
    failed_requests INTEGER DEFAULT 0,
    total_tokens INTEGER DEFAULT 0,
    last_used TIMESTAMP
);

-- Connect to news_db and create schema
\c news_db;

CREATE TABLE IF NOT EXISTS news_sources (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    source_name VARCHAR(100),
    source_url TEXT NOT NULL,
    feed_type VARCHAR(20) DEFAULT 'rss',
    enabled BOOLEAN DEFAULT true,
    last_fetched TIMESTAMP,
    INDEX idx_guild (guild_id)
);

CREATE TABLE IF NOT EXISTS news_articles (
    id SERIAL PRIMARY KEY,
    source_id INTEGER REFERENCES news_sources(id) ON DELETE CASCADE,
    article_title VARCHAR(255),
    article_url TEXT UNIQUE,
    article_content TEXT,
    published_at TIMESTAMP,
    fetched_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS published_articles (
    id SERIAL PRIMARY KEY,
    article_id INTEGER REFERENCES news_articles(id) ON DELETE CASCADE,
    guild_id VARCHAR(32),
    channel_id VARCHAR(32),
    message_id VARCHAR(32),
    published_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_guild (guild_id)
);

-- Connect to whale_db and create schema
\c whale_db;

CREATE TABLE IF NOT EXISTS whale_transactions (
    id SERIAL PRIMARY KEY,
    blockchain VARCHAR(20),
    transaction_hash VARCHAR(128) UNIQUE,
    amount DECIMAL(20,2),
    amount_usd DECIMAL(20,2),
    from_address VARCHAR(128),
    to_address VARCHAR(128),
    timestamp TIMESTAMP,
    detected_at TIMESTAMP DEFAULT NOW(),
    INDEX idx_blockchain (blockchain),
    INDEX idx_amount_usd (amount_usd)
);

CREATE TABLE IF NOT EXISTS whale_alerts (
    id SERIAL PRIMARY KEY,
    guild_id VARCHAR(32) NOT NULL,
    channel_id VARCHAR(32),
    min_amount_usd DECIMAL(20,2) DEFAULT 1000000,
    enabled BOOLEAN DEFAULT true,
    blockchains TEXT[],
    INDEX idx_guild (guild_id)
);

CREATE TABLE IF NOT EXISTS alerted_transactions (
    transaction_id INTEGER REFERENCES whale_transactions(id) ON DELETE CASCADE,
    guild_id VARCHAR(32),
    message_id VARCHAR(32),
    alerted_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (transaction_id, guild_id)
);

-- Grant permissions (optional, for production)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO nerubot;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO nerubot;

-- Return to default database
\c nerubot;

SELECT 'Database initialization completed successfully!' AS status;
