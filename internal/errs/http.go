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
