use chrono::{DateTime, Utc, Datelike, Timelike, Duration, NaiveDate, NaiveTime};
use sqlx::PgPool;
use std::sync::Arc;
use tokio::sync::Mutex;
use crate::lark::bot::LarkBot;

#[derive(Debug)]
struct DueReminder {
    id: i64,
    guild_id: i64,
    channel_id: i64,
    content: String,
    remind_at: DateTime<Utc>,
    recurring: Option<String>,
    reminder_type: String,
    title: String,
}

pub async fn reminder_loop(pool: PgPool, lark_bot: Arc<Mutex<LarkBot>>) {
    tracing::info!("Reminder worker started (checks every 60s)");

    loop {
        tokio::time::sleep(std::time::Duration::from_secs(60)).await;

        match process_reminders(&pool, &lark_bot).await {
            Ok(processed) => {
                if processed > 0 {
                    tracing::info!("Reminder worker processed {} reminders", processed);
                }
            }
            Err(e) => {
                tracing::error!("Reminder worker error: {}", e);
            }
        }
    }
}

async fn process_reminders(pool: &PgPool, lark_bot: &Arc<Mutex<LarkBot>>) -> Result<usize, anyhow::Error> {
    let now = Utc::now();

    // Fetch due reminders (remind_at <= now, active = true)
    let rows = sqlx::query_as::<_, (i64, i64, i64, String, DateTime<Utc>, Option<String>, String, String)>(
        "SELECT id, guild_id, channel_id, content, remind_at, recurring, reminder_type, title
         FROM reminders
         WHERE active = true AND remind_at <= $1
         ORDER BY remind_at ASC
         LIMIT 50"
    )
    .bind(now)
    .fetch_all(pool)
    .await?;

    let count = rows.len();
    let bot = lark_bot.lock().await;

    for (id, guild_id, channel_id, content, remind_at, recurring, reminder_type, title) in &rows {
        let reminder = DueReminder {
            id: *id,
            guild_id: *guild_id,
            channel_id: *channel_id,
            content: content.clone(),
            remind_at: *remind_at,
            recurring: recurring.clone(),
            reminder_type: reminder_type.clone(),
            title: title.clone(),
        };

        // Build message
        let icon = match reminder.reminder_type.as_str() {
            "holiday" => "[Holiday]",
            "work" => "[Work]",
            "standup" => "[Standup]",
            "break" => "[Break]",
            "support" => "[Support]",
            "announcement" => "[Announcement]",
            _ => "[Reminder]",
        };

        let mut message = format!(
            "{} **{}**\n\n{}\n\n_{}_",
            icon,
            if reminder.title.is_empty() { &reminder.reminder_type } else { &reminder.title },
            reminder.content,
            reminder.remind_at.format("%d %b %Y %H:%M WIB")
        );

        // Try to send via Lark bot (primary)
        // Use environment variable for Lark target chat, or try Discord channel
        let lark_chat_id = std::env::var("LARK_REMINDER_CHAT_ID").unwrap_or_default();

        if !lark_chat_id.is_empty() {
            if let Err(e) = bot.send_text(&lark_chat_id, &message).await {
                tracing::error!("Failed to send Lark reminder {}: {}", reminder.id, e);
                continue; // skip this reminder, try again next cycle
            }
        }

        if reminder.channel_id != 0 {
            // For Discord, format differently
            message = format!(
                "{} **{}**\n{}\n<t:{}:F>",
                icon,
                if reminder.title.is_empty() { &reminder.reminder_type } else { &reminder.title },
                reminder.content,
                reminder.remind_at.timestamp()
            );
            // Discord sending would go here — needs Discord HTTP client
            // For now, Discord channel-based reminders are logged
            tracing::info!("Discord reminder {} would be sent to channel {}", reminder.id, reminder.channel_id);
        }

        tracing::info!("Sent reminder #{}: {}", reminder.id, &reminder.title);

        // Handle recurring or deactivate
        if let Some(recur) = &reminder.recurring {
            if let Some(next) = calculate_next(reminder.remind_at, recur) {
                sqlx::query("UPDATE reminders SET remind_at = $1 WHERE id = $2")
                    .bind(next).bind(reminder.id)
                    .execute(pool).await?;
                tracing::info!("Rescheduled reminder #{} to {}", reminder.id, next);
            }
        } else {
            // One-time reminder — mark inactive
            sqlx::query("UPDATE reminders SET active = false WHERE id = $1")
                .bind(reminder.id)
                .execute(pool).await?;
        }
    }

    Ok(count)
}

fn calculate_next(current: DateTime<Utc>, recurring: &str) -> Option<DateTime<Utc>> {
    let naive = current.naive_utc();
    let next_naive = match recurring {
        "daily" => naive.checked_add_signed(Duration::days(1)),
        "weekly" => naive.checked_add_signed(Duration::weeks(1)),
        "monthly" => {
            let month = naive.month() as i32 + 1;
            let year = naive.year() + if month > 12 { 1 } else { 0 };
            let month = if month > 12 { month - 12 } else { month };
            NaiveDate::from_ymd_opt(year, month as u32, naive.day())
                .and_then(|d| {
                    NaiveTime::from_hms_opt(naive.hour(), naive.minute(), naive.second())
                        .map(|t| d.and_time(t))
                })
        }
        _ => None,
    };
    next_naive.map(|dt| DateTime::from_naive_utc_and_offset(dt, Utc))
}
