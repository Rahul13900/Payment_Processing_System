package middleware

import (
	"user_service/pkg/logger"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs all HTTP requests with context enrichment
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		method := "HTTP:" + c.Request.Method + ":" + c.Request.URL.Path

		// Log request entry
		logger.Info(ctx, method, "Request received")

		// Continue to next handler
		c.Next()

		// Log request exit with duration
		logger.Info(ctx, method,
			map[string]interface{}{
				"status": c.Writer.Status(),
				"path":   c.Request.URL.Path,
				"ip":     c.ClientIP(),
			})
	}
}
