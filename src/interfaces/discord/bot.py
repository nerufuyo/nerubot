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
from src.core.utils.messages import (
    BOT_SETUP_COMPLETE,
    BOT_LOGGED_IN,
    BOT_GUILD_COUNT,
    ERROR_COMMAND,
    ERROR_MISSING_ARG,
    ERROR_BAD_ARG,
    HELP_TITLE,
    HELP_DESCRIPTION,
    HELP_MUSIC_COMMANDS_TITLE,
    HELP_JOIN_DESC,
    HELP_LEAVE_DESC,
    HELP_PLAY_DESC,
    HELP_STOP_DESC,
    HELP_PAUSE_DESC,
    HELP_RESUME_DESC,
    HELP_SKIP_DESC,
    HELP_VOLUME_DESC,
    HELP_NOW_DESC,
    HELP_QUEUE_DESC,
    HELP_REMOVE_DESC,
    HELP_SHUFFLE_DESC,
    HELP_LOOP_DESC
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
            description="A Discord music bot with clean architecture"
        )
        self.initial_extensions = [
            # Add more extensions as needed
        ]
    
    async def setup_hook(self) -> None:
        """Initialize the bot when it's ready."""
        # Add music cog
        await self.add_cog(MusicCog(self))
        
        # Add help cog
        try:
            from src.interfaces.discord.help_cog import HelpCog
            await self.add_cog(HelpCog(self))
            logger.info("Loaded HelpCog")
        except Exception as e:
            logger.error(f"Failed to load HelpCog: {e}")
        
        # Add feature cogs - temporarily disabled for testing
        # try:
        #     from src.features.news.cogs.news_cog import NewsCog
        #     await self.add_cog(NewsCog(self))
        #     logger.info("Loaded NewsCog")
        # except Exception as e:
        #     logger.error(f"Failed to load NewsCog: {e}")
        
        # Disabled - Coming Soon
        # try:
        #     from src.features.quotes.cogs.quotes_cog import QuotesCog
        #     await self.add_cog(QuotesCog(self))
        #     logger.info("Loaded QuotesCog")
        # except Exception as e:
        #     logger.error(f"Failed to load QuotesCog: {e}")
            
        # try:
        #     from src.features.profile.cogs.profile_cog import ProfileCog
        #     await self.add_cog(ProfileCog(self))
        #     logger.info("Loaded ProfileCog")
        # except Exception as e:
        #     logger.error(f"Failed to load ProfileCog: {e}")
            
        # try:
        #     from src.features.confession.cogs.confession_cog import ConfessionCog
        #     await self.add_cog(ConfessionCog(self))
        #     logger.info("Loaded ConfessionCog")
        # except Exception as e:
        #     logger.error(f"Failed to load ConfessionCog: {e}")
        
        # Force sync all commands with Discord (globally)
        try:
            synced = await self.tree.sync()
            logger.info(f"Synced {len(synced)} command(s) with Discord")
        except Exception as e:
            logger.error(f"Failed to sync commands: {e}")
        
        logger.info(BOT_SETUP_COMPLETE)
    
    async def on_ready(self) -> None:
        """Called when the bot is ready."""
        if self.user:
            logger.info(BOT_LOGGED_IN.format(name=self.user.name, id=self.user.id))
        logger.info(BOT_GUILD_COUNT.format(count=len(self.guilds)))
        
        # Set bot activity
        await self.change_presence(
            activity=discord.Activity(
                type=discord.ActivityType.listening,
                name=f"{self.command_prefix}help | Music"
            )
        )
    
    async def on_command_error(self, ctx: commands.Context, error) -> None:
        """Handle command errors."""
        if isinstance(error, commands.CommandNotFound):
            # Ignore command not found errors
            return
            
        if isinstance(error, commands.MissingRequiredArgument):
            await ctx.send(ERROR_MISSING_ARG.format(param=error.param.name))
            return
            
        if isinstance(error, commands.BadArgument):
            await ctx.send(ERROR_BAD_ARG.format(error=error))
            return
            
        # Log other errors
        logger.error(f"Error in command {ctx.command}: {error}")
        await ctx.send(ERROR_COMMAND.format(error=error))
