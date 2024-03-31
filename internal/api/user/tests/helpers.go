package tests

import (
	"time"

	"github.com/mistandok/auth/internal/utils"

	"github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	userName      string          = "test"
	userEmail     model.UserEmail = "test"
	userEmailStr  string          = "test"
	userPassword  string          = "password"
	userRoleStr   string          = "USER"
	userRole      user_v1.Role    = 1
	userRoleModel model.UserRole  = "USER"
)

func userCreateRequest() *user_v1.CreateRequest {
	return &user_v1.CreateRequest{
		Name:            userName,
		Email:           userEmailStr,
		Password:        userPassword,
		PasswordConfirm: userPassword,
		Role:            userRole,
	}
}

func userResponseForGet(userID int64, userDate time.Time) *user_v1.GetResponse {
	return &user_v1.GetResponse{
		Id:        userID,
		Name:      userName,
		Email:     userEmailStr,
		Role:      userRole,
		CreatedAt: timestamppb.New(userDate),
		UpdatedAt: timestamppb.New(userDate),
	}
}

func userFromRepo(userID int64, userDate time.Time) *model.User {
	return &model.User{
		ID:        userID,
		Name:      userName,
		Email:     userEmail,
		Role:      userRoleModel,
		CreatedAt: userDate,
		UpdatedAt: userDate,
	}
}

func userUpdateRequest(userID int64) *user_v1.UpdateRequest {
	return &user_v1.UpdateRequest{
		Id:    userID,
		Name:  utils.Pointer[string](userName),
		Email: utils.Pointer[string](userEmailStr),
		Role:  utils.Pointer[user_v1.Role](userRole),
	}
}

func userUpdateRepo(userID int64) *model.UserForUpdate {
	return &model.UserForUpdate{
		ID:    userID,
		Name:  utils.Pointer[string](userName),
		Email: utils.Pointer[model.UserEmail](userEmail),
		Role:  utils.Pointer[model.UserRole](userRoleModel),
	}
}
