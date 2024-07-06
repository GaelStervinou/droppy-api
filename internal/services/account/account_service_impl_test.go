package account

import (
	"errors"
	"go-api/internal/repositories"
	errors2 "go-api/pkg/errors2"
	"go-api/pkg/model"
	"testing"
)

type MockUserRepository struct{}

func (m *MockUserRepository) Update(args model.UserPatchParam) (model.UserModel, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) CanUserBeFollowed(followedId uint) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) GetUsersFromUserIds(userIds []uint) ([]model.UserModel, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) Search(query string) ([]model.UserModel, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) IsActiveUser(userId uint) (bool, error) {
	//TODO implement me
	panic("implement me")
}

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

func TestAccountService_Create(t *testing.T) {
	repo := repositories.Repositories{
		UserRepository:  &MockUserRepository{},
		TokenRepository: nil,
	}

	type fields struct {
		Repo *repositories.Repositories
	}

	type args struct {
		email    string
		password string
		username string
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
				email:    "test",
				password: "test123",
				username: "gae",
			},
			wantErr:     true,
			errorFields: []string{"email", "username", "password"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccountService{
				Repo: tt.fields.Repo,
			}
			err := a.Create(tt.args.email, tt.args.password, tt.args.username)
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
