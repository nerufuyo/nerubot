#!/bin/bash
# Run NeruBot in debug mode with help system info
cd /Users/infantai/nerubot
export PYTHONPATH=/Users/infantai/nerubot
export LOG_LEVEL=DEBUG

# Display colorful debug info
echo -e "\033[1;36m🤖 NeruBot - DEBUG MODE\033[0m"
echo -e "\033[1;36m===============================\033[0m"
echo -e "\033[1;32m✅ Starting NeruBot in debug mode...\033[0m"
echo -e "\033[1;33m📊 System Information:\033[0m"

# Check Python version
PY_VERSION=$(python3 --version)
echo -e "\033[1;37m  • Python: \033[1;34m$PY_VERSION\033[0m"

# Check installed packages
echo -e "\033[1;37m  • Checking dependencies...\033[0m"
pip list | grep -E "discord|yt-dlp|spotipy|beautifulsoup4" | sed 's/^/    /'

echo -e "\033[1;33m🔍 Debug Help Commands:\033[0m"
echo -e "\033[1;37m  • \033[1;34m/help\033[1;37m - Test paginated help system"
echo -e "\033[1;37m  • \033[1;34m/commands\033[1;37m - Test command reference card"
echo -e "\033[1;37m  • \033[1;34m/about\033[1;37m - Test about information"
echo -e "\033[1;37m  • \033[1;34m/features\033[1;37m - Test feature showcase"
echo ""
echo -e "\033[1;32mStarting bot in debug mode...\033[0m"
echo -e "\033[1;36m===============================\033[0m"

# Run the bot with debug flags
/Users/infantai/nerubot/nerubot_env/bin/python src/main.py
