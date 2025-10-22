package usecase

import (
	"context"
	"trainer/internal/application/dto"
	"trainer/internal/domain/user"
)

type ListUser struct {
	userRepository user.Repository
}

func NewListUser(userRepository user.Repository) *ListUser {
	return &ListUser{
		userRepository: userRepository,
	}
}

func (u *ListUser) Execute(ctx context.Context, req dto.ListUserRequest) (*dto.ListUserResponse, error) {
	if err := validateList(req); err != nil {
		return nil, err
	}

	users, err := u.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return dto.NewUsersResponse(users), nil
}

func validateList(v interface{}) error {
	// В реальном приложении используйте github.com/go-playground/validator
	// validator.New().Struct(v)
	return nil
}
