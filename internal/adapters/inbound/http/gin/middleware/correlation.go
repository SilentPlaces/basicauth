package middleware

import (
	"github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/SilentPlaces/basicauth.git/internal/shared/observability"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CorrelationIDHeader = "X-Correlation-ID"

func CorrelationIDMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		ctx := observability.WithCorrelationID(c.Request.Context(), correlationID)
		c.Request = c.Request.WithContext(ctx)
		c.Header(CorrelationIDHeader, correlationID)

		log.Debug(ctx, "correlation id assigned", map[string]interface{}{
			"path":           c.Request.URL.Path,
			"correlation_id": correlationID,
		})

		c.Next()
	}
}
