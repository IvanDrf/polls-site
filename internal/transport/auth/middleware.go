package auth

import (
	"net/http"
	"strings"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type Middleware interface {
	AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	jwter jwt.Jwter
}

func NewMiddleware(cfg *config.Config) Middleware {
	return middleware{jwter: jwt.NewJwter(cfg)}
}

func (middle middleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "auth header is required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "invalid auth format", http.StatusUnauthorized)
			return
		}

		if err := middle.jwter.IsValidJWT(tokenParts[1]); err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
