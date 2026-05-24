use axum::{
    extract::State,
    http::StatusCode,
    Json,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::Mutex;

use super::bot::LarkBot;

/// Decrypt Lark encrypted event body (AES-256-CBC)
/// Encrypted format: base64(IV[16 bytes] + ciphertext)
/// Key is SHA-256 hash of the encrypt key string from Lark Developer Console
fn decrypt_lark_body(encrypted_b64: &str, encrypt_key: &str) -> Result<String, String> {
    use aes::cipher::{BlockDecryptMut, KeyIvInit};
    use base64::Engine;
    use sha2::{Sha256, Digest};

    type Aes256CbcDec = cbc::Decryptor<aes::Aes256>;

    let engine = base64::engine::general_purpose::STANDARD;
    let encrypted = engine.decode(encrypted_b64).map_err(|e| format!("base64 decode: {e}"))?;

    if encrypted.len() < 16 {
        return Err("encrypted data too short".into());
    }

    // Derive AES-256 key: SHA-256(encrypt_key)
    let mut hasher = Sha256::new();
    hasher.update(encrypt_key.as_bytes());
    let key_hash = hasher.finalize();

    let iv: &[u8; 16] = encrypted[..16].try_into().map_err(|_| "bad iv".to_string())?;
    let ciphertext = &encrypted[16..];

    let mut buf = ciphertext.to_vec();
    let decryptor = Aes256CbcDec::new(key_hash.as_slice().into(), iv.into());
    let plaintext = decryptor
        .decrypt_padded_mut::<aes::cipher::block_padding::Pkcs7>(&mut buf)
        .map_err(|e| format!("decrypt: {e}"))?;

    String::from_utf8(plaintext.to_vec()).map_err(|e| format!("utf8: {e}"))
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LarkWebhookEvent {
    #[serde(rename = "type")]
    pub event_type: Option<String>,
    pub challenge: Option<String>,
    pub token: Option<String>,
    pub event: Option<serde_json::Value>,
    pub schema: Option<String>,
    pub header: Option<serde_json::Value>,
    pub encrypt: Option<String>,
}

pub struct LarkState {
    pub bot: Arc<Mutex<LarkBot>>,
    pub verification_token: String,
    pub encrypt_key: String,
}

pub async fn handle_lark_webhook(
    State(state): State<Arc<LarkState>>,
    body: String,
) -> Result<Json<serde_json::Value>, StatusCode> {
    // Try to parse raw body first to check for encryption
    let raw: serde_json::Value = serde_json::from_str(&body)
        .map_err(|_| StatusCode::BAD_REQUEST)?;

    // If encrypted, decrypt first
    let event: LarkWebhookEvent = if let Some(encrypted) = raw.get("encrypt").and_then(|v| v.as_str()) {
        tracing::debug!("Lark encrypted payload detected, decrypting...");
        let decrypted = decrypt_lark_body(encrypted, &state.encrypt_key)
            .map_err(|e| {
                tracing::error!("Lark decrypt failed: {}", e);
                StatusCode::BAD_REQUEST
            })?;
        serde_json::from_str(&decrypted).map_err(|_| StatusCode::BAD_REQUEST)?
    } else {
        serde_json::from_value(raw).map_err(|_| StatusCode::BAD_REQUEST)?
    };

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
