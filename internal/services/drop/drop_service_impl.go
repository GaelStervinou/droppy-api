package drop

import (
	"errors"
	"go-api/internal/repositories"
	"go-api/internal/storage/postgres/drop"
	"go-api/pkg/errors2"
	"go-api/pkg/model"
	"go-api/pkg/validation"
	"gorm.io/gorm"
)

type DropService struct {
	Repo *repositories.Repositories
}

func (s *DropService) CanCreateDrop(dropNotificationId uint, userId uint) (bool, error) {
	alreadyDropped, err := s.Repo.DropRepository.GetDropByDropNotificationAndUser(dropNotificationId, userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if alreadyDropped != nil {
		return false, errors2.CannotDropError{Reason: "User already dropped this notification"}
	}

	currentNotification, err := s.Repo.DropNotificationRepository.GetCurrentDropNotification()

	if err != nil {
		return false, err
	}

	if currentNotification.GetID() != dropNotificationId {
		return false, errors2.CannotDropError{Reason: "Drop notification is not the current one"}
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
	if can, err := s.CanCreateDrop(args.DropNotificationId, userId); !can || err != nil {
		return nil, err
	}

	if can, err := s.IsValidDropCreation(args); !can || err != nil {
		return nil, err
	}

	statusActive := drop.DropStatusActive{}
	dropNotification, err := s.Repo.DropNotificationRepository.GetNotificationByID(args.DropNotificationId)

	if err != nil {
		return nil, err
	}

	if dropNotification == nil {
		return nil, errors.New("Drop notification not found")
	}

	args.Type = dropNotification.GetType()
	return s.Repo.DropRepository.Create(args.DropNotificationId, args.Type, args.Content, args.Description, userId, statusActive.ToInt(), false)
}
