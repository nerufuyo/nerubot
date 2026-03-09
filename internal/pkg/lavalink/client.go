package lavalink

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// EventHandler is called when a Lavalink player event occurs.
type EventHandler func(player disgolink.Player, event lavalink.Event)

// Client wraps DisGoLink and manages the Lavalink connection.
type Client struct {
	Link   disgolink.Client
	logger *logger.Logger
	mu     sync.RWMutex

	// Callbacks
	onTrackStart     func(player disgolink.Player, event lavalink.TrackStartEvent)
	onTrackEnd       func(player disgolink.Player, event lavalink.TrackEndEvent)
	onTrackException func(player disgolink.Player, event lavalink.TrackExceptionEvent)
	onTrackStuck     func(player disgolink.Player, event lavalink.TrackStuckEvent)
}

// New creates a new Lavalink client wrapper.
func New(session *discordgo.Session, log *logger.Logger) *Client {
	botUserID := snowflake.MustParse(session.State.User.ID)

	c := &Client{
		logger: log,
	}

	c.Link = disgolink.New(botUserID,
		disgolink.WithListenerFunc(c.handleTrackStart),
		disgolink.WithListenerFunc(c.handleTrackEnd),
		disgolink.WithListenerFunc(c.handleTrackException),
		disgolink.WithListenerFunc(c.handleTrackStuck),
	)

	return c
}

// AddNode connects a Lavalink node.
func (c *Client) AddNode(ctx context.Context, name, address, password string, secure bool) error {
	node, err := c.Link.AddNode(ctx, disgolink.NodeConfig{
		Name:     name,
		Address:  address,
		Password: password,
		Secure:   secure,
	})
	if err != nil {
		return fmt.Errorf("failed to add Lavalink node: %w", err)
	}
	c.logger.Info("Lavalink node connected", "name", node.Config().Name, "address", address)
	return nil
}

// OnTrackStart registers a callback for track start events.
func (c *Client) OnTrackStart(fn func(player disgolink.Player, event lavalink.TrackStartEvent)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onTrackStart = fn
}

// OnTrackEnd registers a callback for track end events.
func (c *Client) OnTrackEnd(fn func(player disgolink.Player, event lavalink.TrackEndEvent)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onTrackEnd = fn
}

// OnTrackException registers a callback for track exception events.
func (c *Client) OnTrackException(fn func(player disgolink.Player, event lavalink.TrackExceptionEvent)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onTrackException = fn
}

// OnTrackStuck registers a callback for track stuck events.
func (c *Client) OnTrackStuck(fn func(player disgolink.Player, event lavalink.TrackStuckEvent)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onTrackStuck = fn
}

// HandleVoiceStateUpdate forwards voice state updates to DisGoLink.
// Must be called from the discordgo VoiceStateUpdate handler.
func (c *Client) HandleVoiceStateUpdate(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	// Only forward events for the bot itself
	if event.UserID != s.State.User.ID {
		return
	}

	guildID := snowflake.MustParse(event.GuildID)
	var channelID *snowflake.ID
	if event.ChannelID != "" {
		id := snowflake.MustParse(event.ChannelID)
		channelID = &id
	}

	c.Link.OnVoiceStateUpdate(context.TODO(), guildID, channelID, event.SessionID)
}

// HandleVoiceServerUpdate forwards voice server updates to DisGoLink.
// Must be called from the discordgo VoiceServerUpdate handler.
func (c *Client) HandleVoiceServerUpdate(s *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	guildID := snowflake.MustParse(event.GuildID)

	c.Link.OnVoiceServerUpdate(context.TODO(), guildID, event.Token, event.Endpoint)
}

// LoadTracks searches for tracks on Lavalink.
func (c *Client) LoadTracks(ctx context.Context, query string) (*lavalink.LoadResult, error) {
	node := c.Link.BestNode()
	if node == nil {
		return nil, fmt.Errorf("no available Lavalink node")
	}
	result, err := node.LoadTracks(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to load tracks: %w", err)
	}
	return result, nil
}

// Player returns the Lavalink player for a guild.
func (c *Client) Player(guildID string) disgolink.Player {
	id := snowflake.MustParse(guildID)
	return c.Link.Player(id)
}

// ExistingPlayer returns the player for a guild if it exists.
func (c *Client) ExistingPlayer(guildID string) disgolink.Player {
	id := snowflake.MustParse(guildID)
	return c.Link.ExistingPlayer(id)
}

// RemovePlayer removes and destroys the player for a guild.
func (c *Client) RemovePlayer(guildID string) {
	id := snowflake.MustParse(guildID)
	c.Link.RemovePlayer(id)
}

// GetLyrics fetches lyrics for the currently playing track via Lavalink REST API (LavaLyrics plugin).
func (c *Client) GetLyrics(ctx context.Context, guildID string) (string, error) {
	node := c.Link.BestNode()
	if node == nil {
		return "", fmt.Errorf("no available Lavalink node")
	}

	cfg := node.Config()
	url := fmt.Sprintf("%s/v4/sessions/%s/players/%s/track/lyrics", cfg.RestURL(), node.SessionID(), guildID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", cfg.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("lyrics request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("lyrics not available (status %d)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// LavaLyrics returns: { "sourceName": "...", "provider": "...", "text": "...", "lines": [...] }
	var lyricsResp struct {
		Text  string `json:"text"`
		Lines []struct {
			Line string `json:"line"`
		} `json:"lines"`
	}
	if err := json.Unmarshal(body, &lyricsResp); err != nil {
		return "", fmt.Errorf("failed to parse lyrics: %w", err)
	}

	// Prefer full text, fall back to joining lines
	if lyricsResp.Text != "" {
		return lyricsResp.Text, nil
	}

	if len(lyricsResp.Lines) > 0 {
		var text string
		for _, l := range lyricsResp.Lines {
			text += l.Line + "\n"
		}
		return text, nil
	}

	return "", fmt.Errorf("no lyrics found")
}

// --- Internal event handlers ---

func (c *Client) handleTrackStart(p disgolink.Player, event lavalink.TrackStartEvent) {
	c.mu.RLock()
	fn := c.onTrackStart
	c.mu.RUnlock()
	if fn != nil {
		fn(p, event)
	}
}

func (c *Client) handleTrackEnd(p disgolink.Player, event lavalink.TrackEndEvent) {
	c.mu.RLock()
	fn := c.onTrackEnd
	c.mu.RUnlock()
	if fn != nil {
		fn(p, event)
	}
}

func (c *Client) handleTrackException(p disgolink.Player, event lavalink.TrackExceptionEvent) {
	c.mu.RLock()
	fn := c.onTrackException
	c.mu.RUnlock()
	if fn != nil {
		fn(p, event)
	}
}

func (c *Client) handleTrackStuck(p disgolink.Player, event lavalink.TrackStuckEvent) {
	c.mu.RLock()
	fn := c.onTrackStuck
	c.mu.RUnlock()
	if fn != nil {
		fn(p, event)
	}
}
