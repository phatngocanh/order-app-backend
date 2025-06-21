package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pna/order-app-backend/internal/utils/env"
)

func CorsMiddleware() gin.HandlerFunc {
	// Try to get allowed origins from environment, fallback to frontend URL
	allowedOrigins, err := env.GetEnv("ALLOWED_ORIGINS")
	if err != nil {
		// Default to frontend URL for development
		allowedOrigins = "http://localhost:3000"
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
