package service

import (
	"database/sql"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo/auth/tokens"
	"github.com/IvanDrf/polls-site/internal/repo/polls/answers"
	"github.com/IvanDrf/polls-site/internal/repo/polls/questions"
	"github.com/IvanDrf/polls-site/internal/repo/polls/votes"
	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type PollService interface {
	AddPoll(poll *models.Poll, r *http.Request) (models.PollId, error)
	DeletePoll(poll *models.Poll) error
}

type pollService struct {
	answRepo   answers.AnswersRepo
	questRepo  questions.QuestionRepo
	votesRepo  votes.VotesRepo
	tokensRepo tokens.TokensRepo

	jwter jwter.Jwter
}

func NewPollService(cfg *config.Config, db *sql.DB) PollService {
	return pollService{
		answRepo:   answers.NewAnswersRepo(cfg, db),
		questRepo:  questions.NewQuestionRepo(cfg, db),
		votesRepo:  votes.NewVotesRepo(cfg, db),
		tokensRepo: tokens.NewTokensRepo(cfg, db),

		jwter: jwter.NewJwter(cfg),
	}
}

// TODO: add user id in question -> see man, who created poll
func (p pollService) AddPoll(poll *models.Poll, r *http.Request) (models.PollId, error) {
	token, err := p.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return models.PollId{}, err
	}

	poll.UserId, err = p.tokensRepo.FindUserId(token)
	if err != nil {
		return models.PollId{}, errs.ErrCantFindUser()
	}

	questionId, err := p.questRepo.AddQuestionPoll(poll)
	if err != nil {
		return models.PollId{}, errs.ErrCantAddQuestion()
	}

	err = p.answRepo.AddAnswers(poll.Answers, questionId)
	if err != nil {
		p.answRepo.DeleteAnswers(poll.Answers, questionId)
		p.questRepo.DeleteQuestionPollById(questionId)

		return models.PollId{}, err
	}

	return models.PollId{Id: questionId}, nil
}

func (p pollService) DeletePoll(poll *models.Poll) error {
	question, err := p.questRepo.FindQuestionPollById(poll.Id)
	if err != nil {
		return errs.ErrCantFindQuestion()
	}

	err = p.answRepo.DeleteAnswers(poll.Answers, question.Id)
	if err != nil {
		return errs.ErrCantAddAnswer()
	}

	err = p.questRepo.DeleteQuestionPollById(question.Id)
	if err != nil {
		return errs.ErrCantDeleteQuestion()
	}

	return nil
}
