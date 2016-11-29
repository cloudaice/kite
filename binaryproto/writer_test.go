package binaryproto

import (
	"bytes"
	"testing"
)

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	var err error
	w := NewWriter(&buf)
	err = w.WriteHeader(testHeader)
	if err != nil {
		t.Fatalf("TestWriter: WriteHeader Error")
	}
	if buf.Len() != 6 {
		t.Fatalf("TestWriter: WriteHeader buf Error")
	}

	err = w.WriteBody(testBody)
	if err != nil {
		t.Fatalf("TestWriter: WriteBody Error")
	}
	if buf.Len() != 6+16+len(testData) {
		t.Fatalf("TestWriter: WriteBody buf Error")
	}
}

func TestWrite(t *testing.T) {
	testData := []byte("abcdef")
	var buf bytes.Buffer
	var err error
	w := NewWriter(&buf)
	err = w.Write(testData)
	if err != nil {
		t.Fatalf("TestWrite Error")
	}

	if buf.Len() != len(testData) {
		t.Fatalf("TestWrite Error")
	}
}
