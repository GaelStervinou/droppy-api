// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/model/user_iface.go
//
// Generated by this command:
//
//	mockgen -source=pkg/model/user_iface.go
//

// Package mock_model is a generated GoMock package.
package mock_model

import (
	model "go-api/pkg/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserModel is a mock of UserModel interface.
type MockUserModel struct {
	ctrl     *gomock.Controller
	recorder *MockUserModelMockRecorder
}

// MockUserModelMockRecorder is the mock recorder for MockUserModel.
type MockUserModelMockRecorder struct {
	mock *MockUserModel
}

// NewMockUserModel creates a new mock instance.
func NewMockUserModel(ctrl *gomock.Controller) *MockUserModel {
	mock := &MockUserModel{ctrl: ctrl}
	mock.recorder = &MockUserModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserModel) EXPECT() *MockUserModelMockRecorder {
	return m.recorder
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(args model.UserCreationParam) (model.UserModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", args)
	ret0, _ := ret[0].(model.UserModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(args any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), args)
}

// CreateWithGoogle mocks base method.
func (m *MockUserRepository) CreateWithGoogle(args model.UserCreationWithGoogleParam) (model.UserModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithGoogle", args)
	ret0, _ := ret[0].(model.UserModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWithGoogle indicates an expected call of CreateWithGoogle.
func (mr *MockUserRepositoryMockRecorder) CreateWithGoogle(args any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithGoogle", reflect.TypeOf((*MockUserRepository)(nil).CreateWithGoogle), args)
}

// Delete mocks base method.
func (m *MockUserRepository) Delete(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserRepositoryMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserRepository)(nil).Delete), id)
}

// GetAll mocks base method.
func (m *MockUserRepository) GetAll() ([]model.UserModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.UserModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockUserRepositoryMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUserRepository)(nil).GetAll))
}

// GetByEmail mocks base method.
func (m *MockUserRepository) GetByEmail(email string) (model.UserModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(model.UserModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUserRepositoryMockRecorder) GetByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetByEmail), email)
}

// GetByGoogleAuthId mocks base method.
func (m *MockUserRepository) GetByGoogleAuthId(googleID string) (model.UserModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFirebaseUid", googleID)
	ret0, _ := ret[0].(model.UserModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByGoogleAuthId indicates an expected call of GetByGoogleAuthId.
func (mr *MockUserRepositoryMockRecorder) GetByGoogleAuthId(googleID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFirebaseUid", reflect.TypeOf((*MockUserRepository)(nil).GetByGoogleAuthId), googleID)
}

// GetById mocks base method.
func (m *MockUserRepository) GetById(id uint) (model.UserModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(model.UserModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUserRepositoryMockRecorder) GetById(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUserRepository)(nil).GetById), id)
}

// Update mocks base method.
func (m *MockUserRepository) Update(user model.UserModel) (model.UserModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", user)
	ret0, _ := ret[0].(model.UserModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), user)
}
