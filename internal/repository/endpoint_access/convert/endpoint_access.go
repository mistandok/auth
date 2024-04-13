package convert

import (
	serviceModel "github.com/mistandok/auth/internal/model"
	repoModel "github.com/mistandok/auth/internal/repository/endpoint_access/model"
)

// ToServiceEndpointAccessFromRepoEndpointAccess ..
func ToServiceEndpointAccessFromRepoEndpointAccess(endpointAccess *repoModel.EndpointAccess) *serviceModel.EndpointAccess {
	return &serviceModel.EndpointAccess{
		ID:        endpointAccess.ID,
		Address:   endpointAccess.Address,
		Role:      serviceModel.UserRole(endpointAccess.Role),
		CreatedAt: endpointAccess.CreatedAt,
		UpdatedAt: endpointAccess.UpdatedAt,
	}
}
