use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    Json,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

use crate::handlers::BotData;

// ── Dashboard Stats ─────────────────────────────────────
#[derive(Serialize)]
pub struct DashboardStats {
    pub total_guilds: i64,
    pub total_users: i64,
    pub total_messages: i64,
    pub total_commands: i64,
    pub active_reminders: i64,
    pub active_polls: i64,
    pub total_roasts: i64,
    pub total_warnings: i64,
}

pub async fn get_stats(
    State(data): State<Arc<BotData>>,
) -> Result<Json<DashboardStats>, StatusCode> {
    let pool = &data.pool;

    let total_guilds: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM guild_config")
        .fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let total_messages: (i64,) = sqlx::query_as("SELECT COALESCE(SUM(message_count), 0) FROM message_stats")
        .fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let total_commands: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM command_usage")
        .fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let active_reminders: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM reminders WHERE active = true")
        .fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let active_polls: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM polls WHERE active = true")
        .fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let total_roasts: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM roasts")
        .fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let total_warnings: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM warnings WHERE active = true")
        .fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(DashboardStats {
        total_guilds: total_guilds.0,
        total_users: 0, // Calculate from message_stats
        total_messages: total_messages.0,
        total_commands: total_commands.0,
        active_reminders: active_reminders.0,
        active_polls: active_polls.0,
        total_roasts: total_roasts.0,
        total_warnings: total_warnings.0,
    }))
}

pub async fn get_guild_stats(
    State(data): State<Arc<BotData>>,
    Path(guild_id): Path<i64>,
) -> Result<Json<DashboardStats>, StatusCode> {
    let pool = &data.pool;

    let total_messages: (i64,) = sqlx::query_as("SELECT COALESCE(SUM(message_count), 0) FROM message_stats WHERE guild_id = $1")
        .bind(guild_id).fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let total_commands: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM command_usage WHERE guild_id = $1")
        .bind(guild_id).fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let active_reminders: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM reminders WHERE guild_id = $1 AND active = true")
        .bind(guild_id).fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let total_roasts: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM roasts WHERE guild_id = $1")
        .bind(guild_id).fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let total_warnings: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM warnings WHERE guild_id = $1 AND active = true")
        .bind(guild_id).fetch_one(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(DashboardStats {
        total_guilds: 1,
        total_users: 0,
        total_messages: total_messages.0,
        total_commands: total_commands.0,
        active_reminders: active_reminders.0,
        active_polls: 0,
        total_roasts: total_roasts.0,
        total_warnings: total_warnings.0,
    }))
}

// ── Chat History ────────────────────────────────────────
#[derive(Deserialize)]
pub struct ChatQuery {
    pub user_id: Option<i64>,
    pub guild_id: Option<i64>,
    pub limit: Option<i64>,
}

#[derive(Serialize)]
pub struct ChatMessage {
    pub id: i64,
    pub user_id: i64,
    pub guild_id: i64,
    pub role: String,
    pub content: String,
    pub created_at: String,
}

pub async fn get_chat_history(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<Json<Vec<ChatMessage>>, StatusCode> {
    let pool = &data.pool;
    let limit = query.limit.unwrap_or(50);

    let messages = if let (Some(uid), Some(gid)) = (query.user_id, query.guild_id) {
        sqlx::query_as::<_, (i64, i64, i64, String, String, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, user_id, guild_id, role, content, created_at FROM chat_history
             WHERE user_id = $1 AND guild_id = $2 ORDER BY created_at DESC LIMIT $3"
        ).bind(uid).bind(gid).bind(limit).fetch_all(pool).await
    } else if let Some(gid) = query.guild_id {
        sqlx::query_as::<_, (i64, i64, i64, String, String, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, user_id, guild_id, role, content, created_at FROM chat_history
             WHERE guild_id = $1 ORDER BY created_at DESC LIMIT $2"
        ).bind(gid).bind(limit).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (i64, i64, i64, String, String, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, user_id, guild_id, role, content, created_at FROM chat_history
             ORDER BY created_at DESC LIMIT $1"
        ).bind(limit).fetch_all(pool).await
    };

    let messages = messages.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(messages.into_iter().map(|(id, uid, gid, role, content, created_at)| {
        ChatMessage { id, user_id: uid, guild_id: gid, role, content, created_at: created_at.to_rfc3339() }
    }).collect()))
}

pub async fn clear_chat_history(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    if let (Some(uid), Some(gid)) = (query.user_id, query.guild_id) {
        sqlx::query("DELETE FROM chat_history WHERE user_id = $1 AND guild_id = $2")
            .bind(uid).bind(gid).execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    } else if let Some(gid) = query.guild_id {
        sqlx::query("DELETE FROM chat_history WHERE guild_id = $1")
            .bind(gid).execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    }

    Ok(StatusCode::OK)
}

// ── Roasts ──────────────────────────────────────────────
#[derive(Serialize)]
pub struct RoastEntry {
    pub id: i64,
    pub guild_id: i64,
    pub target_id: i64,
    pub author_id: i64,
    pub content: String,
    pub created_at: String,
}

#[derive(Deserialize)]
pub struct CreateRoastRequest {
    pub guild_id: i64,
    pub target_id: i64,
    pub author_id: i64,
    pub content: String,
}

pub async fn get_roasts(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<Json<Vec<RoastEntry>>, StatusCode> {
    let pool = &data.pool;
    let limit = query.limit.unwrap_or(50);

    let roasts = if let Some(gid) = query.guild_id {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, target_id, author_id, content, created_at FROM roasts
             WHERE guild_id = $1 ORDER BY created_at DESC LIMIT $2"
        ).bind(gid).bind(limit).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, target_id, author_id, content, created_at FROM roasts
             ORDER BY created_at DESC LIMIT $1"
        ).bind(limit).fetch_all(pool).await
    };

    let roasts = roasts.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(roasts.into_iter().map(|(id, gid, tid, aid, content, created_at)| {
        RoastEntry { id, guild_id: gid, target_id: tid, author_id: aid, content, created_at: created_at.to_rfc3339() }
    }).collect()))
}

pub async fn create_roast(
    State(data): State<Arc<BotData>>,
    Json(req): Json<CreateRoastRequest>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("INSERT INTO roasts (guild_id, target_id, author_id, content) VALUES ($1, $2, $3, $4)")
        .bind(req.guild_id).bind(req.target_id).bind(req.author_id).bind(req.content)
        .execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::CREATED)
}

pub async fn delete_roast(
    State(data): State<Arc<BotData>>,
    Path(id): Path<i64>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("DELETE FROM roasts WHERE id = $1")
        .bind(id).execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

// ── Moderation ──────────────────────────────────────────
#[derive(Serialize)]
pub struct ModLogEntry {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub moderator_id: i64,
    pub action: String,
    pub reason: Option<String>,
    pub created_at: String,
}

#[derive(Serialize)]
pub struct WarningEntry {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub moderator_id: i64,
    pub reason: String,
    pub active: bool,
    pub created_at: String,
}

pub async fn get_mod_logs(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<Json<Vec<ModLogEntry>>, StatusCode> {
    let pool = &data.pool;
    let limit = query.limit.unwrap_or(50);

    let logs = if let Some(gid) = query.guild_id {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, Option<String>, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, user_id, moderator_id, action, reason, created_at FROM mod_logs
             WHERE guild_id = $1 ORDER BY created_at DESC LIMIT $2"
        ).bind(gid).bind(limit).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, Option<String>, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, user_id, moderator_id, action, reason, created_at FROM mod_logs
             ORDER BY created_at DESC LIMIT $1"
        ).bind(limit).fetch_all(pool).await
    };

    let logs = logs.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(logs.into_iter().map(|(id, gid, uid, mid, action, reason, created_at)| {
        ModLogEntry { id, guild_id: gid, user_id: uid, moderator_id: mid, action, reason, created_at: created_at.to_rfc3339() }
    }).collect()))
}

pub async fn get_warnings(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<Json<Vec<WarningEntry>>, StatusCode> {
    let pool = &data.pool;
    let limit = query.limit.unwrap_or(50);

    let warnings = if let Some(gid) = query.guild_id {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, bool, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, user_id, moderator_id, reason, active, created_at FROM warnings
             WHERE guild_id = $1 ORDER BY created_at DESC LIMIT $2"
        ).bind(gid).bind(limit).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, bool, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, user_id, moderator_id, reason, active, created_at FROM warnings
             ORDER BY created_at DESC LIMIT $1"
        ).bind(limit).fetch_all(pool).await
    };

    let warnings = warnings.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(warnings.into_iter().map(|(id, gid, uid, mid, reason, active, created_at)| {
        WarningEntry { id, guild_id: gid, user_id: uid, moderator_id: mid, reason, active, created_at: created_at.to_rfc3339() }
    }).collect()))
}

pub async fn resolve_warning(
    State(data): State<Arc<BotData>>,
    Path(id): Path<i64>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("UPDATE warnings SET active = false WHERE id = $1")
        .bind(id).execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

// ── Reminders ───────────────────────────────────────────
#[derive(Serialize, Deserialize)]
pub struct ReminderEntry {
    pub id: i64,
    pub guild_id: i64,
    pub user_id: i64,
    pub channel_id: i64,
    pub content: String,
    pub remind_at: String,
    pub recurring: Option<String>,
    pub active: bool,
    pub reminder_type: String,
    pub title: String,
    pub created_at: String,
}

#[derive(Deserialize)]
pub struct CreateReminderRequest {
    pub guild_id: i64,
    pub user_id: i64,
    pub channel_id: i64,
    pub content: String,
    pub remind_at: String,
    pub recurring: Option<String>,
    pub reminder_type: Option<String>,
    pub title: Option<String>,
}

#[derive(Deserialize)]
pub struct ReminderQuery {
    pub guild_id: Option<i64>,
    pub reminder_type: Option<String>,
    pub active: Option<bool>,
    pub limit: Option<i64>,
}

#[derive(Serialize)]
pub struct ReminderTypeInfo {
    pub value: String,
    pub label: String,
    pub icon: String,
    pub count: i64,
}

pub async fn get_reminders(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ReminderQuery>,
) -> Result<Json<Vec<ReminderEntry>>, StatusCode> {
    let pool = &data.pool;
    let limit = query.limit.unwrap_or(100);

    let reminders = if let Some(gid) = query.guild_id {
        if let Some(rt) = &query.reminder_type {
            sqlx::query_as::<_, (i64, i64, i64, i64, String, chrono::DateTime<chrono::Utc>, Option<String>, bool, String, String, chrono::DateTime<chrono::Utc>)>(
                "SELECT id, guild_id, user_id, channel_id, content, remind_at, recurring, active, reminder_type, title, created_at
                 FROM reminders WHERE guild_id = $1 AND reminder_type = $2 ORDER BY remind_at DESC LIMIT $3"
            ).bind(gid).bind(rt).bind(limit).fetch_all(pool).await
        } else {
            sqlx::query_as::<_, (i64, i64, i64, i64, String, chrono::DateTime<chrono::Utc>, Option<String>, bool, String, String, chrono::DateTime<chrono::Utc>)>(
                "SELECT id, guild_id, user_id, channel_id, content, remind_at, recurring, active, reminder_type, title, created_at
                 FROM reminders WHERE guild_id = $1 ORDER BY remind_at DESC LIMIT $2"
            ).bind(gid).bind(limit).fetch_all(pool).await
        }
    } else if let Some(rt) = &query.reminder_type {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, chrono::DateTime<chrono::Utc>, Option<String>, bool, String, String, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, user_id, channel_id, content, remind_at, recurring, active, reminder_type, title, created_at
             FROM reminders WHERE reminder_type = $1 ORDER BY remind_at DESC LIMIT $2"
        ).bind(rt).bind(limit).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, chrono::DateTime<chrono::Utc>, Option<String>, bool, String, String, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, user_id, channel_id, content, remind_at, recurring, active, reminder_type, title, created_at
             FROM reminders ORDER BY remind_at DESC LIMIT $1"
        ).bind(limit).fetch_all(pool).await
    };

    let reminders = reminders.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(reminders.into_iter().map(|(id, gid, uid, cid, content, remind_at, recurring, active, rt, title, created_at)| {
        ReminderEntry {
            id, guild_id: gid, user_id: uid, channel_id: cid, content,
            remind_at: remind_at.to_rfc3339(), recurring, active,
            reminder_type: rt, title,
            created_at: created_at.to_rfc3339(),
        }
    }).collect()))
}

pub async fn create_reminder(
    State(data): State<Arc<BotData>>,
    Json(req): Json<CreateReminderRequest>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    let remind_at = chrono::DateTime::parse_from_rfc3339(&req.remind_at)
        .map_err(|_| StatusCode::BAD_REQUEST)?
        .with_timezone(&chrono::Utc);

    let reminder_type = req.reminder_type.unwrap_or_else(|| "custom".into());
    let title = req.title.unwrap_or_default();

    sqlx::query(
        "INSERT INTO reminders (guild_id, user_id, channel_id, content, remind_at, recurring, reminder_type, title)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
    ).bind(req.guild_id).bind(req.user_id).bind(req.channel_id).bind(req.content)
     .bind(remind_at).bind(req.recurring).bind(reminder_type).bind(title)
     .execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::CREATED)
}

pub async fn update_reminder(
    State(data): State<Arc<BotData>>,
    Path(id): Path<i64>,
    Json(req): Json<CreateReminderRequest>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    let remind_at = chrono::DateTime::parse_from_rfc3339(&req.remind_at)
        .map_err(|_| StatusCode::BAD_REQUEST)?
        .with_timezone(&chrono::Utc);

    let reminder_type = req.reminder_type.unwrap_or_else(|| "custom".into());
    let title = req.title.unwrap_or_default();

    sqlx::query(
        "UPDATE reminders SET content = $1, remind_at = $2, recurring = $3, reminder_type = $4, title = $5,
         guild_id = $6, user_id = $7, channel_id = $8 WHERE id = $9"
    ).bind(req.content).bind(remind_at).bind(req.recurring).bind(reminder_type).bind(title)
     .bind(req.guild_id).bind(req.user_id).bind(req.channel_id).bind(id)
     .execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

pub async fn delete_reminder(
    State(data): State<Arc<BotData>>,
    Path(id): Path<i64>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("DELETE FROM reminders WHERE id = $1")
        .bind(id).execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

pub async fn toggle_reminder(
    State(data): State<Arc<BotData>>,
    Path(id): Path<i64>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("UPDATE reminders SET active = NOT active WHERE id = $1")
        .bind(id).execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

pub async fn get_reminder_types(
    State(data): State<Arc<BotData>>,
) -> Result<Json<Vec<ReminderTypeInfo>>, StatusCode> {
    let pool = &data.pool;

    let types = vec![
        ("holiday", "Holiday", ""),
        ("work", "Work", ""),
        ("standup", "Standup", ""),
        ("break", "Break", ""),
        ("support", "Support", ""),
        ("announcement", "Announcement", ""),
        ("custom", "Custom", ""),
    ];

    let mut result = Vec::new();
    for (value, label, icon) in &types {
        let count: (i64,) = sqlx::query_as(
            "SELECT COUNT(*) FROM reminders WHERE reminder_type = $1 AND active = true"
        ).bind(value).fetch_one(pool).await.unwrap_or((0,));

        result.push(ReminderTypeInfo {
            value: value.to_string(),
            label: label.to_string(),
            icon: icon.to_string(),
            count: count.0,
        });
    }

    Ok(Json(result))
}

// ── Polls ───────────────────────────────────────────────
#[derive(Serialize)]
pub struct PollEntry {
    pub id: i64,
    pub guild_id: i64,
    pub channel_id: i64,
    pub author_id: i64,
    pub question: String,
    pub options: serde_json::Value,
    pub votes: serde_json::Value,
    pub active: bool,
    pub ends_at: Option<String>,
    pub created_at: String,
}

#[derive(Deserialize)]
pub struct CreatePollRequest {
    pub guild_id: i64,
    pub channel_id: i64,
    pub author_id: i64,
    pub question: String,
    pub options: Vec<String>,
    pub ends_at: Option<String>,
}

pub async fn get_polls(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<Json<Vec<PollEntry>>, StatusCode> {
    let pool = &data.pool;
    let limit = query.limit.unwrap_or(50);

    let polls = if let Some(gid) = query.guild_id {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, serde_json::Value, serde_json::Value, bool, Option<chrono::DateTime<chrono::Utc>>, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, channel_id, author_id, question, options, votes, active, ends_at, created_at FROM polls
             WHERE guild_id = $1 ORDER BY created_at DESC LIMIT $2"
        ).bind(gid).bind(limit).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (i64, i64, i64, i64, String, serde_json::Value, serde_json::Value, bool, Option<chrono::DateTime<chrono::Utc>>, chrono::DateTime<chrono::Utc>)>(
            "SELECT id, guild_id, channel_id, author_id, question, options, votes, active, ends_at, created_at FROM polls
             ORDER BY created_at DESC LIMIT $1"
        ).bind(limit).fetch_all(pool).await
    };

    let polls = polls.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(polls.into_iter().map(|(id, gid, cid, aid, question, options, votes, active, ends_at, created_at)| {
        PollEntry { id, guild_id: gid, channel_id: cid, author_id: aid, question, options, votes, active, ends_at: ends_at.map(|e| e.to_rfc3339()), created_at: created_at.to_rfc3339() }
    }).collect()))
}

pub async fn create_poll(
    State(data): State<Arc<BotData>>,
    Json(req): Json<CreatePollRequest>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    let ends_at = req.ends_at.and_then(|e| chrono::DateTime::parse_from_rfc3339(&e).ok().map(|dt| dt.with_timezone(&chrono::Utc)));

    sqlx::query("INSERT INTO polls (guild_id, channel_id, author_id, question, options, ends_at) VALUES ($1, $2, $3, $4, $5, $6)")
        .bind(req.guild_id).bind(req.channel_id).bind(req.author_id).bind(req.question)
        .bind(serde_json::to_value(&req.options).unwrap_or_default())
        .bind(ends_at)
        .execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::CREATED)
}

pub async fn close_poll(
    State(data): State<Arc<BotData>>,
    Path(id): Path<i64>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("UPDATE polls SET active = false WHERE id = $1")
        .bind(id).execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

// ── Guild Config ────────────────────────────────────────
#[derive(Serialize)]
pub struct GuildConfigEntry {
    pub guild_id: i64,
    pub prefix: String,
    pub language: String,
    pub welcome_channel: Option<i64>,
    pub log_channel: Option<i64>,
    pub reminder_channel: Option<i64>,
    pub enable_roast: bool,
    pub enable_reminder: bool,
    pub created_at: String,
    pub updated_at: String,
}

#[derive(Deserialize)]
pub struct UpdateGuildConfig {
    pub prefix: Option<String>,
    pub language: Option<String>,
    pub welcome_channel: Option<i64>,
    pub log_channel: Option<i64>,
    pub reminder_channel: Option<i64>,
    pub enable_roast: Option<bool>,
    pub enable_reminder: Option<bool>,
}

pub async fn get_guilds(
    State(data): State<Arc<BotData>>,
) -> Result<Json<Vec<GuildConfigEntry>>, StatusCode> {
    let pool = &data.pool;

    let guilds = sqlx::query_as::<_, (i64, String, String, Option<i64>, Option<i64>, Option<i64>, bool, bool, chrono::DateTime<chrono::Utc>, chrono::DateTime<chrono::Utc>)>(
        "SELECT guild_id, prefix, language, welcome_channel, log_channel, reminder_channel, enable_roast, enable_reminder, created_at, updated_at FROM guild_config ORDER BY guild_id"
    ).fetch_all(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(guilds.into_iter().map(|(gid, prefix, lang, wc, lc, rc, roast, reminder, created, updated)| {
        GuildConfigEntry { guild_id: gid, prefix, language: lang, welcome_channel: wc, log_channel: lc, reminder_channel: rc, enable_roast: roast, enable_reminder: reminder, created_at: created.to_rfc3339(), updated_at: updated.to_rfc3339() }
    }).collect()))
}

pub async fn get_guild_config(
    State(data): State<Arc<BotData>>,
    Path(guild_id): Path<i64>,
) -> Result<Json<GuildConfigEntry>, StatusCode> {
    let pool = &data.pool;

    let guild = sqlx::query_as::<_, (i64, String, String, Option<i64>, Option<i64>, Option<i64>, bool, bool, chrono::DateTime<chrono::Utc>, chrono::DateTime<chrono::Utc>)>(
        "SELECT guild_id, prefix, language, welcome_channel, log_channel, reminder_channel, enable_roast, enable_reminder, created_at, updated_at FROM guild_config WHERE guild_id = $1"
    ).bind(guild_id).fetch_one(pool).await.map_err(|_| StatusCode::NOT_FOUND)?;

    Ok(Json(GuildConfigEntry {
        guild_id: guild.0, prefix: guild.1, language: guild.2, welcome_channel: guild.3,
        log_channel: guild.4, reminder_channel: guild.5, enable_roast: guild.6, enable_reminder: guild.7,
        created_at: guild.8.to_rfc3339(), updated_at: guild.9.to_rfc3339()
    }))
}

pub async fn update_guild_config(
    State(data): State<Arc<BotData>>,
    Path(guild_id): Path<i64>,
    Json(req): Json<UpdateGuildConfig>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query(
        "INSERT INTO guild_config (guild_id, prefix, language, welcome_channel, log_channel, reminder_channel, enable_roast, enable_reminder)
         VALUES ($1, COALESCE($2, '!'), COALESCE($3, 'id'), $4, $5, $6, COALESCE($7, true), COALESCE($8, true))
         ON CONFLICT (guild_id) DO UPDATE SET
         prefix = COALESCE(EXCLUDED.prefix, guild_config.prefix),
         language = COALESCE(EXCLUDED.language, guild_config.language),
         welcome_channel = COALESCE(EXCLUDED.welcome_channel, guild_config.welcome_channel),
         log_channel = COALESCE(EXCLUDED.log_channel, guild_config.log_channel),
         reminder_channel = COALESCE(EXCLUDED.reminder_channel, guild_config.reminder_channel),
         enable_roast = COALESCE(EXCLUDED.enable_roast, guild_config.enable_roast),
         enable_reminder = COALESCE(EXCLUDED.enable_reminder, guild_config.enable_reminder),
         updated_at = now()"
    ).bind(guild_id).bind(req.prefix).bind(req.language).bind(req.welcome_channel)
     .bind(req.log_channel).bind(req.reminder_channel).bind(req.enable_roast).bind(req.enable_reminder)
     .execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

// ── Blocked Users ───────────────────────────────────────
#[derive(Serialize)]
pub struct BlockedUser {
    pub guild_id: i64,
    pub user_id: i64,
    pub reason: Option<String>,
    pub blocked_at: String,
}

#[derive(Deserialize)]
pub struct BlockUserRequest {
    pub guild_id: i64,
    pub user_id: i64,
    pub reason: Option<String>,
}

pub async fn get_blocked_users(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<Json<Vec<BlockedUser>>, StatusCode> {
    let pool = &data.pool;

    let users = if let Some(gid) = query.guild_id {
        sqlx::query_as::<_, (i64, i64, Option<String>, chrono::DateTime<chrono::Utc>)>(
            "SELECT guild_id, user_id, reason, blocked_at FROM blocked_users WHERE guild_id = $1 ORDER BY blocked_at DESC"
        ).bind(gid).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (i64, i64, Option<String>, chrono::DateTime<chrono::Utc>)>(
            "SELECT guild_id, user_id, reason, blocked_at FROM blocked_users ORDER BY blocked_at DESC"
        ).fetch_all(pool).await
    };

    let users = users.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(users.into_iter().map(|(gid, uid, reason, blocked_at)| {
        BlockedUser { guild_id: gid, user_id: uid, reason, blocked_at: blocked_at.to_rfc3339() }
    }).collect()))
}

pub async fn block_user(
    State(data): State<Arc<BotData>>,
    Json(req): Json<BlockUserRequest>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("INSERT INTO blocked_users (guild_id, user_id, reason) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")
        .bind(req.guild_id).bind(req.user_id).bind(req.reason)
        .execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::CREATED)
}

pub async fn unblock_user(
    State(data): State<Arc<BotData>>,
    Path((guild_id, user_id)): Path<(i64, i64)>,
) -> Result<StatusCode, StatusCode> {
    let pool = &data.pool;

    sqlx::query("DELETE FROM blocked_users WHERE guild_id = $1 AND user_id = $2")
        .bind(guild_id).bind(user_id)
        .execute(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(StatusCode::OK)
}

// ── Bot Settings ────────────────────────────────────────
#[derive(Serialize, Deserialize)]
pub struct BotSettings {
    pub enable_roast: bool,
    pub enable_reminder: bool,
    pub reminder_channel_id: u64,
    pub log_level: String,
}

pub async fn get_settings(
    State(data): State<Arc<BotData>>,
) -> Result<Json<BotSettings>, StatusCode> {
    Ok(Json(BotSettings {
        enable_roast: data.config.enable_roast,
        enable_reminder: data.config.enable_reminder,
        reminder_channel_id: data.config.reminder_channel_id,
        log_level: data.config.log_level.clone(),
    }))
}

pub async fn update_settings(
    State(data): State<Arc<BotData>>,
    Json(req): Json<BotSettings>,
) -> Result<StatusCode, StatusCode> {
    // Note: In production, update the config and persist to .env or database
    Ok(StatusCode::OK)
}

// ── Command Usage ───────────────────────────────────────
#[derive(Serialize)]
pub struct CommandUsage {
    pub command: String,
    pub count: i64,
}

pub async fn get_command_usage(
    State(data): State<Arc<BotData>>,
    Query(query): Query<ChatQuery>,
) -> Result<Json<Vec<CommandUsage>>, StatusCode> {
    let pool = &data.pool;
    let limit = query.limit.unwrap_or(50);

    let usage = if let Some(gid) = query.guild_id {
        sqlx::query_as::<_, (String, i64)>(
            "SELECT command, COUNT(*) as cnt FROM command_usage WHERE guild_id = $1 GROUP BY command ORDER BY cnt DESC LIMIT $2"
        ).bind(gid).bind(limit).fetch_all(pool).await
    } else {
        sqlx::query_as::<_, (String, i64)>(
            "SELECT command, COUNT(*) as cnt FROM command_usage GROUP BY command ORDER BY cnt DESC LIMIT $1"
        ).bind(limit).fetch_all(pool).await
    };

    let usage = usage.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(usage.into_iter().map(|(command, count)| CommandUsage { command, count }).collect()))
}

pub async fn get_top_commands(
    State(data): State<Arc<BotData>>,
) -> Result<Json<Vec<CommandUsage>>, StatusCode> {
    let pool = &data.pool;

    let usage = sqlx::query_as::<_, (String, i64)>(
        "SELECT command, COUNT(*) as cnt FROM command_usage GROUP BY command ORDER BY cnt DESC LIMIT 10"
    ).fetch_all(pool).await.map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(usage.into_iter().map(|(command, count)| CommandUsage { command, count }).collect()))
}
