package service

import (
	"fmt"
	"github.com/SilentPlaces/basicauth.git/internal/config"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	"log"
	"sync"

	"github.com/google/wire"
	consulapi "github.com/hashicorp/consul/api"
)

type ConsulService interface {
	getConfigValue(key string) (string, error)
	getConfigForKeys(keys []string) (map[string]string, error)
	GetMySQLConfig() (map[string]string, error)
	GetRedisConfig() (map[string]string, error)
}

type consulService struct {
	Client *consulapi.Client
}

var (
	service *consulService
	once    sync.Once
)

// NewConsulService creates a singleton ConsulService.
// It panics on error to ensure the application does not continue if the connection fails.
func NewConsulService(cfg *config.AppConfig) ConsulService {
	once.Do(func() {
		consulConfig := consulapi.DefaultConfig()
		consulConfig.Address = cfg.ConsulAddress
		consulConfig.Scheme = cfg.ConsulScheme

		client, err := consulapi.NewClient(consulConfig)
		if err != nil {
			log.Panicf("Error connecting to Consul server: %v", err)
		}

		service = &consulService{Client: client}
	})
	return service
}

// getConfigValue retrieves a single key from Consul's KV store.
func (cs *consulService) getConfigValue(key string) (string, error) {
	kv := cs.Client.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return "", fmt.Errorf("error retrieving key %s: %v", key, err)
	}
	if pair == nil {
		return "", fmt.Errorf("key %s not found", key)
	}
	return string(pair.Value), nil
}

// getConfigForKeys is a function to retrieve multiple configuration values.
func (cs *consulService) getConfigForKeys(keys []string) (map[string]string, error) {
	configMap := make(map[string]string)
	for _, key := range keys {
		value, err := cs.getConfigValue(key)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch config for key %s: %v", key, err)
		}
		configMap[key] = value
	}
	return configMap, nil
}

// GetMySQLConfig retrieves MySQL configuration from Consul.
func (cs *consulService) GetMySQLConfig() (map[string]string, error) {
	keys := []string{
		constants.MySQLHostKey,
		constants.MySQLPortKey,
		constants.MySQLUserKey,
		constants.MySQLPasswordKey,
		constants.MySQLDBKey,
		constants.MySQLMaxLifetimeSecondsKey,
		constants.MySQLMaxOpenConnectionsKey,
		constants.MySQLIdleConnectionsKey,
	}
	return cs.getConfigForKeys(keys)
}

// GetRedisConfig retrieves Redis configuration from Consul.
func (cs *consulService) GetRedisConfig() (map[string]string, error) {
	keys := []string{
		constants.RedisHostKey,
		constants.RedisPortKey,
		constants.RedisPasswordKey,
	}
	return cs.getConfigForKeys(keys)
}

var ConsulProviderSet = wire.NewSet(NewConsulService)
