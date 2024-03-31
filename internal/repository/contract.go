package repository

import (
	"context"
	"time"

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

// WhiteListRepository interface for work with white list for JWT tokens
type WhiteListRepository interface {
	Set(ctx context.Context, userID int64, jwtString string, expireIn time.Duration) error
	Get(ctx context.Context, userID int64) (string, error)
}

// EndpointAccessRepository interface для работы с доступами к эндпоинтам
type EndpointAccessRepository interface {
	Create(context.Context, *serviceModel.EndpointAccess) (int64, error)
	GetByAddressAndRole(ctx context.Context, address string, role serviceModel.UserRole) (*serviceModel.EndpointAccess, error)
}
