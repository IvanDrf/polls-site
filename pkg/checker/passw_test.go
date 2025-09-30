package checker

import (
	"testing"

	"github.com/IvanDrf/polls-site/pkg/test"
)

type passwords struct {
	passw  string
	status bool
}

func TestNewPSWChecker(t *testing.T) {
	checker := NewPSWChecker()

	test.NotEqual(t, checker, nil)
}

func TestValidPassword(t *testing.T) {
	passwords := []passwords{
		{"", false},
		{"123456789", false},
		{"abcdf", false},
		{"printF2f!", true},
		{"scanf_3L", true},
	}

	pswChecker := NewPSWChecker()

	for _, password := range passwords {
		test.Assert(t, pswChecker.ValidPassword(password.passw), password.status)
	}
}
