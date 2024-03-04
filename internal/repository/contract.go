package repository

import (
	"context"

	serviceModel "github.com/mistandok/auth/internal/model"
)

// UserRepository interface for crud user repositories
type UserRepository interface {
	Create(context.Context, *serviceModel.UserForCreate) (int64, error)
	Update(context.Context, *serviceModel.UserForUpdate) error
	Get(context.Context, int64) (*serviceModel.User, error)
	Delete(context.Context, int64) error
}
