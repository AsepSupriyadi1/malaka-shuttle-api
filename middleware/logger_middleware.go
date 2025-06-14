package middleware

import (
	"malakashuttle/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware creates a gin.HandlerFunc for logging HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	logger := config.GetLogger()

	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate request processing time
		duration := time.Since(startTime)

		// Get response details
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		// Get any errors that occurred during request processing
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Create log entry with structured fields
		logEntry := logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"duration":    duration.String(),
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
			"user_agent":  userAgent,
		})

		// Log based on status code
		if len(c.Errors) > 0 {
			logEntry.WithField("errors", errorMessage).Error("Request completed with errors")
		} else if statusCode >= 500 {
			logEntry.Error("Server error")
		} else if statusCode >= 400 {
			logEntry.Warn("Client error")
		} else {
			logEntry.Info("Request completed successfully")
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// Simple request ID generator (you might want to use UUID in production)
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" +
		string(rune(time.Now().Nanosecond()%1000000))
}
