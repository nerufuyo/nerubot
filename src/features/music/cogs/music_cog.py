"""
Music commands cog for the Discord bot.
"""
import asyncio
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional
from src.features.music.services.music_service import MusicService, LoopMode
from src.features.music.services.sources import MusicSource
from src.core.utils.logging_utils import get_logger
from src.core.constants import (
    # Emojis and symbols
    SOURCE_EMOJI, LOOP_EMOJI,
    # Colors
    COLOR_SUCCESS, COLOR_ERROR, COLOR_INFO, COLOR_WARNING,
    # Messages
    MSG_ERROR, MSG_SUCCESS, MSG_INFO, MSG_HELP,
    # Command descriptions
    CMD_DESCRIPTIONS,
    # Timeouts
    TIMEOUT_PLAY_COMMAND
)

# Configure logger
logger = get_logger(__name__)

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
                            description=MSG_INFO['goodbye_disconnect'],
                            color=COLOR_WARNING
                        )
                        await channel.send(embed=embed)
                    except:
                        pass  # Ignore if we can't send the message
            
            # Clean up
            await self.music_service.clear_queue(guild_id)
            await self.music_service.cancel_idle_timer(guild_id)
    
    @app_commands.command(name="join", description=CMD_DESCRIPTIONS['join'])
    async def join(self, interaction: discord.Interaction):
        """Join the user's voice channel."""
        if not interaction.user.voice:
            await interaction.response.send_message(MSG_ERROR['not_in_voice'])
            return
        
        channel = interaction.user.voice.channel
        
        if interaction.guild.voice_client:
            if interaction.guild.voice_client.channel == channel:
                await interaction.response.send_message(MSG_INFO['already_connected'])
                return
            else:
                await interaction.guild.voice_client.move_to(channel)
        else:
            await channel.connect()
        
        embed = discord.Embed(
            title="🔊 Joined Voice Channel",
            description=f"Connected to **{channel.name}**",
            color=COLOR_SUCCESS
        )
        await interaction.response.send_message(embed=embed)
        
        # Set disconnect message channel
        self.music_service.set_disconnect_channel(interaction.guild.id, interaction.channel)
    
    @app_commands.command(name="leave", description=CMD_DESCRIPTIONS['leave'])
    async def leave(self, interaction: discord.Interaction):
        """Leave the voice channel."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message(MSG_ERROR['not_connected'])
            return
        
        channel_name = interaction.guild.voice_client.channel.name
        await interaction.guild.voice_client.disconnect()
        
        # Clear the queue
        await self.music_service.clear_queue(interaction.guild.id)
        
        embed = discord.Embed(
            title="👋 Left Voice Channel",
            description=f"Disconnected from **{channel_name}**",
            color=COLOR_WARNING
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="play", description=CMD_DESCRIPTIONS['play'])
    @app_commands.describe(query="Song name or URL (YouTube, Spotify, SoundCloud)")
    async def play(self, interaction: discord.Interaction, query: str):
        """Play a song from various sources."""
        # Send initial response to prevent timeout
        await interaction.response.defer(thinking=True)
        
        # Check if user is in voice channel
        if not interaction.user.voice:
            await interaction.followup.send(MSG_ERROR['not_in_voice'])
            return
        
        # Join voice channel if not connected
        if not interaction.guild.voice_client:
            try:
                await interaction.user.voice.channel.connect()
                # Set disconnect message channel
                self.music_service.set_disconnect_channel(interaction.guild.id, interaction.channel)
            except Exception as e:
                await interaction.followup.send(f"{MSG_ERROR['join_failed']}: {str(e)}")
                return
        
        try:
            # Add to queue and play with timeout handling
            try:
                result = await asyncio.wait_for(
                    self.music_service.add_to_queue(interaction.guild.id, query, interaction.user),
                    timeout=TIMEOUT_PLAY_COMMAND
                )
            except asyncio.TimeoutError:
                await interaction.followup.send(MSG_ERROR['play_timeout'])
                return
            
            if result['success']:
                # Get source emoji
                source = result.get('source', 'unknown')
                source_emoji = SOURCE_EMOJI.get(source, '🎵')
                
                embed = discord.Embed(
                    title=f"{source_emoji} Added to Queue" if result['queued'] else f"{source_emoji} Now Playing",
                    description=f"**{result['title']}**",
                    color=COLOR_INFO
                )
                embed.add_field(name="Duration", value=result.get('duration', 'Unknown'), inline=True)
                embed.add_field(name="Requested by", value=interaction.user.mention, inline=True)
                embed.add_field(name="Source", value=f"{source_emoji} {source.capitalize()}", inline=True)
                
                if result['queued']:
                    embed.add_field(name="Position in queue", value=result['position'], inline=True)
                
                await interaction.followup.send(embed=embed)
            else:
                # Send more helpful error messages
                error_msg = result.get('error', 'Unknown error occurred')
                if "Could not find any information" in error_msg:
                    error_msg += f"\n\n{MSG_HELP['search_tips']}"
                elif "timed out" in error_msg.lower():
                    error_msg += f"\n\n{MSG_HELP['timeout_tips']}"
                
                await interaction.followup.send(f"❌ {error_msg}")
                
        except Exception as e:
            logger.error(f"Error in play command: {e}")
            await interaction.followup.send(MSG_ERROR['unexpected'])
    
    @app_commands.command(name="stop", description=CMD_DESCRIPTIONS['stop'])
    async def stop(self, interaction: discord.Interaction):
        """Stop music and clear the queue."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message(MSG_ERROR['not_connected'])
            return
        
        if not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message(MSG_ERROR['nothing_playing'])
            return
        
        interaction.guild.voice_client.stop()
        await self.music_service.clear_queue(interaction.guild.id)
        
        embed = discord.Embed(
            title="⏹️ Stopped",
            description="Music stopped and queue cleared",
            color=COLOR_ERROR
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="skip", description=CMD_DESCRIPTIONS['skip'])
    async def skip(self, interaction: discord.Interaction):
        """Skip the current song."""
        if not interaction.guild.voice_client or not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message(MSG_ERROR['nothing_playing'])
            return
        
        # Get current queue to check if there are more songs
        queue = await self.music_service.get_queue(interaction.guild.id)
        has_next_song = len(queue) > 1
        
        # Stop the current song (this will trigger _after_play)
        interaction.guild.voice_client.stop()
        
        if has_next_song:
            embed = discord.Embed(
                title="⏭️ Skipped",
                description="Skipped to the next song",
                color=COLOR_INFO
            )
        else:
            embed = discord.Embed(
                title="⏭️ Skipped",
                description="Skipped the last song. Queue is now empty.",
                color=COLOR_WARNING
            )
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="pause", description=CMD_DESCRIPTIONS['pause'])
    async def pause(self, interaction: discord.Interaction):
        """Pause the current song."""
        if not interaction.guild.voice_client or not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message(MSG_ERROR['nothing_playing'])
            return
        
        interaction.guild.voice_client.pause()
        
        embed = discord.Embed(
            title="⏸️ Paused",
            description="Music paused",
            color=COLOR_WARNING
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="resume", description=CMD_DESCRIPTIONS['resume'])
    async def resume(self, interaction: discord.Interaction):
        """Resume the current song."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message(MSG_ERROR['not_connected'])
            return
        
        if not interaction.guild.voice_client.is_paused():
            await interaction.response.send_message(MSG_ERROR['not_paused'])
            return
        
        interaction.guild.voice_client.resume()
        
        embed = discord.Embed(
            title="▶️ Resumed",
            description="Music resumed",
            color=COLOR_SUCCESS
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="queue", description=CMD_DESCRIPTIONS['queue'])
    async def queue(self, interaction: discord.Interaction):
        """Show the current music queue."""
        queue = await self.music_service.get_queue(interaction.guild.id)
        
        if not queue:
            await interaction.response.send_message(MSG_INFO['queue_empty'])
            return
        
        embed = discord.Embed(
            title="🎵 Music Queue",
            color=COLOR_INFO
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
    
    @app_commands.command(name="loop", description=CMD_DESCRIPTIONS['loop'])
    @app_commands.describe(mode="Loop mode: off, single, or queue")
    async def loop(self, interaction: discord.Interaction, mode: str = None):
        """Toggle loop mode."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message(MSG_ERROR['not_connected'])
            return
        
        current_mode = self.music_service.get_loop_mode(interaction.guild.id)
        
        if mode is None:
            # Cycle through modes
            if current_mode == LoopMode.OFF:
                new_mode = LoopMode.SINGLE
                mode_text = "Single Song"
                emoji = LOOP_EMOJI['single']
            elif current_mode == LoopMode.SINGLE:
                new_mode = LoopMode.QUEUE
                mode_text = "Queue"
                emoji = LOOP_EMOJI['queue']
            else:
                new_mode = LoopMode.OFF
                mode_text = "Off"
                emoji = LOOP_EMOJI['off']
        else:
            # Set specific mode
            mode = mode.lower()
            if mode in ["off", "none", "0"]:
                new_mode = LoopMode.OFF
                mode_text = "Off"
                emoji = LOOP_EMOJI['off']
            elif mode in ["single", "song", "1"]:
                new_mode = LoopMode.SINGLE
                mode_text = "Single Song"
                emoji = LOOP_EMOJI['single']
            elif mode in ["queue", "all", "2"]:
                new_mode = LoopMode.QUEUE
                mode_text = "Queue"
                emoji = LOOP_EMOJI['queue']
            else:
                await interaction.response.send_message(MSG_ERROR['invalid_loop_mode'])
                return
        
        await self.music_service.set_loop_mode(interaction.guild.id, new_mode)
        
        embed = discord.Embed(
            title=f"{emoji} Loop Mode",
            description=f"Loop mode set to: **{mode_text}**",
            color=COLOR_INFO
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="247", description=CMD_DESCRIPTIONS['247'])
    async def twenty_four_seven(self, interaction: discord.Interaction):
        """Toggle 24/7 mode."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message(MSG_ERROR['not_connected'])
            return
        
        is_enabled = await self.music_service.toggle_24_7(interaction.guild.id)
        
        if is_enabled:
            embed = discord.Embed(
                title="🌙 24/7 Mode Enabled",
                description=MSG_SUCCESS['247_enabled'],
                color=COLOR_INFO
            )
        else:
            embed = discord.Embed(
                title="☀️ 24/7 Mode Disabled",
                description=MSG_SUCCESS['247_disabled'],
                color=COLOR_WARNING
            )
        
        await interaction.response.send_message(embed=embed)
        
        # Set disconnect message channel
        self.music_service.set_disconnect_channel(interaction.guild.id, interaction.channel)
    
    @app_commands.command(name="nowplaying", description=CMD_DESCRIPTIONS['nowplaying'])
    async def nowplaying(self, interaction: discord.Interaction):
        """Show the currently playing song."""
        if not interaction.guild.voice_client or not interaction.guild.voice_client.is_playing():
            await interaction.response.send_message(MSG_ERROR['nothing_playing'])
            return
        
        queue = await self.music_service.get_queue(interaction.guild.id)
        if not queue:
            await interaction.response.send_message(MSG_ERROR['nothing_playing'])
            return
        
        current_song = queue[0]
        loop_mode = self.music_service.get_loop_mode(interaction.guild.id)
        is_24_7 = self.music_service.is_24_7_enabled(interaction.guild.id)
        
        # Determine loop emoji
        if loop_mode == LoopMode.SINGLE:
            loop_emoji = LOOP_EMOJI['single']
            loop_text = "Single Song"
        elif loop_mode == LoopMode.QUEUE:
            loop_emoji = LOOP_EMOJI['queue']
            loop_text = "Queue"
        else:
            loop_emoji = LOOP_EMOJI['off']
            loop_text = "Off"
        
        # Get source emoji
        source = current_song.get('source', 'unknown')
        source_emoji = SOURCE_EMOJI.get(source, '🎵')
        
        embed = discord.Embed(
            title=f"{source_emoji} Now Playing",
            description=f"**{current_song['title']}**",
            color=COLOR_INFO
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
    
    @app_commands.command(name="clear", description=CMD_DESCRIPTIONS['clear'])
    async def clear(self, interaction: discord.Interaction):
        """Clear the music queue."""
        if not interaction.guild.voice_client:
            await interaction.response.send_message(MSG_ERROR['not_connected'])
            return
        
        await self.music_service.clear_queue(interaction.guild.id)
        
        embed = discord.Embed(
            title="🗑️ Queue Cleared",
            description=MSG_SUCCESS['queue_cleared'],
            color=COLOR_ERROR
        )
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="sources", description=CMD_DESCRIPTIONS['sources'])
    async def sources(self, interaction: discord.Interaction):
        """Show information about available music sources."""
        embed = discord.Embed(
            title="🎵 Available Music Sources",
            description="You can play music from these sources:",
            color=COLOR_INFO
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['youtube']} YouTube",
            value=MSG_HELP['sources']['youtube'],
            inline=False
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['spotify']} Spotify",
            value=MSG_HELP['sources']['spotify'],
            inline=False
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['soundcloud']} SoundCloud",
            value=MSG_HELP['sources']['soundcloud'],
            inline=False
        )
        
        embed.add_field(
            name=f"{SOURCE_EMOJI['direct']} Direct Links",
            value=MSG_HELP['sources']['direct'],
            inline=False
        )
        
        embed.set_footer(text=MSG_HELP['sources']['footer'])
        
        await interaction.response.send_message(embed=embed)
