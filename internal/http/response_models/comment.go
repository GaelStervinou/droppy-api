package response_models

import "go-api/pkg/model"

type GetCommentResponse struct {
	ID        uint
	Content   string
	CreatedAt int
	CreatedBy GetUserResponseInterface
	Drop      GetDropResponse
}

type GetCommentResponseForDrop struct {
	ID        uint
	Content   string
	CreatedAt int
	CreatedBy GetUserResponseInterface
}

func FormatGetCommentResponse(comment model.CommentModel) GetCommentResponse {
	return GetCommentResponse{
		ID:        comment.GetID(),
		Content:   comment.GetContent(),
		CreatedAt: comment.GetCreatedAt(),
		CreatedBy: FormatGetUserResponse(comment.GetCreatedBy()),
		Drop:      FormatGetDropResponse(comment.GetDrop(), false),
	}
}

func FormatGetCommentResponseForDrop(comment model.CommentModel) GetCommentResponseForDrop {
	return GetCommentResponseForDrop{
		ID:        comment.GetID(),
		Content:   comment.GetContent(),
		CreatedAt: comment.GetCreatedAt(),
		CreatedBy: FormatGetUserResponse(comment.GetCreatedBy()),
	}
}

func FormatGetCommentResponsesForDrop(comments []model.CommentModel) []GetCommentResponseForDrop {
	var result []GetCommentResponseForDrop
	for _, comment := range comments {
		result = append(result)
		result = append(result, FormatGetCommentResponseForDrop(comment))
	}
	return result
}
