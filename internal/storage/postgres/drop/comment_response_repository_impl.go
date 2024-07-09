package drop

import (
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/model"
	"gorm.io/gorm"
)

var _ model.CommentResponseModel = (*CommentResponse)(nil)

type CommentResponse struct {
	gorm.Model
	Content     string    `gorm:"not null"`
	CreatedById uint      `gorm:"not null"`
	CommentId   uint      `gorm:"not null"`
	CreatedBy   user.User `gorm:"foreignKey:CreatedById;references:ID"`
	Comment     Comment   `gorm:"foreignKey:CommentId;references:ID"`
}

func (c *CommentResponse) GetID() uint {
	return c.ID
}

func (c *CommentResponse) GetContent() string {
	return c.Content
}

func (c *CommentResponse) GetCreatedAt() int {
	return int(c.CreatedAt.Unix())
}

func (c *CommentResponse) GetCreatedBy() model.UserModel {
	return &c.CreatedBy
}

func (c *CommentResponse) GetComment() model.CommentModel {
	return &c.Comment
}

type repoCommentResponsePrivate struct {
	db *gorm.DB
}

func NewCommentResponseRepo(db *gorm.DB) model.CommentResponseRepository {
	return &repoCommentResponsePrivate{db: db}
}

func (r *repoCommentResponsePrivate) CreateCommentResponse(content string, commentId uint, userID uint) (model.CommentResponseModel, error) {
	commentResponse := CommentResponse{
		Content:     content,
		CreatedById: userID,
		CommentId:   commentId,
	}
	if err := r.db.Create(&commentResponse).Error; err != nil {
		return nil, err
	}
	return r.GetById(commentResponse.ID)
}

func (r *repoCommentResponsePrivate) DeleteCommentResponse(commentResponseId uint) error {
	return r.db.Delete(&CommentResponse{}, commentResponseId).Error
}

func (r *repoCommentResponsePrivate) GetById(commentResponseId uint) (model.CommentResponseModel, error) {
	var commentResponse CommentResponse
	if err := r.db.Preload("CreatedBy").Preload("Comment").Preload("Comment.Drop").First(&commentResponse, commentResponseId).Error; err != nil {
		return nil, err
	}
	return &commentResponse, nil
}
