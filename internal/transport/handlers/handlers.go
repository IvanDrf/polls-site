package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/service"
	"github.com/IvanDrf/polls-site/internal/transport/auth"
	"github.com/IvanDrf/polls-site/internal/transport/auth/cookies"
)

type Handler interface {
	AuthHandler
	PollHandler
}

type handler struct {
	authService  auth.Auther
	pollServeice service.PollService

	cookier cookies.Cookier

	logger *slog.Logger
}

func NewHandler(cfg *config.Config, db *sql.DB, logger *slog.Logger) Handler {
	return handler{
		authService:  auth.NewAuthService(cfg, db, logger),
		pollServeice: service.NewPollService(cfg, db),
		cookier:      cookies.NewCookier(),
		logger:       logger,
	}
}
