package group

import (
	"errors"
	"fmt"
	"go-api/internal/repositories"
	grouprepository "go-api/internal/storage/postgres/group"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
	"go-api/pkg/validation"
	"slices"
)

type GroupMemberService struct {
	Repo *repositories.Repositories
}

func (s *GroupMemberService) JoinGroup(currentUserId uint, userId uint, args model.GroupMemberCreationParam) (model.GroupMemberModel, error) {
	if currentUserId == userId {
		role := &grouprepository.GroupMemberRoleMember{}
		args.Role = role.ToString()
	}

	if can, err := s.IsValidGroupMemberCreation(args); !can || err != nil {
		return nil, err
	}

	if can, err := s.CanJoinGroup(userId, args.GroupID); !can || err != nil {
		return nil, err
	}

	targetedGroup, err := s.Repo.GroupRepository.GetById(args.GroupID)
	if err != nil {
		return nil, err
	}

	if targetedGroup == nil {
		return nil, errors.New(fmt.Sprintf("Group with id %d not found", args.GroupID))
	}

	var status grouprepository.GroupMemberStatus
	if targetedGroup.IsPrivateGroup() {
		status = &grouprepository.GroupMemberStatusPending{}
	} else {
		status = &grouprepository.GroupMemberStatusActive{}
	}
	_, err = s.Repo.GroupMemberRepository.Create(args.GroupID, userId, args.Role, status.ToIntGroupMemberStatus())

	if err != nil {
		return nil, err
	}

	finalGroupMember, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(args.GroupID, userId)

	if err != nil {
		return nil, err
	}

	return finalGroupMember, nil
}
func (s *GroupMemberService) DeleteGroupMember(actionRequesterID uint, groupID uint, memberID uint) error {
	groupMember, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupID, memberID)

	if err != nil {
		return err
	}

	if groupMember == nil {
		return errors.New(fmt.Sprintf("Group member with id %d not found", memberID))
	}

	actionRequester, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupID, actionRequesterID)

	if err != nil {
		return err
	}

	if actionRequester == nil {
		return errors.New(fmt.Sprintf("Group member with id %d not found", actionRequesterID))
	}

	canMakeAction, err := s.CanMakeActionOnUser(actionRequester, groupMember)

	if actionRequester.GetMemberID() != groupMember.GetMemberID() && !canMakeAction {
		return errors.New("You are not allowed to make this action")
	}
	err = s.Repo.GroupMemberRepository.Delete(groupID, memberID)
	return err
}
func (s *GroupMemberService) AcceptGroupMember(userId uint, groupID uint, memberID uint) (model.GroupMemberModel, error) {
	groupMember, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupID, memberID)

	if err != nil {
		return nil, err
	}

	if groupMember == nil {
		return nil, errors.New(fmt.Sprintf("Group member with id %d not found", memberID))
	}

	actionRequester, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupID, userId)

	if err != nil {
		return nil, err
	}

	if actionRequester == nil {
		return nil, errors.New(fmt.Sprintf("Group member with id %d not found", userId))
	}

	canMakeAction, err := s.CanMakeActionOnUser(actionRequester, groupMember)

	if err != nil {
		return nil, err
	}

	if !canMakeAction {
		return nil, errors.New("You are not allowed to make this action")
	}

	pendingStatus := &grouprepository.GroupMemberStatusPending{}
	if pendingStatus.ToIntGroupMemberStatus() != groupMember.GetStatus() {
		return nil, errors.New("group member is not pending")
	}

	activeStatus := &grouprepository.GroupMemberStatusActive{}
	groupMember, err = s.Repo.GroupMemberRepository.UpdateStatus(groupID, memberID, activeStatus.ToIntGroupMemberStatus())

	finalGroupMember, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupID, memberID)

	if err != nil {
		return nil, err
	}

	return finalGroupMember, nil
}

func (s *GroupMemberService) UpdateGroupMemberRole(requesterId uint, groupID uint, memberID uint, args model.GroupMemberPatchParam) (model.GroupMemberModel, error) {
	targetedGroup, err := s.Repo.GroupRepository.GetById(groupID)

	if err != nil {
		return nil, err
	}

	if targetedGroup == nil {
		return nil, errors.New(fmt.Sprintf("Group with id %d not found", groupID))
	}

	groupMember, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupID, memberID)

	if err != nil {
		return nil, err
	}

	if groupMember == nil {
		return nil, errors.New(fmt.Sprintf("Group member with id %d not found", memberID))
	}

	requester, err := s.Repo.GroupMemberRepository.GetByGroupIDAndMemberID(groupID, requesterId)

	if err != nil {
		return nil, err
	}

	if requester == nil {
		return nil, errors.New(fmt.Sprintf("Requester with id %d not found in group %d", requesterId, groupID))
	}

	canMakeAction, err := s.CanMakeActionOnUser(requester, groupMember)

	if err != nil {
		return nil, errors2.NotAllowedError{
			Reason: err.Error(),
		}
	}

	if !canMakeAction {
		return nil, errors2.NotAllowedError{
			Reason: "You are not allowed to make this action",
		}
	}

	if !slices.Contains(grouprepository.GroupMemberRoles(), args.Role) {
		return nil, errors.New("Invalid role")
	}

	groupMember, err = s.Repo.GroupMemberRepository.UpdateRole(groupID, memberID, args.Role)

	return groupMember, err
}

func (s *GroupMemberService) IsValidGroupMemberCreation(args model.GroupMemberCreationParam) (bool, error) {
	validationError := validation.ValidateGroupMemberCreation(args)

	if len(validationError.Fields) > 0 {
		return false, validationError
	}

	return true, nil
}

func (s *GroupMemberService) CanJoinGroup(userId uint, groupID uint) (bool, error) {
	isGroupMember, _ := s.Repo.GroupMemberRepository.IsGroupMember(groupID, userId)

	if isGroupMember {
		return false, errors2.CannotJoinGroupError{Reason: "User already joined this group"}
	}

	return true, nil
}

func (s *GroupMemberService) CanMakeActionOnUser(actionRequester model.GroupMemberModel, groupMember model.GroupMemberModel) (bool, error) {
	if groupMember.GetMemberID() == groupMember.GetGroup().GetCreatedByID() && groupMember.GetMemberID() != actionRequester.GetMemberID() {
		return false, errors2.NotAllowedError{Reason: "You can not make this action"}
	}

	managerRole := &grouprepository.GroupMemberRoleManager{}

	if actionRequester.GetRole() != managerRole.ToString() {
		return false, errors.New("You are not a manager")
	}

	return true, nil
}

func (s *GroupMemberService) FindAllUserGroups(userId uint) ([]model.GroupModel, error) {
	groupMembers, err := s.Repo.GroupMemberRepository.GetByMemberID(userId)

	if err != nil {
		return nil, err
	}

	var groups []model.GroupModel
	for _, groupMember := range groupMembers {
		groups = append(groups, groupMember.GetGroup())
	}

	return groups, nil
}
