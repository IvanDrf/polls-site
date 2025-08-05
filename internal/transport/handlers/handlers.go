package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/transport/auth"
	"github.com/IvanDrf/polls-site/internal/transport/auth/cookies"
)

type Handler interface {
	AuthHandler
	Private(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	authService auth.Auther
	cookier     cookies.Cookier

	logger *slog.Logger
}

func NewHandler(cfg *config.Config, db *sql.DB, logger *slog.Logger) Handler {
	return handler{
		authService: auth.NewAuthService(cfg, db, logger),
		cookier:     cookies.NewCookier(),
		logger:      logger,
	}
}

func (h handler) Private(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello private\n"))
}
