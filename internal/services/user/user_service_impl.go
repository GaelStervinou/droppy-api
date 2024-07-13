package user

import (
	"go-api/internal/repositories"
	"go-api/pkg/file"
	"go-api/pkg/model"
	"go-api/pkg/validation"
)

type UserService struct {
	Repo *repositories.Repositories
}

func NewUserService(repo *repositories.Repositories) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) UpdateUser(userId uint, userToPatch model.UserPatchParam) (model.UserModel, error) {
	validationError := validation.ValidateUserPatch(userToPatch)

	if len(validationError.Fields) > 0 {
		return nil, validationError
	}

	updates := make(map[string]interface{})

	if userToPatch.Bio != "" {
		updates["Bio"] = userToPatch.Bio
	}

	if userToPatch.Username != "" {
		updates["Username"] = userToPatch.Username
	}

	if userToPatch.Picture != nil {
		filePath, err := file.UploadFile(userToPatch.Picture)
		if err != nil {
			return nil, err
		}
		updates["Avatar"] = filePath
		userToPatch.PicturePath = filePath
	}

	if userToPatch.IsPrivate != nil {
		if *userToPatch.IsPrivate {
			updates["IsPrivate"] = true
		} else {
			updates["IsPrivate"] = false
		}
	}

	updatedUser, err := s.Repo.UserRepository.Update(userId, updates)

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
