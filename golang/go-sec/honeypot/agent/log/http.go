package log

import (
	"github.com/sirupsen/logrus"
	"honeypot/agent/util/api"
	"net/http"
	"time"
)

type (
	HttpHook struct {
		HttpClient http.Client
	}
)

func NewHttpHook() (*HttpHook, error) {
	timeout := time.Duration(1 * time.Second)
	client := http.Client{Timeout: timeout}
	return &HttpHook{HttpClient: client}, nil
}

func (hook *HttpHook) Fire(entry *logrus.Entry) error {
	var err error
	field := entry.Data
	data := entry.Message
	_, ok := field["post"]
	if ok {
		err = api.Post(data)
	}
	return err
}

func (hook *HttpHook) Levells() []logrus.Level {
	return logrus.AllLevels
}
