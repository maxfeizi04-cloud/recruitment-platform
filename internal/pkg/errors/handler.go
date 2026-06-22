package errors

import (
	"log/slog"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Respond 统一错误响应
func Respond(c *gin.Context, err *AppError) {
	traceID := middleware.TraceID(c.Request.Context())
	status := err.WithStatus()

	slog.Warn("request error",
		"trace_id", traceID,
		"code", err.Code,
		"message", err.Message,
		"http_status", status,
		"path", c.Request.URL.Path,
	)

	c.JSON(status, gin.H{
		"code":     err.Code,
		"message":  err.Message,
		"trace_id": traceID,
	})
}

// RespondMsg 自定义消息的错误响应
func RespondMsg(c *gin.Context, httpStatus int, code int, message string) {
	traceID := middleware.TraceID(c.Request.Context())

	slog.Warn("request error",
		"trace_id", traceID,
		"code", code,
		"message", message,
		"http_status", httpStatus,
		"path", c.Request.URL.Path,
	)

	c.JSON(httpStatus, gin.H{
		"code":     code,
		"message":  message,
		"trace_id": traceID,
	})
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	traceID := middleware.TraceID(c.Request.Context())
	c.JSON(200, gin.H{
		"code":     0,
		"message":  "success",
		"data":     data,
		"trace_id": traceID,
	})
}
