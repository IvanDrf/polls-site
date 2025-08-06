package service

import (
	"database/sql"
	"log"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo/polls/answers"
	"github.com/IvanDrf/polls-site/internal/repo/polls/questions"
	"github.com/IvanDrf/polls-site/internal/repo/polls/votes"
)

type PollService interface {
	AddPoll(poll *models.Poll) error
}

type pollService struct {
	answRepo  answers.AnswersRepo
	questRepo questions.QuestionRepo
	votesRepo votes.VotesRepo
}

func NewPollService(cfg *config.Config, db *sql.DB) PollService {
	return pollService{
		answRepo:  answers.NewAnswersRepo(cfg, db),
		questRepo: questions.NewQuestionRepo(cfg, db),
		votesRepo: votes.NewVotesRepo(cfg, db),
	}
}

// TODO: add user id in question -> see man, who created poll
func (p pollService) AddPoll(poll *models.Poll) error {
	err := p.questRepo.AddQuestionPoll(poll.Question)
	if err != nil {
		log.Println("shit 1")
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
