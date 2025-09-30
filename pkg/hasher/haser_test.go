package hasher

import (
	"sync"
	"testing"

	"github.com/IvanDrf/polls-site/pkg/test"
)

func TestNewHasher(t *testing.T) {
	hasher := NewPswHasher()

	test.NotEqual(t, hasher, nil)
}

var passwords = []string{
	"12345",
	"printF",
	"qwertYpl",
}

func TestHashPassword(t *testing.T) {
	hasher := NewPswHasher()

	wg := new(sync.WaitGroup)

	wg.Add(len(passwords))
	for _, password := range passwords {
		go func() {
			defer wg.Done()
			test.NotEqual(t, hasher.HashPassword(password), password)
		}()
	}

	wg.Wait()
}

func TestComparePassword(t *testing.T) {
	hasher := NewPswHasher()

	wg := new(sync.WaitGroup)

	wg.Add(len(passwords))
	for _, password := range passwords {
		go func() {
			defer wg.Done()
			hashed := hasher.HashPassword(password)
			test.Assert(t, hasher.ComparePassword(hashed, password), true)
		}()
	}

	wg.Wait()
}
