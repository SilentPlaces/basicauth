package di

import (
	"context"

	"github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/handlers"
	ginrouter "github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/router"
	"github.com/SilentPlaces/basicauth.git/internal/application/usecase"
	"github.com/SilentPlaces/basicauth.git/internal/config"
	"github.com/SilentPlaces/basicauth.git/internal/db/mysql"
	redisprovider "github.com/SilentPlaces/basicauth.git/internal/db/redis"
	healthinfra "github.com/SilentPlaces/basicauth.git/internal/infrastructure/health"
	"github.com/SilentPlaces/basicauth.git/internal/infrastructure/logging"
	registrationrepo "github.com/SilentPlaces/basicauth.git/internal/repositories/registration"
	userrepo "github.com/SilentPlaces/basicauth.git/internal/repositories/user"
	authservice "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	consulservice "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	mailservice "github.com/SilentPlaces/basicauth.git/internal/services/mail"
	registrationservice "github.com/SilentPlaces/basicauth.git/internal/services/registration"
	userservice "github.com/SilentPlaces/basicauth.git/internal/services/users"
	vaultservice "github.com/SilentPlaces/basicauth.git/internal/services/vault"
	"github.com/gin-gonic/gin"
)

type Container struct {
	Router *gin.Engine
}

func BuildContainer() (*Container, error) {
	appCfg := config.LoadConsulConfig()
	logger := logging.NewZeroLogger(appCfg)
	logger.Info(context.Background(), "building dependency container", nil)

	consul := consulservice.NewConsulService(appCfg)

	mysqlDB, err := mysql.NewMySQLDb(consul)
	if err != nil {
		logger.Error(context.Background(), "mysql initialization failed", err, nil)
		return nil, err
	}

	redisClient, err := redisprovider.NewRedis(consul)
	if err != nil {
		logger.Error(context.Background(), "redis initialization failed", err, nil)
		return nil, err
	}

	userRepository := userrepo.NewUserRepository(mysqlDB, consul)
	registrationRepository := registrationrepo.NewRegistrationRepository(redisClient, consul)
	userService := userservice.NewUserService(userRepository)
	registrationService := registrationservice.NewUserRegistrationService(registrationRepository, userRepository)
	authService := authservice.NewAuthService(vaultservice.NewSecureVaultService())

	mailSvc, err := mailservice.NewMailService(consul, appCfg)
	if err != nil {
		logger.Error(context.Background(), "mail service initialization failed", err, nil)
		return nil, err
	}

	generalCfg, err := consul.GetGeneralConfig()
	if err != nil {
		logger.Error(context.Background(), "general config retrieval failed", err, nil)
		return nil, err
	}

	registrationUseCase := usecase.NewRegistrationUseCase(
		mailSvc,
		registrationService,
		consul.GetRegistrationConfig(),
		consul.GetRegistrationPasswordConfig(),
		generalCfg,
		logger,
	)
	userUseCase := usecase.NewUserUseCase(userService, logger)
	authUseCase := usecase.NewAuthUseCase(userService, authService, logger)

	userHandler := handlers.NewUserHandler(userUseCase, authUseCase, logger)
	registrationHandler := handlers.NewRegistrationHandler(registrationUseCase, logger)
	healthHandler := handlers.NewHealthHandler(healthinfra.NewChecker(mysqlDB, redisClient), logger)

	router := ginrouter.NewRouter(userHandler, registrationHandler, healthHandler, authService, logger)
	logger.Info(context.Background(), "dependency container built", nil)
	return &Container{Router: router}, nil
}
