package usecase

import (
	"context"
	"errors"
	"trainer/internal/application/dto"
	"trainer/internal/domain/user"
)

type CreateUser struct {
	userService    *user.Service
	userRepository user.Repository
}

func NewCreateUser(userService *user.Service, userRepository user.Repository) *CreateUser {
	return &CreateUser{
		userService:    userService,
		userRepository: userRepository,
	}
}

func (u *CreateUser) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	ctx.Value("role")

	if err := validateCreate(req); err != nil {
		return nil, err
	}

	userWithSameEmail, err := u.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if userWithSameEmail != nil {
		return nil, errors.New("Not unique email")
	}

	createdUser, err := u.userService.NewUser(ctx, req.Email, req.FirstName, req.LastName, req.Password, req.Role)

	if err != nil {
		return nil, err
	}

	err = u.userRepository.Save(ctx, createdUser)

	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(createdUser), nil
}

func validateCreate(v interface{}) error {
	// В реальном приложении используйте github.com/go-playground/validator
	// validator.New().Struct(v)
	return nil
}
