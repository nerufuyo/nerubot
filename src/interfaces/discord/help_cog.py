"""
Help commands cog
"""
import discord
from discord.ext import commands
from src.core.utils.messages import (
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


class HelpCog(commands.Cog):
    """Help commands cog."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
    
    @commands.hybrid_command(name="bothelp", description="Show help information")
    async def help_command(self, ctx: commands.Context) -> None:
        """Show help information."""
        embed = discord.Embed(
            title="🤖 NeruBot - Enhanced Discord Bot",
            description="A sophisticated Discord bot with music, news, and AI features!",
            color=discord.Color.blue()
        )
        
        # Music commands
        music_commands = [
            ("/join", "Join your voice channel"),
            ("/leave", "Leave the voice channel"),
            ("/play <song>", "Play a song from YouTube"),
            ("/stop", "Stop music and clear queue"),
            ("/pause", "Pause the current song"),
            ("/resume", "Resume the current song"),
            ("/skip", "Skip the current song"),
            ("/queue", "Show the music queue"),
            ("/nowplaying", "Show current song with status"),
            ("/loop [mode]", "Toggle loop mode (off/single/queue)"),
            ("/247", "Toggle 24/7 mode"),
            ("/clear", "Clear the music queue")
        ]
        
        commands_text = "\\n".join([f"**{cmd}**: {desc}" for cmd, desc in music_commands])
        embed.add_field(name="🎵 Music Commands", value=commands_text, inline=False)
        
        embed.add_field(
            name="🏗️ Architecture", 
            value="Clean, modular design following best practices\\nEasy to maintain and extend", 
            inline=True
        )
        
        embed.add_field(
            name="📊 Status", 
            value="✅ Music Streaming\\n🚧 Additional features coming soon!", 
            inline=True
        )
        
        embed.set_footer(text="NeruBot - Clean Discord Music Bot!")
        
        await ctx.send(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(HelpCog(bot))
