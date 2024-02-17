package user

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit"
	"github.com/mistandok/auth/pkg/user_v1"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server user Server.
type Server struct {
	user_v1.UnimplementedUserV1Server
	logger *zerolog.Logger
}

// NewServer generate instance for user Server.
func NewServer(logger *zerolog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

// Create user by param.
func (s *Server) Create(_ context.Context, request *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try create user: %+v", request))

	return &user_v1.CreateResponse{Id: 1}, nil
}

// Get user by params
func (s *Server) Get(_ context.Context, request *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try get user: %+v", request))

	return &user_v1.GetResponse{
		Id:        1,
		Name:      "Anton",
		Email:     "arti-anton@yandex.ru",
		Role:      0,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

// Update user by params.
func (s *Server) Update(_ context.Context, request *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try update user: %+v", request))

	return &emptypb.Empty{}, nil
}

// Delete user by params.
func (s *Server) Delete(_ context.Context, request *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Msg(fmt.Sprintf("try delete user: %+v", request))

	return &emptypb.Empty{}, nil
}