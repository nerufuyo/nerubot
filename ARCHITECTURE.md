# NeruBot Architecture

## Modular Features Architecture Overview

```
nerubot/
├── src/
│   ├── main.py              # Main entry point
│   ├── core/                # Core utilities and shared components
│   │   └── utils/           # Logging, file utils, messages
│   ├── interfaces/          # External interfaces
│   │   └── discord/         # Discord bot interface
│   │       ├── bot.py       # Main bot class
│   │       └── help_cog.py  # Help system
│   └── features/            # Modular feature system
│       ├── music/           # 🎵 Music feature (COMPLETE)
│       │   ├── cogs/        # Discord commands
│       │   ├── services/    # Music business logic
│       │   └── models/      # Song data models
│       ├── news/            # 📰 News feature (COMPLETE)
│       │   ├── cogs/        # News commands
│       │   ├── services/    # RSS fetching logic
│       │   └── models/      # News data models
│       ├── quotes/          # 🔮 AI Quotes feature (READY)
│       │   ├── cogs/        # Quote commands
│       │   ├── services/    # DeepSeek AI integration
│       │   └── models/      # Quote data models
│       ├── profile/         # 👤 User Profiles (READY)
│       │   ├── cogs/        # Profile commands
│       │   ├── services/    # User data management
│       │   └── models/      # Profile data models
│       └── confession/      # 🤫 Anonymous Confessions (READY)
│           ├── cogs/        # Confession commands
│           ├── services/    # Anonymization logic
│           └── models/      # Confession data models
├── requirements.txt
├── run_nerubot.sh          # Deployment script
└── README.md               # Documentation
```

## Key Principles

1. **Feature Isolation**: Each feature is completely self-contained
2. **Consistent Structure**: All features follow the same cogs/services/models pattern
3. **Easy Maintenance**: Update one feature without affecting others
4. **Scalable Design**: Adding new features takes minutes, not hours
5. **DRY Principle**: Shared utilities prevent code duplication
6. **KISS Principle**: Simple, clean interfaces throughout

## Adding New Features

To add a new feature (e.g., weather commands):

1. Create `src/features/weather/` directory structure:
   ```
   weather/
   ├── __init__.py
   ├── cogs/
   │   ├── __init__.py
   │   └── weather_cog.py     # Discord commands
   ├── services/
   │   ├── __init__.py
   │   └── weather_service.py # Weather API logic
   └── models/
       ├── __init__.py
       └── weather.py         # Weather data models
   ```

2. The bot automatically loads all feature cogs!
3. Features can be easily enabled/disabled independently
4. The bot automatically loads all cogs

## Benefits

- **Modular**: Each feature is completely independent
- **Maintainable**: Clean separation of concerns
- **Scalable**: Easy to add/remove features without affecting others
- **Testable**: Services can be unit tested independently
- **Consistent**: All features follow the same architectural pattern
- **Future-proof**: Architecture supports any type of Discord bot feature

## Current Feature Status

| Feature | Status | Description |
|---------|--------|-------------|
| **Music** | ✅ **ACTIVE** | YouTube music with queue management |
| **News** | ✅ **ACTIVE** | RSS news feeds with 6 sources |
| **Quotes** | 🚧 **READY** | AI-powered quotes via DeepSeek |
| **Profile** | 🚧 **READY** | User profiles and statistics |
| **Confession** | 🚧 **READY** | Anonymous confession system |
