"""
Music service for handling audio playback in Discord voice channels.
"""
import asyncio
import os
import discord
import yt_dlp as youtube_dl
from collections import deque, defaultdict
from typing import Dict, Optional
from src.core.utils.logging_utils import get_logger
from src.core.utils.file_utils import find_ffmpeg

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

# Configure yt-dlp
ytdl_format_options = {
    'format': 'bestaudio/best',
    'outtmpl': '%(extractor)s-%(id)s-%(title)s.%(ext)s',
    'restrictfilenames': True,
    'noplaylist': False,
    'nocheckcertificate': True,
    'ignoreerrors': True,
    'logtostderr': False,
    'quiet': True,
    'no_warnings': True,
    'default_search': 'auto',
    'source_address': '0.0.0.0',
}

ffmpeg_options = {
    'before_options': '-reconnect 1 -reconnect_streamed 1 -reconnect_delay_max 5',
    'options': '-vn'
}

ffmpeg_executable = FFMPEG_PATH if FFMPEG_PATH else 'ffmpeg'

ytdl = youtube_dl.YoutubeDL(ytdl_format_options)

class MusicService:
    """Service for handling music playback in Discord voice channels."""
    
    def __init__(self, bot):
        self.bot = bot
        self.queues = defaultdict(deque)
        self.current_songs = {}
    
    async def add_to_queue(self, guild_id: int, query: str, requester):
        """Add a song to the queue."""
        try:
            # Search for the song
            loop = asyncio.get_event_loop()
            data = await loop.run_in_executor(None, lambda: ytdl.extract_info(query, download=False))
            
            if data is None:
                return {
                    'success': False,
                    'error': "Could not find any information for that song"
                }
            
            if 'entries' in data and data['entries']:
                # Playlist - take first result
                data = data['entries'][0]
                if data is None:
                    return {
                        'success': False,
                        'error': "No valid entries found in the playlist"
                    }
            
            song_info = {
                'title': data.get('title', 'Unknown'),
                'url': data.get('url'),
                'webpage_url': data.get('webpage_url'),
                'duration': self._format_duration(data.get('duration')),
                'requester': requester
            }
            
            # Validate that we have a URL to play
            if not song_info['url']:
                return {
                    'success': False,
                    'error': "No playable URL found for that song"
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
                    'duration': song_info['duration']
                }
            
            return {
                'success': True,
                'queued': True,
                'title': song_info['title'],
                'duration': song_info['duration'],
                'position': len(self.queues[guild_id])
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
            
            logger.info(f"Now playing: {song['title']} in guild {guild_id}")
            
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
        
        # Remove the finished song from queue
        if self.queues[guild_id]:
            self.queues[guild_id].popleft()
        
        # Remove from current songs
        if guild_id in self.current_songs:
            del self.current_songs[guild_id]
        
        # Play next song
        coro = self._play_next(guild_id)
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
