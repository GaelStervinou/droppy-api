package errors2

type AlreadyReportedError struct {
	Entity string
}

func (e AlreadyReportedError) Error() string {
	return e.Entity + " has already been reported"
}