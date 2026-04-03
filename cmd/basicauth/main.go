package main

import (
	"context"

	"github.com/SilentPlaces/basicauth.git/internal/config"
	"github.com/SilentPlaces/basicauth.git/internal/infrastructure/di"
	"github.com/SilentPlaces/basicauth.git/internal/infrastructure/logging"
)

func main() {
	logger := logging.NewZeroLogger(config.LoadConsulConfig())
	container, err := di.BuildContainer()
	if err != nil {
		logger.Error(context.Background(), "failed to initialize application container", err, nil)
		return
	}

	logger.Info(context.Background(), "server starting", map[string]interface{}{"port": "8080"})
	if err := container.Router.Run(":8080"); err != nil {
		logger.Error(context.Background(), "server terminated", err, nil)
	}
}
