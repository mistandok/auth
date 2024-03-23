package tests

import (
	"context"
	"errors"
	"testing"

	userImpl "github.com/mistandok/auth/internal/api/user"
	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCreate_SuccessUpdateUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	var userID int64 = 1
	request := userUpdateRequest(userID)

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Update", ctx, userUpdateRepo(userID)).Return(nil).Once()

	service := userService.NewService(&logger, userRepoMock)

	impl := userImpl.NewImplementation(service)

	_, err := impl.Update(ctx, request)

	require.NoError(t, err)
}

func TestCreate_FailUpdateUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	var userID int64 = 1
	request := userUpdateRequest(userID)
	someErr := errors.New("some error")

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Update", ctx, userUpdateRepo(userID)).Return(someErr).Once()

	service := userService.NewService(&logger, userRepoMock)

	impl := userImpl.NewImplementation(service)

	_, err := impl.Update(ctx, request)

	require.Error(t, err)
}
