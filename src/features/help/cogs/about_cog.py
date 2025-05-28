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
        
        # Get bot stats
        guild_count = len(self.bot.guilds)
        total_members = sum(guild.member_count for guild in self.bot.guilds)
        command_count = len(self.bot.tree.get_commands())
        
        embed = discord.Embed(
            title=f"🎵 About {BOT_CONFIG['name']} - Your Friendly Music Companion!",
            description=(
                "Hi there! I'm **NeruBot** - a powerful, feature-rich Discord bot designed to bring music, "
                "community engagement, and fun to your server! 🎉\n\n"
                "I'm built with love to provide the best experience for your Discord community with "
                "high-quality audio streaming, anonymous confessions, news updates, and much more!"
            ),
            color=0x7289DA
        )
        
        # Set the banner image
        embed.set_image(url="https://imgur.com/yh3j7PK.png")
        
        # Set the bot's profile picture as thumbnail
        embed.set_thumbnail(url="https://imgur.com/7IqhTL0.png")
        
        # Core Features - What makes me awesome!
        embed.add_field(
            name="🌟 What Makes Me Special",
            value=(
                "🎵 **Multi-Platform Music** - Stream from YouTube, Spotify & SoundCloud\n"
                "📝 **Anonymous Confessions** - Safe space for community sharing\n"
                "📰 **News Integration** - Stay updated with RSS feeds\n"
                "🎛️ **Advanced Audio** - High-quality playback with queue management\n"
                "🔄 **24/7 Mode** - I can stay in your voice channel all day!\n"
                "⚡ **Lightning Fast** - Optimized for speed and reliability"
            ),
            inline=False
        )
        
        # Developer & Author Information
        embed.add_field(
            name="👨‍💻 Created By",
            value=(
                "**nerufuyo** - A passionate developer who loves creating amazing Discord experiences!\n\n"
                "🎯 *Vision:* To build the most user-friendly and feature-rich Discord bot\n"
                "💡 *Mission:* Making Discord servers more engaging and entertaining\n"
                "❤️ *Passion:* Combining clean code architecture with awesome user experience"
            ),
            inline=False
        )
        
        # Live Statistics
        embed.add_field(
            name="📊 Live Stats",
            value=(
                f"🏠 **Servers:** {guild_count:,}\n"
                f"👥 **Users:** {total_members:,}\n"
                f"⚡ **Commands:** {command_count}\n"
                f"⏱️ **Uptime:** {uptime_str}\n"
                f"💾 **Memory:** {memory_usage:.1f} MB"
            ),
            inline=True
        )
        
        # Technical Excellence
        embed.add_field(
            name="⚙️ Built With",
            value=(
                f"🐍 **Python** {platform.python_version()}\n"
                f"🔗 **discord.py** {discord.__version__}\n"
                f"🏗️ **Clean Architecture**\n"
                f"🎵 **FFmpeg Audio**\n"
                f"☁️ **Async Programming**"
            ),
            inline=True
        )
        
        # Special Features Highlight
        embed.add_field(
            name="🎉 Why Users Love Me",
            value=(
                "✨ **Easy to Use** - Simple slash commands for everything\n"
                "🛡️ **Reliable** - Built to handle high-traffic servers\n"
                "🎨 **Beautiful UI** - Rich embeds and interactive components\n"
                "🔒 **Privacy First** - Anonymous features with proper moderation\n"
                "🆓 **Completely Free** - No premium features, everything included!"
            ),
            inline=False
        )
        
        # Call to Action & Support
        embed.add_field(
            name="🚀 Get Started",
            value=(
                "Ready to enhance your server? Here's how to begin:\n"
                "• Type `/help` to see all my amazing features\n"
                "• Use `/play` to start jamming with music\n"
                "• Try `/confess` for anonymous community sharing\n"
                "• Check `/features` for detailed capabilities\n\n"
                "Need help? I'm designed to be intuitive and user-friendly!"
            ),
            inline=False
        )
        
        embed.set_footer(
            text="Made with ❤️ by nerufuyo | Thank you for choosing NeruBot!",
            icon_url="https://imgur.com/7IqhTL0.png"
        )
        
        await interaction.response.send_message(embed=embed)


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(AboutCog(bot))
