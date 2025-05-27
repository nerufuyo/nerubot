"""
Main Discord bot class
"""
import os
import discord
from discord.ext import commands
from typing import List, Optional
from dotenv import load_dotenv
from src.features.music.cogs.music_cog import MusicCog
from src.core.utils.logging_utils import get_logger
from src.core.constants import (
    BOT_NAME, BOT_VERSION, BOT_DEFAULT_STATUS, BOT_DEFAULT_ACTIVITY_TYPE,
    MSG_INFO, LOG_MSG, COLOR_SUCCESS
)

# Get logger
logger = get_logger(__name__)

# Load environment variables
load_dotenv()

# Bot config
DEFAULT_PREFIX = "!"
DISCORD_TOKEN = os.getenv("DISCORD_TOKEN")
INTENTS = discord.Intents.default()
INTENTS.message_content = True  # Enable message content intent


class NeruBot(commands.Bot):
    """Main Discord bot class."""
    
    def __init__(self, prefix: str = DEFAULT_PREFIX):
        super().__init__(
            command_prefix=prefix,
            intents=INTENTS,
            help_command=None,  # Disable default help command
            description=f"{BOT_NAME} v{BOT_VERSION} - A Discord music bot with clean architecture"
        )
    
    async def setup_hook(self) -> None:
        """Initialize the bot when it's ready."""
        logger.info(LOG_MSG["bot_starting"].format(bot_name=BOT_NAME, version=BOT_VERSION))
        
        # Add music cog
        await self.add_cog(MusicCog(self))
        logger.info(LOG_MSG["cog_loaded"].format(cog_name="MusicCog"))
        
        # Add news cog
        try:
            from src.features.news.cogs.news_cog import NewsCog
            await self.add_cog(NewsCog(self))
            logger.info(LOG_MSG["cog_loaded"].format(cog_name="NewsCog"))
        except Exception as e:
            logger.error(LOG_MSG["cog_failed"].format(cog_name="NewsCog", error=e))
        
        # Add help system cogs
        try:
            from src.features.help.cogs.help_cog import HelpCog
            from src.features.help.cogs.about_cog import AboutCog
            from src.features.help.cogs.features_cog import FeaturesCog
            from src.features.help.cogs.commands_cog import CommandsCog
            
            await self.add_cog(HelpCog(self))
            await self.add_cog(AboutCog(self))
            await self.add_cog(FeaturesCog(self))
            await self.add_cog(CommandsCog(self))
            logger.info(LOG_MSG["cog_loaded"].format(cog_name="Help system"))
        except Exception as e:
            logger.error(LOG_MSG["cog_failed"].format(cog_name="Help system", error=e))
        
        # Add confession cog
        try:
            from src.features.confession.cogs.confession_cog import ConfessionCog
            await self.add_cog(ConfessionCog(self))
            logger.info(LOG_MSG["cog_loaded"].format(cog_name="ConfessionCog"))
        except Exception as e:
            logger.error(LOG_MSG["cog_failed"].format(cog_name="ConfessionCog", error=e))
        
        # Force sync all commands with Discord (globally)
        try:
            synced = await self.tree.sync()
            logger.info(f"Synced {len(synced)} command(s) with Discord")
        except Exception as e:
            logger.error(f"Failed to sync commands: {e}")
        
        logger.info("Bot setup complete")
    
    async def on_ready(self) -> None:
        """Called when the bot is ready."""
        if self.user:
            logger.info(LOG_MSG["bot_ready"].format(user=f"{self.user.name} (ID: {self.user.id})"))
        logger.info(f"Connected to {len(self.guilds)} guilds")
        
        # Set bot status now that we're connected
        try:
            await self.change_presence(
                activity=discord.Activity(
                    type=getattr(discord.ActivityType, BOT_DEFAULT_ACTIVITY_TYPE, discord.ActivityType.listening),
                    name=BOT_DEFAULT_STATUS
                ),
                status=discord.Status.online
            )
            logger.info(f"Bot status set to: {BOT_DEFAULT_STATUS}")
        except Exception as e:
            logger.error(f"Failed to set bot status: {e}")
        
        # Send ready message to console
        print(MSG_INFO["bot_ready"].format(bot_name=BOT_NAME))
    
    async def on_disconnect(self) -> None:
        """Called when the bot disconnects."""
        logger.info(LOG_MSG["bot_disconnected"])
    
    async def on_command_error(self, ctx: commands.Context, error) -> None:
        """Handle command errors."""
        if isinstance(error, commands.CommandNotFound):
            # Ignore command not found errors
            return
            
        if isinstance(error, commands.MissingRequiredArgument):
            await ctx.send(f"Missing required argument: {error.param.name}")
            return
            
        if isinstance(error, commands.BadArgument):
            await ctx.send(f"Bad argument: {error}")
            return
            
        # Log other errors
        logger.error(f"Error in command {ctx.command}: {error}")
        await ctx.send(f"An error occurred: {error}")


def create_bot() -> NeruBot:
    """Create and return a bot instance."""
    return NeruBot()


def run_bot() -> None:
    """Run the Discord bot."""
    if not DISCORD_TOKEN:
        logger.error("DISCORD_TOKEN not found in environment variables")
        return
    
    bot = create_bot()
    try:
        bot.run(DISCORD_TOKEN)
    except discord.LoginFailure:
        logger.error("Invalid Discord token")
    except Exception as e:
        logger.error(f"Failed to run bot: {e}")
