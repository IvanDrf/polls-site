package errs

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (this Error) Error() string {
	return fmt.Sprintf("code: %v, msg: %s", this.Code, this.Msg)
}

// Can't load config
func ErrCFGLoad() error {
	return Error{Code: -1, Msg: "can't load cfg file .env"}
}

// Can't connect to database
func ErrDBConnection() error {
	return Error{Code: -1, Msg: "can't connect to database"}
}

// Invalid logger level in database
func ErrLoggerLevel() error {
	return Error{Code: -1, Msg: "can't set up logger's level"}
}

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

// Error of starting server, cant start new server
func ErrCantStartServer() error {
	return Error{Code: -1, Msg: "can't start server"}
}

// Error of content-type in request
func ErrInvalidBodyType() error {
	return Error{Code: http.StatusUnsupportedMediaType, Msg: "want json"}
}

// Error of сontent in json
func ErrInvalidBodyReq() error {
	return Error{Code: http.StatusBadRequest, Msg: "invalid json"}
}

// Error of loggin, cant find user in database
func ErrCantFindUser() error {
	return Error{Code: http.StatusUnauthorized, Msg: "user with this email doesn't exist"}
}

// Error of password, incorrect password during logging
func ErrInvalidPswInLog() error {
	return Error{Code: http.StatusUnauthorized, Msg: "incorrect password"}
}

// Error of jwt token, cant create new jwt token
func ErrCantCreateToken() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't create jwt token"}
}

// Error of jwt token, incorrect signing method
func ErrIncorrectJWTMethod() error {
	return Error{Code: http.StatusBadRequest, Msg: "invalid jwt signing method"}
}

// Error of jwt token, invalid token or expired token
func ErrInValidToken() error {
	return Error{Code: http.StatusBadRequest, Msg: "invalid jwt token"}
}
