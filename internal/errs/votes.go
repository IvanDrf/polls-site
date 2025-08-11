package errs

import "net/http"

func ErrCantAddVote() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't add vote"}
}

func ErrAlreadyVoted() error {
	return Error{Code: http.StatusForbidden, Msg: "already voted"}
}

func ErrCantCountVotes() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't count votes in poll"}
}

func ErrCantDeleteAllVotes() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't delete all votes for question"}
}
