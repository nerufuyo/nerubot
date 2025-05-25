"""
Quotes commands cog for the Discord bot.
AI-powered quotes using DeepSeek API with fallback system.
"""
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional, List
import asyncio

from src.features.quotes.services.quotes_service import QuotesService
from src.features.quotes.models.quote import QuoteRequest

import logging
logger = logging.getLogger(__name__)


class QuotesCog(commands.Cog):
    """AI-powered quotes commands using DeepSeek API."""
    
    def __init__(self, bot):
        self.bot = bot
        self.quotes_service = QuotesService()
        logger.info("QuotesCog loaded with DeepSeek AI integration")
    
    @app_commands.command(name="quote", description="Get an AI-generated inspirational quote")
    @app_commands.describe(
        category="Quote category (motivation, wisdom, technology, philosophy, humor)",
        mood="The mood/tone for the quote",
        length="Quote length preference"
    )
    async def quote(
        self, 
        interaction: discord.Interaction, 
        category: Optional[str] = None,
        mood: Optional[str] = None,
        length: Optional[str] = "medium"
    ):
        """Get an AI-generated quote with optional category, mood, and length."""
        await interaction.response.defer()
        
        try:
            # Create quote request
            request = QuoteRequest(
                category=category,
                mood=mood,
                length=length or "medium"
            )
            
            # Get quote from service
            quote = await self.quotes_service.get_quote(request)
            
            # Create embed
            embed = discord.Embed(
                title="💭 AI-Generated Quote",
                description=f"*\"{quote.content}\"*",
                color=discord.Color.purple()
            )
            
            if quote.author:
                embed.add_field(
                    name="Author",
                    value=f"**{quote.author}**",
                    inline=True
                )
            
            if quote.category:
                embed.add_field(
                    name="Category",
                    value=f"**{quote.category.title()}**",
                    inline=True
                )
            
            # Add source indicator
            source_emoji = "🤖" if quote.source == "deepseek_ai" else "📚"
            source_text = "AI Generated" if quote.source == "deepseek_ai" else "Curated Collection"
            
            embed.add_field(
                name="Source",
                value=f"{source_emoji} {source_text}",
                inline=True
            )
            
            embed.set_footer(text=f"Created at {quote.created_at.strftime('%H:%M:%S')}")
            
            await interaction.followup.send(embed=embed)
            
        except Exception as e:
            logger.error(f"Error generating quote: {e}")
            
            error_embed = discord.Embed(
                title="❌ Quote Generation Failed",
                description="Sorry, I couldn't generate a quote right now. Please try again later.",
                color=discord.Color.red()
            )
            
            await interaction.followup.send(embed=error_embed, ephemeral=True)
    
    @app_commands.command(name="random-quote", description="Get a random inspirational quote")
    async def random_quote(self, interaction: discord.Interaction):
        """Get a random inspirational quote."""
        await interaction.response.defer()
        
        try:
            quote = await self.quotes_service.get_random_quote()
            
            embed = discord.Embed(
                title="🎲 Random Quote",
                description=f"*\"{quote.content}\"*",
                color=discord.Color.gold()
            )
            
            if quote.author:
                embed.add_field(
                    name="Author",
                    value=f"**{quote.author}**",
                    inline=True
                )
            
            if quote.category:
                embed.add_field(
                    name="Category",
                    value=f"**{quote.category.title()}**",
                    inline=True
                )
            
            source_emoji = "🤖" if quote.source == "deepseek_ai" else "📚"
            source_text = "AI Generated" if quote.source == "deepseek_ai" else "Curated Collection"
            
            embed.add_field(
                name="Source",
                value=f"{source_emoji} {source_text}",
                inline=True
            )
            
            embed.set_footer(text="Use /quote for customized quotes!")
            
            await interaction.followup.send(embed=embed)
            
        except Exception as e:
            logger.error(f"Error getting random quote: {e}")
            
            error_embed = discord.Embed(
                title="❌ Random Quote Failed",
                description="Sorry, I couldn't get a random quote right now. Please try again later.",
                color=discord.Color.red()
            )
            
            await interaction.followup.send(embed=error_embed, ephemeral=True)
    
    @app_commands.command(name="quote-categories", description="View available quote categories")
    async def quote_categories(self, interaction: discord.Interaction):
        """Show available quote categories."""
        categories = self.quotes_service.get_available_categories()
        moods = self.quotes_service.get_available_moods()
        
        embed = discord.Embed(
            title="💭 Quote System Information",
            description="Available categories, moods, and options for quote generation",
            color=discord.Color.blue()
        )
        
        embed.add_field(
            name="📂 Available Categories",
            value="• " + "\n• ".join([cat.title() for cat in categories]),
            inline=False
        )
        
        embed.add_field(
            name="🎭 Available Moods",
            value="• " + "\n• ".join([mood.title() for mood in moods]),
            inline=False
        )
        
        embed.add_field(
            name="📏 Length Options",
            value="• **Short** - Concise quotes (under 20 words)\n• **Medium** - Moderate length (20-40 words)\n• **Long** - Detailed quotes (40-80 words)",
            inline=False
        )
        
        embed.add_field(
            name="🔧 Usage Examples",
            value="• `/quote` - Get a random quote\n• `/quote category:motivation` - Motivational quote\n• `/quote mood:humorous length:short` - Short funny quote\n• `/random-quote` - Completely random quote",
            inline=False
        )
        
        embed.set_footer(text="Powered by DeepSeek AI with fallback quotes")
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="daily-quote", description="Get multiple quotes for daily inspiration")
    @app_commands.describe(
        category="Category for all quotes",
        count="Number of quotes to generate (1-5)"
    )
    async def daily_quote(
        self, 
        interaction: discord.Interaction, 
        category: Optional[str] = None,
        count: Optional[int] = 3
    ):
        """Get multiple quotes for daily inspiration."""
        await interaction.response.defer()
        
        # Validate count
        count = max(1, min(count or 3, 5))
        
        try:
            if category:
                quotes = await self.quotes_service.get_quotes_by_category(category, count)
            else:
                # Get quotes from different categories
                quotes = []
                for i in range(count):
                    quote = await self.quotes_service.get_random_quote()
                    quotes.append(quote)
            
            embed = discord.Embed(
                title=f"📅 Daily Inspiration ({count} Quotes)",
                description=f"Here are {count} quotes to inspire your day!",
                color=discord.Color.green()
            )
            
            for i, quote in enumerate(quotes, 1):
                author_text = f" - {quote.author}" if quote.author else ""
                source_emoji = "🤖" if quote.source == "deepseek_ai" else "📚"
                
                embed.add_field(
                    name=f"Quote #{i} {source_emoji}",
                    value=f"*\"{quote.content}\"*{author_text}",
                    inline=False
                )
            
            if category:
                embed.set_footer(text=f"Category: {category.title()} • Generated at {quotes[0].created_at.strftime('%H:%M:%S')}")
            else:
                embed.set_footer(text=f"Mixed categories • Generated at {quotes[0].created_at.strftime('%H:%M:%S')}")
            
            await interaction.followup.send(embed=embed)
            
        except Exception as e:
            logger.error(f"Error generating daily quotes: {e}")
            
            error_embed = discord.Embed(
                title="❌ Daily Quote Generation Failed",
                description="Sorry, I couldn't generate daily quotes right now. Please try again later.",
                color=discord.Color.red()
            )
            
            await interaction.followup.send(embed=error_embed, ephemeral=True)
    
    # Autocomplete functions
    @quote.autocomplete('category')
    @daily_quote.autocomplete('category')
    async def category_autocomplete(
        self,
        interaction: discord.Interaction,
        current: str,
    ) -> List[app_commands.Choice[str]]:
        """Autocomplete for quote categories."""
        categories = self.quotes_service.get_available_categories()
        return [
            app_commands.Choice(name=category.title(), value=category)
            for category in categories
            if current.lower() in category.lower()
        ][:25]  # Discord limit
    
    @quote.autocomplete('mood')
    async def mood_autocomplete(
        self,
        interaction: discord.Interaction,
        current: str,
    ) -> List[app_commands.Choice[str]]:
        """Autocomplete for quote moods."""
        moods = self.quotes_service.get_available_moods()
        return [
            app_commands.Choice(name=mood.title(), value=mood)
            for mood in moods
            if current.lower() in mood.lower()
        ][:25]  # Discord limit
    
    @quote.autocomplete('length')
    async def length_autocomplete(
        self,
        interaction: discord.Interaction,
        current: str,
    ) -> List[app_commands.Choice[str]]:
        """Autocomplete for quote lengths."""
        lengths = ["short", "medium", "long"]
        return [
            app_commands.Choice(name=length.title(), value=length)
            for length in lengths
            if current.lower() in length.lower()
        ][:25]  # Discord limit


async def setup(bot):
    """Setup function for the cog."""
    await bot.add_cog(QuotesCog(bot))
