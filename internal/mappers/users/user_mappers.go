package mapper

import (
	dto "github.com/SilentPlaces/basicauth.git/internal/dto/user"
	"github.com/SilentPlaces/basicauth.git/internal/models/models"
)

func MapUserToUserResponse(u *models.User) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		ID:   u.ID,
		Name: u.Name,
	}
}
