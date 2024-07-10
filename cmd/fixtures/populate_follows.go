package fixtures

import (
	"go-api/internal/storage/postgres"
	"gorm.io/gorm"
	"math/rand/v2"
)

func PopulateFollows(db *gorm.DB) error {
	var pvusers []postgres.User
	db.Where("is_private = ?", true).Find(&pvusers)

	var pubusers []postgres.User
	db.Where("is_private = ?", false).Find(&pubusers)

	for i := range pvusers {
		iuint := uint64(uint(i))
		s := rand.NewPCG(iuint, iuint*47329)
		r := rand.New(s)
		db.Create(&postgres.Follow{
			FollowerID: pubusers[r.IntN(len(pubusers))].ID,
			FollowedID: pvusers[i].ID,
			Status:     new(postgres.FollowPendingStatus).ToInt(),
		})
	}

	for i := range pubusers {
		iuint := uint64(uint(i))
		for range 50 {
			s := rand.NewPCG(iuint, iuint*47329)
			r := rand.New(s)
			db.Create(&postgres.Follow{
				FollowerID: pubusers[i].ID,
				FollowedID: pubusers[r.IntN(len(pubusers))].ID,
				Status:     new(postgres.FollowAcceptedStatus).ToInt(),
			})
		}
	}

	return nil
}
