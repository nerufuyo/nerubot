use serenity::all::*;
use sqlx::PgPool;
use crate::utils::{success_embed, error_embed, info_embed};

pub async fn kick(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.unwrap();
    let moderator = &interaction.user;

    let target_user = interaction.data.options.iter()
        .find(|o| o.name == "user")
        .and_then(|o| o.value.as_user_id())
        .ok_or_else(|| anyhow::anyhow!("Missing user"))?;

    let reason = interaction.data.options.iter()
        .find(|o| o.name == "reason")
        .and_then(|o| o.value.as_str())
        .unwrap_or("No reason provided");

    guild_id.kick_with_reason(&ctx.http, target_user, reason).await?;

    sqlx::query("INSERT INTO mod_logs (guild_id, user_id, moderator_id, action, reason) VALUES ($1, $2, $3, 'kick', $4)")
        .bind(guild_id.get() as i64).bind(target_user.get() as i64).bind(moderator.id.get() as i64).bind(reason)
        .execute(pool).await?;

    let embed = success_embed("User Kicked", &format!("**{}** has been kicked.\nReason: {}", target_user, reason));
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}

pub async fn ban(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.unwrap();
    let moderator = &interaction.user;

    let target_user = interaction.data.options.iter()
        .find(|o| o.name == "user")
        .and_then(|o| o.value.as_user_id())
        .ok_or_else(|| anyhow::anyhow!("Missing user"))?;

    let reason = interaction.data.options.iter()
        .find(|o| o.name == "reason")
        .and_then(|o| o.value.as_str())
        .unwrap_or("No reason provided");

    guild_id.ban_with_reason(&ctx.http, target_user, 7, reason).await?;

    sqlx::query("INSERT INTO mod_logs (guild_id, user_id, moderator_id, action, reason) VALUES ($1, $2, $3, 'ban', $4)")
        .bind(guild_id.get() as i64).bind(target_user.get() as i64).bind(moderator.id.get() as i64).bind(reason)
        .execute(pool).await?;

    let embed = success_embed("User Banned", &format!("**{}** has been banned.\nReason: {}", target_user, reason));
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}

pub async fn timeout(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.unwrap();
    let moderator = &interaction.user;

    let target_user = interaction.data.options.iter()
        .find(|o| o.name == "user")
        .and_then(|o| o.value.as_user_id())
        .ok_or_else(|| anyhow::anyhow!("Missing user"))?;

    let duration_mins: i64 = interaction.data.options.iter()
        .find(|o| o.name == "duration")
        .and_then(|o| o.value.as_i64())
        .unwrap_or(10);

    let reason = interaction.data.options.iter()
        .find(|o| o.name == "reason")
        .and_then(|o| o.value.as_str())
        .unwrap_or("No reason provided");

    let until = chrono::Utc::now() + chrono::Duration::minutes(duration_mins);
    let mut member = guild_id.member(&ctx.http, target_user).await?;
    member.disable_communication_until_datetime(&ctx.http, Timestamp::from_unix_timestamp(until.timestamp())?).await?;

    sqlx::query("INSERT INTO mod_logs (guild_id, user_id, moderator_id, action, reason) VALUES ($1, $2, $3, 'timeout', $4)")
        .bind(guild_id.get() as i64).bind(target_user.get() as i64).bind(moderator.id.get() as i64).bind(reason)
        .execute(pool).await?;

    let embed = success_embed("User Timed Out", &format!("**{}** timed out for {} minutes.\nReason: {}", target_user, duration_mins, reason));
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}

pub async fn warn(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.unwrap();
    let moderator = &interaction.user;

    let target_user = interaction.data.options.iter()
        .find(|o| o.name == "user")
        .and_then(|o| o.value.as_user_id())
        .ok_or_else(|| anyhow::anyhow!("Missing user"))?;

    let reason = interaction.data.options.iter()
        .find(|o| o.name == "reason")
        .and_then(|o| o.value.as_str())
        .unwrap_or("No reason provided");

    sqlx::query("INSERT INTO warnings (guild_id, user_id, moderator_id, reason) VALUES ($1, $2, $3, $4)")
        .bind(guild_id.get() as i64).bind(target_user.get() as i64).bind(moderator.id.get() as i64).bind(reason)
        .execute(pool).await?;

    let count: (i64,) = sqlx::query_as("SELECT COUNT(*) FROM warnings WHERE guild_id = $1 AND user_id = $2 AND active = true")
        .bind(guild_id.get() as i64).bind(target_user.get() as i64)
        .fetch_one(pool).await?;

    let embed = success_embed("User Warned", &format!("**{}** has been warned ({} warnings).\nReason: {}", target_user, count.0, reason));
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}

pub async fn purge(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.unwrap();
    let channel_id = interaction.channel_id;
    let moderator = &interaction.user;

    let amount: u64 = interaction.data.options.iter()
        .find(|o| o.name == "amount")
        .and_then(|o| o.value.as_i64())
        .unwrap_or(10) as u64;

    let messages = channel_id.messages(&ctx.http, GetMessages::new().limit(amount as u8)).await?;
    let msg_ids: Vec<MessageId> = messages.iter().map(|m| m.id).collect();

    channel_id.delete_messages(&ctx.http, &msg_ids).await?;

    sqlx::query("INSERT INTO mod_logs (guild_id, user_id, moderator_id, action, reason) VALUES ($1, $2, $3, 'purge', $4)")
        .bind(guild_id.get() as i64).bind(0i64).bind(moderator.id.get() as i64).bind(format!("{} messages", amount))
        .execute(pool).await?;

    let embed = success_embed("Messages Purged", &format!("Deleted {} messages.", msg_ids.len()));
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed).ephemeral(true)
    )).await?;

    Ok(())
}
