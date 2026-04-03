package handlers

import (
	"net/http"

	"github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/middleware"
	"github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/response"
	"github.com/SilentPlaces/basicauth.git/internal/application/usecase"
	logindto "github.com/SilentPlaces/basicauth.git/internal/dto/auth/login"
	refreshtokendto "github.com/SilentPlaces/basicauth.git/internal/dto/auth/refresh_token"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	authUseCase *usecase.AuthUseCase
	logger      appLogger.Logger
}

func NewUserHandler(userUseCase *usecase.UserUseCase, authUseCase *usecase.AuthUseCase, logger appLogger.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, ok := c.Get(middleware.UserContextKey)
	if !ok {
		h.logger.Warn(c.Request.Context(), "get user missing context user id", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusUnauthorized, "User ID not found")
		return
	}

	user, err := h.userUseCase.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error(c.Request.Context(), "get user failed", err, map[string]interface{}{"user_id": userID.(string)})
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.logger.Info(c.Request.Context(), "get user succeeded", map[string]interface{}{"user_id": userID.(string)})
	response.Success(c, http.StatusOK, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req logindto.LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(c.Request.Context(), "login request binding failed", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusBadRequest, "Bad request")
		return
	}

	data, err := h.authUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == usecase.ErrBadRequest {
			h.logger.Warn(c.Request.Context(), "login rejected due to bad request", map[string]interface{}{"email": req.Email})
			response.Error(c, http.StatusBadRequest, "Bad request")
			return
		}
		if err == usecase.ErrWrongCredential {
			h.logger.Warn(c.Request.Context(), "login rejected due to wrong credentials", map[string]interface{}{"email": req.Email})
			response.Error(c, http.StatusUnauthorized, "Wrong Email or Password")
			return
		}
		h.logger.Error(c.Request.Context(), "login failed", err, map[string]interface{}{"email": req.Email})
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.logger.Info(c.Request.Context(), "login succeeded", map[string]interface{}{"user_id": data.User.ID})
	response.Success(c, http.StatusOK, data)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req refreshtokendto.RefreshTokenReqDto
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(c.Request.Context(), "refresh token request binding failed", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusBadRequest, "Bad request")
		return
	}

	data, err := h.authUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Warn(c.Request.Context(), "refresh token failed", map[string]interface{}{"path": c.Request.URL.Path})
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	h.logger.Info(c.Request.Context(), "refresh token succeeded", nil)
	response.Success(c, http.StatusOK, data)
}
