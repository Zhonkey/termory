package dto

import "trainer/internal/domain/user"

type CreateUserRequest struct {
	Role      string
	Email     string
	FirstName string
	LastName  string
	Password  string
}

type UpdateUserRequest struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
	Password  string
}

type GetUserRequest struct {
	Id string
}

type DeleteUserRequest struct {
	Id string
}

type UserResponse struct {
}

func NewUserResponse(user *user.User) *UserResponse {
	return &UserResponse{}
}
