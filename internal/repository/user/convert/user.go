package convert

import (
	"github.com/mistandok/auth/internal/common"
	serviceModel "github.com/mistandok/auth/internal/model"
	repoModel "github.com/mistandok/auth/internal/repository/user/model"
)

// ToRepoUserForUpdateFromServiceUserForUpdate ..
func ToRepoUserForUpdateFromServiceUserForUpdate(user *serviceModel.UserForUpdate) *repoModel.UserForUpdate {
	var (
		name  *string
		email *string
		role  *string
	)

	if user.Name != nil {
		name = common.Pointer[string](*user.Name)
	}

	if user.Email != nil {
		email = common.Pointer[string](string(*user.Email))
	}

	if user.Role != nil {
		role = common.Pointer[string](string(*user.Role))
	}

	return &repoModel.UserForUpdate{
		ID:    user.ID,
		Name:  name,
		Email: email,
		Role:  role,
	}
}

// ToServiceUserFromRepoUser ..
func ToServiceUserFromRepoUser(user *repoModel.User) *serviceModel.User {
	return &serviceModel.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     serviceModel.UserEmail(user.Email),
		Role:      serviceModel.UserRole(user.Role),
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
