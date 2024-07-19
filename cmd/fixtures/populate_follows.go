package fixtures

import (
	"go-api/internal/storage/postgres"
	"gorm.io/gorm"
	"math/rand"
)

func PopulateFollows(db *gorm.DB) error {
	var users []postgres.User
	db.Model(&postgres.User{}).Where("status = 1").Find(&users)

	for _, user := range users {
		var status uint = 1
		if user.IsPrivate {
			status = 0
		}
		fnb := rand.Intn(50)
		for range fnb {
			db.Create(&postgres.Follow{
				FollowerID: users[rand.Intn(len(users))].ID,
				FollowedID: user.ID,
				Status:     status,
			})
		}
	}

	return nil
}
