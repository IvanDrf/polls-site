package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

const userTable = "users"

type UserRepo interface {
	AddUser(ctx context.Context, user *models.User) (int, error)
	ActivateUser(ctx context.Context, user *models.User) error
	DeleteUnverifiedUsers(ctx context.Context) error

	FindUserById(ctx context.Context, id int) (models.User, error)
	FindUserByEmail(ctx context.Context, em string) (models.User, error)
	FindUserByLink(ctx context.Context, link string) (models.User, error)

	ResetPassword(ctx context.Context, password string, userId int) error
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

func (u userRepo) AddUser(ctx context.Context, user *models.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s.%s (email, passw, verificated, expired, veriftoken) VALUES (?, ?, ?, ?, ?)", u.dbName, userTable)
	res, err := u.db.ExecContext(ctx, query, user.Email, user.Password, user.Verificated, user.Expired, user.VerifToken)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (u userRepo) ActivateUser(ctx context.Context, user *models.User) error {
	query := fmt.Sprintf("UPDATE %s.%s SET verificated = 1 WHERE id = ?", u.dbName, userTable)
	_, err := u.db.ExecContext(ctx, query, user.Id)

	return err
}

func (u userRepo) DeleteUnverifiedUsers(ctx context.Context) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE verificated = 0 AND expired < ?", u.dbName, userTable)
	_, err := u.db.ExecContext(ctx, query, time.Now())

	return err
}

func (u userRepo) FindUserById(ctx context.Context, userId int) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE id = ?", u.dbName, userTable)
	res := u.db.QueryRowContext(ctx, query, userId)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}
func (u userRepo) FindUserByEmail(ctx context.Context, em string) (models.User, error) {
	query := fmt.Sprintf("SELECT id, email, passw FROM %s.%s WHERE email= ?", u.dbName, userTable)
	res := u.db.QueryRowContext(ctx, query, em)

	user := models.User{}
	if err := res.Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u userRepo) FindUserByLink(ctx context.Context, link string) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE veriftoken = ?", u.dbName, userTable)
	res := u.db.QueryRowContext(ctx, query, link)

	user := models.User{}
	expired := ""

	if err := res.Scan(&user.Id, &user.Email, &user.Password, &user.Verificated, &expired, &user.VerifToken); err != nil {
		return models.User{}, err
	}

	var err error
	user.Expired, err = time.Parse("2006-01-02 15:04:05", expired)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u userRepo) ResetPassword(ctx context.Context, password string, userId int) error {
	query := fmt.Sprintf("UPDATE TABLE %s.%s SET passw = ? WHERE user_id = ?", u.dbName, userTable)
	_, err := u.db.ExecContext(ctx, query, password, userId)

	return err
}
