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
	PicturePath custom_type.NullString
	CreatedAt   *time.Time
	CreatedBy   GetUserResponseInterface `json:",omitempty"`
}

func FormatGetGroupResponse(group model.GroupModel, createdBy GetUserResponseInterface) GetGroupResponse {
	if nil == group {
		return GetGroupResponse{}
	}

	createdAt := time.Unix(int64(group.GetCreatedAt()), 0)

	picturePath := custom_type.NullString{NullString: group.GetPicturePath()}

	if nil == createdBy {
		return GetGroupResponse{
			ID:          group.GetID(),
			Name:        group.GetName(),
			Description: group.GetDescription(),
			PicturePath: picturePath,
			CreatedAt:   &createdAt,
		}
	}

	return GetGroupResponse{
		ID:          group.GetID(),
		Name:        group.GetName(),
		Description: group.GetDescription(),
		PicturePath: picturePath,
		CreatedAt:   &createdAt,
		CreatedBy:   createdBy,
	}
}
