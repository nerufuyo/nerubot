use serenity::all::*;
use sqlx::PgPool;
use crate::utils::{success_embed, error_embed, random_color};

const ROASTS: &[&str] = &[
    "Ehehe~ Paimon liat kamu tuh kayak treasure chest yang belum di-unlock... penuh potensi tapi kuncinya ilang entah ke mana!",
    "Hmm, Paimon jadi penasaran deh — kamu ini adventurer beneran atau cuma NPC yang suka mondar-mandir nggak jelas ya?",
    "Wah, Paimon liat skill kamu kayak menu emergency food... ada sih, tapi nggak ada yang spesial~",
    "Paimon nggak bilang kamu lemah kok! Tapi kalau party kamu cuma berdua sama Paimon... Paimon mungkin bakal cari healer lain dulu deh~",
    "Kamu tuh kayak quest item yang inventory-nya selalu penuh — orang lain butuh, tapi nggak ada yang mau nyimpen!",
    "Ehehe~ Paimon salut sama percaya dirimu! Tapi kayaknya stat kamu lebih cocok jadi mascot daripada main DPS deh...",
    "Paimon mau ngasih tips nih: kalau hidupmu adalah gacha, kamu pasti dapet 3-star terus ya? Hmm, setidaknya konsisten!",
    "Wow~ Paimon lihat kamu tuh kayak Anemo Slime ya... kadang terbang tinggi, tapi lebih sering kena one-hit KO!",
    "Paimon pikir kamu selevel sama Paimon nih... Tapi Paimon at least bisa jadi emergency food, kamu... ehehe~",
    "Jangan sedih ya! Paimon yakin kamu punya hidden potential... cuma mungkin masih di locked di region yang belum dirilis~",
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
    let full_roast = format!("**{}**: {}", target_name, roast_text);

    sqlx::query("INSERT INTO roasts (guild_id, target_id, author_id, content) VALUES ($1, $2, $3, $4)")
        .bind(guild_id).bind(target_id).bind(interaction.user.id.get() as i64).bind(&full_roast)
        .execute(pool).await?;

    let embed = CreateEmbed::new()
        .title("Roasted!")
        .description(&full_roast)
        .color(0xFF4500)
        .timestamp(chrono::Utc::now());

    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}
