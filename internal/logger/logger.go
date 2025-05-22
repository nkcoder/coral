// Package logger provides a simple logging interface for the application
package logger

import (
	"log"
	"os"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// DebugLevel logs detailed information for debugging
	DebugLevel LogLevel = iota
	// InfoLevel logs general operational information
	InfoLevel
	// WarnLevel logs issues that might need attention
	WarnLevel
	// ErrorLevel logs issues that need addressing
	ErrorLevel
)

var (
	// Current logging level
	currentLevel = InfoLevel

	// Loggers for different levels
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
	infoLogger  = log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime)
	warnLogger  = log.New(os.Stderr, "WARN:  ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
)

// SetLevel sets the minimum logging level
func SetLevel(level LogLevel) {
	currentLevel = level
}

// Debug logs messages at DEBUG level
func Debug(format string, args ...any) {
	if currentLevel <= DebugLevel {
		debugLogger.Printf(format, args...)
	}
}

// Info logs messages at INFO level
func Info(format string, args ...any) {
	if currentLevel <= InfoLevel {
		infoLogger.Printf(format, args...)
	}
}

// Warn logs messages at WARN level
func Warn(format string, args ...any) {
	if currentLevel <= WarnLevel {
		warnLogger.Printf(format, args...)
	}
}

// Error logs messages at ERROR level
func Error(format string, args ...any) {
	if currentLevel <= ErrorLevel {
		errorLogger.Printf(format, args...)
	}
}

// Fatal logs an error message and exits the program
func Fatal(format string, args ...any) {
	errorLogger.Printf(format, args...)
	os.Exit(1)
}
