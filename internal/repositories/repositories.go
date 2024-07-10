package repositories

import (
	"go-api/internal/storage/postgres"
	"go-api/pkg/model"
	"os"
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
	CommentRepository          model.CommentRepository
	CommentResponseRepository  model.CommentResponseRepository
	LikeRepository             model.LikeRepository
}

func Setup() *Repositories {
	sqlDB, err := postgres.Connect()
	if err != nil {
		os.Exit(1)
	}

	return &Repositories{
		//wg:              wg,
		UserRepository:             postgres.NewUserRepo(sqlDB),
		TokenRepository:            postgres.NewTokenRepo(sqlDB),
		DropRepository:             postgres.NewDropRepo(sqlDB),
		DropNotificationRepository: postgres.NewDropNotifRepo(sqlDB),
		FollowRepository:           postgres.NewFollowRepo(sqlDB),
		GroupRepository:            postgres.NewGroupRepo(sqlDB),
		GroupMemberRepository:      postgres.NewGroupMemberRepo(sqlDB),
		CommentRepository:          postgres.NewCommentRepo(sqlDB),
		CommentResponseRepository:  postgres.NewCommentResponseRepo(sqlDB),
		LikeRepository:             postgres.NewLikeRepo(sqlDB),
	}
}

func (r *Repositories) Disconnect() {
	//r.wg.Done()
}
