package service

import (
	"errors"
	repository "github.com/SilentPlaces/basicauth.git/internal/repositories/verification"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"log"
)

type (
	UserVerificationService interface {
		GenerateVerificationToken(email string) (string, error)
		VerifyToken(email, token string) error
	}

	userVerificationService struct {
		verificationRepo repository.VerificationRepository
	}
)

func NewUserVerificationService(verificationRepo repository.VerificationRepository) UserVerificationService {

	return &userVerificationService{
		verificationRepo: verificationRepo,
	}
}

func (s *userVerificationService) GenerateVerificationToken(email string) (string, error) {
	//generate token
	token, err := helpers.GenerateRandomString(256)
	if err != nil {
		log.Printf("Error generating token: %s", err.Error())
		return "", err
	}
	//save token to redis
	err = s.verificationRepo.SetVerifyToken(email, token)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return token, nil
}

func (s *userVerificationService) VerifyToken(email, token string) error {
	result, err := s.verificationRepo.GetVerifyToken(email)
	if errors.Is(err, redis.Nil) {
		return errors.New("token does not exist")
	} else if err != nil {
		return err
	}

	if result != token {
		return errors.New("token does not match")
	}

	err = s.verificationRepo.DeleteToken(email)
	if err != nil {
		log.Printf("Error deleting token: %s", err.Error())
	}
	return nil
}

var UserVerificationServiceProvider = wire.NewSet(NewUserVerificationService)
