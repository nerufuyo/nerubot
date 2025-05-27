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
from src.core.constants import (
    TIMEOUT_SEARCH, TIMEOUT_CONVERSION, IDLE_DISCONNECT_TIME, MAX_QUEUE_SIZE,
    FFMPEG_OPTIONS, OPUS_PATHS, LOG_MSG, MSG_ERROR, MSG_INFO
)

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
            logger.info(f"Adding to queue in guild {guild_id}: {query}")
            
            # Search for the song using the source manager with timeout
            try:
                results = await asyncio.wait_for(
                    self.source_manager.search(query),
                    timeout=TIMEOUT_SEARCH
                )
            except asyncio.TimeoutError:
                logger.error(f"Search timeout for query: {query}")
                return {
                    'success': False,
                    'error': MSG_ERROR['search_timeout']
                }
            
            if not results:
                logger.warning(f"No search results for: {query}")
                return {
                    'success': False,
                    'error': MSG_ERROR['no_results']
                }
            
            # Get the first result
            result = results[0]
            logger.info(f"Found result: {result.title} from {result.source}")
            
            # Get a playable version of the result with timeout
            try:
                playable_result = await asyncio.wait_for(
                    self.source_manager.get_playable(result),
                    timeout=TIMEOUT_CONVERSION
                )
            except asyncio.TimeoutError:
                logger.error(f"Conversion timeout for: {result.title}")
                return {
                    'success': False,
                    'error': MSG_ERROR['conversion_timeout']
                }
            
            if not playable_result or not playable_result.url:
                logger.error(f"No playable URL found for: {result.title}")
                return {
                    'success': False,
                    'error': MSG_ERROR['no_playable_url']
                }
            
            logger.info(f"Successfully converted to playable: {playable_result.title}")
            
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
            logger.info(f"Added to queue: {song_info['title']} in guild {guild_id}")
            
            # If nothing is playing, start playing
            guild = self.bot.get_guild(guild_id)
            if guild and guild.voice_client and not guild.voice_client.is_playing():
                logger.info(f"Starting playback immediately for guild {guild_id}")
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
            logger.error(f"Error adding song to queue for '{query}': {e}")
            return {
                'success': False,
                'error': MSG_ERROR['add_queue_failed']
            }
    
    async def _play_next(self, guild_id: int):
        """Play the next song in the queue."""
        if not self.queues[guild_id]:
            logger.info(f"No songs in queue for guild {guild_id}")
            return
        
        guild = self.bot.get_guild(guild_id)
        if not guild or not guild.voice_client:
            logger.warning(f"Guild or voice client not found for guild {guild_id}")
            return
        
        # Cancel idle timer since we're about to play
        await self.cancel_idle_timer(guild_id)
        
        song = self.queues[guild_id][0]  # Keep the current song at front for queue display
        self.current_songs[guild_id] = song
        
        max_retries = 3
        retry_count = 0
        
        while retry_count < max_retries:
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
                return  # Success, exit the retry loop
                
            except Exception as e:
                retry_count += 1
                logger.error(f"Error playing song (attempt {retry_count}/{max_retries}): {e}")
                
                if retry_count >= max_retries:
                    # Remove the failed song and try next without recursion
                    logger.error(f"Failed to play song after {max_retries} attempts, removing from queue")
                    if self.queues[guild_id]:
                        failed_song = self.queues[guild_id].popleft()
                        logger.info(f"Removed failed song: {failed_song.get('title', 'Unknown')}")
                    
                    # Clean up current song reference
                    if guild_id in self.current_songs:
                        del self.current_songs[guild_id]
                    
                    # Try to play the next song (if any) without recursion
                    if self.queues[guild_id]:
                        # Schedule the next play attempt
                        asyncio.create_task(self._play_next(guild_id))
                    else:
                        # No more songs, start idle timer
                        await self.start_idle_timer(guild_id)
                    return
                else:
                    # Wait a bit before retrying
                    await asyncio.sleep(1)
    
    def _after_play(self, guild_id: int, error):
        """Called after a song finishes playing."""
        if error:
            logger.error(f"Player error: {error}")
        
        try:
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
            
            # Schedule next action using run_coroutine_threadsafe for thread safety
            if self.queues[guild_id]:
                coro = self._play_next(guild_id)
            else:
                # No more songs - clean up and start idle timer
                logger.info(f"Queue empty for guild {guild_id}, starting idle timer")
                if guild_id in self.current_songs:
                    del self.current_songs[guild_id]
                coro = self.start_idle_timer(guild_id)
            
            # Always use run_coroutine_threadsafe for thread safety from Discord callbacks
            try:
                if self.bot.loop and not self.bot.loop.is_closed():
                    future = asyncio.run_coroutine_threadsafe(coro, self.bot.loop)
                    # Add a callback to handle any exceptions
                    def handle_future_exception(fut):
                        try:
                            fut.result()  # This will raise any exception that occurred
                        except Exception as e:
                            logger.error(f"Exception in scheduled coroutine for guild {guild_id}: {e}")
                            # Perform emergency cleanup
                            self._emergency_cleanup(guild_id)
                    
                    future.add_done_callback(handle_future_exception)
                else:
                    logger.warning(f"Bot loop is closed or None, performing emergency cleanup for guild {guild_id}")
                    self._emergency_cleanup(guild_id)
                    
            except Exception as e:
                logger.error(f"Error scheduling coroutine in after_play for guild {guild_id}: {e}")
                self._emergency_cleanup(guild_id)
                
        except Exception as e:
            logger.error(f"Critical error in _after_play for guild {guild_id}: {e}")
            self._emergency_cleanup(guild_id)
    
    def _emergency_cleanup(self, guild_id: int):
        """Perform emergency cleanup when normal scheduling fails."""
        try:
            # Clean up current song reference
            if guild_id in self.current_songs:
                del self.current_songs[guild_id]
                
            # Cancel any existing idle timers
            if guild_id in self.idle_timers:
                self.idle_timers[guild_id].cancel()
                del self.idle_timers[guild_id]
                
            logger.warning(f"Performed emergency cleanup for guild {guild_id}")
            
        except Exception as cleanup_e:
            logger.error(f"Error in emergency cleanup for guild {guild_id}: {cleanup_e}")
    
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
        try:
            if self.is_24_7[guild_id]:
                logger.info(f"24/7 mode enabled for guild {guild_id}, skipping idle timer")
                return  # Don't start timer if 24/7 is enabled
            
            # Cancel existing timer
            if guild_id in self.idle_timers:
                self.idle_timers[guild_id].cancel()
                del self.idle_timers[guild_id]
            
            # Verify we have a valid guild and voice client before starting timer
            guild = self.bot.get_guild(guild_id)
            if not guild or not guild.voice_client:
                logger.info(f"No guild or voice client found for {guild_id}, skipping idle timer")
                return
            
            # Check if voice client is actually connected and in a channel
            if not guild.voice_client.is_connected():
                logger.warning(f"Voice client not connected for guild {guild_id}, cleaning up")
                await self.clear_queue(guild_id)
                return
                
            if not guild.voice_client.channel:
                logger.warning(f"Voice client has no channel for guild {guild_id}, cleaning up")
                await self.clear_queue(guild_id)
                return
            
            logger.info(f"Starting idle timer for guild {guild_id}")
            # Start new timer
            self.idle_timers[guild_id] = asyncio.create_task(
                self._idle_disconnect_timer(guild_id)
            )
        except Exception as e:
            logger.error(f"Error starting idle timer for guild {guild_id}: {e}")
            # Ensure we don't leave a broken timer reference
            if guild_id in self.idle_timers:
                try:
                    self.idle_timers[guild_id].cancel()
                except:
                    pass
                del self.idle_timers[guild_id]
    
    async def cancel_idle_timer(self, guild_id: int):
        """Cancel the idle timer."""
        try:
            if guild_id in self.idle_timers:
                self.idle_timers[guild_id].cancel()
                del self.idle_timers[guild_id]
                logger.debug(f"Cancelled idle timer for guild {guild_id}")
        except Exception as e:
            logger.error(f"Error cancelling idle timer for guild {guild_id}: {e}")
    
    async def _idle_disconnect_timer(self, guild_id: int):
        """Timer that disconnects the bot after 5 minutes of inactivity."""
        try:
            await asyncio.sleep(300)  # 5 minutes
            
            # Check if we should still disconnect
            if self.is_24_7[guild_id]:
                logger.info(f"24/7 mode enabled during idle timer for guild {guild_id}, cancelling disconnect")
                return
                
            guild = self.bot.get_guild(guild_id)
            if guild and guild.voice_client:
                # Validate voice client state before disconnecting
                if not guild.voice_client.is_connected():
                    logger.info(f"Voice client already disconnected for guild {guild_id}")
                    await self.clear_queue(guild_id)
                    return
                
                if not guild.voice_client.channel:
                    logger.info(f"Voice client no longer in channel for guild {guild_id}")
                    await self.clear_queue(guild_id)
                    return
                
                # Check if still not playing and not in 24/7 mode
                if not guild.voice_client.is_playing() and not self.is_24_7[guild_id]:
                    logger.info(LOG_MSG['idle_disconnect'].format(guild_id=guild_id))
                    await self._send_goodbye_message(guild_id)
                    
                    try:
                        await guild.voice_client.disconnect()
                        # Update bot status when leaving voice
                        await self._update_bot_presence_for_voice(joined_voice=False)
                    except Exception as e:
                        logger.error(f"Error disconnecting voice client for guild {guild_id}: {e}")
                    
                    await self.clear_queue(guild_id)
                else:
                    logger.info(f"Bot is playing or 24/7 mode enabled for guild {guild_id}, not disconnecting")
            else:
                logger.info(f"No guild or voice client found for guild {guild_id} during idle timer")
                
        except asyncio.CancelledError:
            logger.debug(f"Idle timer cancelled for guild {guild_id}")
        except Exception as e:
            logger.error(f"Error in idle disconnect timer for guild {guild_id}: {e}")
        finally:
            # Clean up timer reference
            if guild_id in self.idle_timers:
                del self.idle_timers[guild_id]
    
    async def _send_goodbye_message(self, guild_id: int):
        """Send goodbye message when disconnecting."""
        if guild_id in self.disconnect_messages:
            channel = self.disconnect_messages[guild_id]
            if channel:
                try:
                    embed = discord.Embed(
                        title="👋 Goodbye!",
                        description=MSG_INFO['goodbye_idle'],
                        color=discord.Color.orange()
                    )
                    await channel.send(embed=embed)
                except:
                    pass  # Ignore if we can't send the message
    
    async def _update_bot_presence_for_voice(self, joined_voice: bool = False):
        """Update bot presence to show deafened status when in voice channels."""
        try:
            if joined_voice:
                # Show deafened status for privacy
                activity = discord.Activity(
                    type=discord.ActivityType.listening,
                    name="🔇 Deafened for privacy"
                )
                status = discord.Status.online
            else:
                # Check if bot is still in any voice channels
                in_any_voice = any(guild.voice_client for guild in self.bot.guilds if guild.voice_client)
                
                if in_any_voice:
                    # Still in voice channels, keep deafened status
                    activity = discord.Activity(
                        type=discord.ActivityType.listening,
                        name="🔇 Deafened for privacy"
                    )
                else:
                    # Not in any voice channels, return to default status
                    from src.core.constants import BOT_DEFAULT_STATUS, BOT_DEFAULT_ACTIVITY_TYPE
                    activity = discord.Activity(
                        type=getattr(discord.ActivityType, BOT_DEFAULT_ACTIVITY_TYPE, discord.ActivityType.listening),
                        name=BOT_DEFAULT_STATUS
                    )
                status = discord.Status.online
            
            await self.bot.change_presence(activity=activity, status=status)
            logger.info(f"Updated bot presence: {activity.name}")
            
        except Exception as e:
            logger.error(f"Failed to update bot presence: {e}")

    def set_disconnect_channel(self, guild_id: int, channel):
        """Set the channel to send disconnect messages to."""
        self.disconnect_messages[guild_id] = channel
