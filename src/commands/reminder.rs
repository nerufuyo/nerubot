use serenity::all::*;
use sqlx::PgPool;
use crate::utils::{success_embed, error_embed, info_embed};
use chrono::Datelike;

pub async fn reminder(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.map(|g| g.get() as i64).unwrap_or(0);
    let user_id = interaction.user.id.get() as i64;

    // Get upcoming holidays
    let today = chrono::Utc::now().date_naive();
    let holidays = crate::utils::reminder::get_indonesian_holidays(today.year());

    let upcoming: Vec<_> = holidays.iter()
        .filter(|h| h.date >= today)
        .take(5)
        .collect();

    let mut description = String::from("**📅 Upcoming Indonesian Holidays:**\n\n");
    for h in &upcoming {
        description.push_str(&format!("{} **{}** — {}\n", h.emoji, h.name, h.date.format("%d %B %Y")));
    }

    // Check Ramadan
    if crate::utils::reminder::is_ramadan(today) {
        let (sahoor, berbuka) = crate::utils::reminder::get_sahoor_berbuka_times();
        description.push_str(&format!(
            "\n🌙 **Ramadan Schedule:**\n🕌 Sahoor: {}\n🌅 Berbuka: {}",
            sahoor, berbuka
        ));
    }

    let embed = info_embed("⏰ Reminders", &description);
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}
