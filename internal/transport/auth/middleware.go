package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	t "github.com/IvanDrf/polls-site/internal/repo/auth/tokens"
	u "github.com/IvanDrf/polls-site/internal/repo/auth/user"

	"github.com/IvanDrf/polls-site/internal/transport/auth/cookies"
	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
	"github.com/golang-jwt/jwt"
)

type Middleware interface {
	AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	cookier cookies.Cookier

	jwter jwter.Jwter

	userRepo  u.UserRepo
	tokenRepo t.TokensRepo
}

func NewMiddleware(cfg *config.Config, db *sql.DB) Middleware {
	return middleware{
		cookier:   cookies.NewCookier(),
		jwter:     jwter.NewJwter(cfg),
		userRepo:  u.NewRepo(cfg, db),
		tokenRepo: t.NewTokensRepo(cfg, db),
	}
}

func (middle middleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := middle.jwter.GetToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err := middle.jwter.IsValidJWT(accessToken); err == nil {
			next(w, r)
			return
		}

		if errors.Is(err, errs.ErrInValidToken()) {
			cookie, err := r.Cookie(jwter.RefreshToken)
			if err != nil {
				http.Error(w, errs.ErrCantFindToken().Error(), http.StatusUnauthorized)
				return
			}

			refreshToken := cookie.Value
			if err := middle.jwter.IsValidJWT(refreshToken); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			accessToken, refreshToken, err := middle.RefreshTokens(refreshToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			middle.cookier.SetAuthCookies(w, accessToken, refreshToken)

			next(w, r)

		}

		http.Error(w, errs.ErrInValidToken().Error(), http.StatusUnauthorized)

	}
}

// Refresh tokens for user
func (middle middleware) RefreshTokens(refreshToken string) (string, string, error) {
	token, err := middle.jwter.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errs.ErrInValidToken()
	}

	userId, ok := claims[jwter.UserId].(int)
	if !ok {
		return "", "", errs.ErrInValidToken()
	}

	user, err := middle.userRepo.FindUserById(userId)
	if err != nil {
		return "", "", errs.ErrInValidToken()
	}

	accessToken, refreshToken, err := middle.jwter.GenerateTokens(&user)
	if err != nil {
		return "", "", err
	}

	err = middle.tokenRepo.AddRefreshToken(user.Id, refreshToken)
	if err != nil {
		return "", "", errs.ErrCantAddToken()
	}

	return accessToken, refreshToken, nil
}
