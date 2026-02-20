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
	return &NewsService{
		sources:  defaultSources("EN"),
		parser:   gofeed.NewParser(),
		running:  false,
		stopChan: make(chan struct{}),
	}
}

// defaultSources returns news sources for the given language/region
func defaultSources(lang string) []NewsSource {
	switch lang {
	case "ID":
		return []NewsSource{
			{Name: "Kompas", URL: "https://rss.kompas.com/kompas-terpopuler"},
			{Name: "Detik", URL: "https://rss.detik.com/index.php/detikcom"},
			{Name: "CNN Indonesia", URL: "https://www.cnnindonesia.com/nasional/rss"},
			{Name: "Tempo", URL: "https://rss.tempo.co/nasional"},
			{Name: "Liputan6", URL: "https://www.liputan6.com/rss"},
		}
	case "JP":
		return []NewsSource{
			{Name: "NHK News", URL: "https://www3.nhk.or.jp/rss/news/cat0.xml"},
			{Name: "Japan Times", URL: "https://www.japantimes.co.jp/feed/"},
			{Name: "Mainichi", URL: "https://mainichi.jp/rss/etc/mainichi-flash.rss"},
			{Name: "Asahi", URL: "https://www.asahi.com/rss/asahi/newsheadlines.rdf"},
			{Name: "NHK World", URL: "https://www3.nhk.or.jp/nhkworld/en/news/rss/index.xml"},
		}
	case "KR":
		return []NewsSource{
			{Name: "Yonhap News", URL: "https://en.yna.co.kr/RSS/news.xml"},
			{Name: "Korea Herald", URL: "http://www.koreaherald.com/common/rss_xml.php?ct=102"},
			{Name: "Korea Times", URL: "https://www.koreatimes.co.kr/www/rss/nation.xml"},
			{Name: "Arirang", URL: "https://www.arirang.com/rss"},
			{Name: "KBS World", URL: "https://world.kbs.co.kr/rss/rss_news.htm?lang=e"},
		}
	case "ZH":
		return []NewsSource{
			{Name: "BBC Chinese", URL: "https://feeds.bbci.co.uk/zhongwen/simp/rss.xml"},
			{Name: "RFI Chinese", URL: "https://www.rfi.fr/cn/rss"},
			{Name: "VOA Chinese", URL: "https://www.voachinese.com/api/zyrtemoj"},
			{Name: "DW Chinese", URL: "https://rss.dw.com/xml/rss-chi-all"},
			{Name: "FT Chinese", URL: "https://www.ftchinese.com/rss/news"},
		}
	default: // EN
		return []NewsSource{
			{Name: "BBC News", URL: "http://feeds.bbci.co.uk/news/rss.xml"},
			{Name: "CNN", URL: "http://rss.cnn.com/rss/edition.rss"},
			{Name: "Reuters", URL: "https://news.google.com/rss/search?q=reuters+world+news&hl=en-US&gl=US&ceid=US:en"},
			{Name: "TechCrunch", URL: "https://techcrunch.com/feed/"},
			{Name: "The Verge", URL: "https://www.theverge.com/rss/index.xml"},
		}
	}
}

// FetchLatestByLang fetches the latest news from sources for the specified language
func (s *NewsService) FetchLatestByLang(ctx context.Context, limit int, lang string) ([]*entity.NewsArticle, error) {
	sources := defaultSources(lang)
	return s.fetchFromSources(ctx, sources, limit)
}

// FetchLatest fetches the latest news from all default (EN) sources
func (s *NewsService) FetchLatest(ctx context.Context, limit int) ([]*entity.NewsArticle, error) {
	return s.fetchFromSources(ctx, s.sources, limit)
}

// fetchFromSources fetches news from the given sources
func (s *NewsService) fetchFromSources(ctx context.Context, sources []NewsSource, limit int) ([]*entity.NewsArticle, error) {
	items := make([]*entity.NewsArticle, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)
		go func(src NewsSource) {
			defer wg.Done()

			feed, err := s.parser.ParseURLWithContext(src.URL, ctx)
			if err != nil {
				return
			}

			for i, item := range feed.Items {
				if i >= limit/len(sources) {
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
