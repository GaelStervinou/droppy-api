package drop

import (
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/model"
	"gorm.io/gorm"
)

var _ model.CommentModel = (*Comment)(nil)

type Comment struct {
	gorm.Model
	Content     string    `gorm:"not null"`
	CreatedById uint      `gorm:"not null"`
	DropId      uint      `gorm:"not null"`
	CreatedBy   user.User `gorm:"foreignKey:CreatedById;references:ID"`
	Drop        Drop      `gorm:"foreignKey:DropId;references:ID"`
}

func (c *Comment) GetID() uint {
	return c.ID
}

func (c *Comment) GetContent() string {
	return c.Content
}

func (c *Comment) GetCreatedAt() int {
	return int(c.CreatedAt.Unix())
}

func (c *Comment) GetCreatedBy() model.UserModel {
	return &c.CreatedBy
}

func (c *Comment) GetDrop() model.DropModel {
	return &c.Drop
}

type repoCommentPrivate struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) model.CommentRepository {
	return &repoCommentPrivate{db: db}
}

func (r *repoCommentPrivate) CreateComment(content string, userID uint, dropId uint) (model.CommentModel, error) {
	comment := Comment{
		Content:     content,
		CreatedById: userID,
		DropId:      dropId,
	}
	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	return r.GetById(comment.ID)
}

func (r *repoCommentPrivate) DeleteComment(commentId uint) error {
	return r.db.Delete(&Comment{}, commentId).Error
}

func (r *repoCommentPrivate) GetCommentsByDropId(dropId uint) ([]model.CommentModel, error) {
	var comments []Comment
	if err := r.db.Preload("CreatedBy").Preload("Drop").Where("drop_id = ?", dropId).Find(&comments).Error; err != nil {
		return nil, err
	}
	var result []model.CommentModel
	for _, comment := range comments {
		result = append(result, &comment)
	}
	return result, nil
}

func (r *repoCommentPrivate) GetById(commentId uint) (model.CommentModel, error) {
	var comment Comment
	if err := r.db.Preload("CreatedBy").Preload("Drop").Preload("Drop.CreatedBy").First(&comment, commentId).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}
