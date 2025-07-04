"""
Chatbot Discord Cog - Main chat interface with AI integration
"""
import discord
from discord.ext import commands
from discord import app_commands
from typing import Optional, List
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
    
    @app_commands.command(name="chat", description="Start a conversation with NeruBot!")
    @app_commands.describe(message="What would you like to say?")
    async def chat_command(self, interaction: discord.Interaction, message: str):
        """Slash command to chat with the bot"""
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
                embed = discord.Embed(
                    title="🤖 NeruBot Chat",
                    description=response,
                    color=0x00ff9f
                )
                embed.set_author(
                    name=interaction.user.display_name,
                    icon_url=interaction.user.avatar.url if interaction.user.avatar else None
                )
                await safe_followup_send(interaction, embed=embed)
            else:
                await safe_followup_send(interaction, "Sorry, I couldn't process that right now. Try again?")
        
        except Exception as e:
            logger.error(f"Error in chat command: {e}")
            await safe_followup_send(interaction, "Oops! My AI brain had a hiccup. Try again! 🤖")
    
    @app_commands.command(name="ai-provider", description="Set your preferred AI provider")
    @app_commands.describe(provider="Choose your preferred AI provider")
    async def set_ai_provider(
        self,
        interaction: discord.Interaction,
        provider: str
    ):
        """Set user's preferred AI provider"""
        if not await safe_defer_interaction(interaction, ephemeral=True):
            return
        
        try:
            # Convert string to enum
            ai_provider = AIProvider(provider.lower())
            
            # Check if provider is available
            available_providers = ai_service.get_available_providers()
            if ai_provider not in available_providers:
                await interaction.followup.send(
                    f"❌ {ai_provider.value.title()} is not currently available. "
                    f"Available providers: {', '.join(p.value.title() for p in available_providers)}",
                    ephemeral=True
                )
                return
            
            # Set user preference
            await chatbot_service.set_user_preferred_provider(interaction.user.id, ai_provider)
            
            await interaction.followup.send(
                f"✅ Set your preferred AI provider to **{ai_provider.value.title()}**! "
                f"I'll use this for our future chats~ 🤖",
                ephemeral=True
            )
        
        except ValueError:
            await interaction.followup.send(
                "❌ Invalid AI provider. Please choose from the available options.",
                ephemeral=True
            )
        except Exception as e:
            logger.error(f"Error setting AI provider: {e}")
            await interaction.followup.send(
                "❌ Something went wrong while setting your AI provider. Try again later!",
                ephemeral=True
            )
    
    @set_ai_provider.autocomplete('provider')
    async def ai_provider_autocomplete(
        self,
        interaction: discord.Interaction,
        current: str,
    ) -> List[app_commands.Choice[str]]:
        """Autocomplete for AI provider selection"""
        choices = []
        available_providers = ai_service.get_available_providers()
        
        provider_names = {
            AIProvider.OPENAI: "OpenAI GPT",
            AIProvider.CLAUDE: "Anthropic Claude",
            AIProvider.GEMINI: "Google Gemini"
        }
        
        for provider in available_providers:
            name = provider_names.get(provider, provider.value.title())
            if current.lower() in name.lower():
                choices.append(app_commands.Choice(name=name, value=provider.value))
        
        return choices[:25]  # Discord limit
    
    @app_commands.command(name="chat-stats", description="View your chat statistics")
    async def chat_stats(self, interaction: discord.Interaction):
        """Show user's chat statistics"""
        if not await safe_defer_interaction(interaction, ephemeral=True):
            return
        
        try:
            stats = await chatbot_service.get_user_stats(interaction.user.id)
            
            if not stats or stats.total_messages == 0:
                await interaction.followup.send(
                    "📊 No chat statistics yet! Start chatting with me to see your stats~ 🤖",
                    ephemeral=True
                )
                return
            
            embed = discord.Embed(
                title="📊 Your Chat Statistics",
                color=0x00d4ff
            )
            embed.set_author(
                name=interaction.user.display_name,
                icon_url=interaction.user.avatar.url if interaction.user.avatar else None
            )
            
            embed.add_field(
                name="💬 Total Messages",
                value=f"{stats.total_messages:,}",
                inline=True
            )
            embed.add_field(
                name="🤖 AI Responses",
                value=f"{stats.total_ai_responses:,}",
                inline=True
            )
            
            if stats.favorite_provider:
                embed.add_field(
                    name="🎯 Preferred AI",
                    value=stats.favorite_provider.value.title(),
                    inline=True
                )
            
            if stats.first_chat:
                import datetime
                first_chat_date = datetime.datetime.fromtimestamp(stats.first_chat)
                embed.add_field(
                    name="📅 First Chat",
                    value=first_chat_date.strftime("%Y-%m-%d"),
                    inline=True
                )
            
            # Current session info
            session = chatbot_service.active_sessions.get(interaction.user.id)
            if session and session.is_active:
                embed.add_field(
                    name="🟢 Current Session",
                    value=f"{session.message_count} messages",
                    inline=True
                )
            
            embed.set_footer(text="Thanks for chatting with me! 💙")
            
            await interaction.followup.send(embed=embed, ephemeral=True)
        
        except Exception as e:
            logger.error(f"Error showing chat stats: {e}")
            await interaction.followup.send(
                "❌ Couldn't retrieve your chat stats right now. Try again later!",
                ephemeral=True
            )
    
    @app_commands.command(name="ai-status", description="Check AI services status")
    async def ai_status(self, interaction: discord.Interaction):
        """Show AI services status"""
        if not await safe_defer_interaction(interaction, ephemeral=True):
            return
        
        try:
            available_providers = ai_service.get_available_providers()
            active_sessions = await chatbot_service.get_active_sessions_count()
            
            embed = discord.Embed(
                title="🤖 AI Services Status",
                color=0x00ff9f if available_providers else 0xff6b6b
            )
            
            # Provider status
            all_providers = [AIProvider.OPENAI, AIProvider.CLAUDE, AIProvider.GEMINI]
            provider_names = {
                AIProvider.OPENAI: "OpenAI GPT",
                AIProvider.CLAUDE: "Anthropic Claude",
                AIProvider.GEMINI: "Google Gemini"
            }
            
            status_text = ""
            for provider in all_providers:
                name = provider_names[provider]
                status = "🟢 Available" if provider in available_providers else "🔴 Unavailable"
                status_text += f"{name}: {status}\n"
            
            embed.add_field(
                name="🔌 Provider Status",
                value=status_text,
                inline=False
            )
            
            embed.add_field(
                name="📊 Active Sessions",
                value=f"{active_sessions} users currently chatting",
                inline=True
            )
            
            if available_providers:
                embed.set_footer(text="All systems operational! 🚀")
            else:
                embed.set_footer(text="No AI providers available - check API keys!")
            
            await interaction.followup.send(embed=embed, ephemeral=True)
        
        except Exception as e:
            logger.error(f"Error showing AI status: {e}")
            await interaction.followup.send(
                "❌ Couldn't check AI status right now. Try again later!",
                ephemeral=True
            )


async def setup(bot: commands.Bot):
    """Setup function for the cog"""
    await bot.add_cog(ChatbotCog(bot))
