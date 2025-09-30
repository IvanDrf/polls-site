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

const (
	start    = "START TRANSACTION"
	commit   = "COMMIT"
	rollBack = "ROLLBACK"
)

func (t transactioner) StartTransaction() {
	t.db.Exec(start)
}

func (t transactioner) CommitTransaction() {
	t.db.Exec(commit)
}

func (t transactioner) RollBackTransaction() {
	t.db.Exec(rollBack)
}
