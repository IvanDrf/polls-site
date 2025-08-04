package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
)

type AuthHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)

	RefreshTokens(w http.ResponseWriter, r *http.Request)
}

func (h handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> Register")

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

	h.logger.Debug("start users registration")
	token.Access, token.Refresh, err = h.authService.RegisterUser(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)
		return
	}

	h.logger.Debug("end user registration")

	h.cookier.SetAuthCookies(w, token.Access, token.Refresh)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (h handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> Login")

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

	token.Access, token.Refresh, err = h.authService.LoginUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)
		return
	}

	h.cookier.SetAuthCookies(w, token.Access, token.Refresh)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (h handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> Refresh")

	w.Header().Set("Content-Type", "application/json")

	token := models.JWT{}
	var err error

	token.Access, token.Refresh, err = h.authService.RefreshTokens(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(err)
		return
	}

	h.cookier.SetAuthCookies(w, token.Access, token.Refresh)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
