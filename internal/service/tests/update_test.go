package tests

import (
	"context"
	"errors"
	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/utils/password"
	"testing"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/mistandok/auth/internal/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func userForUpdate() *model.UserForUpdate {
	return &model.UserForUpdate{
		ID:    1,
		Name:  utils.Pointer[string]("test"),
		Email: utils.Pointer[model.UserEmail]("test"),
		Role:  utils.Pointer[model.UserRole]("admin"),
	}
}

func TestUpdate_SuccessUpdateUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	passManager := password.NewManager(&config.PasswordConfig{PasswordSalt: "test"})

	user := userForUpdate()

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Update", ctx, user).Return(nil).Once()

	service := userService.NewService(&logger, userRepoMock, passManager)

	err := service.Update(ctx, user)

	require.NoError(t, err)
}

func TestUpdate_FailUpdateUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	passManager := password.NewManager(&config.PasswordConfig{PasswordSalt: "test"})

	user := userForUpdate()
	repoErr := errors.New("some error")

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Update", ctx, user).Return(repoErr).Once()

	service := userService.NewService(&logger, userRepoMock, passManager)

	err := service.Update(ctx, user)

	require.Error(t, err)
}
