package server_v1

import (
	"context"
	"fmt"

	"github.com/mistandok/auth/internal/common"
	"github.com/mistandok/auth/internal/repositories"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mistandok/auth/pkg/user_v1"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CRUDUserRepo interface for crud user repo in server.
type CRUDUserRepo interface {
	Create(context.Context, *repositories.UserCreateIn) (*repositories.UserCreateOut, error)
	Update(context.Context, *repositories.UserUpdateIn) error
	Get(context.Context, *repositories.UserGetIn) (*repositories.UserGetOut, error)
	Delete(context.Context, *repositories.UserDeleteIn) error
}

// Server user Server.
type Server struct {
	user_v1.UnimplementedUserV1Server
	logger   *zerolog.Logger
	userRepo CRUDUserRepo
}

// NewServer generate instance for user Server.
func NewServer(logger *zerolog.Logger, userRepo CRUDUserRepo) *Server {
	return &Server{
		logger:   logger,
		userRepo: userRepo,
	}
}

// Create user by param.
func (s *Server) Create(ctx context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try create user: %+v", request))

	out, err := s.userRepo.Create(ctx, &repositories.UserCreateIn{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     common.RoleNameFromRole(request.Role),
	})
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrEmailIsTaken):
			s.logger.Warn().Err(err).Msg("не удалось создать пользователя")
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			s.logger.Err(err).Msg("не удалось создать пользователя")
			return nil, status.Error(codes.Internal, "прошу понять и простить :(")
		}
	}

	return &user_v1.CreateResponse{Id: out.ID}, nil
}

// Get user by params
func (s *Server) Get(ctx context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try get user: %+v", request))

	out, err := s.userRepo.Get(ctx, &repositories.UserGetIn{ID: request.Id})
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrUserNotFound):
			s.logger.Warn().Msg("не удалось получить пользователя")
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			s.logger.Err(err).Msg("не удалось получить пользователя")
			return nil, status.Error(codes.Internal, "прошу понять и простить :(")
		}
	}

	return &user_v1.GetResponse{
		Id:        out.ID,
		Name:      out.Name,
		Email:     out.Email,
		Role:      common.RoleFromRoleName(out.Role),
		CreatedAt: timestamppb.New(out.CreatedAt),
		UpdatedAt: timestamppb.New(out.UpdatedAt),
	}, nil
}

// Update user by params.
func (s *Server) Update(ctx context.Context, request *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try update user: %+v", request))

	err := s.userRepo.Update(ctx, &repositories.UserUpdateIn{
		ID:    request.Id,
		Name:  request.Name,
		Email: request.Email,
		Role:  common.PointerRoleNameFromRole(request.Role),
	})
	if err != nil {
		s.logger.Err(err).Msg("не удалось обновить пользователя")
		return nil, status.Error(codes.Internal, "прошу понять и простить :(")
	}

	return &emptypb.Empty{}, nil
}

// Delete user by params.
func (s *Server) Delete(ctx context.Context, request *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try delete user: %+v", request))

	err := s.userRepo.Delete(ctx, &repositories.UserDeleteIn{ID: request.Id})
	if err != nil {
		s.logger.Err(err).Msg("не удалось удалить пользователя")
		return nil, status.Error(codes.Internal, "прошу понять и простить :(")
	}

	return &emptypb.Empty{}, nil
}
