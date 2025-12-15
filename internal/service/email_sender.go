package service

import (
	"fmt"
	"log/slog"
	"net/smtp"
)

type EmailSender interface {
	Send(toEmail, subject, body string) error
}

type SMPTSender struct {
	Host     string
	Port     string
	Username string
	Password string
	MailFrom string
}

// TODO: need config
func NewSMPTSender() *SMPTSender {
	return &SMPTSender{
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     "2525",
		Username: "0df92e2a922276",
		Password: "4122fe43d571cb",
		MailFrom: "from@mail.com",
	}
}

func (s *SMPTSender) Send(toEmail, subject, body string) error{
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte(
		"To: " + toEmail + "\r\n" +
        "From: " + s.Username + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
        "\r\n" +
        body,
	)

	err := smtp.SendMail(addr, auth, s.Username, []string{toEmail}, msg)
	if err != nil {
		slog.Error("Failde to send email", "error", err)
		return err
	}
	return nil
}