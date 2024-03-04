package service

import (
	"context"

	"github.com/mistandok/auth/internal/model"
)

// UserService ..
type UserService interface {
	Create(context.Context, *model.UserForCreate) (int64, error)
	Update(context.Context, *model.UserForUpdate) error
	Get(context.Context, int64) (*model.User, error)
	Delete(context.Context, int64) error
}
