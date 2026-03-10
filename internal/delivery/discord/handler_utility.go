package discord

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
)

// handleCoinflip handles the /coinflip command.
func (b *Bot) handleCoinflip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	result := "Heads 🪙"
	if rand.Intn(2) == 0 {
		result = "Tails 🪙"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🪙 Coin Flip",
		Description: fmt.Sprintf("**%s**", result),
		Color:       config.ColorPrimary,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// 8ball responses
var eightBallResponses = []string{
	// Positive
	"It is certain ✅", "It is decidedly so ✅", "Without a doubt ✅",
	"Yes, definitely ✅", "You may rely on it ✅", "As I see it, yes ✅",
	"Most likely ✅", "Outlook good ✅", "Yes ✅", "Signs point to yes ✅",
	// Neutral
	"Reply hazy, try again 🤔", "Ask again later 🤔", "Better not tell you now 🤔",
	"Cannot predict now 🤔", "Concentrate and ask again 🤔",
	// Negative
	"Don't count on it ❌", "My reply is no ❌", "My sources say no ❌",
	"Outlook not so good ❌", "Very doubtful ❌",
}

// handleEightBall handles the /8ball command.
func (b *Bot) handleEightBall(s *discordgo.Session, i *discordgo.InteractionCreate) {
	question := ""
	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Name == "question" {
			question = opt.StringValue()
		}
	}

	response := eightBallResponses[rand.Intn(len(eightBallResponses))]

	embed := &discordgo.MessageEmbed{
		Title: "🎱 Magic 8-Ball",
		Color: config.ColorPrimary,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Question", Value: question, Inline: false},
			{Name: "Answer", Value: response, Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}

// handlePoll handles the /poll command - creates a simple poll.
func (b *Bot) handlePoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	question := ""
	var choices []string

	for _, opt := range options {
		switch opt.Name {
		case "question":
			question = opt.StringValue()
		case "option1":
			choices = append(choices, opt.StringValue())
		case "option2":
			choices = append(choices, opt.StringValue())
		case "option3":
			if opt.StringValue() != "" {
				choices = append(choices, opt.StringValue())
			}
		case "option4":
			if opt.StringValue() != "" {
				choices = append(choices, opt.StringValue())
			}
		case "option5":
			if opt.StringValue() != "" {
				choices = append(choices, opt.StringValue())
			}
		}
	}

	if len(choices) < 2 {
		b.respondError(s, i, "Please provide at least 2 options.")
		return
	}

	// Build poll description
	pollEmojis := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}
	desc := ""
	for idx, choice := range choices {
		desc += fmt.Sprintf("%s %s\n", pollEmojis[idx], choice)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📊 " + question,
		Description: desc,
		Color:       config.ColorPrimary,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Poll by %s • React to vote!", i.Member.User.Username),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Send response and add reactions
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond to poll", "error", err)
		return
	}

	// Get the response message to add reactions
	resp, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		b.logger.Error("Failed to get poll response", "error", err)
		return
	}

	for idx := range choices {
		_ = s.MessageReactionAdd(resp.ChannelID, resp.ID, pollEmojis[idx])
	}
}

// handleCalc handles the /calc command - simple math calculations.
func (b *Bot) handleCalc(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var a, bVal float64
	op := ""

	for _, opt := range options {
		switch opt.Name {
		case "a":
			a = opt.FloatValue()
		case "b":
			bVal = opt.FloatValue()
		case "operation":
			op = opt.StringValue()
		}
	}

	var result float64
	var resultStr string

	switch op {
	case "add":
		result = a + bVal
		resultStr = fmt.Sprintf("%.4g + %.4g = **%.4g**", a, bVal, result)
	case "subtract":
		result = a - bVal
		resultStr = fmt.Sprintf("%.4g - %.4g = **%.4g**", a, bVal, result)
	case "multiply":
		result = a * bVal
		resultStr = fmt.Sprintf("%.4g × %.4g = **%.4g**", a, bVal, result)
	case "divide":
		if bVal == 0 {
			b.respondError(s, i, "Cannot divide by zero!")
			return
		}
		result = a / bVal
		resultStr = fmt.Sprintf("%.4g ÷ %.4g = **%.4g**", a, bVal, result)
	case "power":
		result = math.Pow(a, bVal)
		resultStr = fmt.Sprintf("%.4g ^ %.4g = **%.4g**", a, bVal, result)
	case "modulo":
		if bVal == 0 {
			b.respondError(s, i, "Cannot modulo by zero!")
			return
		}
		result = math.Mod(a, bVal)
		resultStr = fmt.Sprintf("%.4g %% %.4g = **%.4g**", a, bVal, result)
	default:
		b.respondError(s, i, "Unknown operation")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🧮 Calculator",
		Description: resultStr,
		Color:       config.ColorPrimary,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	b.respondEmbed(s, i, embed)
}
