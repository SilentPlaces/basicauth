package service

import (
	service "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	"github.com/google/wire"
	"net/smtp"
)

type (
	MailService interface {
		SendVerificationEmail(from string, to string, subject string, body string) error
	}

	mailService struct {
		auth     smtp.Auth
		smtpHost string
		smtpPort string
	}
)

func NewMailService(consul service.ConsulService) MailService {
	cfg, _ := consul.GetSMTPConfig()
	username := cfg[constants.SMTPUsernameKey]
	password := cfg[constants.SMTPPasswordKey]
	smtpHost := cfg[constants.SMTPHostKey]
	smtpPort := cfg[constants.SMTPPortKey]

	auth := smtp.PlainAuth("", username, password, smtpHost)
	return &mailService{auth: auth, smtpHost: smtpHost, smtpPort: smtpPort}
}

func (ms *mailService) SendVerificationEmail(from string, to string, subject string, body string) error {
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)
	addr := ms.smtpHost + ":" + ms.smtpPort

	err := smtp.SendMail(addr, ms.auth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}

var MailServiceProviderSet = wire.NewSet(NewMailService)
