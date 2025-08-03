package repo

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
)

const tokensTable = "tokens"

type TokensRepo interface {
	AddRefreshToken(userId int, refresh string) error
}

// TODO write find method
type tokensRepo struct {
	dbName string
	db     *sql.DB
}

func NewTokensRepo(cfg *config.Config, db *sql.DB) TokensRepo {
	return tokensRepo{dbName: cfg.DBName, db: db}
}

func (t tokensRepo) AddRefreshToken(userId int, refresh string) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (user_id, token) VALUES(?, ?)", t.dbName, tokensTable)
	_, err := t.db.Exec(query, userId, refresh)

	return err
}
