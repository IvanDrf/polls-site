package repo

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const userTable = "users"

type UserRepo interface {
	AddUser(user *models.UserReq) error

	FindUserByEmail(em string) (models.User, error)
	FindUserById(id int) (models.User, error)
}

type repo struct {
	dbName string
	db     *sql.DB
}

func NewRepo(cfg *config.Config, db *sql.DB) UserRepo {
	return repo{
		dbName: cfg.DBName,
		db:     db,
	}
}

func (r repo) AddUser(user *models.UserReq) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (email, passw) VALUES (?, ?)", r.dbName, userTable)
	_, err := r.db.Exec(query, user.Email, user.Password)

	return err
}

func (r repo) FindUserByEmail(em string) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE email= ?", r.dbName, userTable)
	res := r.db.QueryRow(query, em)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r repo) FindUserById(id int) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE id = ?", r.dbName, userTable)
	res := r.db.QueryRow(query, id)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}
