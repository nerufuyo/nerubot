#!/bin/bash
# Simple startup script for NeruBot

echo "🤖 Starting NeruBot..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    echo "Please create a .env file with your Discord token:"
    echo "DISCORD_TOKEN=your_token_here"
    exit 1
fi

# Check if Discord token is set
if ! grep -q "DISCORD_TOKEN=" .env || grep -q "DISCORD_TOKEN=$" .env; then
    echo "❌ Discord token not set in .env file!"
    echo "Please add your Discord token to the .env file:"
    echo "DISCORD_TOKEN=your_token_here"
    exit 1
fi

# Run the bot
python3 bot.py
