package usecase

import (
	"context"
	"errors"
	"trainer/internal/application/dto"
	"trainer/internal/domain/user"

	"github.com/google/uuid"
)

type GetUser struct {
	userRepository user.Repository
}

func NewGetUser(userRepository user.Repository) *GetUser {
	return &GetUser{
		userRepository: userRepository,
	}
}

func (u *GetUser) Execute(ctx context.Context, req dto.GetUserRequest) (*dto.UserResponse, error) {
	if err := validateGet(req); err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	userModel, err := u.userRepository.FindByID(ctx, userId)

	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, errors.New("Not existed user")
	}

	return dto.NewUserResponse(userModel), nil
}

func validateGet(v interface{}) error {
	// В реальном приложении используйте github.com/go-playground/validator
	// validator.New().Struct(v)
	return nil
}
