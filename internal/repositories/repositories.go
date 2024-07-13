package repositories

import (
	"go-api/internal/storage/postgres"
	"go-api/pkg/model"
)

type Repositories struct {
	//wg              *sync.WaitGroup
	UserRepository             model.UserRepository
	TokenRepository            model.AuthTokenRepository
	DropRepository             model.DropRepository
	DropNotificationRepository model.DropNotificationRepository
	FollowRepository           model.FollowRepository
	GroupRepository            model.GroupRepository
	GroupMemberRepository      model.GroupMemberRepository
	GroupDropRepository        model.GroupDropRepository
	CommentRepository          model.CommentRepository
	CommentResponseRepository  model.CommentResponseRepository
	LikeRepository             model.LikeRepository
}

func Setup() *Repositories {
	sqlDB := postgres.Connect()

	return &Repositories{
		UserRepository:             postgres.NewUserRepo(sqlDB),
		TokenRepository:            postgres.NewTokenRepo(sqlDB),
		DropRepository:             postgres.NewDropRepo(sqlDB),
		DropNotificationRepository: postgres.NewDropNotifRepo(sqlDB),
		FollowRepository:           postgres.NewFollowRepo(sqlDB),
		GroupRepository:            postgres.NewGroupRepo(sqlDB),
		GroupMemberRepository:      postgres.NewGroupMemberRepo(sqlDB),
		GroupDropRepository:        postgres.NewGroupDropRepo(sqlDB),
		CommentRepository:          postgres.NewCommentRepo(sqlDB),
		CommentResponseRepository:  postgres.NewCommentResponseRepo(sqlDB),
		LikeRepository:             postgres.NewLikeRepo(sqlDB),
	}
}

func (r *Repositories) Disconnect() {
	//r.wg.Done()
}
