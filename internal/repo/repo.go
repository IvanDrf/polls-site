package repo

import (
	"database/sql"

	"github.com/IvanDrf/polls-site/config"
)

type Repo interface {
}

type repo struct {
	dbName string
	db     *sql.DB
}

func NewRepo(cfg *config.Config, db *sql.DB) Repo {
	return repo{
		dbName: cfg.DBName,
		db:     db,
	}
}
