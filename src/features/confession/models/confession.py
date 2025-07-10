"""
Data models for the confession feature
"""
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional, List, Dict
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
    thread_id: Optional[int] = None
    attachments: Optional[List[str]] = None  # List of attachment URLs
    status: ConfessionStatus = ConfessionStatus.PENDING
    created_at: datetime = field(default_factory=datetime.now)
    posted_at: Optional[datetime] = None
    reply_count: int = 0


@dataclass
class ConfessionReply:
    """A reply to a confession."""
    reply_id: str  # Format: REPLY-{confession_id}-{letter}
    confession_id: int
    content: str
    author_id: int
    guild_id: int
    message_id: Optional[int] = None
    attachments: Optional[List[str]] = None  # List of attachment URLs
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
    next_confession_id: int = 1
    next_reply_letter: Dict[int, str] = field(default_factory=dict)  # confession_id -> next letter (A, B, C, etc.)
