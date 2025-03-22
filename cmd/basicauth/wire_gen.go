// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"database/sql"
	"github.com/SilentPlaces/basicauth.git/internal/config"
	"github.com/SilentPlaces/basicauth.git/internal/controllers/user"
	"github.com/SilentPlaces/basicauth.git/internal/db"
	"github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	service2 "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	"github.com/SilentPlaces/basicauth.git/internal/services/consul"
	service4 "github.com/SilentPlaces/basicauth.git/internal/services/users"
	service3 "github.com/SilentPlaces/basicauth.git/internal/services/vault"
)

// Injectors from wire.go:

// InitializeConsulService initializes a ConsulService using dependencies from config and services packages.
func InitializeConsulService() service.ConsulService {
	appConfig := config.LoadConfig()
	consulService := service.NewConsulService(appConfig)
	return consulService
}

func InitializeAuthService() service2.AuthService {
	secureVaultService := service3.NewSecureVaultService()
	authService := service2.NewAuthService(secureVaultService)
	return authService
}

func InitializeMySQLDB() (*sql.DB, error) {
	consulService := InitializeConsulService()
	sqlDB, err := db.NewMySQLDb(consulService)
	if err != nil {
		return nil, err
	}
	return sqlDB, nil
}

func InitializeUserService() (service4.UserService, error) {
	sqlDB, err := InitializeMySQLDB()
	if err != nil {
		return nil, err
	}
	userRepository := user.NewUserRepository(sqlDB)
	userService := service4.NewUserService(userRepository)
	return userService, nil
}

func InitializeUserController() (controller.UserController, error) {
	userService, err := InitializeUserService()
	if err != nil {
		return nil, err
	}
	userController := controller.NewUserController(userService)
	return userController, nil
}
