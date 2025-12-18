package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nerufuyo/nerubot/internal/config"
	"github.com/nerufuyo/nerubot/internal/entity"
)

// handlePlay handles the play command
func (b *Bot) handlePlay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a song name or URL")
		return
	}

	query := options[0].StringValue()

	// Check if user is in voice channel
	vs := b.findUserVoiceState(i.GuildID, i.Member.User.ID)
	if vs == nil {
		b.respondError(s, i, "You must be in a voice channel!")
		return
	}

	// Defer response since search might take time
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
		return
	}

	// Search and add song
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	song, err := b.musicService.AddSong(ctx, i.GuildID, vs.ChannelID, i.ChannelID, query, i.Member.User.ID)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("Failed to add song: %s", err.Error()))
		return
	}

	// Get queue to check position
	queue := b.musicService.GetQueue(i.GuildID, vs.ChannelID, i.ChannelID)

	var message string
	if queue.Size() == 1 {
		// First song - try to join voice and start playback via Lavalink
		b.logger.Info("Music playback initiated", 
			"guild", i.GuildID, "channel", vs.ChannelID)
		
		if err := b.musicService.Play(i.GuildID); err != nil {
			b.followUpError(s, i, fmt.Sprintf("Failed to start playback: %s", err.Error()))
			return
		}

		// If Lavalink is enabled, join voice and play
		if b.config.Lavalink.Enabled && b.lavalinkClient != nil {
			go b.playWithLavalink(i.GuildID, vs.ChannelID, song)
		}
		
		message = fmt.Sprintf("🎵 Now playing: **%s** by **%s**", song.Title, song.Artist)
	} else {
		message = fmt.Sprintf("➕ Added to queue: **%s** by **%s** (Position: #%d)", 
			song.Title, song.Artist, queue.Size())
	}

	b.followUp(s, i, message)
}

// handleSkip handles the skip command
func (b *Bot) handleSkip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	song, err := b.musicService.Skip(i.GuildID)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	if song != nil {
		b.respond(s, i, fmt.Sprintf("⏭️ Skipped to: **%s**", song.Title))
	} else {
		b.respond(s, i, "⏭️ Skipped. Queue is now empty.")
	}
}

// handleStop handles the stop command
func (b *Bot) handleStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := b.musicService.Stop(i.GuildID); err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	b.respond(s, i, "⏹️ Stopped playback and cleared queue")
}

// handleQueue handles the queue command
func (b *Bot) handleQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	queue := b.musicService.GetQueue(i.GuildID, "", "")
	
	if queue.IsEmpty() {
		b.respond(s, i, "📭 Queue is empty")
		return
	}

	// Build embed
	embed := &discordgo.MessageEmbed{
		Title: "🎵 Music Queue",
		Color: config.ColorMusic,
		Fields: []*discordgo.MessageEmbedField{},
	}

	// Add current song
	if current := queue.Current(); current != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Now Playing",
			Value:  fmt.Sprintf("**%s** by **%s**\n%s", current.Title, current.Artist, entity.FormatDuration(current.Duration)),
			Inline: false,
		})
	}

	// Add queue
	if queue.Remaining() > 0 {
		queueText := ""
		for i := queue.CurrentIndex + 1; i < queue.Size() && i < queue.CurrentIndex+6; i++ {
			song := queue.Songs[i]
			queueText += fmt.Sprintf("%d. **%s** - %s\n", i-queue.CurrentIndex, song.Title, entity.FormatDuration(song.Duration))
		}
		
		if queue.Remaining() > 5 {
			queueText += fmt.Sprintf("\n...and %d more songs", queue.Remaining()-5)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Up Next",
			Value:  queueText,
			Inline: false,
		})
	}

	// Add footer
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Total: %d songs | Duration: %s | Loop: %s",
			queue.Size(),
			entity.FormatDuration(queue.TotalDuration()),
			queue.LoopMode,
		),
	}

	b.respondEmbed(s, i, embed)
}

// handleConfess handles the confess command
func (b *Bot) handleConfess(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide confession content")
		return
	}

	content := options[0].StringValue()

	ctx := context.Background()
	confession, err := b.confessionService.SubmitConfession(ctx, i.GuildID, i.Member.User.ID, content)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	// Respond privately
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Your confession has been submitted anonymously! (ID: #%d)", confession.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		b.logger.Error("Failed to respond", "error", err)
	}
}

// handleRoast handles the roast command
func (b *Bot) handleRoast(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Determine target user
	targetUser := i.Member.User
	options := i.ApplicationCommandData().Options
	if len(options) > 0 && options[0].Type == discordgo.ApplicationCommandOptionUser {
		targetUser = options[0].UserValue(s)
	}

	if targetUser.Bot {
		b.respondError(s, i, "I don't roast bots, they're already broken enough!")
		return
	}

	ctx := context.Background()
	roast, err := b.roastService.GenerateRoast(ctx, targetUser.ID, i.GuildID, targetUser.Username)
	if err != nil {
		b.respondError(s, i, err.Error())
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🔥 Roast",
		Description: roast,
		Color:       config.ColorError,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Requested by %s", i.Member.User.Username),
		},
	}

	b.respondEmbed(s, i, embed)
}

// handleHelp handles the help command
func (b *Bot) handleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "📚 " + b.config.Bot.Name + " Help",
		Description: b.config.Bot.Description,
		Color:       config.ColorPrimary,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "🎵 Music Commands",
				Value: "`/play <query>` - Play a song or add to queue\n" +
					"`/skip` - Skip current song\n" +
					"`/stop` - Stop playback and clear queue\n" +
					"`/queue` - Show music queue",
				Inline: false,
			},
			{
				Name: "📝 Confession Commands",
				Value: "`/confess <content>` - Submit an anonymous confession",
				Inline: false,
			},
			{
				Name: "🔥 Roast Commands",
				Value: "`/roast [user]` - Get roasted based on Discord activity",
				Inline: false,
			},
			{
				Name: "🤖 AI Chatbot Commands",
				Value: "`/chat <message>` - Chat with AI (Claude, Gemini, or OpenAI)\n" +
					"`/chat-reset` - Reset your chat history",
				Inline: false,
			},
			{
				Name: "📰 News Commands",
				Value: "`/news` - Get latest news from multiple sources",
				Inline: false,
			},
			{
				Name: "🐋 Whale Alert Commands",
				Value: "`/whale` - Get recent whale cryptocurrency transactions",
				Inline: false,
			},
			{
				Name: "ℹ️ Other Commands",
				Value: "`/help` - Show this help message",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s v%s | %s", b.config.Bot.Name, b.config.Bot.Version, b.config.Bot.Author),
		},
	}

	b.respondEmbed(s, i, embed)
}

// Helper methods

func (b *Bot) findUserVoiceState(guildID, userID string) *discordgo.VoiceState {
	guild, err := b.session.State.Guild(guildID)
	if err != nil {
		return nil
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs
		}
	}

	return nil
}

func (b *Bot) followUp(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

func (b *Bot) followUpError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	b.followUp(s, i, "❌ "+message)
}

// handleChat handles the chatbot command
func (b *Bot) handleChat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.chatbotService == nil {
		b.respondError(s, i, "Chatbot service is not available")
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondError(s, i, "Please provide a message")
		return
	}

	message := options[0].StringValue()

	// Defer response since AI might take time
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
		return
	}

	// Get AI response
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := b.chatbotService.Chat(ctx, i.Member.User.ID, message)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("AI error: %s", err.Error()))
		return
	}

	// Send response
	embed := &discordgo.MessageEmbed{
		Title:       "🤖 AI Response",
		Description: response,
		Color:       0x00ff00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Providers: %v", b.chatbotService.GetAvailableProviders()),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

// handleChatReset handles clearing chat history
func (b *Bot) handleChatReset(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.chatbotService == nil {
		b.respondError(s, i, "Chatbot service is not available")
		return
	}

	b.chatbotService.ResetSession(i.Member.User.ID)
	b.respond(s, i, "✅ Chat history cleared!")
}

// handleNews handles fetching latest news
func (b *Bot) handleNews(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.newsService == nil {
		b.respondError(s, i, "News service is not available")
		return
	}

	// Defer response since fetching might take time
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
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
		Title:       "📰 Latest News",
		Color:       0x0099ff,
		Timestamp:   time.Now().Format(time.RFC3339),
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	for _, article := range articles {
		fieldValue := article.Description
		if len(fieldValue) > 200 {
			fieldValue = fieldValue[:197] + "..."
		}
		fieldValue += fmt.Sprintf("\n[Read more](%s)", article.URL)

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("📌 %s", article.Title),
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
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
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
		Title:     "🐋 Whale Transactions",
		Color:     0xffd700,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make([]*discordgo.MessageEmbedField, 0),
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
			Name:   fmt.Sprintf("🔔 %s Transaction", tx.Symbol),
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

// handleStats handles showing server statistics
func (b *Bot) handleStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.analyticsService == nil {
		b.respondError(s, i, "Analytics service is not available")
		return
	}

	// Defer response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
		return
	}

	// Get server stats
	stats, err := b.analyticsService.GetServerStats(i.GuildID)
	if err != nil {
		b.followUpError(s, i, "No statistics available yet. Start using commands to generate stats!")
		return
	}

	// Get most used command and most active user
	topCmd, topCmdCount := stats.GetMostUsedCommand()
	topUser, topUserCount := stats.GetMostActiveUser()

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title:     "📊 Server Statistics",
		Color:     0x00ff00,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "📈 Total Commands",
				Value:  fmt.Sprintf("%d", stats.CommandsUsed),
				Inline: true,
			},
			{
				Name:   "🎵 Songs Played",
				Value:  fmt.Sprintf("%d", stats.SongsPlayed),
				Inline: true,
			},
			{
				Name:   "📝 Confessions",
				Value:  fmt.Sprintf("%d", stats.ConfessionsTotal),
				Inline: true,
			},
			{
				Name:   "🔥 Roasts Generated",
				Value:  fmt.Sprintf("%d", stats.RoastsGenerated),
				Inline: true,
			},
			{
				Name:   "💬 Chat Messages",
				Value:  fmt.Sprintf("%d", stats.ChatMessages),
				Inline: true,
			},
			{
				Name:   "📰 News Requests",
				Value:  fmt.Sprintf("%d", stats.NewsRequests),
				Inline: true,
			},
		},
	}

	if topCmd != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "⭐ Most Used Command",
			Value:  fmt.Sprintf("/%s (%d times)", topCmd, topCmdCount),
			Inline: false,
		})
	}

	if topUser != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🏆 Most Active User",
			Value:  fmt.Sprintf("<@%s> (%d commands)", topUser, topUserCount),
			Inline: false,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Server: %s | Active since %s", stats.GuildName, stats.FirstSeen.Format("Jan 2, 2006")),
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

// handleProfile handles showing user statistics
func (b *Bot) handleProfile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if b.analyticsService == nil {
		b.respondError(s, i, "Analytics service is not available")
		return
	}

	// Get target user (default to command user)
	targetUser := i.Member.User
	options := i.ApplicationCommandData().Options
	if len(options) > 0 && options[0].Type == discordgo.ApplicationCommandOptionUser {
		targetUser = options[0].UserValue(s)
	}

	// Defer response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		b.logger.Error("Failed to defer response", "error", err)
		return
	}

	// Get user stats
	stats, err := b.analyticsService.GetUserStats(targetUser.ID)
	if err != nil {
		b.followUpError(s, i, fmt.Sprintf("No statistics available for <@%s> yet.", targetUser.ID))
		return
	}

	// Get favorite command
	favCmd, favCount := stats.GetFavoriteCommand()

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("📊 User Profile: %s", targetUser.Username),
		Color:     0x0099ff,
		Timestamp: time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: targetUser.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "📈 Total Commands",
				Value:  fmt.Sprintf("%d", stats.CommandsUsed),
				Inline: true,
			},
			{
				Name:   "🎵 Songs Requested",
				Value:  fmt.Sprintf("%d", stats.SongsRequested),
				Inline: true,
			},
			{
				Name:   "📝 Confessions Posted",
				Value:  fmt.Sprintf("%d", stats.ConfessionsPosted),
				Inline: true,
			},
			{
				Name:   "🔥 Roasts Received",
				Value:  fmt.Sprintf("%d", stats.RoastsReceived),
				Inline: true,
			},
			{
				Name:   "💬 Chat Messages",
				Value:  fmt.Sprintf("%d", stats.ChatMessages),
				Inline: true,
			},
			{
				Name:   "📰 News Checks",
				Value:  fmt.Sprintf("%d", stats.NewsRequests),
				Inline: true,
			},
		},
	}

	if favCmd != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "⭐ Favorite Command",
			Value:  fmt.Sprintf("/%s (%d times)", favCmd, favCount),
			Inline: false,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Active since %s", stats.FirstSeen.Format("Jan 2, 2006")),
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		b.logger.Error("Failed to send follow-up", "error", err)
	}
}

// playWithLavalink handles voice joining and playback via Lavalink
func (b *Bot) playWithLavalink(guildID, channelID string, song *entity.Song) {
	b.logger.Info("Attempting to play with Lavalink",
		"guild", guildID, "channel", channelID, "song", song.Title)

	// Get session ID for Lavalink
	sessionID := b.session.State.User.ID

	// Search for track on Lavalink
	query := fmt.Sprintf("ytsearch:%s %s", song.Title, song.Artist)
	tracks, err := b.lavalinkClient.SearchTracks(query)
	if err != nil {
		b.logger.Error("Failed to search tracks on Lavalink", "error", err)
		return
	}

	if len(tracks) == 0 {
		b.logger.Warn("No tracks found on Lavalink", "query", query)
		return
	}

	// Get the first track
	track := tracks[0]

	// Tell Lavalink to join voice and play the track
	// Note: This does NOT use discordgo's ChannelVoiceJoin which has encryption issues
	// Lavalink handles all voice protocol internally
	if err := b.lavalinkClient.JoinVoice(guildID, sessionID, channelID, track.Encoded); err != nil {
		b.logger.Error("Failed to join voice channel with Lavalink", "error", err)
		return
	}

	b.logger.Info("Track playing on Lavalink", "guild", guildID, "track", track.Info.Title)

	// Auto-disconnect after track duration + buffer
	duration := time.Duration(track.Info.Length) * time.Millisecond
	go func() {
		time.Sleep(duration + 3*time.Second)
		if err := b.lavalinkClient.Stop(guildID, sessionID); err != nil {
			b.logger.Warn("Failed to stop playback", "error", err)
		}
		b.logger.Info("Auto-disconnected from voice", "guild", guildID)
	}()
}
