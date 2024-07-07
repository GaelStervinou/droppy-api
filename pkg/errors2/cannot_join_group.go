package errors2

type CannotJoinGroupError struct {
	Reason string
}

func (e CannotJoinGroupError) Error() string {
	return e.Reason
}
