"""
Quotes models for AI-powered quotes feature
"""
from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass
class Quote:
    """Represents a quote from AI or database."""
    
    content: str
    author: Optional[str] = None
    category: Optional[str] = None
    source: str = "ai"
    created_at: datetime = None
    
    def __post_init__(self):
        if self.created_at is None:
            self.created_at = datetime.now()
    
    def __str__(self) -> str:
        if self.author:
            return f'"{self.content}" - {self.author}'
        return f'"{self.content}"'


@dataclass
class QuoteRequest:
    """Represents a request for a quote."""
    
    category: Optional[str] = None
    mood: Optional[str] = None
    length: str = "medium"  # short, medium, long
    language: str = "en"
