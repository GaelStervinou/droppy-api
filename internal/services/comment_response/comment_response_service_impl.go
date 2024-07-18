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
