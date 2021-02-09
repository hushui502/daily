package notify

import "errors"

// Notifier is enabled by default
const defaultDisabled = false

// Notify is the central struct for managing notification services and sending messages to them.
type Notify struct {
	Disabled bool
	notifiers []Notifier
}

// ErrSendNotification signals that the notifier failed to send a notification.
var ErrSendNotification = errors.New("send notification")

func New() *Notify {
	notifier := &Notify{
		Disabled:  defaultDisabled,
		notifiers: []Notifier{},
	}

	return notifier
}