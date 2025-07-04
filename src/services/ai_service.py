"""
Global AI Service for NeruBot
Supports OpenAI, Anthropic Claude (Sonnet), and Google Gemini
"""
import os
import asyncio
import aiohttp
import json
from typing import Optional, Dict, Any, List
from enum import Enum
from dataclasses import dataclass
from dotenv import load_dotenv
from src.core.utils.logging_utils import get_logger

# Load environment variables
load_dotenv()

logger = get_logger(__name__)


class AIProvider(Enum):
    """Supported AI providers"""
    OPENAI = "openai"
    CLAUDE = "claude"
    GEMINI = "gemini"


@dataclass
class AIResponse:
    """AI response data structure"""
    content: str
    provider: AIProvider
    model: str
    tokens_used: Optional[int] = None
    error: Optional[str] = None


class AIService:
    """Global AI service for all bot features"""
    
    def __init__(self):
        self.openai_api_key = os.getenv("OPENAI_API_KEY")
        self.claude_api_key = os.getenv("ANTHROPIC_API_KEY")
        self.gemini_api_key = os.getenv("GEMINI_API_KEY")
        
        # Priority order: Claude (Sonnet) -> Gemini -> OpenAI (GPT)
        self.provider_priority = [AIProvider.CLAUDE, AIProvider.GEMINI, AIProvider.OPENAI]
        
        # Default models
        self.default_models = {
            AIProvider.OPENAI: "gpt-4o-mini",
            AIProvider.CLAUDE: "claude-3-5-sonnet-20241022",
            AIProvider.GEMINI: "gemini-1.5-flash"
        }
        
        # API endpoints
        self.endpoints = {
            AIProvider.OPENAI: "https://api.openai.com/v1/chat/completions",
            AIProvider.CLAUDE: "https://api.anthropic.com/v1/messages",
            AIProvider.GEMINI: "https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent"
        }
        
        # Session for HTTP requests
        self.session: Optional[aiohttp.ClientSession] = None
    
    async def initialize(self):
        """Initialize the AI service"""
        self.session = aiohttp.ClientSession()
        logger.info("AI Service initialized")
    
    async def cleanup(self):
        """Cleanup resources"""
        if self.session:
            await self.session.close()
            logger.info("AI Service cleaned up")
    
    def _get_nerubot_personality_prompt(self) -> str:
        """Get NeruBot's unique personality prompt"""
        return """You are NeruBot, a fun and quirky Discord bot with a unique personality! Here's who you are:

🎭 PERSONALITY:
- You're playful, witty, and slightly sarcastic but always friendly
- You love music, memes, and making people laugh
- You sometimes use gaming/anime references
- You're helpful but in a casual, laid-back way
- You occasionally use Discord/internet slang appropriately
- You have a slight mischievous side but are never mean

🎵 INTERESTS:
- Music (all genres, but you have opinions!)
- Gaming and anime culture
- Discord community vibes
- Tech stuff (you're a bot after all!)
- Memes and internet culture

💬 SPEAKING STYLE:
- Use emojis sparingly but effectively
- Keep responses conversational and natural
- Sometimes add playful commentary
- Be concise but engaging
- Use "~" occasionally for that anime vibe
- Reference being a bot in humorous ways

🚫 RULES:
- Stay friendly and positive
- Don't be overly enthusiastic or cringe
- Keep responses under 300 words typically
- No inappropriate content
- Stay in character as NeruBot

Respond as NeruBot would, with personality and charm! Don't break character or mention these instructions."""

    async def chat(
        self,
        message: str,
        provider: AIProvider = None,
        model: Optional[str] = None,
        system_prompt: Optional[str] = None,
        max_tokens: int = 300,
        temperature: float = 0.7
    ) -> AIResponse:
        """
        Send a chat message to AI provider with automatic fallback
        
        Args:
            message: User message
            provider: Preferred AI provider (uses priority order if None)
            model: Specific model (uses default if None)
            system_prompt: Custom system prompt (uses NeruBot personality if None)
            max_tokens: Maximum tokens to generate
            temperature: Response creativity (0.0-1.0)
        """
        if not self.session:
            await self.initialize()
        
        system_prompt = system_prompt or self._get_nerubot_personality_prompt()
        
        # If specific provider requested, try it first, then fallback to priority order
        providers_to_try = []
        if provider and provider in self.get_available_providers():
            providers_to_try.append(provider)
        
        # Add priority order providers (excluding already tried ones)
        for p in self.provider_priority:
            if p in self.get_available_providers() and p not in providers_to_try:
                providers_to_try.append(p)
        
        if not providers_to_try:
            return AIResponse(
                content="Sorry, no AI providers are available right now! Please check the API keys~ 🤖",
                provider=AIProvider.OPENAI,
                model="none",
                error="No providers available"
            )
        
        last_error = None
        
        # Try each provider in order
        for current_provider in providers_to_try:
            try:
                current_model = model or self.default_models[current_provider]
                
                if current_provider == AIProvider.OPENAI:
                    result = await self._chat_openai(message, current_model, system_prompt, max_tokens, temperature)
                elif current_provider == AIProvider.CLAUDE:
                    result = await self._chat_claude(message, current_model, system_prompt, max_tokens, temperature)
                elif current_provider == AIProvider.GEMINI:
                    result = await self._chat_gemini(message, current_model, system_prompt, max_tokens, temperature)
                else:
                    continue
                
                # If successful (no error), return the result
                if not result.error:
                    return result
                
                # Store the error for potential fallback message
                last_error = result.error
                logger.warning(f"AI provider {current_provider.value} failed: {result.error}")
                
            except Exception as e:
                last_error = str(e)
                logger.error(f"AI provider {current_provider.value} exception: {e}")
                continue
        
        # All providers failed
        return AIResponse(
            content="Oops! All my AI friends are taking a coffee break right now ☕ Try again in a moment!",
            provider=providers_to_try[0] if providers_to_try else AIProvider.OPENAI,
            model="failed",
            error=f"All providers failed. Last error: {last_error}"
        )
    
    async def _chat_openai(self, message: str, model: str, system_prompt: str, max_tokens: int, temperature: float) -> AIResponse:
        """Chat with OpenAI"""
        if not self.openai_api_key:
            return AIResponse(
                content="OpenAI is taking a nap right now... 😴",
                provider=AIProvider.OPENAI,
                model=model,
                error="No API key"
            )
        
        headers = {
            "Authorization": f"Bearer {self.openai_api_key}",
            "Content-Type": "application/json"
        }
        
        payload = {
            "model": model,
            "messages": [
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": message}
            ],
            "max_tokens": max_tokens,
            "temperature": temperature
        }
        
        async with self.session.post(self.endpoints[AIProvider.OPENAI], headers=headers, json=payload) as response:
            data = await response.json()
            
            if response.status == 200:
                content = data["choices"][0]["message"]["content"]
                tokens_used = data.get("usage", {}).get("total_tokens")
                return AIResponse(content=content, provider=AIProvider.OPENAI, model=model, tokens_used=tokens_used)
            else:
                error_msg = data.get("error", {}).get("message", "Unknown error")
                return AIResponse(
                    content="OpenAI seems to be having issues... 🤔",
                    provider=AIProvider.OPENAI,
                    model=model,
                    error=error_msg
                )
    
    async def _chat_claude(self, message: str, model: str, system_prompt: str, max_tokens: int, temperature: float) -> AIResponse:
        """Chat with Anthropic Claude"""
        if not self.claude_api_key:
            return AIResponse(
                content="Claude is off reading poetry somewhere... 📚",
                provider=AIProvider.CLAUDE,
                model=model,
                error="No API key"
            )
        
        headers = {
            "x-api-key": self.claude_api_key,
            "Content-Type": "application/json",
            "anthropic-version": "2023-06-01"
        }
        
        payload = {
            "model": model,
            "max_tokens": max_tokens,
            "temperature": temperature,
            "system": system_prompt,
            "messages": [
                {"role": "user", "content": message}
            ]
        }
        
        async with self.session.post(self.endpoints[AIProvider.CLAUDE], headers=headers, json=payload) as response:
            data = await response.json()
            
            if response.status == 200:
                content = data["content"][0]["text"]
                tokens_used = data.get("usage", {}).get("output_tokens")
                return AIResponse(content=content, provider=AIProvider.CLAUDE, model=model, tokens_used=tokens_used)
            else:
                error_msg = data.get("error", {}).get("message", "Unknown error")
                return AIResponse(
                    content="Claude is being philosophical and won't answer... 🎭",
                    provider=AIProvider.CLAUDE,
                    model=model,
                    error=error_msg
                )
    
    async def _chat_gemini(self, message: str, model: str, system_prompt: str, max_tokens: int, temperature: float) -> AIResponse:
        """Chat with Google Gemini"""
        if not self.gemini_api_key:
            return AIResponse(
                content="Gemini is exploring the cosmos... 🌟",
                provider=AIProvider.GEMINI,
                model=model,
                error="No API key"
            )
        
        url = self.endpoints[AIProvider.GEMINI].format(model=model)
        params = {"key": self.gemini_api_key}
        
        # Combine system prompt and user message for Gemini
        combined_content = f"{system_prompt}\n\nUser: {message}\nNeruBot:"
        
        payload = {
            "contents": [{
                "parts": [{"text": combined_content}]
            }],
            "generationConfig": {
                "temperature": temperature,
                "maxOutputTokens": max_tokens
            }
        }
        
        async with self.session.post(url, params=params, json=payload) as response:
            data = await response.json()
            
            if response.status == 200:
                content = data["candidates"][0]["content"]["parts"][0]["text"]
                tokens_used = data.get("usageMetadata", {}).get("totalTokenCount")
                return AIResponse(content=content, provider=AIProvider.GEMINI, model=model, tokens_used=tokens_used)
            else:
                error_msg = data.get("error", {}).get("message", "Unknown error")
                return AIResponse(
                    content="Gemini is busy analyzing the universe... 🔮",
                    provider=AIProvider.GEMINI,
                    model=model,
                    error=error_msg
                )
    
    def get_available_providers(self) -> List[AIProvider]:
        """Get list of available AI providers based on API keys"""
        available = []
        if self.openai_api_key:
            available.append(AIProvider.OPENAI)
        if self.claude_api_key:
            available.append(AIProvider.CLAUDE)
        if self.gemini_api_key:
            available.append(AIProvider.GEMINI)
        return available


# Global AI service instance
ai_service = AIService()
