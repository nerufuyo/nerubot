"""
Demo script for the Anonymous Confession feature
This script demonstrates the confession system functionality
"""
import asyncio
import discord
from discord.ext import commands
from src.features.confession.services.confession_service import ConfessionService
from src.features.confession.models.confession import ConfessionStatus


async def demo_confession_system():
    """Demonstrate the confession system functionality."""
    print("🔥 Anonymous Confession System Demo")
    print("=" * 50)
    
    # Create service instance
    service = ConfessionService()
    
    # Demo guild and user IDs
    guild_id = 123456789
    user1_id = 111111111
    user2_id = 222222222
    channel_id = 333333333
    
    print("\n1. Setting up guild confession settings...")
    settings = service.update_guild_settings(
        guild_id,
        confession_channel_id=channel_id,
        moderation_enabled=False,
        anonymous_replies=True,
        max_confession_length=2000,
        cooldown_minutes=1  # Shorter cooldown for demo
    )
    print(f"✅ Guild settings configured:")
    print(f"   - Channel ID: {settings.confession_channel_id}")
    print(f"   - Moderation: {settings.moderation_enabled}")
    print(f"   - Anonymous replies: {settings.anonymous_replies}")
    print(f"   - Cooldown: {settings.cooldown_minutes} minutes")
    
    print("\n2. Creating first confession...")
    success, message, confession1 = service.create_confession(
        content="I've been struggling with anxiety lately and don't know who to talk to. Sometimes I feel like I'm drowning in my own thoughts.",
        author_id=user1_id,
        guild_id=guild_id
    )
    
    if success:
        print(f"✅ {message}")
        print(f"   - Confession ID: {confession1.confession_id}")
        print(f"   - Status: {confession1.status.value}")
        print(f"   - Content preview: {confession1.content[:50]}...")
    else:
        print(f"❌ {message}")
    
    print("\n3. Creating second confession...")
    success, message, confession2 = service.create_confession(
        content="I have a huge crush on someone in my class but I'm too shy to say anything. What should I do?",
        author_id=user2_id,
        guild_id=guild_id
    )
    
    if success:
        print(f"✅ {message}")
        print(f"   - Confession ID: {confession2.confession_id}")
        print(f"   - Content preview: {confession2.content[:50]}...")
    
    print("\n4. Testing cooldown system...")
    success, message, confession3 = service.create_confession(
        content="This should fail due to cooldown",
        author_id=user1_id,
        guild_id=guild_id
    )
    print(f"Expected failure: {message}")
    
    print("\n5. Creating replies to first confession...")
    
    # Reply 1
    success, message, reply1 = service.create_reply(
        confession_id=confession1.confession_id,
        content="I understand how you feel. Have you considered talking to a counselor or therapist? Sometimes professional help can make a huge difference.",
        author_id=user2_id,
        guild_id=guild_id
    )
    
    if success:
        print(f"✅ Reply 1: {message}")
        print(f"   - Reply ID: {reply1.reply_id}")
        print(f"   - Content preview: {reply1.content[:50]}...")
    
    # Reply 2  
    success, message, reply2 = service.create_reply(
        confession_id=confession1.confession_id,
        content="You're not alone in feeling this way. Anxiety is more common than you think. Take it one day at a time. 💙",
        author_id=user1_id,  # Same user can reply to their own confession anonymously
        guild_id=guild_id
    )
    
    if success:
        print(f"✅ Reply 2: {message}")
        print(f"   - Reply ID: {reply2.reply_id}")
    
    print("\n6. Testing confession lookup...")
    
    # Test exact ID lookup
    found_confession = service.get_confession(confession1.confession_id)
    if found_confession:
        print(f"✅ Found confession by exact ID: {found_confession.confession_id}")
    
    # Test partial ID lookup
    partial_id = confession1.confession_id[:4]
    found_confession = service.get_confession_by_tag(partial_id, guild_id)
    if found_confession:
        print(f"✅ Found confession by partial ID '{partial_id}': {found_confession.confession_id}")
    
    print("\n7. Getting confession statistics...")
    guild_confessions = service.get_guild_confessions(guild_id)
    total_confessions = len(guild_confessions)
    total_replies = sum(conf.reply_count for conf in guild_confessions)
    
    print(f"📊 Guild Statistics:")
    print(f"   - Total confessions: {total_confessions}")
    print(f"   - Total replies: {total_replies}")
    print(f"   - Average replies per confession: {total_replies/total_confessions:.1f}")
    
    print("\n8. Listing all confessions and replies...")
    for confession in guild_confessions:
        print(f"\n📝 Confession #{confession.confession_id}")
        print(f"   Content: {confession.content[:100]}...")
        print(f"   Replies: {confession.reply_count}")
        print(f"   Created: {confession.created_at.strftime('%Y-%m-%d %H:%M:%S')}")
        
        replies = service.get_confession_replies(confession.confession_id)
        for i, reply in enumerate(replies, 1):
            print(f"   💬 Reply {i}: {reply.content[:60]}...")
    
    print("\n" + "=" * 50)
    print("🎉 Demo completed successfully!")
    print("\nKey Features Demonstrated:")
    print("✅ Anonymous confession submission")
    print("✅ Anonymous replies to confessions") 
    print("✅ Confession ID system for easy reference")
    print("✅ User cooldown system")
    print("✅ Guild-specific settings")
    print("✅ Data persistence")
    print("✅ Statistics and moderation tools")


if __name__ == "__main__":
    asyncio.run(demo_confession_system())