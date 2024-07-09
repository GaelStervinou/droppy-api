package like

import (
	"errors"
	"go-api/internal/repositories"
	"go-api/pkg/model"
)

type LikeService struct {
	Repo *repositories.Repositories
}

func (s *LikeService) LikeDrop(userID uint, args model.LikeParam) (model.LikeModel, error) {
	canLike, err := s.CanLikeDrop(userID, args)
	if err != nil {
		return nil, err
	}
	if !canLike {
		return nil, errors.New("user already liked this drop")
	}
	return s.Repo.LikeRepository.CreateLike(args.DropId, userID)
}

func (s *LikeService) UnlikeDrop(userID uint, args model.LikeParam) error {
	return s.Repo.LikeRepository.DeleteLike(args.DropId, userID)
}

func (s *LikeService) CanLikeDrop(userID uint, args model.LikeParam) (bool, error) {
	dropExists, err := s.Repo.DropRepository.DropExists(args.DropId)
	if err != nil {
		return false, err
	}

	if !dropExists {
		return false, errors.New("drop not found")
	}

	canLike, err := s.Repo.LikeRepository.LikeExists(args.DropId, userID)
	if err != nil {
		return false, err
	}

	return !canLike, nil
}
