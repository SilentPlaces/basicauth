package port

import (
	"github.com/SilentPlaces/basicauth.git/internal/dto/user"
	authservice "github.com/SilentPlaces/basicauth.git/internal/services/auth"
)

type UserReader interface {
	GetUser(id string) (*dto.UserResponseDTO, error)
	VerifyLogin(email string, password string) (*dto.UserResponseDTO, error)
}

type AuthTokenManager interface {
	GenerateToken(userID string) (*authservice.Tokens, error)
	RefreshToken(token string) (*authservice.Tokens, error)
	ValidateToken(token string) error
	ExtractClaims(token string) (*authservice.Claims, error)
}
