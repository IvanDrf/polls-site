package jwter

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/golang-jwt/jwt"
)

const (
	AccessToken  = "access_jwt"
	RefreshToken = "refresh_jwt"

	UserId = "user_id"

	accessTime  = 24 * time.Hour
	refreshTime = 7 * 24 * time.Hour
)

type Jwter interface {
	GetToken(r *http.Request) (string, error)
	ParseToken(tokenSrt string) (*jwt.Token, error)

	GenerateTokens(user *models.User) (string, string, error)
	GenerateAccessJWT(user *models.User) (string, error)
	GenerateRefreshJWT(user *models.User) (string, error)

	IsValidJWT(tokenStr string) error
}

type jwter struct {
	jwtSecret []byte
}

func NewJwter(cfg *config.Config) Jwter {
	return jwter{jwtSecret: cfg.JWT}
}

// Get token from header/cookies
func (j jwter) GetToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			return tokenParts[1], nil
		}
	}

	cookie, err := r.Cookie(AccessToken)
	if err != nil {
		return "", errs.ErrCantFindToken()
	}

	return cookie.Value, nil
}

func (j jwter) ParseToken(tokenSrt string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenSrt, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrIncorrectJWTMethod()
		}

		return j.jwtSecret, nil
	})

	return token, err
}

func (j jwter) GenerateTokens(user *models.User) (string, string, error) {
	wg := new(sync.WaitGroup)

	accessToken, refreshToken := "", ""
	var errAccess, errRefresh error

	wg.Add(1)
	go func() {
		defer wg.Done()
		accessToken, errAccess = j.GenerateAccessJWT(user)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		refreshToken, errRefresh = j.GenerateRefreshJWT(user)
	}()

	wg.Wait()
	if errAccess != nil || errRefresh != nil {
		return "", "", errs.ErrCantCreateToken()
	}

	return accessToken, refreshToken, nil
}

func (j jwter) GenerateAccessJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(accessTime).Unix(),
	})

	return token.SignedString(j.jwtSecret)
}

func (j jwter) GenerateRefreshJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(refreshTime).Unix(),
	})

	return token.SignedString(j.jwtSecret)
}

func (j jwter) IsValidJWT(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrIncorrectJWTMethod()
		}

		return j.jwtSecret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errs.ErrInValidToken()
	}

	return nil
}
