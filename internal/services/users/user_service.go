package service

import (
	dto "github.com/SilentPlaces/basicauth.git/internal/dto/user"
	mapper "github.com/SilentPlaces/basicauth.git/internal/mappers/users"
	"github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	"github.com/google/wire"
)

type (
	UserService interface {
		GetUser(id string) (*dto.UserResponse, error)
	}
	userService struct {
		UserRepo user.UserRepository
	}
)

func NewUserService(userRepo user.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}

func (s *userService) GetUser(id string) (*dto.UserResponse, error) {
	data, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return mapper.MapUserToUserResponse(data), nil
}

var UserServiceProviderSet = wire.NewSet(NewUserService)
