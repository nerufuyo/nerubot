mod commands;
mod config;
mod db;
mod handlers;
mod models;
mod utils;
mod api;
mod lark;

use serenity::all::*;
use std::sync::Arc;
use tokio::sync::Mutex;

use config::Config;
use handlers::BotData;

struct Handler;

#[async_trait]
impl EventHandler for Handler {
    async fn ready(&self, ctx: Context, ready: Ready) {
        tracing::info!("{} is connected!", ready.user.name);

        let commands = vec![
            CreateCommand::new("chat").description("Chat with AI assistant")
                .add_option(CreateCommandOption::new(CommandOptionType::String, "message", "Your message").required(true)),
            CreateCommand::new("chat-reset").description("Reset your chat history"),
            CreateCommand::new("roast").description("Roast yourself or another user")
                .add_option(CreateCommandOption::new(CommandOptionType::User, "user", "User to roast")),
            CreateCommand::new("kick").description("Kick a user").default_member_permissions(Permissions::KICK_MEMBERS)
                .add_option(CreateCommandOption::new(CommandOptionType::User, "user", "User to kick").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "reason", "Reason")),
            CreateCommand::new("ban").description("Ban a user").default_member_permissions(Permissions::BAN_MEMBERS)
                .add_option(CreateCommandOption::new(CommandOptionType::User, "user", "User to ban").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "reason", "Reason")),
            CreateCommand::new("timeout").description("Timeout a user").default_member_permissions(Permissions::MODERATE_MEMBERS)
                .add_option(CreateCommandOption::new(CommandOptionType::User, "user", "User to timeout").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::Integer, "duration", "Duration in minutes").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "reason", "Reason")),
            CreateCommand::new("warn").description("Warn a user").default_member_permissions(Permissions::MODERATE_MEMBERS)
                .add_option(CreateCommandOption::new(CommandOptionType::User, "user", "User to warn").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "reason", "Reason").required(true)),
            CreateCommand::new("purge").description("Delete messages").default_member_permissions(Permissions::MANAGE_MESSAGES)
                .add_option(CreateCommandOption::new(CommandOptionType::Integer, "amount", "Number of messages (1-100)").required(true)),
            CreateCommand::new("coinflip").description("Flip a coin"),
            CreateCommand::new("8ball").description("Ask the magic 8-ball")
                .add_option(CreateCommandOption::new(CommandOptionType::String, "question", "Your question").required(true)),
            CreateCommand::new("meme").description("Get a random meme"),
            CreateCommand::new("dad-joke").description("Get a dad joke"),
            CreateCommand::new("stats").description("View server statistics"),
            CreateCommand::new("profile").description("View user profile")
                .add_option(CreateCommandOption::new(CommandOptionType::User, "user", "User to view")),
            CreateCommand::new("reminder").description("View upcoming reminders and holidays"),
            CreateCommand::new("calc").description("Calculate a math expression")
                .add_option(CreateCommandOption::new(CommandOptionType::String, "expression", "Math expression").required(true)),
            CreateCommand::new("poll").description("Create a poll")
                .add_option(CreateCommandOption::new(CommandOptionType::String, "question", "Poll question").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "option1", "Option 1").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "option2", "Option 2").required(true))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "option3", "Option 3"))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "option4", "Option 4"))
                .add_option(CreateCommandOption::new(CommandOptionType::String, "option5", "Option 5")),
            CreateCommand::new("help").description("Show all available commands"),
        ];

        if let Err(e) = Command::set_global_commands(&ctx.http, commands).await {
            tracing::error!("Failed to register commands: {}", e);
        } else {
            tracing::info!("Slash commands registered globally");
        }
    }

    async fn interaction_create(&self, ctx: Context, interaction: Interaction) {
        let data = ctx.data.read().await;
        let bot_data = data.get::<BotDataKey>().unwrap().clone();
        if let Err(e) = handlers::handle_interaction(&ctx, interaction, bot_data).await {
            tracing::error!("Error handling interaction: {}", e);
        }
    }

    async fn message(&self, ctx: Context, msg: Message) {
        let data = ctx.data.read().await;
        let bot_data = data.get::<BotDataKey>().unwrap().clone();
        if let Err(e) = handlers::track_message(&ctx, &msg, &bot_data.pool).await {
            tracing::error!("Error tracking message: {}", e);
        }
    }
}

struct BotDataKey;
impl serenity::prelude::TypeMapKey for BotDataKey {
    type Value = Arc<BotData>;
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let config = Config::from_env()?;

    tracing_subscriber::fmt()
        .with_env_filter(&config.log_level)
        .init();

    tracing::info!("Starting NeruBot v5.0.1...");

    let pool = db::create_pool(&config.database_url).await?;
    if !config.skip_migrations {
        db::run_migrations(&pool).await?;
    } else {
        tracing::info!("Skipping migrations (SKIP_MIGRATIONS=true)");
    }

    let redis_client = redis::Client::open(config.redis_url.as_str())?;
    let mut redis_conn = redis_client.get_connection_manager().await?;
    tracing::info!("Connected to Redis");

    let bot_data = Arc::new(BotData {
        pool,
        config: config.clone(),
        spam_cache: dashmap::DashMap::new(),
    });

    // Initialize Lark bot
    let lark_config = lark::LarkConfig::from_env();
    let mut lark_bot = lark::bot::LarkBot::new(lark_config.app_id.clone(), lark_config.app_secret.clone());

    if !lark_config.app_id.is_empty() && !lark_config.app_secret.is_empty() {
        match lark_bot.get_access_token().await {
            Ok(_) => tracing::info!("Lark access token obtained"),
            Err(e) => tracing::warn!("Failed to get Lark access token: {}", e),
        }
    }

    let lark_state = Arc::new(lark::handler::LarkState {
        bot: Arc::new(Mutex::new(lark_bot)),
        verification_token: lark_config.verification_token.clone(),
    });

    // Start API server
    let api_data = bot_data.clone();
    let lark_state_clone = lark_state.clone();
    let api_handle = tokio::spawn(async move {
        use axum::{routing::post, Router};

        let lark_router: Router = Router::new()
            .route("/api/lark/webhook", post(lark::handler::handle_lark_webhook))
            .with_state(lark_state_clone);

        let app: Router = api::create_router(api_data)
            .merge(lark_router);

        let listener = tokio::net::TcpListener::bind("0.0.0.0:8082").await.unwrap();
        tracing::info!("Admin API + Lark webhook listening on 0.0.0.0:8082");
        axum::serve(listener, app).await.unwrap();
    });

    // Start Discord bot
    let intents = GatewayIntents::GUILDS
        | GatewayIntents::GUILD_MESSAGES
        | GatewayIntents::MESSAGE_CONTENT
        | GatewayIntents::GUILD_MEMBERS
        | GatewayIntents::GUILD_MESSAGE_REACTIONS;

    let mut client = Client::builder(&config.discord_token, intents)
        .event_handler(Handler)
        .await?;

    {
        let mut data = client.data.write().await;
        data.insert::<BotDataKey>(bot_data);
    }

    tracing::info!("NeruBot is starting...");

    tokio::select! {
        _ = client.start() => {
            tracing::error!("Discord client stopped");
        }
        _ = api_handle => {
            tracing::error!("API server stopped");
        }
    }

    Ok(())
}
