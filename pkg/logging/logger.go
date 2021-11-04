// Inspired by https://github.com/brandur/wanikaniapi/blob/v0.2.0/logger.go
package logging

import (
	"fmt"
	"io"
	"os"
)

const (
	// LevelNull sets a logger to show no messages at all.
	LevelNull Level = 0

	// LevelError sets a logger to show error messages only.
	LevelError Level = 1

	// LevelWarn sets a logger to show warning messages or anything more
	// severe.
	LevelWarn Level = 2

	// LevelInfo sets a logger to show informational messages or anything more
	// severe.
	LevelInfo Level = 3

	// LevelDebug sets a logger to show informational messages or anything more
	// severe.
	LevelDebug Level = 4
)

// Level represents a logging level.
type Level uint32

// LeveledLogger is a leveled logger implementation.
//
// It prints warnings and errors to `os.Stderr` and other messages to
// `os.Stdout`.
type LeveledLogger struct {
	// Level is the minimum logging level that will be emitted by this logger.
	//
	// For example, a Level set to LevelWarn will emit warnings and errors, but
	// not informational or debug messages.
	//
	// Always set this with a constant like LevelWarn because the individual
	// values are not guaranteed to be stable.
	Level Level

	// Internal testing use only.
	StderrOverride io.Writer
	StdoutOverride io.Writer
}

// Debugf logs a debug message using Printf conventions.
func (l *LeveledLogger) Debugf(format string, v ...interface{}) {
	if l.Level >= LevelDebug {
		fmt.Fprintf(l.stdout(), "[DEBUG] "+format+"\n", v...)
	}
}

// Errorf logs a warning message using Printf conventions.
func (l *LeveledLogger) Errorf(format string, v ...interface{}) {
	// Infof logs a debug message using Printf conventions.
	if l.Level >= LevelError {
		fmt.Fprintf(l.stderr(), "[ERROR] "+format+"\n", v...)
	}
}

// Infof logs an informational message using Printf conventions.
func (l *LeveledLogger) Infof(format string, v ...interface{}) {
	if l.Level >= LevelInfo {
		fmt.Fprintf(l.stdout(), "[INFO] "+format+"\n", v...)
	}
}

// Warnf logs a warning message using Printf conventions.
func (l *LeveledLogger) Warnf(format string, v ...interface{}) {
	if l.Level >= LevelWarn {
		fmt.Fprintf(l.stderr(), "[WARN] "+format+"\n", v...)
	}
}

func (l *LeveledLogger) stderr() io.Writer {
	if l.StderrOverride != nil {
		return l.StderrOverride
	}

	return os.Stderr
}

func (l *LeveledLogger) stdout() io.Writer {
	if l.StdoutOverride != nil {
		return l.StdoutOverride
	}

	return os.Stdout
}

// LeveledLoggerInterface provides a basic leveled logging interface for
// printing debug, informational, warning, and error messages.
//
// It's implemented by LeveledLogger and also provides out-of-the-box
// compatibility with a Logrus Logger, but may require a thin shim for use with
// other logging libraries that you use less standard conventions like Zap.
type LeveledLoggerInterface interface {
	// Debugf logs a debug message using Printf conventions.
	Debugf(format string, v ...interface{})

	// Errorf logs a warning message using Printf conventions.
	Errorf(format string, v ...interface{})

	// Infof logs an informational message using Printf conventions.
	Infof(format string, v ...interface{})

	// Warnf logs a warning message using Printf conventions.
	Warnf(format string, v ...interface{})
}

type SimpleLogger interface {
	// Log inserts a log entry.  Arguments may be handled in the manner
	// of fmt.Print, but the underlying logger may also decide to handle
	// them differently.
	Log(v ...interface{})
	// Logf insets a log entry.  Arguments are handled in the manner of
	// fmt.Printf.
	Logf(format string, v ...interface{})
}
