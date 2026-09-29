package logging

import (
	"context"
	"log/slog"
	"os"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Format int8

const (
	FormatText Format = iota
	FormatJSON
)

type KV struct {
	key   string
	value any
}

func NewKV(key string, value any) KV {
	return KV{key: key, value: value}
}

type Logger interface {
	Debug(ctx context.Context, msg string, kv ...KV)
	Info(ctx context.Context, msg string, kv ...KV)
	Warn(ctx context.Context, msg string, kv ...KV)
	Error(ctx context.Context, msg string, kv ...KV)
	ErrorE(ctx context.Context, msg string, err error, kv ...KV)
}

type slogLogger struct {
	l *slog.Logger
}

var _ Logger = (*slogLogger)(nil)

func (s *slogLogger) Debug(ctx context.Context, msg string, kv ...KV) {
	s.l.LogAttrs(ctx, slog.LevelDebug, msg, toAttrs(kv)...)
}
func (s *slogLogger) Info(ctx context.Context, msg string, kv ...KV) {
	s.l.LogAttrs(ctx, slog.LevelInfo, msg, toAttrs(kv)...)
}
func (s *slogLogger) Warn(ctx context.Context, msg string, kv ...KV) {
	s.l.LogAttrs(ctx, slog.LevelWarn, msg, toAttrs(kv)...)
}
func (s *slogLogger) Error(ctx context.Context, msg string, kv ...KV) {
	s.l.LogAttrs(ctx, slog.LevelError, msg, toAttrs(kv)...)
}
func (s *slogLogger) ErrorE(ctx context.Context, msg string, err error, kv ...KV) {
	attrs := toAttrs(kv)
	attrs = append(attrs, slog.String("error", err.Error()))
	s.l.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}

func toAttrs(kv []KV) []slog.Attr {
	attrs := make([]slog.Attr, len(kv))
	for i, pair := range kv {
		attrs[i] = slog.Any(pair.key, pair.value)
	}
	return attrs
}

var defaultLogger Logger = newSlogLogger(LevelInfo, FormatText, os.Stdout)

func Init(level Level, format Format) {
	defaultLogger = newSlogLogger(level, format, os.Stdout)
}

func SetDefault(l Logger) {
	defaultLogger = l
}

func newSlogLogger(level Level, format Format, w *os.File) Logger {
	opts := &slog.HandlerOptions{Level: toSlogLevel(level)}

	var handler slog.Handler
	if format == FormatJSON {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}
	return &slogLogger{l: slog.New(handler)}
}

func toSlogLevel(l Level) slog.Level {
	switch l {
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

func Debug(ctx context.Context, msg string, kv ...KV) { defaultLogger.Debug(ctx, msg, kv...) }
func Info(ctx context.Context, msg string, kv ...KV)  { defaultLogger.Info(ctx, msg, kv...) }
func Warn(ctx context.Context, msg string, kv ...KV)  { defaultLogger.Warn(ctx, msg, kv...) }
func Error(ctx context.Context, msg string, kv ...KV) { defaultLogger.Error(ctx, msg, kv...) }
func ErrorE(ctx context.Context, msg string, err error, kv ...KV) {
	defaultLogger.ErrorE(ctx, msg, err, kv...)
}
func Fatal(msg string, kv ...KV) {
	Error(context.Background(), msg, kv...)
	os.Exit(1)
}
