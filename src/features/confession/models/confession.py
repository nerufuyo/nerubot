"""
Confession models for anonymous confession feature
"""
from dataclasses import dataclass
from datetime import datetime
from typing import Optional
import uuid


@dataclass
class Confession:
    """Represents an anonymous confession."""
    
    id: str
    content: str
    channel_id: int
    guild_id: int
    created_at: datetime = None
    approved: bool = False
    moderated: bool = False
    reported: bool = False
    
    def __post_init__(self):
        if not self.id:
            self.id = str(uuid.uuid4())[:8]
        if self.created_at is None:
            self.created_at = datetime.now()


@dataclass
class ConfessionSettings:
    """Represents confession system settings for a guild."""
    
    guild_id: int
    enabled: bool = True
    moderation_required: bool = True
    max_length: int = 1000
    cooldown_minutes: int = 10
    confession_channel_id: Optional[int] = None
    mod_log_channel_id: Optional[int] = None
    blocked_words: list = None
    
    def __post_init__(self):
        if self.blocked_words is None:
            self.blocked_words = []
