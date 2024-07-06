package group

import (
	"go-api/internal/repositories"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
	"go-api/pkg/validation"
)

type GroupService struct {
	Repo *repositories.Repositories
}

func (s *GroupService) CanCreateGroup(userId uint, args model.GroupCreationParam) (bool, error) {
	userGroupsOwned, err := s.Repo.GroupRepository.FindAllByUserId(userId)

	if err != nil {
		return false, err
	}

	if len(userGroupsOwned) >= 5 {
		return false, errors2.CannotCreateGroupError{Reason: "You can only create 5 groups"}
	}
	return true, nil
}

func (s *GroupService) IsValidGroupCreation(args model.GroupCreationParam) (bool, error) {
	validationError := validation.ValidateGroupCreation(args)

	if len(validationError.Fields) > 0 {
		return false, validationError
	}

	return true, nil
}

func (s *GroupService) CreateGroup(userId uint, args model.GroupCreationParam) (model.GroupModel, error) {
	if can, err := s.CanCreateGroup(userId, args); !can || err != nil {
		return nil, err
	}

	if can, err := s.IsValidGroupCreation(args); !can || err != nil {
		return nil, err
	}

	user, err := s.Repo.UserRepository.GetById(userId)
	if err != nil {
		return nil, err
	}

	createdGroup, err := s.Repo.GroupRepository.Create(args.Name, args.Description, args.IsPrivate, args.Picture, user)
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}
