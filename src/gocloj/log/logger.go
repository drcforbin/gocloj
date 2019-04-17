package log

import (
	"fmt"
	"sync"
)

type Logger struct {
	mu    sync.Mutex
	Name  string
	Level LogLevel
}

func (logger *Logger) Log(level LogLevel, v ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if level < logger.Level {
		return
	}

	/*
		for i, val := range v {
			if val == nil {
				v[i] = "<nil>"
			}
		}
	*/

	str := fmt.Sprint(v...)
	msg := newMessage(logger.Name, level, str)
	logMessage(msg)
}

func (logger *Logger) Logf(level LogLevel, format string, v ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if level < logger.Level {
		return
	}

	/*
		for i, val := range v {
			if val == nil {
				v[i] = "<nil>"
			}
		}
	*/

	str := fmt.Sprintf(format, v...)
	msg := newMessage(logger.Name, level, str)
	logMessage(msg)
}

func (logger *Logger) Trace(v ...interface{}) {
	logger.Log(Trace, v...)
}

func (logger *Logger) Tracef(format string, v ...interface{}) {
	logger.Logf(Trace, format, v...)
}

func (logger *Logger) Debug(v ...interface{}) {
	logger.Log(Debug, v...)
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.Logf(Debug, format, v...)
}

func (logger *Logger) Info(v ...interface{}) {
	logger.Log(Info, v...)
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.Logf(Info, format, v...)
}

func (logger *Logger) Warn(v ...interface{}) {
	logger.Log(Warn, v...)
}

func (logger *Logger) Warnf(format string, v ...interface{}) {
	logger.Logf(Warn, format, v...)
}

func (logger *Logger) Error(v ...interface{}) {
	logger.Log(Error, v...)
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.Logf(Error, format, v...)
}

func (logger *Logger) Fatal(v ...interface{}) {
	logger.Log(Fatal, v...)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.Logf(Fatal, format, v...)
}

func newLogger(name string) *Logger {
	return &Logger{
		Name:  name,
		Level: defaultLevel,
	}
}
