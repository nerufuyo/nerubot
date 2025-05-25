"""
News service for fetching and managing RSS feeds
"""
import asyncio
import aiohttp
import feedparser
from datetime import datetime, timezone
from typing import List, Dict, Optional
from bs4 import BeautifulSoup
import logging

from src.features.news.models.news import NewsArticle, RSSFeed

logger = logging.getLogger(__name__)


class NewsService:
    """Service for managing news feeds and articles."""
    
    def __init__(self):
        self.feeds: Dict[str, RSSFeed] = {}
        self.session: Optional[aiohttp.ClientSession] = None
        self._load_default_feeds()
    
    def _load_default_feeds(self):
        """Load default RSS feeds."""
        default_feeds = [
            RSSFeed("BBC World", "http://feeds.bbci.co.uk/news/world/rss.xml", "world"),
            RSSFeed("CNN Top Stories", "http://rss.cnn.com/rss/edition.rss", "general"),
            RSSFeed("TechCrunch", "https://techcrunch.com/feed/", "technology"),
            RSSFeed("Hacker News", "https://hnrss.org/frontpage", "technology"),
            RSSFeed("Reuters World", "https://feeds.reuters.com/reuters/worldNews", "world"),
            RSSFeed("AP News", "https://feeds.apnews.com/rss/apf-topnews", "general"),
        ]
        
        for feed in default_feeds:
            self.feeds[feed.name.lower()] = feed
    
    async def __aenter__(self):
        """Async context manager entry."""
        self.session = aiohttp.ClientSession(
            timeout=aiohttp.ClientTimeout(total=30),
            headers={'User-Agent': 'NeruBot RSS Reader 1.0'}
        )
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        if self.session:
            await self.session.close()
    
    async def fetch_feed(self, feed: RSSFeed, max_articles: int = 5) -> List[NewsArticle]:
        """Fetch articles from an RSS feed."""
        if not self.session:
            raise RuntimeError("NewsService must be used as async context manager")
        
        try:
            async with self.session.get(feed.url) as response:
                if response.status != 200:
                    logger.warning(f"Failed to fetch {feed.name}: HTTP {response.status}")
                    return []
                
                content = await response.text()
                
            # Parse RSS feed
            parsed_feed = feedparser.parse(content)
            articles = []
            
            for entry in parsed_feed.entries[:max_articles]:
                try:
                    # Parse publication date
                    published = datetime.now(timezone.utc)
                    if hasattr(entry, 'published_parsed') and entry.published_parsed:
                        try:
                            time_struct = entry.published_parsed
                            if time_struct and isinstance(time_struct, (tuple, list)) and len(time_struct) >= 6:
                                # Convert to integers for datetime constructor
                                time_ints = [int(t) for t in time_struct[:6]]  # type: ignore
                                published = datetime(*time_ints, tzinfo=timezone.utc)
                        except (ValueError, TypeError):
                            pass
                    
                    # Clean description
                    description = self._clean_description(
                        getattr(entry, 'description', '') or 
                        getattr(entry, 'summary', '')
                    )
                    
                    # Extract image URL if available
                    image_url = self._extract_image_url(entry)
                    
                    article = NewsArticle(
                        title=str(getattr(entry, 'title', 'No title')),
                        description=description,
                        link=str(getattr(entry, 'link', '')),
                        published=published,
                        source=feed.name,
                        author=getattr(entry, 'author', None),
                        image_url=image_url
                    )
                    articles.append(article)
                    
                except Exception as e:
                    logger.warning(f"Error parsing article from {feed.name}: {e}")
                    continue
            
            logger.info(f"Fetched {len(articles)} articles from {feed.name}")
            return articles
            
        except Exception as e:
            logger.error(f"Error fetching feed {feed.name}: {e}")
            return []
    
    def _clean_description(self, html_description: str) -> str:
        """Clean HTML from description and truncate if needed."""
        if not html_description:
            return "No description available."
        
        # Remove HTML tags
        soup = BeautifulSoup(html_description, 'html.parser')
        text = soup.get_text(strip=True)
        
        # Truncate if too long
        if len(text) > 500:
            text = text[:500] + "..."
        
        return text or "No description available."
    
    def _extract_image_url(self, entry) -> Optional[str]:
        """Extract image URL from feed entry."""
        # Try different common fields for images
        if hasattr(entry, 'media_thumbnail') and entry.media_thumbnail:
            return entry.media_thumbnail[0].get('url')
        
        if hasattr(entry, 'media_content') and entry.media_content:
            for media in entry.media_content:
                if media.get('type', '').startswith('image/'):
                    return media.get('url')
        
        # Try to find image in description/content
        if hasattr(entry, 'description'):
            soup = BeautifulSoup(entry.description, 'html.parser')
            img = soup.find('img')
            if img and hasattr(img, 'get'):
                src = img.get('src')  # type: ignore
                if src and isinstance(src, str):
                    return src
        
        return None
    
    async def get_latest_news(self, category: Optional[str] = None, max_articles: int = 10) -> List[NewsArticle]:
        """Get latest news from all feeds or specific category."""
        if not self.session:
            raise RuntimeError("NewsService must be used as async context manager")
        
        # Filter feeds by category if specified
        feeds_to_fetch = []
        for feed in self.feeds.values():
            if not feed.enabled:
                continue
            if category is None or feed.category.lower() == category.lower():
                feeds_to_fetch.append(feed)
        
        if not feeds_to_fetch:
            return []
        
        # Fetch articles from all feeds concurrently
        tasks = [self.fetch_feed(feed, max_articles // len(feeds_to_fetch) + 1) 
                for feed in feeds_to_fetch]
        
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        # Combine all articles
        all_articles = []
        for result in results:
            if isinstance(result, list):
                all_articles.extend(result)
            else:
                logger.warning(f"Feed fetch error: {result}")
        
        # Sort by publication date (newest first) and limit
        all_articles.sort(key=lambda x: x.published, reverse=True)
        return all_articles[:max_articles]
    
    def get_available_categories(self) -> List[str]:
        """Get list of available news categories."""
        categories = set()
        for feed in self.feeds.values():
            if feed.enabled:
                categories.add(feed.category)
        return sorted(list(categories))
    
    def get_feed_names(self) -> List[str]:
        """Get list of available feed names."""
        return [feed.name for feed in self.feeds.values() if feed.enabled]
    
    def add_feed(self, name: str, url: str, category: str = "general") -> bool:
        """Add a new RSS feed."""
        try:
            feed = RSSFeed(name, url, category)
            self.feeds[name.lower()] = feed
            logger.info(f"Added new feed: {name}")
            return True
        except Exception as e:
            logger.error(f"Error adding feed {name}: {e}")
            return False
    
    def remove_feed(self, name: str) -> bool:
        """Remove an RSS feed."""
        name_lower = name.lower()
        if name_lower in self.feeds:
            del self.feeds[name_lower]
            logger.info(f"Removed feed: {name}")
            return True
        return False
