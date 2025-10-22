package usecase

import (
	"context"
	"errors"
	"trainer/internal/application"
	"trainer/internal/application/dto"
	"trainer/internal/domain/user"

	"github.com/google/uuid"
)

type UpdateUser struct {
	userService    *user.Service
	userRepository user.Repository
}

func NewUpdateUser(userService *user.Service, userRepository user.Repository) *UpdateUser {
	return &UpdateUser{
		userService:    userService,
		userRepository: userRepository,
	}
}

func (u *UpdateUser) Execute(ctx context.Context, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if errValidate := application.ValidateDTO(req); errValidate != nil {
		return nil, errValidate
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

	err = u.userService.UpdateUser(userModel, req.FirstName, req.LastName, req.Email, req.Password)

	if err != nil {
		return nil, err
	}

	err = u.userRepository.Save(ctx, userModel)

	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(userModel), nil
}
