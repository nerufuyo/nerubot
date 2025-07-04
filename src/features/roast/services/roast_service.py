"""
Roast service for tracking user behavior and generating personalized roasts
"""
import asyncio
import json
import os
import random
import re
import time
from collections import Counter, defaultdict
from typing import Dict, List, Optional, Tuple, Set
from datetime import datetime, timedelta

import discord
from discord.ext import tasks

from src.features.roast.models.roast_models import (
    UserActivity, ActivityType, UserBehaviorPattern, 
    RoastTemplate, UserRoastProfile, RoastStats
)
from src.services.ai_service import ai_service
from src.core.utils.logging_utils import get_logger

logger = get_logger(__name__)


class RoastService:
    """Service for tracking user behavior and generating epic roasts"""
    
    def __init__(self):
        self.user_activities: Dict[int, List[UserActivity]] = {}  # user_id -> activities
        self.behavior_patterns: Dict[int, UserBehaviorPattern] = {}  # user_id -> patterns
        self.roast_profiles: Dict[int, UserRoastProfile] = {}  # user_id -> roast profile
        self.roast_stats = RoastStats()
        
        # Data persistence
        self.data_dir = "data/roasts"
        os.makedirs(self.data_dir, exist_ok=True)
        self.activities_file = f"{self.data_dir}/activities.json"
        self.patterns_file = f"{self.data_dir}/patterns.json"
        self.profiles_file = f"{self.data_dir}/profiles.json"
        self.stats_file = f"{self.data_dir}/stats.json"
        
        # Analysis settings
        self.max_activities_per_user = 1000  # Keep last 1000 activities per user
        self.analysis_interval_hours = 6  # Re-analyze patterns every 6 hours
        
        # Load existing data
        self._load_data()
        
        # Roast templates - categorized by behavior patterns
        self.roast_templates = self._initialize_roast_templates()
        
        # Track active voice channels for real-time monitoring
        self.voice_tracking: Dict[int, Set[int]] = {}  # guild_id -> set of user_ids
    
    def _initialize_roast_templates(self) -> Dict[str, List[RoastTemplate]]:
        """Initialize roast templates for different behavioral patterns"""
        return {
            "night_owl": [
                RoastTemplate(
                    category="night_owl",
                    template="Bruh, you're online at {hour}? Do you even know what sunlight looks like? 🌙😴",
                    severity=2,
                    tags=["sleep", "nocturnal"]
                ),
                RoastTemplate(
                    category="night_owl", 
                    template="Imagine having a normal sleep schedule... couldn't be you at {hour} AM! 🦉💀",
                    severity=3,
                    tags=["sleep", "schedule"]
                ),
                RoastTemplate(
                    category="night_owl",
                    template="The sun: exists\nYou: 'Nah, I'll take the 3 AM shift' 🌃🤡",
                    severity=4,
                    tags=["sleep", "vampire"]
                )
            ],
            "spammer": [
                RoastTemplate(
                    category="spammer",
                    template="Bro sent {message_count} messages in {timeframe}... somebody forgot their indoor voice 📢🤐",
                    severity=2,
                    tags=["spam", "chatty"]
                ),
                RoastTemplate(
                    category="spammer",
                    template="You type more than a middle schooler writing their first essay... and that's saying something 📝💀",
                    severity=3,
                    tags=["spam", "typing"]
                ),
                RoastTemplate(
                    category="spammer",
                    template="Discord notifications when you're online: 📳📳📳📳📳📳📳 RIP everyone's sanity",
                    severity=4,
                    tags=["spam", "notifications"]
                )
            ],
            "music_addict": [
                RoastTemplate(
                    category="music_addict",
                    template="You've played {song_count} songs... are you trying to become a human jukebox? 🎵🤖",
                    severity=2,
                    tags=["music", "addiction"]
                ),
                RoastTemplate(
                    category="music_addict",
                    template="Your Spotify Wrapped probably crashes from data overflow 💀🎶",
                    severity=3,
                    tags=["music", "spotify"]
                ),
                RoastTemplate(
                    category="music_addict",
                    template="You treat the music bot like your personal radio station... this isn't your bedroom, chief 🎧👑",
                    severity=4,
                    tags=["music", "selfish"]
                )
            ],
            "lurker": [
                RoastTemplate(
                    category="lurker",
                    template="You lurk more than a serial killer in a horror movie... just say hi sometimes? 👻😭",
                    severity=2,
                    tags=["lurking", "antisocial"]
                ),
                RoastTemplate(
                    category="lurker",
                    template="Online for hours but {message_count} messages... are you collecting data for the FBI? 🕵️‍♂️📊",
                    severity=3,
                    tags=["lurking", "suspicious"]
                ),
                RoastTemplate(
                    category="lurker",
                    template="Your message history is shorter than a TikTok attention span 📱💀",
                    severity=4,
                    tags=["lurking", "silent"]
                )
            ],
            "confessor": [
                RoastTemplate(
                    category="confessor",
                    template="Anonymous confessions: {confession_count}\nEveryone knowing it's you anyway: Priceless 🎭💀",
                    severity=3,
                    tags=["confession", "anonymous"]
                ),
                RoastTemplate(
                    category="confessor",
                    template="You use the confession bot more than a Catholic at Sunday mass 🙏💒",
                    severity=2,
                    tags=["confession", "religious"]
                ),
                RoastTemplate(
                    category="confessor",
                    template="Plot twist: your 'anonymous' confessions have more character development than Netflix shows 📺🎪",
                    severity=4,
                    tags=["confession", "drama"]
                )
            ],
            "emoji_spammer": [
                RoastTemplate(
                    category="emoji_spammer",
                    template="Uses {emoji_count} emojis per message... we get it, you have feelings 😂😭🤡💀✨🎉",
                    severity=2,
                    tags=["emoji", "expressive"]
                ),
                RoastTemplate(
                    category="emoji_spammer",
                    template="Your keyboard must be 90% emoji at this point 🌈🎨🎭",
                    severity=3,
                    tags=["emoji", "overuse"]
                ),
                RoastTemplate(
                    category="emoji_spammer",
                    template="You communicate like a Gen Z hieroglyphic translator 🏺📜😤",
                    severity=4,
                    tags=["emoji", "hieroglyphics"]
                )
            ],
            "early_bird": [
                RoastTemplate(
                    category="early_bird", 
                    template="Online at {hour} AM? What are you, a rooster with WiFi? 🐓📶",
                    severity=2,
                    tags=["morning", "early"]
                ),
                RoastTemplate(
                    category="early_bird",
                    template="5 AM messages... someone's productivity coach is working overtime 📈☕",
                    severity=3,
                    tags=["morning", "productive"]
                )
            ],
            "weekend_warrior": [
                RoastTemplate(
                    category="weekend_warrior",
                    template="Only shows up on weekends... Discord's part-time employee over here 📅💼",
                    severity=2,
                    tags=["weekend", "casual"]
                ),
                RoastTemplate(
                    category="weekend_warrior",
                    template="Your Discord activity schedule is more predictable than a 9-5 job 🕘📊",
                    severity=3,
                    tags=["weekend", "predictable"]
                )
            ],
            "voice_hopper": [
                RoastTemplate(
                    category="voice_hopper",
                    template="Joined and left VC {hop_count} times... make up your mind! 🚪🏃‍♂️",
                    severity=2,
                    tags=["voice", "indecisive"]
                ),
                RoastTemplate(
                    category="voice_hopper",
                    template="Your VC commitment issues are showing 💔🎤",
                    severity=3,
                    tags=["voice", "commitment"]
                )
            ]
        }
    
    async def track_activity(self, user_id: int, activity_type: ActivityType, 
                           channel_id: int, guild_id: int, content: str = None, 
                           metadata: Dict = None) -> None:
        """Track a user activity"""
        try:
            activity = UserActivity(
                user_id=user_id,
                activity_type=activity_type,
                timestamp=time.time(),
                channel_id=channel_id,
                guild_id=guild_id,
                content=content,
                metadata=metadata or {}
            )
            
            # Add activity to user's history
            if user_id not in self.user_activities:
                self.user_activities[user_id] = []
            
            self.user_activities[user_id].append(activity)
            
            # Trim activities if too many
            if len(self.user_activities[user_id]) > self.max_activities_per_user:
                self.user_activities[user_id] = self.user_activities[user_id][-self.max_activities_per_user:]
            
            # Special handling for late night/early morning detection
            if activity.is_late_night():
                await self.track_activity(user_id, ActivityType.LATE_NIGHT, channel_id, guild_id)
            elif activity.is_early_morning():
                await self.track_activity(user_id, ActivityType.EARLY_MORNING, channel_id, guild_id)
            
            logger.debug(f"Tracked {activity_type.value} activity for user {user_id}")
            
        except Exception as e:
            logger.error(f"Error tracking activity for user {user_id}: {e}")
    
    async def analyze_user_behavior(self, user_id: int) -> UserBehaviorPattern:
        """Analyze a user's behavior patterns"""
        try:
            activities = self.user_activities.get(user_id, [])
            if not activities:
                return UserBehaviorPattern(user_id=user_id)
            
            # Analyze activity hours
            hour_counter = Counter(activity.hour_of_day for activity in activities)
            most_active_hours = [hour for hour, _ in hour_counter.most_common(5)]
            
            # Analyze active days
            day_counter = Counter(activity.day_of_week for activity in activities)
            most_active_days = [day for day, _ in day_counter.most_common(3)]
            
            # Analyze favorite channels
            channel_counter = Counter(activity.channel_id for activity in activities)
            favorite_channels = [channel for channel, _ in channel_counter.most_common(3)]
            
            # Analyze activity frequency
            activity_frequency = Counter(activity.activity_type for activity in activities)
            
            # Calculate late night percentage
            late_night_count = len([a for a in activities if a.is_late_night()])
            late_night_percentage = (late_night_count / len(activities)) * 100 if activities else 0
            
            # Analyze message content
            messages = [a.content for a in activities if a.content and a.activity_type == ActivityType.MESSAGE]
            common_words = []
            message_length_avg = 0.0
            emoji_usage = Counter()
            
            if messages:
                # Extract common words (filter out common words)
                all_words = []
                total_length = 0
                
                for msg in messages:
                    # Extract emojis
                    emoji_pattern = r'[\U0001F600-\U0001F64F\U0001F300-\U0001F5FF\U0001F680-\U0001F6FF\U0001F1E0-\U0001F1FF\U00002702-\U000027B0\U000024C2-\U0001F251]+'
                    emojis = re.findall(emoji_pattern, msg)
                    for emoji in emojis:
                        emoji_usage[emoji] += 1
                    
                    # Extract words
                    words = re.findall(r'\b\w+\b', msg.lower())
                    all_words.extend(words)
                    total_length += len(msg)
                
                # Filter common words and get meaningful ones
                common_stop_words = {'the', 'a', 'an', 'and', 'or', 'but', 'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by', 'is', 'are', 'was', 'were', 'be', 'been', 'have', 'has', 'had', 'do', 'does', 'did', 'will', 'would', 'could', 'should', 'may', 'might', 'can', 'i', 'you', 'he', 'she', 'it', 'we', 'they', 'me', 'him', 'her', 'us', 'them', 'my', 'your', 'his', 'her', 'its', 'our', 'their'}
                
                word_counter = Counter(word for word in all_words if word not in common_stop_words and len(word) > 2)
                common_words = [word for word, _ in word_counter.most_common(10)]
                message_length_avg = total_length / len(messages)
            
            pattern = UserBehaviorPattern(
                user_id=user_id,
                most_active_hours=most_active_hours,
                most_active_days=most_active_days,
                favorite_channels=favorite_channels,
                common_words=common_words,
                activity_frequency=dict(activity_frequency),
                late_night_percentage=late_night_percentage,
                message_length_avg=message_length_avg,
                emoji_usage=dict(emoji_usage),
                last_analyzed=time.time()
            )
            
            self.behavior_patterns[user_id] = pattern
            logger.info(f"Analyzed behavior patterns for user {user_id}")
            
            return pattern
            
        except Exception as e:
            logger.error(f"Error analyzing behavior for user {user_id}: {e}")
            return UserBehaviorPattern(user_id=user_id)
    
    async def generate_roast(self, user_id: int, custom_prompt: str = None) -> Tuple[str, str]:
        """
        Generate a personalized roast for a user
        
        Returns:
            Tuple of (roast_text, roast_category)
        """
        try:
            # Check cooldown
            profile = self.roast_profiles.get(user_id, UserRoastProfile(user_id=user_id))
            if not profile.can_be_roasted():
                hours_left = 6 - ((time.time() - profile.last_roasted) / 3600)
                return f"Hold up! You were just roasted {hours_left:.1f} hours ago. Let your ego recover first! 🛡️😤", "cooldown"
            
            # Analyze user behavior if needed
            pattern = self.behavior_patterns.get(user_id)
            if not pattern or (time.time() - pattern.last_analyzed) > (self.analysis_interval_hours * 3600):
                pattern = await self.analyze_user_behavior(user_id)
            
            # Determine roast category based on behavior
            roast_category, template = self._select_roast_template(pattern)
            
            # If custom prompt provided, use AI to generate custom roast
            if custom_prompt:
                roast_text = await self._generate_ai_roast(user_id, pattern, custom_prompt)
                roast_category = "custom"
            else:
                # Use template-based roast with dynamic data
                roast_text = self._format_roast_template(template, pattern)
            
            # Update roast profile
            if user_id not in self.roast_profiles:
                self.roast_profiles[user_id] = profile
            
            self.roast_profiles[user_id].add_roast(roast_text)
            
            # Update stats
            self.roast_stats.total_roasts_delivered += 1
            self.roast_stats.daily_roast_count += 1
            self.roast_stats.roasts_by_category[roast_category] = self.roast_stats.roasts_by_category.get(roast_category, 0) + 1
            
            # Reset daily stats if needed
            if self.roast_stats.should_reset_daily():
                self.roast_stats.reset_daily_stats()
            
            logger.info(f"Generated {roast_category} roast for user {user_id}")
            return roast_text, roast_category
            
        except Exception as e:
            logger.error(f"Error generating roast for user {user_id}: {e}")
            return "My roast generator is having a moment... unlike your social life! 🤖💀", "error"
    
    def _select_roast_template(self, pattern: UserBehaviorPattern) -> Tuple[str, RoastTemplate]:
        """Select appropriate roast template based on user behavior"""
        try:
            activity_freq = pattern.activity_frequency
            
            # Priority-based selection
            if pattern.late_night_percentage > 30:
                category = "night_owl"
            elif activity_freq.get(ActivityType.MESSAGE, 0) > 50:
                category = "spammer"
            elif activity_freq.get(ActivityType.MUSIC_REQUEST, 0) > 20:
                category = "music_addict"
            elif activity_freq.get(ActivityType.CONFESSION, 0) > 5:
                category = "confessor"
            elif len(pattern.emoji_usage) > 20:
                category = "emoji_spammer"
            elif any(hour < 8 for hour in pattern.most_active_hours):
                category = "early_bird"
            elif all(day in [5, 6] for day in pattern.most_active_days):  # Saturday, Sunday
                category = "weekend_warrior"
            elif activity_freq.get(ActivityType.VOICE_JOIN, 0) > activity_freq.get(ActivityType.MESSAGE, 0):
                category = "voice_hopper"
            elif activity_freq.get(ActivityType.MESSAGE, 0) < 10:
                category = "lurker"
            else:
                # Random category if no clear pattern
                category = random.choice(list(self.roast_templates.keys()))
            
            # Select random template from category
            templates = self.roast_templates.get(category, self.roast_templates["lurker"])
            template = random.choice(templates)
            
            return category, template
            
        except Exception as e:
            logger.error(f"Error selecting roast template: {e}")
            # Fallback
            return "lurker", self.roast_templates["lurker"][0]
    
    def _format_roast_template(self, template: RoastTemplate, pattern: UserBehaviorPattern) -> str:
        """Format roast template with user-specific data"""
        try:
            format_data = {
                'hour': random.choice(pattern.most_active_hours) if pattern.most_active_hours else 3,
                'message_count': pattern.activity_frequency.get(ActivityType.MESSAGE, 0),
                'song_count': pattern.activity_frequency.get(ActivityType.MUSIC_REQUEST, 0),
                'confession_count': pattern.activity_frequency.get(ActivityType.CONFESSION, 0),
                'emoji_count': len(pattern.emoji_usage),
                'hop_count': pattern.activity_frequency.get(ActivityType.VOICE_JOIN, 0),
                'timeframe': "this week"
            }
            
            return template.template.format(**format_data)
            
        except Exception as e:
            logger.error(f"Error formatting roast template: {e}")
            return template.template  # Return unformatted if error
    
    async def _generate_ai_roast(self, user_id: int, pattern: UserBehaviorPattern, custom_prompt: str) -> str:
        """Generate AI-powered custom roast"""
        try:
            # Build behavior summary for AI
            behavior_summary = self._build_behavior_summary(pattern)
            
            ai_prompt = f"""You are a master roaster who creates HILARIOUS and SAVAGE but friendly roasts. 
Generate a funny roast based on this user's Discord behavior:

User Behavior Summary:
{behavior_summary}

Custom roast request: {custom_prompt}

Create a roast that is:
- Hilariously savage but not genuinely mean
- Uses internet/Discord culture references
- Includes relevant emojis
- Maximum 200 characters
- Creative and unexpected

Roast:"""

            response = await ai_service.chat(
                message=ai_prompt,
                max_tokens=100,
                temperature=0.9  # High creativity
            )
            
            return response.content.strip()
            
        except Exception as e:
            logger.error(f"Error generating AI roast: {e}")
            return f"My AI brain short-circuited trying to roast you... that's saying something! 🤖⚡"
    
    def _build_behavior_summary(self, pattern: UserBehaviorPattern) -> str:
        """Build a summary of user behavior for AI roasting"""
        summary_parts = []
        
        if pattern.late_night_percentage > 20:
            summary_parts.append(f"Night owl ({pattern.late_night_percentage:.0f}% late night activity)")
        
        if pattern.activity_frequency.get(ActivityType.MESSAGE, 0) > 30:
            summary_parts.append(f"Chatty ({pattern.activity_frequency[ActivityType.MESSAGE]} messages)")
        
        if pattern.activity_frequency.get(ActivityType.MUSIC_REQUEST, 0) > 10:
            summary_parts.append(f"Music lover ({pattern.activity_frequency[ActivityType.MUSIC_REQUEST]} song requests)")
        
        if len(pattern.emoji_usage) > 10:
            summary_parts.append(f"Emoji enthusiast ({len(pattern.emoji_usage)} different emojis)")
        
        if pattern.most_active_hours:
            hours_str = ", ".join(f"{h}:00" for h in pattern.most_active_hours[:3])
            summary_parts.append(f"Most active at {hours_str}")
        
        if pattern.common_words:
            words_str = ", ".join(pattern.common_words[:5])
            summary_parts.append(f"Common words: {words_str}")
        
        return "; ".join(summary_parts) if summary_parts else "Limited activity data"
    
    async def get_user_roast_stats(self, user_id: int) -> Dict:
        """Get roast statistics for a user"""
        profile = self.roast_profiles.get(user_id, UserRoastProfile(user_id=user_id))
        pattern = self.behavior_patterns.get(user_id)
        
        return {
            "roast_count": profile.roast_count,
            "last_roasted": profile.last_roasted,
            "can_be_roasted": profile.can_be_roasted(),
            "behavior_analyzed": pattern is not None,
            "activity_count": len(self.user_activities.get(user_id, [])),
            "immunity_level": profile.immunity_level
        }
    
    async def get_global_roast_stats(self) -> RoastStats:
        """Get global roast statistics"""
        return self.roast_stats
    
    def _load_data(self):
        """Load roast data from files"""
        try:
            # Load activities
            if os.path.exists(self.activities_file):
                with open(self.activities_file, 'r') as f:
                    data = json.load(f)
                    for user_id_str, activities_data in data.items():
                        user_id = int(user_id_str)
                        activities = []
                        for activity_data in activities_data:
                            activity = UserActivity(
                                user_id=activity_data['user_id'],
                                activity_type=ActivityType(activity_data['activity_type']),
                                timestamp=activity_data['timestamp'],
                                channel_id=activity_data['channel_id'],
                                guild_id=activity_data['guild_id'],
                                content=activity_data.get('content'),
                                metadata=activity_data.get('metadata', {})
                            )
                            activities.append(activity)
                        self.user_activities[user_id] = activities
            
            # Load patterns
            if os.path.exists(self.patterns_file):
                with open(self.patterns_file, 'r') as f:
                    data = json.load(f)
                    for user_id_str, pattern_data in data.items():
                        user_id = int(user_id_str)
                        # Convert activity_frequency keys back to ActivityType enums
                        activity_frequency = {}
                        for activity_str, count in pattern_data.get('activity_frequency', {}).items():
                            try:
                                activity_frequency[ActivityType(activity_str)] = count
                            except ValueError:
                                continue  # Skip invalid activity types
                        
                        pattern = UserBehaviorPattern(
                            user_id=user_id,
                            most_active_hours=pattern_data.get('most_active_hours', []),
                            most_active_days=pattern_data.get('most_active_days', []),
                            favorite_channels=pattern_data.get('favorite_channels', []),
                            common_words=pattern_data.get('common_words', []),
                            activity_frequency=activity_frequency,
                            late_night_percentage=pattern_data.get('late_night_percentage', 0.0),
                            message_length_avg=pattern_data.get('message_length_avg', 0.0),
                            emoji_usage=pattern_data.get('emoji_usage', {}),
                            last_analyzed=pattern_data.get('last_analyzed', 0.0)
                        )
                        self.behavior_patterns[user_id] = pattern
            
            # Load profiles
            if os.path.exists(self.profiles_file):
                with open(self.profiles_file, 'r') as f:
                    data = json.load(f)
                    for user_id_str, profile_data in data.items():
                        user_id = int(user_id_str)
                        profile = UserRoastProfile(
                            user_id=user_id,
                            personality_traits=profile_data.get('personality_traits', []),
                            behavior_summary=profile_data.get('behavior_summary', ''),
                            roast_history=profile_data.get('roast_history', []),
                            last_roasted=profile_data.get('last_roasted', 0.0),
                            roast_count=profile_data.get('roast_count', 0),
                            immunity_level=profile_data.get('immunity_level', 0)
                        )
                        self.roast_profiles[user_id] = profile
            
            # Load stats
            if os.path.exists(self.stats_file):
                with open(self.stats_file, 'r') as f:
                    data = json.load(f)
                    self.roast_stats = RoastStats(
                        total_roasts_delivered=data.get('total_roasts_delivered', 0),
                        most_roasted_user=data.get('most_roasted_user'),
                        best_roast_rating=data.get('best_roast_rating', 0.0),
                        roasts_by_category=data.get('roasts_by_category', {}),
                        daily_roast_count=data.get('daily_roast_count', 0),
                        last_reset=data.get('last_reset', 0.0)
                    )
            
            logger.info("Loaded roast data successfully")
            
        except Exception as e:
            logger.error(f"Error loading roast data: {e}")
    
    def _save_data(self):
        """Save roast data to files"""
        try:
            # Save activities
            activities_data = {}
            for user_id, activities in self.user_activities.items():
                activities_data[str(user_id)] = [
                    {
                        'user_id': activity.user_id,
                        'activity_type': activity.activity_type.value,
                        'timestamp': activity.timestamp,
                        'channel_id': activity.channel_id,
                        'guild_id': activity.guild_id,
                        'content': activity.content,
                        'metadata': activity.metadata
                    }
                    for activity in activities
                ]
            
            with open(self.activities_file, 'w') as f:
                json.dump(activities_data, f, indent=2)
            
            # Save patterns
            patterns_data = {}
            for user_id, pattern in self.behavior_patterns.items():
                # Convert ActivityType enums to strings for JSON serialization
                activity_frequency = {
                    activity_type.value: count 
                    for activity_type, count in pattern.activity_frequency.items()
                }
                
                patterns_data[str(user_id)] = {
                    'user_id': pattern.user_id,
                    'most_active_hours': pattern.most_active_hours,
                    'most_active_days': pattern.most_active_days,
                    'favorite_channels': pattern.favorite_channels,
                    'common_words': pattern.common_words,
                    'activity_frequency': activity_frequency,
                    'late_night_percentage': pattern.late_night_percentage,
                    'message_length_avg': pattern.message_length_avg,
                    'emoji_usage': pattern.emoji_usage,
                    'last_analyzed': pattern.last_analyzed
                }
            
            with open(self.patterns_file, 'w') as f:
                json.dump(patterns_data, f, indent=2)
            
            # Save profiles
            profiles_data = {}
            for user_id, profile in self.roast_profiles.items():
                profiles_data[str(user_id)] = {
                    'user_id': profile.user_id,
                    'personality_traits': profile.personality_traits,
                    'behavior_summary': profile.behavior_summary,
                    'roast_history': profile.roast_history,
                    'last_roasted': profile.last_roasted,
                    'roast_count': profile.roast_count,
                    'immunity_level': profile.immunity_level
                }
            
            with open(self.profiles_file, 'w') as f:
                json.dump(profiles_data, f, indent=2)
            
            # Save stats
            stats_data = {
                'total_roasts_delivered': self.roast_stats.total_roasts_delivered,
                'most_roasted_user': self.roast_stats.most_roasted_user,
                'best_roast_rating': self.roast_stats.best_roast_rating,
                'roasts_by_category': self.roast_stats.roasts_by_category,
                'daily_roast_count': self.roast_stats.daily_roast_count,
                'last_reset': self.roast_stats.last_reset
            }
            
            with open(self.stats_file, 'w') as f:
                json.dump(stats_data, f, indent=2)
            
            logger.debug("Saved roast data successfully")
            
        except Exception as e:
            logger.error(f"Error saving roast data: {e}")
    
    @tasks.loop(minutes=30)
    async def auto_save_task(self):
        """Automatically save data every 30 minutes"""
        try:
            self._save_data()
            logger.debug("Auto-saved roast data")
        except Exception as e:
            logger.error(f"Error in auto-save task: {e}")
    
    async def start_auto_save(self):
        """Start the auto-save task"""
        if not self.auto_save_task.is_running():
            self.auto_save_task.start()
            logger.info("Started roast data auto-save task")
    
    async def stop_auto_save(self):
        """Stop the auto-save task"""
        if self.auto_save_task.is_running():
            self.auto_save_task.cancel()
            # Save one final time
            self._save_data()
            logger.info("Stopped roast data auto-save task")


# Create singleton instance
roast_service = RoastService()
