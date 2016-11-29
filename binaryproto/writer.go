package binaryproto

import (
	"io"
)

type Writer struct {
	W io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		W: w,
	}
}

func (w *Writer) WriteHeader(header *Header) error {
	_, err := w.W.Write(header.BinaryCode())
	return err
}

func (w *Writer) WriteBody(body *Body) error {
	_, err := w.W.Write(body.BinaryCode())
	return err
}

func (w *Writer) Write(data []byte) error {
	_, err := w.W.Write(data)
	return err
}
