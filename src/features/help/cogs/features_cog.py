"""
Features command cog for NeruBot to display available and upcoming features
"""
import discord
from discord.ext import commands
from discord import app_commands
from src.features.music.services.sources import MusicSource


class FeaturesCog(commands.Cog):
    """Features command cog."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
    
    @app_commands.command(name="features", description="Show bot features and upcoming additions")
    async def features_command(self, interaction: discord.Interaction) -> None:
        """Show features information."""
        embed = discord.Embed(
            title="🚀 NeruBot Features",
            description="Here's what NeruBot can do for your server!",
            color=discord.Color.blue()
        )
        
        # Current features
        embed.add_field(
            name="✅ Current Features",
            value=(
                "**🎵 Music**\n"
                "• Multi-source playback (YouTube, Spotify, SoundCloud)\n"
                "• Advanced queue management\n"
                "• Loop mode (single/queue)\n"
                "• 24/7 mode\n"
                "• High-quality audio with volume control\n\n"
                
                "**🤖 Bot**\n"
                "• Slash commands support\n"
                "• Interactive help system\n"
                "• Clean error handling\n"
            ),
            inline=False
        )
        
        # Music sources
        embed.add_field(
            name="🎵 Music Sources",
            value=(
                "• ▶️ YouTube\n"
                "• 💚 Spotify\n"
                "• 🧡 SoundCloud\n"
                "• 🔗 Direct audio links\n"
            ),
            inline=True
        )
        
        # Upcoming features
        embed.add_field(
            name="🚧 Coming Soon",
            value=(
                "• 📰 News feeds\n"
                "• 🎮 Game statistics\n"
                "• 📊 Server analytics\n"
                "• 📅 Event scheduling\n"
            ),
            inline=True
        )
        
        embed.set_footer(text="Use /help to see all available commands!")
        
        await interaction.response.send_message(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(FeaturesCog(bot))