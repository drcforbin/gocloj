package log

import (
	"reflect"
)

func (logger *Logger) TypeOf(val interface{}) {
	if val != nil {
		logger.Infof("type: %s", reflect.TypeOf(val).String())
	} else {
		logger.Info("type: <nil>")
	}
}
