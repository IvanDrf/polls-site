package questions

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const (
	questionTable = "questions"
)

type QuestionRepo interface {
	AddQuestionPoll(question string) error
	DeleteQuestionPoll(question *models.Question) error

	FindQuestionPoll(question string) (models.Question, error)
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

func (r questionRepo) AddQuestionPoll(question string) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (question) VALUES (?)", r.dbName, questionTable)
	_, err := r.db.Exec(query, question)

	return err
}

func (r questionRepo) DeleteQuestionPoll(question *models.Question) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE question = ? AND id = ?", r.dbName, questionTable)
	_, err := r.db.Exec(query, question.Question, question.Id)

	return err
}

func (r questionRepo) FindQuestionPoll(question string) (models.Question, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE question = ?", r.dbName, questionTable)
	res := r.db.QueryRow(query, question)

	ques := models.Question{}
	err := res.Scan(&ques.Question, &ques.Id)
	return ques, err
}
