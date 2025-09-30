package checker

import (
	"testing"

	"github.com/IvanDrf/polls-site/pkg/test"
)

type emailTest struct {
	email  string
	status bool
}

func TestNewEmailChecker(t *testing.T) {
	checker := NewEmailChecker()

	test.NotEqual(t, checker, nil)
}

func TestValidEmail(t *testing.T) {
	emails := []emailTest{
		{"normal@gmail.com", true},
		{"bademail.com", false},
		{"", false},
		{"bad!email@gmail.com", false},
		{"-?bad@mail.com", false},
		{"normal@mail.ru", true},
	}

	emailChecker := NewEmailChecker()

	for _, email := range emails {
		test.Assert(t, emailChecker.ValidEmail(email.email), email.status)
	}
}
