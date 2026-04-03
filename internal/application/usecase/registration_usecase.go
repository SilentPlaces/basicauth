package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/SilentPlaces/basicauth.git/internal/config"
	customerror "github.com/SilentPlaces/basicauth.git/internal/errors"
	mailservice "github.com/SilentPlaces/basicauth.git/internal/services/mail"
	registrationservice "github.com/SilentPlaces/basicauth.git/internal/services/registration"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	validation "github.com/SilentPlaces/basicauth.git/internal/validation/user"
)

const (
	queryParamMailKey  = "email"
	queryParamTokenKey = "token"
	verificationLink   = "https://%s/registration/verify-user?%s=%s&%s=%s"
)

type RegistrationUseCase struct {
	mailService         mailservice.MailService
	registrationService registrationservice.RegistrationService
	registrationConfig  *config.RegistrationConfig
	passwordConfig      *config.RegistrationPasswordConfig
	generalConfig       *config.GeneralConfig
	logger              appLogger.Logger
}

func NewRegistrationUseCase(
	mailService mailservice.MailService,
	registrationService registrationservice.RegistrationService,
	registrationConfig *config.RegistrationConfig,
	passwordConfig *config.RegistrationPasswordConfig,
	generalConfig *config.GeneralConfig,
	logger appLogger.Logger,
) *RegistrationUseCase {
	return &RegistrationUseCase{
		mailService:         mailService,
		registrationService: registrationService,
		registrationConfig:  registrationConfig,
		passwordConfig:      passwordConfig,
		generalConfig:       generalConfig,
		logger:              logger,
	}
}

func (u *RegistrationUseCase) SignUp(ctx context.Context, email, name, password string) error {
	u.logger.Info(ctx, "registration signup requested", map[string]interface{}{"email": email})
	if err := validation.ValidateEmail(email); err != nil {
		u.logger.Warn(ctx, "registration signup invalid email", map[string]interface{}{"email": email})
		return fmt.Errorf("%w: invalid email", ErrBadRequest)
	}
	if err := validation.ValidatePassword(password, u.passwordConfig); err != nil {
		u.logger.Warn(ctx, "registration signup invalid password", map[string]interface{}{"email": email})
		return fmt.Errorf("%w: invalid password", ErrBadRequest)
	}

	token, err := u.registrationService.Signup(email, name, password)
	if err != nil {
		u.logger.Error(ctx, "registration signup service failure", err, map[string]interface{}{"email": email})
		return err
	}

	verificationURL := fmt.Sprintf(
		verificationLink,
		u.generalConfig.Domain,
		queryParamTokenKey,
		url.QueryEscape(token),
		queryParamMailKey,
		url.QueryEscape(email),
	)
	emailBody := fmt.Sprintf(u.registrationConfig.VerificationMailText, verificationURL)
	emailSubject := fmt.Sprintf("Registration Verification Email at %s", u.generalConfig.Domain)

	if err := u.mailService.SendVerificationEmail(
		u.registrationConfig.HostVerificationMailAddress,
		email,
		emailSubject,
		emailBody,
	); err != nil {
		u.logger.Error(ctx, "registration signup email send failed", err, map[string]interface{}{"email": email})
		return err
	}

	u.logger.Info(ctx, "registration signup email sent", map[string]interface{}{"email": email})
	return nil
}

func (u *RegistrationUseCase) VerifyEmail(ctx context.Context, email, token string) error {
	decodedMail, _ := url.QueryUnescape(email)
	decodedToken, _ := url.QueryUnescape(token)
	u.logger.Info(ctx, "registration email verification requested", map[string]interface{}{"email": decodedMail})

	if err := u.registrationService.VerifyToken(decodedMail, decodedToken); err != nil {
		u.logger.Warn(ctx, "registration email token verification failed", map[string]interface{}{"email": decodedMail})
		return err
	}
	if err := u.registrationService.SetUserVerified(decodedMail); err != nil {
		u.logger.Error(ctx, "registration set verified failed", err, map[string]interface{}{"email": decodedMail})
		return err
	}

	u.logger.Info(ctx, "registration email verified", map[string]interface{}{"email": decodedMail})
	return nil
}

func (u *RegistrationUseCase) ResendVerification(ctx context.Context, email string) error {
	u.logger.Info(ctx, "registration resend verification requested", map[string]interface{}{"email": email})
	token, err := u.registrationService.ReloadToken(email)
	if err != nil {
		if errors.Is(err, &customerror.TokenGenerationCountError{}) {
			u.logger.Warn(ctx, "registration resend verification limited", map[string]interface{}{"email": email})
			return err
		}
		u.logger.Error(ctx, "registration resend verification failed", err, map[string]interface{}{"email": email})
		return err
	}

	verificationURL := fmt.Sprintf(
		verificationLink,
		u.generalConfig.Domain,
		queryParamTokenKey,
		url.QueryEscape(token),
		queryParamMailKey,
		url.QueryEscape(email),
	)
	emailBody := fmt.Sprintf(u.registrationConfig.VerificationMailText, verificationURL)
	emailSubject := fmt.Sprintf("Registration Verification Email at %s", u.generalConfig.Domain)

	if err := u.mailService.SendVerificationEmail(
		u.registrationConfig.HostVerificationMailAddress,
		email,
		emailSubject,
		emailBody,
	); err != nil {
		u.logger.Error(ctx, "registration resend email send failed", err, map[string]interface{}{"email": email})
		return err
	}

	u.logger.Info(ctx, "registration resend email sent", map[string]interface{}{"email": email})
	return nil
}
