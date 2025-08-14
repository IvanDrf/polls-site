package votes

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const votesTable = "votes"

type VotesRepo interface {
	AddVote(vote *models.Vote) error
	FindVote(questionId, userId int) (int, error)
	CountVotes(questionId int) (models.PollRes, error)

	DeleteVote(questionId int, userId int) error
	DeleteAllVotes(questionId int) error
}

type votesRepo struct {
	dbName string
	db     *sql.DB
}

func NewVotesRepo(cfg *config.Config, db *sql.DB) VotesRepo {
	return votesRepo{
		dbName: cfg.DBName,
		db:     db,
	}
}

func (v votesRepo) AddVote(vote *models.Vote) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (question_id, answ_id, user_id) VALUES (?, ?, ?)", v.dbName, votesTable)
	_, err := v.db.Exec(query, vote.QuestionId, vote.AnswerId, vote.UserId)

	return err
}

func (v votesRepo) FindVote(questionId, userId int) (int, error) {
	query := fmt.Sprintf("SELECT id FROM %s.%s WHERE question_id = ? AND user_id = ?", v.dbName, votesTable)
	rows := v.db.QueryRow(query, questionId, userId)

	id := 0
	err := rows.Scan(&id)

	return id, err
}

func (v votesRepo) CountVotes(questionId int) (models.PollRes, error) {
	query := fmt.Sprintf("SELECT answ_id, user_id FROM %s.%s WHERE question_id = ?", v.dbName, votesTable)
	rows, err := v.db.Query(query, questionId)
	if err != nil {
		return models.PollRes{}, err
	}
	defer rows.Close()

	pollRes := models.PollRes{QuestionId: questionId, Answers: make(map[int]int)}
	for rows.Next() {
		answId, userId := 0, 0
		if err := rows.Scan(&answId, &userId); err != nil {
			return models.PollRes{}, err
		}

		pollRes.Answers[answId]++
	}

	return pollRes, nil
}

func (v votesRepo) DeleteVote(questionId int, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE question_id = ? AND user_id = ?", v.dbName, votesTable)
	_, err := v.db.Exec(query, questionId, userId)

	return err
}

func (v votesRepo) DeleteAllVotes(questionId int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE question_id = ?", v.dbName, votesTable)
	_, err := v.db.Exec(query, questionId)

	return err
}
