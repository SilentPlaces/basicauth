package router

import (
	"github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/handlers"
	"github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/middleware"
	"github.com/SilentPlaces/basicauth.git/internal/application/port"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(
	userHandler *handlers.UserHandler,
	registrationHandler *handlers.RegistrationHandler,
	healthHandler *handlers.HealthHandler,
	authService port.AuthTokenManager,
	logger appLogger.Logger,
) *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.CorrelationIDMiddleware(logger))
	engine.Use(middleware.OTelReadyMiddleware(logger))
	engine.Use(middleware.PrometheusMetricsMiddleware())

	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.GET("/health/live", healthHandler.Liveness)
	engine.GET("/health/ready", healthHandler.Readiness)

	engine.POST("/auth/login", userHandler.Login)
	engine.POST("/auth/refresh-token", userHandler.RefreshToken)

	engine.POST("/register/init", registrationHandler.SignUp)
	engine.POST("/register/verify", registrationHandler.VerifyMail)
	engine.POST("/register/resend-verification", registrationHandler.ResendVerification)

	protected := engine.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(authService, logger))
	protected.GET("/user", userHandler.GetUser)

	return engine
}
