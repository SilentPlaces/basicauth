//go:build wireinject
// +build wireinject

package main

import (
	config "github.com/SilentPlaces/basicauth.git/internal/config"
	auth "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	consul "github.com/SilentPlaces/basicauth.git/internal/services/consul"
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
