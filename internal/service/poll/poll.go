package poll

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
	questionId, err := p.questRepo.AddQuestion(poll)
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

func (p pollService) DeletePoll(poll *models.Poll, r *http.Request) error {
	token, err := p.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		p.transaction.RollBackTransaction()
		return err
	}

	poll.UserId, err = p.tokensRepo.FindUserId(token)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantFindUserId()
	}

	question, err := p.questRepo.FindQuestionById(poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantFindQuestion()
	}

	if poll.UserId != question.UserId {
		p.transaction.RollBackTransaction()
		return errs.ErrNotAdmin()
	}

	// Do not open transaction cuz it's already open in DeleteAllVotes

	err = p.answRepo.DeleteAllAnswers(poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantDeleteAnswer()
	}

	err = p.questRepo.DeleteQuestionById(poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantDeleteQuestion()
	}

	p.transaction.CommitTransaction()

	return nil
}
