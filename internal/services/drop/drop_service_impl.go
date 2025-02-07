package drop

import (
	"errors"
	"go-api/internal/repositories"
	"go-api/internal/storage/postgres"
	"go-api/pkg/errors2"
	"go-api/pkg/file"
	"go-api/pkg/model"
	"go-api/pkg/validation"
	"gorm.io/gorm"
	"slices"
)

type DropService struct {
	Repo *repositories.Repositories
}

func (s *DropService) CanCreateDrop(userId uint) (bool, error) {
	currentNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()
	if err != nil {
		return false, err
	}
	alreadyDropped, err := s.Repo.DropRepository.GetDropByDropNotificationAndUser(currentNotification.GetID(), userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if alreadyDropped != nil {
		return false, errors2.CannotDropError{Reason: "User already dropped this notification"}
	}

	return true, nil
}

func (s *DropService) IsValidDropCreation(args model.DropCreationParam) (bool, error) {
	validationError := validation.ValidateDropCreation(args)

	if len(validationError.Fields) > 0 {
		return false, validationError
	}

	return true, nil
}

func (s *DropService) CreateDrop(userId uint, args model.DropCreationParam) (model.DropModel, error) {
	if can, err := s.CanCreateDrop(userId); !can || err != nil {
		return nil, err
	}

	if can, err := s.IsValidDropCreation(args); !can || err != nil {
		return nil, err
	}

	currentDropNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()
	if err != nil {
		return nil, err
	}

	var picturePath string
	if args.Picture != nil {
		picturePath, err = file.UploadFile(args.Picture)
		if err != nil {
			return nil, err
		}
	}

	filledDrop := model.FilledDropCreation{
		Type:               currentDropNotification.GetType(),
		Content:            args.Content,
		ContentTile:        args.ContentTitle,
		ContentSubTitle:    args.ContentSubTitle,
		ContentPicturePath: args.ContentPicturePath,
		Description:        args.Description,
		DropNotificationId: currentDropNotification.GetID(),
		PicturePath:        picturePath,
		Lat:                args.Lat,
		Lng:                args.Lng,
		Location:           args.Location,
	}

	statusActive := postgres.DropStatusActive{}

	createdDrop, err := s.Repo.DropRepository.Create(
		filledDrop.DropNotificationId,
		filledDrop.Type,
		filledDrop.Content,
		filledDrop.Description,
		filledDrop.ContentPicturePath,
		filledDrop.ContentTile,
		filledDrop.ContentSubTitle,
		userId,
		statusActive.ToInt(),
		false,
		filledDrop.PicturePath,
		filledDrop.Lat,
		filledDrop.Lng,
		filledDrop.Location,
	)

	if err != nil {
		return nil, err
	}

	user, err := s.Repo.UserRepository.GetById(userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	for _, group := range user.GetGroups() {
		if slices.Contains(args.Groups, group.GetID()) {
			_, err = s.Repo.GroupDropRepository.Create(createdDrop.GetID(), group.GetID())
			if err != nil {
				return nil, err
			}
		}
	}

	return s.Repo.DropRepository.GetDropById(createdDrop.GetID())
}

func (s *DropService) GetUserFeed(userId uint) ([]model.DropModel, error) {
	isActiveUser, err := s.Repo.UserRepository.IsActiveUser(userId)

	if err != nil {
		return nil, err
	}

	if !isActiveUser {
		return nil, errors.New("User is not active")
	}

	lastDropNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()

	if err != nil {
		return nil, err
	}

	if lastDropNotification == nil {
		return nil, errors.New("No drop notifications found")
	}

	followingUsers, err := s.Repo.FollowRepository.GetFollowing(userId)

	if err != nil {
		return nil, err
	}

	var followingUserIds []uint
	for _, user := range followingUsers {
		followingUserIds = append(followingUserIds, user.GetFollowedID())
	}

	followingUserIds = append(followingUserIds, userId)

	drops, err := s.Repo.DropRepository.GetDropsByUserIdsAndDropNotificationId(followingUserIds, lastDropNotification.GetID())

	if err != nil {
		return nil, err
	}

	return drops, nil
}

func (s *DropService) GetDropById(dropID uint, requesterID uint) (model.DropModel, error) {
	return s.Repo.DropRepository.GetDropById(dropID)
}

func (s *DropService) GetDropsByUserId(userId uint, currentUser model.UserModel) ([]model.DropModel, error) {
	isActiveUser, err := s.Repo.UserRepository.IsActiveUser(userId)

	if err != nil {
		return nil, err
	}

	if !isActiveUser {
		return nil, errors.New("User not found")
	}

	user, err := s.Repo.UserRepository.GetById(userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("User not found")
	}

	if user.IsPrivateUser() {
		if nil == currentUser {
			return nil, nil
		}

		currentUserIsFollowing, err := s.Repo.FollowRepository.IsActiveFollowing(currentUser.GetID(), userId)

		if err != nil {
			return nil, err
		}

		if !currentUserIsFollowing {
			return nil, nil
		}
	}

	drops, err := s.Repo.DropRepository.GetUserDrops(userId)

	if err != nil {
		return nil, err
	}

	return drops, nil
}

func (s *DropService) HasUserDroppedToday(userId uint) (bool, error) {
	currentDropNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()

	if err != nil {
		return false, err
	}

	if currentDropNotification == nil {
		return false, errors.New("no drop notifications found")
	}

	hasDropped, err := s.Repo.DropRepository.HasUserDropped(currentDropNotification.GetID(), userId)

	if err != nil {
		return false, err
	}

	return hasDropped, nil
}

func (s *DropService) IsCurrentUserLiking(dropId uint, userId uint) (bool, error) {
	likeExists, err := s.Repo.LikeRepository.LikeExists(dropId, userId)

	if err != nil {
		return false, err
	}

	return likeExists, nil
}

func (s *DropService) DeleteDrop(dropID uint, requesterID uint) error {
	drop, err := s.Repo.DropRepository.GetDropById(dropID)

	if err != nil {
		return err
	}

	if drop == nil {
		return errors.New("drop not found")
	}

	if drop.GetCreatedById() != requesterID {
		return errors2.NotAllowedError{Reason: "This drop is not yours"}
	}

	return s.Repo.DropRepository.Delete(dropID)
}

func (s *DropService) PatchDrop(dropID uint, requesterID uint, patch model.DropPatch) (model.DropModel, error) {
	drop, err := s.Repo.DropRepository.GetDropById(dropID)

	if err != nil {
		return nil, err
	}

	if drop == nil {
		return nil, errors.New("drop not found")
	}

	if drop.GetCreatedById() != requesterID {
		return nil, errors2.NotAllowedError{Reason: "This drop is not yours"}
	}

	updates := make(map[string]interface{})

	if patch.IsPinned {
		updates["IsPinned"] = patch.IsPinned
	} else {
		updates["IsPinned"] = false
	}

	if len(updates) == 0 {
		return nil, errors.New("no updates")
	}

	updatedDrop, err := s.Repo.DropRepository.Update(dropID, updates)

	if err != nil {
		return nil, err
	}

	return s.Repo.DropRepository.GetDropById(updatedDrop.GetID())
}
