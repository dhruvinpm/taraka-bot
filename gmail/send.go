package gmail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Sender struct {
	config SMTPConfig
}

func NewSender(config SMTPConfig) *Sender {
	return &Sender{config: config}
}

func (s *Sender) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: s.config.Host}
	return d.DialAndSend(m)
}

func (s *Sender) SendWithTracking(to, subject, body, messageID, inReplyTo string) error {
	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n",
		s.config.From, to, subject,
	)
	if messageID != "" {
		headers += fmt.Sprintf("Message-ID: <%s>\r\n", messageID)
	}
	if inReplyTo != "" {
		headers += fmt.Sprintf("In-Reply-To: <%s>\r\nReferences: <%s>\r\n", inReplyTo, inReplyTo)
	}
	msg := headers + "\r\n" + body

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(msg))
}
