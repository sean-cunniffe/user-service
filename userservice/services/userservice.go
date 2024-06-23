package userservice

import (
	"errors"
	"user-service/dto"
)

// UserService is the intermediary func between grpc request and CRUD activities in the data store
type UserService interface {
	CreateUser(user *dto.User) (dto.User, error)
	GetUser(id string) (dto.User, error)
	UpdateUser(user *dto.User) (dto.User, error)
	DeleteUser(id string) error
	ListUsers() ([]*dto.User, error)
}

var err error = errors.New("unimplemented")

// NewUserService creates a userService
func NewUserService() UserService {
	return &userServiceImpl{}
}

type userServiceImpl struct {
}

func (*userServiceImpl) CreateUser(user *dto.User) (dto.User, error) {
	return dto.User{}, err
}

func (*userServiceImpl) GetUser(id string) (dto.User, error) {
	return dto.User{}, err
}

func (*userServiceImpl) UpdateUser(user *dto.User) (dto.User, error) {
	return dto.User{}, err
}
func (*userServiceImpl) DeleteUser(id string) error {
	return err
}
func (*userServiceImpl) ListUsers() ([]*dto.User, error) {
	return make([]*dto.User, 0), err
}
