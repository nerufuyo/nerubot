use serenity::all::*;
use sqlx::PgPool;
use crate::utils::ai;
use crate::utils::{info_embed, error_embed};

pub async fn chat(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
    deepseek_key: &str,
    deepseek_url: &str,
) -> anyhow::Result<()> {
    let user_id = interaction.user.id.get() as i64;
    let guild_id = interaction.guild_id.map(|g| g.get() as i64).unwrap_or(0);
    let message = interaction.data.options.iter()
        .find(|o| o.name == "message")
        .and_then(|o| o.value.as_str())
        .unwrap_or("");

    interaction.defer_response(&ctx.http).await?;

    // Load history
    let history: Vec<(String, String)> = sqlx::query_as::<_, (String, String)>(
        "SELECT role, content FROM chat_history
         WHERE user_id = $1 AND guild_id = $2
         ORDER BY created_at DESC LIMIT 20"
    )
    .bind(user_id)
    .bind(guild_id)
    .fetch_all(pool)
    .await?
    .into_iter()
    .rev()
    .collect();

    // Get AI response
    let response = ai::chat(deepseek_key, deepseek_url, &history, message).await?;

    // Save to history
    sqlx::query("INSERT INTO chat_history (user_id, guild_id, role, content) VALUES ($1, $2, 'user', $3)")
        .bind(user_id).bind(guild_id).bind(message).execute(pool).await?;
    sqlx::query("INSERT INTO chat_history (user_id, guild_id, role, content) VALUES ($1, $2, 'assistant', $3)")
        .bind(user_id).bind(guild_id).bind(&response).execute(pool).await?;

    // Trim old history (keep last 50)
    sqlx::query(
        "DELETE FROM chat_history WHERE user_id = $1 AND guild_id = $2
         AND id NOT IN (SELECT id FROM chat_history WHERE user_id = $1 AND guild_id = $2 ORDER BY created_at DESC LIMIT 50)"
    ).bind(user_id).bind(guild_id).execute(pool).await?;

    let embed = info_embed("💬 AI Chat", &response);
    interaction.edit_response(&ctx.http, EditInteractionResponse::new().embed(embed)).await?;
    Ok(())
}

pub async fn chat_reset(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &PgPool,
) -> anyhow::Result<()> {
    let user_id = interaction.user.id.get() as i64;
    let guild_id = interaction.guild_id.map(|g| g.get() as i64).unwrap_or(0);

    sqlx::query("DELETE FROM chat_history WHERE user_id = $1 AND guild_id = $2")
        .bind(user_id).bind(guild_id).execute(pool).await?;

    let embed = success_embed("🔄 Chat Reset", "Your conversation history has been cleared.");
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;
    Ok(())
}
