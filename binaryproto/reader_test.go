package binaryproto

import (
	"bytes"
	"strings"
	"testing"
)

func equalBytes(src, dst []byte) bool {
	if len(src) != len(dst) {
		return false
	}

	for idx, b := range src {
		if dst[idx] != b {
			return false
		}
	}

	return true
}

func TestReader(t *testing.T) {
	var data []byte

	head := make([]byte, 6)
	copy(head, testHeaderBinaryCode)
	data = append(data, head...)

	checksum := make([]byte, 16)
	copy(checksum, testCheckSum)
	data = append(data, checksum...)

	body := make([]byte, len(testData))
	copy(body, testData)
	data = append(data, body...)

	buf := bytes.NewReader(data)
	reader := NewReader(buf)

	header, err := reader.ReadHeader()
	if err != nil {
		t.Fatalf("TestReader: ReadHeader Error %q", err)
	}

	if header.Type != 0x01 ||
		header.Define != 0x02 ||
		header.Length != 0x03040506 {
		t.Fatalf("TestReader: invalid header")
	}

	bodys, err := reader.ReadBody(uint32(len(body)+16), true)
	if err != nil {
		t.Fatalf("TestReader: ReadBody Error")
	}

	if !equalBytes(bodys.CheckSum, checksum) {
		t.Fatalf("TestReader: invalid checksum")
	}
	if !equalBytes(bodys.Data, body) {
		t.Fatalf("TestReader: invalid body, %q, %q", bodys.Data, body)
	}

	data = make([]byte, len(testData))
	copy(data, testData)

	buf = bytes.NewReader(data)
	reader = NewReader(buf)

	bodys, _ = reader.ReadBody(uint32(len(testData)), false)
	if len(bodys.CheckSum) != 0 {
		t.Fatalf("TestReader: no checksum error")
	}
	if !equalBytes(bodys.Data, data) {
		t.Fatalf("TestReader: invaild body")
	}

	copy(data, testData)
	buf = bytes.NewReader(data)
	reader = NewReader(buf)
	_, err = reader.ReadBody(uint32(len(testData))+1, false)
	if err == nil {
		t.Fatalf("TestReader: readLength error")
	}

	head = []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	buf = bytes.NewReader(head)
	reader = NewReader(buf)
	_, err = reader.ReadHeader()
	if err == nil {
		t.Fatalf("TestReader: readLength error")
	}

}

func TestReadLength(t *testing.T) {
	buf := strings.NewReader("hello world !!")
	reader := NewReader(buf)
	_, err := reader.readLength(14)
	if err != nil {
		t.Fatalf("TestReader: readLength Error %q", err)
	}
	buf = strings.NewReader("hello world !!")
	reader = NewReader(buf)
	_, err = reader.readLength(15)
	if err == nil {
		t.Fatalf("TestReader: readLength Error %q", err)
	}
}
