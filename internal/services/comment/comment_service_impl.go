package comment

import (
	"errors"
	"go-api/internal/repositories"
	"go-api/pkg/model"
	"go-api/pkg/validation"
)

type CommentService struct {
	Repo *repositories.Repositories
}

func (s *CommentService) CommentDrop(dropId uint, userID uint, args model.CommentCreationParam) (model.CommentModel, error) {
	canComment, err := s.CanCommentDrop(dropId, userID)

	if err != nil || !canComment {
		return nil, err
	}

	isValid, err := s.IsValidCommentCreation(args)
	if err != nil || !isValid {
		return nil, err
	}

	drop, err := s.Repo.DropRepository.GetDropById(dropId)
	if err != nil {
		return nil, err
	}

	comment, err := s.Repo.CommentRepository.CreateComment(args.Content, userID, drop.GetID())
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) IsValidCommentCreation(args model.CommentCreationParam) (bool, error) {
	validationError := validation.ValidateCommentCreation(args)

	if len(validationError.Fields) > 0 {
		return false, validationError
	}

	return true, nil
}

func (s *CommentService) CanCommentDrop(dropID uint, userID uint) (bool, error) {
	drop, err := s.Repo.DropRepository.GetDropById(dropID)
	if err != nil || nil == drop {
		return false, errors.New("drop not found")
	}

	if drop.GetCreatedBy().GetID() == userID {
		return false, errors.New("cannot comment on own drop")
	}

	lastNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()
	if err != nil || nil == lastNotification {
		return false, errors.New("no drop notifications found")
	}

	if lastNotification.GetID() != drop.GetDropNotificationID() {
		return false, errors.New("drop notification is not current")
	}

	hasUserDropped, err := s.Repo.DropRepository.HasUserDropped(drop.GetDropNotificationID(), userID)

	if err != nil {
		return false, err
	}

	if !hasUserDropped {
		return false, errors.New("you must drop before posting comments")
	}

	isFollowing, err := s.Repo.FollowRepository.IsActiveFollowing(userID, drop.GetCreatedBy().GetID())

	if err != nil {
		return false, err
	}

	if !isFollowing {
		availableGroups, err := s.Repo.GroupDropRepository.GetGroupIdsByDropId(drop.GetID())
		if err != nil {
			return false, err
		}
		areUserInSameGroups, err := s.Repo.GroupMemberRepository.IsUserInGroups(availableGroups, userID)
		if err != nil {
			return false, err
		}
		if !areUserInSameGroups {
			return false, errors.New("you must follow the drop creator are be in the same group before posting comments")
		}
	}

	return true, nil
}

func (s *CommentService) DeleteComment(commentId uint) error {
	return s.Repo.CommentRepository.DeleteComment(commentId)
}
