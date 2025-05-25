"""
Profile models for user profile management
"""
from dataclasses import dataclass
from datetime import datetime
from typing import Optional, List, Dict, Any


@dataclass
class UserProfile:
    """Represents a user's profile information."""
    
    user_id: int
    username: str
    display_name: Optional[str] = None
    bio: Optional[str] = None
    avatar_url: Optional[str] = None
    created_at: datetime = None
    updated_at: datetime = None
    preferences: Dict[str, Any] = None
    stats: Dict[str, int] = None
    
    def __post_init__(self):
        if self.created_at is None:
            self.created_at = datetime.now()
        if self.updated_at is None:
            self.updated_at = datetime.now()
        if self.preferences is None:
            self.preferences = {}
        if self.stats is None:
            self.stats = {
                "commands_used": 0,
                "songs_played": 0,
                "messages_sent": 0
            }


@dataclass
class ProfileSettings:
    """Represents user profile settings."""
    
    privacy_level: str = "public"  # public, friends, private
    show_stats: bool = True
    show_activity: bool = True
    timezone: str = "UTC"
    preferred_language: str = "en"
