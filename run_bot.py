#!/usr/bin/env python3
"""
Helper script to run the bot with appropriate settings
"""
import os
import sys
import subprocess
import argparse
import signal
import platform
import logging
from src.core.utils.logging_utils import configure_root_logger, get_logger
from src.core.utils.file_utils import read_env_file, write_env_file, find_ffmpeg
from src.core.utils.messages import (
    BOT_STARTED, CONFIG_TOKEN_MISSING, CONFIG_TOKEN_CREATE, 
    CONFIG_TOKEN_PROMPT, CONFIG_TOKEN_EMPTY, CONFIG_TOKEN_SAVED,
    CONFIG_DEPS_MISSING, CONFIG_DEPS_INSTALL_PROMPT, CONFIG_DEPS_SKIPPED,
    CONFIG_FFMPEG_NOT_FOUND, CONFIG_FFMPEG_INSTALL, CONFIG_FFMPEG_WIN,
    CONFIG_FFMPEG_MAC, CONFIG_FFMPEG_LINUX, CONFIG_FFMPEG_PATH,
    CONFIG_FFMPEG_VERSION, CONFIG_ISSUES, CONFIG_START,
    VPS_SETUP_TITLE, VPS_USERNAME_PROMPT, VPS_GROUP_PROMPT, VPS_PATH_PROMPT,
    VPS_SERVICE_CREATED, VPS_DEPLOY_HEADER, VPS_DEPLOY_COPY, VPS_DEPLOY_DEPS,
    VPS_DEPLOY_FFMPEG, VPS_DEPLOY_SERVICE, VPS_DEPLOY_ENABLE,
    VPS_MONITOR_HEADER, VPS_MONITOR_STATUS, VPS_MONITOR_LOGS
)

# Get the absolute path to the project directory
PROJECT_DIR = os.path.dirname(os.path.abspath(__file__))

# Configure logger
logger = get_logger(__name__)

# Global variables
FFMPEG_PATH = None  # Will be set by check_ffmpeg()


def validate_token():
    """Validate that the Discord token exists in .env file."""
    env_path = os.path.join(PROJECT_DIR, '.env')
    
    if not os.path.exists(env_path):
        logger.warning(CONFIG_TOKEN_MISSING)
        with open(env_path, 'w') as f:
            f.write("# Discord Bot Token\nDISCORD_TOKEN=\n")
        logger.info(CONFIG_TOKEN_CREATE.format(path=env_path))
        return False
    
    env_vars = read_env_file(env_path)
    
    if "DISCORD_TOKEN" in env_vars and env_vars["DISCORD_TOKEN"].strip():
        return True
    
    token = input(CONFIG_TOKEN_PROMPT)
    if not token.strip():
        logger.error(CONFIG_TOKEN_EMPTY)
        return False
    
    # Update the .env file with the token
    env_vars["DISCORD_TOKEN"] = token
    write_env_file(env_vars, env_path)
    
    logger.info(CONFIG_TOKEN_SAVED)
    return True


def check_dependencies():
    """Check if required dependencies are installed."""
    requirements = [
        'discord.py',
        'python-dotenv',
        'yt-dlp',
        'PyNaCl'
    ]
    
    try:
        # Check for installed packages using pip
        import subprocess
        
        # Get list of installed packages
        result = subprocess.run(
            [sys.executable, '-m', 'pip', 'list'], 
            capture_output=True, 
            text=True
        )
        
        installed_packages = result.stdout.lower()
        
        # Use a more accurate approach to check for installed packages
        missing = []
        
        # Check for packages more accurately in pip list output
        if 'discord ' not in installed_packages and 'discord.py' not in installed_packages:
            missing.append('discord.py')
            
        if 'python-dotenv' not in installed_packages and 'dotenv' not in installed_packages:
            missing.append('python-dotenv')
            
        if 'yt-dlp' not in installed_packages and 'yt_dlp' not in installed_packages:
            missing.append('yt-dlp')
            
        if 'pynacl' not in installed_packages:
            missing.append('PyNaCl')
        
        if missing:
            logger.warning(CONFIG_DEPS_MISSING.format(deps=', '.join(missing)))
            install = input(CONFIG_DEPS_INSTALL_PROMPT)
            
            if install.lower() == 'y':
                subprocess.check_call([
                    sys.executable, '-m', 'pip', 'install', *missing
                ])
                return True
            else:
                logger.warning(CONFIG_DEPS_SKIPPED)
                return False
        
        return True
    except Exception as e:
        logger.error(f"Error checking dependencies: {e}")
        return False


def check_ffmpeg():
    """Check if ffmpeg is installed."""
    global FFMPEG_PATH
    
    # Use the centralized file utility
    FFMPEG_PATH = find_ffmpeg()
    
    if FFMPEG_PATH:
        logger.info(CONFIG_FFMPEG_PATH.format(path=FFMPEG_PATH))
        
        # Try to get ffmpeg version
        try:
            result = subprocess.run(
                [FFMPEG_PATH, '-version'],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            version_line = result.stdout.strip().split('\n')[0]
            logger.info(CONFIG_FFMPEG_VERSION.format(version=version_line))
            return True
        except Exception:
            pass
    
    # FFmpeg not found or not working properly
    logger.error(CONFIG_FFMPEG_NOT_FOUND)
    logger.error(CONFIG_FFMPEG_INSTALL)
    
    # Show platform-specific instructions
    platform_name = platform.system().lower()
    if 'windows' in platform_name:
        logger.error(CONFIG_FFMPEG_WIN)
    elif 'darwin' in platform_name:
        logger.error(CONFIG_FFMPEG_MAC)
    else:  # Linux or other
        logger.error(CONFIG_FFMPEG_LINUX)
    
    return False


def run_bot():
    """Run the bot."""
    logger.info(CONFIG_START)
    sys.path.insert(0, PROJECT_DIR)
    try:
        import asyncio
        from src.main import main
        asyncio.run(main())
    except Exception as e:
        logger.error(f"Error running bot: {e}")
        sys.exit(1)


def setup_vps():
    """Generate a systemd service file for VPS deployment."""
    logger.info(VPS_SETUP_TITLE)
    
    service_content = """[Unit]
Description=Nerubot Discord Music Bot
After=network.target

[Service]
User=USER_NAME
Group=USER_GROUP
WorkingDirectory=BOT_DIR
ExecStart=/usr/bin/python3 -m src.main
Restart=on-failure
RestartSec=5
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=nerubot

[Install]
WantedBy=multi-user.target
"""
    
    username = input(VPS_USERNAME_PROMPT)
    group = input(VPS_GROUP_PROMPT)
    if not group:
        group = username
    
    bot_dir = input(VPS_PATH_PROMPT.format(username=username))
    
    service_content = service_content.replace("USER_NAME", username)
    service_content = service_content.replace("USER_GROUP", group)
    service_content = service_content.replace("BOT_DIR", bot_dir)
    
    service_path = os.path.join(PROJECT_DIR, "nerubot.service")
    with open(service_path, "w") as f:
        f.write(service_content)
    
    logger.info(VPS_SERVICE_CREATED.format(path=service_path))
    logger.info(VPS_DEPLOY_HEADER)
    logger.info(VPS_DEPLOY_COPY.format(project_dir=PROJECT_DIR, username=username, bot_dir=bot_dir))
    logger.info(VPS_DEPLOY_DEPS.format(bot_dir=bot_dir))
    logger.info(VPS_DEPLOY_FFMPEG)
    logger.info(VPS_DEPLOY_SERVICE.format(bot_dir=bot_dir))
    logger.info(VPS_DEPLOY_ENABLE)
    logger.info(VPS_MONITOR_HEADER)
    logger.info(VPS_MONITOR_STATUS)
    logger.info(VPS_MONITOR_LOGS)


def signal_handler(sig, frame):
    """Handle Ctrl+C to exit gracefully."""
    logger.info("Shutting down the bot...")
    sys.exit(0)


if __name__ == "__main__":
    # Configure root logger at the entry point (only once)
    # Check if root logger already has handlers before configuring
    root_logger = logging.getLogger()
    if not root_logger.handlers:
        configure_root_logger()
    
    # Set up signal handler for graceful shutdown
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    
    parser = argparse.ArgumentParser(description="NeruBot Discord Music Bot")
    parser.add_argument(
        "--setup-vps", 
        action="store_true", 
        help="Set up the bot for VPS deployment"
    )
    args = parser.parse_args()
    
    if args.setup_vps:
        setup_vps()
        sys.exit(0)
    
    logger.info(BOT_STARTED)
    
    # Perform checks
    token_valid = validate_token()
    deps_valid = check_dependencies()
    ffmpeg_valid = check_ffmpeg()
    
    if not all([token_valid, deps_valid, ffmpeg_valid]):
        logger.error(CONFIG_ISSUES)
        sys.exit(1)
    
    run_bot()
