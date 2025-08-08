package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
)

type PollHandler interface {
	CreatePoll(w http.ResponseWriter, r *http.Request)
}

func (h handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if w.Header().Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType)
		return
	}

	poll := models.Poll{}
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())
		return
	}

	if err := h.pollServeice.AddPoll(&poll, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"success": "true"})
}
