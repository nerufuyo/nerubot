"""
Utility commands cog - useful tools and calculators
"""
import discord
from discord.ext import commands
from discord import app_commands
import aiohttp
import math
import re
from datetime import datetime, timezone

class Utility(commands.Cog):
    """Utility and tool commands."""
    
    def __init__(self, bot):
        self.bot = bot
    
    @app_commands.command(name="calculate", description="Perform basic calculations")
    @app_commands.describe(expression="Mathematical expression to calculate")
    async def calculate(self, interaction: discord.Interaction, expression: str):
        """Perform basic mathematical calculations."""
        # Clean the expression and only allow safe characters
        safe_chars = re.compile(r'^[0-9+\-*/().\s]+$')
        
        if not safe_chars.match(expression):
            await interaction.response.send_message("❌ Invalid characters in expression. Only numbers and +, -, *, /, (, ) are allowed.")
            return
        
        try:
            # Use eval safely (only after validation)
            result = eval(expression)
            
            embed = discord.Embed(
                title="🧮 Calculator",
                color=discord.Color.blue()
            )
            embed.add_field(name="Expression", value=f"`{expression}`", inline=False)
            embed.add_field(name="Result", value=f"`{result}`", inline=False)
            
            await interaction.response.send_message(embed=embed)
            
        except ZeroDivisionError:
            await interaction.response.send_message("❌ Error: Division by zero!")
        except Exception as e:
            await interaction.response.send_message(f"❌ Error: Invalid expression")
    
    @app_commands.command(name="weather", description="Get weather information for a city")
    @app_commands.describe(city="City name to get weather for")
    async def weather(self, interaction: discord.Interaction, city: str):
        """Get weather information for a city."""
        await interaction.response.defer()
        
        # Note: You'll need to add your OpenWeatherMap API key to .env
        # For now, this will show a placeholder
        embed = discord.Embed(
            title="🌤️ Weather Information",
            description="Weather feature coming soon! Add your OpenWeatherMap API key to enable this feature.",
            color=discord.Color.blue()
        )
        embed.add_field(name="City", value=city, inline=False)
        embed.add_field(name="Note", value="To enable weather, add `WEATHER_API_KEY` to your .env file", inline=False)
        
        await interaction.followup.send(embed=embed)
    
    @app_commands.command(name="userinfo", description="Get information about a user")
    @app_commands.describe(user="User to get information about (optional)")
    async def userinfo(self, interaction: discord.Interaction, user: discord.Member = None):
        """Get information about a user."""
        if user is None:
            user = interaction.user
        
        embed = discord.Embed(
            title=f"👤 User Information: {user.display_name}",
            color=user.color if user.color != discord.Color.default() else discord.Color.blue()
        )
        
        embed.set_thumbnail(url=user.display_avatar.url)
        embed.add_field(name="Username", value=str(user), inline=True)
        embed.add_field(name="ID", value=user.id, inline=True)
        embed.add_field(name="Nickname", value=user.nick or "None", inline=True)
        
        embed.add_field(name="Account Created", value=f"<t:{int(user.created_at.timestamp())}:F>", inline=False)
        embed.add_field(name="Joined Server", value=f"<t:{int(user.joined_at.timestamp())}:F>", inline=False)
        
        if user.premium_since:
            embed.add_field(name="Boosting Since", value=f"<t:{int(user.premium_since.timestamp())}:F>", inline=False)
        
        roles = [role.mention for role in user.roles[1:]]  # Exclude @everyone
        if roles:
            embed.add_field(name=f"Roles ({len(roles)})", value=" ".join(roles), inline=False)
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="serverinfo", description="Get information about this server")
    async def serverinfo(self, interaction: discord.Interaction):
        """Get information about the current server."""
        guild = interaction.guild
        
        embed = discord.Embed(
            title=f"🏰 Server Information: {guild.name}",
            color=discord.Color.gold()
        )
        
        if guild.icon:
            embed.set_thumbnail(url=guild.icon.url)
        
        embed.add_field(name="Server ID", value=guild.id, inline=True)
        embed.add_field(name="Owner", value=guild.owner.mention if guild.owner else "Unknown", inline=True)
        embed.add_field(name="Created", value=f"<t:{int(guild.created_at.timestamp())}:F>", inline=False)
        
        embed.add_field(name="Members", value=guild.member_count, inline=True)
        embed.add_field(name="Channels", value=len(guild.channels), inline=True)
        embed.add_field(name="Roles", value=len(guild.roles), inline=True)
        
        embed.add_field(name="Boost Level", value=guild.premium_tier, inline=True)
        embed.add_field(name="Boosts", value=guild.premium_subscription_count, inline=True)
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="avatar", description="Get a user's avatar")
    @app_commands.describe(user="User to get avatar of (optional)")
    async def avatar(self, interaction: discord.Interaction, user: discord.Member = None):
        """Get a user's avatar."""
        if user is None:
            user = interaction.user
        
        embed = discord.Embed(
            title=f"🖼️ {user.display_name}'s Avatar",
            color=discord.Color.blue()
        )
        embed.set_image(url=user.display_avatar.url)
        embed.add_field(name="Download Link", value=f"[Click here]({user.display_avatar.url})", inline=False)
        
        await interaction.response.send_message(embed=embed)

async def setup(bot):
    await bot.add_cog(Utility(bot))
