package comment_response

import (
	"errors"
	"go-api/internal/repositories"
	"go-api/pkg/model"
	"go-api/pkg/validation"
)

type CommentResponseService struct {
	Repo *repositories.Repositories
}

func (s *CommentResponseService) RespondToComment(commentId uint, userID uint, args model.CommentCreationParam) (model.CommentResponseModel, error) {
	canCommentResponse, err := s.CanCommentResponse(commentId, userID)

	if err != nil || !canCommentResponse {
		return nil, err
	}

	isValid, err := s.IsValidCommentResponseCreation(args)
	if err != nil || !isValid {
		return nil, err
	}
	commentResponse, err := s.Repo.CommentResponseRepository.CreateCommentResponse(args.Content, commentId, userID)
	if err != nil {
		return nil, err
	}

	createdCommentResponse, err := s.Repo.CommentResponseRepository.GetById(commentResponse.GetID())

	if err != nil {
		return nil, err
	}

	return createdCommentResponse, nil
}

func (s *CommentResponseService) DeleteCommentResponse(commentResponseId uint, userID uint) error {
	canDelete, err := s.CanDeleteCommentResponse(commentResponseId, userID)

	if err != nil || !canDelete {
		return err
	}

	return s.Repo.CommentResponseRepository.DeleteCommentResponse(commentResponseId)
}

func (s *CommentResponseService) CanCommentResponse(commentId uint, userID uint) (bool, error) {
	comment, err := s.Repo.CommentRepository.GetById(commentId)
	if err != nil {
		return false, err
	}

	if nil == comment {
		return false, errors.New("comment not found")
	}

	lastNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()
	if err != nil || nil == lastNotification {
		return false, errors.New("no drop notifications found")
	}

	if lastNotification.GetID() != comment.GetDrop().GetDropNotificationID() {
		return false, errors.New("drop notification is not current")
	}

	hasUserDropped, err := s.Repo.DropRepository.HasUserDropped(comment.GetDrop().GetDropNotificationID(), userID)

	if err != nil {
		return false, err
	}

	if !hasUserDropped {
		return false, errors.New("you must drop before posting comments")
	}

	isFollowing, err := s.Repo.FollowRepository.IsActiveFollowing(userID, comment.GetDrop().GetCreatedBy().GetID())

	if err != nil {
		return false, err
	}

	if !isFollowing {
		availableGroups, err := s.Repo.GroupDropRepository.GetGroupIdsByDropId(comment.GetDrop().GetID())
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

func (s *CommentResponseService) IsValidCommentResponseCreation(args model.CommentCreationParam) (bool, error) {
	validationError := validation.ValidateCommentCreation(args)

	if len(validationError.Fields) > 0 {
		return false, validationError
	}

	return true, nil
}

func (s *CommentResponseService) CanDeleteCommentResponse(commentResponseId uint, userID uint) (bool, error) {
	commentResponse, err := s.Repo.CommentResponseRepository.GetById(commentResponseId)
	if err != nil {
		return false, err
	}

	if nil == commentResponse {
		return false, errors.New("comment response not found")
	}

	isOwner := commentResponse.GetCreatedBy().GetID() != userID
	isDropOwner := commentResponse.GetComment().GetDrop().GetCreatedBy().GetID() != userID

	if isOwner || isDropOwner {
		return true, nil
	}

	return false, errors.New("you must be the owner of the comment response or the drop owner")
}
