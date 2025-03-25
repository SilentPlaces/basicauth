package service

import (
	"errors"
	"github.com/SilentPlaces/basicauth.git/internal/models/models"
	repository "github.com/SilentPlaces/basicauth.git/internal/repositories/registration"
	userRepo "github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/strings"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"log"
)

type (
	RegistrationService interface {
		Signup(email string, name string, password string) (string, error) //sign up user and generate verification token
		VerifyToken(email, token string) error
	}

	registrationService struct {
		registrationRepository repository.RegistrationRepository
		userRepository         userRepo.UserRepository
	}
)

func NewUserRegistrationService(verificationRepo repository.RegistrationRepository, userRepository userRepo.UserRepository) RegistrationService {
	return &registrationService{
		registrationRepository: verificationRepo,
		userRepository:         userRepository,
	}
}

// Signup do the logic of sign up user in out system
func (s *registrationService) Signup(email string, name string, password string) (string, error) {
	// Check if the user already exists by email
	existingUser, err := s.userRepository.GetUserByMail(email)
	if err != nil {
		log.Printf("Error getting user by mail: %v", err)
		// Return an error if the user exists
		if existingUser != nil {
			return "", errors.New("this email is already in use, please login")
		}
		return "", err
	}

	// Insert the new user
	dbUser, err := s.userRepository.InsertUser(models.User{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	// Generate verification token
	token, err := helpers.GenerateRandomString(64)
	if err != nil {
		log.Printf("Error generating token: %s", err.Error())
		return "", err
	}

	// Save token to Redis
	err = s.registrationRepository.SetVerifyToken(email, token)
	if err != nil {
		_ = s.userRepository.DeleteUserByID(dbUser.ID) // Rollback by deleting user
		log.Println(err)
		return "", err
	}

	// Return the generated token
	return token, nil
}

func (s *registrationService) VerifyToken(email, token string) error {
	result, err := s.registrationRepository.GetVerifyToken(email)
	if errors.Is(err, redis.Nil) {
		return errors.New("token does not exist")
	} else if err != nil {
		return err
	}

	if result != token {
		return errors.New("token does not match")
	}

	err = s.registrationRepository.DeleteToken(email)
	if err != nil {
		log.Printf("Error deleting token: %s", err.Error())
	}
	return nil
}

var UserRegistrationServiceProvider = wire.NewSet(NewUserRegistrationService)
