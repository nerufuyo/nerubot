package lavalink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// Client represents a Lavalink client connection
type Client struct {
	host       string
	port       int
	password   string
	userID     string
	httpClient *http.Client
	logger     *logger.Logger
}

// TrackInfo represents track information from Lavalink
type TrackInfo struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Length   int64  `json:"length"`
	Identifier string `json:"identifier"`
	IsStream bool   `json:"isStream"`
	URI      string `json:"uri"`
}

// Track represents an encoded track from Lavalink
type Track struct {
	Encoded string    `json:"encoded"`
	Info    TrackInfo `json:"info"`
}

// SearchResult represents search results from Lavalink
type SearchResult struct {
	Tracks       []Track `json:"tracks"`
	Playlists    []struct{} `json:"playlists"`
	LoadType     string `json:"loadType"`
}

// PlayRequest represents a play request for Lavalink
type PlayRequest struct {
	Track string `json:"track"`
	NoReplace bool `json:"noReplace,omitempty"`
}

// NewClient creates a new Lavalink client
func NewClient(host string, port int, password, userID string) *Client {
	return &Client{
		host: host,
		port: port,
		password: password,
		userID: userID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger.New("lavalink"),
	}
}

// SearchTracks searches for tracks on Lavalink
func (c *Client) SearchTracks(query string) ([]Track, error) {
	// Properly encode the query parameter
	encodedQuery := url.QueryEscape(query)
	baseURL := fmt.Sprintf("http://%s:%d/v4/loadtracks", c.host, c.port)
	fullURL := fmt.Sprintf("%s?identifier=%s", baseURL, encodedQuery)
	
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Log response for debugging
	c.logger.Info("Lavalink response", "status", resp.StatusCode, "body_length", len(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lavalink returned status %d: %s", resp.StatusCode, string(body))
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error("Failed to unmarshal response", "error", err, "body", string(body))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Tracks, nil
}

// Play sends a play command to Lavalink
func (c *Client) Play(guildID, sessionID, channelID string, track Track) error {
	url := fmt.Sprintf("http://%s:%d/v4/sessions/%s/players?guildId=%s", c.host, c.port, sessionID, guildID)

	payload := map[string]interface{}{
		"track": track.Encoded,
		"guild_id": guildID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("play request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("play request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	c.logger.Info("Playing track", "guild", guildID, "track", track.Info.Title)
	return nil
}

// Stop sends a stop command to Lavalink
func (c *Client) Stop(guildID, sessionID string) error {
	url := fmt.Sprintf("http://%s:%d/v4/sessions/%s/players?guildId=%s", c.host, c.port, sessionID, guildID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("stop request failed: %w", err)
	}
	defer resp.Body.Close()

	c.logger.Info("Stopped playback", "guild", guildID)
	return nil
}

// Pause sends a pause command to Lavalink
func (c *Client) Pause(guildID, sessionID string, pause bool) error {
	url := fmt.Sprintf("http://%s:%d/v4/sessions/%s/players?guildId=%s", c.host, c.port, sessionID, guildID)

	payload := map[string]interface{}{
		"paused": pause,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("pause request failed: %w", err)
	}
	defer resp.Body.Close()

	action := "paused"
	if !pause {
		action = "resumed"
	}
	c.logger.Info("Playback "+action, "guild", guildID)
	return nil
}

// UpdatePlayer updates player settings (like volume)
func (c *Client) UpdatePlayer(guildID, sessionID string, volume int) error {
	url := fmt.Sprintf("http://%s:%d/v4/sessions/%s/players?guildId=%s", c.host, c.port, sessionID, guildID)

	payload := map[string]interface{}{
		"volume": volume,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("update request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// JoinVoice tells Lavalink to join a voice channel
func (c *Client) JoinVoice(guildID, sessionID, channelID, track string) error {
	url := fmt.Sprintf("http://%s:%d/v4/sessions/%s/players?guildId=%s", c.host, c.port, sessionID, guildID)

	payload := map[string]interface{}{
		"voiceChannel": channelID,
		"track":        track,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("join voice request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("join voice request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	c.logger.Info("Joined voice channel via Lavalink", "guild", guildID, "channel", channelID)
	return nil
}
