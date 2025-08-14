package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
)

type AuthHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	VerifyEmail(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)

	RefreshTokens(w http.ResponseWriter, r *http.Request)
}

func (h handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> Register")

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType())

		h.logger.Info("req -> Register -> bad content-type")
		return
	}

	req := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())

		h.logger.Info("req -> Register -> bad json")
		return
	}

	token := models.JWT{}
	var err error

	h.logger.Debug("start users registration")
	token.Access, token.Refresh, err = h.authService.RegisterUser(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)

		h.logger.Info(fmt.Sprintf("req -> Register -> %s", err))
		return
	}

	h.logger.Info("req -> Register -> success")

	h.cookier.SetAuthCookies(w, token.Access, token.Refresh)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (h handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {

}

func (h handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> Login")

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType())

		h.logger.Info("req -> Login -> bad content-type")
		return
	}

	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())

		h.logger.Info("req -> Login -> bad json")
		return
	}

	token := models.JWT{}
	var err error

	token.Access, token.Refresh, err = h.authService.LoginUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(err)

		h.logger.Info(fmt.Sprintf("req -> Login -> %s", err))
		return
	}

	h.logger.Info("req -> Login -> sucess")

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

		h.logger.Info(fmt.Sprintf("req -> Refresh -> %s", err))
		return
	}

	h.logger.Info("req -> Refresh -> success")

	h.cookier.SetAuthCookies(w, token.Access, token.Refresh)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}
