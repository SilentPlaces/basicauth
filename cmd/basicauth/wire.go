//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	config "github.com/SilentPlaces/basicauth.git/internal/config"
	registrationController "github.com/SilentPlaces/basicauth.git/internal/controllers/registration"
	userController "github.com/SilentPlaces/basicauth.git/internal/controllers/user"
	mysql "github.com/SilentPlaces/basicauth.git/internal/db/mysql"
	redisProvider "github.com/SilentPlaces/basicauth.git/internal/db/redis"
	registrationRepository "github.com/SilentPlaces/basicauth.git/internal/repositories/registration"
	userRepository "github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	auth "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	consul "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	mailService "github.com/SilentPlaces/basicauth.git/internal/services/mail"
	registerationService "github.com/SilentPlaces/basicauth.git/internal/services/registration"
	userService "github.com/SilentPlaces/basicauth.git/internal/services/users"
	vault "github.com/SilentPlaces/basicauth.git/internal/services/vault"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// InitializeConsulService initializes a ConsulService using dependencies from config and services packages.
func InitializeConsulService() consul.ConsulService {
	wire.Build(config.ProviderSet, consul.ConsulProviderSet)
	return nil
}

// InitializeAuthService initializes an AuthService using dependencies from auth and vault packages.
func InitializeAuthService() auth.AuthService {
	wire.Build(auth.AuthServiceProviderSet, vault.VaultServiceProviderSet)
	return nil
}

// InitializeMySQLDB initializes a MySQL database connection.
func InitializeMySQLDB() (*sql.DB, error) {
	wire.Build(
		mysql.MySqlProviderSet,
		InitializeConsulService, // InitializeConsulService is called here
	)
	return nil, nil
}

// InitializeRedis initializes a Redis client.
func InitializeRedis() (*redis.Client, error) {
	wire.Build(
		redisProvider.RedisProviderSet,
		InitializeConsulService, // InitializeConsulService is called here
	)
	return nil, nil
}

// InitializeUserService initializes a UserService.
func InitializeUserService() (userService.UserService, error) {
	wire.Build(
		userRepository.UserRepositoryProviderSet,
		userService.UserServiceProviderSet,
		InitializeMySQLDB,
		InitializeConsulService,
	)
	return nil, nil
}

// InitializeUserController initializes a UserController.
func InitializeUserController() (userController.UserController, error) {
	wire.Build(
		InitializeUserService,
		userController.UserControllerProviderSet,
	)
	return nil, nil
}

// InitializeMailService initializes a MailService.
func InitializeMailService() (mailService.MailService, error) {
	wire.Build(mailService.MailServiceProviderSet,
		InitializeConsulService,
		config.ProviderSet)
	return nil, nil
}

// RegistrationRepository initializes a RegistrationRepository.
func RegistrationRepository() (registrationRepository.RegistrationRepository, error) {
	wire.Build(
		registrationRepository.RegistrationRepositoryProviderSet,
		InitializeRedis,
		InitializeConsulService,
	)
	return nil, nil
}

func InitializeUserRepository() (userRepository.UserRepository, error) {
	wire.Build(userRepository.UserRepositoryProviderSet, InitializeMySQLDB, InitializeConsulService)
	return nil, nil
}

// InitializeRegistrationService initializes a RegistrationService.
func InitializeRegistrationService() (registerationService.RegistrationService, error) {
	wire.Build(
		registerationService.UserRegistrationServiceProvider,
		RegistrationRepository,
		InitializeUserRepository,
	)
	return nil, nil
}

// InitializeRegistrationController initializes a RegistrationController.
func InitializeRegistrationController() (registrationController.RegistrationController, error) {
	wire.Build(
		registrationController.RegistrationControllerProvider,
		InitializeMailService,
		InitializeRegistrationService,
		InitializeUserService,
		InitializeConsulService,
	)
	return nil, nil
}
