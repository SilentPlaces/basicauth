package service

import (
	"errors"
	dto "github.com/SilentPlaces/basicauth.git/internal/dto/user"
	mapper "github.com/SilentPlaces/basicauth.git/internal/mappers/users"
	"github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/hash"
	"github.com/google/wire"
)

type (
	UserService interface {
		GetUser(id string) (*dto.UserResponseDTO, error)
		GetUserByEmail(email string) (*dto.UserResponseDTO, error)
		VerifyLogin(email string, password string) (*dto.UserResponseDTO, error)
	}
	userService struct {
		UserRepo user.UserRepository
	}
)

func NewUserService(userRepo user.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}

func (s *userService) GetUser(id string) (*dto.UserResponseDTO, error) {
	data, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return mapper.MapUserToUserResponse(data), nil
}

func (s *userService) GetUserByEmail(mail string) (*dto.UserResponseDTO, error) {
	data, err := s.UserRepo.GetUserByMail(mail)
	if err != nil {
		return nil, err
	}
	return mapper.MapUserToUserResponse(data), nil
}

func (s *userService) VerifyLogin(mail string, password string) (*dto.UserResponseDTO, error) {
	//get u
	u, err := s.UserRepo.GetUserByMail(mail)
	if err != nil {
		return nil, err
	}
	hashedPassword := helpers.TextToSHA1(password)
	if u.Password != hashedPassword {
		return nil, errors.New("wrong password")
	}
	return mapper.MapUserToUserResponse(u), nil
}

var UserServiceProviderSet = wire.NewSet(NewUserService)
