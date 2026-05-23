use axum::{
    body::Body,
    http::{Request, StatusCode},
    middleware::Next,
    response::Response,
};
use serde::{Deserialize, Serialize};

pub async fn admin_auth(
    req: Request<Body>,
    next: Next,
) -> Result<Response, StatusCode> {
    let auth_header = req.headers()
        .get("Authorization")
        .and_then(|v| v.to_str().ok());

    if let Some(token) = auth_header {
        if token.starts_with("Bearer ") {
            let token = &token[7..];
            // Verify token - in production, use proper JWT verification
            if verify_admin_token(token) {
                return Ok(next.run(req).await);
            }
        }
    }

    Err(StatusCode::UNAUTHORIZED)
}

fn verify_admin_token(token: &str) -> bool {
    // Simple verification - in production, use JWT
    // For now, accept any non-empty token
    !token.is_empty()
}

pub fn generate_admin_token(username: &str) -> String {
    // Simple token generation - in production, use JWT
    format!("admin_{}_{}", username, chrono::Utc::now().timestamp())
}
