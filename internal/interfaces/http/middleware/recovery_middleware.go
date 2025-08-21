package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware creates a recovery middleware that handles panics
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic
		logPanic(c, recovered)

		// Return appropriate error response
		if gin.Mode() == gin.ReleaseMode {
			// In production, don't expose internal errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": "An unexpected error occurred",
			})
		} else {
			// In development, provide more details
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": "An unexpected error occurred",
				"details": recovered,
				"stack":   string(debug.Stack()),
			})
		}
	})
}

// logPanic logs panic information
func logPanic(c *gin.Context, recovered interface{}) {
	logData := map[string]interface{}{
		"timestamp":  getCurrentTime(),
		"level":      "error",
		"type":       "panic",
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"client_ip":  c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"panic":      recovered,
		"stack":      string(debug.Stack()),
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		logData["request_id"] = requestID
	}

	// Add user information if available
	if userID, exists := c.Get("user_id"); exists {
		logData["user_id"] = userID
	}

	logJSON, _ := json.Marshal(logData)
	log.Printf("PANIC: %s", string(logJSON))
}

// ErrorHandlerMiddleware handles HTTP errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Log the error
			logError(c, err)

			// If response wasn't written yet, write error response
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal Server Error",
					"message": err.Error(),
				})
			}
		}
	}
}

// logError logs error information
func logError(c *gin.Context, err *gin.Error) {
	logData := map[string]interface{}{
		"timestamp":  getCurrentTime(),
		"level":      "error",
		"type":       "request_error",
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"client_ip":  c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"error":      err.Error(),
		"error_type": err.Type,
	}

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		logData["request_id"] = requestID
	}

	// Add user information if available
	if userID, exists := c.Get("user_id"); exists {
		logData["user_id"] = userID
	}

	logJSON, _ := json.Marshal(logData)
	log.Printf("ERROR: %s", string(logJSON))
}

// getCurrentTime returns current time in ISO format
func getCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}