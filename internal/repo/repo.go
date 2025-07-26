package repo

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const userTable = "users"

type Repo interface {
	RegisterUser(user *models.RegisterReq) error
	FindUserByEmail(em string) (models.User, error)
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

func (this repo) RegisterUser(user *models.RegisterReq) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (email, passw) VALUES (?, ?)", this.dbName, userTable)
	_, err := this.db.Exec(query, user.Email, user.Password)

	return err
}

func (this repo) FindUserByEmail(em string) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE email= ?", this.dbName, userTable)
	res := this.db.QueryRow(query, em)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil

}
