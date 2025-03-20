package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogLevel represents logging levels
type LogLevel string

const (
	// Debug level for detailed debugging information
	DebugLevel LogLevel = "debug"
	// Info level for general information messages
	InfoLevel LogLevel = "info"
	// Warn level for warning conditions
	WarnLevel LogLevel = "warn"
	// Error level for error messages
	ErrorLevel LogLevel = "error"
	// Fatal level for critical errors that cause the application to exit
	FatalLevel LogLevel = "fatal"
)

var (
	// Default logger instance
	defaultLogger zerolog.Logger
)

// Init initializes the logger with specified parameters
func Init(level LogLevel, pretty bool) {
	// Set appropriate log level
	zerolog.SetGlobalLevel(parseLevel(level))

	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	var output io.Writer = os.Stdout

	// Enable pretty logging for development
	if pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
		}
	}

	// Create default logger
	defaultLogger = zerolog.New(output).With().Timestamp().Caller().Logger()

	log.Logger = defaultLogger
}

// parseLevel converts string level to zerolog.Level
func parseLevel(level LogLevel) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// Debug logs a debug message
func Debug() *zerolog.Event {
	return defaultLogger.Debug()
}

// Info logs an info message
func Info() *zerolog.Event {
	return defaultLogger.Info()
}

// Warn logs a warning message
func Warn() *zerolog.Event {
	return defaultLogger.Warn()
}

// Error logs an error message
func Error() *zerolog.Event {
	return defaultLogger.Error()
}

// Fatal logs a fatal message and exits
func Fatal() *zerolog.Event {
	return defaultLogger.Fatal()
}

// With adds a context field to the log
func With() zerolog.Context {
	return defaultLogger.With()
}

// Logger returns the current logger instance
func Logger() zerolog.Logger {
	return defaultLogger
}
