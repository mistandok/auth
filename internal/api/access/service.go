package access

import (
	"context"
	"errors"

	"github.com/mistandok/auth/internal/api"

	"github.com/mistandok/auth/internal/convert"
	"github.com/mistandok/auth/internal/service"
	"github.com/mistandok/auth/pkg/access_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Implementation user Server.
type Implementation struct {
	access_v1.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewImplementation ..
func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{accessService: accessService}
}

// Create ..
func (i *Implementation) Create(ctx context.Context, request *access_v1.CreateRequest) (*access_v1.CreateResponse, error) {
	id, err := i.accessService.Create(ctx, convert.ToServiceEndpointAccessFromCreateRequest(request))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNeedAdminRole):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		default:
			return nil, api.ErrInternal
		}
	}

	return &access_v1.CreateResponse{Id: id}, nil
}

// Check ..
func (i *Implementation) Check(ctx context.Context, request *access_v1.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, request.Address)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.PermissionDenied, err.Error())
	}

	return &emptypb.Empty{}, nil
}
