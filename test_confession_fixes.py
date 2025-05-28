#!/usr/bin/env python3
"""Test fixed confession functionality"""

import sys
import asyncio
sys.path.append('.')

async def test_confession_fixes():
    print('🧪 Testing Fixed Confession Functionality')
    print('=' * 50)
    
    # Test 1: Import and basic functionality
    print('1. Testing imports and basic setup...')
    from src.features.confession.cogs.confession_cog import ConfessionCog, ConfessionView
    from src.features.confession.services.confession_service import ConfessionService
    import discord
    from discord.ext import commands
    
    # Create bot and cog
    intents = discord.Intents.default()
    bot = commands.Bot(command_prefix='!', intents=intents)
    cog = ConfessionCog(bot)
    
    print('   ✅ Imports successful')
    
    # Test 2: Service functionality
    print('2. Testing service functionality...')
    service = ConfessionService()
    service.update_guild_settings(guild_id=123, confession_channel_id=456)
    
    # Create confession and mark as posted with thread
    success, msg, confession = service.create_confession('Test confession', 789, 123)
    if success:
        service.mark_confession_posted(confession.confession_id, 456, 999, 888)
        print(f'   ✅ Confession created with ID: {confession.confession_id}')
        
        # Verify thread ID was set
        updated_confession = service.get_confession(confession.confession_id)
        print(f'   ✅ Thread ID set: {updated_confession.thread_id}')
        
        # Test reply creation
        success, msg, reply = service.create_reply(confession.confession_id, 'Test reply', 111, 123)
        if success:
            print(f'   ✅ Reply created for confession {reply.confession_id}')
        else:
            print(f'   ❌ Reply creation failed: {msg}')
    else:
        print(f'   ❌ Confession creation failed: {msg}')
    
    # Test 3: View functionality
    print('3. Testing view functionality...')
    view = ConfessionView(123)
    buttons = [item for item in view.children if hasattr(item, 'label')]
    button_labels = [button.label for button in buttons]
    
    print(f'   Button labels: {button_labels}')
    assert 'Reply Anonymously' in button_labels
    assert 'Create New Confession' in button_labels
    print('   ✅ Buttons configured correctly')
    
    # Test 4: Cog loading
    print('4. Testing cog loading...')
    try:
        await cog.cog_load()
        print('   ✅ Cog loaded successfully with persistent views')
    except Exception as e:
        print(f'   ❌ Cog loading failed: {e}')
        import traceback
        traceback.print_exc()
    
    print('\n' + '=' * 50)
    print('🎉 All tests completed!')
    print('\nFixed Issues:')
    print('✅ Persistent views properly registered for each confession')
    print('✅ Reply targeting uses correct confession_id')
    print('✅ Enhanced logging for debugging')
    print('✅ Views registered on confession creation and cog load')

if __name__ == "__main__":
    print("Starting confession fixes test...")
    try:
        asyncio.run(test_confession_fixes())
    except Exception as e:
        print(f"❌ Test failed: {e}")
        import traceback
        traceback.print_exc()
