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
        
        # Add help and info cogs
        try:
            # Load help system cogs from features directory
            from src.features.help.cogs.help_cog import HelpCog
            from src.features.help.cogs.about_cog import AboutCog
            from src.features.help.cogs.features_cog import FeaturesCog
            from src.features.help.cogs.commands_cog import CommandsCog
            
            # Register all help-related cogs
            await self.add_cog(HelpCog(self))
            await self.add_cog(AboutCog(self))
            await self.add_cog(FeaturesCog(self))
            await self.add_cog(CommandsCog(self))
            logger.info("Loaded Help, About, Features, and Commands Cogs")
        except Exception as e:
            logger.error(f"Failed to load Help Cogs: {e}")
        
        # Additional cogs can be added here as the bot grows
        
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
            logger.info(f"Logged in as {self.user.name} (ID: {self.user.id})")
        logger.info(f"Connected to {len(self.guilds)} guilds")
        
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
            await ctx.send(f"Missing required argument: {error.param.name}")
            return
            
        if isinstance(error, commands.BadArgument):
            await ctx.send(f"Bad argument: {error}")
            return
            
        # Log other errors
        logger.error(f"Error in command {ctx.command}: {error}")
        await ctx.send(f"An error occurred: {error}")
