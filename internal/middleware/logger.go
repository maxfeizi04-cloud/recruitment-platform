package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// SetLogLevel sets the global log level. Pass "debug" for verbose output.
func SetLogLevel(level string) {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}

// RequestLogger logs every HTTP request with structured fields.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		traceID, _ := c.Get("trace_id")
		clientIP := c.ClientIP()

		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("ip", clientIP),
			slog.String("trace_id", traceID.(string)),
		}

		if status >= 500 {
			Logger.LogAttrs(c.Request.Context(), slog.LevelError, "request completed", attrs...)
		} else if status >= 400 {
			Logger.LogAttrs(c.Request.Context(), slog.LevelWarn, "request completed", attrs...)
		} else {
			Logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "request completed", attrs...)
		}
	}
}
