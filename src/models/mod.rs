use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct GuildConfig {
    pub guild_id: i64,
    pub prefix: String,
    pub language: String,
    pub welcome_channel: Option<i64>,
    pub log_channel: Option<i64>,
    pub reminder_channel: Option<i64>,
    pub enable_confession: bool,
    pub enable_roast: bool,
    pub enable_reminder: bool,
    pub enable_music: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct ChatMessage {
    pub id: i64,
    pub user_id: i64,
    pub guild_id: i64,
    pub role: String,
    pub content: String,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct Confession {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub content: String,
    pub message_id: Option<i64>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct ModLog {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub moderator_id: i64,
    pub action: String,
    pub reason: Option<String>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct Warning {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub moderator_id: i64,
    pub reason: String,
    pub active: bool,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct MessageStat {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub channel_id: i64,
    pub message_count: i64,
    pub date: chrono::NaiveDate,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct Reminder {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub channel_id: i64,
    pub content: String,
    pub remind_at: DateTime<Utc>,
    pub recurring: Option<String>,
    pub active: bool,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct Poll {
    pub id: i64,
    pub guild_id: i64,
    pub channel_id: i64,
    pub message_id: Option<i64>,
    pub author_id: i64,
    pub question: String,
    pub options: serde_json::Value,
    pub votes: serde_json::Value,
    pub active: bool,
    pub ends_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct Playlist {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub name: String,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct PlaylistTrack {
    pub id: i64,
    pub playlist_id: i64,
    pub title: String,
    pub url: String,
    pub duration: Option<i64>,
    pub position: i32,
}
