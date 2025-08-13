package jwt

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const jwtTable = "jwt"

type JWTRepo interface {
	AddRefreshToken(userId int, refresh string) error
	UpdateRefreshToken(userId int, refresh string) error

	FindRefreshToken(userId int) (models.JWT, error)
	FindUserId(refreshToken string) (int, error)
}

// TODO write find method
type jwtRepo struct {
	dbName string
	db     *sql.DB
}

func NewTokensRepo(cfg *config.Config, db *sql.DB) JWTRepo {
	return jwtRepo{
		dbName: cfg.DBName,
		db:     db,
	}
}

func (t jwtRepo) AddRefreshToken(userId int, refresh string) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (user_id, token) VALUES(?, ?)", t.dbName, jwtTable)
	_, err := t.db.Exec(query, userId, refresh)

	return err
}

func (t jwtRepo) UpdateRefreshToken(userId int, refresh string) error {
	query := fmt.Sprintf("UPDATE %s.%s SET token = ? WHERE user_id = ?", t.dbName, jwtTable)
	_, err := t.db.Exec(query, refresh, userId)

	return err
}

func (t jwtRepo) FindRefreshToken(userId int) (models.JWT, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE user_id = ?", t.dbName, jwtTable)
	res := t.db.QueryRow(query, userId)

	token := models.JWT{}
	if err := res.Scan(&token.Id, &token.UserId, &token.Refresh); err != nil {
		return models.JWT{}, err
	}

	return token, nil
}

func (t jwtRepo) FindUserId(refreshToken string) (int, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE token = ?", t.dbName, jwtTable)
	res := t.db.QueryRow(query, refreshToken)

	token := models.JWT{}
	if err := res.Scan(&token.Id, &token.UserId, &token.Refresh); err != nil {
		return -1, err
	}

	return token.UserId, nil
}
