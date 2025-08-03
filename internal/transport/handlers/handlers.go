package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/transport/auth"
	"github.com/IvanDrf/polls-site/internal/transport/auth/cookies"
)

type Handler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	authService auth.Auther
	cookier     cookies.Cookier

	logger *slog.Logger
}

func NewHandler(cfg *config.Config, db *sql.DB, logger *slog.Logger) Handler {
	return handler{
		authService: auth.NewAuthService(cfg, db),
		cookier:     cookies.NewCookier(),
		logger:      logger,
	}
}

func (hand handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	hand.logger.Info("register req")

	w.Header().Set("Content-Type", "application/json")

	if w.Header().Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType())
		return
	}

	req := models.UserReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())
		return
	}

	token := models.JWT{}
	var err error

	hand.logger.Debug("start users registration")
	token.Access, token.Refresh, err = hand.authService.RegisterUser(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)
		return
	}

	hand.logger.Debug("end user registration")

	hand.cookier.SetAuthCookies(w, token.Access, token.Refresh)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (hand handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if w.Header().Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType())
		return
	}

	user := models.UserReq{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())
		return
	}

	token := models.JWT{}
	var err error

	token.Access, token.Refresh, err = hand.authService.LoginUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)
		return
	}

	hand.cookier.SetAuthCookies(w, token.Access, token.Refresh)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
