"""
Music commands cog for the Discord bot.
"""
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional
from src.features.music.services.music_service import MusicService, LoopMode
from src.features.music.services.sources import MusicSource
from src.core.utils.logging_utils import get_logger

# Configure logger
logger = get_logger(__name__)

# Source emoji mapping
SOURCE_EMOJI = {
    'youtube': '▶️',
    'spotify': '💚',
    'soundcloud': '🧡',
    'direct': '🔗',
    'unknown': '🎵'
}

class MusicCog(commands.Cog):
    """Music playback commands."""
    
    def __init__(self, bot):
        self.bot = bot
        self.music_service = MusicService(bot)
    
    @commands.Cog.listener()
    async def on_voice_state_update(self, member, before, after):
        """Handle voice state updates."""
        # Check if the bot was disconnected
        if member == self.bot.user and before.channel and not after.channel:
            guild_id = before.channel.guild.id
            
            # Send goodbye message
            if guild_id in self.music_service.disconnect_messages:
                channel = self.music_service.disconnect_messages[guild_id]
                if channel:
                    try:
                        embed = discord.Embed(
                            title="👋 Goodbye!",
                            description="I've been disconnected from the voice channel. Thanks for listening!",
                            color=discord.Color.orange()
                        )
                        await channel.send(embed=embed)
                    except:
                        pass  # Ignore if we can't send the message
            
            # Clean up
            await self.music_service.clear_queue(guild_id)
            await self.music_service.cancel_idle_timer(guild_id)
    
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
        
        # Set disconnect message channel
        self.music_service.set_disconnect_channel(interaction.guild.id, interaction.channel)
    
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
    @app_commands.describe(query="Song name or URL (YouTube, Spotify, SoundCloud)")
    async def play(self, interaction: discord.Interaction, query: str):
        """Play a song from various sources."""
        await interaction.response.defer()
        
        # Check if user is in voice channel
        if not interaction.user.voice:
            await interaction.followup.send("❌ You need to be in a voice channel!")
            return
        
        # Join voice channel if not connected
        if not interaction.guild.voice_client:
            await interaction.user.voice.channel.connect()
            # Set disconnect message channel
            self.music_service.set_disconnect_channel(interaction.guild.id, interaction.channel)
        
        try:
            # Add to queue and play
            result = await self.music_service.add_to_queue(interaction.guild.id, query, interaction.user)
            
            if result['success']:
                # Get source emoji
                source = result.get('source', 'unknown')
                source_emoji = SOURCE_EMOJI.get(source, '🎵')
                
                embed = discord.Embed(
                    title=f"{source_emoji} Added to Queue" if result['queued'] else f"{source_emoji} Now Playing",
                    description=f"**{result['title']}**",
                    color=discord.Color.blue()
                )
                embed.add_field(name="Duration", value=result.get('duration', 'Unknown'), inline=True)
                embed.add_field(name="Requested by", value=interaction.user.mention, inline=True)
                embed.add_field(name="Source", value=f"{source_emoji} {source.capitalize()}", inline=True)
                
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
                # Get source emoji
                source = current.get('source', 'unknown')
                source_emoji = SOURCE_EMOJI.get(source, '🎵')
                
                embed.add_field(
                    name=f"{source_emoji} Now Playing",
                    value=f"**{current['title']}**\nRequested by {current['requester'].mention}",
                    inline=False
                )
        
        # Show next songs
        if len(queue) > 1:
            next_songs = []
            for i, song in enumerate(queue[1:6], 1):  # Show next 5 songs
                # Get source emoji
                source = song.get('source', 'unknown')
                source_emoji = SOURCE_EMOJI.get(source, '🎵')
                next_songs.append(f"{i}. {source_emoji} **{song['title']}** - {song['requester'].mention}")
            
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
    
    @app_commands.command(name="loop", description="Toggle loop mode (off/single/queue)")
    @app_commands.describe(mode="Loop mode: off, single, or queue")
    async def loop(self, interaction: discord.Interaction, mode: str = None):
        """Toggle loop mode."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message("❌ I'm not connected to a voice channel!")
            return
        
        current_mode = self.music_service.get_loop_mode(interaction.guild.id)
        
        if mode is None:
            # Cycle through modes
            if current_mode == LoopMode.OFF:
                new_mode = LoopMode.SINGLE
                mode_text = "Single Song"
                emoji = "🔂"
            elif current_mode == LoopMode.SINGLE:
                new_mode = LoopMode.QUEUE
                mode_text = "Queue"
                emoji = "🔁"
            else:
                new_mode = LoopMode.OFF
                mode_text = "Off"
                emoji = "▶️"
        else:
            # Set specific mode
            mode = mode.lower()
            if mode in ["off", "none", "0"]:
                new_mode = LoopMode.OFF
                mode_text = "Off"
                emoji = "▶️"
            elif mode in ["single", "song", "1"]:
                new_mode = LoopMode.SINGLE
                mode_text = "Single Song"
                emoji = "🔂"
            elif mode in ["queue", "all", "2"]:
                new_mode = LoopMode.QUEUE
                mode_text = "Queue"
                emoji = "🔁"
            else:
                await interaction.response.send_message("❌ Invalid loop mode! Use: `off`, `single`, or `queue`")
                return
        
        await self.music_service.set_loop_mode(interaction.guild.id, new_mode)
        
        embed = discord.Embed(
            title=f"{emoji} Loop Mode",
            description=f"Loop mode set to: **{mode_text}**",
            color=discord.Color.blue()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="247", description="Toggle 24/7 mode")
    async def twenty_four_seven(self, interaction: discord.Interaction):
        """Toggle 24/7 mode."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message("❌ I'm not connected to a voice channel!")
            return
        
        is_enabled = await self.music_service.toggle_24_7(interaction.guild.id)
        
        if is_enabled:
            embed = discord.Embed(
                title="🌙 24/7 Mode Enabled",
                description="I will stay in the voice channel even when not playing music.",
                color=discord.Color.purple()
            )
        else:
            embed = discord.Embed(
                title="☀️ 24/7 Mode Disabled",
                description="I will leave the voice channel after 5 minutes of inactivity.",
                color=discord.Color.orange()
            )
        
        await interaction.response.send_message(embed=embed)
        
        # Set disconnect message channel
        self.music_service.set_disconnect_channel(interaction.guild.id, interaction.channel)
    
    @app_commands.command(name="nowplaying", description="Show currently playing song")
    async def nowplaying(self, interaction: discord.Interaction):
        """Show the currently playing song."""
        if not interaction.guild.voice_client or not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message("❌ Nothing is currently playing!")
            return
        
        queue = await self.music_service.get_queue(interaction.guild.id)
        if not queue:
            await interaction.response.send_message("❌ Nothing is currently playing!")
            return
        
        current_song = queue[0]
        loop_mode = self.music_service.get_loop_mode(interaction.guild.id)
        is_24_7 = self.music_service.is_24_7_enabled(interaction.guild.id)
        
        # Determine loop emoji
        if loop_mode == LoopMode.SINGLE:
            loop_emoji = "🔂"
            loop_text = "Single Song"
        elif loop_mode == LoopMode.QUEUE:
            loop_emoji = "🔁"
            loop_text = "Queue"
        else:
            loop_emoji = "▶️"
            loop_text = "Off"
        
        # Get source emoji
        source = current_song.get('source', 'unknown')
        source_emoji = SOURCE_EMOJI.get(source, '🎵')
        
        embed = discord.Embed(
            title=f"{source_emoji} Now Playing",
            description=f"**{current_song['title']}**",
            color=discord.Color.blue()
        )
        embed.add_field(name="Duration", value=current_song['duration'], inline=True)
        embed.add_field(name="Requested by", value=current_song['requester'].mention, inline=True)
        embed.add_field(name="Source", value=f"{source_emoji} {source.capitalize()}", inline=True)
        embed.add_field(name="Loop Mode", value=f"{loop_emoji} {loop_text}", inline=True)
        
        if is_24_7:
            embed.add_field(name="24/7 Mode", value="🌙 Enabled", inline=True)
        
        # Add artist and album if available
        if current_song.get('artist'):
            embed.add_field(name="Artist", value=current_song['artist'], inline=True)
        
        if current_song.get('album'):
            embed.add_field(name="Album", value=current_song['album'], inline=True)
        
        # Add thumbnail if available
        if current_song.get('thumbnail'):
            embed.set_thumbnail(url=current_song['thumbnail'])
        
        if len(queue) > 1:
            # Get source emoji for next song
            next_source = queue[1].get('source', 'unknown')
            next_source_emoji = SOURCE_EMOJI.get(next_source, '🎵')
            embed.add_field(name="Up Next", value=f"{next_source_emoji} **{queue[1]['title']}**", inline=False)
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="clear", description="Clear the music queue")
    async def clear(self, interaction: discord.Interaction):
        """Clear the music queue."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message("❌ I'm not connected to a voice channel!")
            return
        
        await self.music_service.clear_queue(interaction.guild.id)
        
        embed = discord.Embed(
            title="🗑️ Queue Cleared",
            description="The music queue has been cleared.",
            color=discord.Color.red()
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="sources", description="Show available music sources")
    async def sources(self, interaction: discord.Interaction):
        """Show information about available music sources."""
        embed = discord.Embed(
            title="🎵 Available Music Sources",
            description="You can play music from these sources:",
            color=discord.Color.blue()
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['youtube']} YouTube",
            value="Play music from YouTube links or search for songs\nExample: `/play despacito` or `/play https://www.youtube.com/watch?v=kJQP7kiw5Fk`",
            inline=False
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['spotify']} Spotify",
            value="Play music from Spotify links (tracks, albums, playlists, artists)\nExample: `/play https://open.spotify.com/track/6habFhsOp2NvshLv26jCFK`",
            inline=False
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['soundcloud']} SoundCloud",
            value="Play music from SoundCloud links (tracks, playlists, users)\nExample: `/play https://soundcloud.com/artist/track`",
            inline=False
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['direct']} Direct Links",
            value="Play music directly from audio file links\nExample: `/play https://example.com/music.mp3`",
            inline=False
        )
        
        embed.set_footer(text="Use /play followed by a search term or URL to play music from any of these sources!")
        
        await interaction.response.send_message(embed=embed)
