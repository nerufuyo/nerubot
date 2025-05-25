"""
News commands cog for the Discord bot.
"""
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional, List
import asyncio

from src.features.news.services.news_service import NewsService
from src.features.news.models.news import NewsArticle

import logging
logger = logging.getLogger(__name__)


class NewsCog(commands.Cog):
    """News commands for fetching RSS feeds."""
    
    def __init__(self, bot):
        self.bot = bot
        self.news_service = NewsService()
    
    async def cog_load(self):
        """Called when the cog is loaded."""
        logger.info("News cog loaded")
    
    async def cog_unload(self):
        """Called when the cog is unloaded."""
        # Clean up any resources if needed
        pass
    
    @app_commands.command(name="news", description="Get latest news from RSS feeds")
    @app_commands.describe(
        category="News category (optional)",
        count="Number of articles to show (1-10, default: 5)"
    )
    async def news(
        self, 
        interaction: discord.Interaction, 
        category: Optional[str] = None,
        count: Optional[int] = 5
    ):
        """Get latest news articles."""
        await interaction.response.defer(thinking=True)
        
        # Validate count
        count = max(1, min(count or 5, 10))
        
        try:
            async with self.news_service as service:
                articles = await service.get_latest_news(category, count)
            
            if not articles:
                embed = discord.Embed(
                    title="📰 No News Found",
                    description="No news articles found for the specified criteria.",
                    color=discord.Color.orange()
                )
                await interaction.followup.send(embed=embed)
                return
            
            # Create news embed
            embed = self._create_news_embed(articles, category)
            await interaction.followup.send(embed=embed)
            
        except Exception as e:
            logger.error(f"Error fetching news: {e}")
            embed = discord.Embed(
                title="❌ Error",
                description="Sorry, I couldn't fetch the news right now. Please try again later.",
                color=discord.Color.red()
            )
            await interaction.followup.send(embed=embed)
    
    @app_commands.command(name="news-categories", description="Show available news categories")
    async def news_categories(self, interaction: discord.Interaction):
        """Show available news categories."""
        try:
            async with self.news_service as service:
                categories = service.get_available_categories()
                feeds = service.get_feed_names()
            
            embed = discord.Embed(
                title="📂 Available News Categories",
                color=discord.Color.blue()
            )
            
            if categories:
                categories_text = "\\n".join([f"• **{cat.title()}**" for cat in categories])
                embed.add_field(
                    name="Categories",
                    value=categories_text,
                    inline=False
                )
            
            if feeds:
                feeds_text = "\\n".join([f"• {feed}" for feed in feeds[:10]])
                if len(feeds) > 10:
                    feeds_text += f"\\n... and {len(feeds) - 10} more"
                
                embed.add_field(
                    name="Available Sources",
                    value=feeds_text,
                    inline=False
                )
            
            embed.add_field(
                name="Usage",
                value="`/news` - Get general news\\n`/news category:technology` - Get tech news",
                inline=False
            )
            
            await interaction.response.send_message(embed=embed)
            
        except Exception as e:
            logger.error(f"Error getting categories: {e}")
            await interaction.response.send_message(
                "❌ Error getting news categories. Please try again later.",
                ephemeral=True
            )
    
    @app_commands.command(name="news-source", description="Get news from a specific source")
    @app_commands.describe(
        source="News source name",
        count="Number of articles to show (1-5, default: 3)"
    )
    async def news_source(
        self, 
        interaction: discord.Interaction, 
        source: str,
        count: Optional[int] = 3
    ):
        """Get news from a specific RSS source."""
        await interaction.response.defer(thinking=True)
        
        # Validate count
        count = max(1, min(count or 3, 5))
        
        try:
            async with self.news_service as service:
                # Find the feed
                feed = None
                for feed_name, feed_obj in service.feeds.items():
                    if source.lower() in feed_name or source.lower() in feed_obj.name.lower():
                        feed = feed_obj
                        break
                
                if not feed:
                    embed = discord.Embed(
                        title="❌ Source Not Found",
                        description=f"Could not find news source: `{source}`\\n\\nUse `/news-categories` to see available sources.",
                        color=discord.Color.red()
                    )
                    await interaction.followup.send(embed=embed)
                    return
                
                articles = await service.fetch_feed(feed, count)
                
                if not articles:
                    embed = discord.Embed(
                        title="📰 No Articles",
                        description=f"No articles found from **{feed.name}**.",
                        color=discord.Color.orange()
                    )
                    await interaction.followup.send(embed=embed)
                    return
                
                # Create source-specific embed
                embed = self._create_source_embed(articles, feed.name)
                await interaction.followup.send(embed=embed)
                
        except Exception as e:
            logger.error(f"Error fetching source news: {e}")
            embed = discord.Embed(
                title="❌ Error",
                description="Sorry, I couldn't fetch news from that source. Please try again later.",
                color=discord.Color.red()
            )
            await interaction.followup.send(embed=embed)
    
    def _create_news_embed(self, articles: List[NewsArticle], category: Optional[str] = None) -> discord.Embed:
        """Create an embed for news articles."""
        title = "📰 Latest News"
        if category:
            title += f" - {category.title()}"
        
        embed = discord.Embed(
            title=title,
            color=discord.Color.blue(),
            timestamp=articles[0].published if articles else None
        )
        
        for i, article in enumerate(articles[:5], 1):
            # Format published time
            time_str = article.published.strftime("%m/%d %H:%M")
            
            field_name = f"{i}. {article.title[:80]}{'...' if len(article.title) > 80 else ''}"
            field_value = (
                f"**Source:** {article.source}\\n"
                f"**Time:** {time_str} UTC\\n"
                f"**Description:** {article.short_description}\\n"
                f"[Read More]({article.link})"
            )
            
            embed.add_field(
                name=field_name,
                value=field_value,
                inline=False
            )
        
        embed.set_footer(text=f"Showing {len(articles)} articles • Use /news-categories for more options")
        
        return embed
    
    def _create_source_embed(self, articles: List[NewsArticle], source_name: str) -> discord.Embed:
        """Create an embed for articles from a specific source."""
        embed = discord.Embed(
            title=f"📰 {source_name}",
            color=discord.Color.green(),
            timestamp=articles[0].published if articles else None
        )
        
        for i, article in enumerate(articles, 1):
            time_str = article.published.strftime("%m/%d %H:%M")
            
            field_name = f"{i}. {article.title[:80]}{'...' if len(article.title) > 80 else ''}"
            field_value = (
                f"**Time:** {time_str} UTC\\n"
                f"{article.short_description}\\n"
                f"[Read More]({article.link})"
            )
            
            embed.add_field(
                name=field_name,
                value=field_value,
                inline=False
            )
        
        embed.set_footer(text=f"From {source_name}")
        
        return embed

    # Add autocomplete for categories
    @news.autocomplete('category')
    async def news_category_autocomplete(self, interaction: discord.Interaction, current: str):
        """Autocomplete for news categories."""
        try:
            async with self.news_service as service:
                categories = service.get_available_categories()
            
            return [
                app_commands.Choice(name=cat.title(), value=cat)
                for cat in categories
                if current.lower() in cat.lower()
            ][:25]
        except:
            return []
    
    # Add autocomplete for news sources
    @news_source.autocomplete('source')
    async def news_source_autocomplete(self, interaction: discord.Interaction, current: str):
        """Autocomplete for news sources."""
        try:
            async with self.news_service as service:
                sources = service.get_feed_names()
            
            return [
                app_commands.Choice(name=source, value=source)
                for source in sources
                if current.lower() in source.lower()
            ][:25]
        except:
            return []


async def setup(bot):
    """Setup function for the cog."""
    await bot.add_cog(NewsCog(bot))
