"""
SoundCloud music source adapter.
"""
import asyncio
import re
import logging
import os
import json
import aiohttp
import bs4
from . import MusicSource, MusicSourceResult
from .youtube import YouTubeAdapter

# Configure logger
logger = logging.getLogger(__name__)

class SoundCloudAdapter:
    """Adapter for SoundCloud music."""
    
    def __init__(self):
        # Using yt-dlp for SoundCloud instead of direct API
        self.initialized = False
        logger.info("Using yt-dlp for SoundCloud playback")
    
    async def search(self, query: str):
        """Search for a song on SoundCloud."""
        try:
            # Use yt-dlp's SoundCloud search capability
            if 'soundcloud.com' in query:
                # If it's a SoundCloud URL, use it directly
                return await self._fallback_search(query)
            else:
                # Add SoundCloud prefix for search
                search_query = f"scsearch:{query}"
                return await self._fallback_search(search_query)
                
        except Exception as e:
            logger.error(f"SoundCloud search error: {e}")
            return None
    
    @staticmethod
    async def _fallback_search(query):
        """Use yt-dlp for SoundCloud search."""
        try:
            # Add SoundCloud prefix if it's not a URL
            if 'soundcloud.com' not in query and not query.startswith('scsearch:'):
                search_query = f"scsearch:{query}"
            else:
                search_query = query
            
            # Use YouTube adapter with SoundCloud search
            youtube_results = await YouTubeAdapter.search(search_query)
            
            if not youtube_results:
                return None
            
            # Convert to SoundCloud source
            soundcloud_results = []
            for result in youtube_results:
                if result:
                    # Create a new result with SoundCloud as source
                    soundcloud_results.append(MusicSourceResult(
                        title=result.title,
                        url=result.url,
                        source=MusicSource.SOUNDCLOUD,
                        duration=result.duration,
                        thumbnail=result.thumbnail,
                        webpage_url=result.webpage_url,
                        artist=result.artist,
                        album=None
                    ))
            
            return soundcloud_results
            
        except Exception as e:
            logger.error(f"Error in SoundCloud fallback search: {e}")
            return None
    
    @staticmethod
    async def convert_to_playable(soundcloud_result):
        """Convert a SoundCloud result to a playable result if needed."""
        if not soundcloud_result:
            return None
        
        try:
            # We're already using yt-dlp for SoundCloud, so the URL should be playable
            return soundcloud_result
            
        except Exception as e:
            logger.error(f"Error converting SoundCloud to playable: {e}")
            return None
