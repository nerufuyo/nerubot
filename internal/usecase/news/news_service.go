package news

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nerufuyo/nerubot/internal/entity"
)

// NewsSource represents a news RSS feed source
type NewsSource struct {
	Name string
	URL  string
}

// NewsService handles news aggregation and publishing
type NewsService struct {
	sources  []NewsSource
	parser   *gofeed.Parser
	running  bool
	stopChan chan struct{}
	mu       sync.RWMutex
}

// NewNewsService creates a new news service
func NewNewsService() *NewsService {
	sources := []NewsSource{
		{Name: "BBC News", URL: "http://feeds.bbci.co.uk/news/rss.xml"},
		{Name: "CNN", URL: "http://rss.cnn.com/rss/edition.rss"},
		{Name: "Reuters", URL: "https://news.google.com/rss/search?q=reuters+world+news&hl=en-US&gl=US&ceid=US:en"},
		{Name: "TechCrunch", URL: "https://techcrunch.com/feed/"},
		{Name: "The Verge", URL: "https://www.theverge.com/rss/index.xml"},
	}

	return &NewsService{
		sources:  sources,
		parser:   gofeed.NewParser(),
		running:  false,
		stopChan: make(chan struct{}),
	}
}

// FetchLatest fetches the latest news from all sources
func (s *NewsService) FetchLatest(ctx context.Context, limit int) ([]*entity.NewsArticle, error) {
	items := make([]*entity.NewsArticle, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, source := range s.sources {
		wg.Add(1)
		go func(src NewsSource) {
			defer wg.Done()

			feed, err := s.parser.ParseURLWithContext(src.URL, ctx)
			if err != nil {
				return
			}

			for i, item := range feed.Items {
				if i >= limit/len(s.sources) {
					break
				}

				newsItem := &entity.NewsArticle{
					Title:       item.Title,
					Description: item.Description,
					URL:         item.Link,
					Source:      src.Name,
					PublishedAt: time.Now(),
					FetchedAt:   time.Now(),
				}

				if item.PublishedParsed != nil {
					newsItem.PublishedAt = *item.PublishedParsed
				}

				mu.Lock()
				items = append(items, newsItem)
				mu.Unlock()
			}
		}(source)
	}

	wg.Wait()

	// Sort by published date (most recent first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].PublishedAt.After(items[j].PublishedAt)
	})

	// Limit to requested count
	if len(items) > limit {
		items = items[:limit]
	}

	return items, nil
}

// Start begins auto-publishing news at specified interval
func (s *NewsService) Start(interval time.Duration, publishFunc func([]*entity.NewsArticle)) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("news service already running")
	}
	s.running = true
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				items, err := s.FetchLatest(ctx, 5)
				cancel()

				if err == nil && len(items) > 0 {
					publishFunc(items)
				}

			case <-s.stopChan:
				return
			}
		}
	}()

	return nil
}

// Stop stops the auto-publishing
func (s *NewsService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		close(s.stopChan)
		s.running = false
	}
}

// IsRunning returns whether the service is currently running
func (s *NewsService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetSources returns the list of configured news sources
func (s *NewsService) GetSources() []NewsSource {
	return s.sources
}

// AddSource adds a new news source
func (s *NewsService) AddSource(name, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sources = append(s.sources, NewsSource{
		Name: name,
		URL:  url,
	})
}
