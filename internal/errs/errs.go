package errs

import (
	"fmt"
	"strconv"
	"strings"
)

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %v, msg: %s", e.Code, e.Msg)
}

func GetCode(e error) int {
	s := strings.Split(e.Error(), " ")
	if len(s) < 2 {
		return 400
	}

	strCode := s[1][:len(s[1])-1]
	code, _ := strconv.Atoi(strCode)

	return code
}
