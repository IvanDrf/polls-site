package errs

import (
	"net/http"
)

// Error of сontent in json
func ErrInvalidBodyReq() error {
	return Error{Code: http.StatusBadRequest, Msg: "invalid json"}
}

// Error of content-type in request
func ErrInvalidBodyType() error {
	return Error{Code: http.StatusUnsupportedMediaType, Msg: "want json"}
}

// Error of adding question in poll
func ErrCantAddQuestion() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't add question in db"}
}

// Error of adding answer in poll
func ErrCantAddAnswer() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't add answer in db"}
}

func ErrCantDeleteAnswer() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't delete answer in db"}
}

func ErrCantFindUserId() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't find user id in db"}
}
