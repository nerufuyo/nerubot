"""
Demo script for the Roast feature
"""
import asyncio
import random
import time
from src.features.roast.services.roast_service import roast_service
from src.features.roast.models.roast_models import ActivityType


async def demo_roast_feature():
    """Demonstrate the roast feature functionality"""
    print("🔥 NeruBot Roast Feature Demo 🔥")
    print("=" * 50)
    
    # Simulate user IDs
    user_ids = [
        123456789,  # Night owl user
        987654321,  # Spammer user  
        555444333,  # Music addict
        111222333,  # Lurker
        999888777   # Emoji spammer
    ]
    
    guild_id = 123456789
    channel_id = 987654321
    
    print("\n1. Simulating user activities...")
    
    # Simulate night owl behavior
    print("📊 Creating night owl user behavior...")
    for i in range(30):
        # Simulate late night messages (11 PM - 3 AM)
        night_hour = random.choice([23, 0, 1, 2, 3])
        timestamp = time.time() - (random.randint(1, 7) * 86400)  # Last week
        
        await roast_service.track_activity(
            user_id=user_ids[0],
            activity_type=ActivityType.MESSAGE,
            channel_id=channel_id,
            guild_id=guild_id,
            content=f"Message at {night_hour}:00 - Can't sleep again...",
            metadata={"simulated_hour": night_hour}
        )
    
    # Simulate spammer behavior
    print("📊 Creating spammer user behavior...")
    for i in range(80):
        await roast_service.track_activity(
            user_id=user_ids[1],
            activity_type=ActivityType.MESSAGE,
            channel_id=channel_id,
            guild_id=guild_id,
            content=f"Short message {i}",
            metadata={"burst_message": True}
        )
    
    # Simulate music addict behavior
    print("📊 Creating music addict user behavior...")
    for i in range(40):
        await roast_service.track_activity(
            user_id=user_ids[2],
            activity_type=ActivityType.MUSIC_REQUEST,
            channel_id=channel_id,
            guild_id=guild_id,
            content=f"Song request #{i}",
            metadata={"song_genre": random.choice(["pop", "rock", "edm", "classical"])}
        )
    
    # Simulate lurker behavior
    print("📊 Creating lurker user behavior...")
    for i in range(5):  # Very few messages
        await roast_service.track_activity(
            user_id=user_ids[3],
            activity_type=ActivityType.MESSAGE,
            channel_id=channel_id,
            guild_id=guild_id,
            content="...",
            metadata={"lurker_activity": True}
        )
    
    # Add some voice activity to show they're present but quiet
    for i in range(20):
        await roast_service.track_activity(
            user_id=user_ids[3],
            activity_type=ActivityType.VOICE_JOIN,
            channel_id=channel_id,
            guild_id=guild_id,
            content="General Voice Channel",
            metadata={"silent_join": True}
        )
    
    # Simulate emoji spammer behavior
    print("📊 Creating emoji spammer user behavior...")
    emojis = ["😂", "😭", "💀", "✨", "🔥", "🎵", "🎉", "💙", "👀", "🤡", "🌟", "⭐", "🎭", "🎪", "🎨"]
    for i in range(60):
        emoji = random.choice(emojis)
        await roast_service.track_activity(
            user_id=user_ids[4],
            activity_type=ActivityType.EMOJI_REACTION,
            channel_id=channel_id,
            guild_id=guild_id,
            content=emoji,
            metadata={"emoji_type": "unicode"}
        )
        
        # Also add some messages with lots of emojis
        if i % 5 == 0:
            emoji_message = " ".join(random.choices(emojis, k=random.randint(3, 8)))
            await roast_service.track_activity(
                user_id=user_ids[4],
                activity_type=ActivityType.MESSAGE,
                channel_id=channel_id,
                guild_id=guild_id,
                content=f"Hey everyone! {emoji_message}",
                metadata={"emoji_heavy": True}
            )
    
    print("\n2. Analyzing user behavior patterns...")
    
    for i, user_id in enumerate(user_ids):
        print(f"\n📊 Analyzing User {i+1} (ID: {user_id})...")
        pattern = await roast_service.analyze_user_behavior(user_id)
        
        print(f"   - Total activities: {sum(pattern.activity_frequency.values())}")
        print(f"   - Most active hours: {pattern.most_active_hours[:3]}")
        print(f"   - Late night percentage: {pattern.late_night_percentage:.1f}%")
        print(f"   - Activity types: {list(pattern.activity_frequency.keys())[:3]}")
        print(f"   - Emoji usage: {len(pattern.emoji_usage)} different emojis")
        print(f"   - Message length avg: {pattern.message_length_avg:.1f} chars")
    
    print("\n3. Generating roasts for each user type...")
    
    user_types = ["Night Owl", "Spammer", "Music Addict", "Lurker", "Emoji Spammer"]
    
    for i, (user_id, user_type) in enumerate(zip(user_ids, user_types)):
        print(f"\n🔥 Roasting {user_type} (User {i+1})...")
        
        roast_text, roast_category = await roast_service.generate_roast(user_id)
        
        print(f"   Category: {roast_category}")
        print(f"   Roast: {roast_text}")
        
        # Show user stats
        stats = await roast_service.get_user_roast_stats(user_id)
        print(f"   Stats: {stats['roast_count']} roasts, {stats['activity_count']} activities tracked")
    
    print("\n4. Testing custom AI roasts...")
    
    custom_roasts = [
        "their terrible music taste",
        "being online way too much",
        "never talking in voice chat",
        "using too many emojis",
        "their weird sleep schedule"
    ]
    
    for i, (user_id, custom_prompt) in enumerate(zip(user_ids, custom_roasts)):
        print(f"\n🤖 AI Custom Roast for User {i+1}...")
        print(f"   Prompt: {custom_prompt}")
        
        try:
            roast_text, roast_category = await roast_service.generate_roast(user_id, custom_prompt)
            print(f"   Category: {roast_category}")
            print(f"   AI Roast: {roast_text}")
        except Exception as e:
            print(f"   Error: {e}")
    
    print("\n5. Global roast statistics...")
    
    global_stats = await roast_service.get_global_roast_stats()
    print(f"   Total roasts delivered: {global_stats.total_roasts_delivered}")
    print(f"   Today's roasts: {global_stats.daily_roast_count}")
    print(f"   Roasts by category: {global_stats.roasts_by_category}")
    
    print("\n6. Testing cooldown system...")
    
    # Try to roast the same user again immediately
    print("\n🛡️ Testing roast cooldown...")
    roast_text, roast_category = await roast_service.generate_roast(user_ids[0])
    print(f"   Cooldown roast: {roast_text}")
    
    print("\n✅ Roast feature demo completed!")
    print("\nThe roast feature is ready to deliver epic burns! 🔥")
    print("\nKey features demonstrated:")
    print("- ✅ Activity tracking and pattern analysis")
    print("- ✅ Behavior-based roast category selection")
    print("- ✅ Template-based roast generation")
    print("- ✅ AI-powered custom roasts")
    print("- ✅ User statistics and insights")
    print("- ✅ Cooldown protection system")
    print("- ✅ Comprehensive behavior analysis")


if __name__ == "__main__":
    print("🎭 Starting Roast Feature Demo...")
    asyncio.run(demo_roast_feature())
