package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ClaudeProvider implements AIProvider for Anthropic Claude
type ClaudeProvider struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

// NewClaudeProvider creates a new Claude AI provider
func NewClaudeProvider(apiKey string) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		model: "claude-3-5-sonnet-20241022",
	}
}

// Name returns the provider name
func (c *ClaudeProvider) Name() string {
	return "Claude"
}

// IsAvailable checks if Claude is configured
func (c *ClaudeProvider) IsAvailable() bool {
	return c.apiKey != ""
}

// Chat sends a message to Claude and returns the response
func (c *ClaudeProvider) Chat(ctx context.Context, messages []Message) (string, error) {
	if !c.IsAvailable() {
		return "", fmt.Errorf("claude API key not configured")
	}

	// Convert messages to Claude format
	claudeMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		claudeMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	requestBody := map[string]interface{}{
		"model":      c.model,
		"max_tokens": 1024,
		"messages":   claudeMessages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API returned status %d", resp.StatusCode)
	}

	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return response.Content[0].Text, nil
}
