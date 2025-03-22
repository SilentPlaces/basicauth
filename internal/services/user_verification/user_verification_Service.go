package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	consul "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type (
	UserVerificationService interface {
		GenerateVerificationToken(email string) (string, error)
		VerifyToken(email, token string) error
	}

	userVerificationService struct {
		redis              *redis.Client
		tokenExpirySeconds time.Duration
	}
)

func NewUserVerificationService(redis *redis.Client, consul consul.ConsulService) UserVerificationService {
	cfg, err := consul.GetRegistrationConfig()
	var tokenExpirySeconds int
	if err != nil {
		tokenExpirySeconds = 6000
	} else {
		tokenExpirySeconds, err = helpers.ParseInt("token expireTime", cfg[constants.GeneralRegisterMailVerificationTimeInSecondsKey])
		if err != nil {
			tokenExpirySeconds = 6000
		}
	}

	return &userVerificationService{
		redis:              redis,
		tokenExpirySeconds: time.Duration(tokenExpirySeconds) * time.Second,
	}
}

func (s *userVerificationService) GenerateVerificationToken(email string) (string, error) {
	tokenBytes := make([]byte, 128)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(tokenBytes)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.redis.Set(ctx, fmt.Sprintf("verify:%s", email), token, s.tokenExpirySeconds).Err()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return token, nil
}

func (s *userVerificationService) VerifyToken(email, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := s.redis.Get(ctx, fmt.Sprintf("verify:%s", email)).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("token does not exist")
	} else if err != nil {
		return err
	}

	if result != token {
		return errors.New("token does not match")
	}

	s.redis.Del(ctx, fmt.Sprintf("verify:%s", email))
	return nil
}

var UserVerificationServiceProvider = wire.NewSet(NewUserVerificationService)
