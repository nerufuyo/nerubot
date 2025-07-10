"""
Chatbot Configuration and Constants
"""

# ============================
# CHATBOT SETTINGS
# ============================

CHATBOT_CONFIG = {
    "session_timeout_minutes": 5,
    "max_response_tokens": 300,
    "response_temperature": 0.8,
    "enable_welcome_messages": True,
    "enable_thanks_messages": True,
    "enable_personality": True,
    "max_message_length": 2000,
}

# ============================
# AI PROVIDER SETTINGS
# ============================

AI_PROVIDER_CONFIG = {
    "default_provider": "openai",
    "fallback_provider": "claude",
    "retry_attempts": 2,
    "request_timeout": 30.0,
    "rate_limit_delay": 1.0,
}

# ============================
# CHATBOT PERSONALITY
# ============================

NERUBOT_PERSONALITY = {
    "base_traits": [
        "playful", "witty", "slightly sarcastic", "friendly",
        "casual", "laid-back", "mischievous but kind"
    ],
    "interests": [
        "music", "gaming", "anime", "memes", "discord culture",
        "technology", "internet culture"
    ],
    "speaking_style": [
        "conversational", "uses emojis sparingly", "gaming references",
        "anime references", "discord slang", "tech humor"
    ]
}

# ============================
# RESPONSE TEMPLATES
# ============================

WELCOME_MESSAGES = [
    "Hey there! 👋 What's on your mind today?",
    "Oh, someone wants to chat! What's up? 😊",
    "Heyyy~ What can I help you with? ✨",
    "Well well, look who's back! What's cooking? 🍳",
    "Yo! Ready for some quality bot conversation? 🤖",
    "Greetings, human! What adventure shall we embark on? 🚀",
    "Oh hi there! I was just organizing my digital thoughts~ 💭",
    "Another chat session? I'm all ears! Well... all code~ 👂",
    "Sup! My neural networks are all warmed up~ 🔥",
    "Hey hey! What brings you to my digital corner? 🌟"
]

THANKS_MESSAGES = [
    "Thanks for chatting with me! Come back anytime~ 🌟",
    "It was fun talking! Don't be a stranger! 👋",
    "Hope I helped! Catch you later~ ✨",
    "Always a pleasure! See you around! 🎵",
    "Thanks for the chat! Until next time! 🤖",
    "That was nice! I'll be here when you need me~ 💙",
    "Enjoyed our conversation! Take care! 🌈",
    "Thanks for hanging out! Stay awesome! ⭐",
    "Great chat! My circuits are happy~ ⚡",
    "See ya! Keep being amazing! 🚀"
]

ERROR_MESSAGES = [
    "Oops! My circuits got a bit tangled there. Try again? 🤖⚡",
    "Uh oh! Something went wrong in my neural networks... 💥",
    "My AI brain had a little hiccup. Give me another shot? 🧠",
    "Error 404: Witty response not found. Retry? 😅",
    "Beep boop! System glitch detected. Try once more? 🔧",
    "Houston, we have a problem... but nothing a retry can't fix! 🚀",
    "My code got its wires crossed. One more time? ⚡",
    "Oops! My digital neurons misfired. Let's try again~ 🤖"
]

# ============================
# FEATURE FLAGS
# ============================

CHATBOT_FEATURES = {
    "respond_to_mentions": True,
    "respond_to_dms": True,
    "session_tracking": True,
    "welcome_messages": True,
    "thanks_messages": True,
    "user_stats": True,
    "provider_rotation": True,
    "personality_mode": True,
    "auto_cleanup": True,
}

# ============================
# RATE LIMITING
# ============================

RATE_LIMITS = {
    "messages_per_minute": 10,
    "messages_per_hour": 100,
    "cooldown_seconds": 2,
    "burst_limit": 3,
}

# ============================
# LOGGING SETTINGS
# ============================

CHATBOT_LOGGING = {
    "log_conversations": True,
    "log_ai_responses": True,
    "log_errors": True,
    "log_stats": True,
    "max_log_entries": 10000,
}
