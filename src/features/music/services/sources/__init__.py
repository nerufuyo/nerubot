"""
Music source adapters for different music platforms.
This module contains adapters for different music platforms like YouTube, Spotify, SoundCloud, etc.
"""
from enum import Enum

class MusicSource(Enum):
    """Enum for different music sources."""
    YOUTUBE = 'youtube'
    SPOTIFY = 'spotify' 
    SOUNDCLOUD = 'soundcloud'
    DIRECT = 'direct'
    UNKNOWN = 'unknown'

class MusicSourceResult:
    """Class to represent a music search result."""
    def __init__(self, 
                 title: str, 
                 url: str, 
                 source: MusicSource,
                 duration: str = "Unknown", 
                 thumbnail: str = None,
                 webpage_url: str = None,
                 artist: str = None,
                 album: str = None):
        self.title = title
        self.url = url  # Playable URL
        self.source = source
        self.duration = duration
        self.thumbnail = thumbnail
        self.webpage_url = webpage_url or url
        self.artist = artist
        self.album = album

    def to_dict(self):
        """Convert to dictionary format."""
        return {
            'title': self.title,
            'url': self.url,
            'source': self.source.value,
            'duration': self.duration,
            'thumbnail': self.thumbnail,
            'webpage_url': self.webpage_url,
            'artist': self.artist,
            'album': self.album
        }
