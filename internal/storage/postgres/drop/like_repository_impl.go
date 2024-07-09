package drop

import (
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/model"
	"gorm.io/gorm"
)

var _ model.LikeModel = (*Like)(nil)

type Like struct {
	gorm.Model
	DropID uint      `gorm:"not null"`
	UserID uint      `gorm:"not null"`
	Drop   Drop      `gorm:"foreignKey:DropID;references:ID"`
	User   user.User `gorm:"foreignKey:UserID;references:ID"`
}

func (l *Like) GetID() uint {
	return l.ID
}

func (l *Like) GetDropID() uint {
	return l.DropID
}

func (l *Like) GetDrop() model.DropModel {
	return &l.Drop
}

func (l *Like) GetUserID() uint {
	return l.UserID
}

func (l *Like) GetUser() model.UserModel {
	return &l.User
}

type repoLikePrivate struct {
	db *gorm.DB
}

func NewLikeRepo(db *gorm.DB) model.LikeRepository {
	return &repoLikePrivate{db: db}
}

func (r *repoLikePrivate) CreateLike(dropId uint, userId uint) (model.LikeModel, error) {
	like := Like{
		DropID: dropId,
		UserID: userId,
	}
	if err := r.db.Create(&like).Error; err != nil {
		return nil, err
	}
	return r.GetById(like.ID)
}

func (r *repoLikePrivate) DeleteLike(dropId uint, userId uint) error {
	return r.db.Delete(&Like{}, "drop_id = ? AND user_id = ?", dropId, userId).Error
}

func (r *repoLikePrivate) GetDropTotalLikes(dropId uint) (int, error) {
	var count int64
	if err := r.db.Model(&Like{}).Where("drop_id = ?", dropId).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *repoLikePrivate) GetById(likeId uint) (model.LikeModel, error) {
	var like Like
	if err := r.db.Preload("Drop").Preload("User").First(&like, likeId).Error; err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *repoLikePrivate) LikeExists(dropId uint, userId uint) (bool, error) {
	var count int64
	if err := r.db.Model(&Like{}).Where("drop_id = ? AND user_id = ?", dropId, userId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
