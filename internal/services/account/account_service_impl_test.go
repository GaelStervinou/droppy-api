package account

import (
	"errors"
	"fmt"
	"go-api/internal/repositories"
	errors2 "go-api/pkg/errors"
	"go-api/pkg/model"
	"testing"
)

type MockUserRepository struct{}

func (m *MockUserRepository) Create(args model.UserCreationParam) (model.UserModel, error) {
	return nil, nil
}
func (m *MockUserRepository) CreateWithGoogle(args model.UserCreationWithGoogleParam) (model.UserModel, error) {
	return nil, nil
}
func (m *MockUserRepository) Delete(id uint) error {
	return nil
}
func (m *MockUserRepository) GetAll() ([]model.UserModel, error) {
	return nil, nil
}
func (m *MockUserRepository) GetByEmail(email string) (model.UserModel, error) {
	return nil, nil
}
func (m *MockUserRepository) GetById(id uint) (model.UserModel, error) {
	return nil, nil
}
func (m *MockUserRepository) GetByGoogleAuthId(googleID string) (model.UserModel, error) {
	return nil, nil
}
func (m *MockUserRepository) Update(user model.UserModel) (model.UserModel, error) {
	return nil, nil
}

func TestAccountService_Create(t *testing.T) {
	repo := repositories.Repositories{
		UserRepository:  &MockUserRepository{},
		TokenRepository: nil,
	}
	a := AccountService{
		Repo: &repo,
	}

	fmt.Println(a)
	type fields struct {
		Repo *repositories.Repositories
	}

	type args struct {
		email     string
		password  string
		firstname string
		lastname  string
		username  string
	}
	tests := []struct {
		name        string
		wantErr     bool
		errorFields []string
		fields      fields
		args        args
	}{
		{
			name: "Test all fields are invalid",
			fields: fields{
				Repo: &repo,
			},
			args: args{
				email:     "test",
				password:  "test123",
				firstname: "G",
				lastname:  "B",
				username:  "gae",
			},
			wantErr:     true,
			errorFields: []string{"email", "firstname", "lastname", "username", "password"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountService{
				Repo: tt.fields.Repo,
			}
			err := a.Create(tt.args.firstname, tt.args.lastname, tt.args.email, tt.args.password, tt.args.username)
			if err == nil && tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errorFields != nil {
				if errors.Is(err, errors2.MultiFieldsError{}) {
					var multiErr errors2.MultiFieldsError
					errors.As(err, &multiErr)
					for field, _ := range multiErr.Fields {
						ok := false
						for index, wanted := range tt.errorFields {
							if wanted == field {
								ok = true
								tt.errorFields = append(tt.errorFields[:index], tt.errorFields[index+1:]...)
								break
							}
						}
						if !ok {
							t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
						}
					}
					if len(tt.errorFields) > 0 {
						t.Errorf("Missing fields: %v", tt.errorFields)
					}
				} else {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
