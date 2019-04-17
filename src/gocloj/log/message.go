package log

import (
	"fmt"
	"time"
)

// using these rather than implementing
// LogLevel.String, because they're fixed
// width strings, not a pretty representation
var levelNames = []string{
	"TRACE",
	"DEBUG",
	" INFO",
	" WARN",
	"ERROR",
	"FATAL",
	"OTHER",
}

type message struct {
	logname string
	level   LogLevel
	ts      time.Time

	msg string
}

func newMessage(logname string, level LogLevel, str string) (msg message) {
	msg = message{
		logname: logname,
		level:   level,
		ts:      time.Now(),
	}

	// TODO: UTC?
	timestr := msg.ts.Format("2006-01-02T15:04:05.000")

	if level < Trace || Off < level {
		level = Off // in levelNames, maps to OTHER
	}

	levelStr := levelNames[level]

	msg.msg = fmt.Sprintf("%s %s %s - %s", timestr, levelStr,
		msg.logname, str)

	return
}
