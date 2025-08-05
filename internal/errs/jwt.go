package errs

import "net/http"

// Error of jwt token, incorrect signing method
func ErrIncorrectJWTMethod() error {
	return Error{Code: http.StatusBadRequest, Msg: "invalid jwt signing method"}
}

// Error of jwt token, invalid token or expired token
func ErrInValidToken() error {
	return Error{Code: http.StatusBadRequest, Msg: "invalid jwt token"}
}

// Error of jwt token, cant add refresh token in database
func ErrCantAddToken() error {
	return Error{Code: http.StatusInternalServerError, Msg: "can't add refresh token to database"}
}

// Error of jwt token, cant find access token in cookies or header 'authorization'
func ErrCantFindToken() error {
	return Error{Code: http.StatusUnauthorized, Msg: "can't find access token in header/cookies"}
}
