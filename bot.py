#!/usr/bin/env python3
"""
NeruBot - Simple and extensible Discord bot
Main entry point with clean architecture
"""
import asyncio
import logging
import os
from pathlib import Path
import discord
from discord.ext import commands
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('bot.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

class NeruBot(commands.Bot):
    """Main bot class with automatic cog loading."""
    
    def __init__(self):
        # Set up intents
        intents = discord.Intents.default()
        intents.message_content = True
        intents.voice_states = True
        
        super().__init__(
            command_prefix='!',
            intents=intents,
            help_command=None,
            description="NeruBot - A modular Discord bot"
        )
        
    async def setup_hook(self):
        """Load all cogs automatically."""
        logger.info("🤖 Loading cogs...")
        
        # Get the cogs directory
        cogs_dir = Path(__file__).parent / 'cogs'
        
        # Load all Python files in cogs directory
        if cogs_dir.exists():
            for cog_file in cogs_dir.glob('*.py'):
                if cog_file.name == '__init__.py':
                    continue
                    
                cog_name = f'cogs.{cog_file.stem}'
                try:
                    await self.load_extension(cog_name)
                    logger.info(f"✅ Loaded cog: {cog_name}")
                except Exception as e:
                    logger.error(f"❌ Failed to load cog {cog_name}: {e}")
        
        # Sync slash commands
        try:
            synced = await self.tree.sync()
            logger.info(f"🔄 Synced {len(synced)} slash commands")
        except Exception as e:
            logger.error(f"❌ Failed to sync commands: {e}")
    
    async def on_ready(self):
        """Called when bot is ready."""
        logger.info(f"🚀 {self.user} is online!")
        logger.info(f"📡 Connected to {len(self.guilds)} servers")
        
        # Set bot status
        await self.change_presence(
            activity=discord.Activity(
                type=discord.ActivityType.listening,
                name="!help | NeruBot"
            )
        )
    
    async def on_command_error(self, ctx, error):
        """Handle command errors."""
        if isinstance(error, commands.CommandNotFound):
            return  # Ignore unknown commands
        
        if isinstance(error, commands.MissingRequiredArgument):
            await ctx.send(f"❌ Missing argument: `{error.param.name}`")
        elif isinstance(error, commands.BadArgument):
            await ctx.send(f"❌ Invalid argument: {error}")
        else:
            logger.error(f"Command error in {ctx.command}: {error}")
            await ctx.send("❌ An error occurred while processing the command.")

async def main():
    """Main function to run the bot."""
    # Check for Discord token
    token = os.getenv('DISCORD_TOKEN')
    if not token:
        logger.error("❌ DISCORD_TOKEN not found in environment variables!")
        logger.error("Please create a .env file with your Discord bot token.")
        return
    
    # Create and start bot
    bot = NeruBot()
    
    try:
        logger.info("🏁 Starting NeruBot...")
        await bot.start(token)
    except KeyboardInterrupt:
        logger.info("👋 Shutting down NeruBot...")
    except Exception as e:
        logger.error(f"💥 Fatal error: {e}")
    finally:
        await bot.close()

if __name__ == '__main__':
    asyncio.run(main())
