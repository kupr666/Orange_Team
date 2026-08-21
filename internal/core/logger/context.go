package core_logger

import (
	"context"
	"log/slog"
)

type contextKey struct{}

func ToContext(
	ctx context.Context,
	log *slog.Logger,
) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

func FromContext(
	ctx context.Context,
	fallback *slog.Logger,
) *slog.Logger {
	log, ok := ctx.Value(contextKey{}).(*slog.Logger)
	if !ok || log == nil {
		return fallback
	}

	return log
}