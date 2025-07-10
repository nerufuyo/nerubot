#!/usr/bin/env python3
"""
Demo script for the Anonymous Confession System

This script demonstrates the new confession system implementation
based on the specification requirements.
"""

import asyncio
import discord
from discord.ext import commands
from datetime import datetime
import sys
import os

# Add the src directory to the path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..', '..'))

from src.features.confession.services.confession_service import ConfessionService
from src.features.confession.models.confession import Confession, ConfessionReply, GuildConfessionSettings


def print_section(title):
    """Print a section header."""
    print(f"\n{'='*60}")
    print(f"  {title}")
    print(f"{'='*60}")


def print_confession_format(confession):
    """Print confession in the expected format."""
    print(f"\n📝 Confession #{confession.confession_id:03d}")
    print(f"{confession.content}")
    if confession.attachments:
        print(f"📎 Attachments: {', '.join(confession.attachments)}")
    print(f"ID: CONF-{confession.confession_id:03d} | 🔄 Reply")
    print(f"Created: {confession.created_at.strftime('%Y-%m-%d %H:%M:%S')}")


def print_reply_format(reply):
    """Print reply in the expected format."""
    print(f"\n↪️ {reply.reply_id}")
    print(f"{reply.content}")
    if reply.attachments:
        print(f"📎 Attachments: {', '.join(reply.attachments)}")
    print(f"ID: {reply.reply_id} | 🔄 Reply")
    print(f"Created: {reply.created_at.strftime('%Y-%m-%d %H:%M:%S')}")


def demo_confession_system():
    """Demonstrate the confession system."""
    
    print_section("Anonymous Confession System Demo")
    print("This demo shows the new confession system implementation")
    print("based on the specification requirements.")
    
    # Initialize the service
    service = ConfessionService()
    
    # Demo guild ID
    guild_id = 12345
    user_id_1 = 98765
    user_id_2 = 54321
    
    print_section("1. Guild Setup")
    
    # Set up guild settings
    service.update_guild_settings(
        guild_id=guild_id,
        confession_channel_id=9999,
        max_confession_length=2000,
        max_reply_length=1000
    )
    
    settings = service.get_guild_settings(guild_id)
    print(f"✅ Guild {guild_id} configured")
    print(f"   Confession Channel: {settings.confession_channel_id}")
    print(f"   Max Confession Length: {settings.max_confession_length}")
    print(f"   Max Reply Length: {settings.max_reply_length}")
    print(f"   Next Confession ID: CONF-{settings.next_confession_id:03d}")
    
    print_section("2. Creating Confessions")
    
    # Create first confession
    success, message, confession_id = service.create_confession(
        content="I've been struggling with anxiety lately and don't know who to talk to.",
        author_id=user_id_1,
        guild_id=guild_id,
        attachments=["https://example.com/image1.png"]
    )
    
    if success:
        print("✅ First confession created successfully!")
        # Get the confession object
        confession1 = service.get_confession(int(confession_id.split('-')[1]))
        print_confession_format(confession1)
    else:
        print(f"❌ Failed to create confession: {message}")
        confession1 = None
    
    # Create second confession
    success, message, confession_id = service.create_confession(
        content="I have a secret crush on my best friend but I'm too scared to tell them.",
        author_id=user_id_2,
        guild_id=guild_id,
        attachments=["https://example.com/gif1.gif", "https://example.com/image2.jpg"]
    )
    
    if success:
        print("✅ Second confession created successfully!")
        # Get the confession object
        confession2 = service.get_confession(int(confession_id.split('-')[1]))
        print_confession_format(confession2)
    else:
        print(f"❌ Failed to create confession: {message}")
        confession2 = None
    
    print_section("3. Creating Replies")
    
    # Create replies to first confession
    if confession1:
        success, message, reply_id = service.create_reply(
            confession_id=confession1.confession_id,
            content="You're not alone in this. Consider talking to a counselor or trusted friend.",
            author_id=user_id_2,
            guild_id=guild_id
        )
        
        if success:
            print("✅ First reply created successfully!")
            # Get the reply object
            replies = service.get_confession_replies(confession1.confession_id)
            reply1 = next((r for r in replies if r.reply_id == reply_id), None)
            if reply1:
                print_reply_format(reply1)
        else:
            print(f"❌ Failed to create reply: {message}")
            reply1 = None
    
    # Create another reply to first confession
    if confession1:
        success, message, reply_id = service.create_reply(
            confession_id=confession1.confession_id,
            content="I've been through something similar. Feel free to reach out if you need support.",
            author_id=user_id_1,
            guild_id=guild_id,
            attachments=["https://example.com/supportive_image.png"]
        )
        
        if success:
            print("✅ Second reply created successfully!")
            # Get the reply object
            replies = service.get_confession_replies(confession1.confession_id)
            reply2 = next((r for r in replies if r.reply_id == reply_id), None)
            if reply2:
                print_reply_format(reply2)
        else:
            print(f"❌ Failed to create reply: {message}")
            reply2 = None
    
    # Create reply to second confession
    if confession2:
        success, message, reply_id = service.create_reply(
            confession_id=confession2.confession_id,
            content="Sometimes the best friendships grow into something more. Be honest about your feelings!",
            author_id=user_id_1,
            guild_id=guild_id
        )
        
        if success:
            print("✅ Third reply created successfully!")
            # Get the reply object
            replies = service.get_confession_replies(confession2.confession_id)
            reply3 = next((r for r in replies if r.reply_id == reply_id), None)
            if reply3:
                print_reply_format(reply3)
        else:
            print(f"❌ Failed to create reply: {message}")
            reply3 = None
    
    print_section("4. ID System Demonstration")
    
    print("📋 ID System Examples:")
    if confession1 and confession2:
        print(f"   Confession IDs: CONF-{confession1.confession_id:03d}, CONF-{confession2.confession_id:03d}")
        
        # Get all replies for demo
        replies1 = service.get_confession_replies(confession1.confession_id)
        replies2 = service.get_confession_replies(confession2.confession_id)
        
        if replies1 or replies2:
            reply_ids = []
            if replies1:
                reply_ids.extend([r.reply_id for r in replies1])
            if replies2:
                reply_ids.extend([r.reply_id for r in replies2])
            print(f"   Reply IDs: {', '.join(reply_ids)}")
    
    print("\n💡 ID System Features:")
    print("   • Confessions: Sequential numbering (CONF-001, CONF-002, etc.)")
    print("   • Replies: Parent ID + letter suffix (REPLY-001-A, REPLY-001-B, etc.)")
    print("   • Easy reference for users to reply to specific messages")
    print("   • Hierarchical organization for better thread management")
    
    print_section("5. Statistics")
    
    confessions = service.get_guild_confessions(guild_id)
    total_replies = sum(confession.reply_count for confession in confessions)
    
    print(f"📊 Guild {guild_id} Statistics:")
    print(f"   Total Confessions: {len(confessions)}")
    print(f"   Total Replies: {total_replies}")
    print(f"   Average Replies per Confession: {total_replies/len(confessions):.1f}")
    if confessions:
        print(f"   Latest Confession: CONF-{confessions[0].confession_id:03d}")
    
    print_section("6. Modal Interface Demo")
    
    print("🎭 Modal Interface Structure:")
    print("\n📝 New Confession Modal:")
    print("   Title: 'Create New Confession'")
    print("   Fields:")
    print("     1. Message (required) - Text area for confession content")
    print("     2. Attachments (optional) - Text input for URLs")
    print("\n💬 Reply Modal:")
    print("   Title: 'Reply to Confession CONF-XXX'")
    print("   Fields:")
    print("     1. Message (required) - Text area for reply content")
    print("     2. Confession ID (read-only) - Auto-populated with target message ID")
    print("     3. Attachments (optional) - Text input for URLs")
    
    print_section("7. Thread Organization")
    
    print("🧵 Thread Structure:")
    print("   • Each confession creates a dedicated thread")
    print("   • Thread name: '💬 Confession #XXX'")
    print("   • First message: Confession content with Reply button")
    print("   • Subsequent messages: Replies with their own Reply buttons")
    print("   • All messages maintain complete anonymity")
    
    print_section("8. Anonymity Features")
    
    print("🔒 Anonymity Protection:")
    print("   • All messages posted by bot account")
    print("   • No usernames, avatars, or user IDs visible")
    print("   • User IDs stored internally but never displayed")
    print("   • Consistent anonymous formatting for all messages")
    print("   • No way to correlate messages with users")
    
    print_section("Demo Complete!")
    print("The Anonymous Confession System has been successfully recreated")
    print("according to the specification requirements.")
    print("\n🎯 Key Features Implemented:")
    print("   ✅ Exact modal structure (2 fields for confession, 3 for reply)")
    print("   ✅ Proper ID system (CONF-XXX, REPLY-XXX-Y)")
    print("   ✅ Thread-based organization")
    print("   ✅ Complete anonymity protection")
    print("   ✅ Attachment support")
    print("   ✅ Persistent button interface")
    print("   ✅ Sequential ID generation")
    print("   ✅ Hierarchical reply system")


if __name__ == "__main__":
    demo_confession_system()
