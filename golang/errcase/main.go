package main

import (
	"errors"
	"github.com/sirupsen/logrus"
)

// field
func example(accountNumber string) error {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	ctxFields := logrus.Fields{
		"accountNumber": accountNumber,
		"appName": "my-app",
	}

	err := errors.New("Some test error while doing happy processing")
	if err != nil {
		logrus.WithFields(ctxFields).WithError(err).Error("Err Msg")
		return err
	}

	return nil
}

func main() {
	example("Act1")
}