package answers

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const answersTable = "answers"

type AnswersRepo interface {
	AddAnswer(answ *models.Answer) error
	FindAnswer(answ string, questionId int) (models.Answer, error)
}

type answersRepo struct {
	dbName string
	db     *sql.DB
}

func NewAnswersRepo(cfg *config.Config, db *sql.DB) AnswersRepo {
	return answersRepo{
		dbName: cfg.DBName,
		db:     db,
	}
}

func (r answersRepo) AddAnswer(answ *models.Answer) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (answ, question_id) VALUES (?, ?)", r.dbName, answersTable)
	_, err := r.db.Exec(query, answ.Answer, answ.QuestionId)

	return err
}

func (r answersRepo) FindAnswer(answ string, questionId int) (models.Answer, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE answ = ? AND question_id = ?", r.dbName, answersTable)
	res, err := r.db.Query(query, answ, questionId)
	if err != nil {
		return models.Answer{}, err
	}

	a := models.Answer{}
	err = res.Scan(&a.Id, &a.Answer, &a.QuestionId)

	return a, err
}
