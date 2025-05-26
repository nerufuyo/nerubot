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
from src.core.constants import LOG_MSG

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
        logger.info(LOG_MSG["source_search_start"].format(query=query))
        
        # Determine source based on the query
        source = self._determine_source(query)
        logger.debug(LOG_MSG["source_determined"].format(source=source))
        
        # Search based on the source
        try:
            if source == MusicSource.SPOTIFY:
                results = await self.spotify.search(query)
            elif source == MusicSource.SOUNDCLOUD:
                results = await self.soundcloud.search(query)
            else:
                # Default to YouTube
                results = await self.youtube.search(query)
            
            if results:
                logger.info(LOG_MSG["source_search_results"].format(count=len(results), query=query, source=source))
            else:
                logger.warning(LOG_MSG["source_no_results"].format(query=query, source=source))
                
            return results or []
            
        except Exception as e:
            logger.error(LOG_MSG["source_search_error"].format(query=query, source=source, error=e))
            return []
    
    async def get_playable(self, result: MusicSourceResult) -> Optional[MusicSourceResult]:
        """Get a playable version of the result."""
        if not result:
            logger.warning(LOG_MSG["source_convert_none"])
            return None
        
        logger.debug(LOG_MSG["source_convert_start"].format(title=result.title, source=result.source))
        
        try:
            # If it's already from YouTube, it's playable
            if result.source == MusicSource.YOUTUBE:
                logger.debug(LOG_MSG["source_youtube_playable"].format(title=result.title))
                return result
            
            # Convert based on source
            if result.source == MusicSource.SPOTIFY:
                playable = await self.spotify.convert_to_playable(result)
                if playable:
                    logger.info(LOG_MSG["source_spotify_converted"].format(title=playable.title))
                else:
                    logger.error(LOG_MSG["source_spotify_failed"].format(title=result.title))
                return playable
            elif result.source == MusicSource.SOUNDCLOUD:
                return await self.soundcloud.convert_to_playable(result)
            
            # Default to the original result
            logger.debug(LOG_MSG["source_using_original"].format(title=result.title))
            return result
            
        except Exception as e:
            logger.error(LOG_MSG["source_convert_error"].format(title=result.title, error=e))
            return None
    
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
