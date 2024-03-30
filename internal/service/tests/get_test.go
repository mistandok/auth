package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/utils/password"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func userForGet(userID int64) *model.User {
	return &model.User{
		ID:        userID,
		Name:      "test",
		Email:     "test",
		Role:      "admin",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func TestGet_SuccessGetUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	passManager := password.NewManager(&config.PasswordConfig{PasswordSalt: "test"})

	var userIDForGet int64 = 1
	expectedUser := userForGet(userIDForGet)

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Get", ctx, userIDForGet).Return(expectedUser, nil).Once()

	service := userService.NewService(&logger, userRepoMock, passManager)

	resultUser, err := service.Get(ctx, userIDForGet)

	require.NoError(t, err)
	require.Equal(t, expectedUser, resultUser)
}

func TestGet_FailGetUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	passManager := password.NewManager(&config.PasswordConfig{PasswordSalt: "test"})

	var userIDForGet int64 = 1
	repoErr := errors.New("some error")

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Get", ctx, userIDForGet).Return(nil, repoErr).Once()

	service := userService.NewService(&logger, userRepoMock, passManager)

	_, err := service.Get(ctx, userIDForGet)

	require.Error(t, err)
}
