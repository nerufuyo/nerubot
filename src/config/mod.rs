use std::env;

#[derive(Debug, Clone)]
pub struct Config {
    pub discord_token: String,
    pub database_url: String,
    pub redis_url: String,
    pub deepseek_api_key: String,
    pub deepseek_base_url: String,
    pub enable_roast: bool,
    pub enable_reminder: bool,
    pub reminder_channel_id: u64,
    pub skip_migrations: bool,
    pub log_level: String,
}

impl Config {
    pub fn from_env() -> anyhow::Result<Self> {
        dotenvy::dotenv().ok();

        Ok(Self {
            discord_token: env::var("DISCORD_TOKEN").expect("DISCORD_TOKEN required"),
            database_url: env::var("DATABASE_URL")
                .unwrap_or_else(|_| "postgres://neruwork:PNJSilkUgW7c4WEkYY6wBP1f@localhost:5432/neruwork".into()),
            redis_url: env::var("REDIS_URL")
                .unwrap_or_else(|_| "redis://:FwiAyyUEE5CuReLg4NdMZymd@localhost:6379".into()),
            deepseek_api_key: env::var("DEEPSEEK_API_KEY").unwrap_or_default(),
            deepseek_base_url: env::var("DEEPSEEK_BASE_URL")
                .unwrap_or_else(|_| "https://api.deepseek.com".into()),
            enable_roast: env::var("ENABLE_ROAST")
                .unwrap_or_else(|_| "true".into()).parse().unwrap_or(true),
            enable_reminder: env::var("ENABLE_REMINDER")
                .unwrap_or_else(|_| "true".into()).parse().unwrap_or(true),
            reminder_channel_id: env::var("REMINDER_CHANNEL_ID")
                .unwrap_or_default().parse().unwrap_or(0),
            skip_migrations: env::var("SKIP_MIGRATIONS")
                .unwrap_or_default().parse().unwrap_or(false),
            log_level: env::var("LOG_LEVEL").unwrap_or_else(|_| "info".into()),
        })
    }
}
