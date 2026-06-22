package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type traceIDKey struct{}

// TraceID returns the trace ID from context, or empty string.
func TraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey{}).(string); ok {
		return id
	}
	return ""
}

// Tracing generates a UUID trace ID per request, sets it in the Gin context
// and as X-Trace-ID response header, and stores it in request context.
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()[:8] // short 8-char trace ID
		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), traceIDKey{}, traceID))
		c.Next()
	}
}
