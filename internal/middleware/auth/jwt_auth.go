package middleware

import (
	auth "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type contextKey string

var UserContextKey contextKey = "user"

func JWTAuthMiddleware(service auth.AuthService, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized request!", http.StatusUnauthorized)
			return
		}

		split := strings.Split(authHeader, " ")
		if len(split) != 2 && split[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format.", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(split[1])
		err := service.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		claims, err := service.ExtractClaims(tokenString)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
		next(w, r.WithContext(ctx), ps)
	}
}
