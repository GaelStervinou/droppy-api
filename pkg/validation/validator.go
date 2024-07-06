package validation

import (
	"errors"
	"go-api/internal/services/drop_type"
	errors2 "go-api/pkg/errors2"
	"go-api/pkg/model"
	"net/mail"
	"slices"
)

func ValidateUserCreation(args model.UserCreationParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}
	if _, err := mail.ParseAddress(args.Email); err != nil {
		finalErrors.Fields["email"] = "Invalid email address"
	}
	if len(args.Password) < 8 {
		finalErrors.Fields["password"] = "Password must be at least 8 characters long"
	}
	if len(args.Username) < 4 {
		finalErrors.Fields["username"] = "Username must be at least 4 characters long"
	}
	if slices.Contains([]string{"user", "admin"}, args.Role) == false {
		finalErrors.Fields["role"] = "Invalid role"
	}
	if len(args.Username) < 4 {
		finalErrors.Fields["username"] = "Username must be at least 4 characters long"
	}

	return finalErrors
}

func ValidateUserPatch(args model.UserPatchParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}
	if _, err := mail.ParseAddress(args.Email); err != nil {
		finalErrors.Fields["email"] = "Invalid email address"
	}
	if len(args.Username) < 4 {
		finalErrors.Fields["username"] = "Username must be at least 4 characters long"
	}

	return finalErrors
}

func ValidateDropCreation(args model.DropCreationParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}

	if len(args.Content) < 1 {
		finalErrors.Fields["content"] = "Content must be at least 1 character long"
	}

	if len(args.Description) < 1 || len(args.Description) > 255 {
		finalErrors.Fields["description"] = "Description must be at least 1 character long and at most 255 characters long"
	}

	if args.Lat < -90 || args.Lat > 90 {
		finalErrors.Fields["lat"] = "Invalid latitude"
	}

	if args.Lng < -180 || args.Lng > 180 {
		finalErrors.Fields["lng"] = "Invalid longitude"
	}

	/*validTypes := []string{"youtube", "spotify", "film"}

	if slices.Contains(validTypes, args.Type) == false {
		finalErrors.Fields["type"] = "Invalid type"
	}*/

	/*err := validateContentByType(args.Content, args.Type)

	if err != nil {
		finalErrors.Fields["content"] = err.Error()
	}*/

	return finalErrors
}

func validateContentByType(content string, dropType string) error {
	dropTypeFactory := drop_type.NewDropTypeFactory()

	dropTypeInstance := dropTypeFactory.CreateDropType(dropType)

	if dropTypeInstance == nil {
		return errors.New("invalid drop type")
	}

	if false == dropTypeInstance.IsValidContent(content) {
		return errors.New("invalid content")
	}

	return nil
}

func ValidateGroupCreation(args model.GroupCreationParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}

	if len(args.Name) < 2 && len(args.Name) > 255 {
		finalErrors.Fields["name"] = "Name must be at least 2 character long and at most 255 characters long"
	}

	if len(args.Description) < 1 && len(args.Description) > 255 {
		finalErrors.Fields["description"] = "Description must be at least 1 character long and at most 255 characters long"
	}

	return finalErrors
}

func ValidateGroupPatch(args model.GroupPatchParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}

	if len(args.Name) < 2 && len(args.Name) > 255 {
		finalErrors.Fields["name"] = "Name must be at least 2 character long and at most 255 characters long"
	}

	if len(args.Description) < 1 && len(args.Description) > 255 {
		finalErrors.Fields["description"] = "Description must be at least 1 character long and at most 255 characters long"
	}

	return finalErrors
}
