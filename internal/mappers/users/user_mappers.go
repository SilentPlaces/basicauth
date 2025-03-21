package users

import (
	dto "github.com/SilentPlaces/basicauth.git/internal/dto/user"
	"github.com/SilentPlaces/basicauth.git/internal/models/user"
)

func MapUserToUserResponse(u *user.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:   u.ID,
		Name: u.Name,
	}
}
