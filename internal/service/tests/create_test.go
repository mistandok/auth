package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func userForCreate() *model.UserForCreate {
	return &model.UserForCreate{
		Name:     "test",
		Email:    "test&email.com",
		Password: "password",
		Role:     "role",
	}
}

func TestCreate_SuccessCreateUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	user := userForCreate()
	var userID int64 = 1

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Create", ctx, user).Return(userID, nil).Once()

	service := userService.NewService(&logger, userRepoMock)

	resultUserID, err := service.Create(ctx, user)

	require.NoError(t, err)
	require.NotEmpty(t, resultUserID)
	require.Equal(t, userID, resultUserID)
}

func TestCreate_FailCreateUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	user := userForCreate()
	repoErr := errors.New("some error")

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Create", ctx, user).Return(int64(0), repoErr).Once()

	service := userService.NewService(&logger, userRepoMock)

	_, err := service.Create(ctx, user)

	require.Error(t, err)
}
