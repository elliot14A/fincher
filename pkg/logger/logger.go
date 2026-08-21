package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

func init() {
	Init("development", os.Stdout)
}

// Init configures the global logger based on environment.
func Init(env string, w io.Writer) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	if env == "production" {
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Debug logs at LevelDebug.
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// DebugContext logs at LevelDebug with context.
func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}

// Info logs at LevelInfo.
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// InfoContext logs at LevelInfo with context.
func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}

// Warn logs at LevelWarn.
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// WarnContext logs at LevelWarn with context.
func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}

// Error logs at LevelError.
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// ErrorContext logs at LevelError with context.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}

// With returns a Logger that includes the given attributes in each output operation.
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}
