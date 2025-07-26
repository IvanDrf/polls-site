package checker

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
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
	'?': struct{}{},
	'#': struct{}{},
	'<': struct{}{},
	'>': struct{}{},
	'%': struct{}{},
	'@': struct{}{},
	'/': struct{}{},
	';': struct{}{},
}

func (this checker) ValidPassword(passw string) bool {
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

type PswHasher interface {
	HashPassword(passw string) string
	ComparePassword(hashed, passw string) bool
}

type hasher struct {
}

func NewPswHasher() PswHasher {
	return hasher{}
}

const hashLen = 14

func (this hasher) HashPassword(passw string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(passw), hashLen)
	return string(bytes)
}

func (this hasher) ComparePassword(hashed, passw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(passw)) == nil
}
