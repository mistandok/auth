package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mistandok/auth/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userImpl "github.com/mistandok/auth/internal/api/user"
	"github.com/mistandok/auth/pkg/user_v1"

	"github.com/mistandok/auth/internal/repository/mocks"
	userService "github.com/mistandok/auth/internal/service/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCreate_SuccessGetUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	var userID int64 = 1
	curTime := time.Now()
	user := userFromRepo(userID, curTime)
	expectedResponse := userResponseForGet(userID, curTime)

	userRepoMock := mocks.NewUserRepository(t)
	userRepoMock.On("Get", ctx, userID).Return(user, nil).Once()

	service := userService.NewService(&logger, userRepoMock)

	impl := userImpl.NewImplementation(service)

	response, err := impl.Get(ctx, &user_v1.GetRequest{Id: userID})

	require.NoError(t, err)
	require.Equal(t, expectedResponse, response)
}

func TestCreate_FailGetUser(t *testing.T) {
	ctx := context.TODO()
	logger := zerolog.Nop()

	var userID int64 = 1
	request := &user_v1.GetRequest{Id: userID}

	errorRepoMockGenerator := func(err error) *mocks.UserRepository {
		userRepoMock := mocks.NewUserRepository(t)
		userRepoMock.
			On("Get", ctx, userID).
			Return(nil, err).
			Once()

		return userRepoMock
	}

	tests := []struct {
		name              string
		getRequest        *user_v1.GetRequest
		internalError     error
		expectedErrorCode codes.Code
	}{
		{
			name:              "fail get because user not found",
			getRequest:        request,
			internalError:     repository.ErrUserNotFound,
			expectedErrorCode: codes.NotFound,
		},
		{
			name:              "fail create",
			getRequest:        request,
			internalError:     errors.New("some error"),
			expectedErrorCode: codes.Internal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoMock := errorRepoMockGenerator(test.internalError)
			service := userService.NewService(&logger, repoMock)
			impl := userImpl.NewImplementation(service)

			_, err := impl.Get(ctx, test.getRequest)

			require.Error(t, err)
			if e, ok := status.FromError(err); ok {
				require.Equal(t, e.Code(), test.expectedErrorCode)
			}
		})
	}
}
