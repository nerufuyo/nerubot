"""
Help commands cog with improved pagination and feature-based organization
"""
import discord
from discord import ui, ButtonStyle, Interaction
from discord.ext import commands
from discord import app_commands
from typing import List, Dict, Any, Optional
from src.config.messages import MSG_HELP, CMD_DESCRIPTIONS
from src.config.settings import DISCORD_CONFIG


class HelpView(ui.View):
    """Custom view for paginated help commands."""
    
    def __init__(self, pages: List[discord.Embed], timeout: int = 60):
        super().__init__(timeout=timeout)
        self.pages = pages
        self.current_page = 0
        self.message: Optional[discord.Message] = None
    
    async def on_timeout(self):
        """Remove buttons when view times out."""
        if self.message:
            for item in self.children:
                item.disabled = True
            await self.message.edit(view=self)
    
    @ui.button(emoji="⬅️", style=ButtonStyle.primary)
    async def previous_button(self, interaction: Interaction, button: ui.Button):
        """Handle previous page button."""
        if self.current_page > 0:
            self.current_page -= 1
        else:
            self.current_page = len(self.pages) - 1
        
        await interaction.response.edit_message(embed=self.pages[self.current_page], view=self)
    
    @ui.button(emoji="❌", style=ButtonStyle.danger)
    async def close_button(self, interaction: Interaction, button: ui.Button):
        """Handle close button."""
        self.stop()
        await interaction.message.delete()
    
    @ui.button(emoji="➡️", style=ButtonStyle.primary)
    async def next_button(self, interaction: Interaction, button: ui.Button):
        """Handle next page button."""
        if self.current_page < len(self.pages) - 1:
            self.current_page += 1
        else:
            self.current_page = 0
        
        await interaction.response.edit_message(embed=self.pages[self.current_page], view=self)


class HelpCog(commands.Cog):
    """Help commands cog with pagination and feature categories."""
    
    def __init__(self, bot: commands.Bot):
        self.bot = bot
        
        # Define all commands by category
        self.command_categories = {
            "🎵 Music - Playback": [
                ("/play <query>", MSG_HELP["commands"]["play"]),
                ("/pause", MSG_HELP["commands"]["pause"]),
                ("/resume", MSG_HELP["commands"]["resume"]),
                ("/stop", MSG_HELP["commands"]["stop"]),
                ("/skip", MSG_HELP["commands"]["skip"]),
            ],
            "🎵 Music - Voice": [
                ("/join", MSG_HELP["commands"]["join"]),
                ("/leave", MSG_HELP["commands"]["leave"]),
                ("/volume <0-100>", MSG_HELP["commands"]["volume"]),
            ],
            "🎵 Music - Queue": [
                ("/queue", MSG_HELP["commands"]["queue"]),
                ("/nowplaying", MSG_HELP["commands"]["nowplaying"]),
                ("/clear", MSG_HELP["commands"]["clear"]),
                ("/loop <off/song/queue>", MSG_HELP["commands"]["loop"]),
                ("/247", MSG_HELP["commands"]["247"]),
            ],
            "🎵 Music - Info": [
                ("/sources", MSG_HELP["commands"]["sources"]),
            ],
            "📝 Confessions": [
                ("/confess", MSG_HELP["commands"]["confess"]),
                ("/reply <confession_id>", MSG_HELP["commands"]["reply"]),
                ("/confession-stats", MSG_HELP["commands"]["confession-stats"]),
            ],
            "📝 Confession Admin": [
                ("/confession-setup <channel>", MSG_HELP["commands"]["confession-setup"]),
                ("/confession-settings", MSG_HELP["commands"]["confession-settings"]),
            ],
            "📰 News": [
                ("/news latest [count]", MSG_HELP["commands"]["news-latest"]),
                ("/news sources", MSG_HELP["commands"]["news-sources"]),
                ("/news status", MSG_HELP["commands"]["news-status"]),
                ("/news help", MSG_HELP["commands"]["news-help"]),
            ],
            "📰 News Admin": [
                ("/news set-channel [channel]", MSG_HELP["commands"]["news-set-channel"]),
                ("/news start", MSG_HELP["commands"]["news-start"]),
                ("/news stop", MSG_HELP["commands"]["news-stop"]),
                ("/news add <name> <url>", MSG_HELP["commands"]["news-add"]),
                ("/news remove <name>", MSG_HELP["commands"]["news-remove"]),
            ],
            "🤖 General": [
                ("/help", MSG_HELP["commands"]["help"]),
                ("/about", MSG_HELP["commands"]["about"]),
                ("/features", MSG_HELP["commands"]["features"]),
                ("/commands", MSG_HELP["commands"]["commands"]),
            ]
        }

    @app_commands.command(name="help", description="Show help information with categories")
    async def help_command(self, interaction: discord.Interaction) -> None:
        """Show paginated help information organized by feature."""
        # Create help pages
        pages = self._create_help_pages()
        
        # Create the view with pagination
        view = HelpView(pages)
        
        # Send the first page
        await interaction.response.send_message(embed=pages[0], view=view)
        
        # Store the message for timeout handling
        view.message = await interaction.original_response()
    
    def _create_help_pages(self) -> List[discord.Embed]:
        """Create help pages organized by category."""
        pages = []
        
        # Main help page
        main_embed = discord.Embed(
            title="🤖 NeruBot Help Menu",
            description=MSG_HELP["main_description"],
            color=DISCORD_CONFIG["colors"]["info"]
        )
        
        main_embed.set_thumbnail(url="https://i.imgur.com/4M34hi2.png")
        main_embed.set_footer(text=f"Page 1 of {len(self.command_categories) + 1} | Main Help Menu")
        pages.append(main_embed)
        
        # Category pages
        page_num = 2
        for category, commands in self.command_categories.items():
            embed = discord.Embed(
                title=f"{category} Commands",
                description="Here are the commands available in this category:",
                color=DISCORD_CONFIG["colors"]["info"]
            )
            
            commands_text = "\n".join([f"**{cmd}**: {desc}" for cmd, desc in commands])
            embed.add_field(name="Commands", value=commands_text, inline=False)
            
            embed.set_footer(text=f"Page {page_num} of {len(self.command_categories) + 1} | {category}")
            pages.append(embed)
            page_num += 1
        
        return pages


async def setup(bot: commands.Bot) -> None:
    """Setup function to add the cog to the bot."""
    await bot.add_cog(HelpCog(bot))
