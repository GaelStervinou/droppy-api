package group

import (
	"go-api/internal/repositories"
	"go-api/internal/storage/postgres/group"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
	"go-api/pkg/validation"
)

type GroupService struct {
	Repo *repositories.Repositories
}

func (s *GroupService) CanCreateGroup(userId uint) (bool, error) {
	userGroupsOwned, err := s.Repo.GroupRepository.FindAllGroupOwnedByUserId(userId)

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
	if can, err := s.CanCreateGroup(userId); !can || err != nil {
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

	role := &group.GroupMemberRoleManager{}
	status := &group.GroupMemberStatusActive{}
	_, err = s.Repo.GroupMemberRepository.Create(createdGroup.GetID(), user.GetID(), role.ToString(), status.ToIntGroupMemberStatus())
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}

func (s *GroupService) PatchGroup(groupId uint, userId uint, args model.GroupPatchParam) (model.GroupModel, error) {
	if can, err := s.CanUpdateGroup(groupId, userId); !can || err != nil {
		return nil, err
	}

	if can, err := s.IsValidGroupUpdate(groupId, args); !can || err != nil {
		return nil, err
	}

	updatedGroup, err := s.Repo.GroupRepository.Update(model.FilledGroupPatchParam{
		ID:          groupId,
		Name:        args.Name,
		Description: args.Description,
		IsPrivate:   args.IsPrivate,
		Picture:     args.Picture,
	})

	if err != nil {
		return nil, err
	}

	if nil == updatedGroup {
		return nil, errors2.CannotUpdateGroupError{Reason: "Group not found"}
	}

	return updatedGroup, nil
}

func (s *GroupService) CanUpdateGroup(groupId uint, userId uint) (bool, error) {
	groupToUpdate, err := s.Repo.GroupRepository.GetById(groupId)
	if err != nil {
		return false, err
	}

	//TODO rajouter les modérateurs qd ils seront dispo
	if groupToUpdate.GetCreatedByID() != userId {
		return false, errors2.CannotUpdateGroupError{Reason: "You are not the owner of the group"}
	}

	return true, nil
}

func (s *GroupService) IsValidGroupUpdate(groupId uint, args model.GroupPatchParam) (bool, error) {
	validationError := validation.ValidateGroupPatch(args)

	if len(validationError.Fields) > 0 {
		return false, validationError
	}

	if res, err := s.Repo.GroupRepository.GetByName(args.Name); res != nil && res.GetID() != groupId && err == nil {
		return false, errors2.CannotUpdateGroupError{Reason: "Group with this name already exists"}
	}

	return true, nil
}
