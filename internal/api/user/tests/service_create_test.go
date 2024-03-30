package tests

import (
	"context"
	"errors"
	"github.com/mistandok/auth/internal/config"
	"github.com/mistandok/auth/internal/utils/password"
	"testing"

	userImpl "github.com/mistandok/auth/internal/api/user"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCreate_SuccessCreateUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	passManager := password.NewManager(&config.PasswordConfig{PasswordSalt: "test"})

	var userID int64 = 1
	request := userCreateRequest()

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Create", ctx, userCreateForRepo()).Return(userID, nil).Once()

	service := userService.NewService(&logger, userRepoMock, passManager)

	impl := userImpl.NewImplementation(service)

	resultUserID, err := impl.Create(ctx, request)

	require.NoError(t, err)
	require.NotEmpty(t, resultUserID)
}

func TestCreate_FailCreateUser(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	passManager := password.NewManager(&config.PasswordConfig{PasswordSalt: "test"})

	request := userCreateRequest()
	repoUserForCreate := userCreateForRepo()

	errorRepoMockGenerator := func(err error) *mocks.UserRepository {
		userRepoMock := mocks.NewUserRepository(t)
		userRepoMock.
			On("Create", ctx, repoUserForCreate).
			Return(int64(0), err).
			Once()

		return userRepoMock
	}

	tests := []struct {
		name              string
		createRequest     *user_v1.CreateRequest
		internalError     error
		expectedErrorCode codes.Code
	}{
		{
			name:              "fail create because email is taken",
			createRequest:     request,
			internalError:     repository.ErrEmailIsTaken,
			expectedErrorCode: codes.AlreadyExists,
		},
		{
			name:              "fail create",
			createRequest:     request,
			internalError:     errors.New("some error"),
			expectedErrorCode: codes.Internal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoMock := errorRepoMockGenerator(test.internalError)
			service := userService.NewService(&logger, repoMock, passManager)
			impl := userImpl.NewImplementation(service)

			_, err := impl.Create(ctx, test.createRequest)

			require.Error(t, err)
			if e, ok := status.FromError(err); ok {
				require.Equal(t, e.Code(), test.expectedErrorCode)
			}
		})
	}
}
