package binaryproto

var (
	HeaderLengthError = Error("Invalid header length.")
)

type BinaryProtoError struct {
	s string
}

func (err *BinaryProtoError) Error() string {
	return err.s
}

func Error(text string) error {
	return &BinaryProtoError{
		"BinaryProtoError: " + text,
	}
}
