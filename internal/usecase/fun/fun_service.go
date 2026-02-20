package fun

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"github.com/nerufuyo/nerubot/internal/repository"
)

// SendFunc is a callback to send a message to a Discord channel.
type SendFunc func(channelID string, embed *FunEmbed)

// FunEmbed holds data for an embed message.
type FunEmbed struct {
	Title       string
	Description string
	ImageURL    string
	Footer      string
	Color       int
	URL         string
}

// FunService manages dad jokes, memes, and their scheduled delivery.
type FunService struct {
	mu          sync.RWMutex
	repo        *repository.GuildConfigRepository
	logger      *logger.Logger
	httpClient  *http.Client
	sendFn      SendFunc
	stopCh      chan struct{}
	wg          sync.WaitGroup
	rng         *rand.Rand

	// Track last fire times per guild to avoid duplicates
	lastJoke map[string]time.Time
	lastMeme map[string]time.Time
}

// NewFunService creates a new fun service.
func NewFunService() *FunService {
	return &FunService{
		repo:   repository.NewGuildConfigRepository(),
		logger: logger.New("fun"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		stopCh:   make(chan struct{}),
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		lastJoke: make(map[string]time.Time),
		lastMeme: make(map[string]time.Time),
	}
}

// SetSendFunc sets the callback used to post embeds.
func (s *FunService) SetSendFunc(fn SendFunc) {
	s.sendFn = fn
}

// Start begins the background scheduler that checks every minute for scheduled jokes/memes.
func (s *FunService) Start() {
	s.wg.Add(1)
	go s.loop()
	s.logger.Info("Fun service scheduler started")
}

// Stop gracefully shuts down the scheduler.
func (s *FunService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	s.logger.Info("Fun service scheduler stopped")
}

func (s *FunService) loop() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.checkScheduled()
		}
	}
}

func (s *FunService) checkScheduled() {
	configs, err := s.repo.GetAll()
	if err != nil {
		s.logger.Warn("Failed to load guild configs for scheduler", "error", err)
		return
	}

	now := time.Now()

	for _, cfg := range configs {
		// Check dad jokes schedule
		if cfg.DadJokeChannelID != "" && cfg.DadJokeInterval > 0 {
			lastFired, ok := s.lastJoke[cfg.GuildID]
			interval := time.Duration(cfg.DadJokeInterval) * time.Minute
			if !ok || now.Sub(lastFired) >= interval {
				s.lastJoke[cfg.GuildID] = now
				go s.sendScheduledJoke(cfg.DadJokeChannelID)
			}
		}

		// Check memes schedule
		if cfg.MemeChannelID != "" && cfg.MemeInterval > 0 {
			lastFired, ok := s.lastMeme[cfg.GuildID]
			interval := time.Duration(cfg.MemeInterval) * time.Minute
			if !ok || now.Sub(lastFired) >= interval {
				s.lastMeme[cfg.GuildID] = now
				go s.sendScheduledMeme(cfg.MemeChannelID)
			}
		}
	}
}

func (s *FunService) sendScheduledJoke(channelID string) {
	joke, err := s.FetchDadJoke()
	if err != nil {
		s.logger.Warn("Scheduled dad joke fetch failed", "error", err)
		return
	}

	embed := &FunEmbed{
		Title:       "🤣 Dad Joke of the Hour",
		Description: joke.Punchline,
		Footer:      "Powered by icanhazdadjoke.com",
		Color:       0xFFD700, // gold
	}
	if s.sendFn != nil {
		s.sendFn(channelID, embed)
	}
}

func (s *FunService) sendScheduledMeme(channelID string) {
	meme, err := s.FetchMeme()
	if err != nil {
		s.logger.Warn("Scheduled meme fetch failed", "error", err)
		return
	}

	embed := &FunEmbed{
		Title:    "😂 " + meme.Title,
		ImageURL: meme.URL,
		Footer:   fmt.Sprintf("r/%s • by u/%s", meme.Subreddit, meme.Author),
		Color:    0xFF4500, // reddit orange
		URL:      meme.PostLink,
	}
	if s.sendFn != nil {
		s.sendFn(channelID, embed)
	}
}

// --- Dad Joke API ---

// icanhazdadjoke response
type dadJokeAPIResponse struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

// FetchDadJoke fetches a random clean dad joke from icanhazdadjoke.com.
func (s *FunService) FetchDadJoke() (*entity.DadJoke, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://icanhazdadjoke.com/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NeruBot Discord Bot (https://github.com/nerufuyo/nerubot)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dad joke: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dad joke API returned status %d", resp.StatusCode)
	}

	var apiResp dadJokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode dad joke response: %w", err)
	}

	return &entity.DadJoke{
		ID:        apiResp.ID,
		Punchline: apiResp.Joke,
		Source:    "icanhazdadjoke.com",
		FetchedAt: time.Now(),
	}, nil
}

// --- Meme API ---

// memeAPIResponse is the response from the meme-api.com endpoint.
type memeAPIResponse struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	NSFW      bool   `json:"nsfw"`
	Spoiler   bool   `json:"spoiler"`
	Author    string `json:"author"`
	Ups       int    `json:"ups"`
}

// memeAPIMultiResponse is the response for multiple memes.
type memeAPIMultiResponse struct {
	Count int               `json:"count"`
	Memes []memeAPIResponse `json:"memes"`
}

// FetchMeme fetches a random SFW meme from Reddit via meme-api.com.
func (s *FunService) FetchMeme() (*entity.Meme, error) {
	// Use clean meme subreddits only
	subreddits := []string{"memes", "dankmemes", "wholesomememes", "me_irl", "ProgrammerHumor"}
	sub := subreddits[s.rng.Intn(len(subreddits))]

	url := fmt.Sprintf("https://meme-api.com/gimme/%s", sub)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "NeruBot Discord Bot (https://github.com/nerufuyo/nerubot)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meme: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meme API returned status %d", resp.StatusCode)
	}

	var apiResp memeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode meme response: %w", err)
	}

	// Reject NSFW content
	if apiResp.NSFW || apiResp.Spoiler {
		// Try once more with wholesomememes as fallback
		return s.fetchMemeFromSubreddit("wholesomememes")
	}

	return &entity.Meme{
		Title:     apiResp.Title,
		URL:       apiResp.URL,
		PostLink:  apiResp.PostLink,
		Subreddit: apiResp.Subreddit,
		Author:    apiResp.Author,
		NSFW:      apiResp.NSFW,
		FetchedAt: time.Now(),
	}, nil
}

func (s *FunService) fetchMemeFromSubreddit(subreddit string) (*entity.Meme, error) {
	url := fmt.Sprintf("https://meme-api.com/gimme/%s", subreddit)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "NeruBot Discord Bot (https://github.com/nerufuyo/nerubot)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meme: %w", err)
	}
	defer resp.Body.Close()

	var apiResp memeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode meme response: %w", err)
	}

	if apiResp.NSFW {
		return nil, fmt.Errorf("could not find a SFW meme")
	}

	return &entity.Meme{
		Title:     apiResp.Title,
		URL:       apiResp.URL,
		PostLink:  apiResp.PostLink,
		Subreddit: apiResp.Subreddit,
		Author:    apiResp.Author,
		NSFW:      false,
		FetchedAt: time.Now(),
	}, nil
}

// --- Guild Config helpers ---

// GetGuildConfig retrieves (or creates) a guild config.
func (s *FunService) GetGuildConfig(guildID, guildName string) (*entity.GuildConfig, error) {
	cfg, err := s.repo.Get(guildID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = entity.NewGuildConfig(guildID, guildName)
	}
	return cfg, nil
}

// SaveGuildConfig persists a guild config.
func (s *FunService) SaveGuildConfig(cfg *entity.GuildConfig) error {
	return s.repo.Save(cfg)
}
