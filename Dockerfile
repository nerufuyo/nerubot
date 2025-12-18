# NeruBot Go Edition - Production Dockerfile
FROM golang:1.25-alpine AS builder

# Install build dependencies and required tools
RUN apk add --no-cache git make ca-certificates tzdata python3 py3-pip ffmpeg

# Install yt-dlp
RUN pip3 install --break-system-packages yt-dlp

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the bot
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o nerubot ./cmd/nerubot

# Final stage - using alpine for smaller image with required tools
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata python3 py3-pip ffmpeg && \
    pip3 install --break-system-packages yt-dlp

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/nerubot .

# Copy data directories
COPY --from=builder /app/data ./data

# Create logs directory
RUN mkdir -p /app/logs

# Create non-root user
RUN addgroup -g 1000 nerubot && \
    adduser -D -u 1000 -G nerubot nerubot && \
    chown -R nerubot:nerubot /app

USER nerubot

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD pgrep -f nerubot || exit 1

CMD ["./nerubot"]
