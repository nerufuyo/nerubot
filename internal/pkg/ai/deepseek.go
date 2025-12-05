package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DeepSeekProvider implements AIProvider for DeepSeek
type DeepSeekProvider struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

// NewDeepSeekProvider creates a new DeepSeek AI provider
func NewDeepSeekProvider(apiKey string) *DeepSeekProvider {
	return &DeepSeekProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		model: "deepseek-chat",
	}
}

// Name returns the provider name
func (d *DeepSeekProvider) Name() string {
	return "DeepSeek"
}

// IsAvailable checks if DeepSeek is configured
func (d *DeepSeekProvider) IsAvailable() bool {
	return d.apiKey != ""
}

// Chat sends a message to DeepSeek and returns the response
func (d *DeepSeekProvider) Chat(ctx context.Context, messages []Message) (string, error) {
	if !d.IsAvailable() {
		return "", fmt.Errorf("deepseek API key not configured")
	}

	// Convert messages to DeepSeek format
	deepseekMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		deepseekMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	requestBody := map[string]interface{}{
		"model":       d.model,
		"messages":    deepseekMessages,
		"max_tokens":  1024,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("deepseek API returned status %d", resp.StatusCode)
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return response.Choices[0].Message.Content, nil
}
