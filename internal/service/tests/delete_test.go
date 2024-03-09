package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestDelete_SuccessDeleteUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	var userIDForDelete int64 = 1

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Delete", ctx, userIDForDelete).Return(nil).Once()

	service := userService.NewService(&logger, userRepoMock)

	err := service.Delete(ctx, userIDForDelete)

	require.NoError(t, err)
}

func TestDelete_FailDeleteUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	var userIDForDelete int64 = 1
	repoErr := errors.New("some error")

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Delete", ctx, userIDForDelete).Return(repoErr).Once()

	service := userService.NewService(&logger, userRepoMock)

	err := service.Delete(ctx, userIDForDelete)

	require.Error(t, err)
}
