package errors2

type NotAllowedError struct {
	Reason string
}

func (e NotAllowedError) Error() string {
	return e.Reason
}
