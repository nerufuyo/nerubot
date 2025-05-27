"""
Spotify music source adapter.
"""
import asyncio
import spotipy
import os
import logging
from spotipy.oauth2 import SpotifyClientCredentials
from . import MusicSource, MusicSourceResult
from .youtube import YouTubeAdapter
from src.core.constants import (
    LOG_MSG, TIMEOUT_SPOTIFY_API, DEFAULT_UNKNOWN_DURATION, 
    DEFAULT_UNKNOWN_ALBUM, SPOTIFY_SEARCH_STRATEGIES
)

# Configure logger
logger = logging.getLogger(__name__)

class SpotifyAdapter:
    """Adapter for Spotify music."""
    
    def __init__(self):
        # Get credentials from environment variables
        client_id = os.getenv("SPOTIFY_CLIENT_ID")
        client_secret = os.getenv("SPOTIFY_CLIENT_SECRET")
        
        if client_id and client_secret:
            auth_manager = SpotifyClientCredentials(
                client_id=client_id,
                client_secret=client_secret
            )
            self.sp = spotipy.Spotify(auth_manager=auth_manager)
            self.initialized = True
        else:
            self.initialized = False
            logger.warning(LOG_MSG["spotify_no_credentials"])
    
    async def search(self, query: str):
        """Search for a song on Spotify."""
        if not self.initialized:
            logger.warning(LOG_MSG["spotify_not_initialized"])
            return None
        
        try:
            # Check if it's a Spotify URL
            if 'open.spotify.com' in query:
                if '/track/' in query:
                    return await self._process_track_url(query)
                elif '/album/' in query:
                    return await self._process_album_url(query)
                elif '/playlist/' in query:
                    return await self._process_playlist_url(query)
                elif '/artist/' in query:
                    return await self._process_artist_url(query)
            
            # Regular search with timeout
            loop = asyncio.get_event_loop()
            try:
                results = await asyncio.wait_for(
                    loop.run_in_executor(None, lambda: self.sp.search(q=query, limit=5)),
                    timeout=TIMEOUT_SPOTIFY_API
                )
            except asyncio.TimeoutError:
                logger.error(LOG_MSG["spotify_search_timeout"].format(query=query))
                return None
            
            if not results or not results['tracks']['items']:
                logger.warning(LOG_MSG["spotify_no_results"].format(query=query))
                return None
            
            spotify_results = []
            for track in results['tracks']['items']:
                result = self._convert_track_to_result(track)
                if result:
                    spotify_results.append(result)
            
            logger.info(LOG_MSG["spotify_results_found"].format(count=len(spotify_results), query=query))
            return spotify_results
                
        except Exception as e:
            logger.error(LOG_MSG["spotify_search_error"].format(query=query, error=e))
            return None
    
    async def _process_track_url(self, url):
        """Process a Spotify track URL."""
        try:
            # Extract track ID
            track_id = url.split('track/')[1].split('?')[0]
            
            loop = asyncio.get_event_loop()
            try:
                track = await asyncio.wait_for(
                    loop.run_in_executor(None, lambda: self.sp.track(track_id)),
                    timeout=TIMEOUT_SPOTIFY_API
                )
            except asyncio.TimeoutError:
                logger.error(LOG_MSG["spotify_track_timeout"].format(track_id=track_id))
                return None
            
            result = self._convert_track_to_result(track)
            return [result] if result else None
            
        except Exception as e:
            logger.error(LOG_MSG["spotify_track_error"].format(url=url, error=e))
            return None
    
    async def _process_album_url(self, url):
        """Process a Spotify album URL."""
        try:
            # Extract album ID
            album_id = url.split('album/')[1].split('?')[0]
            
            loop = asyncio.get_event_loop()
            album = await loop.run_in_executor(None, lambda: self.sp.album(album_id))
            
            if not album or 'tracks' not in album or not album['tracks']['items']:
                return None
            
            album_results = []
            for track in album['tracks']['items']:
                # Add album details to track
                track['album'] = {
                    'name': album['name'],
                    'images': album.get('images', [])
                }
                
                result = self._convert_track_to_result(track)
                if result:
                    album_results.append(result)
            
            return album_results
            
        except Exception as e:
            logger.error(LOG_MSG["spotify_album_error"].format(error=e))
            return None
    
    async def _process_playlist_url(self, url):
        """Process a Spotify playlist URL."""
        try:
            # Extract playlist ID
            playlist_id = url.split('playlist/')[1].split('?')[0]
            
            loop = asyncio.get_event_loop()
            
            # Get playlist tracks (handles pagination automatically)
            tracks = []
            results = await loop.run_in_executor(None, lambda: self.sp.playlist_items(playlist_id))
            
            tracks.extend(results['items'])
            while results['next']:
                results = await loop.run_in_executor(None, lambda: self.sp.next(results))
                tracks.extend(results['items'])
            
            if not tracks:
                return None
            
            playlist_results = []
            for item in tracks:
                track = item.get('track')
                if track:
                    result = self._convert_track_to_result(track)
                    if result:
                        playlist_results.append(result)
            
            return playlist_results
            
        except Exception as e:
            logger.error(LOG_MSG["spotify_playlist_error"].format(error=e))
            return None
    
    async def _process_artist_url(self, url):
        """Process a Spotify artist URL."""
        try:
            # Extract artist ID
            artist_id = url.split('artist/')[1].split('?')[0]
            
            loop = asyncio.get_event_loop()
            
            # Get artist's top tracks
            results = await loop.run_in_executor(None, lambda: self.sp.artist_top_tracks(artist_id))
            
            if not results or not results['tracks']:
                return None
            
            artist_results = []
            for track in results['tracks']:
                result = self._convert_track_to_result(track)
                if result:
                    artist_results.append(result)
            
            return artist_results
            
        except Exception as e:
            logger.error(LOG_MSG["spotify_artist_error"].format(error=e))
            return None
    
    def _convert_track_to_result(self, track):
        """Convert a Spotify track to a MusicSourceResult."""
        try:
            if not track:
                return None
            
            # Extract track details
            title = track['name']
            artists = ", ".join([artist['name'] for artist in track['artists']])
            album_name = track.get('album', {}).get('name', DEFAULT_UNKNOWN_ALBUM)
            
            # Format for searching on YouTube
            search_query = f"{title} {artists} audio"
            
            # Get duration
            duration_ms = track.get('duration_ms', 0)
            duration = DEFAULT_UNKNOWN_DURATION
            if duration_ms:
                seconds = int(duration_ms / 1000)
                minutes, seconds = divmod(seconds, 60)
                hours, minutes = divmod(minutes, 60)
                
                if hours > 0:
                    duration = f"{hours}:{minutes:02d}:{seconds:02d}"
                else:
                    duration = f"{minutes}:{seconds:02d}"
            
            # Get thumbnail
            thumbnail = None
            if 'album' in track and track['album'].get('images'):
                thumbnail = track['album']['images'][0]['url']
            
            # We can't play directly from Spotify, so we'll need to search for this on YouTube later
            return MusicSourceResult(
                title=f"{title} - {artists}",
                url=search_query,  # This is not a playable URL, just a search query for YouTube
                source=MusicSource.SPOTIFY,
                duration=duration,
                thumbnail=thumbnail,
                webpage_url=track.get('external_urls', {}).get('spotify'),
                artist=artists,
                album=album_name
            )
            
        except Exception as e:
            logger.error(LOG_MSG["spotify_convert_error"].format(error=e))
            return None
    
    @staticmethod
    async def convert_to_playable(spotify_result):
        """Convert a Spotify result to a playable YouTube result."""
        if not spotify_result:
            logger.error(LOG_MSG["spotify_convert_none"])
            return None
        
        try:
            logger.info(LOG_MSG["spotify_converting"].format(title=spotify_result.title))
            
            # Search for the song on YouTube with improved search query
            search_query = spotify_result.url  # This is the search query we created
            
            # Try multiple search strategies for better results
            search_queries = SPOTIFY_SEARCH_STRATEGIES
            formatted_queries = []
            for strategy in search_queries:
                formatted_query = strategy.format(
                    title=spotify_result.title.split(' - ')[0],  # Remove artist from title
                    artist=spotify_result.artist
                )
                formatted_queries.append(formatted_query)
            
            # Add the original search query as fallback
            formatted_queries.insert(0, search_query)
            
            youtube_result = None
            for query in formatted_queries:
                try:
                    logger.debug(LOG_MSG["spotify_youtube_trying"].format(query=query))
                    youtube_results = await YouTubeAdapter.search(query)
                    
                    if youtube_results and len(youtube_results) > 0 and youtube_results[0]:
                        youtube_result = youtube_results[0]
                        logger.info(LOG_MSG["spotify_youtube_found"].format(title=youtube_result.title))
                        break
                    else:
                        logger.debug(LOG_MSG["spotify_youtube_no_results"].format(query=query))
                        
                except Exception as search_e:
                    logger.warning(LOG_MSG["spotify_youtube_failed"].format(query=query, error=search_e))
                    continue
            
            if not youtube_result:
                logger.error(LOG_MSG["spotify_no_youtube"].format(title=spotify_result.title))
                return None
            
            # Create a new result with Spotify metadata but YouTube URL
            playable_result = MusicSourceResult(
                title=spotify_result.title,
                url=youtube_result.url,  # Playable YouTube URL
                source=MusicSource.SPOTIFY,  # Keep Spotify as source for UI
                duration=spotify_result.duration,
                thumbnail=spotify_result.thumbnail or youtube_result.thumbnail,
                webpage_url=spotify_result.webpage_url,
                artist=spotify_result.artist,
                album=spotify_result.album
            )
            
            logger.info(LOG_MSG["spotify_conversion_success"].format(title=playable_result.title))
            return playable_result
            
        except Exception as e:
            logger.error(LOG_MSG["spotify_convert_error"].format(error=e))
            return None
