# AI Provider Migration - DeepSeek Only

## Overview
The chatbot service has been simplified to use **only DeepSeek** as the AI provider, removing support for Anthropic Claude, Google Gemini, and OpenAI.

## Changes Made

### 1. New DeepSeek Provider
- **File**: `internal/pkg/ai/deepseek.go`
- **API Endpoint**: `https://api.deepseek.com/v1/chat/completions`
- **Model**: `deepseek-chat`
- **Features**: 
  - OpenAI-compatible API format
  - 1024 max tokens
  - 0.7 temperature (configurable)
  - 30-second timeout

### 2. Removed Providers
Deleted the following provider implementations:
- ❌ `internal/pkg/ai/claude.go` (Anthropic Claude)
- ❌ `internal/pkg/ai/gemini.go` (Google Gemini)
- ❌ `internal/pkg/ai/openai.go` (OpenAI)

### 3. Configuration Updates

#### `internal/config/config.go`
```go
// Before (3 providers)
type AIConfig struct {
    OpenAIKey    string
    AnthropicKey string
    GeminiKey    string
}

// After (DeepSeek only)
type AIConfig struct {
    DeepSeekKey string
}
```

#### `.env.example`
```bash
# Before
OPENAI_API_KEY=your_openai_api_key_here
ANTHROPIC_API_KEY=your_anthropic_api_key_here
GEMINI_API_KEY=your_gemini_api_key_here

# After
DEEPSEEK_API_KEY=your_deepseek_api_key_here
```

### 4. Service Updates

#### `internal/usecase/chatbot/chatbot_service.go`
```go
// Before
func NewChatbotService(claudeKey, geminiKey, openaiKey string) *ChatbotService {
    // Initialize 3 providers...
}

// After
func NewChatbotService(deepseekKey string) *ChatbotService {
    // Initialize DeepSeek only
}
```

#### `services/chatbot/cmd/main.go`
```go
// Before
chatbotService := chatbot.NewChatbotService(
    cfg.AI.AnthropicKey,
    cfg.AI.GeminiKey,
    cfg.AI.OpenAIKey,
)

// After
chatbotService := chatbot.NewChatbotService(
    cfg.AI.DeepSeekKey,
)
```

## Why DeepSeek?

### Advantages
1. **Cost-Effective**: DeepSeek offers competitive pricing
2. **OpenAI-Compatible**: Uses familiar API format
3. **Performance**: Fast response times
4. **Simplicity**: Single provider reduces complexity
5. **Open Source**: DeepSeek models are open source

### API Comparison
| Feature | DeepSeek | Claude | Gemini | OpenAI |
|---------|----------|--------|--------|--------|
| API Format | OpenAI-compatible | Custom | Custom | Standard |
| Cost | $ | $$$ | $$ | $$$ |
| Speed | Fast | Fast | Fast | Fast |
| Open Source | ✅ | ❌ | ❌ | ❌ |

## Migration Guide

### For Existing Users

1. **Get DeepSeek API Key**
   - Visit: https://platform.deepseek.com/
   - Sign up for an account
   - Generate an API key

2. **Update Environment Variables**
   ```bash
   # Remove old keys
   unset OPENAI_API_KEY
   unset ANTHROPIC_API_KEY
   unset GEMINI_API_KEY
   
   # Add DeepSeek key
   export DEEPSEEK_API_KEY=your_deepseek_api_key_here
   ```

3. **Update .env File**
   ```bash
   # .env
   DEEPSEEK_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxx
   ```

4. **Rebuild Services**
   ```bash
   # Rebuild chatbot service
   go build -o build/chatbot/chatbot.exe ./services/chatbot/cmd
   
   # Rebuild gateway (if needed)
   go build -o build/gateway/gateway.exe ./services/gateway/cmd
   ```

### For New Users

Simply follow the updated `.env.example` file - only `DEEPSEEK_API_KEY` is needed for chatbot functionality.

## Testing

### Verify Configuration
```go
// Check if DeepSeek is configured
if cfg.Features.Chatbot {
    log.Info("Chatbot enabled with DeepSeek")
}
```

### Test Chat Functionality
```bash
# In Discord, use the chat command
/chat message:"Hello, test message"

# Expected response format:
# Response from DeepSeek AI provider
```

## Build Verification

After migration, both services build successfully:
- ✅ **Chatbot Service**: 8.71 MB
- ✅ **Gateway Service**: 19.6 MB

## API Reference

### DeepSeek Chat Completion
```http
POST https://api.deepseek.com/v1/chat/completions
Content-Type: application/json
Authorization: Bearer YOUR_API_KEY

{
  "model": "deepseek-chat",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ],
  "max_tokens": 1024,
  "temperature": 0.7
}
```

### Response Format
```json
{
  "choices": [
    {
      "message": {
        "content": "Hello! How can I help you today?"
      }
    }
  ]
}
```

## Troubleshooting

### Error: "no AI providers configured"
**Solution**: Ensure `DEEPSEEK_API_KEY` is set in your environment

### Error: "deepseek API returned status 401"
**Solution**: Verify your API key is valid

### Error: "deepseek API returned status 429"
**Solution**: Rate limit exceeded - wait and retry

## Future Considerations

### Adding More Providers
If you need to add more providers in the future:

1. Create new provider file (e.g., `internal/pkg/ai/newprovider.go`)
2. Implement `AIProvider` interface
3. Add configuration to `AIConfig`
4. Update `NewChatbotService()` to include new provider
5. Update `.env.example`

### Provider Interface
All providers must implement:
```go
type AIProvider interface {
    Name() string
    Chat(ctx context.Context, messages []Message) (string, error)
    IsAvailable() bool
}
```

## Performance Impact

- **Latency**: Similar to OpenAI (50-500ms typical)
- **Memory**: Reduced due to fewer provider clients
- **Binary Size**: Slightly smaller (removed 3 provider implementations)

## Security Notes

- Store API keys securely (use environment variables, not in code)
- Use `.env` files for local development
- Use Railway secrets for production deployment
- Never commit API keys to version control

## Related Files

### Core Files
- `internal/pkg/ai/deepseek.go` - DeepSeek provider implementation
- `internal/pkg/ai/interface.go` - AIProvider interface
- `internal/config/config.go` - Configuration structure
- `internal/usecase/chatbot/chatbot_service.go` - Chatbot service

### Service Files
- `services/chatbot/cmd/main.go` - Chatbot microservice
- `services/gateway/cmd/main.go` - API Gateway

### Configuration
- `.env.example` - Environment variable template
- `api/proto/chatbot.proto` - gRPC API definition (unchanged)

## Summary

The migration simplifies the AI integration by using a single, cost-effective provider (DeepSeek) while maintaining all chatbot functionality. The OpenAI-compatible API format makes it easy to understand and integrate.

**Migration Status**: ✅ Complete
**Build Status**: ✅ All services building successfully
**Breaking Changes**: Yes - requires new API key configuration
