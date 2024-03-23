package log

import (
	"context"
	"log/slog"
)

type loggerCtxKey string

const CONTEXT_LOGGER loggerCtxKey = "ctx.logger"

func LoggerFromCtx(ctx context.Context, fallback *slog.Logger) *slog.Logger {
	requestLogger, ok := ctx.Value(CONTEXT_LOGGER).(*slog.Logger)
	if !ok {
		fallback.Warn("fallback to global logger")
		requestLogger = fallback
	}
	return requestLogger
}
