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
	AddQuestionPoll(poll *models.Poll) (int, error)
	DeleteQuestionPollById(id int) error

	FindQuestionPollById(id int) (models.Question, error)
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

func (r questionRepo) AddQuestionPoll(poll *models.Poll) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s.%s (question, user_id) VALUES (?, ?)", r.dbName, questionTable)
	res, err := r.db.Exec(query, poll.Question, poll.UserId)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (r questionRepo) DeleteQuestionPollById(id int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE id = ?", r.dbName, questionTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r questionRepo) FindQuestionPollById(id int) (models.Question, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE id = ?", r.dbName, questionTable)
	res := r.db.QueryRow(query, id)

	ques := models.Question{}
	err := res.Scan(&ques.Question, &ques.Id, &ques.UserId)

	return ques, err
}
