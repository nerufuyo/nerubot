"""
Main Discord bot class
"""
import os
import discord
from discord.ext import commands
from typing import List, Optional
from dotenv import load_dotenv
from src.interfaces.discord.music_cog import MusicCog
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
        
        # Force sync all commands with Discord (globally)
        try:
            synced = await self.tree.sync()
            logger.info(f"Synced {len(synced)} command(s) with Discord")
        except Exception as e:
            logger.error(f"Failed to sync commands: {e}")
        
        logger.info(BOT_SETUP_COMPLETE)
    
    async def on_ready(self) -> None:
        """Called when the bot is ready."""
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
    
    @commands.hybrid_command(name="help", description="Show help information")
    async def help_command(self, ctx: commands.Context) -> None:
        """Show help information."""
        embed = discord.Embed(
            title=HELP_TITLE,
            description=HELP_DESCRIPTION,
            color=discord.Color.blue()
        )
        
        music_commands = [
            (f"{self.command_prefix}join", HELP_JOIN_DESC),
            (f"{self.command_prefix}leave", HELP_LEAVE_DESC),
            (f"{self.command_prefix}play <song>", HELP_PLAY_DESC),
            (f"{self.command_prefix}stop", HELP_STOP_DESC),
            (f"{self.command_prefix}pause", HELP_PAUSE_DESC),
            (f"{self.command_prefix}resume", HELP_RESUME_DESC),
            (f"{self.command_prefix}skip", HELP_SKIP_DESC),
            (f"{self.command_prefix}volume <0-100>", HELP_VOLUME_DESC),
            (f"{self.command_prefix}now", HELP_NOW_DESC),
            (f"{self.command_prefix}queue [page]", HELP_QUEUE_DESC),
            (f"{self.command_prefix}remove <index>", HELP_REMOVE_DESC),
            (f"{self.command_prefix}shuffle", HELP_SHUFFLE_DESC),
            (f"{self.command_prefix}loop <off/song/queue>", HELP_LOOP_DESC)
        ]
        
        commands_text = "\n".join([f"**{cmd}**: {desc}" for cmd, desc in music_commands])
        embed.add_field(name=HELP_MUSIC_COMMANDS_TITLE, value=commands_text, inline=False)
        
        await ctx.send(embed=embed)
