package login_dto

import dto "github.com/SilentPlaces/basicauth.git/internal/dto/user"

// LoginResponseDTO represents the response for the login endpoint.
type LoginResponseDTO struct {
	User         *dto.UserResponseDTO `json:"user"`
	Token        string               `json:"token"`
	RefreshToken string               `json:"refreshToken"`
}
