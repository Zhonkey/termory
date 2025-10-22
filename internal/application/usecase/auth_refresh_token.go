package usecase

import (
	"context"
	"trainer/internal/application"
	"trainer/internal/application/dto"
	"trainer/internal/domain/user"

	"github.com/google/uuid"
)

type RefreshToken struct {
	userService    *user.Service
	userRepository user.Repository
	jwtManager     application.JwtManager
}

func NewRefreshToken(userService *user.Service, userRepository user.Repository, jwtManager application.JwtManager) *RefreshToken {
	return &RefreshToken{
		userService:    userService,
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

func (r *RefreshToken) Execute(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	//if errValidate := validateCreate(req); errValidate != nil {
	//	return nil, errValidate
	//}

	refreshToken, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	loggedUser, err := r.userRepository.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	newToken, err := r.userService.RenewRefreshToken(ctx, loggedUser, refreshToken)

	if newToken == nil || err != nil {
		return nil, err
	}

	errSave := r.userRepository.Save(ctx, loggedUser)
	if errSave != nil {
		return nil, errSave
	}

	accessToken, err := r.jwtManager.Generate(application.TokenClaim{
		UserID: loggedUser.ID,
		Role:   loggedUser.Role,
	})

	if err != nil {
		return nil, err
	}

	return dto.NewTokenResponse(accessToken, newToken.ID.String()), nil
}
