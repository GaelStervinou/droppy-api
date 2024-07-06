package response_models

import (
	"go-api/pkg/custom_type"
	"go-api/pkg/model"
	"time"
)

type GetGroupResponse struct {
	ID          uint
	Name        string
	Description string
	IsPrivate   bool
	PicturePath custom_type.NullString
	CreatedAt   *time.Time
	CreatedBy   GetUserResponseInterface `json:",omitempty"`
}

func FormatGetGroupResponse(group model.GroupModel) GetGroupResponse {
	if nil == group {
		return GetGroupResponse{}
	}

	createdAt := time.Unix(int64(group.GetCreatedAt()), 0)

	picturePath := custom_type.NullString{NullString: group.GetPicturePath()}

	return GetGroupResponse{
		ID:          group.GetID(),
		Name:        group.GetName(),
		Description: group.GetDescription(),
		IsPrivate:   group.IsPrivateGroup(),
		PicturePath: picturePath,
		CreatedAt:   &createdAt,
		CreatedBy:   FormatGetUserResponse(group.GetCreatedBy()),
	}
}
