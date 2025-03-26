package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	resendVerification "github.com/SilentPlaces/basicauth.git/internal/dto/registeration/resend_verification"
	verify_mail_req_dto "github.com/SilentPlaces/basicauth.git/internal/dto/registeration/verify"
	customerror "github.com/SilentPlaces/basicauth.git/internal/errors"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/http"
	"github.com/google/wire"
	"net/http"
	"net/url"

	"github.com/SilentPlaces/basicauth.git/internal/config"
	registerationdto "github.com/SilentPlaces/basicauth.git/internal/dto/registeration"
	consulService "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	mailService "github.com/SilentPlaces/basicauth.git/internal/services/mail"
	registrationService "github.com/SilentPlaces/basicauth.git/internal/services/registration"
	userService "github.com/SilentPlaces/basicauth.git/internal/services/users"
	validation "github.com/SilentPlaces/basicauth.git/internal/validation/user"
	"github.com/julienschmidt/httprouter"
)

type (
	RegistrationController interface {
		SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
		VerifyMail(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
		ResendVerification(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	}

	registrationController struct {
		mailService         mailService.MailService
		registrationService registrationService.RegistrationService
		userService         userService.UserService
		registrationConfig  *config.RegistrationConfig
		passwordConfig      *config.RegistrationPasswordConfig
		generalConfig       *config.GeneralConfig
	}
)

func NewRegistrationController(
	mailService mailService.MailService,
	registrationService registrationService.RegistrationService,
	userService userService.UserService,
	consul consulService.ConsulService,
) RegistrationController {
	registrationCfg := consul.GetRegistrationConfig()
	passwordCfg := consul.GetRegistrationPasswordConfig()
	generalConfig, _ := consul.GetGeneralConfig()

	return &registrationController{
		mailService:         mailService,
		registrationService: registrationService,
		userService:         userService,
		registrationConfig:  registrationCfg,
		passwordConfig:      passwordCfg,
		generalConfig:       generalConfig,
	}
}

// SignUp handles user registration and sends a resend_verification email.
func (rc *registrationController) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request body.
	var requestData registerationdto.RegistrationRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Validate email and password
	if err := validateRequestData(requestData, rc.passwordConfig); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate user and resend_verification token
	token, err := rc.registrationService.Signup(requestData.Email, requestData.Name, requestData.Password)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Error generating resend_verification token")
		return
	}

	// Generate resend_verification URL and email body
	verificationUrl := fmt.Sprintf(verificationLink, rc.generalConfig.Domain, queryParamTokenKey, url.QueryEscape(token), queryParamMailKey, url.QueryEscape(requestData.Email))
	emailBody := fmt.Sprintf(rc.registrationConfig.VerificationMailText, verificationUrl)
	emailSubject := fmt.Sprintf("Registration Verification Email at %s", rc.generalConfig.Domain)

	// Send resend_verification email
	if err := rc.mailService.SendVerificationEmail(
		rc.registrationConfig.HostVerificationMailAddress,
		requestData.Email,
		emailSubject,
		emailBody,
	); err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error: Cannot send emails")
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
}

// Validate the request data for email and password
func validateRequestData(requestData registerationdto.RegistrationRequestDTO, passwordConfig *config.RegistrationPasswordConfig) error {
	// Validate email.
	if err := validation.ValidateEmail(requestData.Email); err != nil {
		return fmt.Errorf("Email is not valid")
	}

	// Validate password.
	if err := validation.ValidatePassword(requestData.Password, passwordConfig); err != nil {
		return fmt.Errorf("Password is not valid")
	}

	return nil
}

// VerifyMail handles email resend_verification (not implemented yet).
func (rc *registrationController) VerifyMail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//get data
	var requestData verify_mail_req_dto.VerifyMailReqDTO
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	mail, _ := url.QueryUnescape(requestData.Mail)
	token, _ := url.QueryUnescape(requestData.Token)
	//validate if token is correct
	err := rc.registrationService.VerifyToken(mail, token)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Verification Link is not valid")
		return
	}

	//resend_verification user
	err = rc.registrationService.SetUserVerified(mail)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *registrationController) ResendVerification(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestData := resendVerification.ResendVerificationRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	token, err := rc.registrationService.ReloadToken(requestData.Email)
	if err != nil {
		if errors.Is(err, &customerror.TokenGenerationCountError{}) {
			helpers.SendErrorResponse(w, http.StatusInternalServerError, "maximum number of attempts reached")
		} else {
			helpers.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}
	// Generate resend_verification URL and email body
	verificationUrl := fmt.Sprintf(verificationLink, rc.generalConfig.Domain, queryParamTokenKey, url.QueryEscape(token), queryParamMailKey, url.QueryEscape(requestData.Email))
	emailBody := fmt.Sprintf(rc.registrationConfig.VerificationMailText, verificationUrl)
	emailSubject := fmt.Sprintf("Registration Verification Email at %s", rc.generalConfig.Domain)

	// Send resend_verification email
	if err := rc.mailService.SendVerificationEmail(
		rc.registrationConfig.HostVerificationMailAddress,
		requestData.Email,
		emailSubject,
		emailBody,
	); err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error: Cannot send emails")
		return
	}

	w.WriteHeader(http.StatusOK)
}

const (
	queryParamMailKey  = "email"
	queryParamTokenKey = "token"
	//not as the same as service route (/register/resend_verification)and points to front-end page
	verificationLink = "https://%s/registration/verify-user?%s=%s&%s=%s"
)

var RegistrationControllerProvider = wire.NewSet(NewRegistrationController)
