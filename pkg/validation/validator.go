package validation

import (
	errors2 "go-api/pkg/errors2"
	"go-api/pkg/model"
	"net/mail"
)

func ValidateUserCreation(args model.UserCreationParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}
	if _, err := mail.ParseAddress(args.Email); err != nil {
		finalErrors.Fields["email"] = "Invalid email address"
	}
	if len(args.Firstname) < 2 {
		finalErrors.Fields["firstname"] = "Firstname must be at least 2 characters long"
	}
	if len(args.Lastname) < 2 {
		finalErrors.Fields["lastname"] = "Lastname must be at least 2 characters long"
	}
	if len(args.Password) < 8 {
		finalErrors.Fields["password"] = "Password must be at least 8 characters long"
	}
	if len(args.Username) < 4 {
		finalErrors.Fields["username"] = "Username must be at least 4 characters long"
	}
	if len(args.Roles) == 0 {
		finalErrors.Fields["role"] = "Role must not be empty"
	}
	if len(args.Username) < 4 {
		finalErrors.Fields["username"] = "Username must be at least 4 characters long"
	}

	return finalErrors
}
