package ftp

import (
	"fmt"

	log "github.com/fclairamb/go-log"
	"github.com/sirupsen/logrus"
)

type LoggerAdapter struct {
	*logrus.Entry
}

func extractKeyvals(keyvals ...interface{}) (fields logrus.Fields) {
	fields = make(logrus.Fields)
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}

	for i := 0; i < len(keyvals); i += 2 {
		k, v := keyvals[i], keyvals[i+1]
		if key, ok := k.(string); ok {
			fields[key] = v
		} else {
			fields[fmt.Sprintf("%+v", k)] = v
		}
	}

	return
}

func (l *LoggerAdapter) Debug(event string, keyvals ...interface{}) {

	l.WithFields(extractKeyvals(keyvals...)).Debugf("%s", event)
}

func (l *LoggerAdapter) Error(event string, keyvals ...interface{}) {
	l.WithFields(extractKeyvals(keyvals...)).Errorf("%s", event)
}

func (l *LoggerAdapter) Warn(event string, keyvals ...interface{}) {
	l.WithFields(extractKeyvals(keyvals...)).Warnf("%s", event)
}

func (l *LoggerAdapter) Info(event string, keyvals ...interface{}) {
	l.WithFields(extractKeyvals(keyvals...)).Infof("%s", event)
}

func (l *LoggerAdapter) With(keyvals ...interface{}) log.Logger {
	return &LoggerAdapter{l.WithFields(extractKeyvals(keyvals))}
}

func Test() log.Logger {
	l := &LoggerAdapter{}
	return l
}
