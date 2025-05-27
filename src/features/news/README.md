# News Feature

The News feature provides real-time news updates from trusted sources directly to your Discord server.

## Features

- **Real-time Updates**: Automatically fetches news every 10 minutes from trusted sources
- **Trusted Sources**: Uses RSS feeds from 12 international and Indonesian sources including BBC, Reuters, Bloomberg, and major Indonesian media
- **Manual Control**: Users can start/stop automatic updates manually
- **Latest News**: Get the latest news on demand
- **Source Management**: Add or remove news sources (admin only)

## Commands

### User Commands

- `/news latest [count]` - Get the latest news items (default: 5, max: 10)
- `/news sources` - List all configured news sources
- `/news status` - Show current news configuration and status
- `/news help` - Show help for news commands

### Admin Commands

- `/news set-channel [channel]` - Set channel for automatic updates and enable auto-posting
- `/news start` - Start automatic news updates
- `/news stop` - Stop automatic news updates
- `/news add <name> <feed_url>` - Add a news source
- `/news remove <name>` - Remove a news source

## How It Works

1. **Automatic Updates**: When the bot starts, it begins fetching news every 10 minutes
2. **Channel Setup**: Admins can set a specific channel for news updates using `/news set-channel`
3. **Auto-posting**: Once a channel is set, the bot will automatically post breaking news
4. **Manual Control**: Admins can start/stop the automatic posting at any time

## Default News Sources

### International Sources
- **BBC World**: World news from BBC
- **BBC Business**: Business news from BBC
- **Reuters Top News**: Top stories from Reuters
- **Reuters World**: World news from Reuters
- **Reuters Business**: Business news from Reuters
- **CNN Top Stories**: Top stories from CNN
- **NPR News**: News from National Public Radio
- **Al Jazeera**: International news from Al Jazeera
- **Bloomberg Markets**: Financial and market news from Bloomberg

### Indonesian Sources
- **Tempo**: Leading Indonesian news from Tempo.co
- **Antara News**: Official Indonesian news agency
- **Republika**: Indonesian news from Republika

## Setup Instructions

1. The news feature is automatically loaded when the bot starts
2. Use `/news set-channel` in your desired news channel to enable automatic updates
3. The bot will start posting breaking news automatically
4. Use `/news stop` to disable automatic updates if needed

## Technical Details

- **Update Interval**: 10 minutes (configurable in code)
- **RSS Parsing**: Uses feedparser library for reliable RSS feed parsing
- **Duplicate Prevention**: Tracks published articles to avoid duplicates
- **Error Handling**: Gracefully handles network errors and invalid feeds
- **Memory Management**: Keeps only the latest 100 news items in memory

## Customization

Administrators can:
- Add new RSS feeds using `/news add`
- Remove unwanted sources using `/news remove`
- Configure which channel receives updates
- Start/stop automatic posting as needed

The news service runs in the background and will continue fetching updates even when automatic posting is disabled, ensuring that `/news latest` always has fresh content available.
