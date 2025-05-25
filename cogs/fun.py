"""
Fun commands cog - entertainment features
"""
import discord
from discord.ext import commands
from discord import app_commands
import random
import aiohttp
import json

class Fun(commands.Cog):
    """Fun and entertainment commands."""
    
    def __init__(self, bot):
        self.bot = bot
    
    @app_commands.command(name="roll", description="Roll a dice")
    @app_commands.describe(sides="Number of sides on the dice (default: 6)")
    async def roll(self, interaction: discord.Interaction, sides: int = 6):
        """Roll a dice with specified number of sides."""
        if sides < 2:
            await interaction.response.send_message("❌ Dice must have at least 2 sides!")
            return
        
        if sides > 100:
            await interaction.response.send_message("❌ That's too many sides! Maximum is 100.")
            return
        
        result = random.randint(1, sides)
        
        embed = discord.Embed(
            title="🎲 Dice Roll",
            description=f"You rolled a **{result}** out of {sides}!",
            color=discord.Color.random()
        )
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="coinflip", description="Flip a coin")
    async def coinflip(self, interaction: discord.Interaction):
        """Flip a coin."""
        result = random.choice(["Heads", "Tails"])
        emoji = "🪙" if result == "Heads" else "🥈"
        
        embed = discord.Embed(
            title=f"{emoji} Coin Flip",
            description=f"The coin landed on **{result}**!",
            color=discord.Color.gold() if result == "Heads" else discord.Color.light_grey()
        )
        
        await interaction.response.send_message(embed=embed)
    
    @app_commands.command(name="joke", description="Get a random joke")
    async def joke(self, interaction: discord.Interaction):
        """Get a random joke from an API."""
        await interaction.response.defer()
        
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get("https://official-joke-api.appspot.com/random_joke") as response:
                    if response.status == 200:
                        data = await response.json()
                        
                        embed = discord.Embed(
                            title="😂 Random Joke",
                            color=discord.Color.yellow()
                        )
                        embed.add_field(name="Setup", value=data['setup'], inline=False)
                        embed.add_field(name="Punchline", value=data['punchline'], inline=False)
                        
                        await interaction.followup.send(embed=embed)
                    else:
                        await interaction.followup.send("❌ Couldn't fetch a joke right now. Try again later!")
        except Exception as e:
            # Fallback to local jokes if API fails
            local_jokes = [
                ("Why don't scientists trust atoms?", "Because they make up everything!"),
                ("What do you call a fake noodle?", "An impasta!"),
                ("Why did the scarecrow win an award?", "He was outstanding in his field!"),
                ("What do you call a sleeping bull?", "A bulldozer!"),
                ("Why don't eggs tell jokes?", "They'd crack each other up!")
            ]
            
            setup, punchline = random.choice(local_jokes)
            
            embed = discord.Embed(
                title="😂 Random Joke",
                color=discord.Color.yellow()
            )
            embed.add_field(name="Setup", value=setup, inline=False)
            embed.add_field(name="Punchline", value=punchline, inline=False)
            embed.set_footer(text="(Offline joke - API unavailable)")
            
            await interaction.followup.send(embed=embed)
    
    @app_commands.command(name="8ball", description="Ask the magic 8-ball a question")
    @app_commands.describe(question="Your question for the magic 8-ball")
    async def eightball(self, interaction: discord.Interaction, question: str):
        """Ask the magic 8-ball a question."""
        responses = [
            "It is certain", "It is decidedly so", "Without a doubt",
            "Yes definitely", "You may rely on it", "As I see it, yes",
            "Most likely", "Outlook good", "Yes", "Signs point to yes",
            "Reply hazy, try again", "Ask again later", "Better not tell you now",
            "Cannot predict now", "Concentrate and ask again",
            "Don't count on it", "My reply is no", "My sources say no",
            "Outlook not so good", "Very doubtful"
        ]
        
        response = random.choice(responses)
        
        embed = discord.Embed(
            title="🎱 Magic 8-Ball",
            color=discord.Color.purple()
        )
        embed.add_field(name="Question", value=question, inline=False)
        embed.add_field(name="Answer", value=f"*{response}*", inline=False)
        
        await interaction.response.send_message(embed=embed)

async def setup(bot):
    await bot.add_cog(Fun(bot))
