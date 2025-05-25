"""
Configuration settings for NeruBot
"""
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class Settings:
    """Bot configuration settings."""
    
    # Discord settings
    DISCORD_TOKEN = os.getenv('DISCORD_TOKEN')
    COMMAND_PREFIX = os.getenv('COMMAND_PREFIX', '!')
    
    # Bot settings
    BOT_NAME = "NeruBot"
    BOT_VERSION = "2.0.0"
    BOT_DESCRIPTION = "A simple and extensible Discord bot"
    
    # Music settings
    MAX_QUEUE_SIZE = int(os.getenv('MAX_QUEUE_SIZE', '50'))
    DEFAULT_VOLUME = float(os.getenv('DEFAULT_VOLUME', '0.5'))
    
    # External API keys (optional)
    WEATHER_API_KEY = os.getenv('WEATHER_API_KEY')
    
    # Logging settings
    LOG_LEVEL = os.getenv('LOG_LEVEL', 'INFO')
    LOG_FILE = os.getenv('LOG_FILE', 'bot.log')
    
    @classmethod
    def validate(cls):
        """Validate required settings."""
        if not cls.DISCORD_TOKEN:
            raise ValueError("DISCORD_TOKEN is required in .env file")
        
        return True

# Create global settings instance
settings = Settings()
