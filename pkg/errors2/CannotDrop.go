package errors2

type CannotDropError struct {
	Reason string
}

func (e CannotDropError) Error() string {
	return e.Reason
}
