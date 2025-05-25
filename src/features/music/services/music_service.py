"""
Music service for handling audio playback in Discord voice channels.
"""
import asyncio
import os
import discord
from collections import deque, defaultdict
from typing import Dict, Optional, List
from enum import Enum
from src.core.utils.logging_utils import get_logger
from src.core.utils.file_utils import find_ffmpeg
from src.features.music.services.sources import MusicSource, MusicSourceResult
from src.features.music.services.sources.manager import SourceManager

# Configure logger
logger = get_logger(__name__)

# Load Opus library for Discord voice support
if not discord.opus.is_loaded():
    opus_paths = [
        '/opt/homebrew/lib/libopus.dylib',  # Homebrew on Apple Silicon
        '/usr/local/lib/libopus.dylib',     # Homebrew on Intel
        '/opt/homebrew/lib/libopus.0.dylib',
        '/usr/local/lib/libopus.0.dylib',
        'libopus.dylib',  # System path
        'opus'  # Let the system find it
    ]
    
    for opus_path in opus_paths:
        try:
            discord.opus.load_opus(opus_path)
            logger.info(f"Successfully loaded Opus from: {opus_path}")
            break
        except:
            continue
    
    if not discord.opus.is_loaded():
        logger.warning("Failed to load Opus library. Voice functionality may not work.")

# Get FFmpeg path
FFMPEG_PATH = find_ffmpeg()
logger.info(f"Using FFmpeg path: {FFMPEG_PATH or 'system default (in PATH)'}")

# FFmpeg options
ffmpeg_options = {
    'before_options': '-reconnect 1 -reconnect_streamed 1 -reconnect_delay_max 5',
    'options': '-vn'
}

ffmpeg_executable = FFMPEG_PATH if FFMPEG_PATH else 'ffmpeg'

class LoopMode(Enum):
    """Loop modes for music playback."""
    OFF = 0
    SINGLE = 1
    QUEUE = 2

class MusicService:
    """Service for handling music playback in Discord voice channels."""
    
    def __init__(self, bot):
        self.bot = bot
        self.queues = defaultdict(deque)
        self.current_songs = {}
        self.loop_modes = defaultdict(lambda: LoopMode.OFF)
        self.is_24_7 = defaultdict(bool)
        self.idle_timers = {}
        self.disconnect_messages = {}
        self.source_manager = SourceManager()
    
    async def add_to_queue(self, guild_id: int, query: str, requester):
        """Add a song to the queue."""
        try:
            # Search for the song using the source manager
            results = await self.source_manager.search(query)
            
            if not results:
                return {
                    'success': False,
                    'error': "Could not find any songs matching your query"
                }
            
            # Get the first result
            result = results[0]
            
            # Get a playable version of the result
            playable_result = await self.source_manager.get_playable(result)
            
            if not playable_result or not playable_result.url:
                return {
                    'success': False,
                    'error': "No playable URL found for that song"
                }
            
            # Convert to song info format
            song_info = {
                'title': playable_result.title,
                'url': playable_result.url,
                'webpage_url': playable_result.webpage_url,
                'duration': playable_result.duration,
                'requester': requester,
                'source': playable_result.source.value,
                'thumbnail': playable_result.thumbnail,
                'artist': playable_result.artist,
                'album': playable_result.album
            }
            
            # Add to queue
            self.queues[guild_id].append(song_info)
            
            # If nothing is playing, start playing
            guild = self.bot.get_guild(guild_id)
            if guild and guild.voice_client and not guild.voice_client.is_playing():
                await self._play_next(guild_id)
                return {
                    'success': True,
                    'queued': False,
                    'title': song_info['title'],
                    'duration': song_info['duration'],
                    'source': song_info['source']
                }
            
            return {
                'success': True,
                'queued': True,
                'title': song_info['title'],
                'duration': song_info['duration'],
                'position': len(self.queues[guild_id]),
                'source': song_info['source']
            }
            
        except Exception as e:
            logger.error(f"Error adding song to queue: {e}")
            return {
                'success': False,
                'error': f"Could not find or play that song: {str(e)}"
            }
    
    async def _play_next(self, guild_id: int):
        """Play the next song in the queue."""
        if not self.queues[guild_id]:
            return
        
        guild = self.bot.get_guild(guild_id)
        if not guild or not guild.voice_client:
            return
        
        # Cancel idle timer since we're about to play
        await self.cancel_idle_timer(guild_id)
        
        song = self.queues[guild_id][0]  # Keep the current song at front for queue display
        self.current_songs[guild_id] = song
        
        try:
            # Create audio source
            source = discord.FFmpegPCMAudio(
                song['url'], 
                executable=ffmpeg_executable,
                before_options=ffmpeg_options['before_options'],
                options=ffmpeg_options['options']
            )
            
            # Play the song
            guild.voice_client.play(
                source,
                after=lambda e: self._after_play(guild_id, e)
            )
            
            logger.info(f"Now playing: {song['title']} in guild {guild_id} from {song.get('source', 'unknown')}")
            
        except Exception as e:
            logger.error(f"Error playing song: {e}")
            # Remove the failed song and try next
            if self.queues[guild_id]:
                self.queues[guild_id].popleft()
            await self._play_next(guild_id)
    
    def _after_play(self, guild_id: int, error):
        """Called after a song finishes playing."""
        if error:
            logger.error(f"Player error: {error}")
        
        loop_mode = self.loop_modes[guild_id]
        
        # Handle looping
        if loop_mode == LoopMode.SINGLE and self.queues[guild_id]:
            # Single song loop - don't remove the song
            pass
        elif loop_mode == LoopMode.QUEUE and self.queues[guild_id]:
            # Queue loop - move current song to end
            current_song = self.queues[guild_id].popleft()
            self.queues[guild_id].append(current_song)
        else:
            # No loop - remove the finished song
            if self.queues[guild_id]:
                self.queues[guild_id].popleft()
        
        # Remove from current songs if not looping single
        if loop_mode != LoopMode.SINGLE and guild_id in self.current_songs:
            del self.current_songs[guild_id]
        
        # Play next song or start idle timer
        if self.queues[guild_id]:
            coro = self._play_next(guild_id)
        else:
            # No more songs - start idle timer
            coro = self.start_idle_timer(guild_id)
        
        future = asyncio.run_coroutine_threadsafe(coro, self.bot.loop)
        try:
            future.result()
        except Exception as e:
            logger.error(f"Error in after_play: {e}")
    
    async def get_queue(self, guild_id: int):
        """Get the current queue for a guild."""
        return list(self.queues[guild_id])
    
    async def clear_queue(self, guild_id: int):
        """Clear the queue for a guild."""
        self.queues[guild_id].clear()
        if guild_id in self.current_songs:
            del self.current_songs[guild_id]
    
    def _format_duration(self, duration):
        """Format duration in seconds to human readable format."""
        if duration is None:
            return "Unknown"
        
        hours, remainder = divmod(int(duration), 3600)
        minutes, seconds = divmod(remainder, 60)
        
        if hours:
            return f"{hours}:{minutes:02d}:{seconds:02d}"
        else:
            return f"{minutes}:{seconds:02d}"
    
    async def set_loop_mode(self, guild_id: int, mode: LoopMode):
        """Set the loop mode for a guild."""
        self.loop_modes[guild_id] = mode
    
    def get_loop_mode(self, guild_id: int) -> LoopMode:
        """Get the current loop mode for a guild."""
        return self.loop_modes[guild_id]
    
    async def toggle_24_7(self, guild_id: int) -> bool:
        """Toggle 24/7 mode for a guild."""
        self.is_24_7[guild_id] = not self.is_24_7[guild_id]
        
        # Cancel idle timer if 24/7 is enabled
        if self.is_24_7[guild_id] and guild_id in self.idle_timers:
            self.idle_timers[guild_id].cancel()
            del self.idle_timers[guild_id]
        
        return self.is_24_7[guild_id]
    
    def is_24_7_enabled(self, guild_id: int) -> bool:
        """Check if 24/7 mode is enabled for a guild."""
        return self.is_24_7[guild_id]
    
    async def start_idle_timer(self, guild_id: int):
        """Start the idle timer for auto-disconnect."""
        if self.is_24_7[guild_id]:
            return  # Don't start timer if 24/7 is enabled
        
        # Cancel existing timer
        if guild_id in self.idle_timers:
            self.idle_timers[guild_id].cancel()
        
        # Start new timer
        self.idle_timers[guild_id] = asyncio.create_task(
            self._idle_disconnect_timer(guild_id)
        )
    
    async def cancel_idle_timer(self, guild_id: int):
        """Cancel the idle timer."""
        if guild_id in self.idle_timers:
            self.idle_timers[guild_id].cancel()
            del self.idle_timers[guild_id]
    
    async def _idle_disconnect_timer(self, guild_id: int):
        """Timer that disconnects the bot after 5 minutes of inactivity."""
        try:
            await asyncio.sleep(300)  # 5 minutes
            
            guild = self.bot.get_guild(guild_id)
            if guild and guild.voice_client:
                # Check if still not playing
                if not guild.voice_client.is_playing() and not self.is_24_7[guild_id]:
                    await self._send_goodbye_message(guild_id)
                    await guild.voice_client.disconnect()
                    await self.clear_queue(guild_id)
                    
        except asyncio.CancelledError:
            pass  # Timer was cancelled
    
    async def _send_goodbye_message(self, guild_id: int):
        """Send goodbye message when disconnecting."""
        if guild_id in self.disconnect_messages:
            channel = self.disconnect_messages[guild_id]
            if channel:
                try:
                    embed = discord.Embed(
                        title="👋 Goodbye!",
                        description="I've been inactive for 5 minutes, so I'm leaving the voice channel. Thanks for listening!",
                        color=discord.Color.orange()
                    )
                    await channel.send(embed=embed)
                except:
                    pass  # Ignore if we can't send the message
    
    def set_disconnect_channel(self, guild_id: int, channel):
        """Set the channel to send disconnect messages to."""
        self.disconnect_messages[guild_id] = channel
