use axum::{
    extract::State,
    http::StatusCode,
    Json,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::Mutex;

use super::bot::LarkBot;

#[derive(Debug, Serialize, Deserialize)]
pub struct LarkWebhookEvent {
    #[serde(rename = "type")]
    pub event_type: Option<String>,
    pub challenge: Option<String>,
    pub token: Option<String>,
    pub event: Option<serde_json::Value>,
    pub schema: Option<String>,
    pub header: Option<serde_json::Value>,
}

pub struct LarkState {
    pub bot: Arc<Mutex<LarkBot>>,
    pub verification_token: String,
}

pub async fn handle_lark_webhook(
    State(state): State<Arc<LarkState>>,
    Json(event): Json<LarkWebhookEvent>,
) -> Result<Json<serde_json::Value>, StatusCode> {
    tracing::info!("Lark webhook received: {:?}", event);

    // Handle URL verification challenge (Lark requires this for webhook setup)
    if let Some(challenge) = &event.challenge {
        tracing::info!("Lark verification challenge: {}", challenge);
        return Ok(Json(serde_json::json!({
            "challenge": challenge
        })));
    }

    // Verify token
    if let Some(token) = &event.token {
        if token != &state.verification_token {
            tracing::warn!("Invalid Lark verification token");
            return Err(StatusCode::UNAUTHORIZED);
        }
    }

    // Handle events
    if let Some(event_data) = &event.event {
        let bot = state.bot.lock().await;

        let event_type = event.header
            .as_ref()
            .and_then(|h| h["event_type"].as_str())
            .unwrap_or("")
            .to_string();

        let lark_event = super::bot::LarkEvent {
            event_type: if event_type.is_empty() {
                event.event_type.clone().unwrap_or_default()
            } else {
                event_type
            },
            event: event_data.clone(),
        };

        if let Err(e) = bot.handle_event(&lark_event).await {
            tracing::error!("Error handling Lark event: {}", e);
        }
    }

    Ok(Json(serde_json::json!({"msg": "ok"})))
}
