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
        
        # Create confession
        success, message, confession = cog.confession_service.create_confession(
            content=content,
            author_id=interaction.user.id,
            guild_id=interaction.guild.id,
            image_url=self.image_url
        )
        
        if not success:
            await interaction.response.send_message(f"❌ {message}", ephemeral=True)
            return
        
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
        
        # Create reply
        success, message, reply = cog.confession_service.create_reply(
            confession_id=self.confession_id,
            content=content,
            author_id=interaction.user.id,
            guild_id=interaction.guild.id,
            image_url=self.image_url
        )
        
        if not success:
            await interaction.response.send_message(f"❌ {message}", ephemeral=True)
            return
        
        # Post reply to channel
        await cog.post_reply(reply, interaction.guild)
        
        await interaction.response.send_message(
            f"✅ Your reply has been posted anonymously!",
            ephemeral=True
        )


class ConfessionView(ui.View):
    """View with buttons for confession interactions."""
    
    def __init__(self, confession_id: int):
        super().__init__(timeout=None)  # Persistent view
        self.confession_id = confession_id

    @ui.button(label="Reply Anonymously", style=discord.ButtonStyle.secondary, emoji="💬", custom_id="confession_reply")
    async def reply_button(self, interaction: discord.Interaction, button: ui.Button):
        # Get confession ID from the embed
        if interaction.message and interaction.message.embeds:
            embed = interaction.message.embeds[0]
            title = embed.title or ""
            # Extract confession ID from title like "📝 Anonymous Confession #123"
            import re
            match = re.search(r'#(\d+)', title)
            if match:
                confession_id = int(match.group(1))
                modal = ReplyModal(confession_id)
                await interaction.response.send_modal(modal)
            else:
                await interaction.response.send_message("❌ Could not find confession ID.", ephemeral=True)
        else:
            await interaction.response.send_message("❌ Could not find confession information.", ephemeral=True)

    @ui.button(label="View Replies", style=discord.ButtonStyle.primary, emoji="👁️", custom_id="confession_view_replies")
    async def view_replies_button(self, interaction: discord.Interaction, button: ui.Button):
        cog = interaction.client.get_cog("ConfessionCog")
        if not cog:
            await interaction.response.send_message("❌ Confession system is not available.", ephemeral=True)
            return

        # Get confession ID from the embed
        if interaction.message and interaction.message.embeds:
            embed = interaction.message.embeds[0]
            title = embed.title or ""
            # Extract confession ID from title like "📝 Anonymous Confession #123"
            import re
            match = re.search(r'#(\d+)', title)
            if match:
                confession_id = int(match.group(1))
            else:
                await interaction.response.send_message("❌ Could not find confession ID.", ephemeral=True)
                return
        else:
            await interaction.response.send_message("❌ Could not find confession information.", ephemeral=True)
            return

        replies = cog.confession_service.get_confession_replies(confession_id)

        if not replies:
            await interaction.response.send_message("📭 No replies yet for this confession.", ephemeral=True)
            return

        embed = discord.Embed(
            title=f"💬 Replies to Confession #{confession_id}",
            color=DISCORD_CONFIG["colors"]["info"],
            description=f"**{len(replies)}** anonymous replies:"
        )

        for i, reply in enumerate(replies[-5:], 1):  # Show last 5 replies
            embed.add_field(
                name=f"Reply #{i}",
                value=reply.content[:200] + ("..." if len(reply.content) > 200 else ""),
                inline=False
            )

        if len(replies) > 5:
            embed.set_footer(text=f"Showing latest 5 of {len(replies)} replies")

        await interaction.response.send_message(embed=embed, ephemeral=True)


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
        # Add persistent view for confession interactions
        self.bot.add_view(ConfessionView(0))  # ID doesn't matter since we extract it dynamically
        # Add persistent view for setup success button
        self.bot.add_view(SetupSuccessView())
    
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

        embed.set_footer(text="React or use buttons to interact anonymously")

        # Create view with buttons
        view = ConfessionView(confession.confession_id)

        try:
            message = await channel.send(embed=embed, view=view)
            
            # Create a thread for this confession
            thread_name = f"Confession #{confession.confession_id} Replies"
            thread = await message.create_thread(name=thread_name, auto_archive_duration=10080)  # 7 days
            
            # Post a welcome message in the thread
            welcome_embed = discord.Embed(
                title="💬 Reply to this Confession",
                description="Use the button on the original message or the `/reply` command to post anonymous replies here.",
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
    
    async def post_reply(self, reply: ConfessionReply, guild: discord.Guild):
        """Post a reply to the confession thread."""
        settings = self.confession_service.get_guild_settings(guild.id)

        if not settings.confession_channel_id:
            logger.error(f"No confession channel set for guild {guild.id}")
            return

        # Get the confession to find the thread ID
        confession = self.confession_service.get_confession(reply.confession_id)
        if not confession or not confession.thread_id:
            logger.error(f"No thread found for confession {reply.confession_id}")
            return

        # Get the thread
        thread = self.bot.get_channel(confession.thread_id)
        if not thread:
            logger.error(f"Thread {confession.thread_id} not found")
            return

        # Create embed
        embed = discord.Embed(
            title=f"💬 Anonymous Reply",
            description=reply.content,
            color=DISCORD_CONFIG["colors"]["warning"],
            timestamp=reply.created_at
        )

        if reply.image_url:
            embed.set_image(url=reply.image_url)

        embed.set_footer(text=f"Reply to Confession #{reply.confession_id}")

        try:
            message = await thread.send(embed=embed)

            # Mark as posted
            self.confession_service.mark_reply_posted(reply.reply_id, message.id)

            logger.info(f"Posted reply {reply.reply_id} to confession {reply.confession_id} thread")

        except Exception as e:
            logger.error(f"Error posting reply: {e}")


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(ConfessionCog(bot))
