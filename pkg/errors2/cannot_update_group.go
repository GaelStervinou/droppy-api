package errors2

type CannotUpdateGroupError struct {
	Reason string
}

func (e CannotUpdateGroupError) Error() string {
	return e.Reason
}
