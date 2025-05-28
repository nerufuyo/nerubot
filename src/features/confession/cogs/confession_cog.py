"""
Confession commands cog for Discord bot
"""
import discord
from discord import ui
from discord.ext import commands
from discord import app_commands
from typing import Optional, List
from src.features.confession.services.confession_service import ConfessionService
from src.features.confession.models.confession import Confession, ConfessionReply
from src.core.utils.logging_utils import get_logger
from src.config.settings import DISCORD_CONFIG

logger = get_logger(__name__)


class ConfessionModal(ui.Modal):
    """Modal for submitting confessions."""
    
    def __init__(self, image_url: Optional[str] = None):
        super().__init__(title="Submit Anonymous Confession")
        self.image_url = image_url
        
        self.confession_input = ui.TextInput(
            label="Your Confession",
            placeholder="Type your anonymous confession here...",
            style=discord.TextStyle.paragraph,
            max_length=2000,
            required=True
        )
        self.add_item(self.confession_input)
    
    async def on_submit(self, interaction: discord.Interaction):
        # Get the cog to access the service
        cog = interaction.client.get_cog("ConfessionCog")
        if not cog:
            await interaction.response.send_message("❌ Confession system is not available.", ephemeral=True)
            return
        
        content = self.confession_input.value.strip()
        if not content:
            await interaction.response.send_message("❌ Please provide some content for your confession.", ephemeral=True)
            return
        
        logger.info(f"Creating new confession from user {interaction.user.id} in guild {interaction.guild.id}")
        
        # Create confession
        success, message, confession = cog.confession_service.create_confession(
            content=content,
            author_id=interaction.user.id,
            guild_id=interaction.guild.id,
            image_url=self.image_url
        )
        
        if not success:
            logger.warning(f"Failed to create confession: {message}")
            await interaction.response.send_message(f"❌ {message}", ephemeral=True)
            return
        
        logger.info(f"Successfully created confession {confession.confession_id}, now posting to channel")
        
        # Post confession to channel
        await cog.post_confession(confession, interaction.guild)
        
        await interaction.response.send_message(
            f"✅ Your confession has been submitted anonymously! (ID: `{confession.confession_id}`)",
            ephemeral=True
        )


class ReplyModal(ui.Modal):
    """Modal for replying to confessions."""
    
    def __init__(self, confession_id: int, image_url: Optional[str] = None):
        super().__init__(title=f"Reply to Confession #{confession_id}")
        self.confession_id = confession_id
        self.image_url = image_url
        
        self.reply_input = ui.TextInput(
            label="Your Reply",
            placeholder="Type your anonymous reply here...",
            style=discord.TextStyle.paragraph,
            max_length=1000,
            required=True
        )
        self.add_item(self.reply_input)
    
    async def on_submit(self, interaction: discord.Interaction):
        # Get the cog to access the service
        cog = interaction.client.get_cog("ConfessionCog")
        if not cog:
            await interaction.response.send_message("❌ Confession system is not available.", ephemeral=True)
            return
        
        content = self.reply_input.value.strip()
        if not content:
            await interaction.response.send_message("❌ Please provide some content for your reply.", ephemeral=True)
            return
        
        logger.info(f"Creating reply to confession {self.confession_id} from user {interaction.user.id}")
        
        # Create reply
        success, message, reply = cog.confession_service.create_reply(
            confession_id=self.confession_id,
            content=content,
            author_id=interaction.user.id,
            guild_id=interaction.guild.id,
            image_url=self.image_url
        )
        
        if not success:
            logger.warning(f"Failed to create reply: {message}")
            await interaction.response.send_message(f"❌ {message}", ephemeral=True)
            return
        
        logger.info(f"Successfully created reply {reply.reply_id}, now posting to thread")
        
        # Post reply to channel
        await cog.post_reply(reply, interaction.guild)
        
        await interaction.response.send_message(
            f"✅ Your reply has been posted anonymously in the confession thread!",
            ephemeral=True
        )


class ConfessionView(ui.View):
    """View with buttons for confession interactions."""
    
    def __init__(self, confession_id: int = None):
        super().__init__(timeout=None)  # Persistent view
        self.confession_id = confession_id

    @ui.button(label="Reply Anonymously", style=discord.ButtonStyle.secondary, emoji="💬", custom_id="confession_reply")
    async def reply_button(self, interaction: discord.Interaction, button: ui.Button):
        # Extract confession ID from the embed title
        confession_id = self._extract_confession_id(interaction)
        if confession_id is None:
            await interaction.response.send_message("❌ Could not find confession ID.", ephemeral=True)
            return
        
        logger.info(f"User {interaction.user.id} clicked reply button for confession {confession_id}")
        modal = ReplyModal(confession_id)
        await interaction.response.send_modal(modal)

    @ui.button(label="Create New Confession", style=discord.ButtonStyle.primary, emoji="📝", custom_id="confession_create_new")
    async def create_new_confession_button(self, interaction: discord.Interaction, button: ui.Button):
        logger.info(f"User {interaction.user.id} clicked create new confession button")
        modal = ConfessionModal()
        await interaction.response.send_modal(modal)
    
    def _extract_confession_id(self, interaction: discord.Interaction) -> Optional[int]:
        """Extract confession ID from the message embed."""
        if interaction.message and interaction.message.embeds:
            embed = interaction.message.embeds[0]
            title = embed.title or ""
            # Extract confession ID from title like "📝 Anonymous Confession #123"
            import re
            match = re.search(r'#(\d+)', title)
            if match:
                return int(match.group(1))
        return None


class SetupSuccessView(ui.View):
    """View with button for creating confessions after setup."""
    
    def __init__(self):
        super().__init__(timeout=None)  # Persistent view

    @ui.button(label="Create Confession", style=discord.ButtonStyle.green, emoji="📝", custom_id="create_confession")
    async def create_confession_button(self, interaction: discord.Interaction, button: ui.Button):
        modal = ConfessionModal()
        await interaction.response.send_modal(modal)


class ConfessionCog(commands.Cog):
    """Cog for anonymous confession functionality."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
        self.confession_service = ConfessionService()
    
    async def cog_load(self):
        """Called when the cog is loaded. Add persistent views here."""
        # Add persistent view for setup success button
        self.bot.add_view(SetupSuccessView())
        
        # Add one persistent view that can handle all confessions
        self.bot.add_view(ConfessionView())
        
        logger.info("ConfessionCog loaded with persistent views")
    
    @app_commands.command(name="confess", description="Submit an anonymous confession")
    @app_commands.describe(image="Optional image attachment to include with your confession")
    async def confess(self, interaction: discord.Interaction, image: Optional[discord.Attachment] = None):
        """Submit an anonymous confession through a modal."""
        image_url = None
        
        # Validate image if provided
        if image:
            if not image.content_type or not image.content_type.startswith('image/'):
                await interaction.response.send_message(
                    "❌ Please attach a valid image file (PNG, JPG, GIF, etc.)",
                    ephemeral=True
                )
                return
            
            # Check file size (limit to 8MB)
            if image.size > 8 * 1024 * 1024:
                await interaction.response.send_message(
                    "❌ Image too large! Please use an image smaller than 8MB.",
                    ephemeral=True
                )
                return
            
            image_url = image.url
        
        modal = ConfessionModal(image_url)
        await interaction.response.send_modal(modal)
    
    @app_commands.command(name="reply", description="Reply to a confession anonymously")
    @app_commands.describe(
        confession_id="The ID or tag of the confession to reply to",
        image="Optional image attachment to include with your reply"
    )
    async def reply_to_confession(self, interaction: discord.Interaction, confession_id: str, image: Optional[discord.Attachment] = None):
        """Reply to a confession using its ID."""
        # Find confession - the service now handles both int parsing and string prefix matching
        confession = self.confession_service.get_confession_by_tag(confession_id, interaction.guild.id)
        
        if not confession:
            await interaction.response.send_message(
                f"❌ No confession found with ID `{confession_id}` in this server.",
                ephemeral=True
            )
            return
        
        image_url = None
        
        # Validate image if provided
        if image:
            if not image.content_type or not image.content_type.startswith('image/'):
                await interaction.response.send_message(
                    "❌ Please attach a valid image file (PNG, JPG, GIF, etc.)",
                    ephemeral=True
                )
                return
            
            # Check file size (limit to 8MB)
            if image.size > 8 * 1024 * 1024:
                await interaction.response.send_message(
                    "❌ Image too large! Please use an image smaller than 8MB.",
                    ephemeral=True
                )
                return
            
            image_url = image.url
        
        modal = ReplyModal(confession.confession_id, image_url)
        await interaction.response.send_modal(modal)
    
    @app_commands.command(name="confession-setup", description="Set up confession channel (Admin only)")
    @app_commands.describe(channel="The channel where confessions will be posted")
    @app_commands.default_permissions(administrator=True)
    async def setup_confession_channel(self, interaction: discord.Interaction, channel: discord.TextChannel):
        """Set up the confession channel for the server."""
        self.confession_service.update_guild_settings(
            interaction.guild.id,
            confession_channel_id=channel.id
        )

        embed = discord.Embed(
            title="✅ Confession Channel Set",
            description=f"Anonymous confessions will now be posted to {channel.mention}",
            color=DISCORD_CONFIG["colors"]["success"]
        )
        
        embed.add_field(
            name="📋 Channel Information",
            value=(
                f"**Channel:** {channel.mention}\n"
                f"**Channel ID:** `{channel.id}`\n"
                f"**Thread Support:** Enabled\n"
                f"**Anonymous Replies:** Enabled"
            ),
            inline=False
        )
        
        embed.add_field(
            name="📝 How It Works",
            value=(
                "• Confessions will be posted as messages in this channel\n"
                "• Each confession will automatically create a thread\n"
                "• Replies will be posted inside the thread for better organization\n"
                "• Users can use the button below or `/confess` command to create confessions"
            ),
            inline=False
        )
        
        embed.set_footer(text="Use the button below to create your first confession!")

        # Create view with "Create Confession" button
        view = SetupSuccessView()
        
        await interaction.response.send_message(embed=embed, view=view)
    
    @app_commands.command(name="confession-settings", description="View/modify confession settings (Admin only)")
    @app_commands.default_permissions(administrator=True)
    async def confession_settings(self, interaction: discord.Interaction):
        """View confession settings for the server."""
        settings = self.confession_service.get_guild_settings(interaction.guild.id)
        
        embed = discord.Embed(
            title="⚙️ Confession Settings",
            color=DISCORD_CONFIG["colors"]["info"]
        )
        
        # Channel info
        if settings.confession_channel_id:
            channel = self.bot.get_channel(settings.confession_channel_id)
            channel_text = channel.mention if channel else "Channel not found"
        else:
            channel_text = "Not set"
        
        embed.add_field(name="Confession Channel", value=channel_text, inline=True)
        embed.add_field(name="Moderation", value="Enabled" if settings.moderation_enabled else "Disabled", inline=True)
        embed.add_field(name="Anonymous Replies", value="Enabled" if settings.anonymous_replies else "Disabled", inline=True)
        embed.add_field(name="Max Confession Length", value=f"{settings.max_confession_length} chars", inline=True)
        embed.add_field(name="Max Reply Length", value=f"{settings.max_reply_length} chars", inline=True)
        embed.add_field(name="Cooldown", value=f"{settings.cooldown_minutes} minutes", inline=True)
        embed.add_field(name="Next Confession ID", value=f"#{settings.next_confession_id}", inline=True)
        embed.add_field(name="Next Reply ID", value=f"#{settings.next_reply_id}", inline=True)
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="confession-stats", description="View confession statistics")
    async def confession_stats(self, interaction: discord.Interaction):
        """View confession statistics for the server."""
        confessions = self.confession_service.get_guild_confessions(interaction.guild.id, limit=100)
        
        if not confessions:
            await interaction.response.send_message("📊 No confessions found for this server.", ephemeral=True)
            return
        
        total_confessions = len(confessions)
        total_replies = sum(confession.reply_count for confession in confessions)
        
        embed = discord.Embed(
            title="📊 Confession Statistics",
            color=DISCORD_CONFIG["colors"]["info"]
        )
        
        embed.add_field(name="Total Confessions", value=total_confessions, inline=True)
        embed.add_field(name="Total Replies", value=total_replies, inline=True)
        embed.add_field(name="Average Replies", value=f"{total_replies/total_confessions:.1f}" if total_confessions > 0 else "0", inline=True)
        
        if confessions:
            latest = confessions[0]
            embed.add_field(name="Latest Confession", value=f"#{latest.confession_id}", inline=True)
        
        await interaction.response.send_message(embed=embed)
    
    async def post_confession(self, confession: Confession, guild: discord.Guild):
        """Post a confession to the designated channel."""
        settings = self.confession_service.get_guild_settings(guild.id)

        if not settings.confession_channel_id:
            logger.error(f"No confession channel set for guild {guild.id}")
            return

        channel = self.bot.get_channel(settings.confession_channel_id)
        if not channel:
            logger.error(f"Confession channel {settings.confession_channel_id} not found")
            return

        # Create embed
        embed = discord.Embed(
            title=f"📝 Anonymous Confession #{confession.confession_id}",
            description=confession.content,
            color=DISCORD_CONFIG["colors"]["info"],
            timestamp=confession.created_at
        )

        if confession.image_url:
            embed.set_image(url=confession.image_url)

        embed.set_footer(text="Use the buttons below to interact anonymously")

        # Create view with buttons
        view = ConfessionView(confession.confession_id)
        
        # Register the view with the bot for persistence
        self.bot.add_view(view)

        try:
            # Post the confession message first
            message = await channel.send(embed=embed, view=view)
            
            # Create a thread for this confession immediately
            thread_name = f"💬 Confession #{confession.confession_id}"
            thread = await message.create_thread(
                name=thread_name, 
                auto_archive_duration=10080  # 7 days
            )
            
            # Post a welcome message in the thread
            welcome_embed = discord.Embed(
                title="📝 Confession Discussion",
                description=f"This is the discussion thread for **Confession #{confession.confession_id}**.\n\nAnyone can reply anonymously using:\n• The **Reply** button on the original message\n• The `/reply {confession.confession_id}` command\n\nAll replies will appear here as anonymous messages.",
                color=DISCORD_CONFIG["colors"]["secondary"]
            )
            await thread.send(embed=welcome_embed)

            # Mark as posted with thread ID
            self.confession_service.mark_confession_posted(
                confession.confession_id,
                channel.id,
                message.id,
                thread.id
            )

            logger.info(f"Posted confession {confession.confession_id} to {guild.name} with thread {thread.id}")

        except Exception as e:
            logger.error(f"Error posting confession: {e}")
            import traceback
            logger.error(traceback.format_exc())
    
    async def post_reply(self, reply: ConfessionReply, guild: discord.Guild):
        """Post a reply to the confession thread."""
        logger.info(f"Attempting to post reply {reply.reply_id} for confession {reply.confession_id}")
        
        settings = self.confession_service.get_guild_settings(guild.id)

        if not settings.confession_channel_id:
            logger.error(f"No confession channel set for guild {guild.id}")
            return

        # Get the confession to find the thread ID
        confession = self.confession_service.get_confession(reply.confession_id)
        if not confession:
            logger.error(f"Confession {reply.confession_id} not found")
            return
            
        if not confession.thread_id:
            logger.error(f"No thread ID found for confession {reply.confession_id}")
            return

        logger.info(f"Found confession {confession.confession_id} with thread_id {confession.thread_id}")

        # Get the thread and main channel
        try:
            thread = await self.bot.fetch_channel(confession.thread_id)
        except discord.NotFound:
            logger.error(f"Thread {confession.thread_id} not found")
            return
        except discord.Forbidden:
            logger.error(f"No permission to access thread {confession.thread_id}")
            return
        except Exception as e:
            logger.error(f"Error fetching thread {confession.thread_id}: {e}")
            return
            
        main_channel = self.bot.get_channel(settings.confession_channel_id)
        if not main_channel:
            logger.error(f"Main channel {settings.confession_channel_id} not found")
            return

        try:
            # Create an embed for the reply to match confession style
            embed = discord.Embed(
                title=f"💬 Anonymous Reply #{reply.reply_id}",
                description=reply.content,
                color=DISCORD_CONFIG["colors"]["secondary"],
                timestamp=reply.created_at
            )
            
            if reply.image_url:
                embed.set_image(url=reply.image_url)
            
            embed.set_footer(text=f"Reply to Confession #{reply.confession_id}")
            
            logger.info(f"Posting reply to thread {thread.name} (ID: {thread.id})")
            
            # Post reply embed in the thread
            message = await thread.send(embed=embed)

            # Mark as posted
            self.confession_service.mark_reply_posted(reply.reply_id, message.id)

            logger.info(f"Successfully posted reply {reply.reply_id} to confession {reply.confession_id} thread")

        except discord.Forbidden:
            logger.error(f"No permission to post in thread {thread.id}")
        except Exception as e:
            logger.error(f"Error posting reply: {e}")
            import traceback
            logger.error(traceback.format_exc())


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(ConfessionCog(bot))
