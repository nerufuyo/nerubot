package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// handleWhale handles fetching whale transactions
func (b *Bot) handleWhale(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.whaleService == nil {
		b.respondError(s, i, "Whale alert service is not available")
		return
	}

	if !b.whaleService.IsConfigured() {
		b.respondError(s, i, "Whale alert API key not configured")
		return
	}

	// Defer response since fetching might take time
	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Fetch transactions
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	transactions, err := b.whaleService.FetchTransactions(ctx, 5)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to fetch transactions: %s", err.Error()))
		return
	}

	if len(transactions) == 0 {
		b.followUp(s, i, "No whale transactions found")
		return
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title:     "Whale Transactions",
		Color:     0xffd700,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make([]*discordgo.MessageEmbedField, 0, len(transactions)),
	}

	for _, tx := range transactions {
		fieldValue := fmt.Sprintf("**Amount:** %s\n**USD Value:** %s\n**Blockchain:** %s\n**From:** %s → **To:** %s",
			tx.FormatAmount(),
			tx.FormatUSD(),
			tx.Blockchain,
			tx.From.Type,
			tx.To.Type,
		)

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s Transaction", tx.Symbol),
			Value:  fieldValue,
			Inline: false,
		})
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}
