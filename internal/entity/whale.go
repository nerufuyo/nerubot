package entity

import (
	"fmt"
	"time"
)

// WhaleTransaction represents a large cryptocurrency transaction
type WhaleTransaction struct {
	ID            string    `json:"id"`
	Hash          string    `json:"hash"`
	Blockchain    string    `json:"blockchain"`
	Symbol        string    `json:"symbol"`
	Amount        float64   `json:"amount"`
	AmountUSD     float64   `json:"amount_usd"`
	From          Address   `json:"from"`
	To            Address   `json:"to"`
	Timestamp     time.Time `json:"timestamp"`
	TransactionFee float64  `json:"transaction_fee"`
	BlockNumber   int64     `json:"block_number"`
	FetchedAt     time.Time `json:"fetched_at"`
}

// Address represents a blockchain address
type Address struct {
	Address  string `json:"address"`
	Owner    string `json:"owner,omitempty"`
	Label    string `json:"label,omitempty"`
	Type     string `json:"type,omitempty"` // exchange, wallet, contract, etc.
}

// GuildWhaleSettings holds whale alert settings for a guild
type GuildWhaleSettings struct {
	GuildID        string   `json:"guild_id"`
	Enabled        bool     `json:"enabled"`
	ChannelID      string   `json:"channel_id"`
	MinAmount      float64  `json:"min_amount"` // Minimum USD amount to alert
	Blockchains    []string `json:"blockchains"`
	Symbols        []string `json:"symbols"` // Specific crypto symbols to watch
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CryptoGuru represents a crypto influencer/guru
type CryptoGuru struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	TwitterHandle string  `json:"twitter_handle"`
	Verified    bool      `json:"verified"`
	Followers   int       `json:"followers"`
	Category    string    `json:"category"` // trader, analyst, developer, etc.
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GuruTweet represents a tweet from a crypto guru
type GuruTweet struct {
	ID          string    `json:"id"`
	GuruID      string    `json:"guru_id"`
	GuruName    string    `json:"guru_name"`
	Content     string    `json:"content"`
	URL         string    `json:"url"`
	Likes       int       `json:"likes"`
	Retweets    int       `json:"retweets"`
	CreatedAt   time.Time `json:"created_at"`
	FetchedAt   time.Time `json:"fetched_at"`
}

// NewWhaleTransaction creates a new WhaleTransaction instance
func NewWhaleTransaction(blockchain, symbol string, amount, amountUSD float64) *WhaleTransaction {
	return &WhaleTransaction{
		Blockchain: blockchain,
		Symbol:     symbol,
		Amount:     amount,
		AmountUSD:  amountUSD,
		Timestamp:  time.Now(),
		FetchedAt:  time.Now(),
	}
}

// NewGuildWhaleSettings creates new settings with defaults
func NewGuildWhaleSettings(guildID, channelID string) *GuildWhaleSettings {
	return &GuildWhaleSettings{
		GuildID:     guildID,
		Enabled:     true,
		ChannelID:   channelID,
		MinAmount:   1000000, // $1M minimum
		Blockchains: []string{"bitcoin", "ethereum"},
		Symbols:     []string{"BTC", "ETH"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ShouldAlert checks if a transaction should trigger an alert
func (s *GuildWhaleSettings) ShouldAlert(tx *WhaleTransaction) bool {
	if !s.Enabled {
		return false
	}
	
	// Check minimum amount
	if tx.AmountUSD < s.MinAmount {
		return false
	}
	
	// Check blockchain filter
	if len(s.Blockchains) > 0 {
		found := false
		for _, chain := range s.Blockchains {
			if chain == tx.Blockchain {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check symbol filter
	if len(s.Symbols) > 0 {
		found := false
		for _, symbol := range s.Symbols {
			if symbol == tx.Symbol {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	return true
}

// FormatAmount formats the amount with proper units
func (t *WhaleTransaction) FormatAmount() string {
	if t.Amount >= 1000000 {
		return fmt.Sprintf("%.2fM %s", t.Amount/1000000, t.Symbol)
	} else if t.Amount >= 1000 {
		return fmt.Sprintf("%.2fK %s", t.Amount/1000, t.Symbol)
	}
	return fmt.Sprintf("%.2f %s", t.Amount, t.Symbol)
}

// FormatUSD formats the USD amount
func (t *WhaleTransaction) FormatUSD() string {
	if t.AmountUSD >= 1000000000 {
		return fmt.Sprintf("$%.2fB", t.AmountUSD/1000000000)
	} else if t.AmountUSD >= 1000000 {
		return fmt.Sprintf("$%.2fM", t.AmountUSD/1000000)
	} else if t.AmountUSD >= 1000 {
		return fmt.Sprintf("$%.2fK", t.AmountUSD/1000)
	}
	return fmt.Sprintf("$%.2f", t.AmountUSD)
}
