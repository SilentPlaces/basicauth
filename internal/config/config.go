package config

import (
	constants2 "github.com/SilentPlaces/basicauth.git/pkg/constants"
	"log"
	"os"
	"sync"

	"github.com/google/wire"
	"github.com/joho/godotenv"
)

// AppConfig. it only holds consule connection info. other configs are placed in Consul
type AppConfig struct {
	ConsulAddress string
	ConsulScheme  string
}

var (
	appConfig *AppConfig
	once      sync.Once
)

// LoadConfig loads configuration from the specified .env file. It is singleton.
func LoadConfig() *AppConfig {
	once.Do(func() {
		if err := godotenv.Load(constants2.EnvFile); err != nil {
			log.Panic("Error loading .env file:", err)
		}
		appConfig = &AppConfig{
			ConsulAddress: os.Getenv(constants2.EnvKeyConsulAddress),
			ConsulScheme:  os.Getenv(constants2.EnvKeyConsulScheme),
		}

		if appConfig.ConsulAddress == "" || appConfig.ConsulScheme == "" {
			log.Panic("Enter Consul Configuration in .env file to continue")
		}
	})
	return appConfig
}

// Dependency Injection
var ProviderSet = wire.NewSet(LoadConfig)
