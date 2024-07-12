package response_models

import (
	"go-api/pkg/model"
	"time"
)

type GetOnePendingFollowResponse struct {
	ID        uint                     `json:"id"`
	Follower  GetUserResponseInterface `json:"follower"`
	CreatedAt *time.Time               `json:"createdAt"`
	Status    uint                     `json:"status"`
}

func FormatGetOnePendingFollowResponse(follow model.FollowModel) GetOnePendingFollowResponse {
	if nil == follow {
		return GetOnePendingFollowResponse{}
	}

	createdAt := time.Unix(int64(follow.GetCreatedAt()), 0)
	return GetOnePendingFollowResponse{
		ID:        follow.GetID(),
		Follower:  FormatGetUserResponse(follow.GetFollower()),
		CreatedAt: &createdAt,
		Status:    follow.GetStatus(),
	}
}

type GetOneFollowResponse struct {
	ID        uint                     `json:"id"`
	Follower  GetUserResponseInterface `json:"follower"`
	Followed  GetUserResponseInterface `json:"followed"`
	CreatedAt *time.Time               `json:"createdAt"`
	Status    uint                     `json:"status"`
}

func FormatGetOneFollowResponse(follow model.FollowModel) GetOneFollowResponse {
	if nil == follow {
		return GetOneFollowResponse{}
	}

	createdAt := time.Unix(int64(follow.GetCreatedAt()), 0)
	return GetOneFollowResponse{
		ID:        follow.GetID(),
		Follower:  FormatGetUserResponse(follow.GetFollower()),
		Followed:  FormatGetUserResponse(follow.GetFollowed()),
		CreatedAt: &createdAt,
		Status:    follow.GetStatus(),
	}
}
