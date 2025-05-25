"""
Music Service - Handles music playback logic
"""
import discord
import asyncio
import yt_dlp
from collections import defaultdict, deque
import logging

logger = logging.getLogger(__name__)

class MusicService:
    """Service for handling music playback."""
    
    def __init__(self, bot):
        self.bot = bot
        self.queues = defaultdict(deque)
        self.current_songs = {}
        
        # YouTube-DL options
        self.ytdl_options = {
            'format': 'bestaudio/best',
            'outtmpl': '%(extractor)s-%(id)s-%(title)s.%(ext)s',
            'restrictfilenames': True,
            'noplaylist': True,
            'nocheckcertificate': True,
            'ignoreerrors': False,
            'logtostderr': False,
            'quiet': True,
            'no_warnings': True,
            'default_search': 'auto',
            'source_address': '0.0.0.0',
        }
        
        # FFmpeg options
        self.ffmpeg_options = {
            'before_options': '-reconnect 1 -reconnect_streamed 1 -reconnect_delay_max 5',
            'options': '-vn'
        }
        
        self.ytdl = yt_dlp.YoutubeDL(self.ytdl_options)
    
    async def add_to_queue(self, guild_id: int, query: str, requester: discord.Member):
        """Add a song to the queue."""
        try:
            # Search for the song
            loop = asyncio.get_event_loop()
            data = await loop.run_in_executor(None, lambda: self.ytdl.extract_info(query, download=False))
            
            if 'entries' in data:
                # Playlist - take first result
                data = data['entries'][0]
            
            song_info = {
                'title': data.get('title', 'Unknown'),
                'url': data.get('url'),
                'webpage_url': data.get('webpage_url'),
                'duration': self._format_duration(data.get('duration')),
                'requester': requester
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
            source = discord.FFmpegPCMAudio(song['url'], **self.ffmpeg_options)
            
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
