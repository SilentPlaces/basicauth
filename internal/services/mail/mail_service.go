package service

import (
	"net/smtp"

	consulService "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/google/wire"
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

// NewMailService retrieves the SMTP configuration from Consul and creates a new MailService.
// It panics if the configuration retrieval fails.
func NewMailService(consul consulService.ConsulService) MailService {
	cfg, err := consul.GetSMTPConfig()
	if err != nil {
		panic("failed to get SMTP config: " + err.Error())
	}

	username := cfg.Username
	password := cfg.Password
	smtpHost := cfg.Host
	smtpPort := cfg.Port

	auth := smtp.PlainAuth("", username, password, smtpHost)
	return &mailService{auth: auth, smtpHost: smtpHost, smtpPort: smtpPort}
}

// SendVerificationEmail sends an email using the configured SMTP server.
func (ms *mailService) SendVerificationEmail(from string, to string, subject string, body string) error {
	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)
	addr := ms.smtpHost + ":" + ms.smtpPort

	return smtp.SendMail(addr, ms.auth, from, []string{to}, msg)
}

var MailServiceProviderSet = wire.NewSet(NewMailService)
