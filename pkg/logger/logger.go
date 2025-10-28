package logger

import (
	"io"
	"log"
	"os"
)

// Level represents the log level
type Level int

const (
	// LevelDebug is for debug messages
	LevelDebug Level = iota
	// LevelInfo is for informational messages
	LevelInfo
	// LevelWarn is for warning messages
	LevelWarn
	// LevelError is for error messages
	LevelError
)

// Logger provides structured logging capabilities
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	level       Level
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(os.Stdout, LevelInfo)
}

// New creates a new logger instance
func New(output io.Writer, level Level) *Logger {
	return &Logger{
		debugLogger: log.New(output, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(output, "[INFO]  ", log.Ldate|log.Ltime),
		warnLogger:  log.New(output, "[WARN]  ", log.Ldate|log.Ltime),
		errorLogger: log.New(output, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       level,
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= LevelDebug {
		l.debugLogger.Println(v...)
	}
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.debugLogger.Printf(format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(v ...interface{}) {
	if l.level <= LevelInfo {
		l.infoLogger.Println(v...)
	}
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.infoLogger.Printf(format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= LevelWarn {
		l.warnLogger.Println(v...)
	}
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.warnLogger.Printf(format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	if l.level <= LevelError {
		l.errorLogger.Println(v...)
	}
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.errorLogger.Printf(format, v...)
	}
}

// Default logger functions
func Debug(v ...interface{})                 { defaultLogger.Debug(v...) }
func Debugf(format string, v ...interface{}) { defaultLogger.Debugf(format, v...) }
func Info(v ...interface{})                  { defaultLogger.Info(v...) }
func Infof(format string, v ...interface{})  { defaultLogger.Infof(format, v...) }
func Warn(v ...interface{})                  { defaultLogger.Warn(v...) }
func Warnf(format string, v ...interface{})  { defaultLogger.Warnf(format, v...) }
func Error(v ...interface{})                 { defaultLogger.Error(v...) }
func Errorf(format string, v ...interface{}) { defaultLogger.Errorf(format, v...) }
func SetLevel(level Level)                   { defaultLogger.SetLevel(level) }
