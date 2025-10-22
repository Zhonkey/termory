package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type Service struct {
	repo            Repository
	hasher          PasswordHasher
	refreshDuration time.Duration
}

func NewService(repo Repository, hasher PasswordHasher, refreshDuration time.Duration) *Service {
	return &Service{
		repo:            repo,
		hasher:          hasher,
		refreshDuration: refreshDuration,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrNotFound
	}

	if user == nil {
		return nil, ErrNotFound
	}

	if checkPassword := s.checkPassword(user, password); !checkPassword {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *Service) NewUser(ctx context.Context, email, firstName, lastName, password, role string) (*User, error) {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	return newUser(email, firstName, lastName, hashedPassword, RoleMentor)
}

func (s *Service) UpdateUser(user *User, firstName, lastName, email, password string) error {
	err := user.updateProfile(firstName, lastName, email)
	if err != nil {
		return ErrFailedUpdate
	}

	if password != "" {
		hashedPassword, err := s.hasher.Hash(password)
		if err != nil {
			return ErrInvalidPassword
		}
		err = user.updatePassword(hashedPassword)
		if err != nil {
			return ErrInvalidPassword
		}
	}

	return nil
}

func (s *Service) CreateRefreshToken(ctx context.Context, u *User) (*RefreshToken, error) {
	newToken, err := newRefreshToken(s.refreshDuration)
	if err != nil {
		return nil, ErrTokenRefresh
	}

	err = u.addRefreshToken(newToken)
	if err != nil {
		return nil, ErrTokenRefresh
	}

	return newToken, nil
}

func (s *Service) RenewRefreshToken(ctx context.Context, u *User, token uuid.UUID) (*RefreshToken, error) {
	err := u.revokeRefreshToken(token)
	if err != nil {
		return nil, ErrTokenRefresh
	}

	newToken, err := s.CreateRefreshToken(ctx, u)

	if err != nil {
		return nil, ErrTokenRefresh
	}

	return newToken, nil
}

func (s *Service) checkPassword(user *User, password string) bool {
	return s.hasher.Compare(user.Password, password)
}
