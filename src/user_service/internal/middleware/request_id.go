package middleware

import (
	"user_service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds correlation IDs to request context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		// Create enriched context with request ID
		ctx := logger.NewRequestContext(c.Request.Context(), requestID)

		// Store in Gin context
		c.Set("request_id", requestID)

		// Add to response header
		c.Header("X-Request-ID", requestID)

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
