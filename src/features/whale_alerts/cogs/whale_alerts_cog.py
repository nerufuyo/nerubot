"""Discord cog for whale alerts and guru tweets functionality."""
import asyncio
import logging
import os
from datetime import datetime
from typing import Optional, Union

import discord
from discord.ext import commands, tasks
from discord import app_commands

from ..services.twitter_monitor import TwitterMonitor
from ..services.whale_service import WhaleService
from ..services.alerts_broadcaster import AlertsBroadcaster
from ..models.whale_alert import AlertType
from ..models.guru_tweet import TweetSentiment


class WhaleAlertsCog(commands.Cog):
    """Cog for whale alerts and crypto guru tweets monitoring."""
    
    def __init__(self, bot):
        """Initialize the whale alerts cog."""
        self.bot = bot
        self.logger = logging.getLogger(__name__)
        
        # Initialize services
        self.twitter_monitor = None
        self.whale_service = None
        self.broadcaster = AlertsBroadcaster(bot)
        
        # Initialize API services if credentials are available
        self._init_services()
        
        # Start monitoring tasks
        self.monitoring_task = None
        
    def _init_services(self):
        """Initialize services with API credentials from environment."""
        try:
            # Twitter API credentials
            twitter_api_key = os.getenv('TWITTER_API_KEY')
            twitter_api_secret = os.getenv('TWITTER_API_SECRET')
            twitter_access_token = os.getenv('TWITTER_ACCESS_TOKEN')
            twitter_access_secret = os.getenv('TWITTER_ACCESS_TOKEN_SECRET')
            
            if all([twitter_api_key, twitter_api_secret, twitter_access_token, twitter_access_secret]):
                self.twitter_monitor = TwitterMonitor(
                    twitter_api_key, twitter_api_secret,
                    twitter_access_token, twitter_access_secret
                )
                self.logger.info("Twitter monitor initialized")
            else:
                self.logger.warning("Twitter API credentials not found, Twitter monitoring disabled")
            
            # Whale Alert API
            whale_alert_key = os.getenv('WHALE_ALERT_API_KEY')
            self.whale_service = WhaleService(whale_alert_key)
            self.logger.info("Whale service initialized")
            
        except Exception as e:
            self.logger.error(f"Error initializing services: {e}")
    
    async def cog_load(self):
        """Called when the cog is loaded."""
        self.logger.info("Whale alerts cog loaded")
        # Start monitoring if services are available
        if self.twitter_monitor or self.whale_service:
            self.start_monitoring()
    
    async def cog_unload(self):
        """Called when the cog is unloaded."""
        self.stop_monitoring()
        self.logger.info("Whale alerts cog unloaded")
    
    def start_monitoring(self):
        """Start monitoring whale alerts and guru tweets."""
        if self.monitoring_task and not self.monitoring_task.done():
            return
        
        self.monitoring_task = asyncio.create_task(self._monitoring_loop())
        self.logger.info("Started whale alerts monitoring")
    
    def stop_monitoring(self):
        """Stop monitoring."""
        if self.monitoring_task and not self.monitoring_task.done():
            self.monitoring_task.cancel()
        
        if self.twitter_monitor:
            self.twitter_monitor.stop_monitoring()
        
        if self.whale_service:
            self.whale_service.stop_monitoring()
        
        self.logger.info("Stopped whale alerts monitoring")
    
    async def _monitoring_loop(self):
        """Main monitoring loop."""
        try:
            # Start individual monitoring services
            tasks = []
            
            if self.whale_service:
                tasks.append(
                    asyncio.create_task(
                        self.whale_service.start_monitoring(self.broadcaster.broadcast_whale_alerts)
                    )
                )
            
            if self.twitter_monitor:
                tasks.append(
                    asyncio.create_task(
                        self.twitter_monitor.start_monitoring(self.broadcaster.broadcast_guru_tweets)
                    )
                )
            
            if tasks:
                await asyncio.gather(*tasks, return_exceptions=True)
                
        except asyncio.CancelledError:
            self.logger.info("Monitoring loop cancelled")
        except Exception as e:
            self.logger.error(f"Error in monitoring loop: {e}")
    
    # Whale Alerts Commands Group
    whale_group = app_commands.Group(name="whale", description="Whale alerts commands")
    
    @whale_group.command(name="setup", description="Set up whale alerts in this channel")
    @app_commands.describe(channel="Channel for whale alerts (optional, uses current channel)")
    async def whale_setup(self, interaction: discord.Interaction, channel: Optional[discord.TextChannel] = None):
        """Set up whale alerts in a channel."""
        if not interaction.user.guild_permissions.manage_channels:
            await interaction.response.send_message("❌ You need 'Manage Channels' permission to use this command.", ephemeral=True)
            return
        
        target_channel = channel or interaction.channel
        
        self.broadcaster.enable_whale_alerts(interaction.guild.id, target_channel.id)
        
        embed = discord.Embed(
            title="🐋 Whale Alerts Enabled",
            description=f"Whale alerts will now be posted in {target_channel.mention}",
            color=0x00CCFF
        )
        embed.add_field(
            name="What you'll receive:",
            value="• Large cryptocurrency transfers ($100K+)\n• Exchange deposits/withdrawals ($500K+)\n• DEX swaps ($50K+)\n• Whale accumulation/distribution ($1M+)",
            inline=False
        )
        
        await interaction.response.send_message(embed=embed)
    
    @whale_group.command(name="stop", description="Stop whale alerts in this server")
    async def whale_stop(self, interaction: discord.Interaction):
        """Stop whale alerts for this server."""
        if not interaction.user.guild_permissions.manage_channels:
            await interaction.response.send_message("❌ You need 'Manage Channels' permission to use this command.", ephemeral=True)
            return
        
        self.broadcaster.disable_whale_alerts(interaction.guild.id)
        
        embed = discord.Embed(
            title="🐋 Whale Alerts Disabled",
            description="Whale alerts have been disabled for this server.",
            color=0xFF6600
        )
        
        await interaction.response.send_message(embed=embed)
    
    @whale_group.command(name="recent", description="Get recent whale transactions")
    @app_commands.describe(limit="Number of transactions to show (1-10)")
    async def whale_recent(self, interaction: discord.Interaction, limit: int = 5):
        """Get recent whale transactions."""
        if not 1 <= limit <= 10:
            await interaction.response.send_message("❌ Limit must be between 1 and 10.", ephemeral=True)
            return
        
        await interaction.response.defer()
        
        if not self.whale_service:
            await interaction.followup.send("❌ Whale service is not available.", ephemeral=True)
            return
        
        try:
            alerts = await self.whale_service.get_recent_transactions(limit)
            
            if not alerts:
                embed = discord.Embed(
                    title="🐋 Recent Whale Transactions",
                    description="No recent whale transactions found.",
                    color=0x808080
                )
                await interaction.followup.send(embed=embed)
                return
            
            for i, alert in enumerate(alerts):
                embed_data = alert.to_embed()
                embed = discord.Embed.from_dict(embed_data)
                
                if i == 0:
                    await interaction.followup.send(embed=embed)
                else:
                    await interaction.followup.send(embed=embed)
                
                # Small delay between messages
                await asyncio.sleep(0.5)
                
        except Exception as e:
            self.logger.error(f"Error getting recent whale transactions: {e}")
            await interaction.followup.send("❌ Error retrieving whale transactions.", ephemeral=True)
    
    @whale_group.command(name="status", description="Show whale alerts status")
    async def whale_status(self, interaction: discord.Interaction):
        """Show whale alerts status for this server."""
        status = self.broadcaster.get_status(interaction.guild.id)
        service_stats = self.whale_service.get_statistics() if self.whale_service else {}
        
        embed = discord.Embed(
            title="🐋 Whale Alerts Status",
            color=0x00CCFF if status['whale_alerts_enabled'] else 0x808080
        )
        
        # Status
        status_text = "✅ Enabled" if status['whale_alerts_enabled'] else "❌ Disabled"
        embed.add_field(name="Status", value=status_text, inline=True)
        
        # Channel
        if status['whale_channel_id']:
            channel = self.bot.get_channel(status['whale_channel_id'])
            channel_text = channel.mention if channel else "Unknown Channel"
        else:
            channel_text = "Not configured"
        embed.add_field(name="Channel", value=channel_text, inline=True)
        
        # Service status
        api_status = "✅ API Available" if service_stats.get('api_enabled') else "⚠️ Mock Data"
        embed.add_field(name="Data Source", value=api_status, inline=True)
        
        # Last alert
        last_alert = status.get('last_whale_alert')
        if last_alert:
            last_alert_text = f"<t:{int(last_alert.timestamp())}:R>"
        else:
            last_alert_text = "Never"
        embed.add_field(name="Last Alert", value=last_alert_text, inline=True)
        
        await interaction.response.send_message(embed=embed)
    
    # Guru Tweets Commands Group
    guru_group = app_commands.Group(name="guru", description="Crypto guru tweets commands")
    
    @guru_group.command(name="setup", description="Set up guru tweets in this channel")
    @app_commands.describe(channel="Channel for guru tweets (optional, uses current channel)")
    async def guru_setup(self, interaction: discord.Interaction, channel: Optional[discord.TextChannel] = None):
        """Set up guru tweets in a channel."""
        if not interaction.user.guild_permissions.manage_channels:
            await interaction.response.send_message("❌ You need 'Manage Channels' permission to use this command.", ephemeral=True)
            return
        
        target_channel = channel or interaction.channel
        
        self.broadcaster.enable_guru_tweets(interaction.guild.id, target_channel.id)
        
        embed = discord.Embed(
            title="🧙‍♂️ Guru Tweets Enabled",
            description=f"Crypto guru tweets will now be posted in {target_channel.mention}",
            color=0x1DA1F2
        )
        embed.add_field(
            name="Monitored Accounts:",
            value="• Elon Musk • Michael Saylor • Vitalik Buterin\n• Changpeng Zhao • Brian Armstrong • Jack Dorsey\n• Anthony Pompliano • Naval • Raoul Pal • And more!",
            inline=False
        )
        
        if not self.twitter_monitor:
            embed.add_field(
                name="⚠️ Note:",
                value="Twitter API not configured. Contact admin to enable real-time monitoring.",
                inline=False
            )
        
        await interaction.response.send_message(embed=embed)
    
    @guru_group.command(name="stop", description="Stop guru tweets in this server")
    async def guru_stop(self, interaction: discord.Interaction):
        """Stop guru tweets for this server."""
        if not interaction.user.guild_permissions.manage_channels:
            await interaction.response.send_message("❌ You need 'Manage Channels' permission to use this command.", ephemeral=True)
            return
        
        self.broadcaster.disable_guru_tweets(interaction.guild.id)
        
        embed = discord.Embed(
            title="🧙‍♂️ Guru Tweets Disabled",
            description="Guru tweets have been disabled for this server.",
            color=0xFF6600
        )
        
        await interaction.response.send_message(embed=embed)
    
    @guru_group.command(name="accounts", description="List monitored guru accounts")
    async def guru_accounts(self, interaction: discord.Interaction):
        """List monitored crypto guru accounts."""
        if not self.twitter_monitor:
            embed = discord.Embed(
                title="🧙‍♂️ Monitored Accounts",
                description="Twitter monitoring is not available (API not configured).",
                color=0xFF6600
            )
            await interaction.response.send_message(embed=embed, ephemeral=True)
            return
        
        accounts = self.twitter_monitor.get_monitored_accounts()
        
        embed = discord.Embed(
            title="🧙‍♂️ Monitored Crypto Gurus",
            description="Currently monitoring these accounts for crypto-related tweets:",
            color=0x1DA1F2
        )
        
        # Group by priority
        high_priority = []
        medium_priority = []
        low_priority = []
        
        for username, info in accounts.items():
            account_text = f"@{username} ({info['display_name']})"
            priority = info.get('priority', 'medium')
            
            if priority == 'high':
                high_priority.append(account_text)
            elif priority == 'medium':
                medium_priority.append(account_text)
            else:
                low_priority.append(account_text)
        
        if high_priority:
            embed.add_field(
                name="🔥 High Priority",
                value="\n".join(high_priority),
                inline=False
            )
        
        if medium_priority:
            embed.add_field(
                name="⭐ Medium Priority", 
                value="\n".join(medium_priority),
                inline=False
            )
        
        if low_priority:
            embed.add_field(
                name="📝 Low Priority",
                value="\n".join(low_priority),
                inline=False
            )
        
        await interaction.response.send_message(embed=embed)
    
    @guru_group.command(name="status", description="Show guru tweets status")
    async def guru_status(self, interaction: discord.Interaction):
        """Show guru tweets status for this server."""
        status = self.broadcaster.get_status(interaction.guild.id)
        
        embed = discord.Embed(
            title="🧙‍♂️ Guru Tweets Status",
            color=0x1DA1F2 if status['guru_tweets_enabled'] else 0x808080
        )
        
        # Status
        status_text = "✅ Enabled" if status['guru_tweets_enabled'] else "❌ Disabled"
        embed.add_field(name="Status", value=status_text, inline=True)
        
        # Channel
        if status['tweets_channel_id']:
            channel = self.bot.get_channel(status['tweets_channel_id'])
            channel_text = channel.mention if channel else "Unknown Channel"
        else:
            channel_text = "Not configured"
        embed.add_field(name="Channel", value=channel_text, inline=True)
        
        # Twitter API status
        api_status = "✅ API Available" if self.twitter_monitor else "❌ API Not Configured"
        embed.add_field(name="Twitter API", value=api_status, inline=True)
        
        # Last tweet
        last_tweet = status.get('last_guru_tweet')
        if last_tweet:
            last_tweet_text = f"<t:{int(last_tweet.timestamp())}:R>"
        else:
            last_tweet_text = "Never"
        embed.add_field(name="Last Tweet", value=last_tweet_text, inline=True)
        
        await interaction.response.send_message(embed=embed)


async def setup(bot):
    """Set up the whale alerts cog."""
    await bot.add_cog(WhaleAlertsCog(bot))
