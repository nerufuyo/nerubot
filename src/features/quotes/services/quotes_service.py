"""
Quotes service for AI-powered quote generation using DeepSeek API
"""
import os
import logging
import random
from typing import List, Optional
from openai import AsyncOpenAI
from src.features.quotes.models.quote import Quote, QuoteRequest

logger = logging.getLogger(__name__)


class QuotesService:
    """Service for generating AI quotes using DeepSeek API."""
    
    def __init__(self, api_key: Optional[str] = None):
        self.api_key = api_key or os.getenv("DEEPSEEK_API_KEY")
        
        if self.api_key:
            self.client = AsyncOpenAI(
                api_key=self.api_key,
                base_url="https://api.deepseek.com"
            )
            logger.info("QuotesService initialized with DeepSeek API")
        else:
            self.client = None
            logger.warning("QuotesService initialized without API key - using fallback quotes")
        
        # Fallback quotes for when API is not available
        self.fallback_quotes = {
            "motivation": [
                ("The only way to do great work is to love what you do.", "Steve Jobs"),
                ("Innovation distinguishes between a leader and a follower.", "Steve Jobs"),
                ("Your limitation—it's only your imagination.", "Unknown"),
                ("Success is not final, failure is not fatal: it is the courage to continue that counts.", "Winston Churchill"),
            ],
            "wisdom": [
                ("The journey of a thousand miles begins with a single step.", "Lao Tzu"),
                ("In the middle of difficulty lies opportunity.", "Albert Einstein"),
                ("The only true wisdom is in knowing you know nothing.", "Socrates"),
                ("Yesterday is history, tomorrow is a mystery, today is a gift.", "Eleanor Roosevelt"),
            ],
            "technology": [
                ("Technology is best when it brings people together.", "Matt Mullenweg"),
                ("The advance of technology is based on making it fit in so that you don't really even notice it.", "Bill Gates"),
                ("Any sufficiently advanced technology is indistinguishable from magic.", "Arthur C. Clarke"),
                ("The real problem is not whether machines think but whether men do.", "B.F. Skinner"),
            ],
            "philosophy": [
                ("I think, therefore I am.", "René Descartes"),
                ("The unexamined life is not worth living.", "Socrates"),
                ("Be yourself; everyone else is already taken.", "Oscar Wilde"),
                ("Two things are infinite: the universe and human stupidity; and I'm not sure about the universe.", "Albert Einstein"),
            ],
            "humor": [
                ("I have not failed. I've just found 10,000 ways that won't work.", "Thomas Edison"),
                ("A day without sunshine is like, you know, night.", "Steve Martin"),
                ("I'm not arguing, I'm just explaining why I'm right.", "Anonymous"),
                ("Life is what happens to you while you're busy making other plans.", "John Lennon"),
            ]
        }
    
    async def get_quote(self, request: QuoteRequest) -> Quote:
        """Get a quote based on the request parameters."""
        if self.client and self.api_key:
            try:
                return await self._generate_ai_quote(request)
            except Exception as e:
                logger.warning(f"AI quote generation failed: {e}, falling back to predefined quotes")
        
        # Fallback to predefined quotes
        return self._get_fallback_quote(request)
    
    async def _generate_ai_quote(self, request: QuoteRequest) -> Quote:
        """Generate a quote using DeepSeek AI."""
        # Build the prompt based on request parameters
        prompt_parts = ["Generate an inspiring and thoughtful quote"]
        
        if request.category:
            prompt_parts.append(f"about {request.category}")
        
        if request.mood:
            prompt_parts.append(f"with a {request.mood} tone")
        
        length_guide = {
            "short": "Keep it concise (under 20 words)",
            "medium": "Make it moderate length (20-40 words)",
            "long": "Make it longer and more detailed (40-80 words)"
        }
        
        if request.length in length_guide:
            prompt_parts.append(length_guide[request.length])
        
        prompt_parts.extend([
            "Return ONLY the quote text without quotation marks.",
            "Make it original, meaningful, and inspirational.",
            "Do not include author attribution."
        ])
        
        prompt = ". ".join(prompt_parts) + "."
        
        try:
            response = await self.client.chat.completions.create(
                model="deepseek-chat",
                messages=[
                    {"role": "system", "content": "You are a wise and inspirational quote generator. Create original, meaningful quotes that inspire and motivate people."},
                    {"role": "user", "content": prompt}
                ],
                max_tokens=150,
                temperature=0.8,
            )
            
            quote_text = response.choices[0].message.content.strip()
            
            # Clean up the quote (remove any remaining quotes or extra formatting)
            quote_text = quote_text.strip('"\'')
            
            return Quote(
                content=quote_text,
                author="DeepSeek AI",
                category=request.category or "inspiration",
                source="deepseek_ai"
            )
            
        except Exception as e:
            logger.error(f"Error generating AI quote: {e}")
            raise
    
    def _get_fallback_quote(self, request: QuoteRequest) -> Quote:
        """Get a fallback quote when AI is not available."""
        category = request.category or "motivation"
        
        # Get quotes from the requested category, or default to motivation
        available_quotes = self.fallback_quotes.get(category, self.fallback_quotes["motivation"])
        
        # Select a random quote
        quote_text, author = random.choice(available_quotes)
        
        return Quote(
            content=quote_text,
            author=author,
            category=category,
            source="fallback"
        )
    
    async def get_random_quote(self) -> Quote:
        """Get a random inspirational quote."""
        categories = ["motivation", "wisdom", "philosophy", "technology"]
        random_category = random.choice(categories)
        
        request = QuoteRequest(category=random_category)
        return await self.get_quote(request)
    
    async def get_quotes_by_category(self, category: str, count: int = 5) -> List[Quote]:
        """Get multiple quotes by category."""
        quotes = []
        
        for i in range(count):
            request = QuoteRequest(category=category)
            try:
                quote = await self.get_quote(request)
                quotes.append(quote)
            except Exception as e:
                logger.warning(f"Failed to generate quote {i+1}: {e}")
                # Add a fallback quote if AI fails
                fallback_quote = self._get_fallback_quote(request)
                quotes.append(fallback_quote)
        
        return quotes
    
    def get_available_categories(self) -> List[str]:
        """Get list of available quote categories."""
        return list(self.fallback_quotes.keys())
    
    def get_available_moods(self) -> List[str]:
        """Get list of available mood options."""
        return ["inspiring", "thoughtful", "humorous", "philosophical", "motivational", "calm", "energetic"]
