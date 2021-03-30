package log

import (
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	Logger     *logrus.Entry
	HttpLogger *logrus.Entry
)

func init() {
	l := logrus.New()
	l.Formatter = new(prefixed.TextFormatter)
	l.Level = logrus.DebugLevel

	Logger = l.WithFields(logrus.Fields{"prefix": "proxy-agent"})
	HttpLogger = logrus.New().WithFields(logrus.Fields{"post": true, "prefix": "proxy-agent"})
}
