package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/mistandok/auth/internal/common"
	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func userForUpdate() *model.UserForUpdate {
	return &model.UserForUpdate{
		ID:    1,
		Name:  common.Pointer[string]("test"),
		Email: common.Pointer[model.UserEmail]("test"),
		Role:  common.Pointer[model.UserRole]("admin"),
	}
}

func TestUpdate_SuccessUpdateUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	user := userForUpdate()

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Update", ctx, user).Return(nil).Once()

	service := userService.NewService(&logger, userRepoMock)

	err := service.Update(ctx, user)

	require.NoError(t, err)
}

func TestUpdate_FailUpdateUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	user := userForUpdate()
	repoErr := errors.New("some error")

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Update", ctx, user).Return(repoErr).Once()

	service := userService.NewService(&logger, userRepoMock)

	err := service.Update(ctx, user)

	require.Error(t, err)
}
