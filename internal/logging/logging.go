package logging

import (
	"fmt"
	"log"
	"sync"
)

type LogLevel int

const (
	DEBUG_LEVEL LogLevel = iota
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
)

var (
	globalLevel = DEBUG_LEVEL
	mu          sync.Mutex
)

func SetGlobalLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	globalLevel = level
}

func Debugf(format string, args ...interface{}) {
	logMessage(DEBUG_LEVEL, format, args...)
}

func Infof(format string, args ...interface{}) {
	logMessage(INFO_LEVEL, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logMessage(WARN_LEVEL, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logMessage(ERROR_LEVEL, format, args...)
}

func logMessage(level LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	mu.Lock()
	defer mu.Unlock()

	if level >= globalLevel {
		var levelStr string
		switch level {
		case DEBUG_LEVEL:
			levelStr = "DEBUG"
		case INFO_LEVEL:
			levelStr = "INFO"
		case WARN_LEVEL:
			levelStr = "WARN"
		case ERROR_LEVEL:
			levelStr = "ERROR"
		}

		log.Printf("[%s] %s\n", levelStr, message)
	}
}
