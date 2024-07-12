package follow

import (
	"go-api/internal/repositories"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
)

type FollowService struct {
	Repo *repositories.Repositories
}

func (s *FollowService) GetUserFollowing(userID uint, requesterID uint) ([]model.FollowModel, error) {
	targetedUser, err := s.Repo.UserRepository.GetById(userID)
	if err != nil {
		return nil, err
	}

	if targetedUser == nil {
		return nil, nil
	}

	if targetedUser.IsPrivateUser() && requesterID != userID {
		currentUserIsFollowing, err := s.Repo.FollowRepository.IsActiveFollowing(requesterID, userID)
		if err != nil {
			return nil, err
		}

		if !currentUserIsFollowing {
			return nil, errors2.NotAllowedError{Reason: "This user is private"}
		}
	}
	return s.Repo.FollowRepository.GetFollowing(userID)
}

func (s *FollowService) GetUserFollowers(userID uint, requesterID uint) ([]model.FollowModel, error) {
	targetedUser, err := s.Repo.UserRepository.GetById(userID)
	if err != nil {
		return nil, err
	}

	if targetedUser == nil {
		return nil, nil
	}

	if targetedUser.IsPrivateUser() && requesterID != userID {
		currentUserIsFollowing, err := s.Repo.FollowRepository.IsActiveFollowing(requesterID, userID)
		if err != nil {
			return nil, err
		}

		if !currentUserIsFollowing {
			return nil, errors2.NotAllowedError{Reason: "This user is private"}
		}
	}

	return s.Repo.FollowRepository.GetFollowers(userID)
}
