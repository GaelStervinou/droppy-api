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

	return true, nil
}

func (s *CommentService) DeleteComment(commentId uint) error {
	return s.Repo.CommentRepository.DeleteComment(commentId)
}

func (s *CommentService) CanDeleteComment(author uint, user uint) error {
	comment, err := s.Repo.CommentRepository.GetById(author)
	if err != nil || nil == comment {
		return errors.New("comment not found")
	}

	connectedUser, err := s.Repo.UserRepository.GetById(user)
	if err != nil || nil == connectedUser {
		return errors.New("user not found")
	}

	if comment.GetCreatedBy().GetID() != user || connectedUser.GetRole() != "admin" {
		return errors.New("unauthorized")
	}

	return nil
}
