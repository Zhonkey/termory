package app

import (
	"log"
	"os"
	"strconv"
	"time"
	"trainer/internal/application"
	"trainer/internal/application/usecase"
	"trainer/internal/domain/user"
	"trainer/internal/infrastructure"
	"trainer/internal/infrastructure/database"
)

type Container struct {
	DB             *database.DB
	TokenManager   application.TokenManager
	AccessTokenUC  *usecase.AccessToken
	RefreshTokenUC *usecase.RefreshToken
	CreateUserUC   *usecase.CreateUser
	UpdateUserUC   *usecase.UpdateUser
	DeleteUserUC   *usecase.DeleteUser
	GetUserUC      *usecase.GetUser
	ListUserUC     *usecase.ListUser
}

func NewContainer(db *database.DB) (*Container, error) {
	userRepo := database.NewUserRepository(db)

	passwordHasher := infrastructure.NewBcryptHasher(10)

	durationMinutes, err := strconv.Atoi(os.Getenv("JWT_DURATION_IN_MINUTE"))
	if err != nil {
		log.Fatalf("invalid JWT_DURATION_IN_MINUTE: %v", err)
	}
	tokenManager := infrastructure.NewJwtManager(os.Getenv("JWT_SECRET"), time.Duration(durationMinutes)*time.Minute)

	userService := user.NewService(userRepo, passwordHasher, time.Hour*24*30)

	accessTokenUC := usecase.NewAccessToken(userService, userRepo, tokenManager)
	refreshTokenUC := usecase.NewRefreshToken(userService, userRepo, tokenManager)
	createUserUC := usecase.NewCreateUser(userService, userRepo)
	updateUserUC := usecase.NewUpdateUser(userService, userRepo)
	deleteUserUC := usecase.NewDeleteUser(userRepo)
	getUserUC := usecase.NewGetUser(userRepo)
	listUserUC := usecase.NewListUser(userRepo)

	c := Container{
		db,
		tokenManager,
		accessTokenUC,
		refreshTokenUC,
		createUserUC,
		updateUserUC,
		deleteUserUC,
		getUserUC,
		listUserUC,
	}

	return &c, nil
}
