package logger

import (
	"fmt"
	"os"
)

// LogLevel represents the verbosity level of logging
type LogLevel int

const (
	// LevelError is for critical errors
	LevelError LogLevel = iota
	// LevelInfo is for general information
	LevelInfo
	// LevelDebug is for more detailed information
	LevelDebug
	// LevelTrace is for the most detailed information
	LevelTrace
)

var currentLevel LogLevel = LevelError

// SetVerbosity sets the log level based on the verbose flag count
func SetVerbosity(verbose int) {
	if verbose <= 0 {
		currentLevel = LevelError
	} else if verbose == 1 {
		currentLevel = LevelInfo
	} else if verbose == 2 {
		currentLevel = LevelDebug
	} else {
		currentLevel = LevelTrace
	}
}

// Error logs error messages
func Error(format string, args ...interface{}) {
	log(LevelError, format, args...)
}

// Info logs informational messages
func Info(format string, args ...interface{}) {
	log(LevelInfo, format, args...)
}

// Debug logs debug information
func Debug(format string, args ...interface{}) {
	log(LevelDebug, format, args...)
}

// Trace logs the most detailed information
func Trace(format string, args ...interface{}) {
	log(LevelTrace, format, args...)
}

// log prints a message to stderr if the current log level is high enough
func log(level LogLevel, format string, args ...interface{}) {
	if level <= currentLevel {
		prefix := ""
		switch level {
		case LevelError:
			prefix = "ERROR: "
		case LevelInfo:
			prefix = "INFO: "
		case LevelDebug:
			prefix = "DEBUG: "
		case LevelTrace:
			prefix = "TRACE: "
		}
		fmt.Fprintf(os.Stderr, prefix+format+"\n", args...)
	}
}
