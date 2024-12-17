package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// LogLevelEnum defines the severity of the log messages.
type LogLevelEnum int

const (
	DEBUG LogLevelEnum = iota
	INFO
	WARN
	ERROR
)

var (
	logger    *log.Logger
	logLevel  LogLevelEnum
	mutex     sync.Mutex
	logLevels = map[LogLevelEnum]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}
)

// Init initializes the logger with a log level.
func Init(level LogLevelEnum) {
	mutex.Lock()
	defer mutex.Unlock()

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	logLevel = level
}

// Log writes a message to stdout if the level is enabled.
func Log(level LogLevelEnum, format string, v ...interface{}) {
	if level < logLevel {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()

	logMessage := fmt.Sprintf("[%s] %s ", logLevels[level], fmt.Sprintf(format, v...))
	logger.Println(logMessage)
}

// Debug logs a debug message.
func Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}

// Info logs an informational message.
func Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

// Warn logs a warning message.
func Warn(format string, v ...interface{}) {
	Log(WARN, format, v...)
}

// Error logs an error message.
func Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}

// Panic logs a message at the ERROR level and then panics.
func Panic(format string, v ...interface{}) {
	Log(ERROR, format, v...)
	panic(fmt.Sprintf(format, v...))
}

// Fatal logs a message at the ERROR level and then exits the application.
func Fatal(format string, v ...interface{}) {
	Log(ERROR, format, v...)
	os.Exit(1)
}
