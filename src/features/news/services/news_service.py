"""News service for fetching and managing news feeds."""
import asyncio
import feedparser
import logging
from datetime import datetime
from typing import Dict, List, Optional, Set
import aiohttp
from bs4 import BeautifulSoup

from src.features.news.models.news_item import NewsItem
from src.core.utils.logging_utils import get_logger

logger = get_logger(__name__)

# List of trusted news sources with their RSS feeds
DEFAULT_NEWS_SOURCES = {
    "BBC World": "http://feeds.bbci.co.uk/news/world/rss.xml",
    "BBC Business": "http://feeds.bbci.co.uk/news/business/rss.xml",
    "Reuters Top News": "https://feeds.reuters.com/reuters/topNews",
    "Reuters Business": "https://feeds.reuters.com/reuters/businessNews",
    "AP News": "https://feeds.apnews.com/rss/topstories",
    "CNN Top Stories": "http://rss.cnn.com/rss/edition.rss",
    "NPR News": "https://feeds.npr.org/1001/rss.xml",
    "Al Jazeera": "https://www.aljazeera.com/xml/rss/all.xml",
}

class NewsService:
    """Service for fetching and managing news feeds."""
    
    def __init__(self, update_interval: int = 10, max_news_items: int = 100):
        """
        Initialize the news service.
        
        Args:
            update_interval: How often to check for news updates (in minutes)
            max_news_items: Maximum number of news items to keep in memory
        """
        self.update_interval = update_interval
        self.max_news_items = max_news_items
        self.news_sources = DEFAULT_NEWS_SOURCES.copy()
        self.news_items: List[NewsItem] = []
        self.last_fetch_time: Dict[str, datetime] = {}
        self.published_ids: Set[str] = set()  # Track published news to avoid duplicates
        self.running = False
        self.update_task = None
        self.new_items_callback = None  # Callback for when new items are found
        
    async def start(self) -> None:
        """Start the news service."""
        if self.running:
            return
            
        self.running = True
        self.update_task = asyncio.create_task(self._update_loop())
        logger.info("News service started")
        
    async def stop(self) -> None:
        """Stop the news service."""
        if not self.running:
            return
            
        self.running = False
        if self.update_task:
            self.update_task.cancel()
            try:
                await self.update_task
            except asyncio.CancelledError:
                pass
        logger.info("News service stopped")
    
    async def _update_loop(self) -> None:
        """Periodically update news items."""
        while self.running:
            try:
                logger.info("Fetching news updates...")
                new_items = await self.fetch_all_news()
                
                # If there are new items and we have a callback, notify about them
                if new_items and self.new_items_callback:
                    await self.new_items_callback(new_items)
                
                logger.info(f"Fetched {len(self.news_items)} total news items ({len(new_items)} new). Next update in {self.update_interval} minutes.")
            except Exception as e:
                logger.error(f"Error fetching news: {e}")
                
            # Wait for the next update interval
            await asyncio.sleep(self.update_interval * 60)
    
    async def fetch_all_news(self) -> List[NewsItem]:
        """Fetch news from all configured sources."""
        all_new_items = []
        
        for source_name, feed_url in self.news_sources.items():
            try:
                new_items = await self._fetch_news_from_source(source_name, feed_url)
                all_new_items.extend(new_items)
            except Exception as e:
                logger.error(f"Error fetching news from {source_name}: {e}")
        
        # Update the main news items list, keeping the most recent ones up to max_news_items
        self.news_items = sorted(
            self.news_items + all_new_items,
            key=lambda item: item.published_date,
            reverse=True
        )[:self.max_news_items]
        
        return all_new_items
    
    async def _fetch_news_from_source(self, source_name: str, feed_url: str) -> List[NewsItem]:
        """Fetch news from a single source."""
        logger.info(f"Fetching news from {source_name}...")
        
        async with aiohttp.ClientSession() as session:
            try:
                async with session.get(feed_url) as response:
                    if response.status != 200:
                        logger.error(f"Failed to fetch {source_name} feed: HTTP {response.status}")
                        return []
                    
                    content = await response.text()
            except Exception as e:
                logger.error(f"Error connecting to {source_name} feed: {e}")
                return []
        
        # Parse the feed
        feed = feedparser.parse(content)
        if not feed.entries:
            logger.warning(f"No entries found in {source_name} feed")
            return []
        
        # Create NewsItem objects
        new_items = []
        for entry in feed.entries:
            try:
                # Generate a unique ID for this news item
                entry_id = f"{source_name}:{entry.link}"
                
                # Skip if we've already processed this item
                if entry_id in self.published_ids:
                    continue
                
                # Parse published date
                published_date = None
                if hasattr(entry, 'published_parsed'):
                    published_date = datetime(*entry.published_parsed[:6])
                elif hasattr(entry, 'updated_parsed'):
                    published_date = datetime(*entry.updated_parsed[:6])
                else:
                    published_date = datetime.now()
                
                # Get description
                description = entry.summary if hasattr(entry, 'summary') else ""
                if not description and hasattr(entry, 'description'):
                    description = entry.description
                
                # Clean up description (remove HTML)
                description = BeautifulSoup(description, 'html.parser').get_text()
                if len(description) > 500:
                    description = description[:497] + "..."
                
                # Get category
                category = None
                if hasattr(entry, 'tags') and entry.tags:
                    category = entry.tags[0].term
                
                # Create and add the news item
                news_item = NewsItem(
                    title=entry.title,
                    source=source_name,
                    link=entry.link,
                    description=description,
                    published_date=published_date,
                    category=category,
                    image_url=None  # We could try to extract this from the content if needed
                )
                
                new_items.append(news_item)
                self.published_ids.add(entry_id)
                
            except Exception as e:
                logger.error(f"Error processing news entry from {source_name}: {e}")
        
        logger.info(f"Fetched {len(new_items)} new items from {source_name}")
        return new_items
    
    def get_latest_news(self, count: int = 5) -> List[NewsItem]:
        """Get the latest news items."""
        return self.news_items[:count]
    
    def add_news_source(self, name: str, feed_url: str) -> bool:
        """Add a new news source."""
        if name in self.news_sources:
            return False
        
        self.news_sources[name] = feed_url
        return True
    
    def remove_news_source(self, name: str) -> bool:
        """Remove a news source."""
        if name not in self.news_sources:
            return False
        
        del self.news_sources[name]
        return True
    
    def get_news_sources(self) -> Dict[str, str]:
        """Get all configured news sources."""
        return self.news_sources.copy()
    
    def set_new_items_callback(self, callback) -> None:
        """Set a callback function to be called when new items are found."""
        self.new_items_callback = callback
