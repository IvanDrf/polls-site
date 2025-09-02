package jwter

import (
	"net/http"
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

	accessTime  = 15 * time.Minute
	refreshTime = 7 * 24 * time.Hour
)

type Jwter interface {
	// Get access/refresh token from cookies
	GetToken(r *http.Request, tokenType string) (string, error)

	ParseToken(tokenSrt string) (*jwt.Token, error)

	// Return access, refresh, error
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

func (j jwter) GetToken(r *http.Request, tokenType string) (string, error) {
	if tokenType != AccessToken && tokenType != RefreshToken {
		return "", errs.ErrInValidToken()
	}

	cookie, err := r.Cookie(tokenType)
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
