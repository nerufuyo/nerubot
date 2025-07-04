"""
Chatbot data models
"""
import time
from dataclasses import dataclass
from typing import Optional, Dict, Any
from src.services.ai_service import AIProvider


@dataclass
class ChatSession:
    """Represents a chat session with a user"""
    user_id: int
    channel_id: int
    last_interaction: float = 0.0  # timestamp
    message_count: int = 0
    ai_provider: AIProvider = AIProvider.OPENAI
    is_active: bool = True
    
    def __post_init__(self):
        if self.last_interaction == 0.0:
            self.last_interaction = time.time()
    
    def update_interaction(self):
        """Update the last interaction timestamp"""
        self.last_interaction = time.time()
        self.message_count += 1
    
    def is_expired(self, timeout_minutes: int = 5) -> bool:
        """Check if session has expired"""
        return time.time() - self.last_interaction > (timeout_minutes * 60)
    
    def get_idle_time_minutes(self) -> float:
        """Get idle time in minutes"""
        return (time.time() - self.last_interaction) / 60


@dataclass
class ChatMessage:
    """Represents a chat message"""
    user_id: int
    channel_id: int
    content: str
    timestamp: float = 0.0
    ai_response: Optional[str] = None
    ai_provider: Optional[AIProvider] = None
    response_time: Optional[float] = None
    
    def __post_init__(self):
        if self.timestamp == 0.0:
            self.timestamp = time.time()


@dataclass
class ChatStats:
    """Chat statistics for a user"""
    user_id: int
    total_messages: int = 0
    total_ai_responses: int = 0
    favorite_provider: Optional[AIProvider] = None
    first_chat: Optional[float] = None
    last_chat: Optional[float] = None
    
    def update_stats(self, provider: AIProvider):
        """Update chat statistics"""
        self.total_messages += 1
        self.total_ai_responses += 1
        self.last_chat = time.time()
        
        if self.first_chat is None:
            self.first_chat = time.time()
        
        # Simple favorite provider tracking (could be more sophisticated)
        self.favorite_provider = provider
