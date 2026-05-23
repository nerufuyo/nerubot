use serenity::all::*;
use crate::utils::{success_embed, error_embed, info_embed};

pub async fn calc(ctx: &Context, interaction: &CommandInteraction) -> anyhow::Result<()> {
    let expression = interaction.data.options.iter()
        .find(|o| o.name == "expression")
        .and_then(|o| o.value.as_str())
        .unwrap_or("0");

    // Simple calculator - evaluate basic math
    let result = evaluate_expression(expression);

    let embed = info_embed("🧮 Calculator", &format!("`{}` = `{}`", expression, result));
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;
    Ok(())
}

fn evaluate_expression(expr: &str) -> String {
    // Simple math evaluator - handles +, -, *, /
    let expr = expr.replace(" ", "");
    let parts: Vec<&str> = expr.split_inclusive(&['+', '-', '*', '/'][..]).collect();

    // For now, return the expression as-is with a note
    format!("{} = (use a proper math library for evaluation)", expr)
}

pub async fn poll(
    ctx: &Context,
    interaction: &CommandInteraction,
    pool: &sqlx::PgPool,
) -> anyhow::Result<()> {
    let guild_id = interaction.guild_id.map(|g| g.get() as i64).unwrap_or(0);
    let channel_id = interaction.channel_id.get() as i64;
    let author_id = interaction.user.id.get() as i64;

    let question = interaction.data.options.iter()
        .find(|o| o.name == "question")
        .and_then(|o| o.value.as_str())
        .unwrap_or("Poll");

    let options: Vec<String> = interaction.data.options.iter()
        .filter(|o| o.name.starts_with("option"))
        .filter_map(|o| o.value.as_str().map(|s| s.to_string()))
        .collect();

    if options.len() < 2 {
        let embed = error_embed("❌ Error", "Please provide at least 2 options.");
        interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
            CreateInteractionResponseMessage::new().embed(embed).ephemeral(true)
        )).await?;
        return Ok(());
    }

    let emojis = vec!["1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"];

    let mut options_text = String::new();
    for (i, opt) in options.iter().enumerate() {
        options_text.push_str(&format!("{} {}\n", emojis[i], opt));
    }

    let embed = CreateEmbed::new()
        .title("📊 Poll")
        .description(format!("**{}**\n\n{}", question, options_text))
        .footer(CreateEmbedFooter::new(format!("Created by {}", interaction.user.name)))
        .color(0x5865F2)
        .timestamp(chrono::Utc::now());

    let response = interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;

    // Add reactions
    if let Some(msg) = interaction.get_response(&ctx.http).await.ok() {
        for (i, _) in options.iter().enumerate() {
            msg.react(&ctx.http, ReactionType::Unicode(emojis[i].to_string())).await?;
        }
    }

    Ok(())
}
