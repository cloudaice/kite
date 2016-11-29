package kite

var (
	ErrCheck       = PError("CheckSum Fail: ")
	ErrEncryption  = PError("Encryption Fail: ")
	ErrCompression = PError("compression Fail: ")

	ErrWrite = PError("WritePackage Fail: ")
	ErrRead  = PError("ReadPackage Fail: ")
)

type ProtoError struct {
	s string
}

func (pe *ProtoError) Error() string {
	return pe.s
}

func (pe *ProtoError) E(err error) error {
	pe.s = pe.s + err.Error()
	return pe
}

func PError(text string) *ProtoError {
	return &ProtoError{
		"ProtoError: " + text,
	}
}
