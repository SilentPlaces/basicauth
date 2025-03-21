package consul

import (
	"fmt"
	"github.com/SilentPlaces/basicauth.git/internal/config"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	"log"
	"sync"

	"github.com/google/wire"
	consulapi "github.com/hashicorp/consul/api"
)

type ConsulService struct {
	Client *consulapi.Client
}

var (
	consulService *ConsulService
	once          sync.Once
)

// NewConsulService creates a singleton ConsulService.
// It panics on error to ensure the application does not continue if the connection fails.
func NewConsulService(cfg *config.AppConfig) *ConsulService {
	once.Do(func() {
		consulConfig := consulapi.DefaultConfig()
		consulConfig.Address = cfg.ConsulAddress
		consulConfig.Scheme = cfg.ConsulScheme

		client, err := consulapi.NewClient(consulConfig)
		if err != nil {
			log.Panicf("Error connecting to Consul server: %v", err)
		}

		consulService = &ConsulService{Client: client}
	})
	return consulService
}

// getConfigValue retrieves a single key from Consul's KV store.
func (cs *ConsulService) getConfigValue(key string) (string, error) {
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
func (cs *ConsulService) getConfigForKeys(keys []string) (map[string]string, error) {
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
func (cs *ConsulService) GetMySQLConfig() (map[string]string, error) {
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
func (cs *ConsulService) GetRedisConfig() (map[string]string, error) {
	keys := []string{
		constants.RedisHostKey,
		constants.RedisPortKey,
		constants.RedisPasswordKey,
	}
	return cs.getConfigForKeys(keys)
}

var ProviderSet = wire.NewSet(NewConsulService)
