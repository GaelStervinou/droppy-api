package errors2

type CannotCreateGroupError struct {
	Reason string
}

func (e CannotCreateGroupError) Error() string {
	return e.Reason
}
