package mapper

import (
	"github.com/SilentPlaces/basicauth.git/internal/dto/auth/refresh_token"
	dto "github.com/SilentPlaces/basicauth.git/internal/dto/user"
	"github.com/SilentPlaces/basicauth.git/internal/models/models"
	service "github.com/SilentPlaces/basicauth.git/internal/services/auth"
)

func MapUserToUserResponse(u *models.User) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		ID:   u.ID,
		Name: u.Name,
	}
}

func MapTokenToRefreshTokenResDTO(token *service.Tokens) *refresh_token.RefreshTokenResDTO {
	return &refresh_token.RefreshTokenResDTO{
		RefreshToken: token.RefreshToken,
		Token:        token.AccessToken,
	}
}
