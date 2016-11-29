package binaryproto

import (
	"encoding/binary"
)

type Header struct {
	Type   uint8
	Define uint8
	Length uint32
}

func NewHeader(tp, df uint8, lth uint32) *Header {
	return &Header{
		Type:   tp,
		Define: df,
		Length: lth,
	}
}

// ParseHeader parse []byte to a Header, []byte length must be 6.
// first two bytes are type, define, last 4 bytes is a big endian uint32
func ParseHeader(head []byte) (*Header, error) {
	if len(head) != 6 {
		return nil, HeaderLengthError
	}
	tp, df := uint8(head[0]), uint8(head[1])
	lth := binary.BigEndian.Uint32(head[2:6])
	return &Header{tp, df, lth}, nil
}

func (kh *Header) BinaryCode() []byte {
	header := make([]byte, 6)
	header[0], header[1] = byte(kh.Type), byte(kh.Define)
	binary.BigEndian.PutUint32(header[2:], kh.Length)
	return header
}

// CheckSum length is 128 bit, use md5 function to checksum
type Body struct {
	CheckSum []byte
	Data     []byte
}

func (kb *Body) Length() uint32 {
	return uint32(len(kb.BinaryCode()))
}

func (kb *Body) BinaryCode() []byte {
	var body []byte
	body = append(body, kb.CheckSum...)
	body = append(body, kb.Data...)
	return body
}
