"""
Roast feature data models
"""
import time
import json
from dataclasses import dataclass, field
from typing import Dict, List, Optional, Any
from datetime import datetime, timedelta
from enum import Enum


class ActivityType(Enum):
    """Types of user activities we track"""
    MESSAGE = "message"
    VOICE_JOIN = "voice_join"
    VOICE_LEAVE = "voice_leave"
    MUSIC_REQUEST = "music_request"
    CONFESSION = "confession"
    COMMAND_USE = "command_use"
    EMOJI_REACTION = "emoji_reaction"
    LATE_NIGHT = "late_night"
    EARLY_MORNING = "early_morning"


@dataclass
class UserActivity:
    """Represents a single user activity"""
    user_id: int
    activity_type: ActivityType
    timestamp: float
    channel_id: int
    guild_id: int
    content: Optional[str] = None  # Message content, command name, etc.
    metadata: Dict[str, Any] = field(default_factory=dict)  # Extra data
    
    def __post_init__(self):
        if self.timestamp == 0.0:
            self.timestamp = time.time()
    
    @property
    def hour_of_day(self) -> int:
        """Get the hour of day (0-23) when this activity occurred"""
        return datetime.fromtimestamp(self.timestamp).hour
    
    @property
    def day_of_week(self) -> int:
        """Get day of week (0=Monday, 6=Sunday)"""
        return datetime.fromtimestamp(self.timestamp).weekday()
    
    def is_late_night(self) -> bool:
        """Check if activity happened late at night (11 PM - 3 AM)"""
        hour = self.hour_of_day
        return hour >= 23 or hour <= 3
    
    def is_early_morning(self) -> bool:
        """Check if activity happened early morning (4 AM - 7 AM)"""
        hour = self.hour_of_day
        return 4 <= hour <= 7


@dataclass
class UserBehaviorPattern:
    """Represents patterns in user behavior"""
    user_id: int
    most_active_hours: List[int] = field(default_factory=list)  # Hours 0-23
    most_active_days: List[int] = field(default_factory=list)   # 0=Mon, 6=Sun
    favorite_channels: List[int] = field(default_factory=list)
    common_words: List[str] = field(default_factory=list)
    activity_frequency: Dict[ActivityType, int] = field(default_factory=dict)
    late_night_percentage: float = 0.0
    message_length_avg: float = 0.0
    emoji_usage: Dict[str, int] = field(default_factory=dict)
    last_analyzed: float = 0.0
    
    def __post_init__(self):
        if self.last_analyzed == 0.0:
            self.last_analyzed = time.time()


@dataclass
class RoastTemplate:
    """Template for generating roasts"""
    category: str
    template: str
    conditions: Dict[str, Any] = field(default_factory=dict)
    severity: int = 1  # 1-5, 1 being mild, 5 being savage
    tags: List[str] = field(default_factory=list)


@dataclass
class UserRoastProfile:
    """Complete roast profile for a user"""
    user_id: int
    personality_traits: List[str] = field(default_factory=list)
    behavior_summary: str = ""
    roast_history: List[str] = field(default_factory=list)
    last_roasted: float = 0.0
    roast_count: int = 0
    immunity_level: int = 0  # Users can build "immunity" to repeated roasts
    
    def can_be_roasted(self, cooldown_hours: int = 6) -> bool:
        """Check if user can be roasted (cooldown protection)"""
        if self.last_roasted == 0:
            return True
        hours_since_last = (time.time() - self.last_roasted) / 3600
        return hours_since_last >= cooldown_hours
    
    def add_roast(self, roast_text: str):
        """Add a roast to user's history"""
        self.roast_history.append(roast_text)
        self.last_roasted = time.time()
        self.roast_count += 1
        
        # Keep only last 10 roasts
        if len(self.roast_history) > 10:
            self.roast_history = self.roast_history[-10:]


@dataclass
class RoastStats:
    """Statistics about roasting activity"""
    total_roasts_delivered: int = 0
    most_roasted_user: Optional[int] = None
    best_roast_rating: float = 0.0
    roasts_by_category: Dict[str, int] = field(default_factory=dict)
    daily_roast_count: int = 0
    last_reset: float = 0.0
    
    def reset_daily_stats(self):
        """Reset daily statistics"""
        self.daily_roast_count = 0
        self.last_reset = time.time()
    
    def should_reset_daily(self) -> bool:
        """Check if daily stats should be reset"""
        if self.last_reset == 0:
            return True
        hours_since_reset = (time.time() - self.last_reset) / 3600
        return hours_since_reset >= 24
