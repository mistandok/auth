package service

import (
	"context"
	"github.com/mistandok/auth/internal/model"
)

// UserService ..
type UserService interface {
	Create(context.Context, *model.UserForCreate) (model.UserID, error)
	Update(context.Context, *model.UserForUpdate) error
	Get(context.Context, model.UserID) (*model.User, error)
	Delete(context.Context, model.UserID) error
}
