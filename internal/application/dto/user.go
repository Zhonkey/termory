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

type ListUserRequest struct {
}

type DeleteUserRequest struct {
	Id string
}

type UserResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type ListUserResponse struct {
	Users []*UserResponse `json:"users"`
}

func NewUserResponse(user *user.User) *UserResponse {
	return &UserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func NewUsersResponse(users []*user.User) *ListUserResponse {
	responseUsers := make([]*UserResponse, len(users))
	for i, user := range users {
		responseUsers[i] = NewUserResponse(user)
	}

	return &ListUserResponse{
		Users: responseUsers,
	}
}
