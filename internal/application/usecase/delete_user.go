package usecase

import (
	"context"
	"trainer/internal/application/dto"
	"trainer/internal/domain/user"

	"github.com/google/uuid"
)

type DeleteUser struct {
	userRepository user.Repository
}

func NewDeleteUser(userRepository user.Repository) *DeleteUser {
	return &DeleteUser{
		userRepository: userRepository,
	}
}

func (u *DeleteUser) Execute(ctx context.Context, req dto.DeleteUserRequest) error {
	if errValidate := validateDelete(req); errValidate != nil {
		return errValidate
	}

	userId, err := uuid.Parse(req.Id)
	if err != nil {
		return err
	}

	err = u.userRepository.Delete(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func validateDelete(v interface{}) error {
	// В реальном приложении используйте github.com/go-playground/validator
	// validator.New().Struct(v)
	return nil
}
