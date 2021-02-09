package mail

import (
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
)

type Mail struct {
	senderAddress string
	receiverAddress []string
	smtpHostAddr string
	smtpAuth smtp.Auth
}

// New returns a new instance of a Mail notification service.
func New(senderAddress, smtpHostAddress string) *Mail {
	return &Mail{
		senderAddress:   senderAddress,
		smtpHostAddr:    smtpHostAddress,
		receiverAddress: []string{},
	}
}

// e.g. "", "test@gmail.com", "password123", "smtp.gmail.com"
func (m *Mail) AuthenticateSMTP(identify, username, password, host string) {
	m.smtpAuth = smtp.PlainAuth(identify, username, password, host)
}

func (m *Mail) AddReceivers(addresses ...string) {
	m.receiverAddress = append(m.receiverAddress, addresses...)
}

// Send takes email address and adds them to the internal address list.
// The send method will send a given message to all those addresses.
func (m Mail) Send(subject, message string) error {
	msg := &email.Email{
		To:      m.receiverAddress,
		From:    m.senderAddress,
		Subject: subject,
		HTML:    []byte(message),
		Headers: textproto.MIMEHeader{},
	}

	err := msg.Send(m.smtpHostAddr, m.smtpAuth)
	if err != nil {
		err = errors.Wrap(err, "failed to send mail")
	}

	return err
}