# Dockerfile for Music Service
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o music ./services/music/cmd

# Final stage
FROM alpine:latest

# Install FFmpeg and yt-dlp
RUN apk --no-cache add ca-certificates tzdata ffmpeg python3 py3-pip && \
    pip3 install --no-cache-dir yt-dlp

WORKDIR /root/

COPY --from=builder /app/music .
COPY --from=builder /app/internal/config ./internal/config

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

CMD ["./music"]
