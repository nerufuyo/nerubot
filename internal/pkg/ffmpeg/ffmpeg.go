package ffmpeg

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nerufuyo/nerubot/internal/pkg/logger"
)

// FFmpeg represents an FFmpeg wrapper
type FFmpeg struct {
	path   string
	logger *logger.Logger
}

// Options holds FFmpeg execution options
type Options struct {
	Input          string
	Output         string
	BeforeOptions  []string
	Options        []string
	Format         string
	Timeout        time.Duration
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
}

// New creates a new FFmpeg instance
func New() (*FFmpeg, error) {
	log := logger.New("ffmpeg")
	
	// Try to detect FFmpeg
	path, err := detect()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %w", err)
	}
	
	log.Info("FFmpeg detected", "path", path)
	
	return &FFmpeg{
		path:   path,
		logger: log,
	}, nil
}

// NewWithPath creates a new FFmpeg instance with a specific path
func NewWithPath(path string) (*FFmpeg, error) {
	log := logger.New("ffmpeg")
	
	// Verify the path exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("ffmpeg binary not found at %s: %w", path, err)
	}
	
	log.Info("Using custom FFmpeg path", "path", path)
	
	return &FFmpeg{
		path:   path,
		logger: log,
	}, nil
}

// detect attempts to find the FFmpeg binary
func detect() (string, error) {
	// Common paths to check
	commonPaths := []string{
		"ffmpeg",                           // In PATH
		"/usr/bin/ffmpeg",                  // Linux
		"/usr/local/bin/ffmpeg",            // Linux/macOS
		"/opt/homebrew/bin/ffmpeg",         // macOS Homebrew (Apple Silicon)
		"/usr/local/Cellar/ffmpeg",         // macOS Homebrew (Intel)
	}
	
	// On Windows, add .exe extension
	if runtime.GOOS == "windows" {
		for i, path := range commonPaths {
			if !strings.HasSuffix(path, ".exe") {
				commonPaths[i] = path + ".exe"
			}
		}
		commonPaths = append(commonPaths, "C:\\ffmpeg\\bin\\ffmpeg.exe")
	}
	
	// Try to find in PATH first
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}
	
	// Try common paths
	for _, path := range commonPaths {
		// For paths with wildcards or partial paths, try to find them
		if strings.Contains(path, "*") || !filepath.IsAbs(path) && path != "ffmpeg" {
			matches, err := filepath.Glob(path + "/*/bin/ffmpeg")
			if err == nil && len(matches) > 0 {
				return matches[0], nil
			}
			continue
		}
		
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	return "", fmt.Errorf("ffmpeg binary not found in common locations")
}

// Path returns the path to the FFmpeg binary
func (f *FFmpeg) Path() string {
	return f.path
}

// Version returns the FFmpeg version
func (f *FFmpeg) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, f.path, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get ffmpeg version: %w", err)
	}
	
	// Parse version from first line
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	
	return string(output), nil
}

// Run executes FFmpeg with the given options
func (f *FFmpeg) Run(ctx context.Context, opts *Options) error {
	if opts == nil {
		return fmt.Errorf("options cannot be nil")
	}
	
	// Build command arguments
	args := make([]string, 0)
	
	// Add before options (reconnect, etc.)
	if len(opts.BeforeOptions) > 0 {
		args = append(args, opts.BeforeOptions...)
	}
	
	// Add input
	if opts.Input != "" {
		args = append(args, "-i", opts.Input)
	}
	
	// Add format if specified
	if opts.Format != "" {
		args = append(args, "-f", opts.Format)
	}
	
	// Add options (filters, codecs, etc.)
	if len(opts.Options) > 0 {
		args = append(args, opts.Options...)
	}
	
	// Add output
	if opts.Output != "" {
		args = append(args, opts.Output)
	} else {
		// Default to pipe output
		args = append(args, "pipe:1")
	}
	
	f.logger.Debug("Executing FFmpeg", "args", strings.Join(args, " "))
	
	// Create command
	cmd := exec.CommandContext(ctx, f.path, args...)
	
	// Setup stdin/stdout/stderr
	if opts.Stdin != nil {
		cmd.Stdin = opts.Stdin
	}
	if opts.Stdout != nil {
		cmd.Stdout = opts.Stdout
	}
	if opts.Stderr != nil {
		cmd.Stderr = opts.Stderr
	}
	
	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	
	// Wait for completion or timeout
	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Wait()
	}()
	
	// Handle timeout if specified
	if opts.Timeout > 0 {
		select {
		case err := <-errChan:
			return err
		case <-time.After(opts.Timeout):
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			return fmt.Errorf("ffmpeg execution timeout after %v", opts.Timeout)
		case <-ctx.Done():
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			return ctx.Err()
		}
	}
	
	// Wait for completion
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return ctx.Err()
	}
}

// Convert converts audio/video with the given options
func (f *FFmpeg) Convert(ctx context.Context, input, output string, opts []string) error {
	return f.Run(ctx, &Options{
		Input:   input,
		Output:  output,
		Options: opts,
	})
}

// Stream streams audio to the given writer
func (f *FFmpeg) Stream(ctx context.Context, input string, output io.Writer, format string, opts []string) error {
	return f.Run(ctx, &Options{
		Input:  input,
		Format: format,
		Options: opts,
		Stdout: output,
	})
}

// ParseOptions parses a string of FFmpeg options into a slice
func ParseOptions(optString string) []string {
	if optString == "" {
		return nil
	}
	
	// Split by spaces but respect quotes
	var opts []string
	var current strings.Builder
	inQuote := false
	
	for _, r := range optString {
		switch r {
		case '"', '\'':
			inQuote = !inQuote
		case ' ':
			if inQuote {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				opts = append(opts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	
	if current.Len() > 0 {
		opts = append(opts, current.String())
	}
	
	return opts
}
