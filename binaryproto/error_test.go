package binaryproto

import (
	"testing"
)

func TestBinaryProtoError(t *testing.T) {
	errorText := "Test Error"
	err := Error(errorText)
	if err.Error() != "BinaryProtoError: "+errorText {
		t.Fatalf("TestBinaryError Error")
	}
}
