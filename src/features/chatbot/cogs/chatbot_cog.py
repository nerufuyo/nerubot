"""
Chatbot Discord Cog - Main chat interface with AI integration
"""
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional
from src.features.chatbot.services.chatbot_service import chatbot_service
from src.services.ai_service import ai_service, AIProvider
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
        # Interaction expired, try to send a message in the channel instead
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
        # Try to send directly to channel as fallback (only for non-ephemeral messages)
        if not ephemeral:
            try:
                mention = f"{interaction.user.mention} " if content else ""
                if embed:
                    await interaction.channel.send(content=mention + (content or ""), embed=embed)
                else:
                    await interaction.channel.send(mention + content)
            except Exception as fallback_error:
                logger.error(f"Fallback message failed: {fallback_error}")


class ChatbotCog(commands.Cog):
    """Discord cog for AI chatbot functionality"""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
        self._setup_complete = False
    
    async def cog_load(self):
        """Called when the cog is loaded"""
        try:
            await ai_service.initialize()
            await chatbot_service.start_session_monitor(self.bot)
            
            self._setup_complete = True
            logger.info("Chatbot cog loaded successfully")
        except Exception as e:
            logger.error(f"Error loading chatbot cog: {e}")
    
    async def cog_unload(self):
        """Called when the cog is unloaded"""
        try:
            await chatbot_service.stop_session_monitor()
            await ai_service.cleanup()
            
            logger.info("Chatbot cog unloaded")
        except Exception as e:
            logger.error(f"Error unloading chatbot cog: {e}")
    
    @commands.Cog.listener()
    async def on_message(self, message: discord.Message):
        """Listen for messages and respond with AI if mentioned or in DM"""
        # Ignore bot messages
        if message.author.bot:
            return
        
        # Only respond if setup is complete
        if not self._setup_complete:
            return
        
        # Check if bot is mentioned or if it's a DM
        bot_mentioned = self.bot.user.mentioned_in(message) and not message.mention_everyone
        is_dm = isinstance(message.channel, discord.DMChannel)
        
        # For guild messages, only respond to mentions
        # For DMs, respond to all messages
        should_respond = is_dm or bot_mentioned
        
        if should_respond:
            # Remove bot mention from message if present
            content = message.content
            if bot_mentioned:
                content = content.replace(f'<@{self.bot.user.id}>', '').strip()
                content = content.replace(f'<@!{self.bot.user.id}>', '').strip()
            
            # Skip empty messages
            if not content.strip():
                return
            
            # Show typing indicator
            async with message.channel.typing():
                try:
                    # Get AI response
                    response = await chatbot_service.handle_message(
                        user=message.author,
                        channel=message.channel,
                        message=content
                    )
                    
                    if response:
                        # Send response
                        await message.reply(response, mention_author=False)
                    
                except Exception as e:
                    logger.error(f"Error processing chat message: {e}")
                    await message.reply(
                        "Uh oh! Something went wrong in my neural networks... 🤖💥 Try again?",
                        mention_author=False
                    )
    
    @app_commands.command(name="chat", description="Start a conversation with the AI")
    @app_commands.describe(message="Your message to the AI")
    async def chat_command(self, interaction: discord.Interaction, message: str):
        """Direct chat command for AI interaction"""
        if not await safe_defer_interaction(interaction):
            return
        
        try:
            # Get AI response
            response = await chatbot_service.handle_message(
                user=interaction.user,
                channel=interaction.channel,
                message=message
            )
            
            if response:
                await safe_followup_send(interaction, response)
            else:
                await safe_followup_send(
                    interaction,
                    "I couldn't generate a response right now. Please try again!"
                )
                
        except Exception as e:
            logger.error(f"Error in chat command: {e}")
            await safe_followup_send(
                interaction,
                "Something went wrong while processing your message. Please try again!"
            )
    
    @app_commands.command(name="reset-chat", description="Reset your conversation history with the AI")
    async def reset_chat(self, interaction: discord.Interaction):
        """Reset the user's chat session"""
        if not await safe_defer_interaction(interaction, ephemeral=True):
            return
        
        try:
            # Reset the user's session
            await chatbot_service.reset_user_session(interaction.user.id)
            
            embed = discord.Embed(
                title="🔄 Chat Session Reset",
                description="Your conversation history has been cleared. Starting fresh!",
                color=0x00ff00
            )
            
            await safe_followup_send(interaction, embed=embed, ephemeral=True)
            
        except Exception as e:
            logger.error(f"Error resetting chat session: {e}")
            await safe_followup_send(
                interaction,
                "❌ Failed to reset your chat session. Please try again!",
                ephemeral=True
            )


async def setup(bot: commands.Bot):
    """Setup function for the cog"""
    await bot.add_cog(ChatbotCog(bot))
