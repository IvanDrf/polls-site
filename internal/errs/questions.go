package errs

import "net/http"

// Error of adding question in poll
func ErrCantAddQuestion() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't add question in db"}
}

// Error of finding question in database
func ErrCantFindQuestion() error {
	return Error{Code: http.StatusNotFound, Msg: "can't find question in db"}
}

// Error of deletion question in database
func ErrCantDeleteQuestion() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't delete question in db"}
}
