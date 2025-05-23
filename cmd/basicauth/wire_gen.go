// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"database/sql"
	"github.com/SilentPlaces/basicauth.git/internal/config"
	controller2 "github.com/SilentPlaces/basicauth.git/internal/controllers/registration"
	"github.com/SilentPlaces/basicauth.git/internal/controllers/user"
	"github.com/SilentPlaces/basicauth.git/internal/db/mysql"
	redis2 "github.com/SilentPlaces/basicauth.git/internal/db/redis"
	"github.com/SilentPlaces/basicauth.git/internal/repositories/registration"
	"github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	service2 "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	"github.com/SilentPlaces/basicauth.git/internal/services/consul"
	service5 "github.com/SilentPlaces/basicauth.git/internal/services/mail"
	service6 "github.com/SilentPlaces/basicauth.git/internal/services/registration"
	service4 "github.com/SilentPlaces/basicauth.git/internal/services/users"
	service3 "github.com/SilentPlaces/basicauth.git/internal/services/vault"
	"github.com/redis/go-redis/v9"
)

// Injectors from wire.go:

// InitializeConsulService initializes a ConsulService using dependencies from config and services packages.
func InitializeConsulService() service.ConsulService {
	appConfig := config.LoadConsulConfig()
	consulService := service.NewConsulService(appConfig)
	return consulService
}

// InitializeAuthService initializes an AuthService using dependencies from auth and vault packages.
func InitializeAuthService() service2.AuthService {
	secureVaultService := service3.NewSecureVaultService()
	authService := service2.NewAuthService(secureVaultService)
	return authService
}

// InitializeMySQLDB initializes a MySQL database connection.
func InitializeMySQLDB() (*sql.DB, error) {
	consulService := InitializeConsulService()
	db, err := mysql.NewMySQLDb(consulService)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InitializeRedis initializes a Redis client.
func InitializeRedis() (*redis.Client, error) {
	consulService := InitializeConsulService()
	client, err := redis2.NewRedis(consulService)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// InitializeUserService initializes a UserService.
func InitializeUserService() (service4.UserService, error) {
	db, err := InitializeMySQLDB()
	if err != nil {
		return nil, err
	}
	consulService := InitializeConsulService()
	userRepository := user.NewUserRepository(db, consulService)
	userService := service4.NewUserService(userRepository)
	return userService, nil
}

// InitializeUserController initializes a UserController.
func InitializeUserController() (controller.UserController, error) {
	userService, err := InitializeUserService()
	if err != nil {
		return nil, err
	}
	authService := InitializeAuthService()
	userController := controller.NewUserController(userService, authService)
	return userController, nil
}

// InitializeMailService initializes a MailService.
func InitializeMailService() (service5.MailService, error) {
	consulService := InitializeConsulService()
	appConfig := config.LoadConsulConfig()
	mailService, err := service5.NewMailService(consulService, appConfig)
	if err != nil {
		return nil, err
	}
	return mailService, nil
}

// RegistrationRepository initializes a RegistrationRepository.
func RegistrationRepository() (repository.RegistrationRepository, error) {
	client, err := InitializeRedis()
	if err != nil {
		return nil, err
	}
	consulService := InitializeConsulService()
	registrationRepository := repository.NewRegistrationRepository(client, consulService)
	return registrationRepository, nil
}

func InitializeUserRepository() (user.UserRepository, error) {
	db, err := InitializeMySQLDB()
	if err != nil {
		return nil, err
	}
	consulService := InitializeConsulService()
	userRepository := user.NewUserRepository(db, consulService)
	return userRepository, nil
}

// InitializeRegistrationService initializes a RegistrationService.
func InitializeRegistrationService() (service6.RegistrationService, error) {
	registrationRepository, err := RegistrationRepository()
	if err != nil {
		return nil, err
	}
	userRepository, err := InitializeUserRepository()
	if err != nil {
		return nil, err
	}
	registrationService := service6.NewUserRegistrationService(registrationRepository, userRepository)
	return registrationService, nil
}

// InitializeRegistrationController initializes a RegistrationController.
func InitializeRegistrationController() (controller2.RegistrationController, error) {
	mailService, err := InitializeMailService()
	if err != nil {
		return nil, err
	}
	registrationService, err := InitializeRegistrationService()
	if err != nil {
		return nil, err
	}
	userService, err := InitializeUserService()
	if err != nil {
		return nil, err
	}
	consulService := InitializeConsulService()
	registrationController := controller2.NewRegistrationController(mailService, registrationService, userService, consulService)
	return registrationController, nil
}
