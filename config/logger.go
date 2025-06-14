package config

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.Mkdir(logsDir, 0755)
	}

	// Generate date-based log file name
	currentDate := time.Now().Format("2006-01-02")
	logFileName := "app-" + currentDate + ".log"
	logFile := filepath.Join(logsDir, logFileName)

	// Configure log rotation with lumberjack
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // megabytes
		MaxBackups: 5,    // number of old log files to keep
		MaxAge:     30,   // days
		Compress:   true, // compress old log files
	}

	// Set output to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
	Logger.SetOutput(multiWriter)

	// Set log format to JSON for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set log level based on environment
	env := os.Getenv("GIN_MODE")
	if env == "release" {
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		Logger.SetLevel(logrus.DebugLevel)
	}

	Logger.WithField("log_file", logFileName).Info("Logger initialized successfully with date-based file output")
}

func GetLogger() *logrus.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

// GetLogFileName generates a log file name based on the current date
func GetLogFileName() string {
	currentDate := time.Now().Format("2006-01-02")
	return "app-" + currentDate + ".log"
}

// RotateLogDaily can be called to rotate logs daily at midnight
// This is optional and can be used with a scheduler like cron
func RotateLogDaily() {
	if Logger == nil {
		return
	}

	logsDir := "logs"
	newLogFileName := GetLogFileName()
	newLogFile := filepath.Join(logsDir, newLogFileName)

	// Create new lumberjack logger with new filename
	lumberjackLogger := &lumberjack.Logger{
		Filename:   newLogFile,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	// Update logger output
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
	Logger.SetOutput(multiWriter)

	Logger.WithField("new_log_file", newLogFileName).Info("Log file rotated for new day")
}
