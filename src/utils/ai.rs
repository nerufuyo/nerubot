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
            content: "You are Neru, a friendly and helpful AI assistant. You can speak Indonesian and English. Be concise and helpful.".into(),
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
