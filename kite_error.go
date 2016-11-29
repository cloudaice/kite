package kite

var (
	ErrUnexpectedPackage = Error("Unexpected Package.")
)

type KiteError struct {
	s string
}

func (pe *KiteError) Error() string {
	return pe.s
}

func (pe *KiteError) E(err error) error {
	pe.s = pe.s + err.Error()
	return pe
}

func Error(text string) *KiteError {
	return &KiteError{
		"KiteError: " + text,
	}
}
