package poll

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo/auth/tokens"
	"github.com/IvanDrf/polls-site/internal/repo/polls/answers"
	"github.com/IvanDrf/polls-site/internal/repo/polls/questions"
	"github.com/IvanDrf/polls-site/internal/repo/transaction"
	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type PollService interface {
	AddPoll(poll *models.Poll, r *http.Request) (models.PollId, error)
}

type pollService struct {
	answRepo   answers.AnswersRepo
	questRepo  questions.QuestionRepo
	tokensRepo tokens.TokensRepo

	transaction transaction.Transactioner

	jwter jwter.Jwter

	logger *slog.Logger
}

func NewPollService(cfg *config.Config, db *sql.DB, logger *slog.Logger) PollService {
	return pollService{
		answRepo:  answers.NewAnswersRepo(cfg, db),
		questRepo: questions.NewQuestionRepo(cfg, db),

		tokensRepo: tokens.NewTokensRepo(cfg, db),

		transaction: transaction.NewTransactioner(cfg, db),

		jwter: jwter.NewJwter(cfg),

		logger: logger,
	}
}

func (p pollService) AddPoll(poll *models.Poll, r *http.Request) (models.PollId, error) {
	token, err := p.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return models.PollId{}, err
	}

	poll.UserId, err = p.tokensRepo.FindUserId(token)
	if err != nil {
		return models.PollId{}, errs.ErrCantFindUserId()
	}

	p.transaction.StartTransaction()
	questionId, err := p.questRepo.AddQuestionPoll(poll)
	if err != nil {
		p.transaction.RollBackTransaction()
		return models.PollId{}, errs.ErrCantAddQuestion()
	}

	err = p.answRepo.AddAnswers(poll.Answers, questionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return models.PollId{}, err
	}

	answId, err := p.answRepo.FindAnswersId(questionId, len(poll.Answers))
	if err != nil {
		p.transaction.RollBackTransaction()
		return models.PollId{}, errs.ErrCantFindAnswers()
	}

	p.transaction.CommitTransaction()

	return models.PollId{Id: questionId, AnswersId: answId}, nil
}
