pub mod admin;
pub mod auth;

use axum::{
    routing::{get, post, put, delete},
    Router,
    response::Html,
};
use std::sync::Arc;
use tower_http::cors::CorsLayer;

use crate::handlers::BotData;

pub fn create_router(bot_data: Arc<BotData>) -> Router {
    Router::new()
        // ── Dashboard UI ────────────────────────────────
        .route("/dashboard", get(dashboard_html))
        .route("/dashboard/reminders", get(dashboard_html))

        // ── Dashboard Stats ─────────────────────────────
        .route("/api/admin/stats", get(admin::get_stats))
        .route("/api/admin/stats/guild/{guild_id}", get(admin::get_guild_stats))

        // ── Chat History ────────────────────────────────
        .route("/api/admin/chat/history", get(admin::get_chat_history))
        .route("/api/admin/chat/clear", delete(admin::clear_chat_history))

        // ── Roasts ──────────────────────────────────────
        .route("/api/admin/roasts", get(admin::get_roasts))
        .route("/api/admin/roasts", post(admin::create_roast))
        .route("/api/admin/roasts/{id}", delete(admin::delete_roast))

        // ── Moderation ──────────────────────────────────
        .route("/api/admin/mod/logs", get(admin::get_mod_logs))
        .route("/api/admin/mod/warnings", get(admin::get_warnings))
        .route("/api/admin/mod/warnings/{id}/resolve", put(admin::resolve_warning))

        // ── Reminders ───────────────────────────────────
        .route("/api/admin/reminders", get(admin::get_reminders))
        .route("/api/admin/reminders", post(admin::create_reminder))
        .route("/api/admin/reminders/types", get(admin::get_reminder_types))
        .route("/api/admin/reminders/{id}", put(admin::update_reminder))
        .route("/api/admin/reminders/{id}", delete(admin::delete_reminder))
        .route("/api/admin/reminders/{id}/toggle", put(admin::toggle_reminder))

        // ── Polls ───────────────────────────────────────
        .route("/api/admin/polls", get(admin::get_polls))
        .route("/api/admin/polls", post(admin::create_poll))
        .route("/api/admin/polls/{id}/close", put(admin::close_poll))

        // ── Guild Config ────────────────────────────────
        .route("/api/admin/guilds", get(admin::get_guilds))
        .route("/api/admin/guilds/{guild_id}", get(admin::get_guild_config))
        .route("/api/admin/guilds/{guild_id}", put(admin::update_guild_config))

        // ── Blocked Users ───────────────────────────────
        .route("/api/admin/blocked", get(admin::get_blocked_users))
        .route("/api/admin/blocked", post(admin::block_user))
        .route("/api/admin/blocked/{guild_id}/{user_id}", delete(admin::unblock_user))

        // ── Bot Settings ────────────────────────────────
        .route("/api/admin/settings", get(admin::get_settings))
        .route("/api/admin/settings", put(admin::update_settings))

        // ── Command Usage ───────────────────────────────
        .route("/api/admin/commands/usage", get(admin::get_command_usage))
        .route("/api/admin/commands/top", get(admin::get_top_commands))

        // ── Health ──────────────────────────────────────
        .route("/api/health", get(|| async { "{\"status\":\"ok\"}" }))

        .layer(CorsLayer::permissive())
        .with_state(bot_data)
}

async fn dashboard_html() -> Html<&'static str> {
    Html(include_str!("../../static/dashboard.html"))
}
