"""
Demo: Using the Global AI Service in Other Features
This example shows how to integrate the AI service into other parts of NeruBot
"""
import asyncio
import os
from dotenv import load_dotenv
from src.services.ai_service import ai_service, AIProvider

# Load environment variables
load_dotenv()


async def demo_music_ai_integration():
    """Demo: AI-powered music recommendations"""
    print("🎵 Demo: AI Music Recommendations")
    
    await ai_service.initialize()
    
    # Custom prompt for music recommendations
    music_prompt = """You are NeruBot's music expert! You have great taste in music and love helping people discover new songs. 
    Give music recommendations based on the user's request. Be enthusiastic about music and include specific song/artist suggestions.
    Keep responses under 200 words and use your fun personality."""
    
    user_request = "I like electronic music and anime soundtracks. Any recommendations?"
    
    response = await ai_service.chat(
        message=user_request,
        system_prompt=music_prompt,
        max_tokens=200,
        temperature=0.8
    )
    
    print(f"User: {user_request}")
    print(f"NeruBot: {response.content}")
    
    await ai_service.cleanup()


async def demo_news_ai_integration():
    """Demo: AI-powered news summarization"""
    print("\n📰 Demo: AI News Summarization")
    
    await ai_service.initialize()
    
    # Custom prompt for news summarization
    news_prompt = """You are NeruBot's news analyst! Summarize news articles in a fun, accessible way while keeping the important facts. 
    Add your own commentary and make it engaging. Use your witty personality but stay informative."""
    
    fake_news_article = """
    Tech Giant Announces New AI Chip
    
    A major technology company today announced the release of their next-generation AI processing chip, 
    claiming 10x performance improvements over previous models. The chip will power everything from 
    smartphones to data centers, with the company expecting to ship 50 million units in the first year.
    Stock prices jumped 15% following the announcement.
    """
    
    response = await ai_service.chat(
        message=f"Summarize this news article: {fake_news_article}",
        system_prompt=news_prompt,
        max_tokens=150
    )
    
    print(f"Original Article: {fake_news_article[:100]}...")
    print(f"NeruBot Summary: {response.content}")
    
    await ai_service.cleanup()


async def demo_confession_ai_integration():
    """Demo: AI-powered confession response generation"""
    print("\n💭 Demo: AI Confession Responses")
    
    await ai_service.initialize()
    
    # Custom prompt for confession responses
    confession_prompt = """You are NeruBot helping with anonymous confessions. Generate supportive, 
    empathetic responses to confessions. Be understanding and offer gentle advice when appropriate. 
    Keep your fun personality but be more caring and thoughtful."""
    
    confession = "I've been feeling really anxious about starting college next month. What if I don't make any friends?"
    
    response = await ai_service.chat(
        message=f"Someone confessed: '{confession}' - Generate a supportive response",
        system_prompt=confession_prompt,
        max_tokens=180,
        temperature=0.7  # Slightly less random for more thoughtful responses
    )
    
    print(f"Confession: {confession}")
    print(f"NeruBot Response: {response.content}")
    
    await ai_service.cleanup()


async def demo_help_ai_integration():
    """Demo: AI-powered help system"""
    print("\n❓ Demo: AI Help Generation")
    
    await ai_service.initialize()
    
    # Custom prompt for help responses
    help_prompt = """You are NeruBot's help system! Explain bot features and commands in a fun, 
    easy-to-understand way. Be helpful while maintaining your playful personality. 
    Include practical examples and tips."""
    
    help_question = "How do I use the music commands?"
    
    response = await ai_service.chat(
        message=help_question,
        system_prompt=help_prompt,
        max_tokens=200
    )
    
    print(f"Question: {help_question}")
    print(f"NeruBot Help: {response.content}")
    
    await ai_service.cleanup()


async def main():
    """Run all demos"""
    print("🤖 NeruBot Global AI Service Integration Demos")
    print("=" * 60)
    print("This shows how the AI service can be used across different features!\n")
    
    # Run demos
    await demo_music_ai_integration()
    await demo_news_ai_integration()
    await demo_confession_ai_integration()
    await demo_help_ai_integration()
    
    print("\n🎉 Demo complete!")
    print("\n💡 Integration Tips:")
    print("1. Use different system prompts for different features")
    print("2. Adjust temperature based on the use case")
    print("3. Set appropriate max_tokens for response length")
    print("4. The AI service handles provider selection automatically")
    print("5. Always initialize before use and cleanup when done")


if __name__ == "__main__":
    print("Note: This demo requires at least one AI API key to be configured.")
    print("Set OPENAI_API_KEY, ANTHROPIC_API_KEY, or GEMINI_API_KEY in your .env file.\n")
    
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\n⏹️  Demo interrupted by user")
    except Exception as e:
        print(f"\n💥 Demo failed: {e}")
        print("Make sure you have at least one AI API key configured!")
