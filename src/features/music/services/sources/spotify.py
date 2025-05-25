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
            logger.warning("Spotify credentials not found. Spotify support will be limited.")
    
    async def search(self, query: str):
        """Search for a song on Spotify."""
        if not self.initialized:
            logger.warning("Spotify adapter not initialized. Falling back to YouTube.")
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
            
            # Regular search
            loop = asyncio.get_event_loop()
            results = await loop.run_in_executor(None, lambda: self.sp.search(q=query, limit=5))
            
            if not results or not results['tracks']['items']:
                return None
            
            spotify_results = []
            for track in results['tracks']['items']:
                result = self._convert_track_to_result(track)
                if result:
                    spotify_results.append(result)
            
            return spotify_results
                
        except Exception as e:
            logger.error(f"Spotify search error: {e}")
            return None
    
    async def _process_track_url(self, url):
        """Process a Spotify track URL."""
        try:
            # Extract track ID
            track_id = url.split('track/')[1].split('?')[0]
            
            loop = asyncio.get_event_loop()
            track = await loop.run_in_executor(None, lambda: self.sp.track(track_id))
            
            result = self._convert_track_to_result(track)
            return [result] if result else None
            
        except Exception as e:
            logger.error(f"Error processing Spotify track URL: {e}")
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
            logger.error(f"Error processing Spotify album URL: {e}")
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
            logger.error(f"Error processing Spotify playlist URL: {e}")
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
            logger.error(f"Error processing Spotify artist URL: {e}")
            return None
    
    def _convert_track_to_result(self, track):
        """Convert a Spotify track to a MusicSourceResult."""
        try:
            if not track:
                return None
            
            # Extract track details
            title = track['name']
            artists = ", ".join([artist['name'] for artist in track['artists']])
            album_name = track.get('album', {}).get('name', 'Unknown Album')
            
            # Format for searching on YouTube
            search_query = f"{title} {artists} audio"
            
            # Get duration
            duration_ms = track.get('duration_ms', 0)
            duration = "Unknown"
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
            logger.error(f"Error converting Spotify track: {e}")
            return None
    
    @staticmethod
    async def convert_to_playable(spotify_result):
        """Convert a Spotify result to a playable YouTube result."""
        if not spotify_result:
            return None
        
        try:
            # Search for the song on YouTube
            youtube_results = await YouTubeAdapter.search(spotify_result.url)
            
            if not youtube_results or not youtube_results[0]:
                return None
            
            # Get the first YouTube result
            youtube_result = youtube_results[0]
            
            # Create a new result with Spotify metadata but YouTube URL
            return MusicSourceResult(
                title=spotify_result.title,
                url=youtube_result.url,  # Playable YouTube URL
                source=MusicSource.SPOTIFY,  # Keep Spotify as source for UI
                duration=spotify_result.duration,
                thumbnail=spotify_result.thumbnail or youtube_result.thumbnail,
                webpage_url=spotify_result.webpage_url,
                artist=spotify_result.artist,
                album=spotify_result.album
            )
            
        except Exception as e:
            logger.error(f"Error converting Spotify to playable: {e}")
            return None
