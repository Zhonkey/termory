package application

import (
	"trainer/internal/domain/user"

	"github.com/google/uuid"
)

type JwtManager interface {
	Generate(claim TokenClaim) (string, error)
	Parse(tokenStr string) (*TokenClaim, error)
}

type TokenClaim struct {
	UserID uuid.UUID
	Role   user.Role
}
