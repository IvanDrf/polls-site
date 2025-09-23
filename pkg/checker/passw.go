package checker

import (
	"unicode"
)

type PswChecker interface {
	ValidPassword(passw string) bool
}

type checker struct {
}

func NewPSWChecker() PswChecker {
	return checker{}
}

var invalidSymbols = map[rune]struct{}{
	'?': {},
	'#': {},
	'<': {},
	'>': {},
	'%': {},
	'@': {},
	'/': {},
	';': {},
}

func (c checker) ValidPassword(passw string) bool {
	if len(passw) < 5 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasInvalid bool
	for _, val := range passw {
		switch {
		case unicode.IsLower(val):
			hasLower = true

		case unicode.IsUpper(val):
			hasUpper = true

		case unicode.IsDigit(val):
			hasNumber = true

		default:
			if _, ok := invalidSymbols[val]; ok {
				hasInvalid = true
			}

		}
	}

	return hasUpper && hasLower && hasNumber && !hasInvalid
}
