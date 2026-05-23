use serenity::all::*;
use sqlx::PgPool;
use std::sync::Arc;
use dashmap::DashMap;

use crate::commands::*;
use crate::config::Config;
use crate::utils::{success_embed, error_embed, info_embed};

pub struct BotData {
    pub pool: PgPool,
    pub config: Config,
    pub spam_cache: DashMap<u64, Vec<(i64, String)>>,
}

pub async fn handle_interaction(
    ctx: &Context,
    interaction: Interaction,
    data: Arc<BotData>,
) -> anyhow::Result<()> {
    match interaction {
        Interaction::Command(command) => {
            let guild_id = command.guild_id.map(|g| g.get() as i64).unwrap_or(0);
            let user_id = command.user.id.get() as i64;
            let cmd_name = command.data.name.clone();

            let _ = sqlx::query(
                "INSERT INTO command_usage (guild_id, user_id, command) VALUES ($1, $2, $3)"
            ).bind(guild_id).bind(user_id).bind(&cmd_name)
            .execute(&data.pool).await;

            match command.data.name.as_str() {
                "chat" => chat::chat(&ctx, &command, &data.pool, &data.config.deepseek_api_key, &data.config.deepseek_base_url).await?,
                "chat-reset" => chat::chat_reset(&ctx, &command, &data.pool).await?,
                "roast" => roast::roast(&ctx, &command, &data.pool).await?,
                "kick" => moderation::kick(&ctx, &command, &data.pool).await?,
                "ban" => moderation::ban(&ctx, &command, &data.pool).await?,
                "timeout" => moderation::timeout(&ctx, &command, &data.pool).await?,
                "warn" => moderation::warn(&ctx, &command, &data.pool).await?,
                "purge" => moderation::purge(&ctx, &command, &data.pool).await?,
                "coinflip" => fun::coinflip(&ctx, &command).await?,
                "8ball" => fun::eight_ball(&ctx, &command).await?,
                "meme" => fun::meme(&ctx, &command).await?,
                "dad-joke" => fun::dad_joke(&ctx, &command).await?,
                "stats" => analytics::stats(&ctx, &command, &data.pool).await?,
                "profile" => analytics::profile(&ctx, &command, &data.pool).await?,
                "reminder" => reminder::reminder(&ctx, &command, &data.pool).await?,
                "calc" => utility::calc(&ctx, &command).await?,
                "poll" => utility::poll(&ctx, &command, &data.pool).await?,
                "help" => show_help(&ctx, &command).await?,
                _ => {
                    let embed = error_embed("Unknown Command", "This command is not recognized.");
                    command.create_response(&ctx.http, CreateInteractionResponse::Message(
                        CreateInteractionResponseMessage::new().embed(embed).ephemeral(true)
                    )).await?;
                }
            }
        }
        _ => {}
    }
    Ok(())
}

async fn show_help(ctx: &Context, interaction: &CommandInteraction) -> anyhow::Result<()> {
    let embed = CreateEmbed::new()
        .title("🤖 NeruBot Help")
        .description("A feature-rich Discord bot built with Rust 🦀")
        .field("**💬 AI Chat**", "`/chat` `/chat-reset`", false)
        .field("**🔥 Roast**", "`/roast` — Humorous roasts", false)
        .field("**🛡️ Moderation**", "`/kick` `/ban` `/timeout` `/warn` `/purge`", false)
        .field("**😄 Fun**", "`/coinflip` `/8ball` `/meme` `/dad-joke`", false)
        .field("**📊 Analytics**", "`/stats` `/profile`", false)
        .field("**⏰ Reminders**", "`/reminder` — Indonesian holidays & Ramadan", false)
        .field("**🔧 Utility**", "`/calc` `/poll`", false)
        .colour(Colour::from_rgb(88, 101, 242))
        .footer(CreateEmbedFooter::new("NeruBot v5.0.1 — Built with Rust"))
        .timestamp(chrono::Utc::now());

    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed).ephemeral(true)
    )).await?;

    Ok(())
}

pub async fn track_message(
    ctx: &Context,
    msg: &Message,
    pool: &PgPool,
) -> anyhow::Result<()> {
    if msg.author.bot {
        return Ok(());
    }

    let guild_id = msg.guild_id.map(|g| g.get() as i64).unwrap_or(0);
    let user_id = msg.author.id.get() as i64;
    let channel_id = msg.channel_id.get() as i64;

    sqlx::query(
        "INSERT INTO message_stats (guild_id, user_id, channel_id, message_count, date)
         VALUES ($1, $2, $3, 1, CURRENT_DATE)
         ON CONFLICT (guild_id, user_id, channel_id, date)
         DO UPDATE SET message_count = message_stats.message_count + 1"
    ).bind(guild_id).bind(user_id).bind(channel_id)
    .execute(pool).await?;

    Ok(())
}
