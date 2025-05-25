"""
Commands reference cog for NeruBot with compact command listing
"""
import discord
from discord.ext import commands
from discord import app_commands


class CommandsCog(commands.Cog):
    """Commands reference cog."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
    
    @app_commands.command(name="commands", description="Show compact command reference")
    async def commands_command(self, interaction: discord.Interaction) -> None:
        """Show a compact command reference card."""
        embed = discord.Embed(
            title="📋 NeruBot Command Reference",
            description="Quick reference for all available commands",
            color=discord.Color.blue()
        )
        
        # Music commands section
        music_commands = (
            "`/play` - Play a song\n"
            "`/pause` - Pause playback\n"
            "`/resume` - Resume playback\n"
            "`/skip` - Skip current song\n"
            "`/stop` - Stop and clear queue\n"
            "`/queue` - Show music queue\n"
            "`/nowplaying` - Show current song\n"
            "`/clear` - Clear queue\n"
            "`/loop` - Toggle loop mode\n"
            "`/247` - Toggle 24/7 mode\n"
            "`/join` - Join voice channel\n"
            "`/leave` - Leave voice channel\n"
            "`/sources` - Show music sources"
        )
        embed.add_field(name="🎵 Music Commands", value=music_commands, inline=False)
        
        # General commands
        general_commands = (
            "`/help` - Detailed help pages\n"
            "`/commands` - This command card\n"
            "`/about` - Bot information\n"
            "`/features` - Feature showcase"
        )
        embed.add_field(name="🤖 General Commands", value=general_commands, inline=False)
        
        # Tip section
        tips = (
            "**Pro Tips:**\n"
            "• Use `/play` with Spotify, YouTube or SoundCloud links\n"
            "• Try `/loop queue` to repeat your playlist\n"
            "• Use `/sources` to see all supported music sources"
        )
        embed.add_field(name="💡 Tips", value=tips, inline=False)
        
        embed.set_footer(text="Use /help for more detailed command information")
        
        await interaction.response.send_message(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(CommandsCog(bot))
