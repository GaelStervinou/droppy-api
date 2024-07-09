package response_models

import "go-api/pkg/model"

type GetCommentResponseResponse struct {
	ID        uint
	Content   string
	CreatedAt int
	CreatedBy GetUserResponseInterface
}

func FormatGetCommentResponseResponse(commentResponse model.CommentResponseModel) GetCommentResponseResponse {
	return GetCommentResponseResponse{
		ID:        commentResponse.GetID(),
		Content:   commentResponse.GetContent(),
		CreatedAt: commentResponse.GetCreatedAt(),
		CreatedBy: FormatGetUserResponse(commentResponse.GetCreatedBy()),
	}
}

func FormatGetCommentResponsesResponse(commentResponses []model.CommentResponseModel) []GetCommentResponseResponse {
	var result []GetCommentResponseResponse
	for _, commentResponse := range commentResponses {
		result = append(result, FormatGetCommentResponseResponse(commentResponse))
	}
	return result
}
