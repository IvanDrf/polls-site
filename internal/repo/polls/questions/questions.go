package questions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const (
	questionTable = "questions"
)

type QuestionRepo interface {
	AddQuestion(ctx context.Context, poll *models.Poll) (int, error)
	DeleteQuestionById(ctx context.Context, id int) error

	FindQuestionById(ctx context.Context, id int) (models.Question, error)
}

type questionRepo struct {
	dbName string
	db     *sql.DB
}

func NewQuestionRepo(cfg *config.Config, db *sql.DB) QuestionRepo {
	return questionRepo{
		dbName: cfg.DBName,
		db:     db,
	}
}

func (q questionRepo) AddQuestion(ctx context.Context, poll *models.Poll) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s.%s (question, user_id) VALUES (?, ?)", q.dbName, questionTable)
	res, err := q.db.ExecContext(ctx, query, poll.Question, poll.UserId)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (q questionRepo) DeleteQuestionById(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE id = ?", q.dbName, questionTable)
	_, err := q.db.ExecContext(ctx, query, id)

	return err
}

func (q questionRepo) FindQuestionById(ctx context.Context, id int) (models.Question, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE id = ?", q.dbName, questionTable)
	res := q.db.QueryRowContext(ctx, query, id)

	ques := models.Question{}
	err := res.Scan(&ques.Question, &ques.Id, &ques.UserId)

	return ques, err
}
