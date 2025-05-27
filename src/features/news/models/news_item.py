"""News item model."""
from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass
class NewsItem:
    """Class representing a news item."""
    
    title: str
    source: str
    link: str
    description: str
    published_date: datetime
    category: Optional[str] = None
    image_url: Optional[str] = None
    
    def to_embed(self) -> dict:
        """Convert the news item to a Discord embed."""
        embed = {
            "title": self.title,
            "description": self.description,
            "url": self.link,
            "color": 0x3498db,  # Blue color
            "timestamp": self.published_date.isoformat(),
            "footer": {
                "text": f"Source: {self.source}"
            }
        }
        
        if self.image_url:
            embed["image"] = {"url": self.image_url}
            
        if self.category:
            embed["fields"] = [{"name": "Category", "value": self.category, "inline": True}]
            
        return embed
