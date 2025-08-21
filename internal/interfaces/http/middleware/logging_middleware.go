package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logData := map[string]interface{}{
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
			"status_code":  param.StatusCode,
			"latency":      param.Latency.String(),
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
		}

		// Add request ID if available
		if requestID := param.Keys["request_id"]; requestID != nil {
			logData["request_id"] = requestID
		}

		// Add user information if available
		if userID := param.Keys["user_id"]; userID != nil {
			logData["user_id"] = userID
		}

		logJSON, _ := json.Marshal(logData)
		return string(logJSON) + "\n"
	})
}

// RequestResponseLoggingMiddleware logs request and response details
func RequestResponseLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request
		requestBody := ""
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore the body for further processing
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Create a response writer that captures the response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log the request and response
		logData := map[string]interface{}{
			"timestamp":     time.Now().Format(time.RFC3339),
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query":         c.Request.URL.RawQuery,
			"status_code":   c.Writer.Status(),
			"latency":       latency.String(),
			"client_ip":     c.ClientIP(),
			"user_agent":    c.Request.UserAgent(),
			"request_body":  requestBody,
			"response_body": blw.body.String(),
		}

		// Add request ID if available
		if requestID, exists := c.Get("request_id"); exists {
			logData["request_id"] = requestID
		}

		// Add user information if available
		if userID, exists := c.Get("user_id"); exists {
			logData["user_id"] = userID
		}

		// Don't log sensitive information in production
		if gin.Mode() == gin.ReleaseMode {
			delete(logData, "request_body")
			delete(logData, "response_body")
		}

		logJSON, _ := json.Marshal(logData)
		log.Println(string(logJSON))
	}
}

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + 
		   time.Now().Format("000000") // microseconds
}