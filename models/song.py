"""
Song model for music functionality
"""
from dataclasses import dataclass
from typing import Optional
import discord

@dataclass
class Song:
    """Represents a song in the queue."""
    
    title: str
    url: str
    webpage_url: Optional[str] = None
    duration: Optional[str] = None
    requester: Optional[discord.Member] = None
    thumbnail: Optional[str] = None
    
    def __str__(self):
        return self.title
    
    def to_dict(self):
        """Convert to dictionary for serialization."""
        return {
            'title': self.title,
            'url': self.url,
            'webpage_url': self.webpage_url,
            'duration': self.duration,
            'requester_id': self.requester.id if self.requester else None,
            'thumbnail': self.thumbnail
        }
