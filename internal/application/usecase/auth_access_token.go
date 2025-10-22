package usecase

import (
	"context"
	"trainer/internal/application"
	"trainer/internal/application/dto"
	"trainer/internal/domain/user"
)

type AccessToken struct {
	userService    *user.Service
	userRepository user.Repository
	jwtManager     application.TokenManager
}

func NewAccessToken(userService *user.Service, userRepository user.Repository, jwtManager application.TokenManager) *AccessToken {
	return &AccessToken{
		userService:    userService,
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

func (a *AccessToken) Execute(ctx context.Context, req dto.AccessTokenRequest) (*dto.TokenResponse, error) {
	//if errValidate := validateCreate(req); errValidate != nil {
	//	return nil, errValidate
	//}

	loggedUser, err := a.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.jwtManager.Generate(application.TokenClaim{
		UserID: loggedUser.ID,
		Role:   loggedUser.Role,
	})

	if err != nil {
		return nil, err
	}

	refreshToken, err := a.userService.CreateRefreshToken(ctx, loggedUser)
	if err != nil {
		return nil, err
	}

	errSave := a.userRepository.Update(ctx, loggedUser)
	if errSave != nil {
		return nil, errSave
	}

	return dto.NewTokenResponse(accessToken, refreshToken.ID.String(), loggedUser), nil
}
