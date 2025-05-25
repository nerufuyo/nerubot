"""News cog for Discord commands."""
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional, Dict

from src.features.news.services.news_service import NewsService
from src.core.utils.logging_utils import get_logger

logger = get_logger(__name__)

class NewsCog(commands.Cog):
    """Discord commands for news functionality."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
        self.news_service = NewsService()
        self.news_channels: Dict[int, int] = {}  # guild_id -> channel_id
        self.auto_post_enabled: Dict[int, bool] = {}  # guild_id -> enabled status
        
    async def cog_load(self) -> None:
        """Called when the cog is loaded."""
        # Set the callback for new items
        self.news_service.set_new_items_callback(self._handle_new_items)
        await self.news_service.start()
        logger.info("News cog loaded")
    
    async def cog_unload(self) -> None:
        """Called when the cog is unloaded."""
        await self.news_service.stop()
        logger.info("News cog unloaded")
    
    async def _handle_new_items(self, new_items) -> None:
        """Handle new news items by posting them to configured channels."""
        for guild_id, channel_id in self.news_channels.items():
            if not self.auto_post_enabled.get(guild_id, False):
                continue
                
            try:
                channel = self.bot.get_channel(channel_id)
                if not channel:
                    logger.warning(f"Channel {channel_id} not found for guild {guild_id}")
                    continue
                
                # Post only the latest news item to avoid spam
                if new_items:
                    latest_item = new_items[0]  # Get the most recent item
                    embed = discord.Embed.from_dict(latest_item.to_embed())
                    await channel.send("📰 **Breaking News!**", embed=embed)
                    
            except Exception as e:
                logger.error(f"Error posting news to channel {channel_id}: {e}")
    
    @commands.hybrid_group(name="news", description="News commands")
    async def news(self, ctx: commands.Context) -> None:
        """News commands group."""
        if ctx.invoked_subcommand is None:
            await ctx.send("Please specify a news subcommand. Use `/news help` for more information.")
    
    @news.command(name="latest", description="Get the latest news")
    @app_commands.describe(count="Number of news items to show")
    async def latest(self, ctx: commands.Context, count: Optional[int] = 5) -> None:
        """Get the latest news."""
        count = min(max(1, count), 10)  # Limit between 1 and 10
        latest_news = self.news_service.get_latest_news(count)
        
        if not latest_news:
            await ctx.send("No news items available yet. Please try again later.")
            return
        
        for news_item in latest_news:
            embed = discord.Embed.from_dict(news_item.to_embed())
            await ctx.send(embed=embed)
    
    @news.command(name="sources", description="List all news sources")
    async def sources(self, ctx: commands.Context) -> None:
        """List all configured news sources."""
        sources = self.news_service.get_news_sources()
        
        if not sources:
            await ctx.send("No news sources configured.")
            return
        
        embed = discord.Embed(
            title="News Sources",
            description="List of configured news sources",
            color=0x3498db
        )
        
        for source_name, feed_url in sources.items():
            embed.add_field(name=source_name, value=feed_url, inline=False)
        
        await ctx.send(embed=embed)
    
    @news.command(name="add", description="Add a news source")
    @app_commands.describe(
        name="Name of the news source",
        feed_url="URL of the RSS feed"
    )
    @commands.has_permissions(administrator=True)
    async def add_source(self, ctx: commands.Context, name: str, feed_url: str) -> None:
        """Add a news source."""
        success = self.news_service.add_news_source(name, feed_url)
        
        if success:
            await ctx.send(f"Added news source: {name}")
        else:
            await ctx.send(f"News source already exists: {name}")
    
    @news.command(name="remove", description="Remove a news source")
    @app_commands.describe(name="Name of the news source to remove")
    @commands.has_permissions(administrator=True)
    async def remove_source(self, ctx: commands.Context, name: str) -> None:
        """Remove a news source."""
        success = self.news_service.remove_news_source(name)
        
        if success:
            await ctx.send(f"Removed news source: {name}")
        else:
            await ctx.send(f"News source not found: {name}")
    
    @news.command(name="set-channel", description="Set the channel for automatic news updates")
    @app_commands.describe(channel="The channel to send news updates to")
    @commands.has_permissions(administrator=True)
    async def set_channel(self, ctx: commands.Context, channel: Optional[discord.TextChannel] = None) -> None:
        """Set the channel for automatic news updates."""
        channel = channel or ctx.channel
        guild_id = ctx.guild.id
        
        self.news_channels[guild_id] = channel.id
        self.auto_post_enabled[guild_id] = True  # Enable auto-posting when setting channel
        
        await ctx.send(f"News updates will be sent to {channel.mention}. Auto-posting is now enabled.")
    
    @news.command(name="start", description="Start automatic news updates")
    @commands.has_permissions(administrator=True)
    async def start_updates(self, ctx: commands.Context) -> None:
        """Start automatic news updates."""
        guild_id = ctx.guild.id
        
        if guild_id not in self.news_channels:
            await ctx.send("Please set a news channel first using `/news set-channel`.")
            return
        
        self.auto_post_enabled[guild_id] = True
        await ctx.send("Automatic news updates have been started!")
    
    @news.command(name="status", description="Show current news configuration")
    async def status(self, ctx: commands.Context) -> None:
        """Show current news configuration."""
        guild_id = ctx.guild.id
        
        embed = discord.Embed(
            title="News Configuration Status",
            color=0x3498db
        )
        
        # Check if news channel is set
        if guild_id in self.news_channels:
            channel = self.bot.get_channel(self.news_channels[guild_id])
            channel_name = channel.mention if channel else "Channel not found"
            embed.add_field(name="News Channel", value=channel_name, inline=False)
        else:
            embed.add_field(name="News Channel", value="Not set", inline=False)
        
        # Check if auto-posting is enabled
        auto_post_status = "Enabled" if self.auto_post_enabled.get(guild_id, False) else "Disabled"
        embed.add_field(name="Auto-posting", value=auto_post_status, inline=False)
        
        # News service status
        service_status = "Running" if self.news_service.running else "Stopped"
        embed.add_field(name="News Service", value=service_status, inline=False)
        
        # Number of news sources
        source_count = len(self.news_service.get_news_sources())
        embed.add_field(name="News Sources", value=f"{source_count} configured", inline=False)
        
        # Latest news count
        news_count = len(self.news_service.get_latest_news(100))
        embed.add_field(name="Available News", value=f"{news_count} items", inline=False)
        
        await ctx.send(embed=embed)
    
    @news.command(name="stop", description="Stop automatic news updates")
    @commands.has_permissions(administrator=True)
    async def stop_updates(self, ctx: commands.Context) -> None:
        """Stop automatic news updates."""
        guild_id = ctx.guild.id
        
        if guild_id in self.auto_post_enabled:
            self.auto_post_enabled[guild_id] = False
            await ctx.send("Automatic news updates have been stopped.")
        else:
            await ctx.send("Automatic news updates are not currently enabled.")
    
    @news.command(name="help", description="Show help for news commands")
    async def help_command(self, ctx: commands.Context) -> None:
        """Show help for news commands."""
        embed = discord.Embed(
            title="News Commands",
            description="Commands for the news feature",
            color=0x3498db
        )
        
        embed.add_field(
            name="/news latest [count]",
            value="Get the latest news items. Optionally specify how many items to show (default: 5).",
            inline=False
        )
        
        embed.add_field(
            name="/news sources",
            value="List all configured news sources.",
            inline=False
        )
        
        embed.add_field(
            name="/news status",
            value="Show current news configuration and status.",
            inline=False
        )
        
        embed.add_field(
            name="/news set-channel [channel]",
            value="Set the channel for automatic news updates and enable auto-posting. If no channel is specified, the current channel is used.",
            inline=False
        )
        
        embed.add_field(
            name="/news start",
            value="Start automatic news updates (admin only).",
            inline=False
        )
        
        embed.add_field(
            name="/news stop",
            value="Stop automatic news updates (admin only).",
            inline=False
        )
        
        embed.add_field(
            name="/news add <name> <feed_url>",
            value="Add a news source (admin only).",
            inline=False
        )
        
        embed.add_field(
            name="/news remove <name>",
            value="Remove a news source (admin only).",
            inline=False
        )
        
        await ctx.send(embed=embed)

async def setup(bot: commands.Bot) -> None:
    """Setup function for the news cog."""
    await bot.add_cog(NewsCog(bot))
