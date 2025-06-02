"""Alerts broadcaster service for whale alerts and guru tweets."""
import asyncio
import logging
from datetime import datetime, timedelta
from typing import List, Optional, Callable, Dict, Any
import discord
from ..models.whale_alert import WhaleAlert
from ..models.guru_tweet import GuruTweet


class AlertsBroadcaster:
    """Service to broadcast whale alerts and guru tweets to Discord channels."""
    
    def __init__(self, bot):
        """Initialize alerts broadcaster."""
        self.bot = bot
        self.logger = logging.getLogger(__name__)
        
        # Channel configurations per guild
        self.guild_configs = {}
        
        # Rate limiting to prevent spam
        self.last_whale_alert = {}
        self.last_guru_tweet = {}
        self.min_whale_interval = 300  # 5 minutes between whale alerts
        self.min_tweet_interval = 600  # 10 minutes between guru tweets
        
        # Alert filtering
        self.whale_filters = {
            'min_usd_amount': 100_000,
            'priority_tokens': ['BTC', 'ETH', 'SOL', 'ADA', 'MATIC'],
            'high_priority_amount': 1_000_000,
        }
        
        self.tweet_filters = {
            'min_engagement': 100,
            'priority_accounts': ['elonmusk', 'saylor', 'VitalikButerin', 'cz_binance'],
            'urgent_only': False,
        }
    
    async def broadcast_whale_alerts(self, alerts: List[WhaleAlert]):
        """Broadcast whale alerts to configured channels."""
        if not alerts:
            return
        
        for guild_id, config in self.guild_configs.items():
            if not config.get('whale_alerts_enabled', False):
                continue
                
            channel_id = config.get('whale_channel_id')
            if not channel_id:
                continue
            
            channel = self.bot.get_channel(channel_id)
            if not channel:
                continue
            
            # Apply rate limiting
            last_alert_time = self.last_whale_alert.get(guild_id, datetime.min)
            if (datetime.utcnow() - last_alert_time).total_seconds() < self.min_whale_interval:
                continue
            
            # Filter and sort alerts
            filtered_alerts = self._filter_whale_alerts(alerts, guild_id)
            
            if filtered_alerts:
                await self._send_whale_alerts(channel, filtered_alerts)
                self.last_whale_alert[guild_id] = datetime.utcnow()
    
    async def broadcast_guru_tweets(self, tweets: List[GuruTweet]):
        """Broadcast guru tweets to configured channels."""
        if not tweets:
            return
        
        for guild_id, config in self.guild_configs.items():
            if not config.get('guru_tweets_enabled', False):
                continue
                
            channel_id = config.get('tweets_channel_id')
            if not channel_id:
                continue
            
            channel = self.bot.get_channel(channel_id)
            if not channel:
                continue
            
            # Apply rate limiting
            last_tweet_time = self.last_guru_tweet.get(guild_id, datetime.min)
            if (datetime.utcnow() - last_tweet_time).total_seconds() < self.min_tweet_interval:
                continue
            
            # Filter and sort tweets
            filtered_tweets = self._filter_guru_tweets(tweets, guild_id)
            
            if filtered_tweets:
                await self._send_guru_tweets(channel, filtered_tweets)
                self.last_guru_tweet[guild_id] = datetime.utcnow()
    
    def _filter_whale_alerts(self, alerts: List[WhaleAlert], guild_id: int) -> List[WhaleAlert]:
        """Filter whale alerts based on guild configuration."""
        config = self.guild_configs.get(guild_id, {})
        filters = config.get('whale_filters', self.whale_filters)
        
        filtered = []
        
        for alert in alerts:
            # Check minimum USD amount
            if alert.amount_usd < filters.get('min_usd_amount', 100_000):
                continue
            
            # Priority tokens get lower threshold
            if alert.token_symbol in filters.get('priority_tokens', []):
                if alert.amount_usd < filters.get('min_usd_amount', 100_000) * 0.5:
                    continue
            
            # High priority transactions always pass
            if alert.amount_usd >= filters.get('high_priority_amount', 1_000_000):
                filtered.append(alert)
                continue
            
            # Regular filtering
            filtered.append(alert)
        
        # Sort by USD amount (highest first)
        filtered.sort(key=lambda x: x.amount_usd, reverse=True)
        
        # Limit to top 3 alerts to prevent spam
        return filtered[:3]
    
    def _filter_guru_tweets(self, tweets: List[GuruTweet], guild_id: int) -> List[GuruTweet]:
        """Filter guru tweets based on guild configuration."""
        config = self.guild_configs.get(guild_id, {})
        filters = config.get('tweet_filters', self.tweet_filters)
        
        filtered = []
        
        for tweet in tweets:
            # Check minimum engagement
            if tweet.get_engagement_score() < filters.get('min_engagement', 100):
                # Exception for priority accounts
                if tweet.username not in filters.get('priority_accounts', []):
                    continue
            
            # If urgent_only is enabled, only show urgent tweets
            if filters.get('urgent_only', False):
                if tweet.sentiment.value != 'urgent':
                    continue
            
            # Priority accounts get preference
            if tweet.username in filters.get('priority_accounts', []):
                filtered.insert(0, tweet)  # Add to front
            else:
                filtered.append(tweet)
        
        # Limit to top 2 tweets to prevent spam
        return filtered[:2]
    
    async def _send_whale_alerts(self, channel: discord.TextChannel, alerts: List[WhaleAlert]):
        """Send whale alerts to a Discord channel."""
        try:
            for alert in alerts:
                embed_data = alert.to_embed()
                embed = discord.Embed.from_dict(embed_data)
                
                # Add reaction based on severity
                message = await channel.send(embed=embed)
                
                if alert.amount_usd >= 10_000_000:  # $10M+
                    await message.add_reaction("🚨")
                elif alert.amount_usd >= 1_000_000:  # $1M+
                    await message.add_reaction("🐋")
                else:
                    await message.add_reaction("💰")
                
                # Small delay between messages
                await asyncio.sleep(1)
                
        except Exception as e:
            self.logger.error(f"Error sending whale alerts: {e}")
    
    async def _send_guru_tweets(self, channel: discord.TextChannel, tweets: List[GuruTweet]):
        """Send guru tweets to a Discord channel."""
        try:
            for tweet in tweets:
                embed_data = tweet.to_embed()
                embed = discord.Embed.from_dict(embed_data)
                
                # Add reaction based on sentiment
                message = await channel.send(embed=embed)
                
                if tweet.sentiment.value == 'urgent':
                    await message.add_reaction("🚨")
                elif tweet.sentiment.value == 'bullish':
                    await message.add_reaction("🚀")
                elif tweet.sentiment.value == 'bearish':
                    await message.add_reaction("📉")
                else:
                    await message.add_reaction("💭")
                
                # Small delay between messages
                await asyncio.sleep(1)
                
        except Exception as e:
            self.logger.error(f"Error sending guru tweets: {e}")
    
    def configure_guild(self, guild_id: int, config: Dict[str, Any]):
        """Configure alerts for a guild."""
        self.guild_configs[guild_id] = config
        self.logger.info(f"Configured alerts for guild {guild_id}")
    
    def get_guild_config(self, guild_id: int) -> Dict[str, Any]:
        """Get configuration for a guild."""
        return self.guild_configs.get(guild_id, {})
    
    def enable_whale_alerts(self, guild_id: int, channel_id: int):
        """Enable whale alerts for a guild."""
        if guild_id not in self.guild_configs:
            self.guild_configs[guild_id] = {}
        
        self.guild_configs[guild_id].update({
            'whale_alerts_enabled': True,
            'whale_channel_id': channel_id
        })
    
    def enable_guru_tweets(self, guild_id: int, channel_id: int):
        """Enable guru tweets for a guild."""
        if guild_id not in self.guild_configs:
            self.guild_configs[guild_id] = {}
        
        self.guild_configs[guild_id].update({
            'guru_tweets_enabled': True,
            'tweets_channel_id': channel_id
        })
    
    def disable_whale_alerts(self, guild_id: int):
        """Disable whale alerts for a guild."""
        if guild_id in self.guild_configs:
            self.guild_configs[guild_id]['whale_alerts_enabled'] = False
    
    def disable_guru_tweets(self, guild_id: int):
        """Disable guru tweets for a guild."""
        if guild_id in self.guild_configs:
            self.guild_configs[guild_id]['guru_tweets_enabled'] = False
    
    def update_whale_filters(self, guild_id: int, filters: Dict[str, Any]):
        """Update whale alert filters for a guild."""
        if guild_id not in self.guild_configs:
            self.guild_configs[guild_id] = {}
        
        self.guild_configs[guild_id]['whale_filters'] = filters
    
    def update_tweet_filters(self, guild_id: int, filters: Dict[str, Any]):
        """Update tweet filters for a guild."""
        if guild_id not in self.guild_configs:
            self.guild_configs[guild_id] = {}
        
        self.guild_configs[guild_id]['tweet_filters'] = filters
    
    def get_status(self, guild_id: int) -> Dict[str, Any]:
        """Get status of alerts for a guild."""
        config = self.guild_configs.get(guild_id, {})
        
        return {
            'whale_alerts_enabled': config.get('whale_alerts_enabled', False),
            'guru_tweets_enabled': config.get('guru_tweets_enabled', False),
            'whale_channel_id': config.get('whale_channel_id'),
            'tweets_channel_id': config.get('tweets_channel_id'),
            'last_whale_alert': self.last_whale_alert.get(guild_id),
            'last_guru_tweet': self.last_guru_tweet.get(guild_id),
        }
