package fixtures

import (
	"github.com/go-faker/faker/v4"
	"go-api/internal/storage/postgres/drop"
	"go-api/internal/storage/postgres/drop_notification"
	"go-api/internal/storage/postgres/user"
	"gorm.io/gorm"
	"math/rand/v2"
)

func PopulateDrops(db *gorm.DB) error {

	types := []string{"youtube", "twitch", "films", "spotify"}
	var dropNotifications []drop_notification.DropNotification
	for i := range 365 {
		iuint := uint64(uint(i))
		s := rand.NewPCG(iuint, iuint*47329)
		r := rand.New(s)
		dropNotifications = append(dropNotifications, drop_notification.DropNotification{
			Type: types[r.IntN(len(types))],
		})
		db.Create(&dropNotifications[i])
	}

	var activeUsers []user.User
	db.Where("status = ?", 1).Find(&activeUsers)
	statuses := []int{-1, 1}
	for i := range dropNotifications {
		iuint := uint64(uint(i))
		for j := range activeUsers {
			s := rand.NewPCG(iuint, iuint*47329*uint64(j))
			r := rand.New(s)
			isPinned := r.IntN(10) == 0
			status := statuses[r.IntN(len(statuses))]
			db.Create(&drop.Drop{
				Type:               dropNotifications[i].Type,
				Content:            "content",
				Description:        faker.Sentence(),
				CreatedById:        activeUsers[j].ID,
				Status:             uint(status),
				IsPinned:           isPinned,
				DropNotificationID: dropNotifications[i].ID,
			})
		}
	}

	return nil
}
