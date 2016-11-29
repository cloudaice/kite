package binaryproto

import (
	"io"
)

type Reader struct {
	R io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		R: r,
	}
}

func (r *Reader) ReadHeader() (*Header, error) {
	head, err := r.readLength(6)
	if err != nil {
		return nil, err
	}
	return ParseHeader(head)
}

// ReadBody read n bytes, read checksum wether check is true
func (r *Reader) ReadBody(n uint32, check bool) (*Body, error) {
	body, err := r.readLength(n)
	if err != nil {
		return nil, err
	}

	var (
		checksum []byte
		data     []byte
	)
	if check {
		checksum = body[:16]
		data = body[16:]
	} else {
		checksum = []byte{}
		data = body[:]
	}

	return &Body{checksum, data}, nil
}

// readLength read n bytes from io.Reader, the type of n is uint32,
// because the most length in header is 4-bytes uint32
func (r *Reader) readLength(n uint32) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(r.R, buf)
	return buf, err
}
