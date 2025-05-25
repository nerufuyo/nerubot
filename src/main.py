"""
Main entry point for the Discord music bot
"""
import asyncio
import os
import logging
from dotenv import load_dotenv
from src.interfaces.discord.bot import NeruBot
from src.core.utils.messages import (
    BOT_STARTED,
    CONFIG_TOKEN_MISSING,
    BOT_SHUTDOWN
)

# Get logger without reconfiguring (don't use get_logger to avoid circular imports)
logger = logging.getLogger(__name__)

# Load environment variables
load_dotenv()


async def main():
    """Main entry point for the bot."""
    # Get the bot token
    token = os.getenv("DISCORD_TOKEN")
    if not token:
        logger.error(CONFIG_TOKEN_MISSING)
        return
    
    # Create and start the bot
    bot = NeruBot(prefix="!")
    async with bot:
        await bot.start(token)


if __name__ == "__main__":
    try:
        logger.info(BOT_STARTED)
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info(BOT_SHUTDOWN)
    except Exception as e:
        logger.error(f"Fatal error: {e}", exc_info=True)
