package log

import (
	"context"
	"io"
	"log/slog"
)

type loggerCtxKey string

const CONTEXT_LOGGER loggerCtxKey = "ctx.logger"

func LoggerFromCtx(ctx context.Context) *slog.Logger {
	requestLogger, ok := ctx.Value(CONTEXT_LOGGER).(*slog.Logger)
	if !ok {
		slog.Default().Warn("fallback to global logger")
		requestLogger = slog.Default()
	}
	return requestLogger
}

func NoopLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}
