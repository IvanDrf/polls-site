package poll

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
	"github.com/IvanDrf/polls-site/internal/repo/transaction"
	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type PollService interface {
	AddPoll(poll *models.Poll, r *http.Request) (models.PollId, error)
	DeletePoll(poll *models.Poll, r *http.Request) error
}

type pollService struct {
	answRepo   answers.AnswersRepo
	questRepo  questions.QuestionRepo
	tokensRepo jwt.JWTRepo

	transaction transaction.Transactioner

	jwter jwter.Jwter

	logger *slog.Logger
}

func NewPollService(cfg *config.Config, db *sql.DB, logger *slog.Logger) PollService {
	return pollService{
		answRepo:  answers.NewAnswersRepo(cfg, db),
		questRepo: questions.NewQuestionRepo(cfg, db),

		tokensRepo: jwt.NewTokensRepo(cfg, db),

		transaction: transaction.NewTransactioner(cfg, db),

		jwter: jwter.NewJwter(cfg),

		logger: logger,
	}
}

const (
	contextTime = 5 * time.Second
)

func (p pollService) AddPoll(poll *models.Poll, r *http.Request) (models.PollId, error) {
	token, err := p.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return models.PollId{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	poll.UserId, err = p.tokensRepo.FindUserId(ctx, token)
	if err != nil {
		return models.PollId{}, errs.ErrCantFindUserId()
	}

	p.transaction.StartTransaction()

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	questionId, err := p.questRepo.AddQuestion(ctx, poll)
	if err != nil {
		p.transaction.RollBackTransaction()
		return models.PollId{}, errs.ErrCantAddQuestion()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	err = p.answRepo.AddAnswers(ctx, poll.Answers, questionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return models.PollId{}, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	answId, err := p.answRepo.FindAnswersId(ctx, questionId, len(poll.Answers))
	if err != nil {
		p.transaction.RollBackTransaction()
		return models.PollId{}, errs.ErrCantFindAnswers()
	}

	p.transaction.CommitTransaction()

	return models.PollId{Id: questionId, AnswersId: answId}, nil
}

func (p pollService) DeletePoll(poll *models.Poll, r *http.Request) error {
	token, err := p.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		p.transaction.RollBackTransaction()
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	poll.UserId, err = p.tokensRepo.FindUserId(ctx, token)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantFindUserId()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	question, err := p.questRepo.FindQuestionById(ctx, poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantFindQuestion()
	}

	if poll.UserId != question.UserId {
		p.transaction.RollBackTransaction()
		return errs.ErrNotAdmin()
	}

	// Do not open transaction cuz it's already open in DeleteAllVotes

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	err = p.answRepo.DeleteAllAnswers(ctx, poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantDeleteAnswer()
	}

	ctx, cancel = context.WithTimeout(context.Background(), contextTime)
	defer cancel()

	err = p.questRepo.DeleteQuestionById(ctx, poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantDeleteQuestion()
	}

	p.transaction.CommitTransaction()

	return nil
}
