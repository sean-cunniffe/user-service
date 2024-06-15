package servers

import (
	"testing"
)

func TestNewUserServiceServer(t *testing.T) {
	t.Run("create userServiceServer", func(t *testing.T) {})
}

func TestCreateUser(t *testing.T) {
	t.Run("create user successfully", func(t *testing.T) {})
	t.Run("create user that already exists", func(t *testing.T) {})
	t.Run("internal error occurs when creating user", func(t *testing.T) {})
	t.Run("create user this caller is not allowed to create", func(t *testing.T) {})
}

func TestGetUser(t *testing.T) {
	t.Run("get user successfully", func(t *testing.T) {})
	t.Run("get user that does not exist", func(t *testing.T) {})
	t.Run("internal error occurs when getting user", func(t *testing.T) {})

	t.Run("get user this caller is not allowed to get", func(t *testing.T) {})
}

func TestUpdateUser(t *testing.T) {
	t.Run("update user successfully", func(t *testing.T) {})
	t.Run("update user that does not exist", func(t *testing.T) {})
	t.Run("internal error occurs when updating user", func(t *testing.T) {})

	t.Run("update user this caller is not allowed to update", func(t *testing.T) {})
}

func TestDeleteUser(t *testing.T) {
	t.Run("delete user successfully", func(t *testing.T) {})
	t.Run("delete user that does not exist", func(t *testing.T) {})
	t.Run("internal error occurs when deleting user", func(t *testing.T) {})
	t.Run("delete user this caller is not allowed to delete", func(t *testing.T) {})
}
