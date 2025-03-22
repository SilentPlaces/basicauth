package mapper

import (
	dto "github.com/SilentPlaces/basicauth.git/internal/dto/user"
	"github.com/SilentPlaces/basicauth.git/internal/models/models"
)

func MapUserToUserResponse(u *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:   u.ID,
		Name: u.Name,
	}
}
