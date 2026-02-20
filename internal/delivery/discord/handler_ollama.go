package discord

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/pkg/ai"
)

// handleOllamaModels handles the /ollama-models command — lists available models.
func (b *Bot) handleOllamaModels(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.ollamaClient == nil {
		b.respondError(s, i, "Ollama is not configured. Set `OLLAMA_URL` env var.")
		return
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	models, err := b.ollamaClient.ListModels(ctx)
	if err != nil {
		b.followUpError(s, i, "Failed to fetch models: "+err.Error())
		return
	}

	if len(models) == 0 {
		b.followUp(s, i, "No models found on the Ollama server.")
		return
	}

	var fields []*discordgo.MessageEmbedField
	for _, m := range models {
		value := fmt.Sprintf(
			"**Family:** %s\n**Parameters:** %s\n**Quantization:** %s\n**Size:** %s\n**Format:** %s",
			m.Details.Family,
			m.Details.ParameterSize,
			m.Details.QuantizationLevel,
			ai.FormatSize(m.Size),
			m.Details.Format,
		)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("🤖 %s", m.Name),
			Value:  value,
			Inline: true,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🦙 Ollama Models",
		Description: fmt.Sprintf("Found **%d** model(s) on the Ollama server.", len(models)),
		Color:       0x00B4D8,
		Fields:      fields,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ollama • " + b.config.AI.OllamaURL,
		},
	}

	b.followUpEmbed(s, i, embed)
}

// handleOllamaBench handles the /ollama-bench command — benchmarks a model.
func (b *Bot) handleOllamaBench(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.ollamaClient == nil {
		b.respondError(s, i, "Ollama is not configured. Set `OLLAMA_URL` env var.")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please specify a model name.")
		return
	}

	model := options[0].StringValue()
	prompt := "Explain what a Large Language Model is in 3 sentences."
	for _, opt := range options {
		if opt.Name == "prompt" && opt.StringValue() != "" {
			prompt = opt.StringValue()
		}
	}

	if err := b.deferResponse(s, i); err != nil {
		return
	}

	// Use a generous timeout — small models are fast, but large ones can be slow
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	result, err := b.ollamaClient.Benchmark(ctx, model, prompt)
	if err != nil {
		b.followUpError(s, i, "Benchmark failed: "+err.Error())
		return
	}

	// Build performance metrics
	metrics := strings.Builder{}
	metrics.WriteString(fmt.Sprintf("⏱️ **Total Time:** %s\n", result.TotalDuration.Round(time.Millisecond)))
	metrics.WriteString(fmt.Sprintf("🚀 **Time to First Token:** %s\n", result.FirstTokenTime.Round(time.Millisecond)))
	metrics.WriteString(fmt.Sprintf("⚡ **Tokens/sec:** %.1f\n", result.TokensPerSecond))
	metrics.WriteString(fmt.Sprintf("📊 **Tokens Generated:** %d\n", result.EvalCount))

	if result.PromptEvalCount > 0 {
		metrics.WriteString(fmt.Sprintf("📝 **Prompt Tokens Evaluated:** %d\n", result.PromptEvalCount))
	}

	// Ollama server-side durations
	serverMetrics := strings.Builder{}
	if result.LoadDuration > 0 {
		serverMetrics.WriteString(fmt.Sprintf("📦 **Model Load:** %s\n", result.LoadDuration.Round(time.Millisecond)))
	}
	if result.PromptEvalDur > 0 {
		serverMetrics.WriteString(fmt.Sprintf("📋 **Prompt Eval:** %s\n", result.PromptEvalDur.Round(time.Millisecond)))
	}
	if result.EvalDuration > 0 {
		serverMetrics.WriteString(fmt.Sprintf("🔄 **Generation:** %s\n", result.EvalDuration.Round(time.Millisecond)))
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "📊 Performance",
			Value: metrics.String(),
		},
	}

	if serverMetrics.Len() > 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "🖥️ Server Timing",
			Value: serverMetrics.String(),
		})
	}

	// Truncate prompt for display
	displayPrompt := prompt
	if len(displayPrompt) > 100 {
		displayPrompt = displayPrompt[:100] + "..."
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "💬 Prompt",
		Value: fmt.Sprintf("```%s```", displayPrompt),
	})

	if result.Response != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "🤖 Response (preview)",
			Value: fmt.Sprintf("```%s```", result.Response),
		})
	}

	// Color based on speed
	color := 0xFF6B6B // red = slow
	if result.TokensPerSecond >= 30 {
		color = 0x00FF00 // green = fast
	} else if result.TokensPerSecond >= 15 {
		color = 0xFFD700 // yellow = ok
	} else if result.TokensPerSecond >= 5 {
		color = 0xFFA500 // orange = meh
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🦙 Ollama Benchmark — %s", model),
		Description: "Performance benchmark results for the selected model.",
		Color:       color,
		Fields:      fields,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ollama Benchmark • " + b.config.AI.OllamaURL,
		},
	}

	b.followUpEmbed(s, i, embed)
}
