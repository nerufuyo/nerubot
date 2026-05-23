use serenity::all::*;
use sqlx::PgPool;
use crate::utils::{success_embed, error_embed, random_color};

const ROASTS: &[&str] = &[
    "Kamu itu kayak WiFi tetangga — kadang nyambung, kadang nggak, dan selalu bikin kesel.",
    "Kalau kepribadianmu adalah warna, kamu pasti abu-abu.",
    "Kamu bukan bodoh, kamu cuma sedikit berbeda... dari orang pintar.",
    "Aku mau roast kamu, tapi aku nggak mau bully yang lemah.",
    "Kamu kayak software update — selalu datang di waktu yang salah dan nggak ada yang minta.",
    "Kalau kamu jadi presiden, negara ini bakal rame... tapi bukan karena prestasi.",
    "Kamu itu bukan sampah, karena sampah masih bisa didaur ulang.",
    "Otakmu itu kayak browser — banyak tab terbuka tapi nggak ada yang berguna.",
    "Kamu kayak screenshot dari screenshot — makin blur, makin nggak jelas.",
    "Aku nggak bilang kamu jelek, tapi kalau ketampanan adalah kejahatan, kamu pasti orang paling jujur di dunia.",
];

pub async fn roast(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.map(|g| g.get() as i64).unwrap_or(0);
    let target = interaction.data.options.iter()
        .find(|o| o.name == "user")
        .and_then(|o| o.value.as_user_id());

    let (target_name, target_id) = if let Some(uid) = target {
        let user = uid.to_user(&ctx.http).await?;
        (user.name.clone(), uid.get() as i64)
    } else {
        (interaction.user.name.clone(), interaction.user.id.get() as i64)
    };

    let roast_text = ROASTS[rand::random::<usize>() % ROASTS.len()];
    let full_roast = format!("🔥 **{}**: {}", target_name, roast_text);

    sqlx::query("INSERT INTO roasts (guild_id, target_id, author_id, content) VALUES ($1, $2, $3, $4)")
        .bind(guild_id).bind(target_id).bind(interaction.user.id.get() as i64).bind(&full_roast)
        .execute(pool).await?;

    let embed = CreateEmbed::new()
        .title("🔥 Roasted!")
        .description(&full_roast)
        .colour(Colour::from_rgb(255, 69, 0))
        .timestamp(chrono::Utc::now());

    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}
