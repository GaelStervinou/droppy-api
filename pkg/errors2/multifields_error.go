package errors2

import "errors"

type MultiFieldsError struct {
	Fields map[string]string `json:"fields"`
}

func (e MultiFieldsError) Error() string {
	finalError := ""
	for k, v := range e.Fields {
		finalError += k + ": " + v + "\n"
	}

	return finalError
}

func (e MultiFieldsError) Is(target error) bool {
	var multiFieldsError MultiFieldsError
	ok := errors.As(target, &multiFieldsError)
	return ok
}
