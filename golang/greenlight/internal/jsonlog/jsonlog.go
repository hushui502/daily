package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Define a Level type to represent the severity level for a log entry.
type Level int8

// Initialize constants which represent a specific severity level. We use the iota
// keyword as a shortcut to assign successive integer values to the constants.
const (
	LevelInfo  Level = iota // Has the value 0.
	LevelError              // Has the value 1.
	LevelFatal              // Has the value 2.
	LevelOff                // Has the value 3.
)

// Return a human-friendly string for the severity level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Define a custom Logger type.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// Return a new Logger instance which writes log entries at or above a minimum severity
// level to a specific output destination.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// Declare some helper methods for writing log entries at the different levels. Notice
// that these all accept a map as the second parameter which can contain any arbitrary
// 'properties' that you want to appear in the log entry.
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

// Print is an internal method for writing the log entry.
func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	// Declare an anonymous struct holding the data for the log entry.
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	// Declare a line variable for holding the actual log entry text.
	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

// We also implement a Write() method on our Logger type so that it satisfies the
// io.Writer interface. This writes a log entry at the ERROR level with no additional
// properties.
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
