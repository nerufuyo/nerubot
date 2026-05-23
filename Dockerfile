# ── Stage 1: Build ──────────────────────────────────────
FROM rust:slim AS builder

RUN apt-get update && apt-get install -y pkg-config libssl-dev && rm -rf /var/lib/apt/lists/*
RUN rustup toolchain install nightly && rustup default nightly

WORKDIR /app
COPY . .

RUN cargo build --release

# ── Stage 2: Runtime ────────────────────────────────────
FROM debian:trixie-slim

RUN apt-get update && apt-get install -y ca-certificates libssl3 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/nerubot /usr/local/bin/
COPY --from=builder /app/migrations /app/migrations

WORKDIR /app
CMD ["nerubot"]
