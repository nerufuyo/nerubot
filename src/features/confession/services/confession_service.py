"""
Confession service for managing anonymous confessions with queue system
"""
import asyncio
from datetime import datetime, timedelta
from typing import Optional, List, Dict, Tuple
from src.features.confession.models.confession import (
    Confession, ConfessionReply, GuildConfessionSettings, ConfessionStatus
)
from src.features.confession.services.queue_service import confession_queue, QueueItem
from src.core.utils.logging_utils import get_logger
from src.core.constants import (
    CONFESSION_CONSTANTS, CONFESSION_FILE_PATHS, CONFESSION_LOG_MESSAGES
)
import json
import os

logger = get_logger(__name__)


class ConfessionService:
    """Service for managing anonymous confessions with queue system."""
    
    def __init__(self):
        self.confessions: Dict[int, Confession] = {}
        self.replies: Dict[int, List[ConfessionReply]] = {}
        self.guild_settings: Dict[int, GuildConfessionSettings] = {}

        
        # File paths for persistence
        self.data_dir = CONFESSION_FILE_PATHS["data_dir"]
        self.confessions_file = CONFESSION_FILE_PATHS["confessions_file"]
        self.replies_file = CONFESSION_FILE_PATHS["replies_file"]
        self.settings_file = CONFESSION_FILE_PATHS["settings_file"]
        
        # Create data directory if it doesn't exist
        os.makedirs(self.data_dir, exist_ok=True)
        
        # Initialize queue processor for this service BEFORE loading data
        self.queue = confession_queue
        
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
                        confession_id = conf_data['confession_id']
                        if isinstance(confession_id, str):
                            try:
                                confession_id = int(confession_id)
                            except ValueError:
                                confession_id = hash(confession_id) % 1000000 + 1000000
                        
                        confession = Confession(
                            confession_id=confession_id,
                            content=conf_data['content'],
                            author_id=conf_data['author_id'],
                            guild_id=conf_data['guild_id'],
                            channel_id=conf_data.get('channel_id'),
                            message_id=conf_data.get('message_id'),
                            thread_id=conf_data.get('thread_id'),
                            attachments=conf_data.get('attachments', []),
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
                        numeric_conf_id = int(conf_id)
                        
                        replies = []
                        for reply_data in replies_data:
                            reply = ConfessionReply(
                                reply_id=reply_data['reply_id'],
                                confession_id=reply_data['confession_id'],
                                content=reply_data['content'],
                                author_id=reply_data['author_id'],
                                guild_id=reply_data['guild_id'],
                                message_id=reply_data.get('message_id'),
                                attachments=reply_data.get('attachments', []),
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
                            next_confession_id=settings_data.get('next_confession_id', 1),
                            next_reply_letter=settings_data.get('next_reply_letter', {})
                        )
                        # Convert string keys to int keys for next_reply_letter
                        if isinstance(settings.next_reply_letter, dict):
                            settings.next_reply_letter = {int(k): v for k, v in settings.next_reply_letter.items()}
                        self.guild_settings[guild_id] = settings
                        
        except Exception as e:
            logger.error(f"Error loading confession data: {e}")
        
        # Perform migration to set initial ID counters if needed
        self._migrate_id_counters()
    
    def _migrate_id_counters(self):
        """Migrate ID counters from guild settings to queue system."""
        try:
            # Sync ID counters with queue system
            for guild_id, settings in self.guild_settings.items():
                # Update queue counters from settings
                self.queue.guild_confession_counters[guild_id] = settings.next_confession_id
                if settings.next_reply_letter:
                    if guild_id not in self.queue.guild_reply_counters:
                        self.queue.guild_reply_counters[guild_id] = {}
                    self.queue.guild_reply_counters[guild_id].update(settings.next_reply_letter)
            
            # Save queue state after migration
            self.queue._save_queue_state()
            logger.info("ID counter migration completed")
        except Exception as e:
            logger.error(f"Error during ID counter migration: {e}")
        
        # ...existing code...
    
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
                    'thread_id': confession.thread_id,
                    'attachments': confession.attachments or [],
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
                        'attachments': reply.attachments or [],
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
                    'next_confession_id': settings.next_confession_id,
                    'next_reply_letter': {str(k): v for k, v in settings.next_reply_letter.items()}
                }
            
            with open(self.settings_file, 'w') as f:
                json.dump(settings_data, f, indent=2)
                
        except Exception as e:
            logger.error(f"Error saving confession data: {e}")
    
    def get_guild_settings(self, guild_id: int) -> GuildConfessionSettings:
        """Get or create guild settings."""
        if guild_id not in self.guild_settings:
            defaults = CONFESSION_CONSTANTS["settings"]["defaults"]
            self.guild_settings[guild_id] = GuildConfessionSettings(
                guild_id=guild_id,
                max_confession_length=defaults["max_confession_length"],
                max_reply_length=defaults["max_reply_length"]
            )
        return self.guild_settings[guild_id]
    
    def update_guild_settings(self, guild_id: int, **kwargs):
        """Update guild settings."""
        settings = self.get_guild_settings(guild_id)
        for key, value in kwargs.items():
            if hasattr(settings, key):
                setattr(settings, key, value)
        self._save_data()
    
    def create_confession(self, content: str, author_id: int, guild_id: int, 
                         attachments: Optional[List[str]] = None) -> Tuple[bool, str, Optional[str]]:
        """Create a new confession using the queue system."""
        settings = self.get_guild_settings(guild_id)
        
        # Check content length
        max_length = CONFESSION_CONSTANTS["settings"]["defaults"]["max_confession_length"]
        if len(content) > max_length:
            error_msg = CONFESSION_CONSTANTS["messages"]["errors"]["content_too_long"].format(
                type="Confession", max=max_length
            )
            return False, error_msg, None
        
        # Check if confession channel is set
        if not settings.confession_channel_id:
            error_msg = CONFESSION_CONSTANTS["messages"]["errors"]["no_confession_channel"]
            return False, error_msg, None
        
        # Add to queue and get reserved ID
        try:
            item_id, confession_id = self.queue.add_confession(
                guild_id=guild_id,
                user_id=author_id,
                content=content,
                attachments=attachments or []
            )
            
            # Create confession object with reserved ID
            confession = Confession(
                confession_id=confession_id,
                content=content,
                author_id=author_id,
                guild_id=guild_id,
                attachments=attachments or [],
                status=ConfessionStatus.APPROVED  # Auto-approve for now
            )
            
            # Store confession
            self.confessions[confession_id] = confession
            
            # Initialize replies list
            self.replies[confession_id] = []
            
            self._save_data()
            
            assigned_id = CONFESSION_CONSTANTS["ids"]["id_format"].format(
                prefix=CONFESSION_CONSTANTS["ids"]["confession_prefix"],
                id=confession_id
            )
            
            logger.info(CONFESSION_LOG_MESSAGES["confession"]["created"].format(
                id=assigned_id, user_id=author_id, guild_id=guild_id
            ))
            logger.info(CONFESSION_LOG_MESSAGES["confession"]["queued"].format(id=assigned_id))
            
            return True, CONFESSION_CONSTANTS["messages"]["success"]["confession_created"].format(id=assigned_id), assigned_id
            
        except Exception as e:
            logger.error(CONFESSION_LOG_MESSAGES["confession"]["failed_create"].format(error=str(e)))
            return False, f"Failed to create confession: {str(e)}", None
    
    def create_reply(self, confession_id: int, content: str, author_id: int, 
                    guild_id: int, attachments: Optional[List[str]] = None) -> Tuple[bool, str, Optional[str]]:
        """Create a reply to a confession using the queue system."""
        settings = self.get_guild_settings(guild_id)
        
        # Check if confession exists
        if confession_id not in self.confessions:
            error_msg = CONFESSION_CONSTANTS["messages"]["errors"]["confession_not_found"]
            return False, error_msg, None
        
        confession = self.confessions[confession_id]
        if confession.guild_id != guild_id:
            error_msg = CONFESSION_CONSTANTS["messages"]["errors"]["confession_not_in_guild"]
            return False, error_msg, None
        
        # Check content length
        max_length = CONFESSION_CONSTANTS["settings"]["defaults"]["max_reply_length"]
        if len(content) > max_length:
            error_msg = CONFESSION_CONSTANTS["messages"]["errors"]["content_too_long"].format(
                type="Reply", max=max_length
            )
            return False, error_msg, None
        
        # Add to queue and get reserved ID
        try:
            item_id, reply_id = self.queue.add_reply(
                guild_id=guild_id,
                user_id=author_id,
                confession_id=confession_id,
                content=content,
                attachments=attachments or []
            )
            
            # Create reply object
            reply = ConfessionReply(
                reply_id=reply_id,
                confession_id=confession_id,
                content=content,
                author_id=author_id,
                guild_id=guild_id,
                attachments=attachments or []
            )
            
            # Add to replies
            if confession_id not in self.replies:
                self.replies[confession_id] = []
            self.replies[confession_id].append(reply)
            
            # Update confession reply count
            confession.reply_count += 1
            
            self._save_data()
            
            logger.info(CONFESSION_LOG_MESSAGES["reply"]["created"].format(
                id=reply_id, confession_id=confession_id, user_id=author_id
            ))
            logger.info(CONFESSION_LOG_MESSAGES["reply"]["queued"].format(id=reply_id))
            
            return True, CONFESSION_CONSTANTS["messages"]["success"]["reply_created"].format(id=reply_id), reply_id
            
        except Exception as e:
            logger.error(CONFESSION_LOG_MESSAGES["reply"]["failed_create"].format(error=str(e)))
            return False, f"Failed to create reply: {str(e)}", None
    
    def get_confession(self, confession_id: int) -> Optional[Confession]:
        """Get a confession by ID."""
        return self.confessions.get(confession_id)
    
    def get_confession_by_tag(self, tag: str, guild_id: int) -> Optional[Confession]:
        """Get a confession by tag (e.g., 'CONF-001' or '1')."""
        # Try to parse as direct ID first
        try:
            confession_id = int(tag)
            confession = self.get_confession(confession_id)
            if confession and confession.guild_id == guild_id:
                return confession
        except ValueError:
            pass
        
        # Try to parse as CONF-XXX format
        confession_prefix = CONFESSION_CONSTANTS["ids"]["confession_prefix"]
        if tag.upper().startswith(f'{confession_prefix}-'):
            try:
                confession_id = int(tag[len(confession_prefix)+1:])
                confession = self.get_confession(confession_id)
                if confession and confession.guild_id == guild_id:
                    return confession
            except ValueError:
                pass
        
        return None
    
    def get_queue_status(self) -> Dict[str, int]:
        """Get current queue status."""
        return self.queue.get_queue_status()
    
    def process_queue_item(self, item: QueueItem) -> bool:
        """Process a queue item (to be called by queue processor)."""
        try:
            if item.type.value == "confession":
                # The confession was already created in create_confession
                # Here we would handle any additional processing like posting to Discord
                return True
            elif item.type.value == "reply":
                # The reply was already created in create_reply
                # Here we would handle any additional processing like posting to Discord
                return True
            return False
        except Exception as e:
            logger.error(f"Error processing queue item {item.id}: {e}")
            return False
    
    def get_guild_confessions(self, guild_id: int, limit: int = 50) -> List[Confession]:
        """Get confessions for a guild."""
        confessions = [c for c in self.confessions.values() if c.guild_id == guild_id]
        confessions.sort(key=lambda x: x.created_at, reverse=True)
        return confessions[:limit]
    
    def get_confession_replies(self, confession_id: int) -> List[ConfessionReply]:
        """Get replies for a confession."""
        return self.replies.get(confession_id, [])
    
    def mark_confession_posted(self, confession_id: int, channel_id: int, message_id: int, thread_id: int):
        """Mark a confession as posted."""
        if confession_id in self.confessions:
            confession = self.confessions[confession_id]
            confession.channel_id = channel_id
            confession.message_id = message_id
            confession.thread_id = thread_id
            confession.posted_at = datetime.now()
            confession.status = ConfessionStatus.APPROVED
            self._save_data()
    
    def mark_reply_posted(self, reply_id: str, message_id: int):
        """Mark a reply as posted."""
        for replies in self.replies.values():
            for reply in replies:
                if reply.reply_id == reply_id:
                    reply.message_id = message_id
                    reply.posted_at = datetime.now()
                    self._save_data()
                    return
