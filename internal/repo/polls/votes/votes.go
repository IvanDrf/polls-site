package votes

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const votesTable = "votes"

type VotesRepo interface {
	AddVote(ctx context.Context, vote *models.Vote) error
	FindVote(ctx context.Context, questionId, userId int) (int, error)
	CountVotes(ctx context.Context, questionId int) (models.PollRes, error)

	DeleteVote(ctx context.Context, questionId int, userId int) error
	DeleteAllVotes(ctx context.Context, questionId int) error
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

func (v votesRepo) AddVote(ctx context.Context, vote *models.Vote) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (question_id, answ_id, user_id) VALUES (?, ?, ?)", v.dbName, votesTable)
	_, err := v.db.ExecContext(ctx, query, vote.QuestionId, vote.AnswerId, vote.UserId)

	return err
}

func (v votesRepo) FindVote(ctx context.Context, questionId, userId int) (int, error) {
	query := fmt.Sprintf("SELECT id FROM %s.%s WHERE question_id = ? AND user_id = ?", v.dbName, votesTable)
	rows := v.db.QueryRowContext(ctx, query, questionId, userId)

	id := 0
	err := rows.Scan(&id)

	return id, err
}

func (v votesRepo) CountVotes(ctx context.Context, questionId int) (models.PollRes, error) {
	query := fmt.Sprintf("SELECT answ_id, user_id FROM %s.%s WHERE question_id = ?", v.dbName, votesTable)
	rows, err := v.db.QueryContext(ctx, query, questionId)
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

func (v votesRepo) DeleteVote(ctx context.Context, questionId int, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE question_id = ? AND user_id = ?", v.dbName, votesTable)
	_, err := v.db.ExecContext(ctx, query, questionId, userId)

	return err
}

func (v votesRepo) DeleteAllVotes(ctx context.Context, questionId int) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE question_id = ?", v.dbName, votesTable)
	_, err := v.db.ExecContext(ctx, query, questionId)

	return err
}
