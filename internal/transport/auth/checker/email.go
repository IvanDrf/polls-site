package checker

import "regexp"

type EmailChecker interface {
	ValidEmail(em string) bool
}

type emailChecker struct {
}

func NewEmailChecker() EmailChecker {
	return emailChecker{}
}

const re = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func (c emailChecker) ValidEmail(em string) bool {
	return regexp.MustCompile(re).MatchString(em)
}
