package convert

import (
	"github.com/mistandok/auth/internal/common"
	serviceModel "github.com/mistandok/auth/internal/model"
	repoModel "github.com/mistandok/auth/internal/repository/user/model"
)

// ToRepoUserFromServiceUser ..
func ToRepoUserFromServiceUser(user *serviceModel.User) *repoModel.User {
	return &repoModel.User{
		ID:        int64(user.ID),
		Name:      string(user.Name),
		Email:     string(user.Email),
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToRepoUserForCreateFromServiceUserForCreate ..
func ToRepoUserForCreateFromServiceUserForCreate(user *serviceModel.UserForCreate) *repoModel.UserForCreate {
	return &repoModel.UserForCreate{
		Name:     string(user.Name),
		Email:    string(user.Email),
		Password: string(user.Password),
		Role:     string(user.Role),
	}
}

// ToRepoUserForUpdateFromServiceUserForUpdate ..
func ToRepoUserForUpdateFromServiceUserForUpdate(user *serviceModel.UserForUpdate) *repoModel.UserForUpdate {
	var (
		name  *string
		email *string
		role  *string
	)

	if user.Name == nil {
		name = nil
	} else {
		name = common.Pointer[string](string(*user.Name))
	}

	if user.Email == nil {
		email = nil
	} else {
		email = common.Pointer[string](string(*user.Email))
	}

	if user.Role == nil {
		role = nil
	} else {
		role = common.Pointer[string](string(*user.Role))
	}

	return &repoModel.UserForUpdate{
		ID:    int64(user.ID),
		Name:  name,
		Email: email,
		Role:  role,
	}
}

// ToServiceUserFromRepoUser ..
func ToServiceUserFromRepoUser(user *repoModel.User) *serviceModel.User {
	return &serviceModel.User{
		ID:        serviceModel.UserID(user.ID),
		Name:      serviceModel.UserName(user.Name),
		Email:     serviceModel.UserEmail(user.Email),
		Role:      serviceModel.UserRole(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
