package repository

import (
	"context"
	"fmt"
	"github.com/SilentPlaces/basicauth.git/internal/config"
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
		TrackTokenGeneration(email string) error
		CanGenerateToken(email string) (bool, error)
		DeleteVerificationCount(mail string) error
	}

	registrationRepository struct {
		redisClient        *redis.Client
		registrationConfig *config.RegistrationConfig
	}
)

func NewRegistrationRepository(redisClient *redis.Client, consul consul.ConsulService) RegistrationRepository {
	cfg := consul.GetRegistrationConfig()

	return &registrationRepository{
		redisClient:        redisClient,
		registrationConfig: cfg,
	}
}

func (rp *registrationRepository) SetVerifyToken(mail string, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	key := prefixTokenKey + mail
	fmt.Printf("SetVerifyToken key=%s, %s \n", key, rp.registrationConfig.MailVerificationTimeInSeconds)
	err := rp.redisClient.Set(ctx, key, token, rp.registrationConfig.MailVerificationTimeInSeconds).Err()
	if err != nil {
		log.Printf("Failed to set registration token for mail '%s': %v", mail, err)
		return err
	}
	return nil
}

func (rp *registrationRepository) GetVerifyToken(mail string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := rp.redisClient.Get(ctx, prefixTokenKey+mail).Result()
	if err != nil {
		log.Printf("Failed to get registration token for mail '%s': %v", mail, err)
		return "", err
	}
	return token, nil
}

func (rp *registrationRepository) DeleteToken(mail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := rp.redisClient.Del(ctx, prefixTokenKey+mail).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rp *registrationRepository) DeleteVerificationCount(mail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := rp.redisClient.Del(ctx, prefixTokenKey+mail).Err()
	if err != nil {
		return err
	}
	return nil
}

// TrackTokenGeneration stores the timestamp of a token generation attempt.
func (rp *registrationRepository) TrackTokenGeneration(email string) error {
	now := time.Now().Unix()
	dayInSeconds := int64((time.Hour * 24).Seconds())

	// Use a sorted set to store the generation time of tokens
	_, err := rp.redisClient.ZAdd(context.Background(), prefixVerifyCountKey+email, redis.Z{
		Score:  float64(now),
		Member: now,
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to track token generation: %w", err)
	}

	// Remove any timestamps older than 24 hours (86400 seconds).
	_, err = rp.redisClient.ZRemRangeByScore(context.Background(), prefixVerifyCountKey+email, "0", fmt.Sprintf("%d", now-dayInSeconds)).Result()
	if err != nil {
		return fmt.Errorf("failed to clean up old token generations: %w", err)
	}

	return nil
}

// CanGenerateToken checks if a user has generated more than registrationConfig.MaxVerificationMailGenerationInHours tokens in the past 24 hours.
func (rp *registrationRepository) CanGenerateToken(email string) (bool, error) {
	now := time.Now().Unix()
	dayInSeconds := int64((time.Hour * 24).Seconds())
	// Count how many tokens have been generated in the last 24 hours
	count, err := rp.redisClient.ZCount(context.Background(), prefixVerifyCountKey+email, fmt.Sprintf("%d", now-dayInSeconds), fmt.Sprintf("%d", now)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to count token generations: %w", err)
	}

	// Allow generation only if there are fewer than registrationConfig.MaxVerificationMailGenerationInHours token generations in the last 24 hours.
	if count >= rp.registrationConfig.MaxVerificationMailGenerationInHours {
		return false, nil
	}
	return true, nil
}

const (
	prefixTokenKey       = "token-"
	prefixVerifyCountKey = "resend_verification-count-"
)

var RegistrationRepositoryProviderSet = wire.NewSet(NewRegistrationRepository)
