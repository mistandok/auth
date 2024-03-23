package tests

import (
	"context"
	"errors"
	"testing"

	userImpl "github.com/mistandok/auth/internal/api/user"
	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/mistandok/auth/pkg/user_v1"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCreate_SuccessDeleteUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	var userID int64 = 1
	request := &user_v1.DeleteRequest{Id: userID}

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Delete", ctx, userID).Return(nil).Once()

	service := userService.NewService(&logger, userRepoMock)

	impl := userImpl.NewImplementation(service)

	_, err := impl.Delete(ctx, request)

	require.NoError(t, err)
}

func TestCreate_FailDeleteUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	var userID int64 = 1
	request := &user_v1.DeleteRequest{Id: userID}
	someErr := errors.New("some error")

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Delete", ctx, userID).Return(someErr).Once()

	service := userService.NewService(&logger, userRepoMock)

	impl := userImpl.NewImplementation(service)

	_, err := impl.Delete(ctx, request)

	require.Error(t, err)
}
