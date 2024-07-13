package errors2

type NotFoundError struct {
	Entity string
}

func (e NotFoundError) Error() string {
	return e.Entity + " not found"
}