//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	config "github.com/SilentPlaces/basicauth.git/internal/config"
	userController "github.com/SilentPlaces/basicauth.git/internal/controllers/user"
	"github.com/SilentPlaces/basicauth.git/internal/db/mysql"
	userRepository "github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	auth "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	consul "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	userService "github.com/SilentPlaces/basicauth.git/internal/services/users"
	vault "github.com/SilentPlaces/basicauth.git/internal/services/vault"
	"github.com/google/wire"
)

// InitializeConsulService initializes a ConsulService using dependencies from config and services packages.
func InitializeConsulService() consul.ConsulService {
	wire.Build(config.ProviderSet, consul.ConsulProviderSet)
	return nil
}

func InitializeAuthService() auth.AuthService {
	wire.Build(auth.AuthServiceProviderSet, vault.VaultServiceProviderSet)
	return nil
}

func InitializeMySQLDB() (*sql.DB, error) {
	wire.Build(
		InitializeConsulService,
		mysql.MySqlProviderSet,
	)
	return nil, nil
}

func InitializeUserService() (userService.UserService, error) {
	wire.Build(
		userRepository.UserRepositoryProviderSet,
		userService.UserServiceProviderSet,
		InitializeMySQLDB,
	)
	return nil, nil
}

func InitializeUserController() (userController.UserController, error) {
	wire.Build(
		InitializeUserService,
		userController.UserControllerProviderSet,
	)
	return nil, nil
}
