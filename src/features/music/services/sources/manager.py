"""
Music source manager for handling different music sources.
"""
import re
import logging
from typing import List, Optional
from . import MusicSource, MusicSourceResult
from .youtube import YouTubeAdapter
from .spotify import SpotifyAdapter
from .soundcloud import SoundCloudAdapter

# Configure logger
logger = logging.getLogger(__name__)

class SourceManager:
    """Manager for handling different music sources."""
    
    def __init__(self):
        self.youtube = YouTubeAdapter()
        self.spotify = SpotifyAdapter()
        self.soundcloud = SoundCloudAdapter()
    
    async def search(self, query: str) -> List[MusicSourceResult]:
        """Search for music from the appropriate source."""
        # Determine source based on the query
        source = self._determine_source(query)
        
        # Search based on the source
        if source == MusicSource.SPOTIFY:
            results = await self.spotify.search(query)
        elif source == MusicSource.SOUNDCLOUD:
            results = await self.soundcloud.search(query)
        else:
            # Default to YouTube
            results = await self.youtube.search(query)
        
        return results or []
    
    async def get_playable(self, result: MusicSourceResult) -> Optional[MusicSourceResult]:
        """Get a playable version of the result."""
        if not result:
            return None
        
        # If it's already from YouTube, it's playable
        if result.source == MusicSource.YOUTUBE:
            return result
        
        # Convert based on source
        if result.source == MusicSource.SPOTIFY:
            return await self.spotify.convert_to_playable(result)
        elif result.source == MusicSource.SOUNDCLOUD:
            return await self.soundcloud.convert_to_playable(result)
        
        # Default to the original result
        return result
    
    def _determine_source(self, query: str) -> MusicSource:
        """Determine the music source based on the query."""
        # Check for Spotify
        if 'spotify.com' in query or query.startswith('spotify:'):
            return MusicSource.SPOTIFY
        
        # Check for SoundCloud
        if 'soundcloud.com' in query:
            return MusicSource.SOUNDCLOUD
        
        # Check for YouTube and common video sites
        if any(site in query for site in ['youtube.com', 'youtu.be', 'vimeo.com']):
            return MusicSource.YOUTUBE
        
        # Check for direct links
        if re.match(r'https?://.+\.(mp3|wav|ogg|flac)', query, re.IGNORECASE):
            return MusicSource.DIRECT
        
        # Default to YouTube for general searches
        return MusicSource.YOUTUBE
