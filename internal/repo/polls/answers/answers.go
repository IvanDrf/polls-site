package answers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const answersTable = "answers"

type AnswersRepo interface {
	AddAnswer(answ *models.Answer) (int, error)
	AddAnswers(answ []string, questionId int) error

	DeleteAnswer(answ *models.Answer) error
	DeleteAllAnswers(questionId int) error

	FindAnswerById(answId int, questionId int) (models.Answer, error)
	//size - amount of answers
	FindAnswersId(questionId int, size int) ([]int, error)
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

func (a answersRepo) AddAnswer(answ *models.Answer) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s.%s (answ, question_id) VALUES (?, ?)", a.dbName, answersTable)
	res, err := a.db.Exec(query, answ.Answer, answ.QuestionId)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (a answersRepo) AddAnswers(answ []string, questionId int) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (answ, question_id) VALUES", a.dbName, answersTable)

	values := make([]string, 0, len(answ))
	args := make([]any, 2*len(answ))
	k := 0

	for i := range answ {
		args[k] = answ[i]
		args[k+1] = questionId
		k += 2

		values = append(values, "(?, ?)")
	}

	query += strings.Join(values, ", ")

	_, err := a.db.Exec(query, args...)

	return err
}

func (a answersRepo) DeleteAnswer(answ *models.Answer) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE answ = ? AND question_id = ?", a.dbName, answersTable)
	_, err := a.db.Exec(query, answ.Answer, answ.QuestionId)

	return err
}

func (a answersRepo) DeleteAllAnswers(questionId int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE question_id = ?", a.dbName, answersTable)
	_, err := a.db.Exec(query, questionId)

	return err
}

func (a answersRepo) FindAnswerById(answId int, questionId int) (models.Answer, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE id = ? AND question_id = ?", a.dbName, answersTable)
	res := a.db.QueryRow(query, answId, questionId)

	answ := models.Answer{}
	err := res.Scan(&answ.Id, &answ.Answer, &answ.QuestionId)
	return answ, err
}

func (a answersRepo) FindAnswersId(questionId int, size int) ([]int, error) {
	query := fmt.Sprintf("SELECT id FROM %s.%s WHERE question_id = ?", a.dbName, answersTable)
	rows, err := a.db.Query(query, questionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	answId := make([]int, 0, size)
	for rows.Next() {
		id := 0
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		answId = append(answId, id)
	}

	return answId, nil

}
