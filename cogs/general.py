"""
General commands cog - basic bot functionality
"""
import discord
from discord.ext import commands
from discord import app_commands
import time
import psutil
import sys

class General(commands.Cog):
    """General bot commands."""
    
    def __init__(self, bot):
        self.bot = bot
        self.start_time = time.time()
    
    @app_commands.command(name="ping", description="Check bot latency")
    async def ping(self, interaction: discord.Interaction):
        """Check bot latency."""
        latency = round(self.bot.latency * 1000)
        
        embed = discord.Embed(
            title="🏓 Pong!",
            color=discord.Color.green()
        )
        embed.add_field(name="Latency", value=f"{latency}ms", inline=True)
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="info", description="Bot information")
    async def info(self, interaction: discord.Interaction):
        """Show bot information."""
        uptime = time.time() - self.start_time
        hours, remainder = divmod(int(uptime), 3600)
        minutes, seconds = divmod(remainder, 60)
        
        embed = discord.Embed(
            title="🤖 NeruBot Information",
            color=discord.Color.blue()
        )
        embed.add_field(name="Servers", value=len(self.bot.guilds), inline=True)
        embed.add_field(name="Users", value=len(self.bot.users), inline=True)
        embed.add_field(name="Uptime", value=f"{hours}h {minutes}m {seconds}s", inline=True)
        embed.add_field(name="Python Version", value=f"{sys.version_info.major}.{sys.version_info.minor}.{sys.version_info.micro}", inline=True)
        embed.add_field(name="Discord.py Version", value=discord.__version__, inline=True)
        
        try:
            memory = psutil.virtual_memory()
            embed.add_field(name="Memory Usage", value=f"{memory.percent}%", inline=True)
        except:
            pass
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="help", description="Show bot commands")
    async def help(self, interaction: discord.Interaction):
        """Show help information."""
        embed = discord.Embed(
            title="🆘 NeruBot Help",
            description="Here are the available commands:",
            color=discord.Color.gold()
        )
        
        # General commands
        embed.add_field(
            name="📋 General Commands",
            value="`/ping` - Check bot latency\n"
                  "`/info` - Bot information\n"
                  "`/help` - Show this help",
            inline=False
        )
        
        # Music commands
        embed.add_field(
            name="🎵 Music Commands",
            value="`/play <song>` - Play a song\n"
                  "`/stop` - Stop music and clear queue\n"
                  "`/skip` - Skip current song\n"
                  "`/queue` - Show current queue\n"
                  "`/pause` - Pause playback\n"
                  "`/resume` - Resume playback\n"
                  "`/join` - Join voice channel\n"
                  "`/leave` - Leave voice channel",
            inline=False
        )
        
        # Fun commands
        embed.add_field(
            name="🎲 Fun Commands",
            value="`/roll <sides>` - Roll a dice\n"
                  "`/coinflip` - Flip a coin\n"
                  "`/joke` - Get a random joke",
            inline=False
        )
        
        embed.set_footer(text="Use / to access slash commands!")
        await interaction.response.send_message(embed=embed)

async def setup(bot):
    await bot.add_cog(General(bot))
