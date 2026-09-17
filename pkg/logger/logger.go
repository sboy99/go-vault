package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type LogLevelEnum int

const (
	DEBUG LogLevelEnum = iota
	INFO
	WARN
	ERROR
)

var (
	log   *slog.Logger
	level LogLevelEnum
	mu    sync.Mutex
)

// Init initializes a JSON slog logger for container-friendly output.
func Init(lvl LogLevelEnum) {
	mu.Lock()
	defer mu.Unlock()
	level = lvl
	var slogLevel slog.Level
	switch lvl {
	case DEBUG:
		slogLevel = slog.LevelDebug
	case INFO:
		slogLevel = slog.LevelInfo
	case WARN:
		slogLevel = slog.LevelWarn
	default:
		slogLevel = slog.LevelError
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	log = slog.New(handler)
	slog.SetDefault(log)
}

func Debug(format string, v ...interface{}) {
	ensure()
	if level > DEBUG {
		return
	}
	log.Debug(fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	ensure()
	if level > INFO {
		return
	}
	log.Info(fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	ensure()
	if level > WARN {
		return
	}
	log.Warn(fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	ensure()
	log.Error(fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	Error(format, v...)
	os.Exit(1)
}

func Panic(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	Error("%s", msg)
	panic(msg)
}

func Logger() *slog.Logger {
	ensure()
	return log
}

func ensure() {
	mu.Lock()
	defer mu.Unlock()
	if log == nil {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
		log = slog.New(handler)
		level = INFO
	}
}
