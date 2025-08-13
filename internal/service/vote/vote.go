package vote

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo/auth/jwt"
	"github.com/IvanDrf/polls-site/internal/repo/polls/answers"
	"github.com/IvanDrf/polls-site/internal/repo/polls/questions"
	"github.com/IvanDrf/polls-site/internal/repo/polls/votes"
	"github.com/IvanDrf/polls-site/internal/repo/transaction"
	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type VoteService interface {
	VoteInPoll(vote *models.Vote, r *http.Request) (models.PollRes, error)

	DeleteVoteInPoll(vote *models.Vote, r *http.Request) (models.PollRes, error)
	DeleteAllVotesInPoll(poll *models.Poll, r *http.Request) error
}

type voteService struct {
	answRepo   answers.AnswersRepo
	questRepo  questions.QuestionRepo
	votesRepo  votes.VotesRepo
	tokensRepo jwt.JWTRepo

	transaction transaction.Transactioner

	jwter jwter.Jwter

	logger *slog.Logger
}

func NewVoteService(cfg *config.Config, db *sql.DB, logger *slog.Logger) VoteService {
	return voteService{
		answRepo:  answers.NewAnswersRepo(cfg, db),
		questRepo: questions.NewQuestionRepo(cfg, db),
		votesRepo: votes.NewVotesRepo(cfg, db),

		tokensRepo:  jwt.NewTokensRepo(cfg, db),
		transaction: transaction.NewTransactioner(cfg, db),

		jwter:  jwter.NewJwter(cfg),
		logger: logger,
	}
}

func (v voteService) VoteInPoll(vote *models.Vote, r *http.Request) (models.PollRes, error) {
	token, err := v.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return models.PollRes{}, err
	}

	vote.UserId, err = v.tokensRepo.FindUserId(token)
	if err != nil {
		return models.PollRes{}, errs.ErrCantFindUserId()
	}

	question, err := v.questRepo.FindQuestionById(vote.QuestionId)
	if err != nil || question.Id != vote.QuestionId {
		return models.PollRes{}, errs.ErrCantFindQuestion()
	}

	_, err = v.votesRepo.FindVote(vote.QuestionId, vote.UserId)
	if err == nil {
		return models.PollRes{}, errs.ErrAlreadyVoted()
	}

	_, err = v.answRepo.FindAnswerById(vote.AnswerId, vote.QuestionId)
	if err != nil {
		return models.PollRes{}, errs.ErrBadAnswerId()
	}

	v.transaction.StartTransaction()

	err = v.votesRepo.AddVote(vote)
	if err != nil {
		v.transaction.RollBackTransaction()
		return models.PollRes{}, errs.ErrCantAddVote()
	}

	res, err := v.votesRepo.CountVotes(vote.QuestionId)
	if err != nil {
		v.transaction.RollBackTransaction()
		return models.PollRes{}, errs.ErrCantCountVotes()
	}

	v.transaction.CommitTransaction()

	return res, nil
}

func (v voteService) DeleteVoteInPoll(vote *models.Vote, r *http.Request) (models.PollRes, error) {
	token, err := v.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return models.PollRes{}, err
	}

	vote.UserId, err = v.tokensRepo.FindUserId(token)
	if err != nil {
		return models.PollRes{}, errs.ErrCantFindUserId()
	}

	_, err = v.questRepo.FindQuestionById(vote.QuestionId)
	if err != nil {
		return models.PollRes{}, errs.ErrCantFindQuestion()
	}

	_, err = v.votesRepo.FindVote(vote.QuestionId, vote.UserId)
	if err != nil {
		return models.PollRes{}, errs.ErrDidntVote()
	}

	err = v.votesRepo.DeleteVote(vote.QuestionId, vote.UserId)
	if err != nil {
		return models.PollRes{}, errs.ErrCantDeleteVote()
	}

	voteRes, err := v.votesRepo.CountVotes(vote.QuestionId)
	if err != nil {
		return models.PollRes{}, errs.ErrCantCountVotes()
	}

	return voteRes, nil
}

func (v voteService) DeleteAllVotesInPoll(poll *models.Poll, r *http.Request) error {
	token, err := v.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return err
	}

	poll.UserId, err = v.tokensRepo.FindUserId(token)
	if err != nil {
		return errs.ErrCantFindUserId()
	}

	question, err := v.questRepo.FindQuestionById(poll.QuestionId)
	if err != nil {
		return errs.ErrCantFindQuestion()
	}

	if poll.UserId != question.UserId {
		return errs.ErrNotAdmin()
	}

	v.transaction.StartTransaction()

	err = v.votesRepo.DeleteAllVotes(poll.QuestionId)
	if err != nil {
		v.transaction.RollBackTransaction()
		return errs.ErrCantDeleteAllVotes()
	}

	// Don't commit transaction cuz deleting another tables

	return nil
}
