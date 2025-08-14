package errs

import "net/http"

func ErrInvalidEmail() error {
	return Error{Code: http.StatusInternalServerError, Msg: "bad header/body in email"}
}

func ErrCantCreateVerifLink() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't create link for verification"}
}
