package user

import (
	"context"
	"errors"
	"github.com/mistandok/auth/internal/api"

	"github.com/mistandok/auth/internal/service/user"

	"github.com/mistandok/auth/internal/convert"
	"github.com/mistandok/auth/internal/repository"
	"github.com/mistandok/auth/internal/service"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Implementation user Server.
type Implementation struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation ..
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}

// Create ..
func (i *Implementation) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	userID, err := i.userService.Create(ctx, convert.ToServiceUserForCreateFromCreateRequest(request))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailIsTaken):
			return nil, status.Error(codes.AlreadyExists, repository.ErrEmailIsTaken.Error())
		case errors.Is(err, user.ErrPassToLong):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, api.ErrInternal
		}
	}

	return &user_v1.CreateResponse{Id: int64(userID)}, nil
}

// Get user by params
func (i *Implementation) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	user, err := i.userService.Get(ctx, request.Id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, repository.ErrUserNotFound.Error())
		default:
			return nil, api.ErrInternal
		}
	}

	return convert.ToGetResponseFromServiceUser(user), nil
}

// Update user by params.
func (i *Implementation) Update(ctx context.Context, request *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, convert.ToServiceUserForUpdateFromUpdateRequest(request))
	if err != nil {
		return nil, api.ErrInternal
	}

	return &emptypb.Empty{}, nil
}

// Delete user by params.
func (i *Implementation) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, request.Id)
	if err != nil {
		return nil, api.ErrInternal
	}

	return &emptypb.Empty{}, nil
}
