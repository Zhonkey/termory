package dto

import "trainer/internal/domain/user"

type CreateUserRequest struct {
	Role      string `validate:"required" json:"role"`
	Email     string `validate:"required,email" json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `validate:"required"`
}

type UpdateUserRequest struct {
	Id        string `validate:"required" json:"id"`
	Email     string `validate:"email" json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type GetUserRequest struct {
	Id string `validate:"required" json:"id"`
}

type ListUserRequest struct {
}

type DeleteUserRequest struct {
	Id string `validate:"required" json:"id"`
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
	for i, userModel := range users {
		responseUsers[i] = NewUserResponse(userModel)
	}

	return &ListUserResponse{
		Users: responseUsers,
	}
}
