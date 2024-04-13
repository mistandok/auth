package convert

import (
	serviceModel "github.com/mistandok/auth/internal/model"
	repoModel "github.com/mistandok/auth/internal/repository/user/model"
	"github.com/mistandok/auth/internal/utils"
)

// ToRepoUserForUpdateFromServiceUserForUpdate ..
func ToRepoUserForUpdateFromServiceUserForUpdate(user *serviceModel.UserForUpdate) *repoModel.UserForUpdate {
	var (
		name  *string
		email *string
		role  *string
	)

	if user.Name != nil {
		name = utils.Pointer[string](*user.Name)
	}

	if user.Email != nil {
		email = utils.Pointer[string](string(*user.Email))
	}

	if user.Role != nil {
		role = utils.Pointer[string](string(*user.Role))
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
