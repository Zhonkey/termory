package dto

type AccessTokenRequest struct {
	Email    string
	Password string
}
type RefreshTokenRequest struct {
	RefreshToken string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenResponse(accessToken string, refreshToken string) *TokenResponse {
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
