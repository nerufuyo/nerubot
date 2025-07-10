# Anonymous Confession System

## Overview

The Anonymous Confession System allows users to share their thoughts, feelings, and experiences completely anonymously in your Discord server. The system is designed to maintain complete anonymity while providing an organized and interactive environment for confessions and replies.

## Features

### 🔒 Complete Anonymity
- All messages are posted by the bot, not the user
- No usernames, avatars, or user IDs are visible
- No way to trace messages back to original authors
- Consistent anonymous formatting throughout

### 🆔 Unique ID System
- **Confessions**: Sequential numbering (CONF-001, CONF-002, etc.)
- **Replies**: Parent ID + letter suffix (REPLY-001-A, REPLY-001-B, etc.)
- Easy reference system for users to reply to specific messages

### 📝 Modal-Based Interface
- **New Confession Modal**: 2 fields (message, attachments)
- **Reply Modal**: 3 fields (message, confession_id, attachments)
- Clean, intuitive user interface

### 🧵 Thread Organization
- Each confession automatically creates a thread
- All replies are posted in the thread for better organization
- Threads auto-archive after 7 days

### 📎 Attachment Support
- Support for images, GIFs, and videos
- Multiple attachments per confession/reply
- URL-based attachment system

## Setup

### 1. Channel Setup
Use the `/confession-setup` command to set up the confession system:

```
/confession-setup channel:#confessions
```

This will:
- Set the designated channel for confessions
- Send an introduction message explaining the system
- Add a persistent "Create New Confession" button

### 2. Introduction Message
The bot will post an introduction message in the channel explaining:
- How the system works
- Features available
- How to create confessions and replies
- Anonymity guarantees

## Usage

### Creating Confessions

Users can create confessions by:
1. Clicking the "Create New Confession" button
2. Filling out the modal with:
   - **Message**: Their confession content (required)
   - **Attachments**: URLs for images/videos (optional)

### Replying to Confessions

Users can reply to confessions by:
1. Clicking the "Reply" button on any confession
2. Filling out the modal with:
   - **Message**: Their reply content (required)
   - **Confession ID**: Auto-populated (read-only)
   - **Attachments**: URLs for images/videos (optional)

### ID System

- **Confessions**: `CONF-001`, `CONF-002`, etc.
- **Replies**: `REPLY-001-A`, `REPLY-001-B`, etc.
- IDs are displayed in message footers for easy reference

## Message Format

### Confession Format
```
📝 Confession #001
[User's confession content]
[Attachments if any]
ID: CONF-001 | 🔄 Reply
```

### Reply Format
```
↪️ Reply to #001
[User's reply content]
[Attachments if any]
ID: REPLY-001-A | 🔄 Reply
```

## Admin Commands

### `/confession-setup`
- **Description**: Set up confession channel
- **Permission**: Administrator
- **Usage**: `/confession-setup channel:#confessions`

### `/confession-settings`
- **Description**: View current confession settings
- **Permission**: Administrator
- **Shows**: Channel, moderation status, limits, cooldowns, next IDs

### `/confession-stats`
- **Description**: View confession statistics
- **Permission**: Public
- **Shows**: Total confessions, replies, averages, latest confession

## Configuration

### Guild Settings
- **Confession Channel**: Where confessions are posted
- **Moderation**: Enable/disable moderation (currently disabled)
- **Anonymous Replies**: Always enabled
- **Max Confession Length**: 2000 characters (default)
- **Max Reply Length**: 1000 characters (default)
- **Cooldown**: 5 minutes between confessions/replies (default)

### Persistent Views
- Buttons remain functional after bot restarts
- Views are automatically registered on bot startup
- No need to re-setup after restarts

## Technical Implementation

### Models
- `Confession`: Stores confession data with attachments list
- `ConfessionReply`: Stores reply data with string-based IDs
- `GuildConfessionSettings`: Per-guild configuration

### Services
- `ConfessionService`: Handles all confession logic
- Sequential ID generation for confessions
- Letter-based ID generation for replies
- Data persistence via JSON files

### Views
- `SetupView`: Initial setup button
- `ConfessionView`: Main interaction buttons
- Persistent views for reliability

## Data Structure

### Confession IDs
- Start from 1 for each guild
- Format: `CONF-{ID:03d}` (e.g., CONF-001)
- Sequential numbering per guild

### Reply IDs
- Format: `REPLY-{confession_id:03d}-{letter}` (e.g., REPLY-001-A)
- Letters increment per confession (A, B, C, etc.)
- Hierarchical organization

### File Storage
- `confessions.json`: All confession data
- `replies.json`: All reply data
- `settings.json`: Guild settings and ID counters

## Security & Privacy

### Anonymity Protection
- User IDs are stored but never displayed
- All messages posted by bot account
- No correlation between messages and users
- Consistent formatting for all messages

### Moderation
- Currently auto-approves all confessions
- Framework exists for future moderation features
- Cooldown system prevents spam

## Troubleshooting

### Common Issues
1. **Buttons not working**: Views are persistent but may need bot restart
2. **Threads not creating**: Check bot permissions for thread creation
3. **Attachments not showing**: Ensure URLs are valid and accessible

### Required Permissions
- Send Messages
- Create Threads
- Upload Files
- Use Slash Commands
- Add Reactions/Buttons

## Future Enhancements

### Planned Features
- Moderation queue system
- Custom cooldown per guild
- Confession categories
- Report system
- Statistics dashboard

### Extensibility
- Modular design allows easy feature additions
- Clean separation of concerns
- Comprehensive error handling
- Detailed logging system
