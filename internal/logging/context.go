package logging

import "context"

type loggerKeyT struct{}

var loggerKey loggerKeyT

// PushToContext puts a logger into the context.
func PushToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// GetFromContext returns the logger from context, or the default if none.
func GetFromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey).(Logger); ok && l != nil {
		return l
	}
	return defaultLogger
}

// ContextDebug logs at debug level using the logger from context.
func ContextDebug(ctx context.Context, msg string, kv ...KV) {
	GetFromContext(ctx).Debug(ctx, msg, kv...)
}

// ContextInfo logs at info level using the logger from context.
func ContextInfo(ctx context.Context, msg string, kv ...KV) {
	GetFromContext(ctx).Info(ctx, msg, kv...)
}

// ContextWarn logs at warn level using the logger from context.
func ContextWarn(ctx context.Context, msg string, kv ...KV) {
	GetFromContext(ctx).Warn(ctx, msg, kv...)
}

// ContextError logs at error level using the logger from context.
func ContextError(ctx context.Context, msg string, kv ...KV) {
	GetFromContext(ctx).Error(ctx, msg, kv...)
}

// ContextErrorE logs at error level with an error using the logger from context.
func ContextErrorE(ctx context.Context, msg string, err error, kv ...KV) {
	GetFromContext(ctx).ErrorE(ctx, msg, err, kv...)
}
