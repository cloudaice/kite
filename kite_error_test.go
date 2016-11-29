package kite

import (
	"errors"
	"testing"
)

func TestKiteError(t *testing.T) {
	e := Error("error")
	if e.s != "KiteError: "+"error" {
		t.Fatalf("TestKite Error")
	}

	if e.Error() != "KiteError: "+"error" {
		t.Fatalf("TestKite Error")
	}

	if e.E(errors.New("error")).Error() != "KiteError: "+"error"+"error" {
		t.Fatalf("TestKite Error")
	}
}
