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
	"github.com/IvanDrf/polls-site/internal/repo/transaction"
	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type PollService interface {
	AddPoll(poll *models.Poll, r *http.Request) (models.PollId, error)
	DeletePoll(poll *models.Poll) error

	VoteInPoll(vote *models.Vote, r *http.Request) (models.PollRes, error)
}

type pollService struct {
	answRepo   answers.AnswersRepo
	questRepo  questions.QuestionRepo
	votesRepo  votes.VotesRepo
	tokensRepo tokens.TokensRepo

	transaction transaction.Transactioner

	jwter jwter.Jwter
}

func NewPollService(cfg *config.Config, db *sql.DB) PollService {
	return pollService{
		answRepo:   answers.NewAnswersRepo(cfg, db),
		questRepo:  questions.NewQuestionRepo(cfg, db),
		votesRepo:  votes.NewVotesRepo(cfg, db),
		tokensRepo: tokens.NewTokensRepo(cfg, db),

		transaction: transaction.NewTransactioner(cfg, db),

		jwter: jwter.NewJwter(cfg),
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

func (p pollService) DeletePoll(poll *models.Poll) error {
	question, err := p.questRepo.FindQuestionPollById(poll.Id)
	if err != nil {
		return errs.ErrCantFindQuestion()
	}

	p.transaction.StartTransaction()

	err = p.answRepo.DeleteAnswers(poll.Answers, question.Id)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantAddAnswer()
	}

	err = p.questRepo.DeleteQuestionPollById(question.Id)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantDeleteQuestion()
	}

	p.transaction.CommitTransaction()

	return nil
}

// TODO: add check for question and answers id in databases
func (p pollService) VoteInPoll(vote *models.Vote, r *http.Request) (models.PollRes, error) {
	token, err := p.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return models.PollRes{}, err
	}

	vote.UserId, err = p.tokensRepo.FindUserId(token)
	if err != nil {
		return models.PollRes{}, errs.ErrCantFindUserId()
	}

	_, err = p.votesRepo.FindVote(vote.QuestionId, vote.UserId)
	if err == nil {
		return models.PollRes{}, errs.ErrAlreadyVoted()
	}

	_, err = p.votesRepo.AddVote(vote)
	if err != nil {
		return models.PollRes{}, errs.ErrCantAddVote()
	}

	res, err := p.votesRepo.CountVotes(vote.QuestionId)
	if err != nil {
		return models.PollRes{}, errs.ErrCantCountVotes()
	}

	return res, nil
}
