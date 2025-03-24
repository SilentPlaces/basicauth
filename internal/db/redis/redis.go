package redis

import (
	"context"
	"fmt"
	consulService "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedis(consul consulService.ConsulService) (*redis.Client, error) {
	cfg, err := consul.GetRedisConfig()
	if err != nil {
		return nil, err
	}

	host := cfg.Host
	port := cfg.Port
	password := cfg.Password

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

var RedisProviderSet = wire.NewSet(NewRedis)
