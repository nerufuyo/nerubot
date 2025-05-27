# NeruBot Help System

This module provides a comprehensive help system for NeruBot with an interactive, paginated interface and multiple command organization options.

## Features

- **Paginated Help Menu**: Browse through commands with interactive buttons
- **Category Organization**: Commands organized by feature and function
- **Multiple Command Views**:
  - `/help` - Full interactive help with category navigation
  - `/commands` - Compact command reference card
  - `/about` - Bot information and statistics
  - `/features` - Detailed feature showcase

## Components

1. **HelpCog (`help_cog.py`)**:
   - Main help command with pagination
   - Organizes commands by category
   - Interactive UI with navigation buttons

2. **CommandsCog (`commands_cog.py`)**:
   - Compact command reference card
   - Quick overview of all available commands
   - Helpful tips for users

3. **AboutCog (`about_cog.py`)**:
   - Bot information and statistics
   - System information
   - Links and credits

4. **FeaturesCog (`features_cog.py`)**:
   - Showcase of current features
   - Preview of upcoming features
   - Music source information

## Usage

Users can access the help system through the following commands:

- `/help` - Access the main paginated help system
- `/commands` - View a compact command reference card
- `/about` - View bot information and statistics
- `/features` - View available and upcoming features

## Implementation

The help system uses Discord's UI components (buttons) for navigation and interaction. The help content is dynamically generated based on the available commands in the bot.

## Customization

To add new commands to the help system, update the `command_categories` dictionary in the `HelpCog` class.
