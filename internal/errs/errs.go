package errs

import (
	"fmt"
)

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %v, msg: %s", e.Code, e.Msg)
}
