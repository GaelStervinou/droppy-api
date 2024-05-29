package fixtures

import (
	"go-api/internal/storage/postgres/follow"
	"go-api/internal/storage/postgres/user"
	"gorm.io/gorm"
	"math/rand/v2"
)

func PopulateFollows(db *gorm.DB) error {
	var pvusers []user.User
	db.Where("is_private = ?", true).Find(&pvusers)

	var pubusers []user.User
	db.Where("is_private = ?", false).Find(&pubusers)

	for i := range pvusers {
		iuint := uint64(uint(i))
		s := rand.NewPCG(iuint, iuint*47329)
		r := rand.New(s)
		db.Create(&follow.Follow{
			FollowerID: pubusers[r.IntN(len(pubusers))].ID,
			FollowedID: pvusers[i].ID,
			Status:     new(follow.FollowPendingStatus).ToInt(),
		})
	}

	for i := range pubusers {
		iuint := uint64(uint(i))
		for range 50 {
			s := rand.NewPCG(iuint, iuint*47329)
			r := rand.New(s)
			db.Create(&follow.Follow{
				FollowerID: pubusers[i].ID,
				FollowedID: pubusers[r.IntN(len(pubusers))].ID,
				Status:     new(follow.FollowAcceptedStatus).ToInt(),
			})
		}
	}

	return nil
}
