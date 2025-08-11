package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
)

type VoteHandler interface {
	VoteInPoll(w http.ResponseWriter, r *http.Request)
	DeleteVoteInPoll(w http.ResponseWriter, r *http.Request)
}

func (h handler) VoteInPoll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> VoteInPoll")

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType())

		h.logger.Info("req -> VoteInPoll -> bad content-type")
		return
	}

	vote := models.Vote{}
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())

		h.logger.Info("req -> VoteInPoll -> bad json")
		return
	}

	votes, err := h.voteService.VoteInPoll(&vote, r)
	if err != nil {
		w.WriteHeader(errs.GetCode(err))

		json.NewEncoder(w).Encode(err)

		h.logger.Info(fmt.Sprintf("req -> VoteInPoll -> %s", err))
		return
	}

	h.logger.Info("req -> VoteInPoll -> success")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(votes)
}

func (h handler) DeleteVoteInPoll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("req -> DeleteVote")

	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyType())

		h.logger.Info("req -> DeleteVote -> bad content-type")
		return
	}

	vote := models.Vote{}
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(errs.ErrInvalidBodyReq())

		h.logger.Info("req -> DeleteVote -> bad json")
		return
	}

	voteRes, err := h.voteService.DeleteVoteInPoll(&vote, r)
	if err != nil {
		w.WriteHeader(errs.GetCode(err))

		json.NewEncoder(w).Encode(err)

		h.logger.Info(fmt.Sprintf("req -> DeleteVote -> %s", err))
		return
	}

	h.logger.Info("req -> DeleteVote -> success")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(voteRes)
}
