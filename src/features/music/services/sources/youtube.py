"""
YouTube music source adapter.
"""
import asyncio
import yt_dlp as youtube_dl
import logging
from . import MusicSource, MusicSourceResult
from src.core.constants import LOG_MSG, DEFAULT_UNKNOWN_DURATION, DEFAULT_UNKNOWN_ARTIST, DEFAULT_UNKNOWN_TITLE

# Configure logger
logger = logging.getLogger(__name__)

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

ytdl = youtube_dl.YoutubeDL(ytdl_format_options)

class YouTubeAdapter:
    """Adapter for YouTube music."""
    
    @staticmethod
    async def search(query: str):
        """Search for a song on YouTube."""
        try:
            # Run in executor to prevent blocking
            loop = asyncio.get_event_loop()
            data = await loop.run_in_executor(None, lambda: ytdl.extract_info(query, download=False))
            
            if data is None:
                return None
            
            # Check if it's a playlist
            if 'entries' in data:
                # Playlist handling
                if not data['entries']:
                    return None
                
                results = []
                for entry in data['entries']:
                    if entry:
                        result = YouTubeAdapter._process_video_data(entry)
                        if result:
                            results.append(result)
                
                return results
            else:
                # Single video
                return [YouTubeAdapter._process_video_data(data)]
                
        except Exception as e:
            logger.error(LOG_MSG["youtube_search_error"].format(error=e))
            return None
    
    @staticmethod
    def _process_video_data(data):
        """Process video data from YouTube."""
        if not data:
            return None
        
        try:
            # Format duration
            duration = DEFAULT_UNKNOWN_DURATION
            if 'duration' in data and data['duration']:
                minutes, seconds = divmod(int(data['duration']), 60)
                hours, minutes = divmod(minutes, 60)
                
                if hours > 0:
                    duration = f"{hours}:{minutes:02d}:{seconds:02d}"
                else:
                    duration = f"{minutes}:{seconds:02d}"
            
            # Get artist from uploader or channel
            artist = data.get('uploader', data.get('channel', DEFAULT_UNKNOWN_ARTIST))
            
            return MusicSourceResult(
                title=data.get('title', DEFAULT_UNKNOWN_TITLE),
                url=data.get('url'),
                source=MusicSource.YOUTUBE,
                duration=duration,
                thumbnail=data.get('thumbnail'),
                webpage_url=data.get('webpage_url'),
                artist=artist,
                album=None
            )
        except Exception as e:
            logger.error(LOG_MSG["youtube_process_error"].format(error=e))
            return None
