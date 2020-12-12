package bridge

type IMsgSender interface {
	Send(msg string) error
}

type EmailMsgSender struct {
	emails []string
}

func NewEmailMsgSender(emails []string) *EmailMsgSender {
	return &EmailMsgSender{emails:emails}
}

func (s *EmailMsgSender) Send(msg string) error {
	return nil
}

type INotification interface {
	Notify(msg string) error
}

type ErrNotification struct {
	sender IMsgSender
}

func NewErrNotification(sender IMsgSender) *ErrNotification {
	return &ErrNotification{sender:sender}
}

func (c ErrNotification) Notify(msg string) error {
	return c.sender.Send(msg)
}
