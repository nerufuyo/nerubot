package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
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

	// Extract language option
	lang := config.DefaultLang
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "lang" {
			lang = opt.StringValue()
		}
	}

	// Defer response since fetching might take time
	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Fetch news by language
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	articles, err := b.newsService.FetchLatestByLang(ctx, 5, lang)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to fetch news: %s", err.Error()))
		return
	}

	if len(articles) == 0 {
		b.followUp(s, i, "No news articles found")
		return
	}

	// Create embed
	langLabel := config.LanguageNames[lang]
	if langLabel == "" {
		langLabel = "English"
	}

	embed := &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("Latest News (%s)", langLabel),
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
