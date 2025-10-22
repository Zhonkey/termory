package user

import "errors"

var (
	ErrUserNotFound        = errors.New("USER_NOT_FOUND")
	ErrEmailAlreadyUsed    = errors.New("EMAIL_ALREADY_USED")
	ErrInvalidPassword     = errors.New("INVALID_PASSWORD")
	ErrAccessDenied        = errors.New("ACCESS_DENIED")
	ErrInvalidRole         = errors.New("INVALID_ROLE")
	ErrInvalidEmail        = errors.New("INVALID_EMAIL")
	ErrEmptyPassword       = errors.New("EMPTY_PASSWORD")
	ErrEmptyEmail          = errors.New("EMPTY_EMAIL")
	ErrInvalidRefreshToken = errors.New("INVALID_REFRESH_TOKEN")
	ErrNotFound            = errors.New("USER_NOT_FOUND")
	ErrFailedUpdate        = errors.New("UPDATE_FAILED")
	ErrTokenRefresh        = errors.New("REFRESH_TOKEN_FAILED")
)
