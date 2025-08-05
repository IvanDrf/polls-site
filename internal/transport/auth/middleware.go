package auth

import (
	"log/slog"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"

	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type Middleware interface {
	AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	jwter jwter.Jwter

	logger *slog.Logger
}

func NewMiddleware(cfg *config.Config, logger *slog.Logger) Middleware {
	return middleware{
		jwter: jwter.NewJwter(cfg),

		logger: logger,
	}
}

func (middle middleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middle.logger.Info("middle -> start")

		accessToken, err := middle.jwter.GetToken(r, jwter.AccessToken)
		if err != nil {
			middle.logger.Info("middle -> can't get token")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := middle.jwter.IsValidJWT(accessToken); err == nil {
			middle.logger.Info("middle -> valid token")
			next(w, r)
			return
		}

		middle.logger.Info("middle -> invalid token")

		http.Error(w, errs.ErrInValidToken().Error(), http.StatusUnauthorized)

	}
}
