# NeruBot Architecture

## Simplified Architecture Overview

```
nerubot/
├── bot.py                    # Main entry point - simplified
├── config/
│   ├── __init__.py
│   └── settings.py           # All configuration
├── cogs/                     # Feature modules (Discord cogs)
│   ├── __init__.py
│   ├── music.py             # Music functionality
│   ├── general.py           # General commands (help, ping, etc.)
│   ├── moderation.py        # Moderation features (kick, ban, etc.)
│   ├── fun.py               # Fun commands (jokes, memes, etc.)
│   └── utility.py           # Utility commands (weather, calculator, etc.)
├── services/                 # Business logic
│   ├── __init__.py
│   ├── music_service.py     # Music-related logic
│   ├── database_service.py  # Database operations
│   └── api_service.py       # External API calls
├── utils/                    # Shared utilities
│   ├── __init__.py
│   ├── logging_utils.py
│   ├── messages.py
│   ├── decorators.py        # Common decorators
│   └── helpers.py           # Helper functions
├── models/                   # Data models
│   ├── __init__.py
│   ├── song.py
│   ├── user.py
│   └── guild.py
└── requirements.txt
```

## Key Principles

1. **Separation of Concerns**: Each cog handles one feature area
2. **Service Layer**: Business logic separated from Discord interface
3. **Easy Extension**: Adding new features is just adding new cogs
4. **Clean Dependencies**: Clear import structure
5. **Configuration Centralized**: All settings in one place

## Adding New Features

To add a new feature (e.g., weather commands):

1. Create `cogs/weather.py` with Discord commands
2. Create `services/weather_service.py` with weather API logic
3. Add any models in `models/` if needed
4. The bot automatically loads all cogs

## Benefits

- **Modular**: Each feature is independent
- **Scalable**: Easy to add/remove features
- **Maintainable**: Clear code organization
- **Testable**: Services can be tested independently
