package jwt

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const jwtTable = "jwt"

type JWTRepo interface {
	AddRefreshToken(ctx context.Context, userId int, refresh string) error
	UpdateRefreshToken(ctx context.Context, userId int, refresh string) error

	FindRefreshToken(ctx context.Context, userId int) (models.JWT, error)
	FindUserId(ctx context.Context, refreshToken string) (int, error)
}

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

func (t jwtRepo) AddRefreshToken(ctx context.Context, userId int, refresh string) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (user_id, token) VALUES(?, ?)", t.dbName, jwtTable)
	_, err := t.db.ExecContext(ctx, query, userId, refresh)

	return err
}

func (t jwtRepo) UpdateRefreshToken(ctx context.Context, userId int, refresh string) error {
	query := fmt.Sprintf("UPDATE %s.%s SET token = ? WHERE user_id = ?", t.dbName, jwtTable)
	_, err := t.db.ExecContext(ctx, query, refresh, userId)

	return err
}

func (t jwtRepo) FindRefreshToken(ctx context.Context, userId int) (models.JWT, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE user_id = ?", t.dbName, jwtTable)
	res := t.db.QueryRowContext(ctx, query, userId)

	token := models.JWT{}
	if err := res.Scan(&token.Id, &token.UserId, &token.Refresh); err != nil {
		return models.JWT{}, err
	}

	return token, nil
}

func (t jwtRepo) FindUserId(ctx context.Context, refreshToken string) (int, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE token = ?", t.dbName, jwtTable)
	res := t.db.QueryRowContext(ctx, query, refreshToken)

	token := models.JWT{}
	if err := res.Scan(&token.Id, &token.UserId, &token.Refresh); err != nil {
		return -1, err
	}

	return token.UserId, nil
}
