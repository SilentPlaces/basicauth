package handlers

import (
	"errors"
	"net/http"

	"github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/response"
	"github.com/SilentPlaces/basicauth.git/internal/application/usecase"
	registerationdto "github.com/SilentPlaces/basicauth.git/internal/dto/registeration"
	resendverification "github.com/SilentPlaces/basicauth.git/internal/dto/registeration/resend_verification"
	verifymaildto "github.com/SilentPlaces/basicauth.git/internal/dto/registeration/verify"
	customerror "github.com/SilentPlaces/basicauth.git/internal/errors"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	registrationUseCase *usecase.RegistrationUseCase
	logger              appLogger.Logger
}

func NewRegistrationHandler(registrationUseCase *usecase.RegistrationUseCase, logger appLogger.Logger) *RegistrationHandler {
	return &RegistrationHandler{registrationUseCase: registrationUseCase, logger: logger}
}

func (h *RegistrationHandler) SignUp(c *gin.Context) {
	var req registerationdto.RegistrationRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(c.Request.Context(), "signup request binding failed", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := h.registrationUseCase.SignUp(c.Request.Context(), req.Email, req.Name, req.Password); err != nil {
		if errors.Is(err, usecase.ErrBadRequest) {
			h.logger.Warn(c.Request.Context(), "signup rejected due to bad request", map[string]interface{}{"email": req.Email})
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		h.logger.Error(c.Request.Context(), "signup failed", err, map[string]interface{}{"email": req.Email})
		response.Error(c, http.StatusInternalServerError, "Error generating verification token")
		return
	}

	h.logger.Info(c.Request.Context(), "signup succeeded", map[string]interface{}{"email": req.Email})
	response.Success(c, http.StatusCreated, nil)
}

func (h *RegistrationHandler) VerifyMail(c *gin.Context) {
	var req verifymaildto.VerifyMailReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(c.Request.Context(), "verify mail request binding failed", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := h.registrationUseCase.VerifyEmail(c.Request.Context(), req.Mail, req.Token); err != nil {
		h.logger.Warn(c.Request.Context(), "verify mail failed", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusInternalServerError, "Verification Link is not valid")
		return
	}

	h.logger.Info(c.Request.Context(), "verify mail succeeded", nil)
	response.Success(c, http.StatusOK, nil)
}

func (h *RegistrationHandler) ResendVerification(c *gin.Context) {
	var req resendverification.ResendVerificationRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(c.Request.Context(), "resend verification request binding failed", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := h.registrationUseCase.ResendVerification(c.Request.Context(), req.Email); err != nil {
		if errors.Is(err, &customerror.TokenGenerationCountError{}) {
			h.logger.Warn(c.Request.Context(), "resend verification rate limited", map[string]interface{}{"email": req.Email})
			response.Error(c, http.StatusTooManyRequests, "maximum number of attempts reached")
			return
		}
		h.logger.Error(c.Request.Context(), "resend verification failed", err, map[string]interface{}{"email": req.Email})
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.logger.Info(c.Request.Context(), "resend verification succeeded", map[string]interface{}{"email": req.Email})
	response.Success(c, http.StatusOK, nil)
}
