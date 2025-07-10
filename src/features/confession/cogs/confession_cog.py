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
from src.core.constants import CONFESSION_CONSTANTS, CONFESSION_LOG_MESSAGES

logger = get_logger(__name__)


class ConfessionModal(ui.Modal):
    """Modal for submitting confessions with exactly 2 fields as specified."""
    
    def __init__(self):
        super().__init__(title=CONFESSION_CONSTANTS["messages"]["modals"]["confession_title"])
        
        # Field 1: Message (required)
        self.message_input = ui.TextInput(
            label=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["message_label"],
            placeholder=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["message_placeholder"].format(type="confession"),
            style=discord.TextStyle.paragraph,
            max_length=CONFESSION_CONSTANTS["settings"]["defaults"]["max_confession_length"],
            required=True
        )
        self.add_item(self.message_input)
        
        # Field 2: Attachments (optional) - text input for URLs
        self.attachments_input = ui.TextInput(
            label=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["attachments_label"],
            placeholder=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["attachments_placeholder"],
            style=discord.TextStyle.short,
            max_length=500,
            required=False
        )
        self.add_item(self.attachments_input)
    
    async def on_submit(self, interaction: discord.Interaction):
        # Get the cog to access the service
        cog = interaction.client.get_cog("ConfessionCog")
        if not cog:
            await interaction.response.send_message(
                f"❌ {CONFESSION_CONSTANTS['messages']['errors']['system_unavailable']}", 
                ephemeral=True
            )
            return
        
        content = self.message_input.value.strip()
        if not content:
            await interaction.response.send_message(
                f"❌ {CONFESSION_CONSTANTS['messages']['errors']['empty_content'].format(type='confession')}", 
                ephemeral=True
            )
            return
        
        # Parse attachments
        attachments = []
        if self.attachments_input.value.strip():
            attachment_urls = [url.strip() for url in self.attachments_input.value.split(',')]
            # Filter out empty URLs and validate them
            for url in attachment_urls:
                if url and (url.startswith('http://') or url.startswith('https://')):
                    attachments.append(url)
                elif url:
                    logger.warning(f"Invalid attachment URL: {url}")
        
        logger.info(f"Creating new confession from user {interaction.user.id} with {len(attachments)} attachments")
        
        # Create confession
        success, message, confession_id = cog.confession_service.create_confession(
            content=content,
            author_id=interaction.user.id,
            guild_id=interaction.guild.id,
            attachments=attachments
        )
        
        if not success:
            logger.warning(f"Failed to create confession: {message}")
            await interaction.response.send_message(f"❌ {message}", ephemeral=True)
            return
        
        # Get the created confession
        confession = cog.confession_service.get_confession(int(confession_id.split('-')[1]))
        if not confession:
            logger.error(f"Could not retrieve created confession {confession_id}")
            await interaction.response.send_message(
                "❌ Error retrieving created confession", 
                ephemeral=True
            )
            return
        
        logger.info(f"Successfully created confession {confession_id}, now posting to channel")
        
        # Post confession to channel
        await cog.post_confession(confession, interaction.guild)
        
        await interaction.response.send_message(
            f"✅ {CONFESSION_CONSTANTS['messages']['success']['confession_created'].format(id=confession_id)}",
            ephemeral=True
        )


class ReplyModal(ui.Modal):
    """Modal for replying to confessions with exactly 3 fields as specified."""
    
    def __init__(self, confession_id: int):
        super().__init__(title=CONFESSION_CONSTANTS["messages"]["modals"]["reply_title"].format(id=f"CONF-{confession_id:03d}"))
        self.confession_id = confession_id
        
        # Field 1: Message (required)
        self.message_input = ui.TextInput(
            label=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["message_label"],
            placeholder=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["message_placeholder"].format(type="reply"),
            style=discord.TextStyle.paragraph,
            max_length=CONFESSION_CONSTANTS["settings"]["defaults"]["max_reply_length"],
            required=True
        )
        self.add_item(self.message_input)
        
        # Field 2: Confession ID (auto-populated, read-only)
        self.confession_id_input = ui.TextInput(
            label=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["confession_id_label"],
            placeholder=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["confession_id_placeholder"],
            style=discord.TextStyle.short,
            default=f"CONF-{confession_id:03d}",
            max_length=20,
            required=False
        )
        self.add_item(self.confession_id_input)
        
        # Field 3: Attachments (optional)
        self.attachments_input = ui.TextInput(
            label=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["attachments_label"],
            placeholder=CONFESSION_CONSTANTS["messages"]["modals"]["fields"]["attachments_placeholder"],
            style=discord.TextStyle.short,
            max_length=500,
            required=False
        )
        self.add_item(self.attachments_input)
    
    async def on_submit(self, interaction: discord.Interaction):
        # Get the cog to access the service
        cog = interaction.client.get_cog("ConfessionCog")
        if not cog:
            await interaction.response.send_message(
                f"❌ {CONFESSION_CONSTANTS['messages']['errors']['system_unavailable']}", 
                ephemeral=True
            )
            return
        
        content = self.message_input.value.strip()
        if not content:
            await interaction.response.send_message(
                f"❌ {CONFESSION_CONSTANTS['messages']['errors']['empty_content'].format(type='reply')}", 
                ephemeral=True
            )
            return
        
        # Parse attachments
        attachments = []
        if self.attachments_input.value.strip():
            attachment_urls = [url.strip() for url in self.attachments_input.value.split(',')]
            # Filter out empty URLs and validate them
            for url in attachment_urls:
                if url and (url.startswith('http://') or url.startswith('https://')):
                    attachments.append(url)
                elif url:
                    logger.warning(f"Invalid attachment URL: {url}")
        
        logger.info(f"Creating reply to confession {self.confession_id} with {len(attachments)} attachments")
        
        # Create reply
        success, message, reply_id = cog.confession_service.create_reply(
            confession_id=self.confession_id,
            content=content,
            author_id=interaction.user.id,
            guild_id=interaction.guild.id,
            attachments=attachments
        )
        
        if not success:
            logger.warning(f"Failed to create reply: {message}")
            await interaction.response.send_message(f"❌ {message}", ephemeral=True)
            return
        
        # Get the created reply
        replies = cog.confession_service.get_confession_replies(self.confession_id)
        reply = None
        for r in replies:
            if r.reply_id == reply_id:
                reply = r
                break
        
        if not reply:
            logger.error(f"Could not retrieve created reply {reply_id}")
            await interaction.response.send_message(
                "❌ Error retrieving created reply", 
                ephemeral=True
            )
            return
        
        logger.info(f"Successfully created reply {reply_id}, now posting to thread")
        
        # Post reply to channel
        await cog.post_reply(reply, interaction.guild)
        
        await interaction.response.send_message(
            f"✅ {CONFESSION_CONSTANTS['messages']['success']['reply_created'].format(id=reply_id)}",
            ephemeral=True
        )


class ConfessionView(ui.View):
    """View with buttons for confession interactions."""
    
    def __init__(self, confession_id: int = None):
        super().__init__(timeout=None)  # Persistent view
        self.confession_id = confession_id

    @ui.button(label=CONFESSION_CONSTANTS["messages"]["buttons"]["reply"], style=discord.ButtonStyle.secondary, emoji="🔄", custom_id="confession_reply")
    async def reply_button(self, interaction: discord.Interaction, button: ui.Button):
        # First try to use the confession_id stored in the view
        confession_id = self.confession_id
        
        # If not available, try to extract from the embed title
        if confession_id is None:
            confession_id = self._extract_confession_id(interaction)
        
        if confession_id is None:
            await interaction.response.send_message(
                f"❌ {CONFESSION_CONSTANTS['messages']['errors']['confession_not_found']}", 
                ephemeral=True
            )
            return
        
        logger.info(f"User {interaction.user.id} clicked reply button for confession {confession_id}")
        modal = ReplyModal(confession_id)
        await interaction.response.send_modal(modal)

    @ui.button(label=CONFESSION_CONSTANTS["messages"]["buttons"]["create_confession"], style=discord.ButtonStyle.primary, emoji="📝", custom_id="confession_create_new")
    async def create_new_confession_button(self, interaction: discord.Interaction, button: ui.Button):
        logger.info(f"User {interaction.user.id} clicked create new confession button")
        modal = ConfessionModal()
        await interaction.response.send_modal(modal)
    
    def _extract_confession_id(self, interaction: discord.Interaction) -> Optional[int]:
        """Extract confession ID from the message embed."""
        if interaction.message and interaction.message.embeds:
            embed = interaction.message.embeds[0]
            title = embed.title or ""
            # Extract confession ID from title like "📝 Confession #001" or "↪️ Reply to #001"
            import re
            match = re.search(r'#(\d+)', title)
            if match:
                return int(match.group(1))
        
        # Try to extract from footer text like "ID: CONF-001 | 🔄 Reply"
        if interaction.message and interaction.message.embeds:
            embed = interaction.message.embeds[0]
            footer = embed.footer.text if embed.footer else ""
            match = re.search(r'CONF-(\d+)', footer)
            if match:
                return int(match.group(1))
        
        return None


class SetupView(ui.View):
    """View with button for creating confessions after setup."""
    
    def __init__(self):
        super().__init__(timeout=None)  # Persistent view

    @ui.button(label=CONFESSION_CONSTANTS["messages"]["buttons"]["create_confession"], style=discord.ButtonStyle.primary, emoji="📝", custom_id="create_new_confession")
    async def create_new_confession_button(self, interaction: discord.Interaction, button: ui.Button):
        logger.info(f"User {interaction.user.id} clicked create new confession button")
        modal = ConfessionModal()
        await interaction.response.send_modal(modal)


class ConfessionCog(commands.Cog):
    """Cog for anonymous confession functionality."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
        self.confession_service = ConfessionService()
    
    async def cog_load(self):
        """Called when the cog is loaded. Add persistent views here."""
        # Add persistent view for setup button
        self.bot.add_view(SetupView())
        
        # Add one persistent view that can handle all confessions
        self.bot.add_view(ConfessionView())
        
        logger.info(CONFESSION_LOG_MESSAGES["system"]["cog_loaded"])
    
    @app_commands.command(name="confession-setup", description="Set up confession channel (Admin only)")
    @app_commands.describe(channel="The channel where confessions will be posted")
    @app_commands.default_permissions(administrator=True)
    async def setup_confession_channel(self, interaction: discord.Interaction, channel: discord.TextChannel):
        """Set up the confession channel for the server."""
        self.confession_service.update_guild_settings(
            interaction.guild.id,
            confession_channel_id=channel.id
        )

        # Send introduction message to the confession channel
        intro_embed = discord.Embed(
            title=CONFESSION_CONSTANTS["messages"]["titles"]["setup_intro"].format(emoji="📝"),
            description=CONFESSION_CONSTANTS["messages"]["descriptions"]["intro"],
            color=DISCORD_CONFIG["colors"]["primary"]
        )
        
        intro_embed.set_footer(text=CONFESSION_CONSTANTS["messages"]["footers"]["intro"])
        
        # Create persistent view with the button
        view = SetupView()
        
        # Send introduction message to the confession channel
        await channel.send(embed=intro_embed, view=view)

        # Send confirmation to admin
        confirm_embed = discord.Embed(
            title=CONFESSION_CONSTANTS["messages"]["titles"]["setup_success"].format(emoji="✅"),
            description=CONFESSION_CONSTANTS["messages"]["descriptions"]["setup_success"].format(channel=channel.mention),
            color=DISCORD_CONFIG["colors"]["success"]
        )
        
        confirm_embed.add_field(
            name="📋 Channel Information",
            value=(
                f"**Channel:** {channel.mention}\n"
                f"**Channel ID:** `{channel.id}`\n"
                f"**Thread Support:** Enabled\n"
                f"**Anonymous Replies:** Enabled"
            ),
            inline=False
        )
        
        confirm_embed.add_field(
            name="📝 System Features",
            value=(
                "• Anonymous confession creation\n"
                "• Thread-based organization\n"
                "• Anonymous replies in threads\n"
                "• Unique ID system (CONF-001, REPLY-001-A)\n"
                "• Attachment support\n"
                "• Complete anonymity protection"
            ),
            inline=False
        )
        
        await interaction.response.send_message(embed=confirm_embed, ephemeral=True)
    
    @app_commands.command(name="confession-settings", description="View/modify confession settings (Admin only)")
    @app_commands.default_permissions(administrator=True)
    async def confession_settings(self, interaction: discord.Interaction):
        """View confession settings for the server."""
        settings = self.confession_service.get_guild_settings(interaction.guild.id)
        
        embed = discord.Embed(
            title=CONFESSION_CONSTANTS["messages"]["titles"]["settings"].format(emoji="⚙️"),
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
        embed.add_field(name="Next Confession ID", value=f"CONF-{settings.next_confession_id:03d}", inline=True)
        embed.add_field(name="Reply Letters", value=f"{len(settings.next_reply_letter)} confessions tracked", inline=True)
        
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
            title=CONFESSION_CONSTANTS["messages"]["titles"]["stats"].format(emoji="📊"),
            color=DISCORD_CONFIG["colors"]["info"]
        )
        
        embed.add_field(name="Total Confessions", value=total_confessions, inline=True)
        embed.add_field(name="Total Replies", value=total_replies, inline=True)
        embed.add_field(name="Average Replies", value=f"{total_replies/total_confessions:.1f}" if total_confessions > 0 else "0", inline=True)
        
        if confessions:
            latest = confessions[0]
            embed.add_field(name="Latest Confession", value=f"CONF-{latest.confession_id:03d}", inline=True)
        
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

        # Create embed with the specified format
        embed = discord.Embed(
            title=CONFESSION_CONSTANTS["messages"]["titles"]["confession"].format(emoji="📝", id=confession.confession_id),
            description=confession.content,
            color=DISCORD_CONFIG["colors"]["info"],
            timestamp=confession.created_at
        )

        # Add attachments if any
        if confession.attachments:
            # Filter image and non-image attachments
            image_attachments = []
            other_attachments = []
            
            for attachment in confession.attachments:
                if attachment.lower().endswith(('.png', '.jpg', '.jpeg', '.gif', '.webp')):
                    image_attachments.append(attachment)
                else:
                    other_attachments.append(attachment)
            
            # Set first image as main embed image
            if image_attachments:
                embed.set_image(url=image_attachments[0])
            
            # Add remaining images and other attachments as fields
            remaining_images = image_attachments[1:] if len(image_attachments) > 1 else []
            if remaining_images or other_attachments:
                attachment_links = []
                
                # Add remaining images as preview links
                for i, img_url in enumerate(remaining_images, 2):
                    attachment_links.append(f"[🖼️ Image {i}]({img_url})")
                
                # Add other attachments
                for i, att_url in enumerate(other_attachments, 1):
                    attachment_links.append(f"[📎 Attachment {i}]({att_url})")
                
                if attachment_links:
                    embed.add_field(name="📎 Additional Attachments", value="\n".join(attachment_links), inline=False)

        # Create view with buttons
        view = ConfessionView(confession.confession_id)
        
        # Register the view with the bot for persistence
        self.bot.add_view(view)

        try:
            # Post the confession message first
            message = await channel.send(embed=embed, view=view)
            
            # Send additional images as separate messages for better preview
            if confession.attachments:
                image_attachments = [att for att in confession.attachments if att.lower().endswith(('.png', '.jpg', '.jpeg', '.gif', '.webp'))]
                if len(image_attachments) > 1:
                    # Send remaining images as separate messages for preview
                    for img_url in image_attachments[1:]:
                        try:
                            # Create a simple embed for additional images
                            img_embed = discord.Embed(color=DISCORD_CONFIG["colors"]["info"])
                            img_embed.set_image(url=img_url)
                            await channel.send(embed=img_embed)
                        except Exception as e:
                            logger.warning(f"Failed to send additional image {img_url}: {e}")
            
            # Create a thread for this confession immediately
            thread_name = CONFESSION_CONSTANTS["thread"]["name_format"].format(
                emoji="💬", 
                id=confession.confession_id
            )
            thread = await message.create_thread(
                name=thread_name, 
                auto_archive_duration=CONFESSION_CONSTANTS["thread"]["auto_archive_duration"]
            )
            
            # Mark as posted with thread ID
            self.confession_service.mark_confession_posted(
                confession.confession_id,
                channel.id,
                message.id,
                thread.id
            )

            logger.info(CONFESSION_LOG_MESSAGES["confession"]["posted"].format(
                id=f"CONF-{confession.confession_id:03d}",
                guild_name=guild.name,
                thread_id=thread.id
            ))

        except Exception as e:
            logger.error(CONFESSION_LOG_MESSAGES["confession"]["failed_post"].format(error=str(e)))
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

        # Get the thread
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

        try:
            # Create reply embed with the specified format
            embed = discord.Embed(
                title=CONFESSION_CONSTANTS["messages"]["titles"]["reply"].format(emoji="↪️", reply_id=reply.reply_id),
                description=reply.content,
                color=DISCORD_CONFIG["colors"]["secondary"],
                timestamp=reply.created_at
            )
            
            # Add attachments if any
            if reply.attachments:
                # Filter image and non-image attachments
                image_attachments = []
                other_attachments = []
                
                for attachment in reply.attachments:
                    if attachment.lower().endswith(('.png', '.jpg', '.jpeg', '.gif', '.webp')):
                        image_attachments.append(attachment)
                    else:
                        other_attachments.append(attachment)
                
                # Set first image as main embed image
                if image_attachments:
                    embed.set_image(url=image_attachments[0])
                
                # Add remaining images and other attachments as fields
                remaining_images = image_attachments[1:] if len(image_attachments) > 1 else []
                if remaining_images or other_attachments:
                    attachment_links = []
                    
                    # Add remaining images as preview links
                    for i, img_url in enumerate(remaining_images, 2):
                        attachment_links.append(f"[🖼️ Image {i}]({img_url})")
                    
                    # Add other attachments
                    for i, att_url in enumerate(other_attachments, 1):
                        attachment_links.append(f"[📎 Attachment {i}]({att_url})")
                    
                    if attachment_links:
                        embed.add_field(name="📎 Additional Attachments", value="\n".join(attachment_links), inline=False)
            
            # Add Reply button to the reply message
            reply_view = ConfessionView(confession.confession_id)
            self.bot.add_view(reply_view)
            
            logger.info(f"Posting reply to thread {thread.name} (ID: {thread.id})")
            
            # Post reply embed in the thread
            message = await thread.send(embed=embed, view=reply_view)
            
            # Send additional images as separate messages for better preview
            if reply.attachments:
                image_attachments = [att for att in reply.attachments if att.lower().endswith(('.png', '.jpg', '.jpeg', '.gif', '.webp'))]
                if len(image_attachments) > 1:
                    # Send remaining images as separate messages for preview
                    for img_url in image_attachments[1:]:
                        try:
                            # Create a simple embed for additional images
                            img_embed = discord.Embed(color=DISCORD_CONFIG["colors"]["secondary"])
                            img_embed.set_image(url=img_url)
                            await thread.send(embed=img_embed)
                        except Exception as e:
                            logger.warning(f"Failed to send additional image {img_url}: {e}")

            # Mark as posted
            self.confession_service.mark_reply_posted(reply.reply_id, message.id)

            logger.info(CONFESSION_LOG_MESSAGES["reply"]["posted"].format(
                id=reply.reply_id,
                confession_id=reply.confession_id
            ))

        except discord.Forbidden:
            logger.error(f"No permission to post in thread {thread.id}")
        except Exception as e:
            logger.error(CONFESSION_LOG_MESSAGES["reply"]["failed_post"].format(error=str(e)))
            import traceback
            logger.error(traceback.format_exc())


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(ConfessionCog(bot))
