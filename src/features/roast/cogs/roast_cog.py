"""
Roast Discord Cog - Behavior analysis and personalized roasting
"""
import discord
import random
import time
from discord.ext import commands
from discord import app_commands
from typing import Optional

from src.features.roast.services.roast_service import roast_service
from src.features.roast.models.roast_models import ActivityType
from src.core.utils.logging_utils import get_logger

logger = get_logger(__name__)


async def safe_defer_interaction(interaction: discord.Interaction, ephemeral: bool = False) -> bool:
    """
    Safely defer an interaction with proper error handling
    
    Args:
        interaction: Discord interaction to defer
        ephemeral: Whether the response should be ephemeral
        
    Returns:
        True if defer was successful, False if it failed
    """
    try:
        await interaction.response.defer(ephemeral=ephemeral)
        return True
    except discord.NotFound:
        logger.warning(f"Interaction expired for command from {interaction.user}")
        try:
            embed = discord.Embed(
                title="⚠️ Interaction Timeout",
                description="Your command took too long to process. Please try again!",
                color=0xff9500
            )
            await interaction.channel.send(embed=embed)
        except Exception as e:
            logger.error(f"Failed to send timeout message: {e}")
        return False
    except Exception as e:
        logger.error(f"Failed to defer interaction: {e}")
        return False


async def safe_followup_send(interaction: discord.Interaction, content: str = None, embed: discord.Embed = None, ephemeral: bool = False):
    """
    Safely send a followup message with fallback to channel message
    
    Args:
        interaction: Discord interaction
        content: Message content
        embed: Embed to send
        ephemeral: Whether the response should be ephemeral
    """
    try:
        if embed:
            await interaction.followup.send(content=content, embed=embed, ephemeral=ephemeral)
        else:
            await interaction.followup.send(content=content, ephemeral=ephemeral)
    except discord.NotFound:
        logger.warning(f"Followup failed - interaction expired for user {interaction.user}")
        if not ephemeral:
            try:
                mention = f"{interaction.user.mention} " if content else ""
                if embed:
                    await interaction.channel.send(content=mention + (content or ""), embed=embed)
                else:
                    await interaction.channel.send(mention + content)
            except Exception as fallback_error:
                logger.error(f"Fallback message failed: {fallback_error}")


class RoastCog(commands.Cog):
    """Discord cog for roasting functionality"""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
        self._setup_complete = False
    
    async def cog_load(self):
        """Called when the cog is loaded"""
        try:
            # Start roast service auto-save task
            if hasattr(roast_service, 'start_auto_save'):
                await roast_service.start_auto_save()
            
            self._setup_complete = True
            logger.info("Roast cog loaded successfully")
        except Exception as e:
            logger.error(f"Error loading roast cog: {e}")
    
    async def cog_unload(self):
        """Called when the cog is unloaded"""
        try:
            # Stop roast service and save data
            if hasattr(roast_service, 'stop_auto_save'):
                await roast_service.stop_auto_save()
            
            logger.info("Roast cog unloaded")
        except Exception as e:
            logger.error(f"Error unloading roast cog: {e}")
    
    @commands.Cog.listener()
    async def on_message(self, message: discord.Message):
        """Track message activity for roast feature"""
        if message.author.bot or not self._setup_complete:
            return
        
        try:
            await roast_service.track_activity(
                user_id=message.author.id,
                activity_type=ActivityType.MESSAGE,
                channel_id=message.channel.id,
                guild_id=message.guild.id if message.guild else 0,
                content=message.content[:500]  # Limit content length for privacy
            )
        except Exception as e:
            logger.debug(f"Error tracking message activity: {e}")
    
    @commands.Cog.listener()
    async def on_voice_state_update(self, member: discord.Member, before: discord.VoiceState, after: discord.VoiceState):
        """Track voice channel activity for roast feature"""
        if member.bot or not self._setup_complete:
            return
        
        try:
            # User joined a voice channel
            if before.channel is None and after.channel is not None:
                await roast_service.track_activity(
                    user_id=member.id,
                    activity_type=ActivityType.VOICE_JOIN,
                    channel_id=after.channel.id,
                    guild_id=member.guild.id,
                    content=after.channel.name,
                    metadata={"channel_type": str(after.channel.type)}
                )
            
            # User left a voice channel
            elif before.channel is not None and after.channel is None:
                await roast_service.track_activity(
                    user_id=member.id,
                    activity_type=ActivityType.VOICE_LEAVE,
                    channel_id=before.channel.id,
                    guild_id=member.guild.id,
                    content=before.channel.name,
                    metadata={"channel_type": str(before.channel.type)}
                )
                
        except Exception as e:
            logger.debug(f"Error tracking voice activity for user {member.id}: {e}")
    
    @commands.Cog.listener()
    async def on_reaction_add(self, reaction: discord.Reaction, user: discord.User):
        """Track emoji reactions for roast feature"""
        if user.bot or not self._setup_complete:
            return
        
        try:
            await roast_service.track_activity(
                user_id=user.id,
                activity_type=ActivityType.EMOJI_REACTION,
                channel_id=reaction.message.channel.id,
                guild_id=reaction.message.guild.id if reaction.message.guild else 0,
                content=str(reaction.emoji),
                metadata={
                    "message_author_id": reaction.message.author.id,
                    "emoji_type": "custom" if hasattr(reaction.emoji, 'id') and reaction.emoji.id else "unicode"
                }
            )
        except Exception as e:
            logger.debug(f"Error tracking reaction activity for user {user.id}: {e}")
    
    # ==================== ROAST COMMANDS ====================
    
    @app_commands.command(name="roast", description="Get a personalized roast based on your Discord behavior!")
    @app_commands.describe(
        target="Who to roast (mention a user, or leave empty to roast yourself)",
        custom="Custom roast request (optional)"
    )
    async def roast_command(
        self, 
        interaction: discord.Interaction, 
        target: Optional[discord.Member] = None,
        custom: Optional[str] = None
    ):
        """Generate a personalized roast for a user"""
        if not await safe_defer_interaction(interaction):
            return
        
        try:
            # Determine target user
            target_user = target or interaction.user
            
            # Check if trying to roast the bot
            if target_user.bot:
                await safe_followup_send(
                    interaction,
                    "Nice try, but I'm roast-proof! My circuits are too cool to burn 🤖❄️"
                )
                return
            
            # Generate roast
            roast_text, roast_category = await roast_service.generate_roast(
                user_id=target_user.id,
                custom_prompt=custom
            )
            
            # Create roast embed
            embed = discord.Embed(
                title="🔥 ROAST DELIVERED! 🔥",
                description=roast_text,
                color=0xff4444
            )
            
            embed.set_author(
                name=f"Roasting {target_user.display_name}",
                icon_url=target_user.avatar.url if target_user.avatar else None
            )
            
            # Add roast metadata
            if roast_category != "cooldown":
                embed.add_field(
                    name="📊 Roast Category",
                    value=roast_category.replace("_", " ").title(),
                    inline=True
                )
                
                # Add roast stats
                user_stats = await roast_service.get_user_roast_stats(target_user.id)
                embed.add_field(
                    name="🎯 Roast Count",
                    value=f"{user_stats['roast_count']} total",
                    inline=True
                )
            
            # Add fun footer
            roast_footers = [
                "Roasted to perfection! 🍖",
                "Apply ice to burned area 🧊",
                "Destroyed with facts and logic 📊",
                "No survivors 💀",
                "Brought to you by NeruBot™ 🤖"
            ]
            embed.set_footer(text=random.choice(roast_footers))
            
            await safe_followup_send(interaction, embed=embed)
            
            # Track this as a command use
            await roast_service.track_activity(
                user_id=interaction.user.id,
                activity_type=ActivityType.COMMAND_USE,
                channel_id=interaction.channel.id,
                guild_id=interaction.guild.id if interaction.guild else 0,
                content=f"roast {target_user.display_name if target else 'self'}"
            )
            
        except Exception as e:
            logger.error(f"Error in roast command: {e}")
            await safe_followup_send(
                interaction, 
                "My roast generator exploded trying to process your epic behavior! Try again? 🤖💥"
            )
    
    @app_commands.command(name="roast-stats", description="View roasting statistics")
    @app_commands.describe(user="User to view stats for (optional)")
    async def roast_stats(self, interaction: discord.Interaction, user: Optional[discord.Member] = None):
        """Show roast statistics for a user or globally"""
        if not await safe_defer_interaction(interaction, ephemeral=True):
            return
        
        try:
            target_user = user or interaction.user
            
            if target_user.bot:
                await safe_followup_send(
                    interaction,
                    "Bots don't get roasted, they get debugged! 🤖🛠️",
                    ephemeral=True
                )
                return
            
            # Get user stats
            user_stats = await roast_service.get_user_roast_stats(target_user.id)
            global_stats = await roast_service.get_global_roast_stats()
            
            embed = discord.Embed(
                title="📊 Roast Statistics",
                color=0xff6b35
            )
            
            embed.set_author(
                name=target_user.display_name,
                icon_url=target_user.avatar.url if target_user.avatar else None
            )
            
            # User-specific stats
            embed.add_field(
                name="🎯 Personal Roast Stats",
                value=(
                    f"**Total Roasts:** {user_stats['roast_count']}\n"
                    f"**Can Be Roasted:** {'✅ Yes' if user_stats['can_be_roasted'] else '❌ On Cooldown'}\n"
                    f"**Activities Tracked:** {user_stats['activity_count']}\n"
                    f"**Behavior Analyzed:** {'✅ Yes' if user_stats['behavior_analyzed'] else '❌ Need More Data'}"
                ),
                inline=False
            )
            
            # Global stats
            embed.add_field(
                name="🌍 Global Roast Stats",
                value=(
                    f"**Total Roasts Delivered:** {global_stats.total_roasts_delivered:,}\n"
                    f"**Today's Roasts:** {global_stats.daily_roast_count}\n"
                    f"**Top Category:** {max(global_stats.roasts_by_category, key=global_stats.roasts_by_category.get) if global_stats.roasts_by_category else 'None'}"
                ),
                inline=False
            )
            
            # Activity breakdown
            if user_stats['behavior_analyzed']:
                pattern = roast_service.behavior_patterns.get(target_user.id)
                if pattern:
                    activity_text = ""
                    for activity_type, count in list(pattern.activity_frequency.items())[:5]:
                        activity_text += f"**{activity_type.value.title()}:** {count}\n"
                    
                    if activity_text:
                        embed.add_field(
                            name="📈 Recent Activity",
                            value=activity_text,
                            inline=True
                        )
            
            # Time until next roast
            if not user_stats['can_be_roasted'] and user_stats.get('last_roasted'):
                hours_left = 6 - ((time.time() - user_stats['last_roasted']) / 3600)
                embed.add_field(
                    name="⏰ Next Roast Available",
                    value=f"In {hours_left:.1f} hours",
                    inline=True
                )
            
            embed.set_footer(text="Roast responsibly! 🔥")
            
            await safe_followup_send(interaction, embed=embed, ephemeral=True)
            
        except Exception as e:
            logger.error(f"Error showing roast stats: {e}")
            await safe_followup_send(
                interaction,
                "❌ Couldn't retrieve roast stats right now. Try again later!",
                ephemeral=True
            )
    
    @app_commands.command(name="behavior-analysis", description="Get an analysis of your Discord behavior patterns")
    @app_commands.describe(user="User to analyze (optional)")
    async def behavior_analysis(self, interaction: discord.Interaction, user: Optional[discord.Member] = None):
        """Show detailed behavior analysis for a user"""
        if not await safe_defer_interaction(interaction, ephemeral=True):
            return
        
        try:
            target_user = user or interaction.user
            
            if target_user.bot:
                await safe_followup_send(
                    interaction,
                    "Bots are predictably unpredictable! 🤖🎲",
                    ephemeral=True
                )
                return
            
            # Analyze user behavior
            pattern = await roast_service.analyze_user_behavior(target_user.id)
            
            if not pattern.activity_frequency:
                await safe_followup_send(
                    interaction,
                    f"Not enough data to analyze {target_user.display_name}'s behavior yet. Need more Discord activity! 📊❌",
                    ephemeral=True
                )
                return
            
            embed = discord.Embed(
                title="🧠 Behavior Analysis Report",
                description=f"Detailed analysis of {target_user.display_name}'s Discord patterns",
                color=0x7b68ee
            )
            
            embed.set_author(
                name=target_user.display_name,
                icon_url=target_user.avatar.url if target_user.avatar else None
            )
            
            # Activity pattern
            if pattern.most_active_hours:
                hours_str = ", ".join(f"{h}:00" for h in pattern.most_active_hours[:5])
                embed.add_field(
                    name="⏰ Most Active Hours",
                    value=hours_str,
                    inline=False
                )
            
            # Activity frequency
            activity_text = ""
            for activity_type, count in list(pattern.activity_frequency.items())[:6]:
                bar_length = min(10, count // 2)
                bar = "█" * bar_length + "░" * (10 - bar_length)
                activity_text += f"**{activity_type.value.title()}:** {count} {bar}\n"
            
            if activity_text:
                embed.add_field(
                    name="📊 Activity Breakdown",
                    value=activity_text,
                    inline=False
                )
            
            # Behavioral insights
            insights = []
            
            if pattern.late_night_percentage > 30:
                insights.append(f"🌙 Night owl ({pattern.late_night_percentage:.0f}% late night activity)")
            
            if pattern.message_length_avg > 100:
                insights.append(f"📝 Detailed messenger (avg {pattern.message_length_avg:.0f} chars)")
            elif pattern.message_length_avg > 0 and pattern.message_length_avg < 30:
                insights.append(f"💬 Concise communicator (avg {pattern.message_length_avg:.0f} chars)")
            
            if len(pattern.emoji_usage) > 15:
                insights.append(f"😊 Emoji enthusiast ({len(pattern.emoji_usage)} different emojis)")
            
            if pattern.activity_frequency.get(ActivityType.MUSIC_REQUEST, 0) > 10:
                insights.append("🎵 Music lover")
            
            if insights:
                embed.add_field(
                    name="🔍 Behavioral Insights",
                    value="\n".join(insights),
                    inline=False
                )
            
            # Top words
            if pattern.common_words:
                words_str = ", ".join(pattern.common_words[:8])
                embed.add_field(
                    name="💭 Common Words",
                    value=words_str,
                    inline=False
                )
            
            # Most used emojis
            if pattern.emoji_usage:
                top_emojis = sorted(pattern.emoji_usage.items(), key=lambda x: x[1], reverse=True)[:10]
                emoji_str = " ".join(f"{emoji}({count})" for emoji, count in top_emojis)
                embed.add_field(
                    name="😄 Top Emojis",
                    value=emoji_str,
                    inline=False
                )
            
            embed.set_footer(text=f"Analysis based on {sum(pattern.activity_frequency.values())} activities")
            
            await safe_followup_send(interaction, embed=embed, ephemeral=True)
            
        except Exception as e:
            logger.error(f"Error in behavior analysis: {e}")
            await safe_followup_send(
                interaction,
                "❌ My behavior analysis brain had a glitch! Try again later.",
                ephemeral=True
            )


async def setup(bot: commands.Bot):
    """Setup function for the cog"""
    await bot.add_cog(RoastCog(bot))
