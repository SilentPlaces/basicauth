package service

import (
	"fmt"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/convertor"
	"log"
	"sync"
	"time"

	"github.com/SilentPlaces/basicauth.git/internal/config"
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	"github.com/google/wire"
	consulapi "github.com/hashicorp/consul/api"
)

// ConsulService interface updated to return dedicated config structs.
type ConsulService interface {
	getConfigValue(key string) (string, error)
	getConfigForKeys(keys []string) (map[string]string, error)
	GetMySQLConfig() (*config.MySQLConfig, error)
	GetRedisConfig() (*config.RedisConfig, error)
	GetSMTPConfig() (*config.SMTPConfig, error)
	GetGeneralConfig() (*config.GeneralConfig, error)
	GetRegistrationConfig() *config.RegistrationConfig
	GetRegistrationPasswordConfig() *config.RegistrationPasswordConfig
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

// getConfigForKeys retrieves multiple configuration values.
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
func (cs *consulService) GetMySQLConfig() (*config.MySQLConfig, error) {
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
	configMap, err := cs.getConfigForKeys(keys)
	if err != nil {
		return nil, err
	}
	cfg := &config.MySQLConfig{
		Host:               configMap[constants.MySQLHostKey],
		Port:               configMap[constants.MySQLPortKey],
		User:               configMap[constants.MySQLUserKey],
		Password:           configMap[constants.MySQLPasswordKey],
		DB:                 configMap[constants.MySQLDBKey],
		MaxLifetimeSeconds: configMap[constants.MySQLMaxLifetimeSecondsKey],
		MaxOpenConnections: configMap[constants.MySQLMaxOpenConnectionsKey],
		IdleConnections:    configMap[constants.MySQLIdleConnectionsKey],
	}
	return cfg, nil
}

// GetRedisConfig retrieves Redis configuration from Consul.
func (cs *consulService) GetRedisConfig() (*config.RedisConfig, error) {
	keys := []string{
		constants.RedisHostKey,
		constants.RedisPortKey,
		constants.RedisPasswordKey,
	}
	configMap, err := cs.getConfigForKeys(keys)
	if err != nil {
		return nil, err
	}
	cfg := &config.RedisConfig{
		Host:     configMap[constants.RedisHostKey],
		Port:     configMap[constants.RedisPortKey],
		Password: configMap[constants.RedisPasswordKey],
	}
	return cfg, nil
}

// GetSMTPConfig retrieves SMTP configuration from Consul.
func (cs *consulService) GetSMTPConfig() (*config.SMTPConfig, error) {
	keys := []string{
		constants.SMTPHostKey,
		constants.SMTPPortKey,
		constants.SMTPUsernameKey,
		constants.SMTPPasswordKey,
	}
	configMap, err := cs.getConfigForKeys(keys)
	if err != nil {
		return nil, err
	}
	cfg := &config.SMTPConfig{
		Host:     configMap[constants.SMTPHostKey],
		Port:     configMap[constants.SMTPPortKey],
		Username: configMap[constants.SMTPUsernameKey],
		Password: configMap[constants.SMTPPasswordKey],
	}
	return cfg, nil
}

// GetGeneralConfig retrieves general application configuration from Consul.
func (cs *consulService) GetGeneralConfig() (*config.GeneralConfig, error) {
	keys := []string{
		constants.GeneralDomainKey,
		constants.GeneralHTTPListenerPortKey,
	}
	configMap, err := cs.getConfigForKeys(keys)
	if err != nil {
		return nil, err
	}
	cfg := &config.GeneralConfig{
		Domain:           configMap[constants.GeneralDomainKey],
		HTTPListenerPort: configMap[constants.GeneralHTTPListenerPortKey],
	}
	return cfg, nil
}

// GetRegistrationConfig retrieves registration configuration from Consul.
func (cs *consulService) GetRegistrationConfig() *config.RegistrationConfig {
	keys := []string{
		constants.GeneralRegisterMailVerificationTimeInSecondsKey,
		constants.GeneralRegisterHostVerificationMailAddressKey,
		constants.GeneralRegisterVerificationMailTextKey,
	}
	defaultTokenExpirySeconds := 600
	tokenExpirySeconds := defaultTokenExpirySeconds

	configMap, err := cs.getConfigForKeys(keys)
	if err == nil {
		if val, err := helpers.ParseInt("token expireTime", configMap[constants.GeneralRegisterMailVerificationTimeInSecondsKey]); err == nil {
			tokenExpirySeconds = val
		}
	}

	return &config.RegistrationConfig{
		MailVerificationTimeInSeconds: time.Duration(tokenExpirySeconds),
		HostVerificationMailAddress:   configMap[constants.GeneralRegisterHostVerificationMailAddressKey],
		VerificationMailText:          configMap[constants.GeneralRegisterVerificationMailTextKey],
	}
}

// GetRegistrationPasswordConfig retrieves and parses the registration password configuration from Consul.
func (cs *consulService) GetRegistrationPasswordConfig() *config.RegistrationPasswordConfig {
	keys := []string{
		constants.KeyRegistrationPasswordMinLength,
		constants.KeyRegistrationPasswordRequireUpper,
		constants.KeyRegistrationPasswordRequireLower,
		constants.KeyRegistrationPasswordRequireNumber,
		constants.KeyRegistrationPasswordRequireSpecial,
	}

	configMap, err := cs.getConfigForKeys(keys)
	if err != nil {
		return nil
	}

	minLength, err := helpers.ParseInt("minLength", configMap[constants.KeyRegistrationPasswordMinLength])
	if err != nil {
		minLength = 8
	}

	parseBool := func(key, fieldName string) bool {
		if v, err := helpers.ParseBool(fieldName, configMap[key]); err == nil {
			return v
		}
		return false
	}

	return &config.RegistrationPasswordConfig{
		MinLength:      minLength,
		RequireUpper:   parseBool(constants.KeyRegistrationPasswordRequireUpper, "requireUpper"),
		RequireLower:   parseBool(constants.KeyRegistrationPasswordRequireLower, "requireLower"),
		RequireNumber:  parseBool(constants.KeyRegistrationPasswordRequireNumber, "requireNumber"),
		RequireSpecial: parseBool(constants.KeyRegistrationPasswordRequireSpecial, "requireSpecial"),
	}
}

var ConsulProviderSet = wire.NewSet(NewConsulService)
