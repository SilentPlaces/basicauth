package repository

import (
	"context"
	consul "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type (
	RegistrationRepository interface {
		SetVerifyToken(mail string, token string) error
		GetVerifyToken(mail string) (string, error)
		DeleteToken(email string) error
	}

	registrationRepository struct {
		redisClient        *redis.Client
		tokenExpirySeconds time.Duration
	}
)

func NewRegistrationRepository(redisClient *redis.Client, consul consul.ConsulService) RegistrationRepository {
	cfg := consul.GetRegistrationConfig()

	return &registrationRepository{
		redisClient:        redisClient,
		tokenExpirySeconds: cfg.MailVerificationTimeInSeconds,
	}
}

func (v *registrationRepository) SetVerifyToken(mail string, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := v.redisClient.Set(ctx, mail, token, v.tokenExpirySeconds).Err()
	if err != nil {
		log.Printf("Failed to set registration token for mail '%s': %v", mail, err)
		return err
	}
	return nil
}

func (v *registrationRepository) GetVerifyToken(mail string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := v.redisClient.Get(ctx, mail).Result()
	if err != nil {
		log.Printf("Failed to get registration token for mail '%s': %v", mail, err)
		return "", err
	}
	return token, nil
}

func (v *registrationRepository) DeleteToken(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := v.redisClient.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
}

var RegistrationRepositoryProviderSet = wire.NewSet(NewRegistrationRepository)
