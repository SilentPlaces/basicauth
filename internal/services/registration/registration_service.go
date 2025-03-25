package service

import (
	"database/sql"
	"errors"
	"fmt"
	custom_error "github.com/SilentPlaces/basicauth.git/internal/errors"
	"github.com/SilentPlaces/basicauth.git/internal/models/models"
	repository "github.com/SilentPlaces/basicauth.git/internal/repositories/registration"
	userRepo "github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/strings"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type (
	RegistrationService interface {
		Signup(email string, name string, password string) (string, error)
		VerifyToken(email, token string) error
		SetUserVerified(email string) error
		ReloadToken(email string) (string, error)
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

// Signup handles user registration, checks if the email exists, and generates a verify token
func (s *registrationService) Signup(email string, name string, password string) (string, error) {
	// Check if user already exists by email
	existingUser, err := s.userRepository.GetUserByMail(email)
	if err != nil {
		logError("Error getting user by mail: %v", err)
		if existingUser != nil {
			return "", errors.New("this email is already in use, please login")
		}
		return "", err
	}

	// Check if the user has already generated too many tokens in pas 24 hours
	canGenerate, err := s.registrationRepository.CanGenerateToken(email)
	if err != nil {
		return "", fmt.Errorf("error checking token generation limit: %w", err)
	}
	if !canGenerate {
		return "", custom_error.NewTokenGenerationError("500", "you cannot generate more than 5 tokens in the past 24 hours")
	}

	// Insert new user
	dbUser, err := s.userRepository.InsertUser(models.User{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		logError("Error inserting user: %v", err)
		return "", err
	}

	// Generate and save verify token
	token, err := generateToken()
	if err != nil {
		// Rollback user insert on failure
		_ = s.userRepository.DeleteUserByID(dbUser.ID)
		return "", err
	}

	err = s.registrationRepository.SetVerifyToken(email, token)
	if err != nil {
		// Rollback user insert on failure
		_ = s.userRepository.DeleteUserByID(dbUser.ID)
		logError("Error setting verify token: %v", err)
		return "", err
	}

	// Track the token generation timestamp
	err = s.registrationRepository.TrackTokenGeneration(email)
	if err != nil {
		return "", fmt.Errorf("failed to track token generation: %w", err)
	}

	return token, nil
}

// VerifyToken validates the provided token for the given email
func (s *registrationService) VerifyToken(email, token string) error {
	result, err := s.registrationRepository.GetVerifyToken(email)
	if errors.Is(err, redis.Nil) {
		return errors.New("token does not exist")
	}
	if err != nil {
		logError("Error getting verify token: %v", err)
		return err
	}

	if result != token {
		return errors.New("token does not match")
	}

	// Delete the token after verify
	err = s.registrationRepository.DeleteToken(email)
	if err != nil {
		logError("Error deleting token: %v", err)
	}

	err = s.registrationRepository.DeleteVerificationCount(email)
	if err != nil {
		logError("Error deleting verification count: %v", err)
	}
	return nil
}

// SetUserVerified sets the user's verify status to true
func (s *registrationService) SetUserVerified(email string) error {
	user, err := s.userRepository.GetUserByMail(email)
	if err != nil {
		logError("Error getting user by mail: %v", err)
		return err
	}

	user.VerifiedAt = sql.NullTime{Time: time.Now(), Valid: true}
	user.IsVerified = true

	_, err = s.userRepository.UpdateUser(user)
	if err != nil {
		logError("Error updating user verify status: %v", err)
	}
	return err
}

// ReloadToken generates a new verify token and resets it
func (s *registrationService) ReloadToken(mail string) (string, error) {
	// Check if the user has already generated too many tokens in pas 24 hours
	canGenerate, err := s.registrationRepository.CanGenerateToken(mail)
	if err != nil {
		return "", fmt.Errorf("error checking token generation limit: %w", err)
	}
	if !canGenerate {
		return "", custom_error.NewTokenGenerationError("500", "you cannot generate more than 5 tokens in the past 24 hours")
	}
	// Generate a new token
	token, err := helpers.GenerateRandomString(64)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Update the token in Redis
	err = s.registrationRepository.SetVerifyToken(mail, token)
	if err != nil {
		return "", fmt.Errorf("failed to reset token: %w", err)
	}

	// Track the token generation timestamp
	err = s.registrationRepository.TrackTokenGeneration(mail)
	if err != nil {
		return "", fmt.Errorf("failed to track token generation: %w", err)
	}

	return token, nil
}

// Utility function to log errors
func logError(message string, err error) {
	if err != nil {
		log.Printf(message, err)
	}
}

// Utility function to handle token generation
func generateToken() (string, error) {
	token, err := helpers.GenerateRandomString(64)
	if err != nil {
		logError("Error generating token: %v", err)
	}
	return token, err
}

var UserRegistrationServiceProvider = wire.NewSet(NewUserRegistrationService)
