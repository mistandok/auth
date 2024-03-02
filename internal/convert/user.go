package convert

import (
	"github.com/mistandok/auth/internal/common"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToServiceUserForCreateFromCreateRequest ..
func ToServiceUserForCreateFromCreateRequest(request *user_v1.CreateRequest) *serviceModel.UserForCreate {
	return &serviceModel.UserForCreate{
		Name:     serviceModel.UserName(request.Name),
		Email:    serviceModel.UserEmail(request.Email),
		Password: serviceModel.Password(request.Password),
		Role:     ToServiceRoleFromRole(request.Role),
	}
}

// ToServiceUserForUpdateFromUpdateRequest ..
func ToServiceUserForUpdateFromUpdateRequest(request *user_v1.UpdateRequest) *serviceModel.UserForUpdate {
	var (
		name  *serviceModel.UserName
		email *serviceModel.UserEmail
		role  *serviceModel.UserRole
	)

	if request.Name == nil {
		name = nil
	} else {
		name = common.Pointer[serviceModel.UserName](serviceModel.UserName(*request.Name))
	}

	if request.Email == nil {
		email = nil
	} else {
		email = common.Pointer[serviceModel.UserEmail](serviceModel.UserEmail(*request.Email))
	}

	if request.Role == nil {
		role = nil
	} else {
		role = common.Pointer[serviceModel.UserRole](ToServiceRoleFromRole(*request.Role))
	}

	return &serviceModel.UserForUpdate{
		ID:    serviceModel.UserID(request.Id),
		Name:  name,
		Email: email,
		Role:  role,
	}
}

// ToServiceRoleFromRole ..
func ToServiceRoleFromRole(role user_v1.Role) serviceModel.UserRole {
	roleName := user_v1.Role_name[int32(role)]

	return serviceModel.UserRole(roleName)
}

func ToRoleFromServiceRole(role serviceModel.UserRole) user_v1.Role {
	resultRole := user_v1.Role_value[string(role)]

	return user_v1.Role(resultRole)
}

func ToGetResponseFromServiceUser(user *serviceModel.User) *user_v1.GetResponse {
	return &user_v1.GetResponse{
		Id:        int64(user.ID),
		Name:      string(user.Name),
		Email:     string(user.Email),
		Role:      ToRoleFromServiceRole(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
