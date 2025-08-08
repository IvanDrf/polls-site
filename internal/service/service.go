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
	AddPoll(poll *models.Poll, r *http.Request) error
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
func (p pollService) AddPoll(poll *models.Poll, r *http.Request) error {
	token, err := p.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return err
	}

	poll.UserId, err = p.tokensRepo.FindUserId(token)
	if err != nil {
		return errs.ErrCantFindUser()
	}

	err = p.questRepo.AddQuestionPoll(poll)
	if err != nil {
		return errs.ErrCantAddQuestion()
	}

	question, err := p.questRepo.FindQuestionPoll(poll.Question)
	if err != nil {
		return errs.ErrCantAddQuestion()
	}

	err = p.answRepo.AddAnswers(poll.Answers, question.Id)
	if err != nil {
		p.answRepo.DeleteAnswers(poll.Answers, question.Id)
		p.questRepo.DeleteQuestionPoll(&question)

		return err
	}

	return nil
}
