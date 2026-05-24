use serenity::all::*;
use sqlx::PgPool;
use crate::utils::info_embed;

pub async fn stats(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.map(|g| g.get() as i64).unwrap_or(0);

    let total_messages: (i64,) = sqlx::query_as(
        "SELECT COALESCE(SUM(message_count)::BIGINT, 0) FROM message_stats WHERE guild_id = $1"
    ).bind(guild_id).fetch_one(pool).await?;

    let active_users: (i64,) = sqlx::query_as(
        "SELECT COUNT(DISTINCT user_id) FROM message_stats WHERE guild_id = $1 AND date >= CURRENT_DATE - INTERVAL '7 days'"
    ).bind(guild_id).fetch_one(pool).await?;

    let total_commands: (i64,) = sqlx::query_as(
        "SELECT COUNT(*) FROM command_usage WHERE guild_id = $1"
    ).bind(guild_id).fetch_one(pool).await?;

    let top_commands: Vec<(String, i64)> = sqlx::query_as(
        "SELECT command, COUNT(*) as cnt FROM command_usage WHERE guild_id = $1 GROUP BY command ORDER BY cnt DESC LIMIT 5"
    ).bind(guild_id).fetch_all(pool).await?;

    let mut top_cmds_text = String::new();
    for (i, (cmd, count)) in top_commands.iter().enumerate() {
        top_cmds_text.push_str(&format!("{}. `/ {}` — {} uses\n", i + 1, cmd, count));
    }
    if top_cmds_text.is_empty() {
        top_cmds_text = "No commands used yet.".into();
    }

    let embed = info_embed("Server Statistics", &format!(
        "**Total Messages:** {}\n**Active Users (7d):** {}\n**Total Commands:** {}\n\n**Top Commands:**\n{}",
        total_messages.0, active_users.0, total_commands.0, top_cmds_text
    ));

    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}

pub async fn profile(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.map(|g| g.get() as i64).unwrap_or(0);
    let target_user = interaction.data.options.iter()
        .find(|o| o.name == "user")
        .and_then(|o| o.value.as_user_id())
        .unwrap_or(interaction.user.id);

    let user = target_user.to_user(&ctx.http).await?;

    let total_messages: (i64,) = sqlx::query_as(
        "SELECT COALESCE(SUM(message_count)::BIGINT, 0) FROM message_stats WHERE guild_id = $1 AND user_id = $2"
    ).bind(guild_id).bind(target_user.get() as i64).fetch_one(pool).await?;

    let total_commands: (i64,) = sqlx::query_as(
        "SELECT COUNT(*) FROM command_usage WHERE guild_id = $1 AND user_id = $2"
    ).bind(guild_id).bind(target_user.get() as i64).fetch_one(pool).await?;

    let warnings: (i64,) = sqlx::query_as(
        "SELECT COUNT(*) FROM warnings WHERE guild_id = $1 AND user_id = $2 AND active = true"
    ).bind(guild_id).bind(target_user.get() as i64).fetch_one(pool).await?;

    let embed = CreateEmbed::new()
        .title(format!("{}'s Profile", user.name))
        .thumbnail(user.face())
        .field("Messages", format!("{}", total_messages.0), true)
        .field("Commands Used", format!("{}", total_commands.0), true)
        .field("Active Warnings", format!("{}", warnings.0), true)
        .color(0x5865F2)
        .timestamp(chrono::Utc::now());

    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    Ok(())
}
