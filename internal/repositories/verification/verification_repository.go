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
	VerificationRepository interface {
		SetVerifyToken(mail string, token string) error
		GetVerifyToken(mail string) (string, error)
		DeleteToken(email string) error
	}

	verificationRepository struct {
		redisClient        *redis.Client
		tokenExpirySeconds time.Duration
	}
)

func NewVerificationRepository(redisClient *redis.Client, consul consul.ConsulService) VerificationRepository {
	cfg, err := consul.GetRegistrationConfig()
	var expireTime = 24 * time.Hour
	if err == nil {
		expireTime = cfg.MailVerificationTimeInSeconds
	}

	return &verificationRepository{
		redisClient:        redisClient,
		tokenExpirySeconds: expireTime,
	}
}

func (v *verificationRepository) SetVerifyToken(mail string, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := v.redisClient.Set(ctx, mail, token, v.tokenExpirySeconds).Err()
	if err != nil {
		log.Printf("Failed to set verification token for mail '%s': %v", mail, err)
		return err
	}
	return nil
}

func (v *verificationRepository) GetVerifyToken(mail string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := v.redisClient.Get(ctx, mail).Result()
	if err != nil {
		log.Printf("Failed to get verification token for mail '%s': %v", mail, err)
		return "", err
	}
	return token, nil
}

func (v *verificationRepository) DeleteToken(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := v.redisClient.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
}

var VerificationRepositorySet = wire.NewSet(NewVerificationRepository)
