package response_models

import "go-api/pkg/model"

type GetDropNotificationResponse struct {
	ID        uint
	Type      string
	CreatedAt string
}

func FormatGetDropNotificationResponse(dropNotification model.DropNotificationModel) GetDropNotificationResponse {
	if nil == dropNotification {
		return GetDropNotificationResponse{}
	}

	return GetDropNotificationResponse{
		ID:        dropNotification.GetID(),
		Type:      dropNotification.GetType(),
		CreatedAt: dropNotification.GetCreatedAt(),
	}
}
