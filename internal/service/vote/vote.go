package vote

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

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

const (
	contextTime = 5 * time.Second
)

func (v voteService) VoteInPoll(vote *models.Vote, r *http.Request) (models.PollRes, error) {
	token, err := v.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return models.PollRes{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	vote.UserId, err = v.tokensRepo.FindUserId(ctx, token)
	if err != nil {
		return models.PollRes{}, errs.ErrCantFindUserId()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	question, err := v.questRepo.FindQuestionById(ctx, vote.QuestionId)
	if err != nil || question.Id != vote.QuestionId {
		return models.PollRes{}, errs.ErrCantFindQuestion()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	_, err = v.votesRepo.FindVote(ctx, vote.QuestionId, vote.UserId)
	if err == nil {
		return models.PollRes{}, errs.ErrAlreadyVoted()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	_, err = v.answRepo.FindAnswerById(ctx, vote.AnswerId, vote.QuestionId)
	if err != nil {
		return models.PollRes{}, errs.ErrBadAnswerId()
	}

	v.transaction.StartTransaction()

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	err = v.votesRepo.AddVote(ctx, vote)
	if err != nil {
		v.transaction.RollBackTransaction()
		return models.PollRes{}, errs.ErrCantAddVote()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	res, err := v.votesRepo.CountVotes(ctx, vote.QuestionId)
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

	ctx, cancel := context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	vote.UserId, err = v.tokensRepo.FindUserId(ctx, token)
	if err != nil {
		return models.PollRes{}, errs.ErrCantFindUserId()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	_, err = v.questRepo.FindQuestionById(ctx, vote.QuestionId)
	if err != nil {
		return models.PollRes{}, errs.ErrCantFindQuestion()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	_, err = v.votesRepo.FindVote(ctx, vote.QuestionId, vote.UserId)
	if err != nil {
		return models.PollRes{}, errs.ErrDidntVote()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	err = v.votesRepo.DeleteVote(ctx, vote.QuestionId, vote.UserId)
	if err != nil {
		return models.PollRes{}, errs.ErrCantDeleteVote()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	voteRes, err := v.votesRepo.CountVotes(ctx, vote.QuestionId)
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

	ctx, cancel := context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	poll.UserId, err = v.tokensRepo.FindUserId(ctx, token)
	if err != nil {
		return errs.ErrCantFindUserId()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	question, err := v.questRepo.FindQuestionById(ctx, poll.QuestionId)
	if err != nil {
		return errs.ErrCantFindQuestion()
	}

	if poll.UserId != question.UserId {
		return errs.ErrNotAdmin()
	}

	v.transaction.StartTransaction()

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	err = v.votesRepo.DeleteAllVotes(ctx, poll.QuestionId)
	if err != nil {
		v.transaction.RollBackTransaction()
		return errs.ErrCantDeleteAllVotes()
	}

	// Don't commit transaction cuz deleting another tables

	return nil
}
