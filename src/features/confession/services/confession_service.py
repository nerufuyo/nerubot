"""
Confession service for managing anonymous confessions
"""
import asyncio
from datetime import datetime, timedelta
from typing import Optional, List, Dict, Tuple
from src.features.confession.models.confession import (
    Confession, ConfessionReply, GuildConfessionSettings, ConfessionStatus
)
from src.core.utils.logging_utils import get_logger
import json
import os

logger = get_logger(__name__)


class ConfessionService:
    """Service for managing anonymous confessions."""
    
    def __init__(self):
        self.confessions: Dict[int, Confession] = {}
        self.replies: Dict[int, List[ConfessionReply]] = {}
        self.guild_settings: Dict[int, GuildConfessionSettings] = {}
        self.user_cooldowns: Dict[Tuple[int, int], datetime] = {}  # (user_id, guild_id) -> last_confession_time
        
        # File paths for persistence
        self.data_dir = "data/confessions"
        self.confessions_file = f"{self.data_dir}/confessions.json"
        self.replies_file = f"{self.data_dir}/replies.json"
        self.settings_file = f"{self.data_dir}/settings.json"
        
        # Create data directory if it doesn't exist
        os.makedirs(self.data_dir, exist_ok=True)
        
        # Load data
        self._load_data()
    
    def _load_data(self):
        """Load data from files."""
        try:
            # Load confessions
            if os.path.exists(self.confessions_file):
                with open(self.confessions_file, 'r') as f:
                    data = json.load(f)
                    for conf_id, conf_data in data.items():
                        # Handle both old string IDs and new numeric IDs
                        confession_id = conf_data['confession_id']
                        if isinstance(confession_id, str):
                            # Try to convert old string ID to numeric, or use a high number for migration
                            try:
                                confession_id = int(confession_id)
                            except ValueError:
                                # For old UUID-style IDs, we'll assign them high numbers during migration
                                confession_id = hash(confession_id) % 1000000 + 1000000
                        
                        confession = Confession(
                            confession_id=confession_id,
                            content=conf_data['content'],
                            author_id=conf_data['author_id'],
                            guild_id=conf_data['guild_id'],
                            channel_id=conf_data.get('channel_id'),
                            message_id=conf_data.get('message_id'),
                            image_url=conf_data.get('image_url'),
                            status=ConfessionStatus(conf_data.get('status', 'pending')),
                            created_at=datetime.fromisoformat(conf_data['created_at']),
                            posted_at=datetime.fromisoformat(conf_data['posted_at']) if conf_data.get('posted_at') else None,
                            reply_count=conf_data.get('reply_count', 0)
                        )
                        self.confessions[confession_id] = confession
            
            # Load replies
            if os.path.exists(self.replies_file):
                with open(self.replies_file, 'r') as f:
                    data = json.load(f)
                    for conf_id, replies_data in data.items():
                        # Convert confession ID to numeric
                        numeric_conf_id = conf_id
                        if isinstance(conf_id, str):
                            try:
                                numeric_conf_id = int(conf_id)
                            except ValueError:
                                numeric_conf_id = hash(conf_id) % 1000000 + 1000000
                        else:
                            numeric_conf_id = int(conf_id)
                        
                        replies = []
                        for reply_data in replies_data:
                            # Handle reply ID conversion
                            reply_id = reply_data['reply_id']
                            confession_id = reply_data['confession_id']
                            
                            if isinstance(reply_id, str):
                                try:
                                    reply_id = int(reply_id)
                                except ValueError:
                                    reply_id = hash(reply_id) % 1000000 + 1000000
                            
                            if isinstance(confession_id, str):
                                try:
                                    confession_id = int(confession_id)
                                except ValueError:
                                    confession_id = hash(confession_id) % 1000000 + 1000000
                            
                            reply = ConfessionReply(
                                reply_id=reply_id,
                                confession_id=confession_id,
                                content=reply_data['content'],
                                author_id=reply_data['author_id'],
                                guild_id=reply_data['guild_id'],
                                message_id=reply_data.get('message_id'),
                                image_url=reply_data.get('image_url'),
                                created_at=datetime.fromisoformat(reply_data['created_at']),
                                posted_at=datetime.fromisoformat(reply_data['posted_at']) if reply_data.get('posted_at') else None
                            )
                            replies.append(reply)
                        self.replies[numeric_conf_id] = replies
            
            # Load settings
            if os.path.exists(self.settings_file):
                with open(self.settings_file, 'r') as f:
                    data = json.load(f)
                    for guild_id_str, settings_data in data.items():
                        guild_id = int(guild_id_str)
                        settings = GuildConfessionSettings(
                            guild_id=guild_id,
                            confession_channel_id=settings_data.get('confession_channel_id'),
                            moderation_enabled=settings_data.get('moderation_enabled', False),
                            moderation_channel_id=settings_data.get('moderation_channel_id'),
                            anonymous_replies=settings_data.get('anonymous_replies', True),
                            max_confession_length=settings_data.get('max_confession_length', 2000),
                            max_reply_length=settings_data.get('max_reply_length', 1000),
                            cooldown_minutes=settings_data.get('cooldown_minutes', 5),
                            next_confession_id=settings_data.get('next_confession_id', 0),
                            next_reply_id=settings_data.get('next_reply_id', 0)
                        )
                        self.guild_settings[guild_id] = settings
                        
        except Exception as e:
            logger.error(f"Error loading confession data: {e}")
        
        # Perform migration to set initial ID counters if needed
        self._migrate_id_counters()
    
    def _save_data(self):
        """Save data to files."""
        try:
            # Save confessions
            confessions_data = {}
            for conf_id, confession in self.confessions.items():
                confessions_data[conf_id] = {
                    'confession_id': confession.confession_id,
                    'content': confession.content,
                    'author_id': confession.author_id,
                    'guild_id': confession.guild_id,
                    'channel_id': confession.channel_id,
                    'message_id': confession.message_id,
                    'image_url': confession.image_url,
                    'status': confession.status.value,
                    'created_at': confession.created_at.isoformat(),
                    'posted_at': confession.posted_at.isoformat() if confession.posted_at else None,
                    'reply_count': confession.reply_count
                }
            
            with open(self.confessions_file, 'w') as f:
                json.dump(confessions_data, f, indent=2)
            
            # Save replies
            replies_data = {}
            for conf_id, replies in self.replies.items():
                replies_data[conf_id] = []
                for reply in replies:
                    replies_data[conf_id].append({
                        'reply_id': reply.reply_id,
                        'confession_id': reply.confession_id,
                        'content': reply.content,
                        'author_id': reply.author_id,
                        'guild_id': reply.guild_id,
                        'message_id': reply.message_id,
                        'image_url': reply.image_url,
                        'created_at': reply.created_at.isoformat(),
                        'posted_at': reply.posted_at.isoformat() if reply.posted_at else None
                    })
            
            with open(self.replies_file, 'w') as f:
                json.dump(replies_data, f, indent=2)
            
            # Save settings
            settings_data = {}
            for guild_id, settings in self.guild_settings.items():
                settings_data[str(guild_id)] = {
                    'confession_channel_id': settings.confession_channel_id,
                    'moderation_enabled': settings.moderation_enabled,
                    'moderation_channel_id': settings.moderation_channel_id,
                    'anonymous_replies': settings.anonymous_replies,
                    'max_confession_length': settings.max_confession_length,
                    'max_reply_length': settings.max_reply_length,
                    'cooldown_minutes': settings.cooldown_minutes,
                    'next_confession_id': settings.next_confession_id,
                    'next_reply_id': settings.next_reply_id
                }
            
            with open(self.settings_file, 'w') as f:
                json.dump(settings_data, f, indent=2)
                
        except Exception as e:
            logger.error(f"Error saving confession data: {e}")
    
    def get_guild_settings(self, guild_id: int) -> GuildConfessionSettings:
        """Get guild confession settings."""
        if guild_id not in self.guild_settings:
            self.guild_settings[guild_id] = GuildConfessionSettings(guild_id=guild_id)
            self._save_data()
        return self.guild_settings[guild_id]
    
    def update_guild_settings(self, guild_id: int, **kwargs) -> GuildConfessionSettings:
        """Update guild confession settings."""
        settings = self.get_guild_settings(guild_id)
        
        for key, value in kwargs.items():
            if hasattr(settings, key):
                setattr(settings, key, value)
        
        self.guild_settings[guild_id] = settings
        self._save_data()
        return settings
    
    def check_user_cooldown(self, user_id: int, guild_id: int) -> Tuple[bool, Optional[int]]:
        """Check if user is on cooldown. Returns (is_on_cooldown, seconds_remaining)."""
        settings = self.get_guild_settings(guild_id)
        key = (user_id, guild_id)
        
        if key not in self.user_cooldowns:
            return False, None
        
        last_confession = self.user_cooldowns[key]
        cooldown_end = last_confession + timedelta(minutes=settings.cooldown_minutes)
        
        if datetime.now() < cooldown_end:
            remaining = (cooldown_end - datetime.now()).total_seconds()
            return True, int(remaining)
        
        return False, None
    
    def create_confession(self, content: str, author_id: int, guild_id: int, image_url: Optional[str] = None) -> Tuple[bool, str, Optional[Confession]]:
        """Create a new confession. Returns (success, message, confession)."""
        settings = self.get_guild_settings(guild_id)
        
        # Check content length
        if len(content) > settings.max_confession_length:
            return False, f"Confession too long! Maximum {settings.max_confession_length} characters allowed.", None
        
        # Check cooldown
        on_cooldown, remaining = self.check_user_cooldown(author_id, guild_id)
        if on_cooldown:
            minutes = remaining // 60
            seconds = remaining % 60
            return False, f"You're on cooldown! Please wait {minutes}m {seconds}s before submitting another confession.", None
        
        # Check if confession channel is set
        if not settings.confession_channel_id:
            return False, "Confession channel is not set up for this server. Please ask an admin to set it up.", None
        
        # Generate sequential confession ID
        confession_id = settings.next_confession_id
        
        # Increment the counter for next confession
        settings.next_confession_id += 1
        
        # Create confession
        confession = Confession(
            confession_id=confession_id,
            content=content,
            author_id=author_id,
            guild_id=guild_id,
            image_url=image_url,
            status=ConfessionStatus.APPROVED if not settings.moderation_enabled else ConfessionStatus.PENDING
        )
        
        self.confessions[confession_id] = confession
        self.user_cooldowns[(author_id, guild_id)] = datetime.now()
        self._save_data()
        
        return True, "Confession created successfully!", confession
    
    def create_reply(self, confession_id: int, content: str, author_id: int, guild_id: int, image_url: Optional[str] = None) -> Tuple[bool, str, Optional[ConfessionReply]]:
        """Create a reply to a confession. Returns (success, message, reply)."""
        settings = self.get_guild_settings(guild_id)
        
        # Check if confession exists
        if confession_id not in self.confessions:
            return False, "Confession not found!", None
        
        confession = self.confessions[confession_id]
        
        # Check if confession is in the same guild
        if confession.guild_id != guild_id:
            return False, "Confession not found in this server!", None
        
        # Check content length
        if len(content) > settings.max_reply_length:
            return False, f"Reply too long! Maximum {settings.max_reply_length} characters allowed.", None
        
        # Generate sequential reply ID
        reply_id = settings.next_reply_id
        
        # Increment the counter for next reply
        settings.next_reply_id += 1
        
        # Create reply
        reply = ConfessionReply(
            reply_id=reply_id,
            confession_id=confession_id,
            content=content,
            author_id=author_id,
            guild_id=guild_id,
            image_url=image_url
        )
        
        # Add to replies
        if confession_id not in self.replies:
            self.replies[confession_id] = []
        self.replies[confession_id].append(reply)
        
        # Update reply count
        confession.reply_count += 1
        
        self._save_data()
        
        return True, "Reply created successfully!", reply
    
    def get_confession(self, confession_id: int) -> Optional[Confession]:
        """Get a confession by ID."""
        return self.confessions.get(confession_id)
    
    def get_confession_replies(self, confession_id: int) -> List[ConfessionReply]:
        """Get all replies for a confession."""
        return self.replies.get(confession_id, [])
    
    def get_guild_confessions(self, guild_id: int, limit: int = 10) -> List[Confession]:
        """Get recent confessions for a guild."""
        guild_confessions = [
            confession for confession in self.confessions.values()
            if confession.guild_id == guild_id and confession.status == ConfessionStatus.APPROVED
        ]
        # Sort by creation time, newest first
        guild_confessions.sort(key=lambda x: x.created_at, reverse=True)
        return guild_confessions[:limit]
    
    def mark_confession_posted(self, confession_id: int, channel_id: int, message_id: int):
        """Mark a confession as posted."""
        if confession_id in self.confessions:
            confession = self.confessions[confession_id]
            confession.channel_id = channel_id
            confession.message_id = message_id
            confession.posted_at = datetime.now()
            self._save_data()
    
    def mark_reply_posted(self, reply_id: int, message_id: int):
        """Mark a reply as posted."""
        for replies in self.replies.values():
            for reply in replies:
                if reply.reply_id == reply_id:
                    reply.message_id = message_id
                    reply.posted_at = datetime.now()
                    self._save_data()
                    return
    
    def get_confession_by_tag(self, tag: str, guild_id: int) -> Optional[Confession]:
        """Get confession by tag (ID or partial ID)."""
        # Try to parse as exact numeric ID first
        try:
            confession_id = int(tag)
            confession = self.get_confession(confession_id)
            if confession and confession.guild_id == guild_id:
                return confession
        except ValueError:
            pass
        
        # For partial matches with numeric IDs, try string prefix matching
        matches = [
            confession for confession in self.confessions.values()
            if str(confession.confession_id).startswith(tag) and confession.guild_id == guild_id
        ]
        
        if len(matches) == 1:
            return matches[0]
        
        return None
    
    def _migrate_id_counters(self):
        """Migrate existing data to set initial ID counters for guilds."""
        for guild_id in self.guild_settings:
            settings = self.guild_settings[guild_id]
            
            # Find the highest confession ID for this guild
            max_confession_id = -1
            for confession in self.confessions.values():
                if confession.guild_id == guild_id and confession.confession_id > max_confession_id:
                    max_confession_id = confession.confession_id
            
            # Find the highest reply ID for this guild
            max_reply_id = -1
            for replies in self.replies.values():
                for reply in replies:
                    if reply.guild_id == guild_id and reply.reply_id > max_reply_id:
                        max_reply_id = reply.reply_id
            
            # Set next IDs only if they haven't been set yet (for backward compatibility)
            if settings.next_confession_id == 0 and max_confession_id >= 0:
                settings.next_confession_id = max_confession_id + 1
                
            if settings.next_reply_id == 0 and max_reply_id >= 0:
                settings.next_reply_id = max_reply_id + 1
        
        # Save the updated settings
        self._save_data()
