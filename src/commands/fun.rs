use serenity::all::*;
use crate::utils::{info_embed, random_color};

pub async fn coinflip(ctx: &Context, interaction: &CommandInteraction) -> anyhow::Result<()> {
    let result = if rand::random::<bool>() { "🪙 **Heads!**" } else { "🪙 **Tails!**" };
    let embed = info_embed("Coin Flip", result);
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;
    Ok(())
}

pub async fn eight_ball(ctx: &Context, interaction: &CommandInteraction) -> anyhow::Result<()> {
    let responses = vec![
        "✅ Yes, definitely!", "✅ It is certain.", "✅ Without a doubt.",
        "✅ You may rely on it.", "✅ As I see it, yes.", "✅ Most likely.",
        "✅ Outlook good.", "✅ Yes.", "🤔 Reply hazy, try again.",
        "🤔 Ask again later.", "🤔 Better not tell you now.",
        "🤔 Cannot predict now.", "🤔 Concentrate and ask again.",
        "❌ Don't count on it.", "❌ My reply is no.", "❌ My sources say no.",
        "❌ Outlook not so good.", "❌ Very doubtful.",
    ];
    let response = responses[rand::random::<usize>() % responses.len()];
    let embed = info_embed("🎱 Magic 8-Ball", response);
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;
    Ok(())
}

pub async fn meme(ctx: &Context, interaction: &CommandInteraction) -> anyhow::Result<()> {
    interaction.defer_response(&ctx.http).await?;

    let resp: serde_json::Value = reqwest::get("https://meme-api.com/gimme").await?.json().await?;

    let title = resp["title"].as_str().unwrap_or("Meme");
    let url = resp["url"].as_str().unwrap_or("");
    let author = resp["author"].as_str().unwrap_or("unknown");
    let subreddit = resp["subreddit"].as_str().unwrap_or("memes");

    let embed = CreateEmbed::new()
        .title(title)
        .image(url)
        .footer(CreateEmbedFooter::new(format!("by {} | r/{}", author, subreddit)))
        .colour(random_color())
        .timestamp(chrono::Utc::now());

    interaction.edit_response(&ctx.http, EditInteractionResponse::new().embed(embed)).await?;
    Ok(())
}

pub async fn dad_joke(ctx: &Context, interaction: &CommandInteraction) -> anyhow::Result<()> {
    let client = reqwest::Client::new();
    let resp = client.get("https://icanhazdadjoke.com/")
        .header("Accept", "application/json")
        .header("User-Agent", "NeruBot/5.0")
        .send().await?
        .json::<serde_json::Value>().await?;

    let joke = resp["joke"].as_str().unwrap_or("Why did the chicken cross the road? To get to the other side!");
    let embed = info_embed("😄 Dad Joke", joke);
    interaction.create_response(&ctx.http, CreateInteractionResponse::Message(
        CreateInteractionResponseMessage::new().embed(embed)
    )).await?;
    Ok(())
}
