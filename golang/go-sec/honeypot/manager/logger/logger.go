package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	Logger *logrus.Entry
)

func init() {
	l := logrus.New()
	l.Formatter = new(prefixed.TextFormatter)
	l.Level = logrus.DebugLevel
	Logger = l.WithFields(logrus.Fields{"prefix": "proxy-manager"})
}
