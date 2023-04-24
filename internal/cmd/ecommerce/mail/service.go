package mail

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"net/mail"
	"time"

	mail2 "github.com/go-mail/mail/v2"

	"git.sr.ht/~mariusor/lw"
)

type MailService struct {
	dialer      *mail2.Dialer
	from        string
	configError bool
	logger      lw.Logger
}

//go:embed "template"
var emailTemplates embed.FS

// NewMailer creates new eMail client
func NewMailer(host string, port int, username, password, from string, l lw.Logger) *MailService {
	configError := false
	if host == "" || port == 0 || username == "" || password == "" || from == "" {
		l.Warnf("smtp not properly configured")
		configError = true
	}

	dialer := mail2.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return &MailService{
		dialer:      dialer,
		from:        from,
		configError: configError,
		logger:      l,
	}
}

// isEmailValid checks if email address is correct
func isEmailValid(e string) bool {
	_, err := mail.ParseAddress(e)
	return err == nil
}

// Send creats a mail using the given pattern and sends smtp message
func (m *MailService) Send(recipient string, data any, patterns ...string) error {
	//checking if smtp configuration is ok
	if m.configError {
		return errors.New("smtp not properly configured")
	}

	//checking if recipient email is ok
	if !isEmailValid(recipient) {
		return errors.New("recipient email address is incorrect")
	}

	//patterns
	for i := range patterns {
		patterns[i] = "template/" + patterns[i]
	}
	ts, err := template.New("").Funcs(templateFuncs).ParseFS(emailTemplates, patterns...)
	if err != nil {
		return err
	}

	//creating a message
	msg := mail2.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.from)

	//setting email subject
	subject := new(bytes.Buffer)
	err = ts.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	msg.SetHeader("Subject", subject.String())

	//setting email body
	htmlBody := new(bytes.Buffer)
	err = ts.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	msg.SetBody("text/html", htmlBody.String())

	//sending email to the smtp server
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
