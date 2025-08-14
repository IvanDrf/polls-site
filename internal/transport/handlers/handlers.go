package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/service/auth"
	"github.com/IvanDrf/polls-site/internal/service/poll"
	"github.com/IvanDrf/polls-site/internal/service/vote"
	"github.com/IvanDrf/polls-site/internal/transport/auth/cookies"
)

type Handler interface {
	AuthHandler
	PollHandler
	VoteHandler
}

type handler struct {
	authService auth.Auther
	pollService poll.PollService
	voteService vote.VoteService

	cookier cookies.Cookier

	logger *slog.Logger
}

func NewHandler(cfg *config.Config, db *sql.DB, logger *slog.Logger) Handler {
	return handler{
		authService: auth.NewAuthService(cfg, db, logger),
		pollService: poll.NewPollService(cfg, db, logger),
		voteService: vote.NewVoteService(cfg, db, logger),
		cookier:     cookies.NewCookier(),
		logger:      logger,
	}
}
