package errs

import "net/http"

// Error of registration, user already in database
func ErrAlreadyInDB() error {
	return Error{Code: http.StatusUnauthorized, Msg: "user with this email already exists"}
}

// Error of registration, invalid email, doesn't mathc regual expr
func ErrInvalidEmailInReg() error {
	return Error{Code: http.StatusUnauthorized, Msg: "incorrect email during registration"}
}

// Error of registration, invalid passw, have bad symbols
func ErrInvalidPswInReg() error {
	return Error{Code: http.StatusUnauthorized, Msg: "incorrect symbols in password"}
}

// Error of registration, cant register new user
func ErrCantRegister() error {
	return Error{Code: http.StatusUnauthorized, Msg: "can't register new user"}
}

// Error of loggin, cant find user in database
func ErrCantFindUser() error {
	return Error{Code: http.StatusUnauthorized, Msg: "user with this email doesn't exist"}
}

// Error of finding user's id in database
func ErrCantFindUserId() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't find user id in db"}
}

// Error of password, incorrect password during logging
func ErrInvalidPswInLog() error {
	return Error{Code: http.StatusUnauthorized, Msg: "incorrect password"}
}
