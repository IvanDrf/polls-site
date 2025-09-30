package test

import "testing"

func Assert[T comparable](t *testing.T, got T, excepted T) {
	if excepted != got {
		t.Errorf("TEST FAILED ASSERT(got, excepted) %v != %v", got, excepted)
	}
}

func NotEqual[T comparable](t *testing.T, got T, excepted T) {
	if excepted == got {
		t.Errorf("TEST FAILED NOT_EQUAL(got, excepted), %v == %v", got, excepted)
	}
}
