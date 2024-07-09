package model

type CommentResponseModel interface {
	GetID() uint
	GetContent() string
	GetCreatedAt() int
	GetCreatedBy() UserModel
	GetComment() CommentModel
}

type CommentResponseRepository interface {
	CreateCommentResponse(content string, userID uint, commentId uint) (CommentResponseModel, error)
	DeleteCommentResponse(commentResponseId uint) error
	GetById(commentResponseId uint) (CommentResponseModel, error)
}

type CommentResponseService interface {
	RespondToComment(commentId uint, userID uint, args CommentCreationParam) (CommentResponseModel, error)
	DeleteCommentResponse(commentResponseId uint, userID uint) error
}
