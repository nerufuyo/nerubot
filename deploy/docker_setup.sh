#!/bin/bash
# Docker Deployment for NeruBot
# Alternative deployment method using Docker containers

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[NeruBot] $1${NC}"
}

print_success() {
    echo -e "${GREEN}[NeruBot] $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[NeruBot] $1${NC}"
}

print_error() {
    echo -e "${RED}[NeruBot] $1${NC}"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        print_info "Ubuntu/Debian: curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

# Create Dockerfile
create_dockerfile() {
    print_info "Creating Dockerfile..."
    cat > Dockerfile << 'EOF'
# NeruBot Docker Image
FROM python:3.11-slim-bullseye

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    libopus0 \
    libopus-dev \
    libffi-dev \
    libnacl-dev \
    git \
    && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Create non-root user
RUN groupadd -r nerubot && useradd -r -g nerubot nerubot

# Copy requirements first for better caching
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY src/ ./src/
COPY .env* ./

# Create logs directory
RUN mkdir -p /app/logs && chown nerubot:nerubot /app/logs

# Switch to non-root user
USER nerubot

# Health check
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD python -c "import asyncio; import aiohttp; print('Bot is healthy')" || exit 1

# Run the bot
CMD ["python", "src/main.py"]
EOF
    print_success "Dockerfile created"
}

# Create docker-compose.yml
create_docker_compose() {
    print_info "Creating docker-compose.yml..."
    cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  nerubot:
    build: .
    container_name: nerubot
    restart: unless-stopped
    env_file:
      - .env
    volumes:
      - ./logs:/app/logs
      - ./data:/app/data
    environment:
      - PYTHONUNBUFFERED=1
      - LOG_LEVEL=${LOG_LEVEL:-INFO}
    networks:
      - nerubot_network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  # Optional: Add a monitoring container
  watchtower:
    image: containrrr/watchtower
    container_name: nerubot_watchtower
    restart: unless-stopped
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - WATCHTOWER_CLEANUP=true
      - WATCHTOWER_POLL_INTERVAL=86400  # Check daily
      - WATCHTOWER_INCLUDE_STOPPED=true
    command: nerubot

networks:
  nerubot_network:
    driver: bridge
EOF
    print_success "docker-compose.yml created"
}

# Create .dockerignore
create_dockerignore() {
    print_info "Creating .dockerignore..."
    cat > .dockerignore << 'EOF'
# Git
.git
.gitignore

# Python
__pycache__
*.pyc
*.pyo
*.pyd
.Python
env
pip-log.txt
pip-delete-this-directory.txt
.tox
.coverage
.coverage.*
.cache
nosetests.xml
coverage.xml
*.cover
*.log
.git
.mypy_cache
.pytest_cache
.hypothesis

# Virtual environments
nerubot_env/
venv/
ENV/

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Deployment
deploy/
*.sh
!docker-entrypoint.sh

# Logs
*.log
logs/

# Temporary files
*.tmp
*.temp
EOF
    print_success ".dockerignore created"
}

# Create environment template for Docker
create_env_template() {
    if [ ! -f .env ]; then
        print_info "Creating .env template for Docker..."
        cat > .env.docker << 'EOF'
# NeruBot Environment Configuration for Docker

# Required: Discord Bot Token
DISCORD_TOKEN=your_discord_bot_token_here

# Optional: Logging Level
LOG_LEVEL=INFO

# Optional: Bot Prefix
BOT_PREFIX=!

# Optional: Spotify Integration
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# Optional: Feature Flags
ENABLE_NEWS=true
ENABLE_24_7=true
EOF
        print_warning "Created .env.docker template. Copy it to .env and configure your tokens."
    fi
}

# Create Docker deployment scripts
create_docker_scripts() {
    print_info "Creating Docker management scripts..."
    
    # Build script
    cat > docker-build.sh << 'EOF'
#!/bin/bash
# Build NeruBot Docker image

echo "Building NeruBot Docker image..."
docker-compose build --no-cache nerubot
echo "Build complete!"
EOF
    chmod +x docker-build.sh

    # Start script
    cat > docker-start.sh << 'EOF'
#!/bin/bash
# Start NeruBot with Docker Compose

echo "Starting NeruBot..."
docker-compose up -d nerubot
echo "NeruBot started!"
echo "Check logs with: docker-compose logs -f nerubot"
EOF
    chmod +x docker-start.sh

    # Stop script
    cat > docker-stop.sh << 'EOF'
#!/bin/bash
# Stop NeruBot

echo "Stopping NeruBot..."
docker-compose down
echo "NeruBot stopped!"
EOF
    chmod +x docker-stop.sh

    # Logs script
    cat > docker-logs.sh << 'EOF'
#!/bin/bash
# View NeruBot logs

docker-compose logs -f nerubot
EOF
    chmod +x docker-logs.sh

    # Update script
    cat > docker-update.sh << 'EOF'
#!/bin/bash
# Update and restart NeruBot

echo "Updating NeruBot..."
git pull origin main
docker-compose build --no-cache nerubot
docker-compose down
docker-compose up -d nerubot
echo "Update complete!"
EOF
    chmod +x docker-update.sh

    print_success "Docker management scripts created"
}

# Main function
main() {
    print_success "=== NeruBot Docker Setup ==="
    echo

    check_docker
    create_dockerfile
    create_docker_compose
    create_dockerignore
    create_env_template
    create_docker_scripts

    print_success "=== Docker Setup Complete! ==="
    echo
    print_info "Next steps:"
    print_info "1. Configure your .env file with Discord token"
    print_info "2. Build the image: ./docker-build.sh"
    print_info "3. Start the bot: ./docker-start.sh"
    echo
    print_info "Management commands:"
    print_info "- Start: ./docker-start.sh"
    print_info "- Stop: ./docker-stop.sh"
    print_info "- View logs: ./docker-logs.sh"
    print_info "- Update: ./docker-update.sh"
    echo
    print_info "Or use docker-compose directly:"
    print_info "- docker-compose up -d nerubot"
    print_info "- docker-compose logs -f nerubot"
    print_info "- docker-compose down"
}

main "$@"
