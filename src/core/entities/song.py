"""
Core Music entity for the Discord music bot
"""
from dataclasses import dataclass
from typing import Optional


@dataclass
class Song:
    """Represents a song to be played."""
    title: str
    url: str
    duration: int  # Duration in seconds
    requested_by: str
    thumbnail: Optional[str] = None
    
    @property
    def formatted_duration(self) -> str:
        """Returns the duration formatted as mm:ss or hh:mm:ss."""
        minutes, seconds = divmod(self.duration, 60)
        hours, minutes = divmod(minutes, 60)
        
        if hours > 0:
            return f"{hours:02d}:{minutes:02d}:{seconds:02d}"
        else:
            return f"{minutes:02d}:{seconds:02d}"
