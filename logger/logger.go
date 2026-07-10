package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const requestLoggerKey = "slog"

var Log *slog.Logger

type Config struct {
	Level  slog.Level
	Format string
}

func Init(cfg Config) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: true,
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		reqLogger := Log.With(
			slog.String("request_id", fmt.Sprintf("%d", time.Now().UnixNano())),
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.String("ip", ctx.ClientIP()),
		)

		ctx.Set(requestLoggerKey, reqLogger)
		ctx.Next()

		reqLogger.Info("request completed",
			slog.Int("status", ctx.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.Int("bytes", ctx.Writer.Size()),
		)
	}
}

func FromContext(ctx *gin.Context) *slog.Logger {
	if l, exists := ctx.Get(requestLoggerKey); exists {
		if log, ok := l.(*slog.Logger); ok {
			return log
		}
	}
	return Log
}
