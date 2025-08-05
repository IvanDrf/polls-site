package questions

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const (
	pollsTable = "polls"
)

type QuestionRepo interface {
	AddQuestionPoll(question string) error
	FindQuestionPoll(question string) (models.Question, error)
}

type questionRepo struct {
	dbName string
	db     *sql.DB
}

func NewPollsRepo(cfg *config.Config, db *sql.DB) QuestionRepo {
	return questionRepo{
		dbName: cfg.DBName,
		db:     db,
	}
}

func (r questionRepo) AddQuestionPoll(question string) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (question) VALUES (?)", r.dbName, pollsTable)
	_, err := r.db.Exec(query)

	return err
}

func (r questionRepo) FindQuestionPoll(question string) (models.Question, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE question = ?", r.dbName, pollsTable)
	res := r.db.QueryRow(query, question)

	ques := models.Question{}
	err := res.Scan(&ques.Id, &ques.Question)

	return ques, err
}
