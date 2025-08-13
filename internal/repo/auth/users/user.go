package users

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

	ResetPassword(password string, userId int) error
}

type userRepo struct {
	dbName string
	db     *sql.DB
}

func NewRepo(cfg *config.Config, db *sql.DB) UserRepo {
	return userRepo{
		dbName: cfg.DBName,
		db:     db,
	}
}

func (u userRepo) AddUser(user *models.UserReq) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (email, passw) VALUES (?, ?)", u.dbName, userTable)
	_, err := u.db.Exec(query, user.Email, user.Password)

	return err
}

func (u userRepo) FindUserByEmail(em string) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE email= ?", u.dbName, userTable)
	res := u.db.QueryRow(query, em)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u userRepo) FindUserById(userId int) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE id = ?", u.dbName, userTable)
	res := u.db.QueryRow(query, userId)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u userRepo) ResetPassword(password string, userId int) error {
	query := fmt.Sprintf("UPDATE TABLE %s.%s SET passw = ? WHERE user_id = ?", u.dbName, userTable)
	_, err := u.db.Exec(query, password, userId)

	return err
}
