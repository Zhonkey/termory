package dto

import "trainer/internal/domain/user"

type AccessTokenRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}
type RefreshTokenRequest struct {
	RefreshToken string `validate:"required" json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

func NewTokenResponse(accessToken string, refreshToken string, user *user.User) *TokenResponse {
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         NewUserResponse(user),
	}
}
