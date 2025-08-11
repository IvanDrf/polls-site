package errs

import "net/http"

func ErrCantAddVote() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't add vote"}
}

func ErrDidntVote() error {
	return Error{Code: http.StatusNotFound, Msg: "this user didnt vote in this poll"}
}

func ErrAlreadyVoted() error {
	return Error{Code: http.StatusForbidden, Msg: "already voted"}
}

func ErrCantCountVotes() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't count votes in poll"}
}

func ErrCantDeleteVote() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't delete vote"}
}

func ErrCantDeleteAllVotes() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't delete all votes for question"}
}
