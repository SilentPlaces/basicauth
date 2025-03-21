//go:build wireinject
// +build wireinject

package main

import (
	"github.com/SilentPlaces/basicauth.git/internal/config"
	"github.com/SilentPlaces/basicauth.git/internal/services/consul"
	"github.com/google/wire"
)

// InitializeConsulService initializes a ConsulService using dependencies from config and services packages.
func InitializeConsulService() *consul.ConsulService {
	wire.Build(config.ProviderSet, consul.ProviderSet)
	return nil
}
