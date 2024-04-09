package repositories

import (
	"fmt"
	"go-api/internal/storage/postgres"
	"go-api/internal/storage/postgres/token"
	"go-api/internal/storage/postgres/user"
	"go-api/pkg/model"
	"os"
	"sync"
)

type Repositories struct {
	wg              *sync.WaitGroup
	UserRepository  model.UserRepository
	TokenRepository model.AuthTokenRepository
}

func Setup(wg *sync.WaitGroup) *Repositories {
	sqlDB, err := postgres.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &Repositories{
		wg:              wg,
		UserRepository:  user.NewRepo(sqlDB),
		TokenRepository: token.NewRepo(sqlDB),
	}
}

func (r *Repositories) Disconnect() {
	r.wg.Done()
}
