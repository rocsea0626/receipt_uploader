package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	logger *log.Logger
}

var (
	globalLevel LogLevel
	mu          sync.Mutex
)

var logger = NewLogger(INFO, os.Stdout)

func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		logger: log.New(output, "", log.LstdFlags),
	}
}

func SetGlobalLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	globalLevel = level
}

func Debugf(format string, args ...interface{}) {
	logger.logMessage(DEBUG, format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.logMessage(INFO, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.logMessage(WARN, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.logMessage(ERROR, format, args...)
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Println(a ...any) {
	fmt.Println(a...)
}

func (l *Logger) logMessage(level LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	mu.Lock()
	defer mu.Unlock()

	if level >= globalLevel {
		var levelStr string
		switch level {
		case DEBUG:
			levelStr = "DEBUG"
		case INFO:
			levelStr = "INFO"
		case WARN:
			levelStr = "WARN"
		case ERROR:
			levelStr = "ERROR"
		}

		l.logger.Printf("[%s] %s\n", levelStr, message)
	}
}
