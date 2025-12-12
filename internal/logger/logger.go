package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	*log.Logger
	file *os.File
}

var defaultLogger *Logger

func init() {
	var err error
	defaultLogger, err = NewLogger()
	if err != nil {
		log.Printf("Failed to initialize logger: %v", err)
		defaultLogger = &Logger{Logger: log.New(os.Stderr, "", log.LstdFlags)}
	}
}

// NewLogger creates a new logger instance
func NewLogger() (*Logger, error) {
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with timestamp
	logFile := filepath.Join(logDir, fmt.Sprintf("lazyservice_%s.log", 
		time.Now().Format("2006-01-02")))
	
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New(file, "", log.LstdFlags|log.Lshortfile)
	return &Logger{Logger: logger, file: file}, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Global logging functions
func Info(v ...interface{}) {
	defaultLogger.Printf("[INFO] %s", fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Printf("[INFO] "+format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.Printf("[ERROR] %s", fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Printf("[ERROR] "+format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Printf("[WARN] %s", fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Printf("[WARN] "+format, v...)
}

func Debug(v ...interface{}) {
	defaultLogger.Printf("[DEBUG] %s", fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Printf("[DEBUG] "+format, v...)
}

// Close the default logger
func Close() error {
	return defaultLogger.Close()
}
