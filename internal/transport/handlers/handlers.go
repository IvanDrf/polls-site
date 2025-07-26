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
)

type Handler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	authService auth.Auther
	logger      *slog.Logger
}

func NewHandler(cfg *config.Config, db *sql.DB, logger *slog.Logger) Handler {
	return handler{
		authService: auth.NewAuthService(cfg, db),
		logger:      logger,
	}
}

func (this handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
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

	if err := this.authService.RegisterUser(&req); err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "true"})
}

func (this handler) LoginUser(w http.ResponseWriter, r *http.Request) {
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
	token.Token, err = this.authService.LoginUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
