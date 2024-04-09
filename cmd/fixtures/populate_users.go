package fixtures

import (
	"github.com/bxcodec/faker"
	"go-api/pkg/model"
)

func PopulateUsers(userRepo model.UserRepository) {
	for i := range 1000 {
		userRepo.Create(
			nil,
			model.UserCreationParam{
				Firstname: faker.FirstName,
				Lastname:  faker.LastName,
				Email:     faker.Email,
				Password:  faker.PASSWORD,
			},
		)
		i++
	}
}
