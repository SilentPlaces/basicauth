package config

import (
	"github.com/SilentPlaces/basicauth.git/pkg/constants"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/joho/godotenv"
)

// AppConfig. it only holds consule connection info. other configs are placed in Consul
type AppConfig struct {
	ConsulAddress string
	ConsulScheme  string
	Environment   string
}

// Configuration Structs
type MySQLConfig struct {
	Host               string
	Port               string
	User               string
	Password           string
	DB                 string
	MaxLifetimeSeconds string
	MaxOpenConnections string
	IdleConnections    string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type GeneralConfig struct {
	Domain           string
	HTTPListenerPort string
}

type RegistrationPasswordConfig struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

type RegistrationConfig struct {
	MailVerificationTimeInSeconds time.Duration
	HostVerificationMailAddress   string
	VerificationMailText          string
}

var (
	appConfig *AppConfig
	once      sync.Once
)

// LoadConsulConfig loads consul configuration from the specified .env file. It is singleton.
func LoadConsulConfig() *AppConfig {
	once.Do(func() {
		if err := godotenv.Load(constants.EnvFile); err != nil {
			log.Panic("Error loading .env file:", err)
		}
		appConfig = &AppConfig{
			ConsulAddress: os.Getenv(constants.EnvKeyConsulAddress),
			ConsulScheme:  os.Getenv(constants.EnvKeyConsulScheme),
			Environment:   os.Getenv(constants.EnvKeyAppEnvironment),
		}

		if appConfig.ConsulAddress == "" || appConfig.ConsulScheme == "" {
			log.Panic("Enter Consul Configuration in .env file to continue")
		}
	})
	return appConfig
}

// Dependency Injection
var ProviderSet = wire.NewSet(LoadConsulConfig)
