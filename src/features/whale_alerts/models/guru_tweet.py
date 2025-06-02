"""Guru tweet model."""
from dataclasses import dataclass
from datetime import datetime
from typing import Optional, List
from enum import Enum


class TweetSentiment(Enum):
    """Tweet sentiment classification."""
    BULLISH = "bullish"
    BEARISH = "bearish"
    NEUTRAL = "neutral"
    URGENT = "urgent"


@dataclass
class GuruTweet:
    """Class representing a tweet from a crypto guru."""
    
    tweet_id: str
    username: str
    display_name: str
    content: str
    timestamp: datetime
    retweet_count: int = 0
    like_count: int = 0
    reply_count: int = 0
    sentiment: Optional[TweetSentiment] = None
    mentioned_tokens: Optional[List[str]] = None
    is_verified: bool = False
    profile_image_url: Optional[str] = None
    tweet_url: Optional[str] = None
    
    def __post_init__(self):
        """Process tweet content for analysis."""
        if self.mentioned_tokens is None:
            self.mentioned_tokens = []
            
        if self.sentiment is None:
            self.sentiment = self._analyze_sentiment()
        
        if self.tweet_url is None:
            self.tweet_url = f"https://twitter.com/{self.username}/status/{self.tweet_id}"
    
    def _analyze_sentiment(self) -> TweetSentiment:
        """Simple sentiment analysis based on keywords."""
        content_lower = self.content.lower()
        
        # Urgent indicators
        urgent_words = ['urgent', 'breaking', 'alert', 'warning', 'emergency', 'now', '🚨']
        if any(word in content_lower for word in urgent_words):
            return TweetSentiment.URGENT
        
        # Bullish indicators
        bullish_words = ['moon', 'pump', 'bull', 'up', 'rise', 'buy', 'long', '🚀', '📈', '💎']
        bullish_count = sum(1 for word in bullish_words if word in content_lower)
        
        # Bearish indicators
        bearish_words = ['dump', 'bear', 'down', 'fall', 'sell', 'short', 'crash', '📉', '💀']
        bearish_count = sum(1 for word in bearish_words if word in content_lower)
        
        if bullish_count > bearish_count:
            return TweetSentiment.BULLISH
        elif bearish_count > bullish_count:
            return TweetSentiment.BEARISH
        else:
            return TweetSentiment.NEUTRAL
    
    def get_sentiment_color(self) -> int:
        """Get color based on sentiment."""
        color_map = {
            TweetSentiment.BULLISH: 0x00FF00,  # Green
            TweetSentiment.BEARISH: 0xFF0000,  # Red
            TweetSentiment.NEUTRAL: 0x808080,  # Gray
            TweetSentiment.URGENT: 0xFF6600,   # Orange
        }
        return color_map.get(self.sentiment, 0x808080)
    
    def get_sentiment_emoji(self) -> str:
        """Get emoji based on sentiment."""
        emoji_map = {
            TweetSentiment.BULLISH: "🟢",
            TweetSentiment.BEARISH: "🔴",
            TweetSentiment.NEUTRAL: "⚪",
            TweetSentiment.URGENT: "🟠",
        }
        return emoji_map.get(self.sentiment, "⚪")
    
    def get_engagement_score(self) -> int:
        """Calculate engagement score based on likes, retweets, and replies."""
        return self.like_count + (self.retweet_count * 2) + (self.reply_count * 1.5)
    
    def to_embed(self) -> dict:
        """Convert the guru tweet to a Discord embed."""
        title = f"{self.get_sentiment_emoji()} {self.display_name}"
        if self.is_verified:
            title += " ✅"
        
        # Truncate content if too long
        content = self.content
        if len(content) > 1000:
            content = content[:997] + "..."
        
        embed = {
            "title": title,
            "description": content,
            "url": self.tweet_url,
            "color": self.get_sentiment_color(),
            "timestamp": self.timestamp.isoformat(),
            "author": {
                "name": f"@{self.username}",
                "url": f"https://twitter.com/{self.username}",
            },
            "footer": {
                "text": f"Sentiment: {self.sentiment.value.title()} • Engagement: {int(self.get_engagement_score())}"
            },
            "fields": []
        }
        
        if self.profile_image_url:
            embed["author"]["icon_url"] = self.profile_image_url
        
        # Add mentioned tokens as a field
        if self.mentioned_tokens:
            tokens_text = ", ".join([f"${token}" for token in self.mentioned_tokens])
            embed["fields"].append({
                "name": "💰 Mentioned Tokens",
                "value": tokens_text,
                "inline": False
            })
        
        # Add engagement stats
        if self.like_count > 0 or self.retweet_count > 0:
            engagement_text = f"❤️ {self.like_count} • 🔄 {self.retweet_count}"
            if self.reply_count > 0:
                engagement_text += f" • 💬 {self.reply_count}"
            
            embed["fields"].append({
                "name": "📊 Engagement",
                "value": engagement_text,
                "inline": True
            })
        
        return embed
