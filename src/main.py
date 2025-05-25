"""
Main entry point for the Discord music bot
"""
import asyncio
import os
import logging
from dotenv import load_dotenv
from src.interfaces.discord.bot import NeruBot
from src.core.utils.logging_utils import setup_logger
from src.core.utils.messages import (
    BOT_STARTED,
    CONFIG_TOKEN_MISSING,
    BOT_SHUTDOWN
)

# Load environment variables
load_dotenv()

# Setup logging
logger = setup_logger(__name__, level=logging.INFO)

# Enable Discord.py debug logging
discord_logger = logging.getLogger('discord')
discord_logger.setLevel(logging.INFO)
handler = logging.StreamHandler()
handler.setFormatter(logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
discord_logger.addHandler(handler)


async def main():
    """Main entry point for the bot."""
    # Get the bot token
    token = os.getenv("DISCORD_TOKEN")
    if not token:
        logger.error(CONFIG_TOKEN_MISSING)
        return
    
    logger.info("Discord token found, creating bot...")
    
    # Create and start the bot
    bot = NeruBot(prefix="/")
    logger.info("Bot instance created, starting...")
    
    async with bot:
        logger.info("Starting bot connection...")
        await bot.start(token)


if __name__ == "__main__":
    try:
        logger.info(BOT_STARTED)
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info(BOT_SHUTDOWN)
    except Exception as e:
        logger.error(f"Fatal error: {e}", exc_info=True)
