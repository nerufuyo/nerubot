"""
Commands reference cog for NeruBot with compact command listing
"""
import discord
from discord.ext import commands
from discord import app_commands
from src.config.messages import MSG_HELP
from src.config.settings import DISCORD_CONFIG


class CommandsCog(commands.Cog):
    """Commands reference cog."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
    
    @app_commands.command(name="commands", description=MSG_HELP["commands"]["commands"])
    async def commands_command(self, interaction: discord.Interaction) -> None:
        """Show a compact command reference card."""
        embed = discord.Embed(
            title=MSG_HELP["command_card"]["title"],
            description=MSG_HELP["command_card"]["description"],
            color=DISCORD_CONFIG["colors"]["info"]
        )
        
        embed.set_thumbnail(url="https://imgur.com/7IqhTL0.png")
        
        # Music commands section
        embed.add_field(
            name="🎵 Music Commands", 
            value=MSG_HELP["command_card"]["music_commands"], 
            inline=False
        )
        
        # Confession commands section
        embed.add_field(
            name="📝 Confession Commands", 
            value=MSG_HELP["command_card"]["confession_commands"], 
            inline=False
        )
        
        # News commands section
        embed.add_field(
            name="📰 News Commands", 
            value=MSG_HELP["command_card"]["news_commands"], 
            inline=False
        )
        
        # General commands
        embed.add_field(
            name="🤖 General Commands", 
            value=MSG_HELP["command_card"]["general_commands"], 
            inline=False
        )
        
        # Tip section
        embed.add_field(
            name="💡 Tips", 
            value=MSG_HELP["command_card"]["tips"], 
            inline=False
        )
        
        embed.set_footer(text=MSG_HELP["command_card"]["footer"], icon_url="https://imgur.com/7IqhTL0.png")
        
        await interaction.response.send_message(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(CommandsCog(bot))
