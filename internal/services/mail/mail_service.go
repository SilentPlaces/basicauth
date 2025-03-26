package service

import (
	"fmt"
	"net/smtp"

	"github.com/SilentPlaces/basicauth.git/internal/config"
	consulService "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/google/wire"
)

type (
	MailService interface {
		SendVerificationEmail(from string, to string, subject string, body string) error
	}

	mailService struct {
		auth      smtp.Auth
		smtpHost  string
		smtpPort  string
		appConfig *config.AppConfig
	}
)

// NewMailService retrieves SMTP configuration from Consul and creates a new MailService.
func NewMailService(consul consulService.ConsulService, appConfig *config.AppConfig) (MailService, error) {
	cfg, err := consul.GetSMTPConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get SMTP config: %w", err)
	}

	var auth smtp.Auth

	// Conditionally set authentication only in production
	if appConfig.Environment == "production" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	return &mailService{
		auth:      auth,
		smtpHost:  cfg.Host,
		smtpPort:  cfg.Port,
		appConfig: appConfig,
	}, nil
}

// SendVerificationEmail sends an email using the configured SMTP server.
func (ms *mailService) SendVerificationEmail(from string, to string, subject string, body string) error {
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%s", ms.smtpHost, ms.smtpPort)

	return smtp.SendMail(addr, ms.auth, from, []string{to}, msg)
}

var MailServiceProviderSet = wire.NewSet(NewMailService)
