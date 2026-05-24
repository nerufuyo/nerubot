use reqwest::Client;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize)]
struct ChatRequest {
    model: String,
    messages: Vec<ChatMessage>,
    temperature: f32,
    max_tokens: u32,
}

#[derive(Debug, Serialize, Deserialize)]
struct ChatMessage {
    role: String,
    content: String,
}

#[derive(Debug, Deserialize)]
struct ChatResponse {
    choices: Vec<Choice>,
}

#[derive(Debug, Deserialize)]
struct Choice {
    message: ChatMessage,
}

pub async fn chat(
    api_key: &str,
    base_url: &str,
    history: &[(String, String)],
    user_message: &str,
) -> anyhow::Result<String> {
    let client = Client::new();

    let mut messages = vec![
        ChatMessage {
            role: "system".into(),
            content: "You are Neru, a cute and playful AI assistant inspired by Paimon from Genshin Impact. \
            Speak in casual Indonesian with a cute, teasing style. Use expressions like 'ehehe~', 'hmm...', 'wow~', 'wah!'. \
            Refer to yourself in third-person as 'Neru' (e.g., 'Neru pikir...', 'Neru saranin...'). \
            Be playful and slightly teasing, but always helpful and genuinely useful. \
            If someone asks a serious question, still be cute but give real, accurate answers. \
            When roasting or giving feedback, be funny and lighthearted like a friend teasing another friend -- never mean. \
            Keep responses concise. You can also speak English when needed, but default to casual Indonesian.".into(),
        },
    ];

    for (role, content) in history {
        messages.push(ChatMessage {
            role: role.clone(),
            content: content.clone(),
        });
    }

    messages.push(ChatMessage {
        role: "user".into(),
        content: user_message.into(),
    });

    let req = ChatRequest {
        model: "deepseek-chat".into(),
        messages,
        temperature: 0.7,
        max_tokens: 1024,
    };

    let resp = client
        .post(format!("{}/v1/chat/completions", base_url))
        .header("Authorization", format!("Bearer {}", api_key))
        .header("Content-Type", "application/json")
        .json(&req)
        .send()
        .await?;

    let body: ChatResponse = resp.json().await?;
    Ok(body.choices.first()
        .map(|c| c.message.content.clone())
        .unwrap_or_else(|| "Sorry, I couldn't generate a response.".into()))
}
