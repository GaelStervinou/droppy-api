package group

import (
	"errors"
	"go-api/internal/repositories"
	"go-api/internal/storage/postgres"
	"go-api/pkg/errors2"
	"go-api/pkg/file"
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

	if len(userGroupsOwned) >= 200 {
		return false, errors2.CannotCreateGroupError{Reason: "You can only create 200 groups"}
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
	isPrivate := false
	if args.IsPrivate != nil && *args.IsPrivate {
		isPrivate = true
	}

	createdGroup, err := s.Repo.GroupRepository.Create(args.Name, args.Description, isPrivate, args.PicturePath, user)
	if err != nil {
		return nil, err
	}

	role := &postgres.GroupMemberRoleManager{}
	status := &postgres.GroupMemberStatusActive{}
	_, err = s.Repo.GroupMemberRepository.Create(createdGroup.GetID(), user.GetID(), role.ToString(), status.ToIntGroupMemberStatus())
	if err != nil {
		return nil, err
	}

	return s.Repo.GroupRepository.GetById(createdGroup.GetID())
}

func (s *GroupService) PatchGroup(groupId uint, userId uint, args model.GroupPatchParam) (model.GroupModel, error) {
	if can, err := s.CanUpdateGroup(groupId, userId); !can || err != nil {
		return nil, err
	}

	if can, err := s.IsValidGroupUpdate(groupId, args); !can || err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if args.Name != "" {
		updates["Name"] = args.Name
	}

	if args.Description != "" {
		updates["Description"] = args.Description
	}

	if args.IsPrivate != nil {
		if *args.IsPrivate {
			updates["IsPrivate"] = true
		} else {
			updates["IsPrivate"] = false
		}
	}

	if args.Picture != nil {
		filePath, err := file.UploadFile(args.Picture)
		if err != nil {
			return nil, err
		}
		updates["PicturePath"] = filePath
		args.PicturePath = filePath
	}

	updatedGroup, err := s.Repo.GroupRepository.Update(groupId, updates)

	if err != nil {
		return nil, err
	}

	if nil == updatedGroup {
		return nil, errors2.CannotUpdateGroupError{Reason: "Group not found"}
	}

	return s.Repo.GroupRepository.GetById(updatedGroup.GetID())
}

func (s *GroupService) CanUpdateGroup(groupId uint, userId uint) (bool, error) {
	requester, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupId, userId)
	if err != nil {
		return false, err
	}

	if requester == nil {
		return false, errors2.NotAllowedError{Reason: "You are not a manager"}
	}
	groupToUpdate, err := s.Repo.GroupRepository.GetById(groupId)
	if err != nil {
		return false, err
	}

	managerRole := &postgres.GroupMemberRoleManager{}
	if requester.GetRole() != managerRole.ToString() && groupToUpdate.GetCreatedByID() != userId {
		return false, errors2.NotAllowedError{Reason: "You are not a manager"}
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

func (s *GroupService) GetGroupDrops(groupId uint, requesterID uint) ([]model.DropModel, error) {
	targetedGroup, err := s.Repo.GroupRepository.GetById(groupId)
	if err != nil {
		return nil, err
	}

	if targetedGroup == nil {
		return nil, errors.New("Group not found")
	}

	if targetedGroup.IsPrivateGroup() {
		requester, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupId, requesterID)
		if err != nil {
			return nil, err
		}

		if requester == nil {
			return nil, errors2.NotAllowedError{Reason: "You are not a member of this group"}
		}
	}

	lastNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()
	if err != nil {
		return nil, err
	}

	if lastNotification == nil {
		return nil, errors.New("no drop notifications found")
	}
	groupDrops, err := s.Repo.GroupDropRepository.GetByGroupIdAndLastNotificationId(groupId, lastNotification.GetID())

	if err != nil {
		return nil, err
	}

	var drops []model.DropModel
	for _, gd := range groupDrops {
		drops = append(drops, gd.GetDrop())
	}

	return drops, nil
}

func (s *GroupService) DeleteGroup(groupId uint, userId uint) error {
	groupToDelete, err := s.Repo.GroupRepository.GetById(groupId)
	if err != nil {
		return err
	}

	if groupToDelete == nil {
		return errors.New("group not found")
	}

	user, err := s.Repo.UserRepository.GetById(userId)
	if err != nil {
		return err
	}

	if groupToDelete.GetCreatedByID() != userId || user.GetRole() != "admin" {
		return errors2.NotAllowedError{Reason: "You are not allowed to delete this group"}
	}

	err = s.Repo.GroupRepository.DeleteGroup(groupId)
	if err != nil {
		return err
	}

	return s.Repo.GroupMemberRepository.DeleteGroupMembers(groupId)
}
