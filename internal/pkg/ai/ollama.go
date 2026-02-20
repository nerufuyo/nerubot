package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient communicates with an Ollama server.
type OllamaClient struct {
	baseURL    string
	httpClient *http.Client
}

// OllamaModel represents one model returned by /api/tags.
type OllamaModel struct {
	Name       string             `json:"name"`
	Model      string             `json:"model"`
	ModifiedAt string             `json:"modified_at"`
	Size       int64              `json:"size"`
	Digest     string             `json:"digest"`
	Details    OllamaModelDetails `json:"details"`
}

// OllamaModelDetails holds the nested details of a model.
type OllamaModelDetails struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

// OllamaBenchResult holds benchmark metrics for a single run.
type OllamaBenchResult struct {
	Model            string
	Prompt           string
	Response         string
	TotalDuration    time.Duration // wall-clock time from request to last token
	FirstTokenTime   time.Duration // time to first token (TTFT)
	TokenCount       int           // total tokens generated
	TokensPerSecond  float64       // generation speed
	PromptEvalCount  int           // tokens evaluated in the prompt
	EvalCount        int           // tokens generated
	LoadDuration     time.Duration // model load time reported by Ollama
	PromptEvalDur    time.Duration // prompt evaluation time reported by Ollama
	EvalDuration     time.Duration // generation time reported by Ollama
}

// NewOllamaClient creates a new OllamaClient.
func NewOllamaClient(baseURL string) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// IsAvailable checks whether the Ollama server is reachable.
func (o *OllamaClient) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ListModels fetches available models from the Ollama server.
func (o *OllamaClient) ListModels(ctx context.Context) ([]OllamaModel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var result struct {
		Models []OllamaModel `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Models, nil
}

// ollamaStreamChunk is a single line from the streaming /api/generate response.
type ollamaStreamChunk struct {
	Model              string `json:"model"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	TotalDuration      int64  `json:"total_duration"`       // nanoseconds
	LoadDuration       int64  `json:"load_duration"`        // nanoseconds
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"` // nanoseconds
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`        // nanoseconds
}

// Benchmark runs a benchmark against the specified model and returns metrics.
// It uses the streaming /api/generate endpoint to measure time-to-first-token.
func (o *OllamaClient) Benchmark(ctx context.Context, model, prompt string) (*OllamaBenchResult, error) {
	body := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": true,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var (
		fullResponse  string
		firstTokenAt  time.Time
		tokenCount    int
		lastChunk     ollamaStreamChunk
	)

	for {
		var chunk ollamaStreamChunk
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decode stream chunk: %w", err)
		}

		if chunk.Response != "" {
			tokenCount++
			if firstTokenAt.IsZero() {
				firstTokenAt = time.Now()
			}
			fullResponse += chunk.Response
		}

		if chunk.Done {
			lastChunk = chunk
			break
		}
	}

	totalDuration := time.Since(start)
	ttft := time.Duration(0)
	if !firstTokenAt.IsZero() {
		ttft = firstTokenAt.Sub(start)
	}

	// Prefer Ollama's own eval_count if available (more accurate token count)
	evalCount := lastChunk.EvalCount
	if evalCount == 0 {
		evalCount = tokenCount
	}

	tokPerSec := float64(0)
	if lastChunk.EvalDuration > 0 {
		tokPerSec = float64(lastChunk.EvalCount) / (float64(lastChunk.EvalDuration) / 1e9)
	} else if totalDuration.Seconds() > 0 {
		tokPerSec = float64(evalCount) / totalDuration.Seconds()
	}

	// Truncate response for display (max 200 chars)
	displayResp := fullResponse
	if len(displayResp) > 200 {
		displayResp = displayResp[:200] + "..."
	}

	return &OllamaBenchResult{
		Model:           model,
		Prompt:          prompt,
		Response:        displayResp,
		TotalDuration:   totalDuration,
		FirstTokenTime:  ttft,
		TokenCount:      evalCount,
		TokensPerSecond: tokPerSec,
		PromptEvalCount: lastChunk.PromptEvalCount,
		EvalCount:       lastChunk.EvalCount,
		LoadDuration:    time.Duration(lastChunk.LoadDuration),
		PromptEvalDur:   time.Duration(lastChunk.PromptEvalDuration),
		EvalDuration:    time.Duration(lastChunk.EvalDuration),
	}, nil
}

// FormatSize formats bytes into a human-readable string.
func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
