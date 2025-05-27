"""
About command cog for NeruBot
"""
import discord
from discord.ext import commands
from discord import app_commands
import platform
import psutil
import time
from src.config.messages import MSG_HELP, BOT_INFO
from src.config.settings import BOT_CONFIG, DISCORD_CONFIG


class AboutCog(commands.Cog):
    """About command cog."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
        self.start_time = time.time()
    
    @app_commands.command(name="about", description=MSG_HELP["commands"]["about"])
    async def about_command(self, interaction: discord.Interaction) -> None:
        """Show information about the bot."""
        # Calculate uptime
        uptime = time.time() - self.start_time
        days, remainder = divmod(int(uptime), 86400)
        hours, remainder = divmod(remainder, 3600)
        minutes, seconds = divmod(remainder, 60)
        uptime_str = f"{days}d {hours}h {minutes}m {seconds}s"
        
        # Get resource usage
        process = psutil.Process()
        memory_usage = process.memory_info().rss / 1024 / 1024  # Convert to MB
        
        embed = discord.Embed(
            title=f"🤖 About {BOT_CONFIG['name']}",
            description=BOT_CONFIG['description'],
            color=DISCORD_CONFIG["colors"]["info"]
        )
        
        # Bot information
        embed.add_field(
            name="💡 Features",
            value=MSG_HELP["about"]["features"],
            inline=False
        )
        
        # System information
        embed.add_field(
            name="⚙️ System",
            value=f"• Python: {platform.python_version()}\n"
                 f"• discord.py: {discord.__version__}\n"
                 f"• Memory: {memory_usage:.2f} MB\n"
                 f"• Uptime: {uptime_str}",
            inline=True
        )
        
        # Stats information
        guild_count = len(self.bot.guilds)
        total_members = sum(guild.member_count for guild in self.bot.guilds)
        
        embed.add_field(
            name="📊 Stats",
            value=f"• Servers: {guild_count}\n"
                 f"• Users: {total_members}\n"
                 f"• Commands: {len(self.bot.tree.get_commands())}",
            inline=True
        )
        
        # Links and credits
        embed.add_field(
            name="🔗 Links",
            value=MSG_HELP["about"]["links"],
            inline=False
        )
        
        embed.set_footer(text=MSG_HELP["about"]["footer"])
        
        await interaction.response.send_message(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(AboutCog(bot))
