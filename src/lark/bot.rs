use reqwest::Client;
use serde::{Deserialize, Serialize};
use anyhow::Result;
use chrono::Datelike;

#[derive(Debug, Serialize, Deserialize)]
pub struct LarkMessage {
    pub message_id: String,
    pub chat_id: String,
    pub sender_id: String,
    pub message_type: String,
    pub content: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LarkEvent {
    pub event_type: String,
    pub event: serde_json::Value,
}

pub struct LarkBot {
    app_id: String,
    app_secret: String,
    client: Client,
    access_token: Option<String>,
    allowed_chat_ids: Vec<String>,
}

impl LarkBot {
    pub fn new(app_id: String, app_secret: String, allowed_chat_ids: Vec<String>) -> Self {
        Self {
            app_id,
            app_secret,
            client: Client::new(),
            access_token: None,
            allowed_chat_ids,
        }
    }

    pub async fn get_access_token(&mut self) -> Result<String> {
        let resp: serde_json::Value = self.client
            .post("https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal")
            .json(&serde_json::json!({
                "app_id": self.app_id,
                "app_secret": self.app_secret,
            }))
            .send().await?
            .json().await?;

        let token = resp["tenant_access_token"].as_str().unwrap_or("").to_string();
        self.access_token = Some(token.clone());
        Ok(token)
    }

    pub async fn send_message(&self, chat_id: &str, msg_type: &str, content: &str) -> Result<()> {
        let token = self.access_token.as_ref().ok_or_else(|| anyhow::anyhow!("No access token"))?;

        self.client
            .post("https://open.feishu.cn/open-apis/im/v1/messages")
            .header("Authorization", format!("Bearer {}", token))
            .query(&[("receive_id_type", "chat_id")])
            .json(&serde_json::json!({
                "receive_id": chat_id,
                "msg_type": msg_type,
                "content": content,
            }))
            .send().await?;

        Ok(())
    }

    pub async fn send_text(&self, chat_id: &str, text: &str) -> Result<()> {
        let content = serde_json::json!({"text": text}).to_string();
        self.send_message(chat_id, "text", &content).await
    }

    pub async fn send_interactive(&self, chat_id: &str, card: &serde_json::Value) -> Result<()> {
        let content = card.to_string();
        self.send_message(chat_id, "interactive", &content).await
    }

    pub async fn handle_event(&self, event: &LarkEvent) -> Result<()> {
        match event.event_type.as_str() {
            "im.message.receive_v1" => {
                let msg = &event.event;
                let chat_id = msg["chat_id"].as_str().unwrap_or("");
                let content = msg["content"].as_str().unwrap_or("");
                let msg_type = msg["message_type"].as_str().unwrap_or("text");

                if !self.allowed_chat_ids.is_empty()
                    && !self.allowed_chat_ids.iter().any(|id| id == chat_id)
                {
                    tracing::debug!("Ignoring message from unauthorized chat_id: {}", chat_id);
                    return Ok(());
                }

                if msg_type == "text" {
                    let text: serde_json::Value = serde_json::from_str(content)?;
                    let text = text["text"].as_str().unwrap_or("");

                    // Handle commands
                    if text.starts_with("/") {
                        self.handle_command(chat_id, text).await?;
                    } else {
                        // Default: echo back
                        self.send_text(chat_id, &format!("You said: {}", text)).await?;
                    }
                }
            }
            _ => {
                tracing::info!("Unhandled Lark event: {}", event.event_type);
            }
        }

        Ok(())
    }

    async fn handle_command(&self, chat_id: &str, command: &str) -> Result<()> {
        let parts: Vec<&str> = command.splitn(2, ' ').collect();
        let cmd = parts[0];
        let args = parts.get(1).unwrap_or(&"");

        match cmd {
            "/help" => {
                self.send_text(chat_id, "🤖 NeruBot Commands:\n\n/help - Show this help\n/chat <message> - Chat with AI\n/roast - Get roasted\n/stats - Server stats\n/reminder - View reminders").await?;
            }
            "/chat" => {
                if args.is_empty() {
                    self.send_text(chat_id, "Usage: /chat <message>").await?;
                } else {
                    // TODO: Integrate with AI
                    self.send_text(chat_id, &format!("AI: I received your message: {}", args)).await?;
                }
            }
            "/roast" => {
                let roasts = [
                    "Kamu itu kayak WiFi tetangga — kadang nyambung, kadang nggak.",
                    "Kalau kepribadianmu adalah warna, kamu pasti abu-abu.",
                    "Aku mau roast kamu, tapi aku nggak mau bully yang lemah.",
                ];
                let roast = roasts[rand::random::<usize>() % roasts.len()];
                self.send_text(chat_id, &format!("🔥 {}", roast)).await?;
            }
            "/stats" => {
                self.send_text(chat_id, "📊 Stats feature coming soon!").await?;
            }
            "/reminder" => {
                let today = chrono::Utc::now().date_naive();
                let holidays = crate::utils::reminder::get_indonesian_holidays(today.year());
                let upcoming: Vec<_> = holidays.iter()
                    .filter(|h| h.date >= today)
                    .take(5)
                    .collect();

                let mut text = String::from("📅 Upcoming Indonesian Holidays:\n\n");
                for h in &upcoming {
                    text.push_str(&format!("{} {} — {}\n", h.emoji, h.name, h.date.format("%d %B %Y")));
                }

                if crate::utils::reminder::is_ramadan(today) {
                    let (sahoor, berbuka) = crate::utils::reminder::get_sahoor_berbuka_times();
                    text.push_str(&format!(
                        "\n🌙 Ramadan Schedule:\n🕌 Sahoor: {}\n🌅 Berbuka: {}",
                        sahoor, berbuka
                    ));
                }

                self.send_text(chat_id, &text).await?;
            }
            _ => {
                self.send_text(chat_id, &format!("Unknown command: {}. Type /help for available commands.", cmd)).await?;
            }
        }

        Ok(())
    }
}
