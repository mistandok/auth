package convert

import (
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/pkg/access_v1"
	"github.com/mistandok/auth/pkg/user_v1"
)

// ToServiceEndpointAccessFromCreateRequest ..
func ToServiceEndpointAccessFromCreateRequest(request *access_v1.CreateRequest) *serviceModel.EndpointAccess {
	return &serviceModel.EndpointAccess{
		Address: request.Address,
		Role:    ToServiceRoleFromAccessRole(request.Role),
	}
}

// ToServiceRoleFromAccessRole ..
func ToServiceRoleFromAccessRole(role user_v1.Role) serviceModel.UserRole {
	roleName := user_v1.Role_name[int32(role)]

	return serviceModel.UserRole(roleName)
}
