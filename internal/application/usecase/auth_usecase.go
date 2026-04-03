package usecase

import (
	"context"
	"errors"

	"github.com/SilentPlaces/basicauth.git/internal/application/port"
	logindto "github.com/SilentPlaces/basicauth.git/internal/dto/auth/login"
	refreshtokendto "github.com/SilentPlaces/basicauth.git/internal/dto/auth/refresh_token"
	mapper "github.com/SilentPlaces/basicauth.git/internal/mappers/users"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	validation "github.com/SilentPlaces/basicauth.git/internal/validation/user"
)

var (
	ErrBadRequest      = errors.New("bad request")
	ErrWrongCredential = errors.New("wrong email or password")
	ErrUnauthorized    = errors.New("unauthorized")
)

type AuthUseCase struct {
	userService port.UserReader
	authService port.AuthTokenManager
	logger      appLogger.Logger
}

func NewAuthUseCase(userService port.UserReader, authService port.AuthTokenManager, logger appLogger.Logger) *AuthUseCase {
	return &AuthUseCase{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

func (u *AuthUseCase) Login(ctx context.Context, email, password string) (*logindto.LoginResponseDTO, error) {
	u.logger.Info(ctx, "auth login requested", map[string]interface{}{"email": email})
	if err := validation.ValidateEmail(email); err != nil {
		u.logger.Warn(ctx, "auth login validation failed", map[string]interface{}{"email": email})
		return nil, ErrBadRequest
	}

	userData, err := u.userService.VerifyLogin(email, password)
	if err != nil {
		u.logger.Warn(ctx, "auth login credentials rejected", map[string]interface{}{"email": email})
		return nil, ErrWrongCredential
	}

	token, err := u.authService.GenerateToken(userData.ID)
	if err != nil {
		u.logger.Error(ctx, "auth login token generation failed", err, map[string]interface{}{"user_id": userData.ID})
		return nil, err
	}

	u.logger.Info(ctx, "auth login succeeded", map[string]interface{}{"user_id": userData.ID})
	return &logindto.LoginResponseDTO{
		User:         userData,
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (u *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*refreshtokendto.RefreshTokenResDTO, error) {
	u.logger.Info(ctx, "auth refresh token requested", nil)
	tokens, err := u.authService.RefreshToken(refreshToken)
	if err != nil {
		u.logger.Warn(ctx, "auth refresh token failed", map[string]interface{}{"reason": "invalid_or_expired_refresh_token"})
		return nil, err
	}

	mapped := mapper.MapTokenToRefreshTokenResDTO(tokens)
	u.logger.Info(ctx, "auth refresh token succeeded", nil)
	return mapped, nil
}
