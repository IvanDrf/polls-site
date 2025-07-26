package database

import (
	"database/sql"
	"log"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/go-sql-driver/mysql"
)

func InitDB(cfg *config.Config) *sql.DB {
	dbCFG := &mysql.Config{
		User:   cfg.DBUser,
		Passwd: cfg.DBPassword,
		Net:    "tcp",
		Addr:   cfg.DBHost + cfg.DBPort,
		DBName: cfg.DBName,
	}

	db, err := sql.Open("mysql", dbCFG.FormatDSN())
	if err != nil {
		log.Fatal(errs.ErrDBConnection())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(errs.ErrDBConnection())
	}

	return db
}
