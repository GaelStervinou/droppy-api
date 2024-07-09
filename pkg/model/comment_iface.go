package model

type CommentModel interface {
	GetID() uint
	GetContent() string
	GetCreatedAt() int
	GetCreatedBy() UserModel
	GetDrop() DropModel
	GetResponses() []CommentResponseModel
}

type CommentRepository interface {
	CreateComment(content string, userID uint, dropId uint) (CommentModel, error)
	DeleteComment(commentId uint) error
	GetCommentsByDropId(dropId uint) ([]CommentModel, error)
	GetById(commentId uint) (CommentModel, error)
}

type CommentService interface {
	CommentDrop(dropId uint, userID uint, args CommentCreationParam) (CommentModel, error)
	DeleteComment(commentId uint) error
}

type CommentCreationParam struct {
	Content string `json:"content"`
}
