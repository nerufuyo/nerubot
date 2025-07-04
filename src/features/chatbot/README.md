# 🤖 NeruBot Chatbot Feature

An intelligent AI-powered chatbot that brings NeruBot's personality to life using multiple AI providers (OpenAI, Claude, and Gemini).

## ✨ Features

### 🎭 Unique Personality
- **Playful & Witty**: NeruBot has a fun, slightly sarcastic but always friendly personality
- **Gaming/Anime References**: Uses gaming and anime culture references naturally
- **Discord Culture**: Understands Discord slang and internet culture
- **Tech Humor**: Makes jokes about being a bot and tech-related topics

### 🧠 Multi-AI Support
- **OpenAI GPT**: Fast and reliable responses
- **Anthropic Claude**: Thoughtful and nuanced conversations  
- **Google Gemini**: Creative and versatile interactions
- **Provider Rotation**: Automatically rotates between available providers for variety
- **User Preferences**: Users can set their preferred AI provider

### 💬 Smart Interaction
- **Mention Response**: Responds when mentioned in servers
- **DM Support**: Full conversation support in direct messages
- **Natural Conversations**: Maintains context and personality throughout chats
- **Welcome Messages**: Greets users when they start new conversations
- **Auto Thanks**: Sends appreciation messages after 5 minutes of inactivity

### 📊 Session Management
- **Active Sessions**: Tracks ongoing conversations
- **Timeout Handling**: Automatically manages inactive sessions
- **User Statistics**: Tracks chat history and preferences
- **Background Cleanup**: Manages resources efficiently

## 🎯 Usage

### Basic Chat
- **Mention the bot**: `@NeruBot Hey, how are you?`
- **Direct Messages**: Just send a message directly to NeruBot
- **Slash Command**: `/chat message:Hello there!`

### Commands

#### `/chat`
Start a conversation with NeruBot
```
/chat message:What's your favorite music genre?
```

#### `/ai-provider`
Set your preferred AI provider
```
/ai-provider provider:OpenAI GPT
```

#### `/chat-stats`
View your chat statistics
```
/chat-stats
```

#### `/ai-status`
Check AI services status
```
/ai-status
```

## ⚙️ Configuration

### Environment Variables
```bash
# Required for AI functionality
OPENAI_API_KEY=your_openai_api_key_here
ANTHROPIC_API_KEY=your_claude_api_key_here
GEMINI_API_KEY=your_gemini_api_key_here

# At least one AI provider is required
```

### Settings
The chatbot behavior can be configured in `src/features/chatbot/config.py`:

- **Session timeout**: How long before sending thanks message (default: 5 minutes)
- **Response creativity**: Temperature setting for AI responses (default: 0.8)
- **Message limits**: Maximum response length and rate limits
- **Feature toggles**: Enable/disable specific features

## 🔧 Technical Details

### Architecture
```
src/features/chatbot/
├── config.py              # Configuration and constants
├── cogs/
│   └── chatbot_cog.py     # Discord cog (main interface)
├── models/
│   └── chat_models.py     # Data models for sessions/stats
└── services/
    └── chatbot_service.py # Core chatbot logic

src/services/
└── ai_service.py          # Global AI service (used by other features)
```

### Key Components

#### AI Service (`src/services/ai_service.py`)
- **Global Service**: Used across all bot features needing AI
- **Multi-Provider**: Supports OpenAI, Claude, and Gemini
- **Error Handling**: Graceful fallbacks and error messages
- **Personality Injection**: Adds NeruBot's personality to all interactions

#### Chatbot Service (`src/features/chatbot/services/chatbot_service.py`)
- **Session Management**: Tracks active conversations
- **Welcome/Thanks Logic**: Handles greeting and farewell messages
- **Provider Selection**: Chooses AI provider based on user preference
- **Background Tasks**: Monitors for inactive sessions

#### Discord Cog (`src/features/chatbot/cogs/chatbot_cog.py`)
- **Message Listener**: Responds to mentions and DMs
- **Slash Commands**: Provides interactive commands
- **Error Handling**: User-friendly error messages
- **Typing Indicators**: Shows bot is "thinking"

## 🎨 Personality System

NeruBot's personality is defined by:

### Core Traits
- Playful and witty
- Slightly sarcastic but friendly
- Casual and laid-back
- Mischievous but never mean

### Interests
- Music (all genres)
- Gaming and anime culture
- Discord community vibes
- Technology and memes

### Speaking Style
- Conversational and natural
- Strategic emoji usage
- Gaming/anime references
- Discord slang when appropriate
- Self-aware bot humor

## 🔄 Session Flow

1. **User Interaction**: User mentions bot or sends DM
2. **Session Check**: Check for existing active session
3. **Welcome**: Send welcome message for new sessions
4. **AI Processing**: Generate response using selected AI provider
5. **Response**: Send AI response with NeruBot personality
6. **Monitoring**: Track session activity
7. **Timeout**: Send thanks message after 5 minutes of inactivity
8. **Cleanup**: Remove old inactive sessions

## 🛡️ Error Handling

The chatbot includes comprehensive error handling:

- **AI Provider Failures**: Falls back to other providers or shows friendly error
- **Rate Limiting**: Respects API limits and shows appropriate messages
- **Network Issues**: Handles timeouts and connection problems
- **Invalid Inputs**: Validates and sanitizes user input
- **Resource Cleanup**: Properly manages sessions and connections

## 📈 Statistics & Analytics

Users can view their chat statistics including:
- Total messages sent
- AI responses received
- Preferred AI provider
- First chat date
- Current session info

## 🚀 Future Enhancements

Potential improvements:
- **Conversation Memory**: Remember previous conversations
- **Custom Personalities**: Allow users to customize bot personality
- **Voice Responses**: Text-to-speech integration
- **Image Generation**: AI-powered image creation
- **Context Awareness**: Better understanding of server context
- **Multi-language**: Support for different languages

## 🐛 Troubleshooting

### Common Issues

**Bot not responding to mentions:**
- Check bot permissions (Read Messages, Send Messages)
- Verify Message Content Intent is enabled
- Ensure bot is in the channel

**AI providers not working:**
- Verify API keys are set correctly
- Check API quota and billing
- Use `/ai-status` to check provider status

**Sessions not timing out:**
- Check if background task is running
- Verify bot has necessary permissions
- Check logs for task errors

### Debug Commands
- `/ai-status` - Check AI provider availability
- `/chat-stats` - View user statistics
- Check bot logs for detailed error information
