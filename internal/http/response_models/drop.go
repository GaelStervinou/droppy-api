package response_models

import (
	"go-api/pkg/model"
	"time"
)

type GetDropResponse struct {
	ID          uint
	Type        string
	Content     string
	Description string
	Lat         *float64
	Lng         *float64
	PicturePath *string
	CreatedAt   *time.Time
	CreatedBy   GetUserResponseInterface    `json:",omitempty"`
	Comments    []GetCommentResponseForDrop `json:",omitempty"`
}

func FormatGetDropResponse(drop model.DropModel) GetDropResponse {
	if nil == drop {
		return GetDropResponse{}
	}
	lat := drop.GetLat()
	latPointer := &lat
	lng := drop.GetLng()
	lngPointer := &lng
	picturePath := drop.GetPicturePath()
	picturePathPointer := &picturePath

	createdAt := time.Unix(int64(drop.GetCreatedAt()), 0)

	return GetDropResponse{
		ID:          drop.GetID(),
		Type:        drop.GetType(),
		Content:     drop.GetContent(),
		Description: drop.GetDescription(),
		Lat:         latPointer,
		Lng:         lngPointer,
		PicturePath: picturePathPointer,
		CreatedAt:   &createdAt,
		CreatedBy:   FormatGetUserResponse(drop.GetCreatedBy()),
		Comments:    FormatGetCommentResponsesForDrop(drop.GetComments()),
	}
}
