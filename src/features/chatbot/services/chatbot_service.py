"""
Chatbot service for managing chat sessions and interactions
"""
import asyncio
import time
import random
from typing import Dict, Optional, List
from discord import TextChannel, User, Embed, File
from discord.ext import tasks
from src.features.chatbot.models.chat_models import ChatSession, ChatMessage, ChatStats
from src.services.ai_service import ai_service, AIProvider, AIResponse
from src.core.utils.logging_utils import get_logger

logger = get_logger(__name__)


class ChatbotService:
    """Service for managing chatbot interactions"""
    
    def __init__(self):
        self.active_sessions: Dict[int, ChatSession] = {}  # user_id -> ChatSession
        self.user_stats: Dict[int, ChatStats] = {}  # user_id -> ChatStats
        self.session_timeout_minutes = 5
        
        # Welcome messages for new interactions
        self.welcome_messages = [
            "Hey there! 👋 What's on your mind today?",
            "Oh, someone wants to chat! What's up? 😊",
            "Heyyy~ What can I help you with? ✨",
            "Well well, look who's back! What's cooking? 🍳",
            "Yo! Ready for some quality bot conversation? 🤖",
            "Greetings, human! What adventure shall we embark on? 🚀",
            "Oh hi there! I was just organizing my digital thoughts~ 💭",
            "Another chat session? I'm all ears! Well... all code~ 👂"
        ]
        
        # Thanks messages for when users go idle
        self.thanks_messages = [
            "Thanks for chatting with me! Come back anytime~ 🌟",
            "It was fun talking! Don't be a stranger! 👋",
            "Hope I helped! Catch you later~ ✨",
            "Always a pleasure! See you around! 🎵",
            "Thanks for the chat! Until next time! 🤖",
            "That was nice! I'll be here when you need me~ 💙",
            "Enjoyed our conversation! Take care! 🌈",
            "Thanks for hanging out! Stay awesome! ⭐"
        ]
    
    async def start_session_monitor(self, bot=None):
        """Start monitoring for expired sessions"""
        if not self.session_cleanup_task.is_running():
            if bot:
                self.session_cleanup_task._bot = bot
            self.session_cleanup_task.start()
            logger.info("Chatbot session monitor started")
    
    async def stop_session_monitor(self):
        """Stop monitoring sessions"""
        if self.session_cleanup_task.is_running():
            self.session_cleanup_task.cancel()
            logger.info("Chatbot session monitor stopped")
    
    @tasks.loop(minutes=1)
    async def session_cleanup_task(self):
        """Clean up expired sessions and send thanks messages"""
        try:
            expired_sessions = []
            
            for user_id, session in self.active_sessions.items():
                if session.is_expired(self.session_timeout_minutes) and session.is_active:
                    expired_sessions.append((user_id, session))
            
            # Send thanks messages and mark sessions as inactive
            for user_id, session in expired_sessions:
                try:
                    from discord.utils import get
                    bot = getattr(self.session_cleanup_task, '_bot', None)
                    if bot:
                        channel = bot.get_channel(session.channel_id)
                        if channel:
                            await self._send_thanks_message(channel, user_id)
                            session.is_active = False
                            logger.info(f"Sent thanks message to user {user_id} after {session.get_idle_time_minutes():.1f} minutes of inactivity")
                except Exception as e:
                    logger.error(f"Error sending thanks message to user {user_id}: {e}")
            
            # Clean up old inactive sessions (after 30 minutes)
            cleanup_cutoff = time.time() - (30 * 60)  # 30 minutes ago
            sessions_to_remove = [
                user_id for user_id, session in self.active_sessions.items()
                if not session.is_active and session.last_interaction < cleanup_cutoff
            ]
            
            for user_id in sessions_to_remove:
                del self.active_sessions[user_id]
                logger.debug(f"Cleaned up old session for user {user_id}")
                
        except Exception as e:
            logger.error(f"Error in session cleanup task: {e}")
    
    async def handle_message(self, user: User, channel: TextChannel, message: str) -> Optional[str]:
        """
        Handle a chat message from a user
        
        Args:
            user: Discord user who sent the message
            channel: Discord channel where message was sent
            message: The message content
            
        Returns:
            AI response or None if there was an error
        """
        user_id = user.id
        channel_id = channel.id
        
        # Get or create session
        session = self.active_sessions.get(user_id)
        is_new_session = session is None or not session.is_active
        
        if session is None:
            session = ChatSession(user_id=user_id, channel_id=channel_id)
            self.active_sessions[user_id] = session
        else:
            session.is_active = True
            session.channel_id = channel_id  # Update channel in case they moved
        
        session.update_interaction()
        
        # Send welcome message for new sessions
        if is_new_session:
            await self._send_welcome_message(channel, user)
        
        # Get AI response
        try:
            # Use automatic provider fallback (priority: Claude -> Gemini -> OpenAI)
            ai_response = await ai_service.chat(
                message=message,
                max_tokens=300,
                temperature=0.8  # Higher creativity for fun responses
            )
            
            # Update user stats
            if user_id not in self.user_stats:
                self.user_stats[user_id] = ChatStats(user_id=user_id)
            self.user_stats[user_id].update_stats(ai_response.provider)
            
            # Log the interaction
            chat_msg = ChatMessage(
                user_id=user_id,
                channel_id=channel_id,
                content=message,
                ai_response=ai_response.content,
                ai_provider=ai_response.provider,
                response_time=time.time()
            )
            
            logger.info(f"Chat response generated for user {user_id} using {ai_response.provider.value}")
            
            return ai_response.content
            
        except Exception as e:
            logger.error(f"Error handling chat message from user {user_id}: {e}")
            return "Oops! My circuits got a bit tangled there. Try again? 🤖⚡"
    
    async def set_user_preferred_provider(self, user_id: int, provider: AIProvider):
        """Set a user's preferred AI provider (for future use)"""
        if user_id not in self.user_stats:
            self.user_stats[user_id] = ChatStats(user_id=user_id)
        self.user_stats[user_id].favorite_provider = provider
        logger.info(f"Set preferred AI provider for user {user_id} to {provider.value}")
    
    async def _send_welcome_message(self, channel: TextChannel, user: User):
        """Send a welcome message to the user"""
        try:
            welcome_msg = random.choice(self.welcome_messages)
            
            embed = Embed(
                title="🤖 NeruBot Chat",
                description=welcome_msg,
                color=0x00ff9f
            )
            embed.set_author(name=f"Hello, {user.display_name}!", icon_url=user.avatar.url if user.avatar else None)
            embed.set_footer(text="💡 Just chat naturally - I'll respond with my AI brain!")
            
            await channel.send(embed=embed)
            
        except Exception as e:
            logger.error(f"Error sending welcome message: {e}")
            # Fallback to simple message
            try:
                await channel.send(f"Hey {user.mention}! 👋 What's up?")
            except:
                pass
    
    async def _send_thanks_message(self, channel: TextChannel, user_id: int):
        """Send a thanks message when user goes idle"""
        try:
            thanks_msg = random.choice(self.thanks_messages)
            
            embed = Embed(
                title="✨ Thanks for Chatting!",
                description=thanks_msg,
                color=0xff6b9d
            )
            
            # Try to get bot's avatar for the thanks message
            try:
                bot = channel.guild.me
                if bot and bot.avatar:
                    embed.set_thumbnail(url=bot.avatar.url)
            except:
                pass
            
            embed.set_footer(text="I'll be here when you want to chat again! 💙")
            
            await channel.send(embed=embed)
            
        except Exception as e:
            logger.error(f"Error sending thanks message: {e}")
    
    async def get_user_stats(self, user_id: int) -> Optional[ChatStats]:
        """Get chat statistics for a user"""
        return self.user_stats.get(user_id)
    
    async def get_active_sessions_count(self) -> int:
        """Get count of active chat sessions"""
        return len([s for s in self.active_sessions.values() if s.is_active])
    
    async def set_user_preferred_provider(self, user_id: int, provider: AIProvider):
        """Set a user's preferred AI provider (for future use)"""
        if user_id not in self.user_stats:
            self.user_stats[user_id] = ChatStats(user_id=user_id)
        self.user_stats[user_id].favorite_provider = provider
        logger.info(f"Set preferred AI provider for user {user_id} to {provider.value}")


# Global chatbot service instance
chatbot_service = ChatbotService()
