package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
	name string
}

// LogLevel represents the log level
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// Config holds logger configuration
type Config struct {
	Level          LogLevel
	OutputFile     string
	MaxSizeMB      int
	MaxBackups     int
	MaxAgeDays     int
	Compress       bool
	EnableConsole  bool
	EnableFile     bool
	JSONFormat     bool
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:         LevelInfo,
		OutputFile:    "logs/bot.log",
		MaxSizeMB:     10,
		MaxBackups:    5,
		MaxAgeDays:    30,
		Compress:      true,
		EnableConsole: true,
		EnableFile:    true,
		JSONFormat:    false,
	}
}

var defaultLogger *Logger

// Init initializes the global logger
func Init(cfg *Config) (*Logger, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Create log directory if needed
	if cfg.EnableFile && cfg.OutputFile != "" {
		dir := filepath.Dir(cfg.OutputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	// Create multi-writer for console and file
	var writers []io.Writer

	if cfg.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	if cfg.EnableFile && cfg.OutputFile != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.OutputFile,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
		writers = append(writers, fileWriter)
	}

	multiWriter := io.MultiWriter(writers...)

	// Create handler options
	handlerOpts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
	}

	// Create handler based on format
	var handler slog.Handler
	if cfg.JSONFormat {
		handler = slog.NewJSONHandler(multiWriter, handlerOpts)
	} else {
		handler = slog.NewTextHandler(multiWriter, handlerOpts)
	}

	// Create logger
	logger := &Logger{
		Logger: slog.New(handler),
		name:   "root",
	}

	defaultLogger = logger
	return logger, nil
}

// New creates a new named logger
func New(name string) *Logger {
	if defaultLogger == nil {
		// Initialize with defaults if not already done
		defaultLogger, _ = Init(nil)
	}
	return &Logger{
		Logger: defaultLogger.With(slog.String("logger", name)),
		name:   name,
	}
}

// GetDefault returns the default logger
func GetDefault() *Logger {
	if defaultLogger == nil {
		defaultLogger, _ = Init(nil)
	}
	return defaultLogger
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	attrs := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	return &Logger{
		Logger: l.Logger.With(attrs...),
		name:   l.name,
	}
}

// WithField creates a new logger with an additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With(slog.Any(key, value)),
		name:   l.name,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debug(sprintf(format, args...))
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Info(sprintf(format, args...))
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warn(sprintf(format, args...))
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Error(sprintf(format, args...))
}

// Helper functions

func parseLevel(level LogLevel) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func sprintf(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

// Package-level functions for convenience

// Debug logs a debug message using the default logger
func Debug(msg string, args ...any) {
	GetDefault().Debug(msg, args...)
}

// Info logs an info message using the default logger
func Info(msg string, args ...any) {
	GetDefault().Info(msg, args...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, args ...any) {
	GetDefault().Warn(msg, args...)
}

// Error logs an error message using the default logger
func Error(msg string, args ...any) {
	GetDefault().Error(msg, args...)
}
