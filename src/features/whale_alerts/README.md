# Whale Alerts Feature

The Whale Alerts feature provides real-time monitoring of large cryptocurrency transactions (whale alerts) and tweets from crypto gurus/influencers directly to your Discord server.

## Features

### 🐋 Whale Alerts
- **Real-time Monitoring**: Track large cryptocurrency transactions across multiple blockchains
- **Smart Filtering**: Configurable minimum amounts and priority tokens
- **Multiple Alert Types**: 
  - Large transfers ($100K+)
  - Exchange deposits/withdrawals ($500K+)
  - DEX swaps ($50K+)  
  - Whale accumulation/distribution ($1M+)
- **Rich Information**: Transaction details, blockchain info, wallet labels
- **Visual Indicators**: Color-coded embeds and reaction emojis based on transaction size

### 🧙‍♂️ Crypto Guru Tweets
- **Influencer Monitoring**: Track tweets from major crypto personalities
- **Sentiment Analysis**: Automatic classification of tweets (bullish/bearish/neutral/urgent)
- **Token Detection**: Identify mentioned cryptocurrencies in tweets
- **Engagement Filtering**: Focus on high-engagement content
- **Real-time Updates**: Get notified when gurus tweet about crypto

## Monitored Crypto Gurus

### High Priority
- **Elon Musk** (@elonmusk) - Tesla & SpaceX CEO
- **Michael Saylor** (@saylor) - MicroStrategy Executive Chairman  
- **Vitalik Buterin** (@VitalikButerin) - Ethereum Co-founder
- **Changpeng Zhao** (@cz_binance) - Binance Co-founder
- **Nayib Bukele** (@nayibbukele) - President of El Salvador

### Medium Priority
- **Brian Armstrong** (@brian_armstrong) - Coinbase CEO
- **Jack Dorsey** (@jack) - Block (formerly Square) CEO
- **Anthony Pompliano** (@APompliano) - Investor & Host
- **Naval** (@naval) - AngelList Co-founder
- **Raoul Pal** (@RaoulGMI) - Real Vision CEO
- **Andreas Antonopoulos** (@aantonop) - Bitcoin Educator

## Commands

### Whale Alerts Commands (`/whale`)

- `/whale setup [channel]` - Enable whale alerts in the specified channel (admin only)
- `/whale stop` - Disable whale alerts for this server (admin only)
- `/whale recent [limit]` - Get recent whale transactions (1-10, default: 5)
- `/whale status` - Show whale alerts configuration and status

### Guru Tweets Commands (`/guru`)

- `/guru setup [channel]` - Enable guru tweets in the specified channel (admin only)
- `/guru stop` - Disable guru tweets for this server (admin only)
- `/guru accounts` - List all monitored crypto guru accounts
- `/guru status` - Show guru tweets configuration and status

## Setup Requirements

### Environment Variables

To enable real-time monitoring, add these environment variables to your `.env` file:

```env
# Twitter API v2 (for guru tweets)
TWITTER_API_KEY=your_twitter_api_key
TWITTER_API_SECRET=your_twitter_api_secret
TWITTER_ACCESS_TOKEN=your_twitter_access_token
TWITTER_ACCESS_TOKEN_SECRET=your_twitter_access_token_secret

# Whale Alert API (for whale transactions)
WHALE_ALERT_API_KEY=your_whale_alert_api_key
```

### Getting API Keys

#### Twitter API v2
1. Apply for Twitter Developer access at [developer.twitter.com](https://developer.twitter.com)
2. Create a new app and generate API keys
3. Ensure your app has read permissions

#### Whale Alert API
1. Sign up at [whale-alert.io](https://whale-alert.io)
2. Subscribe to their API service
3. Get your API key from the dashboard

> **Note**: Without API keys, the feature will work with mock/demo data for testing purposes.

## How It Works

### Whale Alerts Monitoring
1. **API Integration**: Connects to Whale Alert API to fetch real-time transaction data
2. **Smart Filtering**: Filters transactions based on USD value and transaction type
3. **Rate Limiting**: Prevents spam with configurable intervals between alerts
4. **Rich Embeds**: Displays transaction details with color-coded severity

### Guru Tweets Monitoring  
1. **Twitter Integration**: Monitors specified crypto influencer accounts
2. **Content Analysis**: Analyzes tweets for crypto mentions and sentiment
3. **Engagement Scoring**: Prioritizes tweets based on likes, retweets, and replies
4. **Real-time Updates**: Checks for new tweets every 2 minutes

## Configuration

### Whale Alert Filtering
- **Minimum Amounts**: Configurable thresholds per alert type
- **Priority Tokens**: Lower thresholds for BTC, ETH, SOL, etc.
- **High Priority**: Transactions over $1M always trigger alerts

### Tweet Filtering
- **Engagement Threshold**: Minimum likes/retweets required
- **Priority Accounts**: Lower thresholds for high-priority influencers
- **Sentiment Focus**: Option to show only urgent/breaking news tweets
- **Token Mentions**: Must mention cryptocurrency tokens or have high engagement

## Rate Limiting

To prevent spam and respect Discord's limits:
- **Whale Alerts**: Maximum 1 alert per 5 minutes per server
- **Guru Tweets**: Maximum 1 tweet per 10 minutes per server
- **API Calls**: Respects Twitter and Whale Alert rate limits

## Mock Data Mode

When API keys are not configured, the feature operates in mock data mode:
- Generates realistic sample whale transactions
- Shows demo guru tweets with various sentiments
- Maintains all functionality for testing and demonstration

## Technical Details

### Dependencies
- `tweepy` - Twitter API v2 client
- `aiohttp` - Async HTTP requests for Whale Alert API
- `discord.py` - Discord bot framework

### Architecture
- **Models**: `WhaleAlert`, `GuruTweet` with rich embed generation
- **Services**: `TwitterMonitor`, `WhaleService`, `AlertsBroadcaster`
- **Cogs**: Discord command handlers and UI

### Error Handling
- Automatic retry on API failures
- Graceful degradation when services are unavailable
- Comprehensive logging for debugging

## Security & Privacy

- API keys are stored securely in environment variables
- No personal data is stored or logged
- Only public tweets and blockchain data are monitored
- Rate limiting protects against abuse

## Support

For issues or feature requests related to whale alerts:
1. Check the bot logs for error messages
2. Verify API keys are correctly configured
3. Ensure the bot has proper Discord permissions
4. Contact the bot administrator for assistance
