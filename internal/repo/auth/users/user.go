package users

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const userTable = "users"

type UserRepo interface {
	AddUser(user *models.User) (int, error)
	ActivateUser(user *models.User) error

	FindUserById(id int) (models.User, error)
	FindUserByEmail(em string) (models.User, error)
	FindUserByLink(link string) (models.User, error)

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

func (u userRepo) AddUser(user *models.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s.%s (email, passw, verificated, expired, veriflink) VALUES (?, ?, ?, ?, ?)", u.dbName, userTable)
	res, err := u.db.Exec(query, user.Email, user.Password, user.Verificated, user.Expired, user.VerifLink)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (u userRepo) ActivateUser(user *models.User) error {
	query := fmt.Sprintf("UPDATE TABLE %s.%s SET verificated = 1U T WHERE id = ?", u.dbName, userTable)
	_, err := u.db.Exec(query, user.Id)

	return err
}

func (u userRepo) DeleteUnverifiedUsers() error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE verificated = 0 AND expired < ?", u.dbName, userTable)
	_, err := u.db.Exec(query, time.Now())

	return err
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
func (u userRepo) FindUserByEmail(em string) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE email= ?", u.dbName, userTable)
	res := u.db.QueryRow(query, em)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u userRepo) FindUserByLink(link string) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE veriflink = ?", u.dbName, userTable)
	res := u.db.QueryRow(query, link)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password, &user.Verificated, &user.Expired, &user.VerifLink); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u userRepo) ResetPassword(password string, userId int) error {
	query := fmt.Sprintf("UPDATE TABLE %s.%s SET passw = ? WHERE user_id = ?", u.dbName, userTable)
	_, err := u.db.Exec(query, password, userId)

	return err
}
