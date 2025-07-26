package jwt

import (
	"time"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/golang-jwt/jwt"
)

type Jwter interface {
	GenerateJWT(user *models.User) (string, error)
	IsValidJWT(tokenStr string) error
}

type jwter struct {
	jwtSecret []byte
}

func NewJwter(cfg *config.Config) Jwter {
	return jwter{jwtSecret: cfg.JWT}
}

func (this jwter) GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(this.jwtSecret)
}

func (this jwter) IsValidJWT(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrIncorrectJWTMethod()
		}

		return this.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return errs.ErrInValidToken()
	}

	return nil
}
