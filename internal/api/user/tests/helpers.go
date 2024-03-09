package tests

import (
	"time"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func userCreateRequest() *user_v1.CreateRequest {
	return &user_v1.CreateRequest{
		Name:            "test",
		Email:           "test",
		Password:        "test",
		PasswordConfirm: "test",
		Role:            1,
	}
}

func userCreateForRepo() *model.UserForCreate {
	return &model.UserForCreate{
		Name:     "test",
		Email:    "test",
		Password: "test",
		Role:     "USER",
	}
}

func userResponseForGet(userID int64, userDate time.Time) *user_v1.GetResponse {
	return &user_v1.GetResponse{
		Id:        userID,
		Name:      "test",
		Email:     "test",
		Role:      1,
		CreatedAt: timestamppb.New(userDate),
		UpdatedAt: timestamppb.New(userDate),
	}
}

func userFromRepo(userID int64, userDate time.Time) *model.User {
	return &model.User{
		ID:        userID,
		Name:      "test",
		Email:     "test",
		Role:      "USER",
		CreatedAt: userDate,
		UpdatedAt: userDate,
	}
}
