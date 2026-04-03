package handlers

import (
	"net/http"

	"github.com/SilentPlaces/basicauth.git/internal/infrastructure/health"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	checker health.Checker
	logger  appLogger.Logger
}

func NewHealthHandler(checker health.Checker, logger appLogger.Logger) *HealthHandler {
	return &HealthHandler{checker: checker, logger: logger}
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	h.logger.Debug(c.Request.Context(), "liveness probe", nil)
	c.JSON(http.StatusOK, h.checker.Liveness())
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	status, ready := h.checker.Readiness(c.Request.Context())
	logFields := map[string]interface{}{
		"status": status["status"],
		"mysql":  status["mysql"],
		"redis":  status["redis"],
	}
	if !ready {
		h.logger.Warn(c.Request.Context(), "readiness probe degraded", logFields)
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	h.logger.Debug(c.Request.Context(), "readiness probe ok", logFields)
	c.JSON(http.StatusOK, status)
}
