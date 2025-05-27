#!/bin/bash
# Run NeruBot with initial help message
cd /Users/infantai/nerubot
export PYTHONPATH=/Users/infantai/nerubot

# Display colorful help message
echo -e "\033[1;36m🤖 NeruBot - Discord Music & News Bot\033[0m"
echo -e "\033[1;36m=====================================\033[0m"
echo -e "\033[1;32m✅ Starting NeruBot...\033[0m"
echo -e "\033[1;33m💡 Available commands:\033[0m"
echo -e "\033[1;37m  • \033[1;34m/help\033[1;37m - Show detailed help menu with categories"
echo -e "\033[1;37m  • \033[1;34m/commands\033[1;37m - Show compact command reference"
echo -e "\033[1;37m  • \033[1;34m/play <song>\033[1;37m - Play music from YouTube, Spotify, or SoundCloud"
echo -e "\033[1;37m  • \033[1;34m/news latest\033[1;37m - Get the latest news from trusted sources"
echo -e "\033[1;37m  • \033[1;34m/news set-channel\033[1;37m - Enable automatic news updates"
echo -e "\033[1;37m  • \033[1;34m/about\033[1;37m - Show bot information"
echo -e "\033[1;37m  • \033[1;34m/features\033[1;37m - Show available features"
echo ""
echo -e "\033[1;35m🎵 Music sources:\033[0m YouTube, Spotify, SoundCloud"
echo -e "\033[1;35m📰 News sources:\033[0m BBC, Reuters, AP News, CNN, NPR, Al Jazeera"
echo -e "\033[1;35m📘 Help System:\033[0m Interactive with pagination"
echo ""
echo -e "\033[1;32mStarting bot now...\033[0m"
echo -e "\033[1;36m=====================================\033[0m"

# Run the bot
/Users/infantai/nerubot/nerubot_env/bin/python src/main.py
