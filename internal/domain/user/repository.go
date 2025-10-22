package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	FindAll(ctx context.Context) ([]*User, error)

	FindByEmail(ctx context.Context, email string) (*User, error)

	FindByToken(ctx context.Context, token uuid.UUID) (*User, error)

	Save(ctx context.Context, user *User) error

	Update(ctx context.Context, user *User) error

	Delete(ctx context.Context, id uuid.UUID) error
}
