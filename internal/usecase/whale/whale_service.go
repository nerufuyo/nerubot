package whale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
)

// WhaleService handles whale transaction monitoring
type WhaleService struct {
	apiKey      string
	httpClient  *http.Client
	running     bool
	stopChan    chan struct{}
	mu          sync.RWMutex
	minAmount   float64
	seenTx       map[string]time.Time
	seenTTL      time.Duration
}

// NewWhaleService creates a new whale alert service
func NewWhaleService(apiKey string) *WhaleService {
	return &WhaleService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		running:   false,
		stopChan:  make(chan struct{}),
		minAmount: 1000000, // $1M minimum
		seenTx:    make(map[string]time.Time),
		seenTTL:   24 * time.Hour,
	}
}

func (s *WhaleService) filterNewTransactions(transactions []*entity.WhaleTransaction) []*entity.WhaleTransaction {
	now := time.Now()
	threshold := now.Add(-s.seenTTL)

	s.mu.Lock()
	defer s.mu.Unlock()

	for hash, ts := range s.seenTx {
		if ts.Before(threshold) {
			delete(s.seenTx, hash)
		}
	}

	fresh := make([]*entity.WhaleTransaction, 0, len(transactions))
	for _, tx := range transactions {
		if tx.Hash == "" {
			fresh = append(fresh, tx)
			continue
		}
		if _, exists := s.seenTx[tx.Hash]; exists {
			continue
		}
		s.seenTx[tx.Hash] = now
		fresh = append(fresh, tx)
	}

	return fresh
}

// FetchTransactions fetches recent whale transactions
func (s *WhaleService) FetchTransactions(ctx context.Context, limit int) ([]*entity.WhaleTransaction, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("whale alert API key not configured")
	}

	url := fmt.Sprintf("https://api.whale-alert.io/v1/transactions?api_key=%s&min_value=%d&limit=%d",
		s.apiKey, int64(s.minAmount), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response struct {
		Transactions []struct {
			Blockchain      string  `json:"blockchain"`
			Symbol          string  `json:"symbol"`
			Amount          float64 `json:"amount"`
			AmountUSD       float64 `json:"amount_usd"`
			FromOwner       string  `json:"from_owner_type"`
			ToOwner         string  `json:"to_owner_type"`
			Hash            string  `json:"hash"`
			Timestamp       int64   `json:"timestamp"`
		} `json:"transactions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	transactions := make([]*entity.WhaleTransaction, 0, len(response.Transactions))
	for _, tx := range response.Transactions {
		transactions = append(transactions, &entity.WhaleTransaction{
			Blockchain: tx.Blockchain,
			Symbol:     tx.Symbol,
			Amount:     tx.Amount,
			AmountUSD:  tx.AmountUSD,
			From: entity.Address{
				Type: tx.FromOwner,
			},
			To: entity.Address{
				Type: tx.ToOwner,
			},
			Hash:      tx.Hash,
			Timestamp: time.Unix(tx.Timestamp, 0),
			FetchedAt: time.Now(),
		})
	}

	return transactions, nil
}

// Start begins monitoring whale transactions
func (s *WhaleService) Start(interval time.Duration, alertFunc func([]*entity.WhaleTransaction)) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("whale service already running")
	}
	s.running = true
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	effectiveInterval := interval
	if effectiveInterval < 5*time.Minute {
		effectiveInterval = 5 * time.Minute
	}

	go func() {
		ticker := time.NewTicker(effectiveInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				transactions, err := s.FetchTransactions(ctx, 10)
				cancel()

				if err == nil && len(transactions) > 0 {
					newTx := s.filterNewTransactions(transactions)
					if len(newTx) > 0 {
						alertFunc(newTx)
					}
				}

			case <-s.stopChan:
				return
			}
		}
	}()

	return nil
}

// Stop stops the monitoring
func (s *WhaleService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		close(s.stopChan)
		s.running = false
	}
}

// IsRunning returns whether the service is currently running
func (s *WhaleService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// SetMinAmount sets the minimum transaction amount to alert on
func (s *WhaleService) SetMinAmount(amount float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.minAmount = amount
}

// GetMinAmount returns the current minimum amount
func (s *WhaleService) GetMinAmount() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.minAmount
}

// IsConfigured checks if the service has an API key
func (s *WhaleService) IsConfigured() bool {
	return s.apiKey != ""
}
