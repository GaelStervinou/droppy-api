package validation

import (
	"go-api/pkg/errors2"
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
	if len(args.Username) < 4 || len(args.Username) > 255 {
		finalErrors.Fields["username"] = "Username must be at least 4 characters long and at most 255 characters long"
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

	if args.Lat < -90 || args.Lat > 90 {
		finalErrors.Fields["lat"] = "Invalid latitude"
	}

	if args.Lng < -180 || args.Lng > 180 {
		finalErrors.Fields["lng"] = "Invalid longitude"
	}

	return finalErrors
}

func ValidateGroupCreation(args model.GroupCreationParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}

	if len(args.Name) < 2 || len(args.Name) > 255 {
		finalErrors.Fields["name"] = "Name must be at least 2 character long and at most 255 characters long"
	}

	return finalErrors
}

func ValidateGroupPatch(args model.GroupPatchParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}

	if len(args.Name) < 2 || len(args.Name) > 255 {
		finalErrors.Fields["name"] = "Name must be at least 2 character long and at most 255 characters long"
	}

	return finalErrors
}

func ValidateGroupMemberCreation(args model.GroupMemberCreationParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}

	if slices.Contains([]string{"manager", "member"}, args.Role) == false {
		finalErrors.Fields["role"] = "Invalid role"
	}

	return finalErrors
}

func ValidateCommentCreation(args model.CommentCreationParam) errors2.MultiFieldsError {
	finalErrors := errors2.MultiFieldsError{
		Fields: map[string]string{},
	}

	if len(args.Content) < 1 || len(args.Content) > 255 {
		finalErrors.Fields["content"] = "Content must be at least 1 character long and at most 255 characters long"
	}

	return finalErrors
}
