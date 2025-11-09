package ytdlp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// YtDlp represents a yt-dlp wrapper
type YtDlp struct {
	path   string
	logger *logger.Logger
}

// VideoInfo holds video metadata from yt-dlp
type VideoInfo struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Duration    int      `json:"duration"`
	URL         string   `json:"url"`
	Webpage     string   `json:"webpage_url"`
	Thumbnail   string   `json:"thumbnail"`
	Description string   `json:"description"`
	Uploader    string   `json:"uploader"`
	Channel     string   `json:"channel"`
	ViewCount   int64    `json:"view_count"`
	LikeCount   int64    `json:"like_count"`
	Formats     []Format `json:"formats"`
	ExtractorKey string  `json:"extractor_key"`
	IsLive      bool     `json:"is_live"`
}

// Format holds format information
type Format struct {
	FormatID   string  `json:"format_id"`
	URL        string  `json:"url"`
	Ext        string  `json:"ext"`
	Quality    float64 `json:"quality"`
	Filesize   int64   `json:"filesize"`
	Bitrate    float64 `json:"tbr"`
	AudioCodec string  `json:"acodec"`
	VideoCodec string  `json:"vcodec"`
}

// PlaylistInfo holds playlist metadata
type PlaylistInfo struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Entries []VideoInfo `json:"entries"`
}

// ExtractOptions holds options for extraction
type ExtractOptions struct {
	Format       string
	AudioOnly    bool
	NoPlaylist   bool
	PlaylistItems string
	Cookies      string
	UserAgent    string
	Timeout      time.Duration
}

// New creates a new YtDlp instance
func New() (*YtDlp, error) {
	log := logger.New("ytdlp")
	
	// Try to detect yt-dlp
	path, err := detect()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp not found: %w", err)
	}
	
	log.Info("yt-dlp detected", "path", path)
	
	return &YtDlp{
		path:   path,
		logger: log,
	}, nil
}

// NewWithPath creates a new YtDlp instance with a specific path
func NewWithPath(path string) (*YtDlp, error) {
	log := logger.New("ytdlp")
	
	// Verify the path exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("yt-dlp binary not found at %s: %w", path, err)
	}
	
	log.Info("Using custom yt-dlp path", "path", path)
	
	return &YtDlp{
		path:   path,
		logger: log,
	}, nil
}

// detect attempts to find the yt-dlp binary
func detect() (string, error) {
	// Common names
	names := []string{"yt-dlp", "youtube-dl"}
	
	// On Windows, add .exe extension
	if runtime.GOOS == "windows" {
		for i, name := range names {
			names[i] = name + ".exe"
		}
	}
	
	// Try to find in PATH
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	
	// Common paths to check
	commonPaths := []string{
		"/usr/bin/yt-dlp",
		"/usr/local/bin/yt-dlp",
		"/opt/homebrew/bin/yt-dlp",
		"/usr/bin/youtube-dl",
		"/usr/local/bin/youtube-dl",
	}
	
	if runtime.GOOS == "windows" {
		commonPaths = append(commonPaths, "C:\\yt-dlp\\yt-dlp.exe")
	}
	
	// Try common paths
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	return "", fmt.Errorf("yt-dlp binary not found")
}

// Path returns the path to the yt-dlp binary
func (y *YtDlp) Path() string {
	return y.path
}

// Version returns the yt-dlp version
func (y *YtDlp) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, y.path, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get yt-dlp version: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// ExtractInfo extracts video/playlist information
func (y *YtDlp) ExtractInfo(ctx context.Context, url string, opts *ExtractOptions) (*VideoInfo, error) {
	if opts == nil {
		opts = &ExtractOptions{}
	}
	
	// Build command arguments
	args := []string{
		"--dump-json",
		"--no-warnings",
		"--no-call-home",
		"--no-check-certificate",
	}
	
	// Add playlist handling
	if opts.NoPlaylist {
		args = append(args, "--no-playlist")
	}
	if opts.PlaylistItems != "" {
		args = append(args, "--playlist-items", opts.PlaylistItems)
	}
	
	// Add format
	if opts.Format != "" {
		args = append(args, "-f", opts.Format)
	} else if opts.AudioOnly {
		args = append(args, "-f", "bestaudio/best")
	}
	
	// Add cookies if provided
	if opts.Cookies != "" {
		args = append(args, "--cookies", opts.Cookies)
	}
	
	// Add user agent if provided
	if opts.UserAgent != "" {
		args = append(args, "--user-agent", opts.UserAgent)
	}
	
	// Add URL
	args = append(args, url)
	
	y.logger.Debug("Extracting info", "url", url)
	
	// Create command with timeout
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}
	
	cmd := exec.CommandContext(ctx, y.path, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp extraction failed: %w, output: %s", err, string(output))
	}
	
	// Parse JSON output
	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse yt-dlp output: %w", err)
	}
	
	return &info, nil
}

// ExtractPlaylist extracts playlist information
func (y *YtDlp) ExtractPlaylist(ctx context.Context, url string, opts *ExtractOptions) (*PlaylistInfo, error) {
	if opts == nil {
		opts = &ExtractOptions{}
	}
	
	// Build command arguments
	args := []string{
		"--dump-json",
		"--flat-playlist",
		"--no-warnings",
		"--no-call-home",
		"--no-check-certificate",
	}
	
	if opts.PlaylistItems != "" {
		args = append(args, "--playlist-items", opts.PlaylistItems)
	}
	
	args = append(args, url)
	
	y.logger.Debug("Extracting playlist", "url", url)
	
	// Create command
	cmd := exec.CommandContext(ctx, y.path, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp playlist extraction failed: %w", err)
	}
	
	// Parse JSON output (yt-dlp outputs one JSON per line for playlists)
	lines := strings.Split(string(output), "\n")
	var playlist PlaylistInfo
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		var entry VideoInfo
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		playlist.Entries = append(playlist.Entries, entry)
	}
	
	if len(playlist.Entries) > 0 {
		// Use first entry's info for playlist metadata
		playlist.ID = playlist.Entries[0].ID
		playlist.Title = "Playlist" // yt-dlp doesn't provide playlist title in flat mode
	}
	
	return &playlist, nil
}

// GetStreamURL gets the direct stream URL for a video
func (y *YtDlp) GetStreamURL(ctx context.Context, url string, opts *ExtractOptions) (string, error) {
	if opts == nil {
		opts = &ExtractOptions{
			AudioOnly: true,
		}
	}
	
	// Build command arguments
	args := []string{
		"--get-url",
		"--no-warnings",
		"--no-call-home",
		"--no-check-certificate",
	}
	
	if opts.NoPlaylist {
		args = append(args, "--no-playlist")
	}
	
	// Add format
	if opts.Format != "" {
		args = append(args, "-f", opts.Format)
	} else if opts.AudioOnly {
		args = append(args, "-f", "bestaudio/best")
	}
	
	args = append(args, url)
	
	y.logger.Debug("Getting stream URL", "url", url)
	
	// Create command
	cmd := exec.CommandContext(ctx, y.path, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get stream URL: %w", err)
	}
	
	// Return first line (the URL)
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 && lines[0] != "" {
		return strings.TrimSpace(lines[0]), nil
	}
	
	return "", fmt.Errorf("no stream URL found")
}

// Download downloads a video to a file
func (y *YtDlp) Download(ctx context.Context, url, output string, opts *ExtractOptions) error {
	if opts == nil {
		opts = &ExtractOptions{}
	}
	
	// Build command arguments
	args := []string{
		"-o", output,
		"--no-warnings",
		"--no-call-home",
		"--no-check-certificate",
	}
	
	if opts.NoPlaylist {
		args = append(args, "--no-playlist")
	}
	
	if opts.Format != "" {
		args = append(args, "-f", opts.Format)
	}
	
	args = append(args, url)
	
	y.logger.Info("Downloading", "url", url, "output", output)
	
	// Create output directory if needed
	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Create command
	cmd := exec.CommandContext(ctx, y.path, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	
	y.logger.Info("Download complete", "output", output)
	return nil
}

// Search searches for videos
func (y *YtDlp) Search(ctx context.Context, query string, maxResults int) ([]VideoInfo, error) {
	// Use YouTube search URL
	searchURL := fmt.Sprintf("ytsearch%d:%s", maxResults, query)
	
	// Build command arguments
	args := []string{
		"--dump-json",
		"--flat-playlist",
		"--no-warnings",
		"--no-call-home",
		"--no-check-certificate",
		searchURL,
	}
	
	y.logger.Debug("Searching", "query", query, "max_results", maxResults)
	
	// Create command
	cmd := exec.CommandContext(ctx, y.path, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	
	// Parse results (one JSON per line)
	var results []VideoInfo
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		var info VideoInfo
		if err := json.Unmarshal([]byte(line), &info); err != nil {
			y.logger.Warn("Failed to parse search result", "error", err)
			continue
		}
		results = append(results, info)
	}
	
	return results, nil
}

// IsPlaylist checks if a URL is a playlist
func IsPlaylist(url string) bool {
	playlistIndicators := []string{
		"playlist?list=",
		"&list=",
		"/playlists/",
		"/sets/", // SoundCloud
		"/album/", // Spotify
	}
	
	for _, indicator := range playlistIndicators {
		if strings.Contains(url, indicator) {
			return true
		}
	}
	
	return false
}

// GetSource determines the source from URL
func GetSource(url string) string {
	if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		return "youtube"
	}
	if strings.Contains(url, "spotify.com") {
		return "spotify"
	}
	if strings.Contains(url, "soundcloud.com") {
		return "soundcloud"
	}
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return "direct"
	}
	return "unknown"
}
