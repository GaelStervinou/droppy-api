package repositories

import (
	"fmt"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/drop"
	"go-api/internal/storage/postgres/drop_notification"
	"go-api/internal/storage/postgres/token"
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/model"
	"os"
)

type Repositories struct {
	//wg              *sync.WaitGroup
	UserRepository             model.UserRepository
	TokenRepository            model.AuthTokenRepository
	DropRepository             model.DropRepository
	DropNotificationRepository model.DropNotificationRepository
}

func Setup() *Repositories {
	sqlDB, err := postgres.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &Repositories{
		//wg:              wg,
		UserRepository:             user.NewRepo(sqlDB),
		TokenRepository:            token.NewRepo(sqlDB),
		DropRepository:             drop.NewRepo(sqlDB),
		DropNotificationRepository: drop_notification.NewRepo(sqlDB),
	}
}

func (r *Repositories) Disconnect() {
	//r.wg.Done()
}
