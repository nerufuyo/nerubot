package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nerufuyo/nerubot/internal/entity"
	"github.com/nerufuyo/nerubot/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MusicRepository handles MongoDB operations for music playlists.
type MusicRepository struct {
	logger *logger.Logger
}

// NewMusicRepository creates a new MusicRepository.
func NewMusicRepository() *MusicRepository {
	return &MusicRepository{
		logger: logger.New("music-repo"),
	}
}

func (r *MusicRepository) playlistCol() *mongo.Collection {
	return MongoDB.Collection("playlists")
}

func (r *MusicRepository) settingsCol() *mongo.Collection {
	return MongoDB.Collection("guild_music_settings")
}

// SavePlaylist creates or updates a playlist.
func (r *MusicRepository) SavePlaylist(ctx context.Context, p *entity.Playlist) error {
	p.UpdatedAt = time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	filter := bson.M{"_id": p.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.playlistCol().ReplaceOne(ctx, filter, p, opts)
	if err != nil {
		return fmt.Errorf("failed to save playlist: %w", err)
	}
	return nil
}

// GetPlaylist retrieves a single playlist by ID.
func (r *MusicRepository) GetPlaylist(ctx context.Context, id string) (*entity.Playlist, error) {
	var p entity.Playlist
	err := r.playlistCol().FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetUserPlaylists retrieves all playlists for a user.
func (r *MusicRepository) GetUserPlaylists(ctx context.Context, userID string) ([]*entity.Playlist, error) {
	cursor, err := r.playlistCol().Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to query playlists: %w", err)
	}
	defer cursor.Close(ctx)

	var playlists []*entity.Playlist
	if err := cursor.All(ctx, &playlists); err != nil {
		return nil, fmt.Errorf("failed to decode playlists: %w", err)
	}
	return playlists, nil
}

// DeletePlaylist removes a playlist by ID.
func (r *MusicRepository) DeletePlaylist(ctx context.Context, id string) error {
	_, err := r.playlistCol().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}
	return nil
}

// SaveGuildMusicSettings saves per-guild music settings.
func (r *MusicRepository) SaveGuildMusicSettings(ctx context.Context, s *entity.GuildMusicSettings) error {
	filter := bson.M{"guildId": s.GuildID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.settingsCol().ReplaceOne(ctx, filter, s, opts)
	if err != nil {
		return fmt.Errorf("failed to save guild music settings: %w", err)
	}
	return nil
}

// GetGuildMusicSettings retrieves music settings for a guild.
func (r *MusicRepository) GetGuildMusicSettings(ctx context.Context, guildID string) (*entity.GuildMusicSettings, error) {
	var s entity.GuildMusicSettings
	err := r.settingsCol().FindOne(ctx, bson.M{"guildId": guildID}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
