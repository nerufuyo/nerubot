"""
Music commands cog for the Discord bot.
"""
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional
from src.features.music.services.music_service import MusicService
from src.core.utils.logging_utils import get_logger
from src.core.utils.messages import (
    VOICE_JOINED, VOICE_NOT_CONNECTED, VOICE_JOIN_FAILED, VOICE_DISCONNECTED, 
    VOICE_LEAVE_FAILED, USER_NOT_IN_CHANNEL, PLAYBACK_ERROR,
    SONG_SKIPPED_NEXT, SONG_SKIPPED_NO_MORE, NOTHING_PLAYING, NOTHING_PAUSED,
    PAUSED, RESUMED, QUEUE_CLEARED
)

# Configure logger
logger = get_logger(__name__)

class MusicCog(commands.Cog):
    """Music playback commands."""
    
    def __init__(self, bot):
        self.bot = bot
        self.music_service = MusicService(bot)
    
    @app_commands.command(name="join", description="Join your voice channel")
    async def join(self, interaction: discord.Interaction):
        """Join the user's voice channel."""
        if not interaction.user.voice:
            await interaction.response.send_message("❌ You need to be in a voice channel!")
            return
        
        channel = interaction.user.voice.channel
        
        if interaction.guild.voice_client:
            if interaction.guild.voice_client.channel == channel:
                await interaction.response.send_message("✅ Already connected to your voice channel!")
                return
            else:
                await interaction.guild.voice_client.move_to(channel)
        else:
            await channel.connect()
        
        embed = discord.Embed(
            title="🔊 Joined Voice Channel",
            description=f"Connected to **{channel.name}**",
            color=discord.Color.green()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="leave", description="Leave the voice channel")
    async def leave(self, interaction: discord.Interaction):
        """Leave the voice channel."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message("❌ I'm not connected to a voice channel!")
            return
        
        channel_name = interaction.guild.voice_client.channel.name
        await interaction.guild.voice_client.disconnect()
        
        # Clear the queue
        await self.music_service.clear_queue(interaction.guild.id)
        
        embed = discord.Embed(
            title="👋 Left Voice Channel",
            description=f"Disconnected from **{channel_name}**",
            color=discord.Color.orange()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="play", description="Play a song")
    @app_commands.describe(query="Song name or YouTube URL")
    async def play(self, interaction: discord.Interaction, query: str):
        """Play a song from YouTube."""
        await interaction.response.defer()
        
        # Check if user is in voice channel
        if not interaction.user.voice:
            await interaction.followup.send("❌ You need to be in a voice channel!")
            return
        
        # Join voice channel if not connected
        if not interaction.guild.voice_client:
            await interaction.user.voice.channel.connect()
        
        try:
            # Add to queue and play
            result = await self.music_service.add_to_queue(interaction.guild.id, query, interaction.user)
            
            if result['success']:
                embed = discord.Embed(
                    title="🎵 Added to Queue" if result['queued'] else "🎵 Now Playing",
                    description=f"**{result['title']}**",
                    color=discord.Color.blue()
                )
                embed.add_field(name="Duration", value=result.get('duration', 'Unknown'), inline=True)
                embed.add_field(name="Requested by", value=interaction.user.mention, inline=True)
                
                if result['queued']:
                    embed.add_field(name="Position in queue", value=result['position'], inline=True)
                
                await interaction.followup.send(embed=embed)
            else:
                await interaction.followup.send(f"❌ {result['error']}")
                
        except Exception as e:
            await interaction.followup.send(f"❌ An error occurred: {str(e)}")
    
    @app_commands.command(name="stop", description="Stop music and clear queue")
    async def stop(self, interaction: discord.Interaction):
        """Stop music and clear the queue."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message("❌ I'm not connected to a voice channel!")
            return
        
        if not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message("❌ Nothing is currently playing!")
            return
        
        interaction.guild.voice_client.stop()
        await self.music_service.clear_queue(interaction.guild.id)
        
        embed = discord.Embed(
            title="⏹️ Stopped",
            description="Music stopped and queue cleared",
            color=discord.Color.red()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="skip", description="Skip the current song")
    async def skip(self, interaction: discord.Interaction):
        """Skip the current song."""
        if not interaction.guild.voice_client or not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message("❌ Nothing is currently playing!")
            return
        
        interaction.guild.voice_client.stop()  # This will trigger the next song
        
        embed = discord.Embed(
            title="⏭️ Skipped",
            description="Skipped to the next song",
            color=discord.Color.blue()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="pause", description="Pause the current song")
    async def pause(self, interaction: discord.Interaction):
        """Pause the current song."""
        if not interaction.guild.voice_client or not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message("❌ Nothing is currently playing!")
            return
        
        interaction.guild.voice_client.pause()
        
        embed = discord.Embed(
            title="⏸️ Paused",
            description="Music paused",
            color=discord.Color.yellow()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="resume", description="Resume the current song")
    async def resume(self, interaction: discord.Interaction):
        """Resume the current song."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message("❌ I'm not connected to a voice channel!")
            return
        
        if not interaction.guild.voice_client.is_paused():
            await interaction.response.send_message("❌ Music is not paused!")
            return
        
        interaction.guild.voice_client.resume()
        
        embed = discord.Embed(
            title="▶️ Resumed",
            description="Music resumed",
            color=discord.Color.green()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="queue", description="Show the music queue")
    async def queue(self, interaction: discord.Interaction):
        """Show the current music queue."""
        queue = await self.music_service.get_queue(interaction.guild.id)
        
        if not queue:
            await interaction.response.send_message("❌ The queue is empty!")
            return
        
        embed = discord.Embed(
            title="🎵 Music Queue",
            color=discord.Color.blue()
        )
        
        # Show currently playing
        if interaction.guild.voice_client and interaction.guild.voice_client.is_playing():
            current = queue[0] if queue else None
            if current:
                embed.add_field(
                    name="🎵 Now Playing",
                    value=f"**{current['title']}**\nRequested by {current['requester'].mention}",
                    inline=False
                )
        
        # Show next songs
        if len(queue) > 1:
            next_songs = []
            for i, song in enumerate(queue[1:6], 1):  # Show next 5 songs
                next_songs.append(f"{i}. **{song['title']}** - {song['requester'].mention}")
            
            if next_songs:
                embed.add_field(
                    name="⏭️ Up Next",
                    value="\n".join(next_songs),
                    inline=False
                )
            
            if len(queue) > 6:
                embed.add_field(
                    name="📊 Queue Info",
                    value=f"And {len(queue) - 6} more songs...",
                    inline=False
                )
        
        await interaction.response.send_message(embed=embed)
