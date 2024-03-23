package tests

import (
	"time"

	"github.com/mistandok/auth/internal/common"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	userName     string = "test"
	userEmail           = "test"
	userPassword        = "test"
	userRoleInt         = 1
	userRoleStr         = "USER"
)

func userCreateRequest() *user_v1.CreateRequest {
	return &user_v1.CreateRequest{
		Name:            userName,
		Email:           userEmail,
		Password:        userPassword,
		PasswordConfirm: userPassword,
		Role:            userRoleInt,
	}
}

func userCreateForRepo() *model.UserForCreate {
	return &model.UserForCreate{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
		Role:     userRoleStr,
	}
}

func userResponseForGet(userID int64, userDate time.Time) *user_v1.GetResponse {
	return &user_v1.GetResponse{
		Id:        userID,
		Name:      userName,
		Email:     userEmail,
		Role:      userRoleInt,
		CreatedAt: timestamppb.New(userDate),
		UpdatedAt: timestamppb.New(userDate),
	}
}

func userFromRepo(userID int64, userDate time.Time) *model.User {
	return &model.User{
		ID:        userID,
		Name:      userName,
		Email:     userEmail,
		Role:      userRoleStr,
		CreatedAt: userDate,
		UpdatedAt: userDate,
	}
}

func userUpdateRequest(userID int64) *user_v1.UpdateRequest {
	return &user_v1.UpdateRequest{
		Id:    userID,
		Name:  common.Pointer[string](userName),
		Email: common.Pointer[string](userEmail),
		Role:  common.Pointer[user_v1.Role](userRoleInt),
	}
}

func userUpdateRepo(userID int64) *model.UserForUpdate {
	return &model.UserForUpdate{
		ID:    userID,
		Name:  common.Pointer[string](userName),
		Email: common.Pointer[model.UserEmail](userEmail),
		Role:  common.Pointer[model.UserRole](userRoleStr),
	}
}
