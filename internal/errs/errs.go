package errs

import "fmt"

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (this Error) Error() string {
	return fmt.Sprintf("code: %v, msg: %s", this.Code, this.Msg)
}

func ErrCFGLoad() error {
	return Error{Code: 0, Msg: "can't load cfg file .env"}
}

func ErrDBConnection() error {
	return Error{Code: 0, Msg: "can't connect to database"}
}
