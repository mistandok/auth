package repository

import (
	"context"

	serviceModel "github.com/mistandok/auth/internal/model"
)

//go:generate ../../bin/mockery --output ./mocks  --inpackage-suffix --all --case snake

// UserRepository interface for crud user repositories
type UserRepository interface {
	Create(context.Context, *serviceModel.UserForCreate) (int64, error)
	Update(context.Context, *serviceModel.UserForUpdate) error
	Get(context.Context, int64) (*serviceModel.User, error)
	Delete(context.Context, int64) error
	GetByEmail(ctx context.Context, email string) (*serviceModel.User, error)
}
