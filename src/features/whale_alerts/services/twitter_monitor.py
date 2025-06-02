"""Twitter monitoring service for crypto gurus."""
import asyncio
import logging
import re
from datetime import datetime, timedelta
from typing import List, Optional, Dict, Set
import tweepy
import aiohttp
from ..models.guru_tweet import GuruTweet, TweetSentiment


class TwitterMonitor:
    """Service to monitor Twitter accounts of crypto gurus."""
    
    def __init__(self, api_key: str, api_secret: str, access_token: str, access_token_secret: str):
        """Initialize Twitter monitor with API credentials."""
        self.logger = logging.getLogger(__name__)
        
        # Twitter API v2 client
        self.client = tweepy.Client(
            bearer_token=None,
            consumer_key=api_key,
            consumer_secret=api_secret,
            access_token=access_token,
            access_token_secret=access_token_secret,
            wait_on_rate_limit=True
        )
        
        # Default list of crypto gurus to monitor
        self.monitored_accounts = {
            # Add popular crypto influencers/gurus
            "elonmusk": {"display_name": "Elon Musk", "priority": "high"},
            "saylor": {"display_name": "Michael Saylor", "priority": "high"},
            "VitalikButerin": {"display_name": "Vitalik Buterin", "priority": "high"},
            "cz_binance": {"display_name": "Changpeng Zhao", "priority": "high"},
            "aantonop": {"display_name": "Andreas M. Antonopoulos", "priority": "medium"},
            "brian_armstrong": {"display_name": "Brian Armstrong", "priority": "medium"},
            "jack": {"display_name": "Jack Dorsey", "priority": "medium"},
            "APompliano": {"display_name": "Anthony Pompliano", "priority": "medium"},
            "naval": {"display_name": "Naval", "priority": "medium"},
            "RaoulGMI": {"display_name": "Raoul Pal", "priority": "medium"},
            "PeterSchiff": {"display_name": "Peter Schiff", "priority": "low"},
            "nayibbukele": {"display_name": "Nayib Bukele", "priority": "high"},
        }
        
        # Token patterns for detection
        self.token_patterns = [
            r'\$[A-Z]{3,10}',  # Standard token format like $BTC, $ETH
            r'#[A-Z]{3,10}',   # Hashtag format
            r'\b(bitcoin|btc|ethereum|eth|solana|sol|cardano|ada|polygon|matic|avalanche|avax|chainlink|link|uniswap|uni|aave|compound|sushiswap|sushi|pancakeswap|cake|dogecoin|doge|shiba|shib|pepe|floki)\b',
        ]
        
        self.last_check = {}
        self.is_monitoring = False
        
    async def start_monitoring(self, callback=None):
        """Start monitoring crypto gurus' tweets."""
        self.is_monitoring = True
        self.logger.info("Starting Twitter monitoring for crypto gurus")
        
        while self.is_monitoring:
            try:
                new_tweets = await self.fetch_new_tweets()
                
                if new_tweets and callback:
                    await callback(new_tweets)
                    
                # Wait 2 minutes before next check (to respect rate limits)
                await asyncio.sleep(120)
                
            except Exception as e:
                self.logger.error(f"Error in Twitter monitoring: {e}")
                await asyncio.sleep(300)  # Wait 5 minutes on error
    
    def stop_monitoring(self):
        """Stop monitoring."""
        self.is_monitoring = False
        self.logger.info("Stopped Twitter monitoring")
    
    async def fetch_new_tweets(self) -> List[GuruTweet]:
        """Fetch new tweets from monitored accounts."""
        new_tweets = []
        
        for username, account_info in self.monitored_accounts.items():
            try:
                tweets = await self._get_user_tweets(username, account_info)
                new_tweets.extend(tweets)
                
                # Small delay between requests
                await asyncio.sleep(1)
                
            except Exception as e:
                self.logger.error(f"Error fetching tweets for {username}: {e}")
                continue
        
        # Sort by engagement score and recency
        new_tweets.sort(key=lambda t: (t.get_engagement_score(), t.timestamp), reverse=True)
        
        return new_tweets
    
    async def _get_user_tweets(self, username: str, account_info: Dict) -> List[GuruTweet]:
        """Get recent tweets from a specific user."""
        tweets = []
        
        try:
            # Get user info
            user = self.client.get_user(username=username, user_fields=['verified', 'profile_image_url'])
            if not user.data:
                return tweets
            
            user_data = user.data
            
            # Calculate since_id based on last check
            since_time = self.last_check.get(username, datetime.utcnow() - timedelta(hours=1))
            
            # Get recent tweets
            tweets_response = self.client.get_users_tweets(
                id=user_data.id,
                max_results=10,
                tweet_fields=['created_at', 'public_metrics', 'text'],
                exclude=['retweets', 'replies']
            )
            
            if not tweets_response.data:
                return tweets
            
            for tweet in tweets_response.data:
                tweet_time = tweet.created_at.replace(tzinfo=None)
                
                # Skip if we've already processed this tweet
                if tweet_time <= since_time:
                    continue
                
                # Extract mentioned tokens
                mentioned_tokens = self._extract_tokens(tweet.text)
                
                # Only include tweets that mention crypto tokens or have high engagement
                if mentioned_tokens or tweet.public_metrics['like_count'] > 100:
                    guru_tweet = GuruTweet(
                        tweet_id=tweet.id,
                        username=username,
                        display_name=account_info['display_name'],
                        content=tweet.text,
                        timestamp=tweet_time,
                        retweet_count=tweet.public_metrics['retweet_count'],
                        like_count=tweet.public_metrics['like_count'],
                        reply_count=tweet.public_metrics['reply_count'],
                        mentioned_tokens=mentioned_tokens,
                        is_verified=user_data.verified or False,
                        profile_image_url=user_data.profile_image_url
                    )
                    tweets.append(guru_tweet)
            
            # Update last check time
            self.last_check[username] = datetime.utcnow()
            
        except Exception as e:
            self.logger.error(f"Error getting tweets for {username}: {e}")
        
        return tweets
    
    def _extract_tokens(self, text: str) -> List[str]:
        """Extract cryptocurrency tokens mentioned in the text."""
        tokens = set()
        text_lower = text.lower()
        
        for pattern in self.token_patterns:
            matches = re.findall(pattern, text_lower, re.IGNORECASE)
            for match in matches:
                # Clean up the token
                token = match.strip('$#').upper()
                if len(token) >= 3:
                    tokens.add(token)
        
        return list(tokens)
    
    def add_account(self, username: str, display_name: str, priority: str = "medium"):
        """Add a new account to monitor."""
        self.monitored_accounts[username] = {
            "display_name": display_name,
            "priority": priority
        }
        self.logger.info(f"Added {username} to monitoring list")
    
    def remove_account(self, username: str):
        """Remove an account from monitoring."""
        if username in self.monitored_accounts:
            del self.monitored_accounts[username]
            self.logger.info(f"Removed {username} from monitoring list")
    
    def get_monitored_accounts(self) -> Dict[str, Dict]:
        """Get list of monitored accounts."""
        return self.monitored_accounts.copy()
    
    async def get_account_info(self, username: str) -> Optional[Dict]:
        """Get information about a Twitter account."""
        try:
            user = self.client.get_user(
                username=username,
                user_fields=['verified', 'profile_image_url', 'public_metrics', 'description']
            )
            
            if user.data:
                return {
                    'id': user.data.id,
                    'username': user.data.username,
                    'name': user.data.name,
                    'verified': user.data.verified,
                    'followers_count': user.data.public_metrics['followers_count'],
                    'description': user.data.description,
                    'profile_image_url': user.data.profile_image_url
                }
        except Exception as e:
            self.logger.error(f"Error getting account info for {username}: {e}")
        
        return None
