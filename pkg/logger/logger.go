package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents log level
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger implements domain.Logger interface
type Logger struct {
	level  Level
	prefix string
}

// New creates a new logger instance
func New(prefix string) *Logger {
	level := InfoLevel
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = DebugLevel
	}
	return &Logger{
		level:  level,
		prefix: prefix,
	}
}

// Debug logs debug level messages
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.level <= DebugLevel {
		l.log("DEBUG", msg, fields...)
	}
}

// Info logs info level messages
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.level <= InfoLevel {
		l.log("INFO", msg, fields...)
	}
}

// Warn logs warn level messages
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.level <= WarnLevel {
		l.log("WARN", msg, fields...)
	}
}

// Error logs error level messages
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.level <= ErrorLevel {
		l.log("ERROR", msg, fields...)
	}
}

// Fatal logs fatal level messages and exits
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}

func (l *Logger) log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := fmt.Sprintf("[%s] [%s] [%s]", timestamp, level, l.prefix)

	if len(fields) > 0 {
		log.Printf("%s %s | %v", prefix, msg, fields)
	} else {
		log.Printf("%s %s", prefix, msg)
	}
}

// WithPrefix creates a new logger with a different prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		level:  l.level,
		prefix: prefix,
	}
}
