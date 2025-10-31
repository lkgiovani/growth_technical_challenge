package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lkgiovani/growth_technical_challenge/pkg/logger"
	"go.uber.org/zap"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			log.Error("Request completed with errors", append(fields, zap.String("errors", c.Errors.String()))...)
		} else if statusCode >= 500 {
			log.Error("Request failed with server error", fields...)
		} else if statusCode >= 400 {
			log.Warn("Request failed with client error", fields...)
		} else {
			log.Info("Request completed successfully", fields...)
		}
	}
}
