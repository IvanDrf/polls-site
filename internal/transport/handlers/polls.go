package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
)

type PollHandler interface {
	CreatePoll(w http.ResponseWriter, r *http.Request)
	DeletePoll(w http.ResponseWriter, r *http.Request)
}

func (h handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> CreatePoll")

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType)

		h.logger.Info("req -> CreatePoll -> bad content-type")
		return
	}

	poll := models.Poll{}
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())

		h.logger.Info("req -> CreatePoll - > bad json")
		return
	}

	pollId, err := h.pollService.AddPoll(&poll, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.GetCode(err))

		json.NewEncoder(w).Encode(err)

		h.logger.Info("req -> CreatePoll -> ", "error", err)
		return
	}

	h.logger.Info("req -> CreatePoll -> sucess")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pollId)
}

func (h handler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> DeletePoll")

	if r.Header.Get("Content-Type") != "application/json" {
		h.logger.Error(w.Header().Get("Content-Type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType())

		h.logger.Info("req -> DeletePoll -> bad content-type")
		return
	}

	poll := models.Poll{}
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())

		h.logger.Info("req -> DeletePoll -> bad json")
		return
	}

	if err := h.voteService.DeleteAllVotesInPoll(&poll, r); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.GetCode(err))

		json.NewEncoder(w).Encode(err)

		h.logger.Info("req -> DeletePoll -> can't delete votes")
		return
	}

	if err := h.pollService.DeletePoll(&poll, r); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errs.GetCode(err))

		json.NewEncoder(w).Encode(err)

		h.logger.Info("req -> DeletePoll -> ", "error", err)
		return
	}

	h.logger.Info("req -> DeletePoll -> success")

	w.WriteHeader(http.StatusNoContent)
}
