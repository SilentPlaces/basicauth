package middleware

import (
	"github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/SilentPlaces/basicauth.git/internal/shared/observability"
	"github.com/gin-gonic/gin"
)

const TraceParentHeader = "traceparent"

// OTelReadyMiddleware is a lightweight skeleton for future OpenTelemetry integration.
// It preserves incoming W3C trace context and puts it into request context for logs.
func OTelReadyMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceParent := c.GetHeader(TraceParentHeader)
		if traceParent != "" {
			ctx := observability.WithTraceParent(c.Request.Context(), traceParent)
			c.Request = c.Request.WithContext(ctx)
			log.Debug(ctx, "traceparent detected", map[string]interface{}{"path": c.Request.URL.Path})
		}

		// Future hook:
		// - Start OTel span from context here
		// - Add attributes (method, route, status)
		// - End span after c.Next()
		c.Next()
	}
}
