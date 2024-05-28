package fixtures

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
	"go-api/pkg/model"
)

func PopulateUsers(userRepo model.UserRepository) error {
	for i := range 1000 {
		user := model.UserCreationParam{
			Firstname: faker.FirstName(),
			Lastname:  faker.LastName(),
			Email:     faker.Email(),
			Password:  faker.Password(),
			Username:  faker.FirstName(),
			Roles:     []string{"user"},
		}
		_, err := userRepo.Create(
			user,
		)
		if err != nil {
			fmt.Println(user)
			return err
		}
		i++
	}
	return nil
}
