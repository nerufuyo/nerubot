"""
Data models for the confession feature
"""
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional, List
from enum import Enum


class ConfessionStatus(Enum):
    """Status of a confession."""
    PENDING = "pending"
    APPROVED = "approved"
    REJECTED = "rejected"


@dataclass
class Confession:
    """A confession data model."""
    confession_id: int
    content: str
    author_id: int
    guild_id: int
    channel_id: Optional[int] = None
    message_id: Optional[int] = None
    image_url: Optional[str] = None
    status: ConfessionStatus = ConfessionStatus.PENDING
    created_at: datetime = field(default_factory=datetime.now)
    posted_at: Optional[datetime] = None
    reply_count: int = 0


@dataclass
class ConfessionReply:
    """A reply to a confession."""
    reply_id: int
    confession_id: int
    content: str
    author_id: int
    guild_id: int
    message_id: Optional[int] = None
    image_url: Optional[str] = None
    created_at: datetime = field(default_factory=datetime.now)
    posted_at: Optional[datetime] = None


@dataclass
class GuildConfessionSettings:
    """Guild-specific confession settings."""
    guild_id: int
    confession_channel_id: Optional[int] = None
    moderation_enabled: bool = False
    moderation_channel_id: Optional[int] = None
    anonymous_replies: bool = True
    max_confession_length: int = 2000
    max_reply_length: int = 1000
    cooldown_minutes: int = 5
    next_confession_id: int = 0
    next_reply_id: int = 0
