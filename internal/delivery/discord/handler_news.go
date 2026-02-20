package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// handleNews handles fetching latest news
func (b *Bot) handleNews(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if news is enabled from dashboard
	if b.backendClient != nil && !b.backendClient.GetSettings().Features.NewsEnabled {
		b.respondError(s, i, "News feature is currently disabled by the admin.")
		return
	}

	if b.newsService == nil {
		b.respondError(s, i, "News service is not available")
		return
	}

	// Defer response since fetching might take time
	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Fetch news
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	articles, err := b.newsService.FetchLatest(ctx, 5)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to fetch news: %s", err.Error()))
		return
	}

	if len(articles) == 0 {
		b.followUp(s, i, "No news articles found")
		return
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title:     "Latest News",
		Color:     0x0099ff,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make([]*discordgo.MessageEmbedField, 0, len(articles)),
	}

	for _, article := range articles {
		fieldValue := article.Description
		if len(fieldValue) > 200 {
			fieldValue = fieldValue[:197] + "..."
		}
		fieldValue += fmt.Sprintf("\n[Read more](%s)", article.URL)

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   article.Title,
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
