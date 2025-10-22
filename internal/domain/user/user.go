package user

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	Email         string
	Password      string
	Role          Role
	CreatedAt     time.Time
	UpdatedAt     time.Time
	refreshTokens map[uuid.UUID]*RefreshToken
	revokedTokens map[uuid.UUID]*RefreshToken
}

type RefreshToken struct {
	ID        uuid.UUID
	userID    uuid.UUID
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

type Role string

const (
	RoleStudent Role = "student"
	RoleMentor  Role = "mentor"
	RoleAdmin   Role = "admin"
)

func newUser(email, firstName, lastName, hashedPassword string, role Role) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if firstName == "" {
		return nil, errors.New("name must not be empty")
	}

	if lastName == "" {
		return nil, errors.New("name must not be empty")
	}

	if hashedPassword == "" {
		return nil, errors.New("password must not be empty")
	}

	if !isValidRole(role) {
		return nil, errors.New("invalid role")
	}

	now := time.Now()

	return &User{
		ID:            uuid.New(),
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Password:      hashedPassword,
		Role:          role,
		CreatedAt:     now,
		UpdatedAt:     now,
		refreshTokens: make(map[uuid.UUID]*RefreshToken),
		revokedTokens: make(map[uuid.UUID]*RefreshToken),
	}, nil
}

func (u *User) updateProfile(firstName, lastName, email string) error {
	if firstName != "" {
		u.FirstName = firstName
	}

	if lastName != "" {
		u.LastName = lastName
	}

	if email != "" && email != u.Email {
		if err := validateEmail(email); err != nil {
			return err
		}

		u.Email = email
	}

	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) updatePassword(newHashedPassword string) error {
	if newHashedPassword == "" {
		return errors.New("password cannot be empty")
	}

	u.Password = newHashedPassword
	u.UpdatedAt = time.Now()

	return nil
}

func (u *User) addRefreshToken(newToken *RefreshToken) error {
	u.refreshTokens[newToken.ID] = newToken
	return nil
}

func (u *User) revokeRefreshToken(tokenID uuid.UUID) error {
	token, exists := u.refreshTokens[tokenID]
	if !exists {
		return errors.New("refresh token not found")
	}

	u.revokedTokens[tokenID] = token
	delete(u.refreshTokens, tokenID)
	return nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsMentor() bool {
	return u.Role == RoleMentor
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) GetRefreshTokens() []*RefreshToken {
	tokens := make([]*RefreshToken, 0, len(u.refreshTokens))
	for _, token := range u.refreshTokens {
		tokens = append(tokens, token)
	}
	return tokens
}

func (u *User) GetRevokedTokens() []*RefreshToken {
	tokens := make([]*RefreshToken, 0, len(u.revokedTokens))
	for _, token := range u.revokedTokens {
		tokens = append(tokens, token)
	}
	return tokens
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func isValidRole(role Role) bool {
	switch role {
	case RoleStudent, RoleMentor, RoleAdmin:
		return true
	default:
		return false
	}
}

func NewUserFromStorage(id uuid.UUID, email, firstName, lastName, password string, role Role, createdAt, updatedAt time.Time, tokens []*RefreshToken) *User {
	user := &User{
		ID:            id,
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Password:      password,
		Role:          role,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		refreshTokens: make(map[uuid.UUID]*RefreshToken),
		revokedTokens: make(map[uuid.UUID]*RefreshToken),
	}

	for _, token := range tokens {
		if token.IsExpired() {
			user.revokedTokens[token.ID] = token
		} else {
			user.refreshTokens[token.ID] = token
		}
	}

	return user
}

func NewRefreshTokenFromStorage(id uuid.UUID, expiresAt, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		ID:        id,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}
}

func newRefreshToken(duration time.Duration) (*RefreshToken, error) {
	now := time.Now()

	return &RefreshToken{
		ID:        uuid.New(),
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}, nil
}
