package log

import (
	"sync"
)

type LogLevel int

const (
	Trace LogLevel = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Off
)

// not using sync.Map, because we
// need to lock while iterating in
// SetLevel and when accessing
var once sync.Once
var loggersMu sync.Mutex
var loggers map[string]*Logger
var defaultLevel = Info

var sinksMu sync.Mutex
var sinks []Sink

func initLog() {
	// not using init, to guarantee ordering
	once.Do(func() {
		loggers = map[string]*Logger{}
		sinks = []Sink{}
	})
}

func SetLevel(level LogLevel) {
	initLog()

	loggersMu.Lock()
	defaultLevel = level
	for _, logger := range loggers {
		logger.Level = level
	}
	loggersMu.Unlock()
}

func GetLevel() (level LogLevel) {
	loggersMu.Lock()
	level = defaultLevel
	loggersMu.Unlock()

	return
}

func Get(name string) (logger *Logger) {
	initLog()

	loggersMu.Lock()
	var ok bool
	logger, ok = loggers[name]
	if !ok {
		logger = newLogger(name)
		loggers[name] = logger
	}
	loggersMu.Unlock()

	return
}

func AddSink(sink Sink) {
	initLog()

	sinksMu.Lock()
	sinks = append(sinks, sink)
	sinksMu.Unlock()

	return
}

func logMessage(msg message) {
	initLog()

	if msg.level != Fatal {
		sinksMu.Lock()
		for _, sink := range sinks {
			sink.write(msg)
		}
		sinksMu.Unlock()
	} else {
		panic(msg.msg)
	}
}
