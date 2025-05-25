"""
News models for the news feature
"""
from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass
class NewsArticle:
    """Represents a news article from RSS feed."""
    
    title: str
    description: str
    link: str
    published: datetime
    source: str
    author: Optional[str] = None
    image_url: Optional[str] = None
    
    def __str__(self) -> str:
        return f"{self.title} - {self.source}"
    
    @property
    def short_description(self) -> str:
        """Get a truncated description for embed display."""
        if len(self.description) > 200:
            return self.description[:200] + "..."
        return self.description


@dataclass
class RSSFeed:
    """Represents an RSS feed configuration."""
    
    name: str
    url: str
    category: str = "general"
    enabled: bool = True
    
    def __str__(self) -> str:
        return f"{self.name} ({self.category})"
