use serenity::builder::CreateEmbed;

pub mod ai;
pub mod reminder;

pub fn success_embed(title: &str, description: &str) -> CreateEmbed {
    CreateEmbed::new()
        .title(title)
        .description(description)
        .color(0x57F287)
        .timestamp(chrono::Utc::now())
}

pub fn error_embed(title: &str, description: &str) -> CreateEmbed {
    CreateEmbed::new()
        .title(title)
        .description(description)
        .color(0xED4245)
        .timestamp(chrono::Utc::now())
}

pub fn info_embed(title: &str, description: &str) -> CreateEmbed {
    CreateEmbed::new()
        .title(title)
        .description(description)
        .color(0x5865F2)
        .timestamp(chrono::Utc::now())
}

pub fn random_color() -> u32 {
    let r = rand::random::<u8>() as u32;
    let g = rand::random::<u8>() as u32;
    let b = rand::random::<u8>() as u32;
    (r << 16) | (g << 8) | b
}
