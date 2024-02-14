package user

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	user_v1.UnimplementedUserV1Server
}

func (s *Server) Get(context.Context, *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	return &user_v1.GetResponse{
		Id:        1,
		Name:      "Anton",
		Email:     "arti-anton@yandex.ru",
		Role:      0,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}
