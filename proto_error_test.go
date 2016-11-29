package kite

import (
	"errors"
	"testing"
)

func TestProtoError(t *testing.T) {
	e := PError("error")
	if e.s != "ProtoError: "+"error" {
		t.Fatalf("TestProto Error")
	}

	if e.Error() != "ProtoError: "+"error" {
		t.Fatalf("TestProto Error")
	}

	if e.E(errors.New("error")).Error() != "ProtoError: "+"error"+"error" {
		t.Fatalf("TestProto Error")
	}
}
