package fixtures

import (
	"github.com/go-faker/faker/v4"
	"go-api/internal/storage/postgres"
	"gorm.io/gorm"
	"math/rand/v2"
)

func PopulateUsers(db *gorm.DB) error {
	statuses := []int{-2, -1, 0, 1}
	avatars := []string{"profile/random1", "profile/random2", "profile/random3"}
	roles := []string{"user", "admin"}

	for i := range 1000 {
		iuint := uint64(uint(i))
		s := rand.NewPCG(iuint, iuint*47329)
		r := rand.New(s)
		var avatar string
		private := i%2 == 0
		// initialize local pseudorandom generator
		status := statuses[r.IntN(len(statuses))]
		if i%10 != 0 {
			avatar = avatars[r.IntN(len(avatars))]
		}
		role := roles[r.IntN(len(roles))]
		db.Create(&postgres.User{
			Email:       faker.Email(),
			Password:    faker.Password(),
			Username:    faker.FirstName(),
			Role:        role,
			Status:      status,
			IsPrivate:   private,
			Bio:         faker.Sentence(),
			PhoneNumber: faker.Phonenumber(),
			Avatar:      avatar,
		})
		i++
	}
	return nil
}
