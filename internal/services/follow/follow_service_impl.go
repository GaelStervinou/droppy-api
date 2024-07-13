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

func (s *FollowService) DeleteFollow(requesterID uint, followID uint) error {
	follow, err := s.Repo.FollowRepository.GetFollowByID(followID)
	if err != nil {
		return err
	}

	if follow == nil {
		return errors2.NotAllowedError{Reason: "Follow not found"}
	}

	if follow.GetFollowerID() != requesterID && follow.GetFollowedID() != requesterID {
		return errors2.NotAllowedError{Reason: "You are not allowed to delete this follow"}
	}

	return s.Repo.FollowRepository.Delete(follow.GetID())
}
