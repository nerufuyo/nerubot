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
            (f"{self.bot.command_prefix}join", HELP_JOIN_DESC),
            (f"{self.bot.command_prefix}leave", HELP_LEAVE_DESC),
            (f"{self.bot.command_prefix}play <song>", HELP_PLAY_DESC),
            (f"{self.bot.command_prefix}stop", HELP_STOP_DESC),
            (f"{self.bot.command_prefix}pause", HELP_PAUSE_DESC),
            (f"{self.bot.command_prefix}resume", HELP_RESUME_DESC),
            (f"{self.bot.command_prefix}skip", HELP_SKIP_DESC),
            (f"{self.bot.command_prefix}volume <0-100>", HELP_VOLUME_DESC),
            (f"{self.bot.command_prefix}now", HELP_NOW_DESC),
            (f"{self.bot.command_prefix}queue [page]", HELP_QUEUE_DESC),
            (f"{self.bot.command_prefix}remove <index>", HELP_REMOVE_DESC),
            (f"{self.bot.command_prefix}shuffle", HELP_SHUFFLE_DESC),
            (f"{self.bot.command_prefix}loop <off/song/queue>", HELP_LOOP_DESC)
        ]
        
        commands_text = "\\n".join([f"**{cmd}**: {desc}" for cmd, desc in music_commands])
        embed.add_field(name="🎵 Music Commands", value=commands_text, inline=False)
        
        # News commands
        news_commands = [
            ("/news", "Get latest news from RSS feeds"),
            ("/news category:tech", "Get technology news"),
            ("/news-categories", "Show available categories"),
            ("/news-source source:BBC", "Get news from specific source")
        ]
        
        news_text = "\\n".join([f"**{cmd}**: {desc}" for cmd, desc in news_commands])
        embed.add_field(name="📰 News Commands", value=news_text, inline=False)
        
        # Future features
        future_commands = [
            ("/quote", "AI-generated inspirational quotes (Coming Soon)"),
            ("/profile", "User profile management (Coming Soon)"),
            ("/confess", "Anonymous confessions (Coming Soon)")
        ]
        
        future_text = "\\n".join([f"**{cmd}**: {desc}" for cmd, desc in future_commands])
        embed.add_field(name="🔮 Coming Soon", value=future_text, inline=False)
        
        embed.add_field(
            name="🏗️ Architecture", 
            value="Modular design following DRY & KISS principles\\nEasy to maintain and extend", 
            inline=True
        )
        
        embed.add_field(
            name="📊 Features", 
            value="✅ Music Streaming\\n✅ RSS News Feeds\\n🚧 AI Quotes (Coming Soon)\\n🚧 User Profiles\\n🚧 Confessions", 
            inline=True
        )
        
        embed.set_footer(text="NeruBot v2.0 - Enhanced with modular architecture!")
        
        await ctx.send(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(HelpCog(bot))
