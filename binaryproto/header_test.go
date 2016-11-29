package binaryproto

import (
	"testing"
)

var (
	testHeader = &Header{
		0x01,
		0x02,
		0x03040506,
	}

	testHeaderBinaryCode = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

	testCheckSum = []byte("0123456789abcdef")
	testData     = []byte("girls like flowers")

	testBody = &Body{testCheckSum, testData}
)

func TestParseHeader(t *testing.T) {
	head := make([]byte, 6)
	copy(head, testHeaderBinaryCode)

	header, err := ParseHeader(head)
	if err != nil {
		t.Fatalf("%q %s", head, err.Error())
	}

	if header.Type != 0x01 {
		t.Fatalf("TestParseHeader Error header Type %q", header.Type)
	}

	if header.Define != 0x02 {
		t.Fatalf("TestParseHeader Error header Define  %q", header.Define)
	}

	if header.Length != 0x03040506 {
		t.Fatalf("TestParseHeader Error header Length %q", header.Length)
	}
	head = append(head, 0x00)
	header, err = ParseHeader(head)
	if err == nil {
		t.Fatalf("TestParseHeader Error ParseHeader")
	}
}

func TestHeader(t *testing.T) {
	header := NewHeader(0x01, 0x02, 0x01020304)
	head := header.BinaryCode()
	cmpHead := []byte{0x01, 0x02, 0x01, 0x02, 0x03, 0x04}
	for idx, b := range head {
		if cmpHead[idx] != b {
			t.Fatalf("TestHeader Error")
		}
	}
}

func TestBody(t *testing.T) {
	body := &Body{
		[]byte("love"),
		[]byte("a lovely day"),
	}

	data := body.BinaryCode()
	cmpBody := []byte("lovea lovely day")
	for idx, b := range data {
		if cmpBody[idx] != b {
			t.Fatalf("TestBody Error")
		}
	}
	if int(body.Length()) != len([]byte("lovea lovely day")) {
		t.Fatalf("TestBdoy Length Error")
	}
	body = &Body{
		nil,
		nil,
	}
	_ = body.BinaryCode()
}
