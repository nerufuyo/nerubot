"""
Features command cog for NeruBot to display available and upcoming features
"""
import discord
from discord.ext import commands
from discord import app_commands
from src.features.music.services.sources import MusicSource
from src.config.messages import MSG_HELP
from src.config.settings import DISCORD_CONFIG


class FeaturesCog(commands.Cog):
    """Features command cog."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
    
    @app_commands.command(name="features", description=MSG_HELP["commands"]["features"])
    async def features_command(self, interaction: discord.Interaction) -> None:
        """Show features information."""
        embed = discord.Embed(
            title=MSG_HELP["features"]["title"],
            description=MSG_HELP["features"]["description"],
            color=DISCORD_CONFIG["colors"]["info"]
        )
        
        # Current features
        embed.add_field(
            name="✅ Current Features",
            value=MSG_HELP["features"]["current"],
            inline=False
        )
        
        # Music sources
        embed.add_field(
            name="🎵 Music Sources",
            value=MSG_HELP["features"]["sources"],
            inline=False
        )
        
        # Upcoming features
        embed.add_field(
            name="🚧 Coming Soon",
            value=MSG_HELP["features"]["upcoming"],
            inline=False
        )
        
        embed.set_footer(text=MSG_HELP["features"]["footer"])
        
        await interaction.response.send_message(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(FeaturesCog(bot))