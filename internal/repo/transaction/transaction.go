package transaction

import (
	"database/sql"

	"github.com/IvanDrf/polls-site/config"
)

type Transactioner interface {
	StartTransaction()
	CommitTransaction()
	RollBackTransaction()
}

type transactioner struct {
	dbName string
	db     *sql.DB
}

func NewTransactioner(cfg *config.Config, db *sql.DB) Transactioner {
	return transactioner{dbName: cfg.DBName, db: db}
}

func (t transactioner) StartTransaction() {
	query := "START TRANSACTION"

	t.db.Exec(query)
}

func (t transactioner) CommitTransaction() {
	query := "COMMIT"

	t.db.Exec(query)
}

func (t transactioner) RollBackTransaction() {
	query := "ROLLBACK"

	t.db.Exec(query)
}
