package repository

import (
	"context"
	serviceModel "github.com/mistandok/auth/internal/model"
)

// UserRepository interface for crud user repositories
type UserRepository interface {
	Create(context.Context, *serviceModel.UserForCreate) (serviceModel.UserID, error)
	Update(context.Context, *serviceModel.UserForUpdate) error
	Get(context.Context, serviceModel.UserID) (*serviceModel.User, error)
	Delete(context.Context, serviceModel.UserID) error
}
