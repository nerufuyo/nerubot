"""
Queue system for managing confession and reply creation to prevent duplicate IDs
"""
import asyncio
import json
import os
import threading
import time
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass, asdict
from enum import Enum
from src.core.utils.logging_utils import get_logger
from src.core.constants import QUEUE_CONSTANTS, CONFESSION_FILE_PATHS

logger = get_logger(__name__)


class QueueItemType(Enum):
    """Type of queue item."""
    CONFESSION = "confession"
    REPLY = "reply"


class QueueItemStatus(Enum):
    """Status of queue item."""
    PENDING = "pending"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"
    RETRYING = "retrying"


@dataclass
class QueueItem:
    """An item in the processing queue."""
    id: str  # Unique item ID
    type: QueueItemType
    guild_id: int
    user_id: int
    content: str
    attachments: Optional[List[str]] = None
    confession_id: Optional[int] = None  # For replies
    status: QueueItemStatus = QueueItemStatus.PENDING
    created_at: datetime = None
    processing_started_at: Optional[datetime] = None
    completed_at: Optional[datetime] = None
    retry_count: int = 0
    error_message: Optional[str] = None
    assigned_id: Optional[str] = None  # The final assigned CONF-XXX or REPLY-XXX-Y
    
    def __post_init__(self):
        if self.created_at is None:
            self.created_at = datetime.now()


class ConfessionQueue:
    """Thread-safe queue system for managing confession and reply creation."""
    
    def __init__(self):
        self.queue: Dict[str, QueueItem] = {}
        self.processing_items: Dict[str, QueueItem] = {}
        self.completed_items: Dict[str, QueueItem] = {}
        self.lock = threading.RLock()
        self.running = False
        self.processor_task = None
        
        # ID counters to prevent duplicates
        self.guild_confession_counters: Dict[int, int] = {}
        self.guild_reply_counters: Dict[int, Dict[int, str]] = {}  # guild_id -> {confession_id -> next_letter}
        
        # Load existing queue state
        self._load_queue_state()
        
        # Start the queue processor
        self.start_processor()
    
    def _load_queue_state(self):
        """Load queue state from file."""
        try:
            queue_file = CONFESSION_FILE_PATHS["queue_file"]
            if os.path.exists(queue_file):
                with open(queue_file, 'r') as f:
                    data = json.load(f)
                    
                # Load queue items
                for item_id, item_data in data.get('queue', {}).items():
                    item = QueueItem(
                        id=item_data['id'],
                        type=QueueItemType(item_data['type']),
                        guild_id=item_data['guild_id'],
                        user_id=item_data['user_id'],
                        content=item_data['content'],
                        attachments=item_data.get('attachments'),
                        confession_id=item_data.get('confession_id'),
                        status=QueueItemStatus(item_data.get('status', 'pending')),
                        created_at=datetime.fromisoformat(item_data['created_at']),
                        processing_started_at=datetime.fromisoformat(item_data['processing_started_at']) if item_data.get('processing_started_at') else None,
                        completed_at=datetime.fromisoformat(item_data['completed_at']) if item_data.get('completed_at') else None,
                        retry_count=item_data.get('retry_count', 0),
                        error_message=item_data.get('error_message'),
                        assigned_id=item_data.get('assigned_id')
                    )
                    
                    if item.status == QueueItemStatus.PENDING:
                        self.queue[item_id] = item
                    elif item.status == QueueItemStatus.PROCESSING:
                        # Reset processing items to pending on restart
                        item.status = QueueItemStatus.PENDING
                        item.processing_started_at = None
                        self.queue[item_id] = item
                    elif item.status == QueueItemStatus.COMPLETED:
                        self.completed_items[item_id] = item
                
                # Load counters
                self.guild_confession_counters = data.get('confession_counters', {})
                # Convert string keys to int
                self.guild_confession_counters = {int(k): v for k, v in self.guild_confession_counters.items()}
                
                self.guild_reply_counters = data.get('reply_counters', {})
                # Convert string keys to int
                self.guild_reply_counters = {
                    int(guild_id): {int(conf_id): letter for conf_id, letter in counters.items()}
                    for guild_id, counters in self.guild_reply_counters.items()
                }
                
        except Exception as e:
            logger.error(f"Error loading queue state: {e}")
    
    def _save_queue_state(self):
        """Save queue state to file."""
        try:
            os.makedirs(os.path.dirname(CONFESSION_FILE_PATHS["queue_file"]), exist_ok=True)
            
            # Combine all items for saving
            all_items = {}
            all_items.update(self.queue)
            all_items.update(self.processing_items)
            all_items.update(self.completed_items)
            
            # Convert to serializable format
            queue_data = {}
            for item_id, item in all_items.items():
                item_dict = asdict(item)
                item_dict['type'] = item.type.value
                item_dict['status'] = item.status.value
                item_dict['created_at'] = item.created_at.isoformat()
                if item.processing_started_at:
                    item_dict['processing_started_at'] = item.processing_started_at.isoformat()
                if item.completed_at:
                    item_dict['completed_at'] = item.completed_at.isoformat()
                queue_data[item_id] = item_dict
            
            data = {
                'queue': queue_data,
                'confession_counters': {str(k): v for k, v in self.guild_confession_counters.items()},
                'reply_counters': {
                    str(guild_id): {str(conf_id): letter for conf_id, letter in counters.items()}
                    for guild_id, counters in self.guild_reply_counters.items()
                }
            }
            
            with open(CONFESSION_FILE_PATHS["queue_file"], 'w') as f:
                json.dump(data, f, indent=2)
                
        except Exception as e:
            logger.error(f"Error saving queue state: {e}")
    
    def generate_unique_id(self) -> str:
        """Generate a unique ID for queue items."""
        timestamp = int(time.time() * 1000000)  # microseconds
        return f"queue_{timestamp}"
    
    def reserve_confession_id(self, guild_id: int) -> int:
        """Reserve the next confession ID for a guild."""
        with self.lock:
            if guild_id not in self.guild_confession_counters:
                self.guild_confession_counters[guild_id] = 1
            
            confession_id = self.guild_confession_counters[guild_id]
            self.guild_confession_counters[guild_id] += 1
            
            self._save_queue_state()
            return confession_id
    
    def reserve_reply_id(self, guild_id: int, confession_id: int) -> str:
        """Reserve the next reply ID for a confession."""
        with self.lock:
            if guild_id not in self.guild_reply_counters:
                self.guild_reply_counters[guild_id] = {}
            
            if confession_id not in self.guild_reply_counters[guild_id]:
                self.guild_reply_counters[guild_id][confession_id] = 'A'
            
            letter = self.guild_reply_counters[guild_id][confession_id]
            reply_id = f"REPLY-{confession_id:03d}-{letter}"
            
            # Increment to next letter
            self.guild_reply_counters[guild_id][confession_id] = chr(ord(letter) + 1)
            
            self._save_queue_state()
            return reply_id
    
    def add_confession(self, guild_id: int, user_id: int, content: str, 
                      attachments: Optional[List[str]] = None) -> Tuple[str, int]:
        """Add a confession to the queue and reserve its ID."""
        with self.lock:
            item_id = self.generate_unique_id()
            confession_id = self.reserve_confession_id(guild_id)
            assigned_id = f"CONF-{confession_id:03d}"
            
            item = QueueItem(
                id=item_id,
                type=QueueItemType.CONFESSION,
                guild_id=guild_id,
                user_id=user_id,
                content=content,
                attachments=attachments,
                assigned_id=assigned_id
            )
            
            self.queue[item_id] = item
            self._save_queue_state()
            
            logger.info(f"Queued confession {assigned_id} for guild {guild_id}")
            return item_id, confession_id
    
    def add_reply(self, guild_id: int, user_id: int, confession_id: int, content: str,
                 attachments: Optional[List[str]] = None) -> Tuple[str, str]:
        """Add a reply to the queue and reserve its ID."""
        with self.lock:
            item_id = self.generate_unique_id()
            reply_id = self.reserve_reply_id(guild_id, confession_id)
            
            item = QueueItem(
                id=item_id,
                type=QueueItemType.REPLY,
                guild_id=guild_id,
                user_id=user_id,
                content=content,
                attachments=attachments,
                confession_id=confession_id,
                assigned_id=reply_id
            )
            
            self.queue[item_id] = item
            self._save_queue_state()
            
            logger.info(f"Queued reply {reply_id} for confession {confession_id}")
            return item_id, reply_id
    
    def get_next_item(self) -> Optional[QueueItem]:
        """Get the next item to process from the queue."""
        with self.lock:
            if not self.queue:
                return None
            
            # Get oldest item
            oldest_item_id = min(self.queue.keys(), key=lambda x: self.queue[x].created_at)
            item = self.queue.pop(oldest_item_id)
            
            # Mark as processing
            item.status = QueueItemStatus.PROCESSING
            item.processing_started_at = datetime.now()
            self.processing_items[item.id] = item
            
            self._save_queue_state()
            return item
    
    def mark_completed(self, item_id: str, success: bool = True, error_message: str = None):
        """Mark an item as completed or failed."""
        with self.lock:
            if item_id in self.processing_items:
                item = self.processing_items.pop(item_id)
                
                if success:
                    item.status = QueueItemStatus.COMPLETED
                    item.completed_at = datetime.now()
                    self.completed_items[item_id] = item
                    logger.info(f"Completed processing {item.assigned_id}")
                else:
                    item.retry_count += 1
                    item.error_message = error_message
                    
                    if item.retry_count < QUEUE_CONSTANTS["retry_attempts"]:
                        item.status = QueueItemStatus.RETRYING
                        # Add back to queue with delay
                        time.sleep(QUEUE_CONSTANTS["retry_delay"])
                        item.status = QueueItemStatus.PENDING
                        item.processing_started_at = None
                        self.queue[item_id] = item
                        logger.warning(f"Retrying {item.assigned_id} (attempt {item.retry_count})")
                    else:
                        item.status = QueueItemStatus.FAILED
                        self.completed_items[item_id] = item
                        logger.error(f"Failed to process {item.assigned_id} after {item.retry_count} attempts: {error_message}")
                
                self._save_queue_state()
    
    def get_queue_status(self) -> Dict[str, int]:
        """Get current queue status."""
        with self.lock:
            return {
                'pending': len(self.queue),
                'processing': len(self.processing_items),
                'completed': len([item for item in self.completed_items.values() if item.status == QueueItemStatus.COMPLETED]),
                'failed': len([item for item in self.completed_items.values() if item.status == QueueItemStatus.FAILED])
            }
    
    def cleanup_old_items(self):
        """Clean up old completed items."""
        with self.lock:
            cutoff_time = datetime.now() - timedelta(hours=24)
            
            items_to_remove = []
            for item_id, item in self.completed_items.items():
                if item.completed_at and item.completed_at < cutoff_time:
                    items_to_remove.append(item_id)
            
            for item_id in items_to_remove:
                del self.completed_items[item_id]
            
            if items_to_remove:
                logger.info(f"Cleaned up {len(items_to_remove)} old queue items")
                self._save_queue_state()
    
    def start_processor(self):
        """Start the queue processor."""
        if not self.running:
            self.running = True
            logger.info("Starting confession queue processor")
    
    def stop_processor(self):
        """Stop the queue processor."""
        if self.running:
            self.running = False
            logger.info("Stopping confession queue processor")
    
    def get_item_by_id(self, item_id: str) -> Optional[QueueItem]:
        """Get a queue item by its ID."""
        with self.lock:
            return (self.queue.get(item_id) or 
                   self.processing_items.get(item_id) or 
                   self.completed_items.get(item_id))


# Global queue instance
confession_queue = ConfessionQueue()
