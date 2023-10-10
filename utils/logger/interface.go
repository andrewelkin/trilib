package logger

import (
	"context"
	"github.com/dlclark/regexp2"
	"io"
)

// LogLevel represents one of five possible logging levels (DEBUG, INFO, WARN, ERROR, FATAL)
type LogLevel uint

const (
	// LogLevelDebug is for verbose &/or development logging; omitted in production
	LogLevelDebug = iota

	// LogLevelInfo is for standard operating information that should be recorded in a production setting
	LogLevelInfo

	// LogLevelWarn is for alerting to unusual conditions that can be handled by application logic
	LogLevelWarn

	// LogLevelError is for failures that may require special handling
	LogLevelError

	// LogLevelFatal is for failures that require immediate proccess termination
	LogLevelFatal
)

// Logger represents a standard logger interface
type Logger interface {
	// Debugf logs formatted arguments with log level DEBUG
	Debugf(namespace, format string, a ...interface{})

	// Infof logs formatted arguments with log level INFO
	Infof(namespace, format string, a ...interface{})

	// Warnf logs formatted arguments with log level WARNING
	Warnf(namespace, format string, a ...interface{})

	// Errorf logs formatted arguments with log level ERROR
	Errorf(namespace, format string, a ...interface{})

	// Fatalf logs formatted arguments with log level CRITICAL
	Fatalf(namespace, format string, a ...interface{})

	// AddOutput adds a log output that receives messages where level is >= minlevel and the namespace matches filter
	AddOutput(filter *regexp2.Regexp, output io.Writer, minLevel LogLevel, ansi bool, trailCR bool, opts ...interface{})

	// Flush flushes the logger and clears any pending messages
	Flush()

	// NewLine inserts \n before next output
	NewLine()

	// NoDateNextLine starts next line without date/debug/servie label
	NoDateNextLine()
}

var globalLogger Logger

// GetOrCreateGlobalLogger retrieves an initialized global logger, or creates one if it has not yet been created with ctx
func GetOrCreateGlobalLogger(ctx context.Context, baseLevel LogLevel) Logger {
	if globalLogger == nil {
		globalLogger = NewAsyncLogger(ctx, baseLevel, nil)
	}
	return globalLogger
}

// GetOrCreateGlobalLoggerEx sames as above but with the filter option
func GetOrCreateGlobalLoggerEx(ctx context.Context, baseLevel LogLevel, stdOutFilter *regexp2.Regexp) Logger {
	if globalLogger == nil {
		globalLogger = NewAsyncLogger(ctx, baseLevel, stdOutFilter)
	}
	return globalLogger
}

// GetGlobalLogger retrieves an initialized global logger, or nil
func GetGlobalLogger() Logger {
	return globalLogger
}

type Formatter interface {
	String(logMessage) string
	NewLine()
	NoDateNextLine()
}
