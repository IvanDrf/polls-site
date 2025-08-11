package answers

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
)

const answersTable = "answers"

type AnswersRepo interface {
	AddAnswer(answ *models.Answer) (int, error)
	AddAnswers(answ []string, questionId int) error

	DeleteAnswer(answ *models.Answer) error
	DeleteAnswers(answ []string, questionId int) error

	FindAnswer(answ string, questionId int) (models.Answer, error)
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

func (r answersRepo) AddAnswer(answ *models.Answer) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s.%s (answ, question_id) VALUES (?, ?)", r.dbName, answersTable)
	res, err := r.db.Exec(query, answ.Answer, answ.QuestionId)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (r answersRepo) AddAnswers(answ []string, questionId int) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (answ, question_id) VALUES", r.dbName, answersTable)

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

	_, err := r.db.Exec(query, args...)

	return err
}

func (r answersRepo) deleteAnswer(questionId int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE question_id = ?", r.dbName, answersTable)
	_, err := r.db.Exec(query, questionId)

	return err
}

func (r answersRepo) DeleteAnswer(answ *models.Answer) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE answ = ? AND question_id = ?", r.dbName, answersTable)
	_, err := r.db.Exec(query, answ.Answer, answ.QuestionId)

	return err
}

func (r answersRepo) DeleteAnswers(answ []string, questionId int) error {
	fail := false

	wg := new(sync.WaitGroup)
	for i := range answ {
		wg.Add(1)
		go func(i int, fail *bool) {
			defer wg.Done()
			err := r.deleteAnswer(questionId)

			if err != nil {
				*fail = true
			}
		}(i, &fail)
	}

	wg.Wait()

	if fail {
		return errs.ErrCantDeleteAnswer()
	}

	return nil
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

func (r answersRepo) FindAnswersId(questionId int, size int) ([]int, error) {
	query := fmt.Sprintf("SELECT id FROM %s.%s WHERE question_id = ?", r.dbName, answersTable)
	rows, err := r.db.Query(query, questionId)
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
