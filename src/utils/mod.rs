use serenity::builder::CreateEmbed;
use serenity::utils::Colour;

pub mod ai;
pub mod reminder;

pub fn success_embed(title: &str, description: &str) -> CreateEmbed {
    CreateEmbed::new()
        .title(title)
        .description(description)
        .colour(Colour::from_rgb(87, 242, 135))
        .timestamp(chrono::Utc::now())
}

pub fn error_embed(title: &str, description: &str) -> CreateEmbed {
    CreateEmbed::new()
        .title(title)
        .description(description)
        .colour(Colour::from_rgb(237, 66, 69))
        .timestamp(chrono::Utc::now())
}

pub fn info_embed(title: &str, description: &str) -> CreateEmbed {
    CreateEmbed::new()
        .title(title)
        .description(description)
        .colour(Colour::from_rgb(88, 101, 242))
        .timestamp(chrono::Utc::now())
}

pub fn random_color() -> Colour {
    let r = rand::random::<u8>();
    let g = rand::random::<u8>();
    let b = rand::random::<u8>();
    Colour::from_rgb(r, g, b)
}
