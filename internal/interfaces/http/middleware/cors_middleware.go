package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/product-management/internal/config"
)

// CORSMiddleware creates CORS middleware based on configuration
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// Set other CORS headers
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// Set allowed methods
		methods := ""
		for i, method := range cfg.CORS.AllowedMethods {
			if i > 0 {
				methods += ", "
			}
			methods += method
		}
		c.Header("Access-Control-Allow-Methods", methods)

		// Set allowed headers
		headers := ""
		for i, header := range cfg.CORS.AllowedHeaders {
			if i > 0 {
				headers += ", "
			}
			headers += header
		}
		c.Header("Access-Control-Allow-Headers", headers)

		// Set max age for preflight requests
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}