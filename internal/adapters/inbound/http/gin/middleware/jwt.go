package middleware

import (
	"net/http"
	"strings"

	"github.com/SilentPlaces/basicauth.git/internal/adapters/inbound/http/gin/response"
	"github.com/SilentPlaces/basicauth.git/internal/application/port"
	appLogger "github.com/SilentPlaces/basicauth.git/internal/shared/logger"
	"github.com/gin-gonic/gin"
)

const UserContextKey = "user"

func JWTAuthMiddleware(authService port.AuthTokenManager, logger appLogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn(c.Request.Context(), "jwt middleware missing authorization header", map[string]interface{}{"path": c.Request.URL.Path})
			response.Error(c, http.StatusUnauthorized, "Unauthorized request!")
			c.Abort()
			return
		}

		split := strings.Split(authHeader, " ")
		if len(split) != 2 || split[0] != "Bearer" {
			logger.Warn(c.Request.Context(), "jwt middleware invalid authorization format", map[string]interface{}{"path": c.Request.URL.Path})
			response.Error(c, http.StatusUnauthorized, "Invalid Authorization header format.")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(split[1])
		if err := authService.ValidateToken(tokenString); err != nil {
			logger.Warn(c.Request.Context(), "jwt middleware token validation failed", map[string]interface{}{"path": c.Request.URL.Path})
			response.Error(c, http.StatusUnauthorized, "Invalid Token")
			c.Abort()
			return
		}

		claims, err := authService.ExtractClaims(tokenString)
		if err != nil {
			logger.Warn(c.Request.Context(), "jwt middleware claims extraction failed", map[string]interface{}{"path": c.Request.URL.Path})
			response.Error(c, http.StatusUnauthorized, "Invalid Token")
			c.Abort()
			return
		}

		logger.Debug(c.Request.Context(), "jwt middleware authenticated request", map[string]interface{}{"user_id": claims.UserID, "path": c.Request.URL.Path})
		c.Set(UserContextKey, claims.UserID)
		c.Next()
	}
}
