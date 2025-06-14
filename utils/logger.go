package utils

import (
	"malakashuttle/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LogWithContext creates a logger with context information from Gin context
func LogWithContext(c *gin.Context) *logrus.Entry {
	logger := config.GetLogger()

	entry := logger.WithFields(logrus.Fields{
		"client_ip": c.ClientIP(),
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
	})

	// Add request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		entry = entry.WithField("request_id", requestID)
	}

	return entry
}

// LogInfo logs an info message with context
func LogInfo(c *gin.Context, message string, fields ...map[string]interface{}) {
	entry := LogWithContext(c)
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Info(message)
}

// LogError logs an error message with context
func LogError(c *gin.Context, message string, err error, fields ...map[string]interface{}) {
	entry := LogWithContext(c)
	if err != nil {
		entry = entry.WithError(err)
	}
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Error(message)
}

// LogWarn logs a warning message with context
func LogWarn(c *gin.Context, message string, fields ...map[string]interface{}) {
	entry := LogWithContext(c)
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Warn(message)
}

// LogDebug logs a debug message with context
func LogDebug(c *gin.Context, message string, fields ...map[string]interface{}) {
	entry := LogWithContext(c)
	if len(fields) > 0 {
		entry = entry.WithFields(logrus.Fields(fields[0]))
	}
	entry.Debug(message)
}
