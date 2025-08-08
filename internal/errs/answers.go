package errs

import "net/http"

// Error of adding answer in poll in database
func ErrCantAddAnswer() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't add answer in db"}
}

// Error of deletion answer in database
func ErrCantDeleteAnswer() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't delete answer in db"}
}
