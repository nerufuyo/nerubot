pub mod bot;
pub mod handler;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LarkConfig {
    pub app_id: String,
    pub app_secret: String,
    pub verification_token: String,
    pub encrypt_key: String,
}

impl LarkConfig {
    pub fn from_env() -> Self {
        Self {
            app_id: std::env::var("LARK_APP_ID").unwrap_or_default(),
            app_secret: std::env::var("LARK_APP_SECRET").unwrap_or_default(),
            verification_token: std::env::var("LARK_VERIFICATION_TOKEN").unwrap_or_default(),
            encrypt_key: std::env::var("LARK_ENCRYPT_KEY").unwrap_or_default(),
        }
    }
}
