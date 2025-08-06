package votes

import (
	"database/sql"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
)

const votesTable = "polls_results"

type VotesRepo interface {
	AddVote(vote *models.Vote) error
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

func (r votesRepo) AddVote(vote *models.Vote) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (question_id, answ, user_id) VALUES (?, ?, ?)", r.dbName, votesTable)
	_, err := r.db.Exec(query, vote.QuestionId, vote.Answer, vote.UserId)

	return err
}
